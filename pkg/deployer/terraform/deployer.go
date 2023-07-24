package terraform

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
	"k8s.io/utils/pointer"

	"github.com/seal-io/seal/pkg/auths"
	"github.com/seal-io/seal/pkg/auths/session"
	revisionbus "github.com/seal-io/seal/pkg/bus/servicerevision"
	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/model/environmentconnectorrelationship"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/service"
	"github.com/seal-io/seal/pkg/dao/model/servicerelationship"
	"github.com/seal-io/seal/pkg/dao/model/serviceresource"
	"github.com/seal-io/seal/pkg/dao/model/servicerevision"
	"github.com/seal-io/seal/pkg/dao/model/templateversion"
	"github.com/seal-io/seal/pkg/dao/model/variable"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/crypto"
	"github.com/seal-io/seal/pkg/dao/types/object"
	"github.com/seal-io/seal/pkg/dao/types/property"
	"github.com/seal-io/seal/pkg/dao/types/status"
	deptypes "github.com/seal-io/seal/pkg/deployer/types"
	pkgenv "github.com/seal-io/seal/pkg/environment"
	opk8s "github.com/seal-io/seal/pkg/operator/k8s"
	pkgservice "github.com/seal-io/seal/pkg/service"
	"github.com/seal-io/seal/pkg/settings"
	"github.com/seal-io/seal/pkg/terraform/config"
	"github.com/seal-io/seal/pkg/terraform/parser"
	"github.com/seal-io/seal/pkg/terraform/util"
	"github.com/seal-io/seal/utils/json"
	"github.com/seal-io/seal/utils/log"
)

// DeployerType the type of deployer.
const DeployerType = types.DeployerTypeTF

// Deployer terraform deployer to deploy the service.
type Deployer struct {
	logger      log.Logger
	modelClient model.ClientSet
	clientSet   *kubernetes.Clientset
}

type CreateRevisionOptions struct {
	// JobType indicates the type of the job.
	JobType string
	Tags    []string
	Service *model.Service
}

