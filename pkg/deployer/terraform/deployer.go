package terraform

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	revisionbus "github.com/seal-io/walrus/pkg/bus/resourcerevision"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	deptypes "github.com/seal-io/walrus/pkg/deployer/types"
	pkgrevision "github.com/seal-io/walrus/pkg/resourcerevision"
	"github.com/seal-io/walrus/pkg/settings"
	"github.com/seal-io/walrus/utils/log"
)

// DeployerType the type of deployer.
const DeployerType = types.DeployerTypeTF

// Deployer terraform deployer to deploy the resource.
type Deployer struct {
	logger log.Logger

	clientSet       *kubernetes.Clientset
	revisionManager *pkgrevision.Manager
}

func NewDeployer(_ context.Context, opts deptypes.CreateOptions) (deptypes.Deployer, error) {
	clientSet, err := kubernetes.NewForConfig(opts.KubeConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client set: %w", err)
	}

	return &Deployer{
		logger:          log.WithName("deployer").WithName("tf"),
		clientSet:       clientSet,
		revisionManager: pkgrevision.NewManager(),
	}, nil
}

func (d Deployer) Type() deptypes.Type {
	return DeployerType
}

// Apply creates a new resource revision by the given resource,
// and drives the Kubernetes Job to create components of the resource.
func (d Deployer) Apply(
	ctx context.Context,
	mc model.ClientSet,
	resource *model.Resource,
	opts deptypes.ApplyOptions,
) (err error) {
	revision, err := d.revisionManager.Create(ctx, mc, pkgrevision.CreateOptions{
		ResourceID: resource.ID,
		JobType:    JobTypeApply,
	})
	if err != nil {
		return err
	}

	defer func() {
		if err == nil {
			return
		}

		// Update a failure status.
		status.ResourceRevisionStatusReady.False(revision, err.Error())

		// Report to resource revision.
		_ = d.updateRevisionStatus(ctx, mc, revision)
	}()

	return d.createK8sJob(ctx, mc, createK8sJobOptions{
		Type:             JobTypeApply,
		ResourceRevision: revision,
	})
}

// Destroy creates a new resource revision by the given resource,
// and drives the Kubernetes Job to clean the components of the resource.
func (d Deployer) Destroy(
	ctx context.Context,
	mc model.ClientSet,
	resource *model.Resource,
	opts deptypes.DestroyOptions,
) (err error) {
	revision, err := d.revisionManager.Create(ctx, mc, pkgrevision.CreateOptions{
		ResourceID: resource.ID,
		JobType:    JobTypeDestroy,
	})
	if err != nil {
		return err
	}

	defer func() {
		if err == nil {
			return
		}

		// Update a failure status.
		status.ResourceRevisionStatusReady.False(revision, err.Error())

		// Report to resource revision.
		_ = d.updateRevisionStatus(ctx, mc, revision)
		_ = d.updateRevisionStatus(ctx, mc, revision)
	}()

	return d.createK8sJob(ctx, mc, createK8sJobOptions{
		Type:             JobTypeDestroy,
		ResourceRevision: revision,
	})
}

type createK8sJobOptions struct {
	// Type indicates the type of the job.
	Type string
	// ResourceRevision indicates the resource revision to create the deployment job.
	ResourceRevision *model.ResourceRevision
}

// createK8sJob creates a k8s job to deploy, destroy or rollback the resource.
func (d Deployer) createK8sJob(ctx context.Context, mc model.ClientSet, opts createK8sJobOptions) error {
	// Prepare tfConfig for deployment.
	secretOpts, err := pkgrevision.GetPlanOptions(ctx, mc, opts.ResourceRevision, SecretMountPath)
	if err != nil {
		return err
	}

	if err := d.createK8sSecrets(ctx, mc, secretOpts); err != nil {
		return err
	}

	jobImage, err := settings.DeployerImage.Value(ctx, mc)
	if err != nil {
		return err
	}

	jobEnv, err := d.getEnv(ctx, mc)
	if err != nil {
		return err
	}

	// Create deployment job.
	jobOpts := JobCreateOptions{
		Type:               opts.Type,
		ResourceRevisionID: opts.ResourceRevision.ID.String(),
		Image:              jobImage,
		Env:                jobEnv,
	}

	return CreateJob(ctx, d.clientSet, jobOpts)
}

func (d Deployer) getEnv(ctx context.Context, mc model.ClientSet) ([]corev1.EnvVar, error) {
	var env []corev1.EnvVar

	allProxy, err := settings.DeployerAllProxy.Value(ctx, mc)
	if err != nil {
		return nil, err
	}

	if allProxy != "" {
		env = append(env, corev1.EnvVar{
			Name:  "ALL_PROXY",
			Value: allProxy,
		})
	}

	httpProxy, err := settings.DeployerHttpProxy.Value(ctx, mc)
	if err != nil {
		return nil, err
	}

	if httpProxy != "" {
		env = append(env, corev1.EnvVar{
			Name:  "HTTP_PROXY",
			Value: httpProxy,
		})
	}

	httpsProxy, err := settings.DeployerHttpsProxy.Value(ctx, mc)
	if err != nil {
		return nil, err
	}

	if httpsProxy != "" {
		env = append(env, corev1.EnvVar{
			Name:  "HTTPS_PROXY",
			Value: httpsProxy,
		})
	}

	noProxy, err := settings.DeployerNoProxy.Value(ctx, mc)
	if err != nil {
		return nil, err
	}

	if noProxy != "" {
		env = append(env, corev1.EnvVar{
			Name:  "NO_PROXY",
			Value: noProxy,
		})
	}

	if settings.SkipRemoteTLSVerify.ShouldValueBool(ctx, mc) {
		env = append(env, corev1.EnvVar{
			Name:  "GIT_SSL_NO_VERIFY",
			Value: "true",
		})
	}

	return env, nil
}

func (d Deployer) updateRevisionStatus(ctx context.Context, mc model.ClientSet, ar *model.ResourceRevision) error {
	// Report to resource revision.
	ar.Status.SetSummary(status.WalkResourceRevision(&ar.Status))

	ar, err := mc.ResourceRevisions().UpdateOne(ar).
		SetStatus(ar.Status).
		Save(ctx)
	if err != nil {
		return err
	}

	if err = revisionbus.Notify(ctx, mc, ar); err != nil {
		d.logger.Error(err)
		return err
	}

	return nil
}

// createK8sSecrets creates the k8s secrets for deployment.
func (d Deployer) createK8sSecrets(ctx context.Context, mc model.ClientSet, opts *pkgrevision.PlanOptions) error {
	secretData := make(map[string][]byte)
	// SecretName terraform tfConfig name.
	secretName := _jobSecretPrefix + string(opts.ResourceRevision.ID)

	// Prepare terraform config files bytes for deployment.
	tfConfigs, err := d.revisionManager.LoadConfigs(ctx, mc, opts)
	if err != nil {
		return err
	}

	for k, v := range tfConfigs {
		secretData[k] = v
	}

	// Mount the provider configs(e.g. kubeconfig) to secret.
	providerData, err := d.revisionManager.LoadConnectorConfigs(opts.Connectors)
	if err != nil {
		return err
	}

	for k, v := range providerData {
		secretData[k] = v
	}

	// Create deployment secret.
	if err = CreateSecret(ctx, d.clientSet, secretName, secretData); err != nil {
		return err
	}

	return nil
}
