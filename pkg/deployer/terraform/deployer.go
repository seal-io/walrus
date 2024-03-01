package terraform

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	deptypes "github.com/seal-io/walrus/pkg/deployer/types"
	runconfig "github.com/seal-io/walrus/pkg/resourceruns/config"
	runstatus "github.com/seal-io/walrus/pkg/resourceruns/status"
	"github.com/seal-io/walrus/pkg/settings"
	"github.com/seal-io/walrus/utils/log"
)

// Deployer terraform deployer to deploy the resource.
type Deployer struct {
	logger          log.Logger
	clientSet       *kubernetes.Clientset
	runConfigurator runconfig.Configurator
}

func NewDeployer(_ context.Context, opts deptypes.CreateOptions) (deptypes.Deployer, error) {
	clientSet, err := kubernetes.NewForConfig(opts.KubeConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client set: %w", err)
	}

	return &Deployer{
		clientSet:       clientSet,
		logger:          log.WithName("deployer").WithName("tf"),
		runConfigurator: runconfig.NewConfigurator(types.DeployerTypeTF),
	}, nil
}

func (d Deployer) Type() deptypes.Type {
	return types.DeployerTypeTF
}

// Apply creates a new resource run by the given resource,
// and drives the Kubernetes Job to create components of the resource.
func (d Deployer) Apply(
	ctx context.Context,
	mc model.ClientSet,
	run *model.ResourceRun,
	opts deptypes.ApplyOptions,
) (err error) {
	defer d.errorHandle(mc, run, err)

	if !status.ResourceRunStatusPlanned.IsTrue(run) {
		err = fmt.Errorf("resource run %s is not planned", run.ID)
		return
	}

	status.ResourceRunStatusPlanned.True(run, "")
	status.ResourceRunStatusApplied.Unknown(run, "")

	run, err = runstatus.UpdateStatus(ctx, mc, run)
	if err != nil {
		return
	}

	err = d.createK8sJob(ctx, mc, createK8sJobOptions{
		Type:        types.RunTaskTypeApply,
		ResourceRun: run,
	})

	return
}

func (d Deployer) Plan(
	ctx context.Context,
	mc model.ClientSet,
	run *model.ResourceRun,
	opts deptypes.PlanOptions,
) (err error) {
	defer d.errorHandle(mc, run, err)

	if !status.ResourceRunStatusPending.IsUnknown(run) {
		err = fmt.Errorf("resource run %s is not pending", run.ID)
		return
	}

	status.ResourceRunStatusPending.True(run, "")
	status.ResourceRunStatusPlanned.Unknown(run, "")

	run, err = runstatus.UpdateStatus(ctx, mc, run)
	if err != nil {
		return
	}

	err = d.createK8sJob(ctx, mc, createK8sJobOptions{
		Type:        types.RunTaskTypePlan,
		ResourceRun: run,
	})

	return err
}

// Destroy creates a new resource run by the given resource,
// and drives the Kubernetes Job to clean the components of the resource.
func (d Deployer) Destroy(
	ctx context.Context,
	mc model.ClientSet,
	run *model.ResourceRun,
	opts deptypes.DestroyOptions,
) (err error) {
	defer d.errorHandle(mc, run, err)

	if !status.ResourceRunStatusPlanned.IsTrue(run) {
		err = fmt.Errorf("resource run %s is not planned", run.ID)
		return
	}

	status.ResourceRunStatusPlanned.True(run, "")
	status.ResourceRunStatusApplied.Unknown(run, "")

	run, err = runstatus.UpdateStatus(ctx, mc, run)
	if err != nil {
		return
	}

	err = d.createK8sJob(ctx, mc, createK8sJobOptions{
		Type:        types.RunTaskTypeDestroy,
		ResourceRun: run,
	})

	return
}

// errorHandle handles the error of the deployer operation.
func (d Deployer) errorHandle(mc model.ClientSet, run *model.ResourceRun, err error) {
	if err == nil {
		return
	}

	// Update a failure status.
	runstatus.SetStatusFalse(run, err.Error())

	// Report to resource run.
	_, updateErr := runstatus.UpdateStatus(context.Background(), mc, run)
	if updateErr != nil {
		d.logger.Errorf("failed to update the status of the resource run: %v", updateErr)
	}
}

type createK8sJobOptions struct {
	// Type indicates the type of the job.
	Type types.RunJobType
	// ResourceRun indicates the resource run to create the deployment job.
	ResourceRun *model.ResourceRun
}

