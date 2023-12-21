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
	"sigs.k8s.io/controller-runtime/pkg/client"

	apiconfig "github.com/seal-io/walrus/pkg/apis/config"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	opk8s "github.com/seal-io/walrus/pkg/operator/k8s"
	"github.com/seal-io/walrus/pkg/operator/k8s/kube"
	"github.com/seal-io/walrus/utils/log"
)

type JobCreateOptions struct {
	// Type is the deployment type of job, apply or destroy or other.
	Type  string
	Image string
	Env   []corev1.EnvVar

	Token     string
	ServerURL string
}

type StreamJobLogsOptions struct {
	Cli        *coreclient.CoreV1Client
	RevisionID object.ID
	JobType    string
	Out        io.Writer
}

type JobReconciler struct {
	Logger      logr.Logger
	Kubeconfig  *rest.Config
	KubeClient  client.Client
	ModelClient *model.Client
}

// createJob create a job to run terraform deployment.
func createJob(
	ctx context.Context,
	clientSet *kubernetes.Clientset,
	revision *model.ResourceRevision,
	opts JobCreateOptions,
) error {
	var (
		logger = log.WithName("deployer").WithName("tf")

		backoffLimit            int32 = 0
		ttlSecondsAfterFinished int32 = 900
		revisionID                    = revision.ID.String()
		name                          = getK8sJobName(_jobNameFormat, opts.Type, revisionID)
		configName                    = _jobSecretPrefix + revisionID
		labels                        = map[string]string{
			_resourceRevisionIDLabel: revisionID,
		}
	)

	if revision.Type == types.ResourceRevisionTypeDetect {
		labels[types.LabelWalrusDriftDetection] = "true"
	}

	secret, err := clientSet.CoreV1().Secrets(types.WalrusSystemNamespace).Get(ctx, configName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	podTemplate := getPodTemplate(revision, configName, opts)
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
		Spec: batchv1.JobSpec{
			Template:                podTemplate,
			BackoffLimit:            &backoffLimit,
			TTLSecondsAfterFinished: &ttlSecondsAfterFinished,
		},
	}

	job, err = clientSet.BatchV1().Jobs(types.WalrusSystemNamespace).Create(ctx, job, metav1.CreateOptions{})
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

	_, err = clientSet.CoreV1().Secrets(types.WalrusSystemNamespace).Update(ctx, secret, metav1.UpdateOptions{})
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

	_, err := clientSet.CoreV1().Secrets(types.WalrusSystemNamespace).Create(ctx, secret, metav1.CreateOptions{})
	if err != nil && !kerrors.IsAlreadyExists(err) {
		return err
	}

	return nil
}

// getPodTemplate returns a pod template for deployment.
func getPodTemplate(revision *model.ResourceRevision, configName string, opts JobCreateOptions) corev1.PodTemplateSpec {
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
	case JobTypeSync:
		deployCommand += _syncCommands + varfile
	case JobTypeDetect:
		deployCommand += fmt.Sprintf(_detectCommands, varfile)

		driftAPI := fmt.Sprintf("%s%s", opts.ServerURL,
			fmt.Sprintf(_driftAPI,
				revision.ProjectID,
				revision.EnvironmentID,
				revision.ResourceID,
				revision.ID))

		deployCommand += fmt.Sprintf(
			" && curl -s -f -X PUT -H \"Content-Type: application/json\" -H \"Authorization: Bearer %s\" %s -d @%s",
			opts.Token,
			driftAPI,
			_driftFile,
		)

		if !apiconfig.TlsCertified.Get() {
			deployCommand += " -k"
		}
	}

	command = append(command, deployCommand)

	return corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Name: _podName,
			Labels: map[string]string{
				_resourceRevisionIDLabel: revision.ID.String(),
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

// getK8sJobName returns the kubernetes job name for the given application revision id.
func getK8sJobName(format, jobType, applicationRevisionID string) string {
	return fmt.Sprintf(format, jobType, applicationRevisionID)
}

// StreamJobLogs streams the logs of a job.
func StreamJobLogs(ctx context.Context, opts StreamJobLogsOptions) error {
	var (
		jobName       = getK8sJobName(_jobNameFormat, opts.JobType, opts.RevisionID.String())
		labelSelector = "job-name=" + jobName
	)

	podList, err := opts.Cli.Pods(types.WalrusSystemNamespace).
		List(ctx, metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		return err
	}

	if len(podList.Items) == 0 {
		return nil
	}

	jobPod := &podList.Items[0]

	err = wait.PollUntilContextTimeout(ctx, 1*time.Second, 1*time.Minute, true,
		func(ctx context.Context) (bool, error) {
			var getErr error

			jobPod, getErr = opts.Cli.Pods(types.WalrusSystemNamespace).Get(ctx, jobPod.Name, metav1.GetOptions{
				ResourceVersion: "0",
			})
			if getErr != nil {
				return false, getErr
			}

			return kube.IsPodReady(jobPod), nil
		})
	if err != nil {
		return err
	}

	states := kube.GetContainerStates(jobPod)
	if len(states) == 0 {
		return nil
	}

	var (
		containerName, containerType = states[0].Name, states[0].Type
		follow                       = kube.IsContainerRunning(jobPod, kube.Container{
			Type: containerType,
			Name: containerName,
		})
		podLogOpts = &corev1.PodLogOptions{
			Container: containerName,
			Follow:    follow,
		}
	)

	return opk8s.GetPodLogs(ctx, opts.Cli, types.WalrusSystemNamespace, jobPod.Name, opts.Out, podLogOpts)
}
