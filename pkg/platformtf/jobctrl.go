package platformtf

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/go-logr/logr"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	coreclient "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	revisionbus "github.com/seal-io/seal/pkg/bus/applicationrevision"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/platformk8s"
	"github.com/seal-io/seal/pkg/platformk8s/kube"
	"github.com/seal-io/seal/utils/log"
)

type JobCreateOptions struct {
	// Type is the deployment type of job, apply or destroy or other.
	Type                  string
	ApplicationRevisionID string
	Image                 string
}

const (
	// _podName the name of the pod.
	_podName = "seal-system"

	// _applicationRevisionIDLabel pod template label key for application revision id.
	_applicationRevisionIDLabel = "seal.io/application-revision-id"
	// _jobNameFormat the format of job name.
	_jobNameFormat = "tf-job-%s-%s"
	// _secretPrefix the prefix of secret name.
	_secretPrefix = "tf-secret-"
	// _secretMountPath the path to mount the secret.
	_secretMountPath = "/seal/secrets"
	// _workdir the working directory of the job.
	_workdir = "/seal/deployment"
)

const (
	_jobTypeApply   = "apply"
	_jobTypeDestroy = "destroy"
)

const (
	// _applyCommands the commands to apply deployment of the application.
	_applyCommands = "terraform init -no-color && terraform apply -auto-approve -no-color"
	// _destroyCommands the commands to destroy deployment of the application.
	_destroyCommands = "terraform init -no-color && terraform destroy -auto-approve -no-color"
)

type JobReconciler struct {
	Logger      logr.Logger
	Kubeconfig  *rest.Config
	KubeClient  client.Client
	ModelClient *model.Client
}

func (r JobReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	job := &batchv1.Job{}
	err := r.KubeClient.Get(ctx, req.NamespacedName, job)
	if err != nil {
		if kerrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if err = r.syncApplicationRevisionStatus(ctx, job); err != nil {
		// ignores error, since they can't be fixed by an immediate requeue.
		var ignore = model.IsNotFound(err)
		return ctrl.Result{Requeue: !ignore}, err
	}

	return ctrl.Result{}, nil
}

func (r JobReconciler) Setup(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&batchv1.Job{}).
		Complete(r)
}

// syncApplicationRevisionStatus sync the application revision status.
func (r JobReconciler) syncApplicationRevisionStatus(ctx context.Context, job *batchv1.Job) error {
	var (
		logger         = log.WithName("platformtf").WithName("jobctrl")
		revisionStatus = status.ApplicationRevisionStatusSucceeded
	)
	appRevisionID, ok := job.Labels[_applicationRevisionIDLabel]
	if !ok {
		// not a deployer job
		return nil
	}
	appRevision, err := r.ModelClient.ApplicationRevisions().Get(ctx, types.ID(appRevisionID))
	if err != nil {
		return err
	}
	// if the application revision status is not running, then skip it.
	if appRevision.Status != status.ApplicationRevisionStatusRunning {
		return nil
	}

	if job.Status.Succeeded == 0 && job.Status.Failed == 0 {
		return nil
	}

	// get job pods logs.
	revisionStatusMessage, err := r.getJobPodsLogs(ctx, job.Name)
	if err != nil {
		return err
	}

	if job.Status.Succeeded > 0 {
		logger.Debugf("job is succeeded, job: %s", job.Name)
	}
	if job.Status.Failed > 0 {
		revisionStatus = status.ApplicationRevisionStatusFailed
		logger.Debugf("job is failed, job: %s", job.Name)
	}

	// report to application revision.
	duration := time.Since(*appRevision.CreateTime).Seconds()
	appRevision, err = r.ModelClient.ApplicationRevisions().UpdateOneID(types.ID(appRevisionID)).
		SetStatus(revisionStatus).
		SetStatusMessage(revisionStatusMessage).
		SetDuration(int(duration)).
		Save(ctx)
	if err != nil {
		return err
	}

	if err = revisionbus.Notify(ctx, r.ModelClient, appRevision); err != nil {
		return err
	}

	// if the job is complete, then delete the secret.
	return r.deleteSecret(ctx, job.Name)
}

// getJobPodsLogs returns the logs of all pods of a job.
func (r JobReconciler) getJobPodsLogs(ctx context.Context, jobName string) (string, error) {
	clientSet, err := kubernetes.NewForConfig(r.Kubeconfig)
	if err != nil {
		return "", err
	}
	var ls = "job-name=" + jobName
	pods, err := clientSet.CoreV1().Pods(types.SealSystemNamespace).
		List(ctx, metav1.ListOptions{LabelSelector: ls})
	if err != nil {
		return "", err
	}

	var logs string
	for _, pod := range pods.Items {
		var podLogs []byte
		podLogs, err = clientSet.CoreV1().Pods(types.SealSystemNamespace).GetLogs(pod.Name, &corev1.PodLogOptions{}).DoRaw(ctx)
		if err != nil {
			return "", err
		}
		logs += string(podLogs)
	}

	return logs, nil
}

