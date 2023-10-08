package terraform

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"

	"github.com/seal-io/walrus/pkg/auths"
	"github.com/seal-io/walrus/pkg/auths/session"
	revisionbus "github.com/seal-io/walrus/pkg/bus/servicerevision"
	"github.com/seal-io/walrus/pkg/dao"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/connector"
	"github.com/seal-io/walrus/pkg/dao/model/environmentconnectorrelationship"
	"github.com/seal-io/walrus/pkg/dao/model/predicate"
	"github.com/seal-io/walrus/pkg/dao/model/service"
	"github.com/seal-io/walrus/pkg/dao/model/serviceresource"
	"github.com/seal-io/walrus/pkg/dao/model/servicerevision"
	"github.com/seal-io/walrus/pkg/dao/model/templateversion"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/crypto"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/dao/types/property"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	deptypes "github.com/seal-io/walrus/pkg/deployer/types"
	pkgenv "github.com/seal-io/walrus/pkg/environment"
	opk8s "github.com/seal-io/walrus/pkg/operator/k8s"
	pkgservice "github.com/seal-io/walrus/pkg/service"
	"github.com/seal-io/walrus/pkg/settings"
	"github.com/seal-io/walrus/pkg/terraform/config"
	"github.com/seal-io/walrus/pkg/terraform/parser"
	"github.com/seal-io/walrus/pkg/terraform/util"
	"github.com/seal-io/walrus/utils/json"
	"github.com/seal-io/walrus/utils/log"
	"github.com/seal-io/walrus/utils/pointer"
)

// DeployerType the type of deployer.
const DeployerType = types.DeployerTypeTF

const (
	// _backendAPI the API path to terraform deploy backend.
	// Terraform will get and update deployment states from this API.
	_backendAPI = "/v1/projects/%s/environments/%s/services/%s/revisions/%s/terraform-states"

	// _variablePrefix the prefix of the variable name.
	_variablePrefix = "_walrus_var_"

	// _servicePrefix the prefix of the service output name.
	_servicePrefix = "_walrus_service_"
)

var (
	// _variableReg the regexp to match the variable.
	_variableReg = regexp.MustCompile(`\${var\.([a-zA-Z0-9_-]+)}`)
	// _serviceReg the regexp to match the service output.
	_serviceReg = regexp.MustCompile(`\${service\.([^.}]+)\.([^.}]+)}`)
)

// Deployer terraform deployer to deploy the service.
type Deployer struct {
	logger      log.Logger
	modelClient model.ClientSet
	clientSet   *kubernetes.Clientset
}

func NewDeployer(_ context.Context, opts deptypes.CreateOptions) (deptypes.Deployer, error) {
	clientSet, err := kubernetes.NewForConfig(opts.KubeConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client set: %w", err)
	}

	return &Deployer{
		modelClient: opts.ModelClient,
		clientSet:   clientSet,
		logger:      log.WithName("deployer").WithName("tf"),
	}, nil
}

func (d Deployer) Type() deptypes.Type {
	return DeployerType
}

// Apply creates a new service revision by the given service,
// and drives the Kubernetes Job to create resources of the service.
func (d Deployer) Apply(ctx context.Context, service *model.Service, opts deptypes.ApplyOptions) (err error) {
	revision, err := d.createRevision(ctx, createRevisionOptions{
		ServiceID: service.ID,
		JobType:   JobTypeApply,
	})
	if err != nil {
		return err
	}

	defer func() {
		if err == nil {
			return
		}

		// Update a failure status.
		status.ServiceRevisionStatusReady.False(revision, err.Error())

		// Report to service revision.
		_ = d.updateRevisionStatus(ctx, revision)
	}()

	return d.createK8sJob(ctx, createK8sJobOptions{
		Type:            JobTypeApply,
		SkipTLSVerify:   opts.SkipTLSVerify,
		ServiceRevision: revision,
	})
}

