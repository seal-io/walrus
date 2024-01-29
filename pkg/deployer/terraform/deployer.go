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
	pkgrun "github.com/seal-io/walrus/pkg/resourcerun"
	runconfig "github.com/seal-io/walrus/pkg/resourcerun/config"
	"github.com/seal-io/walrus/pkg/settings"
	"github.com/seal-io/walrus/utils/log"
)

// Deployer terraform deployer to deploy the resource.
type Deployer struct {
	logger          log.Logger
	clientSet       *kubernetes.Clientset
	runConfigLoader runconfig.Loader
}

func NewDeployer(_ context.Context, opts deptypes.CreateOptions) (deptypes.Deployer, error) {
	clientSet, err := kubernetes.NewForConfig(opts.KubeConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client set: %w", err)
	}

	return &Deployer{
		clientSet:       clientSet,
		logger:          log.WithName("deployer").WithName("tf"),
		runConfigLoader: runconfig.NewInputLoader(types.DeployerTypeTF),
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
	resource *model.Resource,
	opts deptypes.ApplyOptions,
) (err error) {
	run, err := pkgrun.Create(ctx, mc, pkgrun.CreateOptions{
		ResourceID:    resource.ID,
		ChangeComment: resource.ChangeComment,
		DeployerType:  types.DeployerTypeTF,
		JobType:       types.RunJobTypeApply,
	})
	if err != nil {
		return err
	}

	defer func() {
		if err == nil {
			return
		}

		// Update a failure status.
		status.ResourceRunStatusReady.False(run, err.Error())

		// Report to resource run.
		_ = pkgrun.UpdateStatus(ctx, mc, run)
	}()

	return d.createK8sJob(ctx, mc, createK8sJobOptions{
		Type:        types.RunJobTypeApply,
		ResourceRun: run,
	})
}

// Destroy creates a new resource run by the given resource,
// and drives the Kubernetes Job to clean the components of the resource.
func (d Deployer) Destroy(
	ctx context.Context,
	mc model.ClientSet,
	resource *model.Resource,
	opts deptypes.DestroyOptions,
) (err error) {
	run, err := pkgrun.Create(ctx, mc, pkgrun.CreateOptions{
		ResourceID:   resource.ID,
		DeployerType: types.DeployerTypeTF,
		JobType:      types.RunJobTypeDestroy,
	})
	if err != nil {
		return err
	}

	defer func() {
		if err == nil {
			return
		}

		// Update a failure status.
		status.ResourceRunStatusReady.False(run, err.Error())

		// Report to resource run.
		_ = pkgrun.UpdateStatus(ctx, mc, run)
	}()

	return d.createK8sJob(ctx, mc, createK8sJobOptions{
		Type:        types.RunJobTypeDestroy,
		ResourceRun: run,
	})
}

type createK8sJobOptions struct {
	// Type indicates the type of the job.
	Type string
	// ResourceRun indicates the resource run to create the deployment job.
	ResourceRun *model.ResourceRun
}

// createK8sJob creates a k8s job to deploy, destroy or rollback the resource.
func (d Deployer) createK8sJob(ctx context.Context, mc model.ClientSet, opts createK8sJobOptions) error {
	// Prepare tfConfig for deployment.
	secretOpts, err := runconfig.GetConfigLoaderOptions(ctx, mc, opts.ResourceRun, _secretMountPath)
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

	jobEnv := d.getEnv(ctx, mc)

	localEnvironmentMode, err := settings.LocalEnvironmentMode.Value(ctx, mc)
	if err != nil {
		return err
	}

	// Create a deployment job.
	jobOpts := JobCreateOptions{
		Type:          opts.Type,
		ResourceRunID: opts.ResourceRun.ID.String(),
		Image:         jobImage,
		Env:           jobEnv,
		DockerMode:    localEnvironmentMode == "docker",
	}

	return CreateJob(ctx, d.clientSet, jobOpts)
}

func (d Deployer) getEnv(ctx context.Context, mc model.ClientSet) (env []corev1.EnvVar) {
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
func (d Deployer) createK8sSecrets(ctx context.Context, mc model.ClientSet, opts *runconfig.LoaderOptions) error {
	secretData := make(map[string][]byte)
	// SecretName terraform tfConfig name.
	secretName := _jobSecretPrefix + string(opts.ResourceRun.ID)

	// Prepare terraform config files bytes for deployment.
	inputConfigs, err := d.runConfigLoader.LoadAll(ctx, mc, opts)
	if err != nil {
		return err
	}

	for k, v := range inputConfigs {
		secretData[k] = v
	}

	// Mount the provider configs(e.g. kubeconfig) to secret.
	providerConfigs, err := d.runConfigLoader.LoadProviders(opts.Connectors)
	if err != nil {
		return err
	}

	for k, v := range providerConfigs {
		secretData[k] = v
	}

	// Create deployment secret.
	if err = CreateSecret(ctx, d.clientSet, secretName, secretData); err != nil {
		return err
	}

	return nil
}