// CreateSecretsOptions options for creating deployment job secrets.
type CreateSecretsOptions struct {
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

// CreateJobOptions options for do job action.
type CreateJobOptions struct {
	Type          string
	SkipTLSVerify bool
	Service       *model.Service
	// ServiceRevision indicates the service revision to create the deploy job.
	ServiceRevision *model.ServiceRevision
}

// _backendAPI the API path to terraform deploy backend.
// Terraform will get and update deployment states from this API.
const _backendAPI = "/v1/service-revisions/%s/terraform-states?projectID=%s"

// _variablePrefix the prefix of the variable name.
const _variablePrefix = "_seal_var_"

// _servicePrefix the prefix of the service output name.
const _servicePrefix = "_seal_service_"

var (
	// _variableReg the regexp to match the variable.
	_variableReg = regexp.MustCompile(`\${var\.([a-zA-Z0-9_-]+)}`)
	// _serviceReg the regexp to match the service output.
	_serviceReg = regexp.MustCompile(`\${service\.([^.}]+)\.([^.}]+)}`)
)

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

// Apply deploys the service.
func (d Deployer) Apply(ctx context.Context, service *model.Service, opts deptypes.ApplyOptions) (err error) {
	revision, err := d.CreateServiceRevision(ctx, CreateRevisionOptions{
		JobType: JobTypeApply,
		Tags:    opts.Tags,
		Service: service,
	})
	if err != nil {
		return err
	}

	defer func() {
		if err == nil {
			return
		}

		// Report to service revision.
		_ = d.updateRevisionStatus(ctx, revision, status.ServiceRevisionStatusFailed, err.Error())
	}()

	return d.CreateK8sJob(ctx, CreateJobOptions{
		Type:            JobTypeApply,
		SkipTLSVerify:   opts.SkipTLSVerify,
		Service:         service,
		ServiceRevision: revision,
	})
}

// Destroy will destroy the resource of the service.
// 1. Get the latest revision, and checkAppRevision it if it is running.
// 2. If not running, then destroy resources.
func (d Deployer) Destroy(
	ctx context.Context,
	service *model.Service,
	destroyOpts deptypes.DestroyOptions,
) (err error) {
	sr, err := d.CreateServiceRevision(ctx, CreateRevisionOptions{
		JobType: JobTypeDestroy,
		Service: service,
	})
	if err != nil {
		return err
	}

	defer func() {
		if err == nil {
			return
		}
		// Report to service revision.
		_ = d.updateRevisionStatus(ctx, sr, status.ServiceRevisionStatusFailed, err.Error())
	}()

	// If no resource exists, skip job and set revision status succeed.
	exist, err := d.modelClient.ServiceResources().Query().
		Where(serviceresource.ServiceID(service.ID)).
		Exist(ctx)
	if err != nil {
		return err
	}

	if !exist {
		return d.updateRevisionStatus(ctx, sr, status.ServiceRevisionStatusSucceeded, sr.StatusMessage)
	}

	return d.CreateK8sJob(ctx, CreateJobOptions{
		Type:            JobTypeDestroy,
		SkipTLSVerify:   destroyOpts.SkipTLSVerify,
		Service:         service,
		ServiceRevision: sr,
	})
}

// CreateK8sJob will create a k8s job to deploy、destroy or rollback the service.
func (d Deployer) CreateK8sJob(ctx context.Context, opts CreateJobOptions) error {
	connectors, err := d.getConnectors(ctx, opts.Service)
	if err != nil {
		return err
	}

	project, err := d.modelClient.Projects().Get(ctx, opts.Service.ProjectID)
	if err != nil {
		return err
	}

	environment, err := dao.GetEnvironmentByID(ctx, d.modelClient, opts.Service.EnvironmentID)
	if err != nil {
		return err
	}

	var subjectID object.ID

	sj, _ := session.GetSubject(ctx)
	if sj.ID != "" {
		subjectID = sj.ID
	} else {
		subjectID, err = pkgservice.GetSubjectID(opts.Service)
		if err != nil {
			return err
		}
	}

	if subjectID == "" {
		return errors.New("subject id is empty")
	}

	// Prepare tfConfig for deployment.
	secretOpts := CreateSecretsOptions{
		SkipTLSVerify:   opts.SkipTLSVerify,
		ServiceRevision: opts.ServiceRevision,
		Connectors:      connectors,
		ProjectID:       opts.Service.ProjectID,
		EnvironmentID:   opts.Service.EnvironmentID,
		SubjectID:       subjectID,
		// Metadata.
		ProjectName:          project.Name,
		EnvironmentName:      environment.Name,
		ServiceName:          opts.Service.Name,
		ServiceID:            opts.Service.ID,
		ManagedNamespaceName: pkgenv.GetManagedNamespaceName(environment),
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

func (d Deployer) updateRevisionStatus(ctx context.Context, ar *model.ServiceRevision, s, m string) error {
	// Report to service revision.
	ar.Status = s
	ar.StatusMessage = m

	ar, err := d.modelClient.ServiceRevisions().UpdateOne(ar).
		SetStatus(ar.Status).
		SetStatusMessage(ar.StatusMessage).
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

// createK8sSecrets will create the k8s secrets for deployment.
func (d Deployer) createK8sSecrets(ctx context.Context, opts CreateSecretsOptions) error {
	secretData := make(map[string][]byte)
	// SecretName terraform tfConfig name.
	secretName := _jobSecretPrefix + string(opts.ServiceRevision.ID)

	// Prepare terraform config files bytes for deployment.
	terraformData, err := d.LoadConfigsBytes(ctx, opts)
	if err != nil {
		return err
	}

	for k, v := range terraformData {
		secretData[k] = v
	}

	// Mount the provider configs(e.g. kubeconfig) to secret.
	providerData, err := d.GetProviderSecretData(opts.Connectors)
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

// CreateServiceRevision will create a new service revision.
// Get the latest revision, and check it if it is running.
// If not running, then apply the latest revision.
// If running, then wait for the latest revision to be applied.
func (d Deployer) CreateServiceRevision(
	ctx context.Context,
	opts CreateRevisionOptions,
) (*model.ServiceRevision, error) {
	entity := &model.ServiceRevision{
		ProjectID:       opts.Service.ProjectID,
		ServiceID:       opts.Service.ID,
		EnvironmentID:   opts.Service.EnvironmentID,
		TemplateID:      opts.Service.Template.ID,
		TemplateVersion: opts.Service.Template.Version,
		Attributes:      opts.Service.Attributes,
		Tags:            opts.Tags,
		DeployerType:    DeployerType,
		Status:          status.ServiceRevisionStatusRunning,
	}

	// Output of the previous revision should be inherited to the new one
	// when creating a new revision.
	prevEntity, err := d.modelClient.ServiceRevisions().Query().
		Where(servicerevision.And(
			servicerevision.ServiceID(opts.Service.ID),
			servicerevision.DeployerType(DeployerType))).
		Order(model.Desc(servicerevision.FieldCreateTime)).
		First(ctx)
	if err != nil && !model.IsNotFound(err) {
		return nil, err
	}

	if prevEntity != nil {
		if prevEntity.Status == status.ServiceRevisionStatusRunning {
			return nil, errors.New("service deployment is running")
		}

		// Inherit the output of previous revision.
		entity.Output = prevEntity.Output

		// Get required providers.
		requiredProviders, err := d.getRequiredProviders(ctx, opts.Service.ID, entity.Output)
		if err != nil {
			return nil, err
		}
		entity.PreviousRequiredProviders = requiredProviders
	}

	if opts.JobType == JobTypeDestroy &&
		prevEntity != nil &&
		prevEntity.Status == status.ServiceRevisionStatusSucceeded {
		entity.TemplateID = prevEntity.TemplateID
		entity.TemplateVersion = prevEntity.TemplateVersion
		entity.Attributes = prevEntity.Attributes
		entity.InputPlan = prevEntity.InputPlan
		entity.Output = prevEntity.Output
		entity.PreviousRequiredProviders = prevEntity.PreviousRequiredProviders
	}

	// Create revision, mark status to running.
	entity, err = d.modelClient.ServiceRevisions().Create().
		Set(entity).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	entity.Edges.Service = opts.Service

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

// LoadConfigsBytes returns terraform main.tf and terraform.tfvars for deployment.
func (d Deployer) LoadConfigsBytes(ctx context.Context, opts CreateSecretsOptions) (map[string][]byte, error) {
	logger := log.WithName("deployer").WithName("tf")
	// Prepare terraform tfConfig.
	//  get module configs from service revision.
	moduleConfig, providerRequirements, err := d.GetModuleConfig(ctx, opts)
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

	// Parse module attributes.
	variables, dependencyOutputs, err := d.parseModuleAttributes(ctx, moduleConfig, opts)
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
		fmt.Sprintf(_backendAPI, opts.ServiceRevision.ID, opts.ProjectID))

	// Prepare API token for terraform backend.
	const _30mins = 1800

	at, err := auths.CreateAccessToken(ctx,
		d.modelClient, opts.SubjectID, types.TokenKindDeployment, string(opts.ServiceRevision.ID), pointer.Int(_30mins))
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
				Token:                at.Value,
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

// GetProviderSecretData returns provider kubeconfig secret data mount into terraform container.
func (d Deployer) GetProviderSecretData(connectors model.Connectors) (map[string][]byte, error) {
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

// GetModuleConfig returns module configs and required connectors to
// get terraform module config block from service revision.
func (d Deployer) GetModuleConfig(
	ctx context.Context,
	opts CreateSecretsOptions,
) (*config.ModuleConfig, []types.ProviderRequirement, error) {
	var (
		requiredProviders = make([]types.ProviderRequirement, 0)
		predicates        = make([]predicate.TemplateVersion, 0)
	)

	predicates = append(predicates, templateversion.And(
		templateversion.TemplateID(opts.ServiceRevision.TemplateID),
		templateversion.Version(opts.ServiceRevision.TemplateVersion),
	))

	moduleVersion, err := d.modelClient.TemplateVersions().
		Query().
		Select(
			templateversion.FieldID,
			templateversion.FieldTemplateID,
			templateversion.FieldVersion,
			templateversion.FieldSource,
			templateversion.FieldSchema,
		).
		Where(templateversion.Or(predicates...)).
		Only(ctx)
	if err != nil {
		return nil, nil, err
	}

	if moduleVersion.Schema != nil {
		requiredProviders = append(requiredProviders, moduleVersion.Schema.RequiredProviders...)
	}

	mc, err := getModuleConfig(opts.ServiceRevision, moduleVersion, opts)
	if err != nil {
		return nil, nil, err
	}

	return mc, requiredProviders, err
}

func (d Deployer) getConnectors(ctx context.Context, ai *model.Service) (model.Connectors, error) {
	rs, err := d.modelClient.EnvironmentConnectorRelationships().Query().
		Where(environmentconnectorrelationship.EnvironmentID(ai.EnvironmentID)).
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

// parseModuleAttributes parse module variables and dependencies, return matched model.Variables and service output.
func (d Deployer) parseModuleAttributes(
	ctx context.Context,
	templateConfig *config.ModuleConfig,
	opts CreateSecretsOptions,
) (variables model.Variables, outputs map[string]parser.OutputState, err error) {
	var (
		templateVariables        []string
		dependencyServiceOutputs []string
	)

	templateVariables, dependencyServiceOutputs = parseAttributeReplace(
		templateConfig.Attributes,
		templateVariables,
		dependencyServiceOutputs,
	)

	// If service revision has variables that inherit from cloned revision, use them directly.
	if len(opts.ServiceRevision.Variables) > 0 {
		for k, v := range opts.ServiceRevision.Variables {
			variables = append(variables, &model.Variable{
				Name:  k,
				Value: crypto.String(v),
			})
		}
	} else {
		variables, err = d.getVariables(ctx, templateVariables, opts.ProjectID, opts.EnvironmentID)
		if err != nil {
			return nil, nil, err
		}
	}

	outputs, err = d.getServiceDependencyOutputs(ctx, opts.ServiceRevision.ServiceID, dependencyServiceOutputs)
	if err != nil {
		return nil, nil, err
	}

	// Check if all dependency service outputs are found.
	for _, o := range dependencyServiceOutputs {
		if _, ok := outputs[o]; !ok {
			return nil, nil, fmt.Errorf("service %s dependency output %s not found", opts.ServiceName, o)
		}
	}

	return variables, outputs, nil
}

func (d Deployer) getVariables(
	ctx context.Context,
	variableNames []string,
	projectID, environmentID object.ID,
) (model.Variables, error) {
	nameIn := make([]any, len(variableNames))
	for i, name := range variableNames {
		nameIn[i] = name
	}

	s, err := session.GetSubject(ctx)
	if err == nil {
		s.IncognitoOn()
		defer s.IncognitoOff()
	}

	type scanVariable struct {
		Name      string        `json:"name"`
		Value     crypto.String `json:"value"`
		Sensitive bool          `json:"sensitive"`
		Scope     int           `json:"scope"`
	}

	var vars []scanVariable

	err = d.modelClient.Variables().Query().
		Modify(func(s *sql.Selector) {
			var (
				envPs = sql.And(
					sql.EQ(variable.FieldProjectID, projectID),
					sql.EQ(variable.FieldEnvironmentID, environmentID),
				)

				projPs = sql.And(
					sql.EQ(variable.FieldProjectID, projectID),
					sql.IsNull(variable.FieldEnvironmentID),
				)

				globalPs = sql.IsNull(variable.FieldProjectID)
			)

			s.Where(
				sql.And(
					sql.In(variable.FieldName, nameIn...),
					sql.Or(
						envPs,
						projPs,
						globalPs,
					),
				),
			).SelectExpr(
				sql.Expr("CASE "+
					"WHEN project_id IS NOT NULL AND environment_id IS NOT NULL THEN 3 "+
					"WHEN project_id IS NOT NULL AND environment_id IS NULL THEN 2 "+
					"ELSE 1 "+
					"END AS scope"),
			).AppendSelect(
				variable.FieldName,
				variable.FieldValue,
				variable.FieldSensitive,
			)
		}).
		Scan(ctx, &vars)
	if err != nil {
		return nil, err
	}

	found := make(map[string]scanVariable)
	for _, v := range vars {
		ev, ok := found[v.Name]
		if !ok {
			found[v.Name] = v
			continue
		}

		if v.Scope > ev.Scope {
			found[v.Name] = v
		}
	}

	var variables model.Variables

	// Validate module variable are all exist.
	foundSet := sets.NewString()
	for n, e := range found {
		foundSet.Insert(n)
		variables = append(variables, &model.Variable{
			Name:      n,
			Value:     e.Value,
			Sensitive: e.Sensitive,
		})
	}
	requiredSet := sets.NewString(variableNames...)

	missingSet := requiredSet.Difference(foundSet)
	if missingSet.Len() > 0 {
		return nil, fmt.Errorf("missing variables: %s", missingSet.List())
	}

	return variables, nil
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

// getServiceDependencyOutputs gets the dependency outputs of the service.
func (d Deployer) getServiceDependencyOutputs(
	ctx context.Context,
	serviceID object.ID,
	dependOutputs []string,
) (map[string]parser.OutputState, error) {
	entity, err := d.modelClient.Services().Query().
		Where(service.ID(serviceID)).
		WithDependencies(func(sq *model.ServiceRelationshipQuery) {
			sq.Where(func(s *sql.Selector) {
				s.Where(sql.ColumnsNEQ(servicerelationship.FieldServiceID, servicerelationship.FieldDependencyID))
			})
		}).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	dependencyServiceIDs := make([]object.ID, 0, len(entity.Edges.Dependencies))

	for _, d := range entity.Edges.Dependencies {
		if d.Type != types.ServiceRelationshipTypeImplicit {
			continue
		}

		dependencyServiceIDs = append(dependencyServiceIDs, d.DependencyID)
	}

	dependencyRevisions, err := d.modelClient.ServiceRevisions().Query().
		Select(
			servicerevision.FieldID,
			servicerevision.FieldAttributes,
			servicerevision.FieldOutput,
			servicerevision.FieldServiceID,
			servicerevision.FieldProjectID,
		).
		Where(func(s *sql.Selector) {
			sq := s.Clone().
				AppendSelectExprAs(
					sql.RowNumber().
						PartitionBy(servicerevision.FieldServiceID).
						OrderBy(sql.Desc(servicerevision.FieldCreateTime)),
					"row_number",
				).
				Where(s.P()).
				From(s.Table()).
				As(servicerevision.Table)

			// Query the latest revision of the service.
			s.Where(sql.EQ(s.C("row_number"), 1)).
				From(sq)
		}).Where(servicerevision.ServiceIDIn(dependencyServiceIDs...)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	outputs := make(map[string]parser.OutputState, 0)
	dependSets := sets.NewString(dependOutputs...)

	for _, r := range dependencyRevisions {
		revisionOutput, err := parser.ParseStateOutputRawMap(r)
		if err != nil {
			return nil, err
		}

		for n, o := range revisionOutput {
			if dependSets.Has(n) {
				outputs[n] = o
			}
		}
	}

	return outputs, nil
}

func SyncServiceRevisionStatus(ctx context.Context, bm revisionbus.BusMessage) (err error) {
	var (
		mc       = bm.TransactionalModelClient
		revision = bm.Refer
	)

	// Report to service.
	entity, err := mc.Services().Query().
		Where(service.ID(revision.ServiceID)).
		Select(
			service.FieldID,
			service.FieldStatus,
		).
		Only(ctx)
	if err != nil {
		return err
	}

	switch revision.Status {
	case status.ServiceRevisionStatusSucceeded:
		if status.ServiceStatusDeleted.IsUnknown(entity) {
			return mc.Services().DeleteOne(entity).
				Exec(ctx)
		}

		status.ServiceStatusDeployed.True(entity, "")
		status.ServiceStatusReady.Unknown(entity, "")
	case status.ServiceRevisionStatusFailed:
		if status.ServiceStatusDeleted.IsUnknown(entity) {
			status.ServiceStatusDeleted.False(entity, "")
		} else {
			status.ServiceStatusDeployed.False(entity, "")
		}

		entity.Status.SummaryStatusMessage = revision.StatusMessage
	}

	entity.Status.SetSummary(status.WalkService(&entity.Status))

	return mc.Services().UpdateOne(entity).
		SetStatus(entity.Status).
		Exec(ctx)
}

// parseAttributeReplace parses attribute variable ${var.name} replaces it with ${var._variablePrefix+name},
// service reference ${service.name.output} replaces it with ${var._servicePrefix+name}
// and returns variable names and service names.
func parseAttributeReplace(
	attributes map[string]interface{},
	variableNames []string,
	serviceOutputs []string,
) ([]string, []string) {
	for key, value := range attributes {
		if value == nil {
			continue
		}

		switch reflect.TypeOf(value).Kind() {
		case reflect.Map:
			if _, ok := value.(map[string]interface{}); !ok {
				continue
			}

			variableNames, serviceOutputs = parseAttributeReplace(
				value.(map[string]interface{}),
				variableNames,
				serviceOutputs,
			)
		case reflect.String:
			str := value.(string)
			matches := _variableReg.FindAllStringSubmatch(str, -1)
			serviceMatches := _serviceReg.FindAllStringSubmatch(str, -1)

			var matched []string

			for _, match := range matches {
				if len(match) > 1 {
					matched = append(matched, match[1])
				}
			}

			var serviceMatched []string

			for _, match := range serviceMatches {
				if len(match) > 1 {
					serviceMatched = append(serviceMatched, match[1]+"_"+match[2])
				}
			}

			variableNames = append(variableNames, matched...)
			variableRepl := "${var." + _variablePrefix + "${1}}"
			str = _variableReg.ReplaceAllString(str, variableRepl)

			serviceOutputs = append(serviceOutputs, serviceMatched...)
			serviceRepl := "${var." + _servicePrefix + "${1}_${2}}"

			attributes[key] = _serviceReg.ReplaceAllString(str, serviceRepl)
		case reflect.Slice:
			if _, ok := value.([]interface{}); !ok {
				continue
			}

			for _, v := range value.([]interface{}) {
				if _, ok := v.(map[string]interface{}); !ok {
					continue
				}
				variableNames, serviceOutputs = parseAttributeReplace(
					v.(map[string]interface{}),
					variableNames,
					serviceOutputs,
				)
			}
		}
	}

	return variableNames, serviceOutputs
}

func getVarConfigOptions(variables model.Variables, serviceOutputs map[string]parser.OutputState) config.CreateOptions {
	varsConfigOpts := config.CreateOptions{
		Attributes: map[string]interface{}{},
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
	modVer *model.TemplateVersion,
	ops CreateSecretsOptions,
) (*config.ModuleConfig, error) {
	var (
		props              = make(property.Properties, len(revision.Attributes))
		typesWith          = revision.Attributes.TypesWith(modVer.Schema.Variables)
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
		Name:       revision.Edges.Service.Name,
		Source:     modVer.Source,
		Schema:     modVer.Schema,
		Attributes: attrs,
	}

	if modVer.Schema == nil {
		return mc, nil
	}

	for _, v := range modVer.Schema.Variables {
		// Add sensitive from schema variable.
		if v.Sensitive {
			sensitiveVariables.Insert(fmt.Sprintf(`var\.%s`, v.Name))
		}

		// Add seal metadata.
		var attrValue string

		switch v.Name {
		case SealMetadataProjectName:
			attrValue = ops.ProjectName
		case SealMetadataEnvironmentName:
			attrValue = ops.EnvironmentName
		case SealMetadataServiceName:
			attrValue = ops.ServiceName
		case SealMetadataProjectID:
			attrValue = ops.ProjectID.String()
		case SealMetadataEnvironmentID:
			attrValue = ops.EnvironmentID.String()
		case SealMetadataServiceID:
			attrValue = ops.ServiceID.String()
		case SealMetadataNamespaceName:
			attrValue = ops.ManagedNamespaceName
		}

		if attrValue != "" {
			mc.Attributes[v.Name] = attrValue
		}
	}

	sensitiveVariableRegex, err := matchAnyRegex(sensitiveVariables.UnsortedList())
	if err != nil {
		return nil, err
	}

	mc.Outputs = make([]config.Output, len(modVer.Schema.Outputs))
	for i, v := range modVer.Schema.Outputs {
		mc.Outputs[i].Sensitive = v.Sensitive
		mc.Outputs[i].Name = v.Name
		mc.Outputs[i].ServiceName = revision.Edges.Service.Name
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