// Destroy creates a new service revision by the given service,
// and drives the Kubernetes Job to clean the resources of the service.
func (d Deployer) Destroy(ctx context.Context, service *model.Service, opts deptypes.DestroyOptions) (err error) {
	revision, err := d.createRevision(ctx, createRevisionOptions{
		ServiceID: service.ID,
		JobType:   JobTypeDestroy,
	})
	if err != nil {
		return err
	}

	defer func() {
		if err == nil {
			return
		}

		// Update a failure status.
		status.ServiceRevisionStatusReady.False(revision, err.Error())

		// Report to service revision.
		_ = d.updateRevisionStatus(ctx, revision)
	}()

	// If no resource exists, skip job and set revision status succeed.
	exist, err := d.modelClient.ServiceResources().Query().
		Where(serviceresource.ServiceID(service.ID)).
		Exist(ctx)
	if err != nil {
		return err
	}

	if !exist {
		status.ServiceRevisionStatusReady.True(revision, "")
		return d.updateRevisionStatus(ctx, revision)
	}

	return d.createK8sJob(ctx, createK8sJobOptions{
		Type:            JobTypeDestroy,
		SkipTLSVerify:   opts.SkipTLSVerify,
		ServiceRevision: revision,
	})
}

type createK8sJobOptions struct {
	// Type indicates the type of the job.
	Type string
	// SkipTLSVerify indicates to skip TLS verification.
	SkipTLSVerify bool
	// ServiceRevision indicates the service revision to create the deployment job.
	ServiceRevision *model.ServiceRevision
}

// createK8sJob creates a k8s job to deploy, destroy or rollback the service.
func (d Deployer) createK8sJob(ctx context.Context, opts createK8sJobOptions) error {
	revision := opts.ServiceRevision

	connectors, err := d.getConnectors(ctx, revision.EnvironmentID)
	if err != nil {
		return err
	}

	proj, err := d.modelClient.Projects().Get(ctx, revision.ProjectID)
	if err != nil {
		return err
	}

	env, err := dao.GetEnvironmentByID(ctx, d.modelClient, revision.EnvironmentID)
	if err != nil {
		return err
	}

	svc, err := d.modelClient.Services().Get(ctx, revision.ServiceID)
	if err != nil {
		return err
	}

	var subjectID object.ID

	sj, _ := session.GetSubject(ctx)
	if sj.ID != "" {
		subjectID = sj.ID
	} else {
		subjectID, err = pkgservice.GetSubjectID(svc)
		if err != nil {
			return err
		}
	}

	if subjectID == "" {
		return errors.New("subject id is empty")
	}

	// Prepare tfConfig for deployment.
	secretOpts := createK8sSecretsOptions{
		SkipTLSVerify:   opts.SkipTLSVerify,
		ServiceRevision: opts.ServiceRevision,
		Connectors:      connectors,
		ProjectID:       proj.ID,
		EnvironmentID:   env.ID,
		SubjectID:       subjectID,
		// Metadata.
		ProjectName:          proj.Name,
		EnvironmentName:      env.Name,
		ServiceName:          svc.Name,
		ServiceID:            svc.ID,
		ManagedNamespaceName: pkgenv.GetManagedNamespaceName(env),
	}
	if err = d.createK8sSecrets(ctx, secretOpts); err != nil {
		return err
	}

	jobImage, err := settings.DeployerImage.Value(ctx, d.modelClient)
	if err != nil {
		return err
	}

	jobEnv, err := d.getProxyEnv(ctx)
	if err != nil {
		return err
	}

	// Create deployment job.
	jobOpts := JobCreateOptions{
		Type:              opts.Type,
		ServiceRevisionID: opts.ServiceRevision.ID.String(),
		Image:             jobImage,
		Env:               jobEnv,
	}

	return CreateJob(ctx, d.clientSet, jobOpts)
}