func (r JobReconciler) deleteSecret(ctx context.Context, secretName string) error {
	clientSet, err := kubernetes.NewForConfig(r.Kubeconfig)
	if err != nil {
		return err
	}
	err = clientSet.CoreV1().Secrets(types.SealSystemNamespace).Delete(ctx, secretName, metav1.DeleteOptions{})
	if err != nil && !kerrors.IsNotFound(err) {
		return err
	}

	return nil
}

// CreateJob create a job to run terraform deployment.
func CreateJob(ctx context.Context, clientSet *kubernetes.Clientset, opts JobCreateOptions) error {
	var (
		logger = log.WithName("platformtf").WithName("jobctrl")

		backoffLimit            int32 = 0
		ttlSecondsAfterFinished int32 = 60
		name                          = getK8sJobName(_jobNameFormat, opts.Type, opts.ApplicationRevisionID)
		configName                    = _secretPrefix + opts.ApplicationRevisionID
	)
	podTemplate := getPodTemplate(opts.ApplicationRevisionID, configName, opts)
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: batchv1.JobSpec{
			Template:                podTemplate,
			BackoffLimit:            &backoffLimit,
			TTLSecondsAfterFinished: &ttlSecondsAfterFinished,
		},
	}

	_, err := clientSet.BatchV1().Jobs(types.SealSystemNamespace).Create(ctx, job, metav1.CreateOptions{})
	if err != nil {
		if kerrors.IsAlreadyExists(err) {
			logger.Warnf("job %s already exists", name)
		} else {
			return err
		}
	}
	logger.Debugf("job %s created", name)

	return nil
}

// CreateSecret create a secret to store terraform config.
func CreateSecret(ctx context.Context, clientSet *kubernetes.Clientset, name string, data map[string][]byte) error {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Data:       data,
	}
	_, err := clientSet.CoreV1().Secrets(types.SealSystemNamespace).Create(ctx, secret, metav1.CreateOptions{})
	if err != nil && !kerrors.IsAlreadyExists(err) {
		return err
	}

	return nil
}

// getPodTemplate returns a pod template for deployment.
func getPodTemplate(applicationRevisionID, configName string, opts JobCreateOptions) corev1.PodTemplateSpec {
	var (
		command       = []string{"/bin/sh", "-c"}
		deployCommand = fmt.Sprintf("cp %s/main.tf main.tf && ", _secretMountPath)
		varfile       = fmt.Sprintf(" -var-file=%s/terraform.tfvars", _secretMountPath)
	)
	switch opts.Type {
	case _jobTypeApply:
		deployCommand += _applyCommands + varfile
	case _jobTypeDestroy:
		deployCommand += _destroyCommands + varfile
	}
	command = append(command, deployCommand)

	return corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Name: _podName,
			Labels: map[string]string{
				_applicationRevisionIDLabel: applicationRevisionID,
			},
		},
		Spec: corev1.PodSpec{
			HostNetwork:        true,
			ServiceAccountName: types.DeployerServiceAccountName,
			RestartPolicy:      corev1.RestartPolicyNever,
			Containers: []corev1.Container{
				{
					Name:            "deployment",
					Image:           opts.Image,
					WorkingDir:      _workdir,
					Command:         command,
					ImagePullPolicy: corev1.PullIfNotPresent,
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      configName,
							MountPath: _secretMountPath,
							ReadOnly:  false,
						},
					},
				},
			},
			Volumes: []corev1.Volume{
				{
					Name: configName,
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: configName,
						},
					},
				},
			},
		},
	}
}

// getK8sJobName returns the kubernetes job name for the given application revision id.
func getK8sJobName(format, jobType, applicationRevisionID string) string {
	return fmt.Sprintf(_jobNameFormat, jobType, applicationRevisionID)
}

// StreamJobLogs streams the logs of a job.
func StreamJobLogs(ctx context.Context, cli *coreclient.CoreV1Client, revisionID types.ID, out io.Writer) error {
	var (
		jobName       = getK8sJobName(_jobNameFormat, _jobTypeApply, revisionID.String())
		labelSelector = "job-name=" + jobName
	)

	podList, err := cli.Pods(types.SealSystemNamespace).
		List(ctx, metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		return err
	}
	if len(podList.Items) == 0 {
		return nil
	}

	var jobPod = podList.Items[0]
	err = wait.PollImmediate(1*time.Second, 1*time.Minute, func() (bool, error) {
		pod, getErr := cli.Pods(types.SealSystemNamespace).Get(ctx, jobPod.Name, metav1.GetOptions{ResourceVersion: "0"})
		if getErr != nil {
			return false, getErr
		}
		return kube.IsPodReady(pod), nil
	})
	if err != nil {
		return err
	}
	var states = kube.GetContainerStates(&jobPod)
	if len(states) == 0 {
		return nil
	}

	var (
		containerName, containerType = states[0].Name, states[0].Type
		follow                       = kube.IsContainerRunning(&jobPod, kube.Container{Type: containerType, Name: containerName})
		podLogOpts                   = &corev1.PodLogOptions{
			Container: containerName,
			Follow:    follow,
		}
	)

	return platformk8s.GetPodLogs(ctx, cli, types.SealSystemNamespace, jobPod.Name, out, podLogOpts)
}
