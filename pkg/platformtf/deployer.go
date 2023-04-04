package platformtf

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"regexp"

	"entgo.io/ent/dialect/sql"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"

	revisionbus "github.com/seal-io/seal/pkg/bus/applicationrevision"
	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/applicationinstance"
	"github.com/seal-io/seal/pkg/dao/model/applicationmodulerelationship"
	"github.com/seal-io/seal/pkg/dao/model/applicationresource"
	"github.com/seal-io/seal/pkg/dao/model/applicationrevision"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/model/environmentconnectorrelationship"
	"github.com/seal-io/seal/pkg/dao/model/moduleversion"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/secret"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/platform/deployer"
	"github.com/seal-io/seal/pkg/platformk8s"
	"github.com/seal-io/seal/pkg/platformtf/config"
	"github.com/seal-io/seal/pkg/platformtf/util"
	"github.com/seal-io/seal/pkg/settings"
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

// CreateSecretsOptions options for creating deployment job secrets.
type CreateSecretsOptions struct {
	SkipTLSVerify       bool
	ApplicationRevision *model.ApplicationRevision
	Connectors          model.Connectors
	ProjectID           types.ID
}

// _backendAPI the API path to terraform deploy backend.
// terraform will get and update deployment states from this API.
const _backendAPI = "/v1/application-revisions/%s/terraform-states"

// _varPrefix the prefix of the variable name.
const _varPrefix = "_seal_"

var (
	// _secretReg the regexp to match the secret variable.
	_secretReg = regexp.MustCompile(`\${secret\.([a-zA-Z0-9_]+)}`)
	// _varReg the regexp to match the variable.
	_varReg = regexp.MustCompile(`\${var\.([a-zA-Z0-9_]+)}`)
)

func NewDeployer(_ context.Context, opts deployer.CreateOptions) (deployer.Deployer, error) {
	clientSet, err := kubernetes.NewForConfig(opts.KubeConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client set: %w", err)
	}
	return &Deployer{
		modelClient: opts.ModelClient,
		clientSet:   clientSet,
		logger:      log.WithName("deployer").WithName("terraform"),
	}, nil
}

func (d Deployer) Type() deployer.Type {
	return DeployerType
}

// Apply will apply the application to deploy the application.
func (d Deployer) Apply(ctx context.Context, ai *model.ApplicationInstance, applyOpts deployer.ApplyOptions) error {
	connectors, err := d.getConnectors(ctx, ai)
	if err != nil {
		return err
	}

	applicationRevision, err := d.CreateApplicationRevision(ctx, ai)
	if err != nil {
		return err
	}

	defer func() {
		if err == nil {
			return
		}
		var errMsg = err.Error()
		// report to application revision.
		entity, rerr := d.modelClient.ApplicationRevisions().UpdateOne(applicationRevision).
			SetStatus(status.ApplicationRevisionStatusFailed).
			SetStatusMessage(errMsg).
			Save(ctx)
		if rerr != nil {
			d.logger.Errorf("failed to update application revision status: %v", rerr)
			return
		}

		if err = revisionbus.Notify(ctx, d.modelClient, entity); err != nil {
			d.logger.Errorf("add application revision update notify err: %w", err)
		}
	}()

	// get application, we need the project id to fetch available secrets.
	app, err := d.modelClient.Applications().Get(ctx, ai.ApplicationID)
	if err != nil {
		return err
	}

	// prepare tfConfig for deployment.
	secretOpts := CreateSecretsOptions{
		SkipTLSVerify:       applyOpts.SkipTLSVerify,
		ApplicationRevision: applicationRevision,
		Connectors:          connectors,
		ProjectID:           app.ProjectID,
	}
	if err = d.createK8sSecrets(ctx, secretOpts); err != nil {
		return err
	}

	jobImage, err := settings.TerraformDeployerImage.Value(ctx, d.modelClient)
	if err != nil {
		return err
	}

	// create deployment job.
	jobCreateOpts := JobCreateOptions{
		Type:                  _jobTypeApply,
		ApplicationRevisionID: applicationRevision.ID.String(),
		Image:                 jobImage,
	}
	if err = CreateJob(ctx, d.clientSet, jobCreateOpts); err != nil {
		return err
	}

	return nil
}

