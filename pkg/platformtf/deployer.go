package platformtf

import (
	"context"
	"errors"
	"fmt"
	"io"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"

	revisionbus "github.com/seal-io/seal/pkg/bus/applicationrevision"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/applicationinstance"
	"github.com/seal-io/seal/pkg/dao/model/applicationmodulerelationship"
	"github.com/seal-io/seal/pkg/dao/model/applicationrevision"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/model/environmentconnectorrelationship"
	"github.com/seal-io/seal/pkg/dao/model/moduleversion"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/platform/deployer"
	"github.com/seal-io/seal/pkg/platformk8s"
	"github.com/seal-io/seal/pkg/platformtf/config"
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
}

// _backendAPI the API path to terraform deploy backend.
// terraform will get and update deployment states from this API.
const _backendAPI = "/v1/application-revisions/%s/terraform-states"

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

	// prepare tfConfig for deployment.
	secretOpts := CreateSecretsOptions{
		SkipTLSVerify:       applyOpts.SkipTLSVerify,
		ApplicationRevision: applicationRevision,
		Connectors:          connectors,
	}
	if err = d.createSecrets(ctx, secretOpts); err != nil {
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

	connectors, err := d.getConnectors(ctx, ai)
	if err != nil {
		return err
	}

	// prepare tfConfig for deployment
	secretOpts := CreateSecretsOptions{
		SkipTLSVerify:       destroyOpts.SkipTLSVerify,
		ApplicationRevision: applicationRevision,
		Connectors:          connectors,
	}
	err = d.createSecrets(ctx, secretOpts)
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

// createSecrets will create the k8s secrets for deployment.
func (d Deployer) createSecrets(ctx context.Context, opts CreateSecretsOptions) error {
	// secretName terraform tfConfig name
	secretName := _secretPrefix + string(opts.ApplicationRevision.ID)
	// prepare tfConfig for deployment
	tfConfig, err := d.LoadConfig(ctx, opts)
	if err != nil {
		return err
	}

	// secretData the data of terraform tfConfig.
	secretData := map[string][]byte{
		"main.tf": []byte(tfConfig),
	}
	// mount the provider tfConfig to secret.
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
	// output of the previous revision should be inherited to the new one
	// when creating a new revision.
	var prevOutput string
	applicationRevision, err := d.modelClient.ApplicationRevisions().Query().
		Order(model.Desc(applicationrevision.FieldCreateTime)).
		Where(applicationrevision.InstanceID(ai.ID)).
		First(ctx)
	if err != nil && !model.IsNotFound(err) {
		return nil, err
	}
	ok, err := d.checkRevisionStatus(applicationRevision)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("application deployment is running")
	}
	if applicationRevision != nil {
		prevOutput = applicationRevision.Output
	}

	// create new revision with modules.
	var modules []types.ApplicationModule
	amrs, err := d.modelClient.ApplicationModuleRelationships().Query().
		Where(applicationmodulerelationship.ApplicationID(ai.ApplicationID)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	for _, amr := range amrs {
		modules = append(modules, types.ApplicationModule{
			ModuleID:   amr.ModuleID,
			Version:    amr.Version,
			Name:       amr.Name,
			Attributes: amr.Attributes,
		})
	}

	newRevision, err := d.modelClient.ApplicationRevisions().Create().
		SetInstanceID(ai.ID).
		SetEnvironmentID(ai.EnvironmentID).
		SetModules(modules).
		SetInputVariables(ai.Variables).
		SetInputPlan("").
		SetOutput(prevOutput).
		SetStatus(status.ApplicationRevisionStatusRunning).
		SetDeployerType(DeployerType).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return newRevision, nil
}

// checkRevisionStatus check the revision status.
func (d Deployer) checkRevisionStatus(applicationRevision *model.ApplicationRevision) (bool, error) {
	if applicationRevision != nil && applicationRevision.Status == status.ApplicationRevisionStatusRunning {
		return false, fmt.Errorf("the deployment is running, please wait for it to be applied")
	}
	return true, nil
}

// LoadConfig returns terraform config for deployment.
func (d Deployer) LoadConfig(ctx context.Context, opts CreateSecretsOptions) (string, error) {
	if opts.ApplicationRevision.InputPlan != "" {
		return opts.ApplicationRevision.InputPlan, nil
	}

	// prepare terraform tfConfig.
	//  get module configs from app revision.
	moduleConfigs, requiredConnTypes, err := d.GetModuleConfigs(ctx, opts.ApplicationRevision)
	if err != nil {
		return "", err
	}

	// prepare address for terraform backend.
	serverAddress, err := settings.ServeUrl.Value(ctx, d.modelClient)
	if err != nil {
		return "", err
	}
	if serverAddress == "" {
		return "", fmt.Errorf("server address is empty")
	}
	address := fmt.Sprintf("%s%s", serverAddress, fmt.Sprintf(_backendAPI, opts.ApplicationRevision.ID))
	// prepare API token for terraform backend.
	token, err := settings.PrivilegeApiToken.Value(ctx, d.modelClient)
	if err != nil {
		return "", err
	}

	configOpts := config.CreateOptions{
		SecretMountPath:    _secretMountPath,
		ConnectorSeparator: connectorSeparator,
		Address:            address,
		Token:              token,
		SkipTLSVerify:      opts.SkipTLSVerify,
		Connectors:         opts.Connectors,
		ModuleConfigs:      moduleConfigs,
		RequiredProviders:  requiredConnTypes,
	}
	conf, err := config.NewConfig(configOpts)
	if err != nil {
		return "", err
	}
	tfReader, err := conf.Reader()
	if err != nil {
		return "", err
	}
	tfConfigBytes, err := io.ReadAll(tfReader)
	if err != nil {
		return "", err
	}

	if err != nil {
		return "", err
	}
	// save input plan to app revision
	revision, err := d.modelClient.ApplicationRevisions().UpdateOne(opts.ApplicationRevision).
		SetInputPlan(string(tfConfigBytes)).
		Save(ctx)
	if err != nil {
		return "", err
	}

	if err = revisionbus.Notify(ctx, d.modelClient, revision); err != nil {
		return "", err
	}

	return string(tfConfigBytes), nil
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
		secretFileName := config.GetSecretK8sConfigName(c.ID.String())
		secretData[secretFileName] = []byte(s)
	}
	return secretData, nil
}

