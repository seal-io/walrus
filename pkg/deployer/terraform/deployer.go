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
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
	"k8s.io/utils/pointer"

	"github.com/seal-io/seal/pkg/auths"
	"github.com/seal-io/seal/pkg/auths/session"
	revisionbus "github.com/seal-io/seal/pkg/bus/applicationrevision"
	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/application"
	"github.com/seal-io/seal/pkg/dao/model/applicationinstance"
	"github.com/seal-io/seal/pkg/dao/model/applicationmodulerelationship"
	"github.com/seal-io/seal/pkg/dao/model/applicationresource"
	"github.com/seal-io/seal/pkg/dao/model/applicationrevision"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/model/environmentconnectorrelationship"
	"github.com/seal-io/seal/pkg/dao/model/moduleversion"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/project"
	"github.com/seal-io/seal/pkg/dao/model/secret"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/crypto"
	"github.com/seal-io/seal/pkg/dao/types/oid"
	"github.com/seal-io/seal/pkg/dao/types/property"
	"github.com/seal-io/seal/pkg/dao/types/status"
	deptypes "github.com/seal-io/seal/pkg/deployer/types"
	opk8s "github.com/seal-io/seal/pkg/operator/k8s"
	"github.com/seal-io/seal/pkg/settings"
	"github.com/seal-io/seal/pkg/terraform/config"
	"github.com/seal-io/seal/pkg/terraform/parser"
	"github.com/seal-io/seal/pkg/terraform/util"
	"github.com/seal-io/seal/pkg/topic/datamessage"
	"github.com/seal-io/seal/utils/log"
)

// DeployerType the type of deployer.
const DeployerType = types.DeployerTypeTF

// Deployer terraform deployer to deploy the application.
type Deployer struct {
	logger      log.Logger
	modelClient model.ClientSet
	clientSet   *kubernetes.Clientset
}

type CreateRevisionOptions struct {
	// JobType indicates the type of the job.
	JobType             string
	Tags                []string
	Application         *model.Application
	ApplicationInstance *model.ApplicationInstance
	// CloneFrom indicates the application revision to clone from.
	CloneFrom *model.ApplicationRevision
}

// CreateSecretsOptions options for creating deployment job secrets.
type CreateSecretsOptions struct {
	SkipTLSVerify       bool
	ApplicationRevision *model.ApplicationRevision
	Connectors          model.Connectors
	ProjectID           oid.ID
	// Metadata.
	ProjectName             string
	ApplicationName         string
	ApplicationInstanceName string
}

// CreateJobOptions options for do job action.
type CreateJobOptions struct {
	Type                string
	SkipTLSVerify       bool
	Application         *model.Application
	ApplicationInstance *model.ApplicationInstance
	// ApplicationRevision indicates the application revision
	// to create the deploy job.
	ApplicationRevision *model.ApplicationRevision
}

// _backendAPI the API path to terraform deploy backend.
// Terraform will get and update deployment states from this API.
const _backendAPI = "/v1/application-revisions/%s/terraform-states?projectID=%s"

// _varPrefix the prefix of the variable name.
const _varPrefix = "_seal_var_"

// _secretPrefix the prefix of the secret name.
const _secretPrefix = "_seal_secret_"