// Destroy will destroy the resource of the application.
// 1. get the latest revision, and checkAppRevision it if it is running.
// 2. if not running, then destroy resources.
func (d Deployer) Destroy(ctx context.Context, ai *model.ApplicationInstance, destroyOpts deployer.DestroyOptions) error {
	applicationRevision, err := d.CreateApplicationRevision(ctx, ai)
	if err != nil {
		return err
	}

	defer func() {
		if err == nil {
			return
		}
		var errMsg = err.Error()
		// report to application revision.
		applicationRevision, rerr := d.modelClient.ApplicationRevisions().UpdateOne(applicationRevision).
			SetStatus(status.ApplicationRevisionStatusFailed).
			SetStatusMessage(errMsg).
			Save(ctx)
		if rerr != nil {
			d.logger.Errorf("failed to update application revision status: %v", rerr)
			return
		}

		if err = revisionbus.Notify(ctx, d.modelClient, applicationRevision); err != nil {
			d.logger.Errorf("add application revision update notify err: %w", err)
		}
	}()

	// if no resource exists, skip job and set revision status succeed.
	exist, err := d.modelClient.ApplicationResources().Query().
		Where(applicationresource.InstanceID(ai.ID)).
		Exist(ctx)
	if err != nil {
		return err
	}
	if !exist {
		applicationRevision, err = d.modelClient.ApplicationRevisions().UpdateOne(applicationRevision).
			SetStatus(status.ApplicationRevisionStatusSucceeded).
			Save(ctx)
		if err != nil {
			return err
		}

		return revisionbus.Notify(ctx, d.modelClient, applicationRevision)
	}

	app, err := d.modelClient.Applications().Get(ctx, ai.ApplicationID)
	if err != nil {
		return err
	}

	connectors, err := d.getConnectors(ctx, ai)
	if err != nil {
		return err
	}

	// prepare tfConfig for deployment
	secretOpts := CreateSecretsOptions{
		SkipTLSVerify:       destroyOpts.SkipTLSVerify,
		ApplicationRevision: applicationRevision,
		Connectors:          connectors,
		ProjectID:           app.ProjectID,
	}
	err = d.createK8sSecrets(ctx, secretOpts)
	if err != nil {
		return err
	}

	jobImage, err := settings.TerraformDeployerImage.Value(ctx, d.modelClient)
	if err != nil {
		return err
	}

	// create deployment job
	jobOpts := JobCreateOptions{
		Type:                  _jobTypeDestroy,
		ApplicationRevisionID: applicationRevision.ID.String(),
		Image:                 jobImage,
	}
	err = CreateJob(ctx, d.clientSet, jobOpts)
	if err != nil {
		return err
	}

	return nil
}

// createK8sSecrets will create the k8s secrets for deployment.
func (d Deployer) createK8sSecrets(ctx context.Context, opts CreateSecretsOptions) error {
	var secretData = make(map[string][]byte)
	// secretName terraform tfConfig name
	secretName := _secretPrefix + string(opts.ApplicationRevision.ID)

	// prepare terraform config files bytes for deployment.
	terraformData, err := d.LoadConfigsBytes(ctx, opts)
	if err != nil {
		return err
	}
	for k, v := range terraformData {
		secretData[k] = v
	}

	// mount the provider configs(e.g. kubeconfig) to secret.
	providerData, err := d.GetProviderSecretData(opts.Connectors)
	if err != nil {
		return err
	}
	for k, v := range providerData {
		secretData[k] = v
	}

	// create deployment secret
	if err = CreateSecret(ctx, d.clientSet, secretName, secretData); err != nil {
		return err
	}

	return nil
}