// createK8sJob creates a k8s job to deploy, destroy or rollback the resource.
func (d Deployer) createK8sJob(ctx context.Context, mc model.ClientSet, opts createK8sJobOptions) error {
	// Prepare tfConfig for deployment.
	secretOpts, err := runconfig.GetConfigOptions(ctx, mc, opts.ResourceRun, _secretMountPath)
	if err != nil {
		return err
	}

	if err = d.createK8sSecrets(ctx, mc, secretOpts); err != nil {
		return err
	}

	jobImage, err := settings.DeployerImage.Value(ctx, mc)
	if err != nil {
		return err
	}

	jobEnv := d.getEnv(ctx, mc, opts)

	localEnvironmentMode, err := settings.LocalEnvironmentMode.Value(ctx, mc)
	if err != nil {
		return err
	}

	// Create a deployment job.
	jobOpts := JobCreateOptions{
		Type:        opts.Type,
		Image:       jobImage,
		Env:         jobEnv,
		DockerMode:  localEnvironmentMode == "docker",
		ResourceRun: opts.ResourceRun,
		ServerURL:   secretOpts.SeverULR,
		Token:       secretOpts.Token,
	}

	return CreateJob(ctx, d.clientSet, jobOpts)
}

func (d Deployer) getEnv(ctx context.Context, mc model.ClientSet, opts createK8sJobOptions) (env []corev1.EnvVar) {
	env = append(env, corev1.EnvVar{
		Name: "ACCESS_TOKEN",
		ValueFrom: &corev1.EnvVarSource{
			SecretKeyRef: &corev1.SecretKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: _jobSecretPrefix + string(opts.ResourceRun.ID),
				},
				Key: _accessTokenkey,
			},
		},
	})

	if v := settings.DeployerAllProxy.ShouldValue(ctx, mc); v != "" {
		env = append(env, corev1.EnvVar{
			Name:  "ALL_PROXY",
			Value: v,
		})
	}

	if v := settings.DeployerHttpProxy.ShouldValue(ctx, mc); v != "" {
		env = append(env, corev1.EnvVar{
			Name:  "HTTP_PROXY",
			Value: v,
		})
	}

	if v := settings.DeployerHttpsProxy.ShouldValue(ctx, mc); v != "" {
		env = append(env, corev1.EnvVar{
			Name:  "HTTPS_PROXY",
			Value: v,
		})
	}

	if v := settings.DeployerNoProxy.ShouldValue(ctx, mc); v != "" {
		env = append(env, corev1.EnvVar{
			Name:  "NO_PROXY",
			Value: v,
		})
	}

	if settings.SkipRemoteTLSVerify.ShouldValueBool(ctx, mc) {
		env = append(env, corev1.EnvVar{
			Name:  "GIT_SSL_NO_VERIFY",
			Value: "true",
		})
	}

	if v := settings.DeployerNetworkMirrorUrl.ShouldValue(ctx, mc); v != "" {
		env = append(env,
			corev1.EnvVar{
				Name:  "TF_CLI_NETWORK_MIRROR_URL",
				Value: v,
			},
			corev1.EnvVar{
				Name:  "TF_CLI_NETWORK_MIRROR_INSECURE_SKIP_VERIFY",
				Value: "true",
			})
	}

	return env
}

// createK8sSecrets creates the k8s secrets for deployment.
func (d Deployer) createK8sSecrets(ctx context.Context, mc model.ClientSet, opts *runconfig.Options) error {
	secretData := make(map[string][]byte)
	// SecretName terraform tfConfig name.
	secretName := _jobSecretPrefix + string(opts.ResourceRun.ID)

	// Prepare terraform config files bytes for deployment.
	inputConfigs, err := d.runConfigurator.LoadAll(ctx, mc, opts)
	if err != nil {
		return err
	}

	for k, v := range inputConfigs {
		secretData[k] = v
	}

	// Mount the provider configs(e.g. kubeconfig) to secret.
	providerConfigs, err := d.runConfigurator.LoadProviders(opts.Connectors)
	if err != nil {
		return err
	}

	for k, v := range providerConfigs {
		secretData[k] = v
	}

	// Mount deploy access token to secret.
	secretData[_accessTokenkey] = []byte(opts.Token)

	// Create deployment secret.
	if err = CreateSecret(ctx, d.clientSet, secretName, secretData); err != nil {
		return err
	}

	return nil
}