var (
	// _secretReg the regexp to match the secret variable.
	_secretReg = regexp.MustCompile(`\${secret\.([a-zA-Z0-9_]+)}`)
	// _varReg the regexp to match the variable.
	_varReg = regexp.MustCompile(`\${var\.([a-zA-Z0-9_]+)}`)
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

// Apply will apply the application to deploy the application.
func (d Deployer) Apply(ctx context.Context, ai *model.ApplicationInstance, opts deptypes.ApplyOptions) (err error) {
	app, err := d.getApplication(ctx, ai.ApplicationID)
	if err != nil {
		return err
	}

	ar, err := d.CreateApplicationRevision(ctx, CreateRevisionOptions{
		JobType:             JobTypeApply,
		Tags:                opts.Tags,
		Application:         app,
		ApplicationInstance: ai,
		CloneFrom:           opts.CloneFrom,
	})
	if err != nil {
		return err
	}

	defer func() {
		if err == nil {
			return
		}
		// Report to application revision.
		_ = d.updateRevisionStatus(ctx, ar, status.ApplicationRevisionStatusFailed, err.Error())
	}()

	return d.CreateK8sJob(ctx, CreateJobOptions{
		Type:                JobTypeApply,
		SkipTLSVerify:       opts.SkipTLSVerify,
		Application:         app,
		ApplicationInstance: ai,
		ApplicationRevision: ar,
	})
}

// Destroy will destroy the resource of the application.
// 1. Get the latest revision, and checkAppRevision it if it is running.
// 2. If not running, then destroy resources.
func (d Deployer) Destroy(
	ctx context.Context,
	ai *model.ApplicationInstance,
	destroyOpts deptypes.DestroyOptions,
) (err error) {
	app, err := d.getApplication(ctx, ai.ApplicationID)
	if err != nil {
		return err
	}

	ar, err := d.CreateApplicationRevision(ctx, CreateRevisionOptions{
		JobType:             JobTypeDestroy,
		Application:         app,
		ApplicationInstance: ai,
	})
	if err != nil {
		return err
	}

	defer func() {
		if err == nil {
			return
		}
		// Report to application revision.
		_ = d.updateRevisionStatus(ctx, ar, status.ApplicationRevisionStatusFailed, err.Error())
	}()

	// If no resource exists, skip job and set revision status succeed.
	exist, err := d.modelClient.ApplicationResources().Query().
		Where(applicationresource.InstanceID(ai.ID)).
		Exist(ctx)
	if err != nil {
		return err
	}

	if !exist {
		return d.updateRevisionStatus(ctx, ar, status.ApplicationRevisionStatusSucceeded, ar.StatusMessage)
	}

	return d.CreateK8sJob(ctx, CreateJobOptions{
		Type:                JobTypeDestroy,
		SkipTLSVerify:       destroyOpts.SkipTLSVerify,
		Application:         app,
		ApplicationInstance: ai,
		ApplicationRevision: ar,
	})
}

// Rollback instance to a specific revision.
func (d Deployer) Rollback(
	ctx context.Context,
	ai *model.ApplicationInstance,
	opts deptypes.RollbackOptions,
) (err error) {
	if opts.CloneFrom == nil || opts.CloneFrom.InstanceID != ai.ID {
		return errors.New("rollback failed: invalid revision")
	}

	app, err := d.getApplication(ctx, ai.ApplicationID)
	if err != nil {
		return err
	}

	status.ApplicationInstanceStatusDeployed.Reset(ai, "Rolling back")

	update, err := dao.ApplicationInstanceUpdate(d.modelClient, ai)
	if err != nil {
		return err
	}

	ai, err = update.Save(ctx)
	if err != nil {
		return err
	}

	ar, err := d.CreateApplicationRevision(ctx,
		CreateRevisionOptions{
			JobType:             JobTypeApply,
			Tags:                opts.Tags,
			Application:         app,
			ApplicationInstance: ai,
			CloneFrom:           opts.CloneFrom,
		},
	)
	if err != nil {
		return err
	}

	defer func() {
		if err == nil {
			return
		}

		if ar != nil {
			// Report to application revision.
			_ = d.updateRevisionStatus(ctx, ar, status.ApplicationRevisionStatusFailed, err.Error())
			return
		}

		status.ApplicationInstanceStatusDeployed.False(ai, err.Error())

		instanceUpdate, updateErr := dao.ApplicationInstanceUpdate(d.modelClient, ai)
		if updateErr != nil {
			d.logger.Error(err)
			return
		}

		updateErr = instanceUpdate.Exec(ctx)
		if updateErr != nil {
			d.logger.Errorf("update application instance status failed: %v", updateErr)
		}
	}()

	return d.CreateK8sJob(ctx, CreateJobOptions{
		Type:                JobTypeApply,
		SkipTLSVerify:       opts.SkipTLSVerify,
		Application:         app,
		ApplicationInstance: ai,
		ApplicationRevision: ar,
	})
}

// getApplication will get the application by id.
func (d Deployer) getApplication(ctx context.Context, id oid.ID) (*model.Application, error) {
	return d.modelClient.Applications().Query().
		Where(application.ID(id)).
		WithProject(func(pq *model.ProjectQuery) {
			pq.Select(
				project.FieldID,
				project.FieldName,
			)
		}).
		Only(ctx)
}

// CreateK8sJob will create a k8s job to deploy、destroy or rollback the application instance.
func (d Deployer) CreateK8sJob(ctx context.Context, opts CreateJobOptions) error {
	connectors, err := d.getConnectors(ctx, opts.ApplicationInstance)
	if err != nil {
		return err
	}

	// Prepare tfConfig for deployment.
	secretOpts := CreateSecretsOptions{
		SkipTLSVerify:       opts.SkipTLSVerify,
		ApplicationRevision: opts.ApplicationRevision,
		Connectors:          connectors,
		ProjectID:           opts.Application.ProjectID,
		// Metadata.
		ProjectName:             opts.Application.Edges.Project.Name,
		ApplicationName:         opts.Application.Name,
		ApplicationInstanceName: opts.ApplicationInstance.Name,
	}
	if err = d.createK8sSecrets(ctx, secretOpts); err != nil {
		return err
	}

	jobImage, err := settings.TerraformDeployerImage.Value(ctx, d.modelClient)
	if err != nil {
		return err
	}

	// Create deployment job.
	jobOpts := JobCreateOptions{
		Type:                  opts.Type,
		ApplicationRevisionID: opts.ApplicationRevision.ID.String(),
		Image:                 jobImage,
	}

	return CreateJob(ctx, d.clientSet, jobOpts)
}

func (d Deployer) updateRevisionStatus(ctx context.Context, ar *model.ApplicationRevision, s, m string) error {
	// Report to application revision.
	ar.Status = s
	ar.StatusMessage = m

	update, err := dao.ApplicationRevisionUpdate(d.modelClient, ar)
	if err != nil {
		return err
	}

	ar, err = update.Save(ctx)
	if err != nil {
		d.logger.Error(err)
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
	secretName := _jobSecretPrefix + string(opts.ApplicationRevision.ID)

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

// CreateApplicationRevision will create a new application revision.
// Get the latest revision, and check it if it is running.
// If not running, then apply the latest revision.
// If running, then wait for the latest revision to be applied.
func (d Deployer) CreateApplicationRevision(
	ctx context.Context,
	opts CreateRevisionOptions,
) (*model.ApplicationRevision, error) {
	var entity *model.ApplicationRevision
	if opts.CloneFrom != nil {
		entity = &model.ApplicationRevision{
			Modules:        opts.CloneFrom.Modules,
			Secrets:        opts.CloneFrom.Secrets,
			Variables:      opts.CloneFrom.Variables,
			InputVariables: opts.CloneFrom.InputVariables,
		}
	} else {
		entity = &model.ApplicationRevision{
			Variables:      opts.Application.Variables,
			InputVariables: opts.ApplicationInstance.Variables,
		}
		// Get modules for new revision.
		amrs, err := d.modelClient.ApplicationModuleRelationships().Query().
			Where(applicationmodulerelationship.ApplicationID(opts.Application.ID)).
			All(ctx)
		if err != nil {
			return nil, err
		}
		entity.Modules = make([]types.ApplicationModule, len(amrs))
		for i := range amrs {
			entity.Modules[i] = types.ApplicationModule{
				ModuleID:   amrs[i].ModuleID,
				Version:    amrs[i].Version,
				Name:       amrs[i].Name,
				Attributes: amrs[i].Attributes,
			}
		}
	}

	entity.Tags = opts.Tags
	entity.Status = status.ApplicationRevisionStatusRunning
	entity.DeployerType = DeployerType
	entity.InstanceID = opts.ApplicationInstance.ID
	entity.EnvironmentID = opts.ApplicationInstance.EnvironmentID

	// Output of the previous revision should be inherited to the new one
	// when creating a new revision.
	prevEntity, err := d.modelClient.ApplicationRevisions().Query().
		Where(applicationrevision.And(
			applicationrevision.InstanceID(opts.ApplicationInstance.ID),
			applicationrevision.DeployerType(DeployerType))).
		Order(model.Desc(applicationrevision.FieldCreateTime)).
		First(ctx)
	if err != nil && !model.IsNotFound(err) {
		return nil, err
	}

	if prevEntity != nil {
		if prevEntity.Status == status.ApplicationRevisionStatusRunning {
			return nil, errors.New("application deployment is running")
		}
		// Inherit the output of previous revision.
		entity.Output = prevEntity.Output

		requiredProviders, err := d.getRequiredProviders(ctx, opts.ApplicationInstance.ID, entity.Output)
		if err != nil {
			return nil, err
		}
		entity.PreviousRequiredProviders = requiredProviders
	}

	if opts.JobType == JobTypeDestroy &&
		prevEntity != nil &&
		prevEntity.Status == status.ApplicationRevisionStatusSucceeded {
		entity.Modules = prevEntity.Modules
		entity.InputVariables = prevEntity.InputVariables
		entity.InputPlan = prevEntity.InputPlan
		entity.Output = prevEntity.Output
		entity.PreviousRequiredProviders = prevEntity.PreviousRequiredProviders
	}

	// Create revision, mark status to running.
	creates, err := dao.ApplicationRevisionCreates(d.modelClient, entity)
	if err != nil {
		return nil, err
	}

	return creates[0].Save(ctx)
}

func (d Deployer) getRequiredProviders(
	ctx context.Context,
	instanceID oid.ID,
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
	//  get module configs from app revision.
	moduleConfigs, providerRequirements, err := d.GetModuleConfigs(ctx, opts)
	if err != nil {
		return nil, err
	}
	// Merge current and previous required providers.
	providerRequirements = append(providerRequirements,
		opts.ApplicationRevision.PreviousRequiredProviders...)

	requiredProviders := make(map[string]*tfconfig.ProviderRequirement, 0)
	for _, p := range providerRequirements {
		if _, ok := requiredProviders[p.Name]; !ok {
			requiredProviders[p.Name] = p.ProviderRequirement
		} else {
			logger.Warnf("duplicate provider requirement: %s", p.Name)
		}
	}

	// Parse module secrets.
	secrets, err := d.parseModuleSecrets(ctx, moduleConfigs, opts)
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
		fmt.Sprintf(_backendAPI, opts.ApplicationRevision.ID, opts.ProjectID))

	// Prepare API token for terraform backend.
	const _30mins = 1800

	at, err := auths.CreateAccessToken(ctx,
		d.modelClient, types.TokenKindDeployment, string(opts.ApplicationRevision.ID), pointer.Int(_30mins))
	if err != nil {
		return nil, err
	}

	// Prepare terraform config files to be mounted to secret.
	requiredProviderNames := sets.NewString()
	for _, p := range providerRequirements {
		requiredProviderNames = requiredProviderNames.Insert(p.Name)
	}

	secretNames := make([]string, 0, len(secrets))
	for _, s := range secrets {
		secretNames = append(secretNames, s.Name)
	}

	// Prepare outputs.
	var outputCount int
	for _, v := range moduleConfigs {
		outputCount += len(v.Outputs)
	}

	outputs := make([]config.Output, 0, outputCount)
	for _, v := range moduleConfigs {
		outputs = append(outputs, v.Outputs...)
	}
	variableNameAndTypes := opts.ApplicationRevision.InputVariables.StringTypesWith(opts.ApplicationRevision.Variables)
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
				ModuleConfigs: moduleConfigs,
			},
			VariableOptions: &config.VariableOptions{
				VarPrefix:            _varPrefix,
				SecretPrefix:         _secretPrefix,
				SecretNames:          secretNames,
				VariableNameAndTypes: variableNameAndTypes,
			},
			OutputOptions: outputs,
		},
		config.FileVars: getVarConfigOptions(secrets, opts.ApplicationRevision.InputVariables),
	}
	secretMaps := make(map[string][]byte, 0)

	for k, v := range secretOptionMaps {
		secretMaps[k], err = config.CreateConfigToBytes(v)
		if err != nil {
			return nil, err
		}
	}

	// Save input plan to app revision.
	opts.ApplicationRevision.InputPlan = string(secretMaps[config.FileMain])
	// If application revision does not inherit secrets from cloned revision,
	// then save the parsed secrets to app revision.
	if len(opts.ApplicationRevision.Secrets) == 0 {
		secretMap := make(crypto.Map[string, string], len(secrets))
		for _, s := range secrets {
			secretMap[s.Name] = string(s.Value)
		}
		opts.ApplicationRevision.Secrets = secretMap
	}

	update, err := dao.ApplicationRevisionUpdate(d.modelClient, opts.ApplicationRevision)
	if err != nil {
		return nil, err
	}

	revision, err := update.Save(ctx)
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