// CreateApplicationRevision will create a new application revision.
// get the latest revision, and check it if it is running.
// if not running, then apply the latest revision.
// if running, then wait for the latest revision to be applied.
func (d Deployer) CreateApplicationRevision(ctx context.Context, ai *model.ApplicationInstance) (*model.ApplicationRevision, error) {
	var entity = &model.ApplicationRevision{
		DeployerType:   DeployerType,
		InstanceID:     ai.ID,
		EnvironmentID:  ai.EnvironmentID,
		InputVariables: ai.Variables,
		Status:         status.ApplicationRevisionStatusRunning,
	}

	// output of the previous revision should be inherited to the new one
	// when creating a new revision.
	var prevEntity, err = d.modelClient.ApplicationRevisions().Query().
		Order(model.Desc(applicationrevision.FieldCreateTime)).
		Where(applicationrevision.And(
			applicationrevision.InstanceID(ai.ID),
			applicationrevision.DeployerType(DeployerType))).
		Select(applicationrevision.FieldOutput, applicationrevision.FieldStatus).
		First(ctx)
	if err != nil && !model.IsNotFound(err) {
		return nil, err
	}
	if prevEntity != nil {
		if prevEntity.Status == status.ApplicationRevisionStatusRunning {
			return nil, errors.New("application deployment is running")
		}
		// inherit the output of previous revision.
		entity.Output = prevEntity.Output
		// inherit the required providers of previous succeeded revision.
		entity.PreviousRequiredProviders, err = d.getPreviousRequiredProviders(ctx, ai.ID)
		if err != nil {
			return nil, err
		}
	}

	// get modules for new revision.
	amrs, err := d.modelClient.ApplicationModuleRelationships().Query().
		Where(applicationmodulerelationship.ApplicationID(ai.ApplicationID)).
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

	// create revision, mark status to running.
	creates, err := dao.ApplicationRevisionCreates(d.modelClient, entity)
	if err != nil {
		return nil, err
	}
	return creates[0].Save(ctx)
}

// LoadConfigsBytes returns terraform main.tf and terraform.tfvars for deployment.
func (d Deployer) LoadConfigsBytes(ctx context.Context, opts CreateSecretsOptions) (map[string][]byte, error) {
	var logger = log.WithName("platformtf")
	// prepare terraform tfConfig.
	//  get module configs from app revision.
	moduleConfigs, providerRequirements, err := d.GetModuleConfigs(ctx, opts.ApplicationRevision)
	if err != nil {
		return nil, err
	}
	// merge current and previous required providers.
	providerRequirements = append(providerRequirements, opts.ApplicationRevision.PreviousRequiredProviders...)

	requiredProviders := make(map[string]*tfconfig.ProviderRequirement, 0)
	for _, p := range providerRequirements {
		if _, ok := requiredProviders[p.Name]; !ok {
			requiredProviders[p.Name] = p.ProviderRequirement
		} else {
			logger.Warnf("duplicate provider requirement: %s", p.Name)
		}
	}

	// parse module secrets
	secrets, err := d.parseModuleSecrets(ctx, moduleConfigs, opts)
	if err != nil {
		return nil, err
	}

	// prepare address for terraform backend.
	serverAddress, err := settings.ServeUrl.Value(ctx, d.modelClient)
	if err != nil {
		return nil, err
	}
	if serverAddress == "" {
		return nil, fmt.Errorf("server address is empty")
	}
	address := fmt.Sprintf("%s%s", serverAddress, fmt.Sprintf(_backendAPI, opts.ApplicationRevision.ID))
	// prepare API token for terraform backend.
	token, err := settings.PrivilegeApiToken.Value(ctx, d.modelClient)
	if err != nil {
		return nil, err
	}

	// prepare terraform config files to be mounted to secret.
	var requiredProviderNames = sets.NewString()
	for _, p := range providerRequirements {
		requiredProviderNames = requiredProviderNames.Insert(p.Name)
	}
	var secretOptionMaps = map[string]config.CreateOptions{
		config.FileMain: {
			TerraformOptions: &config.TerraformOptions{
				Token:                token,
				Address:              address,
				SkipTLSVerify:        opts.SkipTLSVerify,
				ProviderRequirements: requiredProviders,
			},
			ProviderOptions: &config.ProviderOptions{
				RequiredProviderNames: requiredProviderNames.List(),
				Connectors:            opts.Connectors,
				SecretMonthPath:       _secretMountPath,
				ConnectorSeparator:    connectorSeparator,
			},
			ModuleOptions: &config.ModuleOptions{
				ModuleConfigs: moduleConfigs,
			},
			VariableOptions: &config.VariableOptions{
				Variables: opts.ApplicationRevision.InputVariables,
				Secrets:   secrets,
				Prefix:    _varPrefix,
			},
		},
		config.FileVars: getVarConfigOptions(secrets, opts.ApplicationRevision.InputVariables),
	}
	var secretMaps = make(map[string][]byte, 0)
	for k, v := range secretOptionMaps {
		secretMaps[k], err = config.CreateConfigToBytes(v)
		if err != nil {
			return nil, err
		}
	}

	// save input plan to app revision
	revision, err := d.modelClient.ApplicationRevisions().UpdateOne(opts.ApplicationRevision).
		SetInputPlan(string(secretMaps[config.FileMain])).
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
		_, s, err := platformk8s.LoadApiConfig(*c)
		if err != nil {
			return nil, err
		}

		// NB(alex) the secret file name must be config + connector id to match with terraform provider in config convert.
		secretFileName := util.GetK8sSecretName(c.ID.String())
		secretData[secretFileName] = []byte(s)
	}
	return secretData, nil
}