func (d Deployer) getProxyEnv(ctx context.Context) ([]corev1.EnvVar, error) {
	var env []corev1.EnvVar

	allProxy, err := settings.DeployerAllProxy.Value(ctx, d.modelClient)
	if err != nil {
		return nil, err
	}

	if allProxy != "" {
		env = append(env, corev1.EnvVar{
			Name:  "ALL_PROXY",
			Value: allProxy,
		})
	}

	httpProxy, err := settings.DeployerHttpProxy.Value(ctx, d.modelClient)
	if err != nil {
		return nil, err
	}

	if httpProxy != "" {
		env = append(env, corev1.EnvVar{
			Name:  "HTTP_PROXY",
			Value: httpProxy,
		})
	}

	httpsProxy, err := settings.DeployerHttpsProxy.Value(ctx, d.modelClient)
	if err != nil {
		return nil, err
	}

	if httpsProxy != "" {
		env = append(env, corev1.EnvVar{
			Name:  "HTTPS_PROXY",
			Value: httpsProxy,
		})
	}

	noProxy, err := settings.DeployerNoProxy.Value(ctx, d.modelClient)
	if err != nil {
		return nil, err
	}

	if noProxy != "" {
		env = append(env, corev1.EnvVar{
			Name:  "NO_PROXY",
			Value: noProxy,
		})
	}

	return env, nil
}

func (d Deployer) updateRevisionStatus(ctx context.Context, ar *model.ServiceRevision) error {
	// Report to service revision.
	ar.Status.SetSummary(status.WalkServiceRevision(&ar.Status))

	ar, err := d.modelClient.ServiceRevisions().UpdateOne(ar).
		SetStatus(ar.Status).
		Save(ctx)
	if err != nil {
		return err
	}

	if err = revisionbus.Notify(ctx, d.modelClient, ar); err != nil {
		d.logger.Error(err)
		return err
	}

	return nil
}

type createK8sSecretsOptions struct {
	SkipTLSVerify   bool
	ServiceRevision *model.ServiceRevision
	Connectors      model.Connectors
	ProjectID       object.ID
	EnvironmentID   object.ID
	SubjectID       object.ID
	// Metadata.
	ProjectName          string
	EnvironmentName      string
	ServiceName          string
	ServiceID            object.ID
	ManagedNamespaceName string
}

