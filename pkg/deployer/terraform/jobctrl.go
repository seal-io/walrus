package terraform

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

	revisionbus "github.com/seal-io/seal/pkg/bus/servicerevision"
	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/pkg/dao/types/status"
	opk8s "github.com/seal-io/seal/pkg/operator/k8s"
	"github.com/seal-io/seal/pkg/operator/k8s/kube"
	"github.com/seal-io/seal/utils/log"
)

const (
	JobTypeApply   = "apply"
	JobTypeDestroy = "destroy"
)

type JobCreateOptions struct {
	// Type is the deployment type of job, apply or destroy or other.
	Type              string
	ServiceRevisionID string
	Image             string
	Env               []corev1.EnvVar
}

type StreamJobLogsOptions struct {
	Cli        *coreclient.CoreV1Client
	RevisionID oid.ID
	JobType    string
	Out        io.Writer
}

type JobReconciler struct {
	Logger      logr.Logger
	Kubeconfig  *rest.Config
	KubeClient  client.Client
	ModelClient *model.Client
}

const (
	// _podName the name of the pod.
	_podName = "seal-system"

	// _serviceRevisionIDLabel pod template label key for service revision id.
	_serviceRevisionIDLabel = "seal.io/service-revision-id"
	// _jobNameFormat the format of job name.
	_jobNameFormat = "tf-job-%s-%s"
	// _jobSecretPrefix the prefix of secret name.
	_jobSecretPrefix = "tf-secret-"
	// _secretMountPath the path to mount the secret.
	_secretMountPath = "/var/terraform/secrets"
	// _workdir the working directory of the job.
	_workdir = "/var/terraform/workspace"
)

const (
	// _applyCommands the commands to apply deployment of the service.
	_applyCommands = "terraform init -no-color && terraform apply -auto-approve -no-color "
	// _destroyCommands the commands to destroy deployment of the service.
	_destroyCommands = "terraform init -no-color && terraform destroy -auto-approve -no-color "
)

func (r JobReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	job := &batchv1.Job{}

	err := r.KubeClient.Get(ctx, req.NamespacedName, job)
	if err != nil {
		if kerrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, err
	}

	err = r.syncServiceRevisionStatus(ctx, job)
	if err != nil && !model.IsNotFound(err) {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r JobReconciler) Setup(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&batchv1.Job{}).
		Complete(r)
}

// syncServiceRevisionStatus sync the service revision status.
func (r JobReconciler) syncServiceRevisionStatus(ctx context.Context, job *batchv1.Job) (err error) {
	svcRevisionID, ok := job.Labels[_serviceRevisionIDLabel]
	if !ok {
		// Not a deployer job.
		return nil
	}

	svcRevision, err := r.ModelClient.ServiceRevisions().Get(ctx, oid.ID(svcRevisionID))
	if err != nil {
		return err
	}
	// If the service revision status is not running, then skip it.
	if svcRevision.Status != status.ServiceRevisionStatusRunning {
		return nil
	}

	if job.Status.Succeeded == 0 && job.Status.Failed == 0 {
		return nil
	}

	revisionStatus := status.ServiceRevisionStatusSucceeded
	// Get job pods logs.
	revisionStatusMessage, rerr := r.getJobPodsLogs(ctx, job.Name)
	if rerr != nil {
		r.Logger.Error(rerr, "failed to get job pod logs", "service-revision", svcRevisionID)
		revisionStatusMessage = rerr.Error()
	}

	if job.Status.Succeeded > 0 {
		r.Logger.Info("succeed", "service-revision", svcRevisionID)
	}

	if job.Status.Failed > 0 {
		r.Logger.Info("failed", "service-revision", svcRevisionID)
		revisionStatus = status.ServiceRevisionStatusFailed
	}

	// Report to service revision.
	svcRevision.Status = revisionStatus
	svcRevision.StatusMessage = revisionStatusMessage
	svcRevision.Duration = int(time.Since(*svcRevision.CreateTime).Seconds())

	update, err := dao.ServiceRevisionUpdate(r.ModelClient, svcRevision)
	if err != nil {
		return err
	}

	svcRevision, err = update.Save(ctx)
	if err != nil {
		return err
	}

	return revisionbus.Notify(ctx, r.ModelClient, svcRevision)
}