// GetModuleConfigs returns module configs and required connectors to get terraform module config block from application revision.
func (d Deployer) GetModuleConfigs(ctx context.Context, ar *model.ApplicationRevision) ([]*config.ModuleConfig, []types.ProviderRequirement, error) {
	var (
		moduleConfigs     []*config.ModuleConfig
		requiredProviders = make([]types.ProviderRequirement, 0)
		// module id -> module source
		moduleVersionMap = make(map[string]*model.ModuleVersion, 0)
		predicates       = make([]predicate.ModuleVersion, 0)
	)

	for _, m := range ar.Modules {
		predicates = append(predicates, moduleversion.And(moduleversion.ModuleID(m.ModuleID), moduleversion.Version(m.Version)))
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

	modVerMapKey := func(moduleID, version string) string {
		return fmt.Sprintf("%s/%s", moduleID, version)
	}

	for _, m := range moduleVersions {
		moduleVersionMap[modVerMapKey(m.ModuleID, m.Version)] = m
	}
	for _, m := range ar.Modules {
		modVer, ok := moduleVersionMap[modVerMapKey(m.ModuleID, m.Version)]
		if !ok {
			return nil, nil, fmt.Errorf("version %s of module %s not found", m.Version, m.ModuleID)
		}

		// TODO (alex): add module config validation here.
		// verify module config attributes value type are valid.
		moduleConfig := &config.ModuleConfig{
			Name:          m.Name,
			ModuleVersion: modVer,
			Attributes:    m.Attributes,
		}
		moduleConfigs = append(moduleConfigs, moduleConfig)

		if modVer.Schema != nil {
			requiredProviders = append(requiredProviders, modVer.Schema.RequiredProviders...)
		}
	}

	return moduleConfigs, requiredProviders, err
}

func (d Deployer) getConnectors(ctx context.Context, ai *model.ApplicationInstance) (model.Connectors, error) {
	var rs, err = d.modelClient.EnvironmentConnectorRelationships().Query().
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
func (d Deployer) parseModuleSecrets(ctx context.Context, moduleConfigs []*config.ModuleConfig, opts CreateSecretsOptions) (model.Secrets, error) {
	var (
		moduleSecrets []string
		secrets       model.Secrets
	)
	for _, moduleConfig := range moduleConfigs {
		moduleSecrets = parseAttributeReplace(moduleConfig.Attributes, moduleSecrets)
	}
	nameIn := make([]interface{}, len(moduleSecrets))
	for i, name := range moduleSecrets {
		nameIn[i] = name
	}
	// this query is used to distinct the secrets with the same name.
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
	err := d.modelClient.Secrets().Query().
		Modify(func(s *sql.Selector) {
			// select secrets without project id or not in project.
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
	if err != nil {
		return nil, err
	}

	// validate module secret are all exist.
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
func (d Deployer) getPreviousRequiredProviders(ctx context.Context, instanceID types.ID) ([]types.ProviderRequirement, error) {
	var (
		prevRequiredProviders = make([]types.ProviderRequirement, 0)
		predicates            []predicate.ModuleVersion
	)

	var entity, err = d.modelClient.ApplicationRevisions().Query().
		Where(applicationrevision.InstanceID(instanceID)).
		Where(applicationrevision.Status(status.ApplicationRevisionStatusSucceeded)).
		Order(model.Desc(applicationrevision.FieldCreateTime)).
		First(ctx)
	if err != nil && !model.IsNotFound(err) {
		return nil, err
	}
	for _, m := range entity.Modules {
		predicates = append(predicates, moduleversion.And(
			moduleversion.ModuleID(m.ModuleID),
			moduleversion.Version(m.Version)))
	}

	if len(predicates) != 0 {
		mvs, err := d.modelClient.ModuleVersions().Query().
			Where(moduleversion.And(predicates...)).
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

	return prevRequiredProviders, nil
}

func SyncApplicationRevisionStatus(ctx context.Context, bm revisionbus.BusMessage) error {
	var (
		mc       = bm.TransactionalModelClient
		revision = bm.Refer
	)

	// report to application instance.
	appInstance, err := mc.ApplicationInstances().Query().
		Where(applicationinstance.ID(revision.InstanceID)).
		Select(
			applicationinstance.FieldID,
			applicationinstance.FieldStatus).
		Only(ctx)
	if err != nil {
		return err
	}
	switch revision.Status {
	case status.ApplicationRevisionStatusSucceeded:
		if appInstance.Status == status.ApplicationInstanceStatusDeleting {
			// delete application instance.
			err = mc.ApplicationInstances().DeleteOne(appInstance).
				Exec(ctx)
		} else {
			err = mc.ApplicationInstances().UpdateOne(appInstance).
				SetStatus(status.ApplicationInstanceStatusDeployed).
				Exec(ctx)
		}
	case status.ApplicationRevisionStatusFailed:
		if appInstance.Status == status.ApplicationInstanceStatusDeleting {
			appInstance.Status = status.ApplicationInstanceStatusDeleteFailed
		} else {
			appInstance.Status = status.ApplicationInstanceStatusDeployFailed
		}
		err = mc.ApplicationInstances().UpdateOne(appInstance).
			SetStatus(appInstance.Status).
			SetStatusMessage(revision.StatusMessage).
			Exec(ctx)
	}

	return err
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
			repl := "${var." + _varPrefix + "${1}}"
			str = _varReg.ReplaceAllString(str, repl)
			attributes[key] = _secretReg.ReplaceAllString(str, repl)
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

func getVarConfigOptions(secrets model.Secrets, variables map[string]interface{}) config.CreateOptions {
	varsConfigOpts := config.CreateOptions{
		Attributes: map[string]interface{}{},
	}

	for _, v := range secrets {
		varsConfigOpts.Attributes[_varPrefix+v.Name] = v.Value
	}

	for k, v := range variables {
		varsConfigOpts.Attributes[_varPrefix+k] = v
	}

	return varsConfigOpts
}