// createK8sSecrets creates the k8s secrets for deployment.
func (d Deployer) createK8sSecrets(ctx context.Context, opts createK8sSecretsOptions) error {
	secretData := make(map[string][]byte)
	// SecretName terraform tfConfig name.
	secretName := _jobSecretPrefix + string(opts.ServiceRevision.ID)

	// Prepare terraform config files bytes for deployment.
	terraformData, err := d.loadConfigsBytes(ctx, opts)
	if err != nil {
		return err
	}

	for k, v := range terraformData {
		secretData[k] = v
	}

	// Mount the provider configs(e.g. kubeconfig) to secret.
	providerData, err := d.getProviderSecretData(opts.Connectors)
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

type createRevisionOptions struct {
	// ServiceID indicates the ID of service which is for create the revision.
	ServiceID object.ID
	// JobType indicates the type of the job.
	JobType string
}

// createRevision creates a new service revision.
// Get the latest revision, and check it if it is running.
// If not running, then apply the latest revision.
// If running, then wait for the latest revision to be applied.
func (d Deployer) createRevision(
	ctx context.Context,
	opts createRevisionOptions,
) (*model.ServiceRevision, error) {
	// Validate if there is a running revision.
	prevEntity, err := d.modelClient.ServiceRevisions().Query().
		Where(servicerevision.And(
			servicerevision.ServiceID(opts.ServiceID),
			servicerevision.DeployerType(DeployerType))).
		Order(model.Desc(servicerevision.FieldCreateTime)).
		First(ctx)
	if err != nil && !model.IsNotFound(err) {
		return nil, err
	}

	if prevEntity != nil && status.ServiceRevisionStatusReady.IsUnknown(prevEntity) {
		return nil, errors.New("service deployment is running")
	}

	// Get the corresponding service and template version.
	svc, err := d.modelClient.Services().Query().
		Where(service.ID(opts.ServiceID)).
		WithTemplate(func(tvq *model.TemplateVersionQuery) {
			tvq.Select(
				templateversion.FieldName,
				templateversion.FieldVersion,
				templateversion.FieldTemplateID)
		}).
		Select(
			service.FieldID,
			service.FieldProjectID,
			service.FieldEnvironmentID,
			service.FieldAttributes).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	entity := &model.ServiceRevision{
		ProjectID:       svc.ProjectID,
		EnvironmentID:   svc.EnvironmentID,
		ServiceID:       svc.ID,
		TemplateID:      svc.Edges.Template.TemplateID,
		TemplateName:    svc.Edges.Template.Name,
		TemplateVersion: svc.Edges.Template.Version,
		Attributes:      svc.Attributes,
		DeployerType:    DeployerType,
	}

	status.ServiceRevisionStatusReady.Unknown(entity, "")
	entity.Status.SetSummary(status.WalkServiceRevision(&entity.Status))

	// Inherit the output of previous revision to create a new one.
	if prevEntity != nil {
		entity.Output = prevEntity.Output
	}

	switch {
	case opts.JobType == JobTypeApply && entity.Output != "":
		// Get required providers from the previous output after first deployment.
		requiredProviders, err := d.getRequiredProviders(ctx, opts.ServiceID, entity.Output)
		if err != nil {
			return nil, err
		}
		entity.PreviousRequiredProviders = requiredProviders
	case opts.JobType == JobTypeDestroy && entity.Output != "":
		if status.ServiceRevisionStatusReady.IsFalse(prevEntity) {
			// Get required providers from the previous output after first deployment.
			requiredProviders, err := d.getRequiredProviders(ctx, opts.ServiceID, entity.Output)
			if err != nil {
				return nil, err
			}
			entity.PreviousRequiredProviders = requiredProviders
		} else {
			// Copy required providers from the previous revision.
			entity.PreviousRequiredProviders = prevEntity.PreviousRequiredProviders
			// Reuse other fields from the previous revision.
			entity.TemplateID = prevEntity.TemplateID
			entity.TemplateName = prevEntity.TemplateName
			entity.TemplateVersion = prevEntity.TemplateVersion
			entity.Attributes = prevEntity.Attributes
			entity.InputPlan = prevEntity.InputPlan
		}
	}

	// Create revision.
	entity, err = d.modelClient.ServiceRevisions().Create().
		Set(entity).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (d Deployer) getRequiredProviders(
	ctx context.Context,
	instanceID object.ID,
	previousOutput string,
) ([]types.ProviderRequirement, error) {
	stateRequiredProviderSet := sets.NewString()

	previousRequiredProviders, err := d.getPreviousRequiredProviders(ctx, instanceID)
	if err != nil {
		return nil, err
	}

	stateRequiredProviders, err := parser.ParseStateProviders(previousOutput)
	if err != nil {
		return nil, err
	}

	stateRequiredProviderSet.Insert(stateRequiredProviders...)

	requiredProviders := make([]types.ProviderRequirement, 0, len(previousRequiredProviders))

	for _, p := range previousRequiredProviders {
		if stateRequiredProviderSet.Has(p.Name) {
			requiredProviders = append(requiredProviders, p)
		}
	}

	return requiredProviders, nil
}

// loadConfigsBytes returns terraform main.tf and terraform.tfvars for deployment.
func (d Deployer) loadConfigsBytes(ctx context.Context, opts createK8sSecretsOptions) (map[string][]byte, error) {
	logger := log.WithName("deployer").WithName("tf")
	// Prepare terraform tfConfig.
	//  get module configs from service revision.
	moduleConfig, providerRequirements, err := d.getModuleConfig(ctx, opts)
	if err != nil {
		return nil, err
	}
	// Merge current and previous required providers.
	providerRequirements = append(providerRequirements,
		opts.ServiceRevision.PreviousRequiredProviders...)

	requiredProviders := make(map[string]*tfconfig.ProviderRequirement, 0)
	for _, p := range providerRequirements {
		if _, ok := requiredProviders[p.Name]; !ok {
			requiredProviders[p.Name] = p.ProviderRequirement
		} else {
			logger.Warnf("duplicate provider requirement: %s", p.Name)
		}
	}

	serviceOpts := ServiceOpts{
		ServiceRevision: opts.ServiceRevision,
		ServiceName:     opts.ServiceName,
		ProjectID:       opts.ProjectID,
		EnvironmentID:   opts.EnvironmentID,
	}
	// Parse module attributes.
	variables, dependencyOutputs, err := ParseModuleAttributes(
		ctx,
		d.modelClient,
		moduleConfig.Attributes,
		false,
		serviceOpts,
	)
	if err != nil {
		return nil, err
	}

	// Update output sensitive with variables.
	wrapVariables, err := updateOutputWithVariables(variables, moduleConfig)
	if err != nil {
		return nil, err
	}

	// Prepare address for terraform backend.
	serverAddress, err := settings.ServeUrl.Value(ctx, d.modelClient)
	if err != nil {
		return nil, err
	}

	if serverAddress == "" {
		return nil, errors.New("server address is empty")
	}
	address := fmt.Sprintf("%s%s", serverAddress,
		fmt.Sprintf(_backendAPI,
			opts.ProjectID,
			opts.EnvironmentID,
			opts.ServiceID,
			opts.ServiceRevision.ID))

	// Prepare API token for terraform backend.
	const _1Day = 60 * 60 * 24

	at, err := auths.CreateAccessToken(ctx,
		d.modelClient, opts.SubjectID, types.TokenKindDeployment, string(opts.ServiceRevision.ID), pointer.Int(_1Day))
	if err != nil {
		return nil, err
	}

	// Prepare terraform config files to be mounted to secret.
	requiredProviderNames := sets.NewString()
	for _, p := range providerRequirements {
		requiredProviderNames = requiredProviderNames.Insert(p.Name)
	}

	secretOptionMaps := map[string]config.CreateOptions{
		config.FileMain: {
			TerraformOptions: &config.TerraformOptions{
				Token:                at.AccessToken,
				Address:              address,
				SkipTLSVerify:        opts.SkipTLSVerify,
				ProviderRequirements: requiredProviders,
			},
			ProviderOptions: &config.ProviderOptions{
				RequiredProviderNames: requiredProviderNames.List(),
				Connectors:            opts.Connectors,
				SecretMonthPath:       _secretMountPath,
				ConnectorSeparator:    parser.ConnectorSeparator,
			},
			ModuleOptions: &config.ModuleOptions{
				ModuleConfigs: []*config.ModuleConfig{moduleConfig},
			},
			VariableOptions: &config.VariableOptions{
				VariablePrefix:    _variablePrefix,
				ServicePrefix:     _servicePrefix,
				Variables:         wrapVariables,
				DependencyOutputs: dependencyOutputs,
			},
			OutputOptions: moduleConfig.Outputs,
		},
		config.FileVars: getVarConfigOptions(variables, dependencyOutputs),
	}
	secretMaps := make(map[string][]byte, 0)

	for k, v := range secretOptionMaps {
		secretMaps[k], err = config.CreateConfigToBytes(v)
		if err != nil {
			return nil, err
		}
	}

	// Save input plan to service revision.
	opts.ServiceRevision.InputPlan = string(secretMaps[config.FileMain])
	// If service revision does not inherit variables from cloned revision,
	// then save the parsed variables to service revision.
	if len(opts.ServiceRevision.Variables) == 0 {
		variableMap := make(crypto.Map[string, string], len(variables))
		for _, s := range variables {
			variableMap[s.Name] = string(s.Value)
		}
		opts.ServiceRevision.Variables = variableMap
	}

	revision, err := d.modelClient.ServiceRevisions().UpdateOne(opts.ServiceRevision).
		Set(opts.ServiceRevision).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	if err = revisionbus.Notify(ctx, d.modelClient, revision); err != nil {
		return nil, err
	}

	return secretMaps, nil
}

// getProviderSecretData returns provider kubeconfig secret data mount into terraform container.
func (d Deployer) getProviderSecretData(connectors model.Connectors) (map[string][]byte, error) {
	secretData := make(map[string][]byte)

	for _, c := range connectors {
		if c.Type != types.ConnectorTypeK8s {
			continue
		}

		_, s, err := opk8s.LoadApiConfig(*c)
		if err != nil {
			return nil, err
		}

		// NB(alex) the secret file name must be config + connector id to
		// match with terraform provider in config convert.
		secretFileName := util.GetK8sSecretName(c.ID.String())
		secretData[secretFileName] = []byte(s)
	}

	return secretData, nil
}

// getModuleConfig returns module configs and required connectors to
// get terraform module config block from service revision.
func (d Deployer) getModuleConfig(
	ctx context.Context,
	opts createK8sSecretsOptions,
) (*config.ModuleConfig, []types.ProviderRequirement, error) {
	var (
		requiredProviders = make([]types.ProviderRequirement, 0)
		predicates        = make([]predicate.TemplateVersion, 0)
	)

	predicates = append(predicates, templateversion.And(
		templateversion.Version(opts.ServiceRevision.TemplateVersion),
		templateversion.TemplateID(opts.ServiceRevision.TemplateID),
	))

	templateVersion, err := d.modelClient.TemplateVersions().
		Query().
		Select(
			templateversion.FieldID,
			templateversion.FieldTemplateID,
			templateversion.FieldName,
			templateversion.FieldVersion,
			templateversion.FieldSource,
			templateversion.FieldSchema,
		).
		Where(templateversion.Or(predicates...)).
		Only(ctx)
	if err != nil {
		return nil, nil, err
	}

	if templateVersion.Schema != nil {
		requiredProviders = append(requiredProviders, templateVersion.Schema.RequiredProviders...)
	}

	mc, err := getModuleConfig(opts.ServiceRevision, templateVersion, opts)
	if err != nil {
		return nil, nil, err
	}

	return mc, requiredProviders, err
}

func (d Deployer) getConnectors(ctx context.Context, environmentID object.ID) (model.Connectors, error) {
	rs, err := d.modelClient.EnvironmentConnectorRelationships().Query().
		Where(environmentconnectorrelationship.EnvironmentID(environmentID)).
		WithConnector(func(cq *model.ConnectorQuery) {
			cq.Select(
				connector.FieldID,
				connector.FieldName,
				connector.FieldType,
				connector.FieldCategory,
				connector.FieldConfigVersion,
				connector.FieldConfigData)
		}).
		All(ctx)
	if err != nil {
		return nil, err
	}

	var cs model.Connectors
	for i := range rs {
		cs = append(cs, rs[i].Edges.Connector)
	}

	return cs, nil
}

// getPreviousRequiredProviders get previous succeed revision required providers.
// NB(alex): the previous revision may be failed, the failed revision may not contain required providers of states.
func (d Deployer) getPreviousRequiredProviders(
	ctx context.Context,
	serviceID object.ID,
) ([]types.ProviderRequirement, error) {
	prevRequiredProviders := make([]types.ProviderRequirement, 0)

	entity, err := d.modelClient.ServiceRevisions().Query().
		Where(servicerevision.ServiceID(serviceID)).
		Order(model.Desc(servicerevision.FieldCreateTime)).
		First(ctx)
	if err != nil && !model.IsNotFound(err) {
		return nil, err
	}

	if entity == nil {
		return prevRequiredProviders, nil
	}

	templateVersion, err := d.modelClient.TemplateVersions().Query().
		Where(
			templateversion.TemplateID(entity.TemplateID),
			templateversion.Version(entity.TemplateVersion),
		).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	if templateVersion.Schema != nil {
		prevRequiredProviders = append(prevRequiredProviders, templateVersion.Schema.RequiredProviders...)
	}

	return prevRequiredProviders, nil
}

func getVarConfigOptions(variables model.Variables, serviceOutputs map[string]parser.OutputState) config.CreateOptions {
	varsConfigOpts := config.CreateOptions{
		Attributes: map[string]any{},
	}

	for _, v := range variables {
		varsConfigOpts.Attributes[_variablePrefix+v.Name] = v.Value
	}

	// Setup service outputs.
	for n, v := range serviceOutputs {
		varsConfigOpts.Attributes[_servicePrefix+n] = v.Value
	}

	return varsConfigOpts
}

func getModuleConfig(
	revision *model.ServiceRevision,
	template *model.TemplateVersion,
	opts createK8sSecretsOptions,
) (*config.ModuleConfig, error) {
	var (
		props              = make(property.Properties, len(revision.Attributes))
		typesWith          = revision.Attributes.TypesWith(template.Schema.Variables)
		sensitiveVariables = sets.Set[string]{}
	)

	for k, v := range revision.Attributes {
		props[k] = property.Property{
			Type:  typesWith[k],
			Value: v,
		}
	}

	attrs, err := props.TypedValues()
	if err != nil {
		return nil, err
	}

	mc := &config.ModuleConfig{
		Name:       opts.ServiceName,
		Source:     template.Source,
		Schema:     template.Schema,
		Attributes: attrs,
	}

	if template.Schema == nil {
		return mc, nil
	}

	for _, v := range template.Schema.Variables {
		// Add sensitive from schema variable.
		if v.Sensitive {
			sensitiveVariables.Insert(fmt.Sprintf(`var\.%s`, v.Name))
		}

		// Add seal metadata.
		var attrValue string

		switch v.Name {
		case WalrusMetadataProjectName:
			attrValue = opts.ProjectName
		case WalrusMetadataEnvironmentName:
			attrValue = opts.EnvironmentName
		case WalrusMetadataServiceName:
			attrValue = opts.ServiceName
		case WalrusMetadataProjectID:
			attrValue = opts.ProjectID.String()
		case WalrusMetadataEnvironmentID:
			attrValue = opts.EnvironmentID.String()
		case WalrusMetadataServiceID:
			attrValue = opts.ServiceID.String()
		case WalrusMetadataNamespaceName:
			attrValue = opts.ManagedNamespaceName
		}

		if attrValue != "" {
			mc.Attributes[v.Name] = attrValue
		}
	}

	sensitiveVariableRegex, err := matchAnyRegex(sensitiveVariables.UnsortedList())
	if err != nil {
		return nil, err
	}

	mc.Outputs = make([]config.Output, len(template.Schema.Outputs))
	for i, v := range template.Schema.Outputs {
		mc.Outputs[i].Sensitive = v.Sensitive
		mc.Outputs[i].Name = v.Name
		mc.Outputs[i].ServiceName = opts.ServiceName
		mc.Outputs[i].Value = v.Value

		if v.Sensitive {
			continue
		}

		// Update sensitive while output is from sensitive data, like secret.
		if sensitiveVariables.Len() != 0 && sensitiveVariableRegex.Match(v.Value) {
			mc.Outputs[i].Sensitive = true
		}
	}

	return mc, nil
}

func updateOutputWithVariables(variables model.Variables, moduleConfig *config.ModuleConfig) (map[string]bool, error) {
	var (
		variableOpts         = make(map[string]bool)
		encryptVariableNames = sets.NewString()
	)

	for _, s := range variables {
		variableOpts[s.Name] = s.Sensitive

		if s.Sensitive {
			encryptVariableNames.Insert(_variablePrefix + s.Name)
		}
	}

	if encryptVariableNames.Len() == 0 {
		return variableOpts, nil
	}

	reg, err := matchAnyRegex(encryptVariableNames.UnsortedList())
	if err != nil {
		return nil, err
	}

	var shouldEncryptAttr []string

	for k, v := range moduleConfig.Attributes {
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}

		matches := reg.FindAllString(string(b), -1)
		if len(matches) != 0 {
			shouldEncryptAttr = append(shouldEncryptAttr, fmt.Sprintf(`var\.%s`, k))
		}
	}

	// Outputs use encrypted variable should set to sensitive.
	for i, v := range moduleConfig.Outputs {
		if v.Sensitive {
			continue
		}

		reg, err := matchAnyRegex(shouldEncryptAttr)
		if err != nil {
			return nil, err
		}

		if reg.MatchString(string(v.Value)) {
			moduleConfig.Outputs[i].Sensitive = true
		}
	}

	return variableOpts, nil
}

func matchAnyRegex(list []string) (*regexp.Regexp, error) {
	var sb strings.Builder

	sb.WriteString("(")

	for i, v := range list {
		sb.WriteString(v)

		if i < len(list)-1 {
			sb.WriteString("|")
		}
	}

	sb.WriteString(")")

	return regexp.Compile(sb.String())
}