// getJobPodsLogs returns the logs of all pods of a job.
func (r JobReconciler) getJobPodsLogs(ctx context.Context, jobName string) (string, error) {
	clientSet, err := kubernetes.NewForConfig(r.Kubeconfig)
	if err != nil {
		return "", err
	}
	ls := "job-name=" + jobName

	pods, err := clientSet.CoreV1().Pods(types.SealSystemNamespace).
		List(ctx, metav1.ListOptions{LabelSelector: ls})
	if err != nil {
		return "", err
	}

	var logs string

	for _, pod := range pods.Items {
		var podLogs []byte

		podLogs, err = clientSet.CoreV1().Pods(types.SealSystemNamespace).
			GetLogs(pod.Name, &corev1.PodLogOptions{}).
			DoRaw(ctx)
		if err != nil {
			return "", err
		}
		logs += string(podLogs)
	}

	return logs, nil
}

// CreateJob create a job to run terraform deployment.
func CreateJob(ctx context.Context, clientSet *kubernetes.Clientset, opts JobCreateOptions) error {
	var (
		logger = log.WithName("deployer").WithName("tf")

		backoffLimit            int32 = 0
		ttlSecondsAfterFinished int32 = 3600
		name                          = getK8sJobName(_jobNameFormat, opts.Type, opts.ServiceRevisionID)
		configName                    = _jobSecretPrefix + opts.ServiceRevisionID
	)

	secret, err := clientSet.CoreV1().Secrets(types.SealSystemNamespace).Get(ctx, configName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	podTemplate := getPodTemplate(opts.ServiceRevisionID, configName, opts)
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

	job, err = clientSet.BatchV1().Jobs(types.SealSystemNamespace).Create(ctx, job, metav1.CreateOptions{})
	if err != nil {
		if kerrors.IsAlreadyExists(err) {
			logger.Warnf("k8s job %s already exists", name)
		} else {
			return err
		}
	}

	// Set ownerReferences to secret with the job name.
	secret.ObjectMeta = metav1.ObjectMeta{
		Name: configName,
		OwnerReferences: []metav1.OwnerReference{
			{
				APIVersion: batchv1.SchemeGroupVersion.String(),
				Kind:       "Job",
				Name:       name,
				UID:        job.UID,
			},
		},
	}

	_, err = clientSet.CoreV1().Secrets(types.SealSystemNamespace).Update(ctx, secret, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	logger.Debugf("k8s job %s created", name)

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
func getPodTemplate(serviceRevisionID, configName string, opts JobCreateOptions) corev1.PodTemplateSpec {
	var (
		command       = []string{"/bin/sh", "-c"}
		deployCommand = fmt.Sprintf("cp %s/main.tf main.tf && ", _secretMountPath)
		varfile       = fmt.Sprintf(" -var-file=%s/terraform.tfvars", _secretMountPath)
	)

	switch opts.Type {
	case JobTypeApply:
		deployCommand += _applyCommands + varfile
	case JobTypeDestroy:
		deployCommand += _destroyCommands + varfile
	}

	command = append(command, deployCommand)

	return corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Name: _podName,
			Labels: map[string]string{
				_serviceRevisionIDLabel: serviceRevisionID,
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
					Env: opts.Env,
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

// getK8sJobName returns the kubernetes job name for the given service revision id.
func getK8sJobName(format, jobType, serviceRevisionID string) string {
	return fmt.Sprintf(format, jobType, serviceRevisionID)
}

// StreamJobLogs streams the logs of a job.
func StreamJobLogs(ctx context.Context, opts StreamJobLogsOptions) error {
	var (
		jobName       = getK8sJobName(_jobNameFormat, opts.JobType, opts.RevisionID.String())
		labelSelector = "job-name=" + jobName
	)

	podList, err := opts.Cli.Pods(types.SealSystemNamespace).
		List(ctx, metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		return err
	}

	if len(podList.Items) == 0 {
		return nil
	}

	jobPod := podList.Items[0]

	err = wait.PollUntilContextTimeout(ctx, 1*time.Second, 1*time.Minute, true,
		func(ctx context.Context) (bool, error) {
			pod, getErr := opts.Cli.Pods(types.SealSystemNamespace).Get(ctx, jobPod.Name, metav1.GetOptions{
				ResourceVersion: "0",
			})
			if getErr != nil {
				return false, getErr
			}

			return kube.IsPodReady(pod), nil
		})
	if err != nil {
		return err
	}

	states := kube.GetContainerStates(&jobPod)
	if len(states) == 0 {
		return nil
	}

	var (
		containerName, containerType = states[0].Name, states[0].Type
		follow                       = kube.IsContainerRunning(&jobPod, kube.Container{
			Type: containerType,
			Name: containerName,
		})
		podLogOpts = &corev1.PodLogOptions{
			Container: containerName,
			Follow:    follow,
		}
	)

	return opk8s.GetPodLogs(ctx, opts.Cli, types.SealSystemNamespace, jobPod.Name, opts.Out, podLogOpts)
}