// GetModuleConfigs returns module configs and required connectors to
// get terraform module config block from application revision.
func (d Deployer) GetModuleConfigs(
	ctx context.Context,
	opts CreateSecretsOptions,
) ([]*config.ModuleConfig, []types.ProviderRequirement, error) {
	var (
		requiredProviders = make([]types.ProviderRequirement, 0)
		// Module id -> module source.
		moduleVersionMap = make(map[string]*model.ModuleVersion, 0)
		predicates       = make([]predicate.ModuleVersion, 0)
		ar               = opts.ApplicationRevision
		moduleConfigs    = make([]*config.ModuleConfig, 0, len(ar.Modules))
	)

	for _, m := range ar.Modules {
		predicates = append(predicates, moduleversion.And(
			moduleversion.ModuleID(m.ModuleID),
			moduleversion.Version(m.Version),
		))
	}

	moduleVersions, err := d.modelClient.ModuleVersions().
		Query().
		Select(
			moduleversion.FieldID,
			moduleversion.FieldModuleID,
			moduleversion.FieldVersion,
			moduleversion.FieldSource,
			moduleversion.FieldSchema,
		).
		Where(moduleversion.Or(predicates...)).
		All(ctx)
	if err != nil {
		return nil, nil, err
	}

	moduleVersionKey := func(moduleID, version string) string {
		return fmt.Sprintf("%s/%s", moduleID, version)
	}

	for _, m := range moduleVersions {
		moduleVersionMap[moduleVersionKey(m.ModuleID, m.Version)] = m
	}

	for _, m := range ar.Modules {
		modVer, ok := moduleVersionMap[moduleVersionKey(m.ModuleID, m.Version)]
		if !ok {
			return nil, nil, fmt.Errorf("version %s of module %s not found", m.Version, m.ModuleID)
		}

		var mc *config.ModuleConfig

		mc, err = getModuleConfig(m, modVer, opts)
		if err != nil {
			return nil, nil, err
		}

		moduleConfigs = append(moduleConfigs, mc)

		if modVer.Schema != nil {
			requiredProviders = append(requiredProviders, modVer.Schema.RequiredProviders...)
		}
	}

	return moduleConfigs, requiredProviders, err
}