// GetModuleConfigs returns module configs and required connectors to get terraform module config block from application revision.
func (d Deployer) GetModuleConfigs(ctx context.Context, ar *model.ApplicationRevision) ([]*config.ModuleConfig, []string, error) {
	var (
		requiredConnectorSet = sets.NewString()
		moduleConfigs        []*config.ModuleConfig
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
		Where(predicates...).
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

		moduleConfig := &config.ModuleConfig{
			Name:          m.Name,
			ModuleVersion: modVer,
			Attributes:    m.Attributes,
		}
		moduleConfigs = append(moduleConfigs, moduleConfig)

		if modVer.Schema != nil {
			requiredConnectorSet.Insert(modVer.Schema.RequiredConnectorTypes...)
		}
	}

	return moduleConfigs, requiredConnectorSet.List(), err
}

func (d Deployer) getConnectors(ctx context.Context, ai *model.ApplicationInstance) (model.Connectors, error) {
	var rs, err = d.modelClient.EnvironmentConnectorRelationships().Query().
		Where(environmentconnectorrelationship.EnvironmentID(ai.EnvironmentID)).
		WithConnector(func(cq *model.ConnectorQuery) {
			cq.Select(
				connector.FieldID,
				connector.FieldName,
				connector.FieldType,
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

func SyncApplicationRevisionStatus(ctx context.Context, message revisionbus.BusMessage) error {
	var (
		mc       = message.ModelClient
		revision = message.Refer
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