func (d Deployer) getConnectors(ctx context.Context, ai *model.ApplicationInstance) (model.Connectors, error) {
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

// parseModuleSecrets parse module secrets, and return matched model.Secrets.
func (d Deployer) parseModuleSecrets(
	ctx context.Context,
	moduleConfigs []*config.ModuleConfig,
	opts CreateSecretsOptions,
) (model.Secrets, error) {
	var (
		moduleSecrets []string
		secrets       model.Secrets
	)

	for _, moduleConfig := range moduleConfigs {
		moduleSecrets = parseAttributeReplace(moduleConfig.Attributes, moduleSecrets)
	}
	// If application revision has secrets that inherit from cloned revision, use them directly.
	if len(opts.ApplicationRevision.Secrets) > 0 {
		for k, v := range opts.ApplicationRevision.Secrets {
			secrets = append(secrets, &model.Secret{
				Name:  k,
				Value: crypto.String(v),
			})
		}

		return secrets, nil
	}

	nameIn := make([]interface{}, len(moduleSecrets))
	for i, name := range moduleSecrets {
		nameIn[i] = name
	}
	// This query is used to distinct the secrets with the same name.
	//  SELECT
	//    "id",
	//    "name",
	//    "value"
	//  FROM
	//    "secrets"
	//  WHERE
	//    (
	//      (
	//        "project_id" IS NULL
	//        AND "name" NOT IN (
	//          SELECT
	//            "name"
	//          FROM
	//            "secrets"
	//          WHERE
	//            "project_id" = opts.ProjectID
	//        )
	//      )
	//      OR "project_id" = opts.ProjectID
	//    )
	//    AND NAME IN (moduleSecrets)
	err := func() error {
		s, err := session.GetSubject(ctx)
		if err == nil {
			s.IncognitoOn()
			defer s.IncognitoOff()
		}

		return d.modelClient.Secrets().Query().
			Modify(func(s *sql.Selector) {
				// Select secrets without project id or not in project.
				subQuery := sql.Select(secret.FieldName).
					From(sql.Table(secret.Table)).
					Where(sql.EQ(secret.FieldProjectID, opts.ProjectID))
				s.Select(secret.FieldID, secret.FieldName, secret.FieldValue).
					Where(
						sql.And(
							sql.Or(
								sql.And(
									sql.IsNull(secret.FieldProjectID),
									sql.NotIn(secret.FieldName, subQuery),
								),
								sql.EQ(secret.FieldProjectID, opts.ProjectID),
							),
							sql.In(secret.FieldName, nameIn...),
						),
					)
			}).
			Scan(ctx, &secrets)
	}()
	if err != nil {
		return nil, err
	}

	// Validate module secret are all exist.
	foundSecretSet := sets.NewString()
	for _, s := range secrets {
		foundSecretSet.Insert(s.Name)
	}
	requiredSecretSet := sets.NewString(moduleSecrets...)

	missingSecretSet := requiredSecretSet.Difference(foundSecretSet)
	if missingSecretSet.Len() > 0 {
		return nil, fmt.Errorf("missing secrets: %s", missingSecretSet.List())
	}

	return secrets, nil
}

// getPreviousRequiredProviders get previous succeed revision required providers.
// NB(alex): the previous revision may be failed, the failed revision may not contain required providers of states.
func (d Deployer) getPreviousRequiredProviders(
	ctx context.Context,
	instanceID oid.ID,
) ([]types.ProviderRequirement, error) {
	prevRequiredProviders := make([]types.ProviderRequirement, 0)

	entity, err := d.modelClient.ApplicationRevisions().Query().
		Where(applicationrevision.InstanceID(instanceID)).
		Order(model.Desc(applicationrevision.FieldCreateTime)).
		First(ctx)
	if err != nil && !model.IsNotFound(err) {
		return nil, err
	}

	if entity == nil {
		return prevRequiredProviders, nil
	}

	predicates := make([]predicate.ModuleVersion, 0, len(entity.Modules))
	for _, m := range entity.Modules {
		predicates = append(predicates, moduleversion.And(
			moduleversion.ModuleID(m.ModuleID),
			moduleversion.Version(m.Version)))
	}

	if len(predicates) != 0 {
		mvs, err := d.modelClient.ModuleVersions().Query().
			Where(moduleversion.Or(predicates...)).
			All(ctx)
		if err != nil {
			return nil, err
		}

		for _, mv := range mvs {
			if mv.Schema == nil {
				continue
			}

			prevRequiredProviders = append(prevRequiredProviders, mv.Schema.RequiredProviders...)
		}
	}

	prevRequiredProviders = append(prevRequiredProviders, entity.PreviousRequiredProviders...)

	return prevRequiredProviders, nil
}

func SyncApplicationRevisionStatus(ctx context.Context, bm revisionbus.BusMessage) (err error) {
	var (
		mc       = bm.TransactionalModelClient
		revision = bm.Refer
	)

	// Report to application instance.
	appInstance, err := mc.ApplicationInstances().Query().
		Where(applicationinstance.ID(revision.InstanceID)).
		Select(
			applicationinstance.FieldID,
			applicationinstance.FieldStatus,
			applicationinstance.FieldApplicationID).
		Only(ctx)
	if err != nil {
		return err
	}

	var instanceUpdate *model.ApplicationInstanceUpdateOne

	switch revision.Status {
	case status.ApplicationRevisionStatusSucceeded:
		if status.ApplicationInstanceStatusDeleted.IsUnknown(appInstance) {
			err = mc.ApplicationInstances().DeleteOne(appInstance).
				Exec(ctx)
		} else {
			status.ApplicationInstanceStatusDeployed.True(appInstance, "")
			status.ApplicationInstanceStatusReady.Unknown(appInstance, "")
			instanceUpdate, err = dao.ApplicationInstanceUpdate(mc, appInstance)
			if err != nil {
				return err
			}
			err = instanceUpdate.Exec(ctx)
		}
	case status.ApplicationRevisionStatusFailed:
		if status.ApplicationInstanceStatusDeleted.IsUnknown(appInstance) {
			status.ApplicationInstanceStatusDeleted.False(appInstance, "")
		} else {
			status.ApplicationInstanceStatusDeployed.False(appInstance, "")
		}
		appInstance.Status.SummaryStatusMessage = revision.StatusMessage

		instanceUpdate, err = dao.ApplicationInstanceUpdate(mc, appInstance)
		if err != nil {
			return err
		}
		err = instanceUpdate.Exec(ctx)
	}

	if err != nil {
		return err
	}

	return datamessage.Publish(
		ctx,
		string(datamessage.Application),
		model.OpUpdate,
		[]oid.ID{appInstance.ApplicationID},
	)
}

// parseAttributeReplace parses attribute secrets ${secret.name} replaces it with ${var._varprefix+name},
// and returns secret names.
func parseAttributeReplace(
	attributes map[string]interface{},
	secretNames []string,
) []string {
	for key, value := range attributes {
		if value == nil {
			continue
		}

		switch reflect.TypeOf(value).Kind() {
		case reflect.Map:
			if _, ok := value.(map[string]interface{}); !ok {
				continue
			}

			secretNames = parseAttributeReplace(value.(map[string]interface{}), secretNames)
		case reflect.String:
			str := value.(string)
			matches := _secretReg.FindAllStringSubmatch(str, -1)

			var matched []string

			for _, match := range matches {
				if len(match) > 1 {
					matched = append(matched, match[1])
				}
			}

			secretNames = append(secretNames, matched...)
			varRepl := "${var." + _varPrefix + "${1}}"
			str = _varReg.ReplaceAllString(str, varRepl)

			secretRepl := "${var." + _secretPrefix + "${1}}"
			attributes[key] = _secretReg.ReplaceAllString(str, secretRepl)
		case reflect.Slice:
			if _, ok := value.([]interface{}); !ok {
				continue
			}

			for _, v := range value.([]interface{}) {
				if _, ok := v.(map[string]interface{}); !ok {
					continue
				}
				secretNames = parseAttributeReplace(v.(map[string]interface{}), secretNames)
			}
		}
	}

	return secretNames
}

func getVarConfigOptions(secrets model.Secrets, variables property.Values) config.CreateOptions {
	varsConfigOpts := config.CreateOptions{
		Attributes: map[string]interface{}{},
	}

	for _, v := range secrets {
		varsConfigOpts.Attributes[_secretPrefix+v.Name] = v.Value
	}

	for k, v := range variables {
		varsConfigOpts.Attributes[_varPrefix+k] = v
	}

	return varsConfigOpts
}

func getModuleConfig(
	appMod types.ApplicationModule,
	modVer *model.ModuleVersion,
	ops CreateSecretsOptions,
) (*config.ModuleConfig, error) {
	var (
		props              = make(property.Properties, len(appMod.Attributes))
		typesWith          = appMod.Attributes.TypesWith(modVer.Schema.Variables)
		sensitiveVariables = sets.Set[string]{}
	)

	for k, v := range appMod.Attributes {
		props[k] = property.Property{
			Type:  typesWith[k],
			Value: v,
		}

		// Add sensitive from app attributes.
		val, err := v.MarshalJSON()
		if err != nil {
			return nil, err
		}

		if _secretReg.Match(val) {
			sensitiveVariables.Insert(k)
		}
	}

	attrs, err := props.TypedValues()
	if err != nil {
		return nil, err
	}

	mc := &config.ModuleConfig{
		Name:          appMod.Name,
		ModuleVersion: modVer,
		Attributes:    attrs,
	}

	if modVer.Schema == nil {
		return mc, nil
	}

	for _, v := range modVer.Schema.Variables {
		// Add sensitive from schema variable.
		if v.Sensitive {
			sensitiveVariables.Insert(v.Name)
		}

		// Add seal metadata.
		var attrValue string

		switch v.Name {
		case SealMetadataProjectName:
			attrValue = ops.ProjectName
		case SealMetadataApplicationName:
			attrValue = ops.ApplicationName
		case SealMetadataApplicationInstanceName:
			attrValue = ops.ApplicationInstanceName
		case SealMetadataModuleName:
			attrValue = appMod.Name
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
		mc.Outputs[i].ModuleName = appMod.Name
		mc.Outputs[i].Sensitive = v.Sensitive
		mc.Outputs[i].Name = v.Name

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

func matchAnyRegex(list []string) (*regexp.Regexp, error) {
	var sb strings.Builder

	sb.WriteString("(")

	for i, v := range list {
		sb.WriteString(fmt.Sprintf(`var\.%s`, v))

		if i < len(list)-1 {
			sb.WriteString("|")
		}
	}

	sb.WriteString(")")

	return regexp.Compile(sb.String())
}
