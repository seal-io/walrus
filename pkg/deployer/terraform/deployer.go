package terraform

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"

	apiconfig "github.com/seal-io/walrus/pkg/apis/config"
	"github.com/seal-io/walrus/pkg/auths/session"
	revisionbus "github.com/seal-io/walrus/pkg/bus/resourcerevision"
	"github.com/seal-io/walrus/pkg/dao"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/environment"
	"github.com/seal-io/walrus/pkg/dao/model/predicate"
	"github.com/seal-io/walrus/pkg/dao/model/project"
	"github.com/seal-io/walrus/pkg/dao/model/resource"
	"github.com/seal-io/walrus/pkg/dao/model/resourcedefinition"
	"github.com/seal-io/walrus/pkg/dao/model/resourcedefinitionmatchingrule"
	"github.com/seal-io/walrus/pkg/dao/model/resourcerevision"
	"github.com/seal-io/walrus/pkg/dao/model/template"
	"github.com/seal-io/walrus/pkg/dao/model/templateversion"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/crypto"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/dao/types/property"
	"github.com/seal-io/walrus/pkg/dao/types/status"
	deptypes "github.com/seal-io/walrus/pkg/deployer/types"
	pkgenv "github.com/seal-io/walrus/pkg/environment"
	opk8s "github.com/seal-io/walrus/pkg/operator/k8s"
	pkgresource "github.com/seal-io/walrus/pkg/resource"
	"github.com/seal-io/walrus/pkg/resourcedefinitions"
	"github.com/seal-io/walrus/pkg/settings"
	"github.com/seal-io/walrus/pkg/terraform/config"
	"github.com/seal-io/walrus/pkg/terraform/parser"
	"github.com/seal-io/walrus/pkg/terraform/util"
	"github.com/seal-io/walrus/utils/log"
)

var (
	// _variableReg the regexp to match the variable.
	_variableReg = regexp.MustCompile(`\${var\.([a-zA-Z0-9_-]+)}`)
	// _resourceReg the regexp to match the service/resource output.
	_resourceReg = regexp.MustCompile(`\${(svc|res)\.([^.}]+)\.([^.}]+)}`)
	// _interpolationReg is the regular expression for matching non-reference or non-variable expressions.
	// Reference: https://developer.hashicorp.com/terraform/language/expressions/strings#escape-sequences-1
	// To handle escape sequences, ${xxx} is converted to $${xxx}.
	// If there are more than two consecutive $ symbols, like $${xxx}, they are further converted to $$${xxx}.
	// During Terraform processing, $${} is ultimately transformed back to ${},
	// this interpolation is used to ensuring a WYSIWYG user experience.
	_interpolationReg = regexp.MustCompile(`\$\{((var\.)?([^.}]+)(?:\.([^.}]+))?)[^\}]*\}`)
)

// Deployer terraform deployer to deploy the resource.
type Deployer struct {
	logger      log.Logger
	modelClient model.ClientSet
	clientSet   *kubernetes.Clientset
}

type createRevisionOptions struct {
	// ResourceID indicates the ID of resource which is for create the revision.
	ResourceID object.ID
	// JobType indicates the type of the job.
	JobType string
}

type createK8sJobOptions struct {
	// Type indicates the type of the job.
	Type string
	// ResourceRevision indicates the resource revision to create the deployment job.
	ResourceRevision *model.ResourceRevision
}

type createK8sSecretsOptions struct {
	ResourceRevision *model.ResourceRevision
	Connectors       model.Connectors
	SubjectID        object.ID
	// Walrus Context.
	Context Context

	Token     string
	ServerURL string
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

// Apply creates a new resource revision by the given resource,
// and drives the Kubernetes Job to create components of the resource.
func (d Deployer) Apply(ctx context.Context, resource *model.Resource, opts deptypes.ApplyOptions) (err error) {
	revision, err := d.createRevision(ctx, createRevisionOptions{
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
		_ = d.updateRevisionStatus(ctx, revision)
	}()

	err = d.createK8sJob(ctx, createK8sJobOptions{
		Type:             JobTypeApply,
		ResourceRevision: revision,
	})

	return err
}

// Destroy creates a new resource revision by the given resource,
// and drives the Kubernetes Job to clean the components of the resource.
func (d Deployer) Destroy(ctx context.Context, resource *model.Resource, opts deptypes.DestroyOptions) (err error) {
	revision, err := d.createRevision(ctx, createRevisionOptions{
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
		_ = d.updateRevisionStatus(ctx, revision)
	}()

	err = d.createK8sJob(ctx, createK8sJobOptions{
		Type:             JobTypeDestroy,
		ResourceRevision: revision,
	})

	return err
}

// Sync creates a new resource revision by the given resource,
// and drives the Kubernetes Job to sync the components of the resource.
func (d Deployer) Sync(ctx context.Context, resource *model.Resource, opts deptypes.SyncOptions) (err error) {
	revision, err := d.createRevision(ctx, createRevisionOptions{
		ResourceID: resource.ID,
		JobType:    JobTypeSync,
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
		_ = d.updateRevisionStatus(ctx, revision)
	}()

	err = d.createK8sJob(ctx, createK8sJobOptions{
		Type:             JobTypeSync,
		ResourceRevision: revision,
	})

	return err
}

// Detect will detect resource changes from remote system of given service.
func (d Deployer) Detect(
	ctx context.Context,
	resource *model.Resource,
	opts deptypes.DetectOptions,
) error {
	revision, err := d.createRevision(ctx, createRevisionOptions{
		ResourceID: resource.ID,
		JobType:    JobTypeDetect,
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
		_ = d.updateRevisionStatus(ctx, revision)
	}()

	err = d.createK8sJob(ctx, createK8sJobOptions{
		Type:             JobTypeDetect,
		ResourceRevision: revision,
	})

	return err
}

// createK8sJob creates a k8s job to deploy, destroy or rollback the resource.
func (d Deployer) createK8sJob(ctx context.Context, opts createK8sJobOptions) error {
	revision := opts.ResourceRevision

	connectors, err := dao.GetEnvironmentConnectors(ctx, d.modelClient, revision.EnvironmentID)
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

	res, err := d.modelClient.Resources().Get(ctx, revision.ResourceID)
	if err != nil {
		return err
	}

	var subjectID object.ID

	sj, _ := session.GetSubject(ctx)
	if sj.ID != "" {
		subjectID = sj.ID
	} else {
		subjectID, err = pkgresource.GetSubjectID(res)
		if err != nil {
			return err
		}
	}

	if subjectID == "" {
		return errors.New("subject id is empty")
	}

	// Prepare tfConfig for deployment.
	secretOpts := createK8sSecretsOptions{
		ResourceRevision: opts.ResourceRevision,
		Connectors:       connectors,
		SubjectID:        subjectID,
		// Walrus Context.
		Context: *NewContext().
			SetProject(proj.ID, proj.Name).
			SetEnvironment(env.ID, env.Name, pkgenv.GetManagedNamespaceName(env)).
			SetResource(res.ID, res.Name),
	}

	secretOpts.ServerURL, err = getServerURL(ctx, d.modelClient)
	if err != nil {
		return err
	}

	secretOpts.Token, err = getToken(ctx, d.modelClient, secretOpts.SubjectID, opts.ResourceRevision.ID)
	if err != nil {
		return err
	}

	if err = d.createK8sSecrets(ctx, secretOpts); err != nil {
		return err
	}

	jobImage, err := settings.DeployerImage.Value(ctx, d.modelClient)
	if err != nil {
		return err
	}

	jobEnv, err := d.getEnv(ctx)
	if err != nil {
		return err
	}

	// Create deployment job.
	jobOpts := JobCreateOptions{
		Type:      opts.Type,
		Image:     jobImage,
		Env:       jobEnv,
		ServerURL: secretOpts.ServerURL,
		Token:     secretOpts.Token,
	}

	return createJob(ctx, d.clientSet, opts.ResourceRevision, jobOpts)
}

func (d Deployer) getEnv(ctx context.Context) ([]corev1.EnvVar, error) {
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

	if settings.SkipRemoteTLSVerify.ShouldValueBool(ctx, d.modelClient) {
		env = append(env, corev1.EnvVar{
			Name:  "GIT_SSL_NO_VERIFY",
			Value: "true",
		})
	}

	return env, nil
}

func (d Deployer) updateRevisionStatus(ctx context.Context, ar *model.ResourceRevision) error {
	// Report to resource revision.
	ar.Status.SetSummary(status.WalkResourceRevision(&ar.Status))

	ar, err := d.modelClient.ResourceRevisions().UpdateOne(ar).
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

// createK8sSecrets creates the k8s secrets for deployment.
func (d Deployer) createK8sSecrets(ctx context.Context, opts createK8sSecretsOptions) error {
	secretData := make(map[string][]byte)
	// SecretName terraform tfConfig name.
	secretName := _jobSecretPrefix + string(opts.ResourceRevision.ID)

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

// createRevision creates a new resource revision.
// Get the latest revision, and check it if it is running.
// If not running, then apply the latest revision.
// If running, then wait for the latest revision to be applied.
func (d Deployer) createRevision(
	ctx context.Context,
	opts createRevisionOptions,
) (*model.ResourceRevision, error) {
	// Validate if there is a running revision.
	prevEntity, err := d.modelClient.ResourceRevisions().Query().
		Where(resourcerevision.And(
			resourcerevision.ResourceID(opts.ResourceID),
			resourcerevision.DeployerType(DeployerType))).
		Order(model.Desc(resourcerevision.FieldCreateTime)).
		First(ctx)
	if err != nil && !model.IsNotFound(err) {
		return nil, err
	}

	if prevEntity != nil && status.ResourceRevisionStatusReady.IsUnknown(prevEntity) {
		return nil, errors.New("deployment is running")
	}

	// Get the corresponding resource and template version.
	res, err := d.modelClient.Resources().Query().
		Where(resource.ID(opts.ResourceID)).
		WithTemplate(func(tvq *model.TemplateVersionQuery) {
			tvq.Select(
				templateversion.FieldName,
				templateversion.FieldVersion,
				templateversion.FieldTemplateID)
		}).
		WithProject(func(pq *model.ProjectQuery) {
			pq.Select(project.FieldName)
		}).
		WithEnvironment(func(env *model.EnvironmentQuery) {
			env.Select(environment.FieldLabels)
			env.Select(environment.FieldName)
			env.Select(environment.FieldType)
		}).
		WithResourceDefinition(func(rd *model.ResourceDefinitionQuery) {
			rd.Select(resourcedefinition.FieldType)
			rd.WithMatchingRules(func(mrq *model.ResourceDefinitionMatchingRuleQuery) {
				mrq.Order(model.Asc(resourcedefinitionmatchingrule.FieldOrder)).
					Select(
						resourcedefinitionmatchingrule.FieldName,
						resourcedefinitionmatchingrule.FieldSelector,
						resourcedefinitionmatchingrule.FieldAttributes,
					).
					WithTemplate(func(tvq *model.TemplateVersionQuery) {
						tvq.Select(
							templateversion.FieldID,
							templateversion.FieldVersion,
							templateversion.FieldName,
						)
					})
			})
		}).
		Select(
			resource.FieldID,
			resource.FieldProjectID,
			resource.FieldEnvironmentID,
			resource.FieldType,
			resource.FieldLabels,
			resource.FieldAttributes).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	var (
		templateID                    object.ID
		templateName, templateVersion string
		attributes                    property.Values
	)

	switch {
	case res.TemplateID != nil:
		templateID = res.Edges.Template.TemplateID
		templateName = res.Edges.Template.Name
		templateVersion = res.Edges.Template.Version
		attributes = res.Attributes
	case res.ResourceDefinitionID != nil:
		rd := res.Edges.ResourceDefinition
		matchRule := resourcedefinitions.Match(
			rd.Edges.MatchingRules,
			res.Edges.Project.Name,
			res.Edges.Environment.Name,
			res.Edges.Environment.Type,
			res.Edges.Environment.Labels,
			res.Labels,
		)

		if matchRule == nil {
			return nil, fmt.Errorf("resource definition %s does not match resource %s", rd.Name, res.Name)
		}
		templateName = matchRule.Edges.Template.Name
		templateVersion = matchRule.Edges.Template.Version

		templateID, err = d.modelClient.Templates().Query().
			Where(
				template.Name(templateName),
				// Now we only support resource definition globally.
				template.ProjectIDIsNil(),
			).
			OnlyID(ctx)
		if err != nil {
			return nil, err
		}

		// Merge attributes. Resource attributes take precedence over resource definition attributes.
		attributes = matchRule.Attributes
		if attributes == nil {
			attributes = make(property.Values)
		}

		for k, v := range res.Attributes {
			attributes[k] = v
		}
	default:
		return nil, errors.New("service has no template or resource definition")
	}

	entity := &model.ResourceRevision{
		Type:            opts.JobType,
		ProjectID:       res.ProjectID,
		EnvironmentID:   res.EnvironmentID,
		ResourceID:      res.ID,
		TemplateID:      templateID,
		TemplateName:    templateName,
		TemplateVersion: templateVersion,
		Attributes:      attributes,
		DeployerType:    DeployerType,
	}

	status.ResourceRevisionStatusReady.Unknown(entity, "")
	entity.Status.SetSummary(status.WalkResourceRevision(&entity.Status))

	// Inherit the output of previous revision to create a new one.
	if prevEntity != nil {
		entity.Output = prevEntity.Output
	}

	switch {
	case opts.JobType == JobTypeApply && entity.Output != "":
		// Get required providers from the previous output after first deployment.
		requiredProviders, err := d.getRequiredProviders(ctx, opts.ResourceID, entity.Output)
		if err != nil {
			return nil, err
		}
		entity.PreviousRequiredProviders = requiredProviders
	case opts.JobType == JobTypeDestroy && entity.Output != "":
		if status.ResourceRevisionStatusReady.IsFalse(prevEntity) {
			// Get required providers from the previous output after first deployment.
			requiredProviders, err := d.getRequiredProviders(ctx, opts.ResourceID, entity.Output)
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
	entity, err = d.modelClient.ResourceRevisions().Create().
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
	// Prepare terraform tfConfig.
	//  get module configs from resource revision.
	moduleConfig, providerRequirements, err := d.getModuleConfig(ctx, opts)
	if err != nil {
		return nil, err
	}
	// Merge current and previous required providers.
	providerRequirements = append(providerRequirements,
		opts.ResourceRevision.PreviousRequiredProviders...)

	requiredProviders := make(map[string]*tfconfig.ProviderRequirement, 0)
	for _, p := range providerRequirements {
		if _, ok := requiredProviders[p.Name]; !ok {
			requiredProviders[p.Name] = p.ProviderRequirement
		} else {
			d.logger.Warnf("duplicate provider requirement: %s", p.Name)
		}
	}

	revisionOpts := RevisionOpts{
		ResourceRevision: opts.ResourceRevision,
		ResourceName:     opts.Context.Resource.Name,
		ProjectID:        opts.Context.Project.ID,
		EnvironmentID:    opts.Context.Environment.ID,
	}
	// Parse module attributes.
	attrs, variables, dependencyOutputs, err := ParseModuleAttributes(
		ctx,
		d.modelClient,
		moduleConfig.Attributes,
		false,
		revisionOpts,
	)
	if err != nil {
		return nil, err
	}

	moduleConfig.Attributes = attrs

	// Update output sensitive with variables.
	wrapVariables, err := updateOutputWithVariables(variables, moduleConfig)
	if err != nil {
		return nil, err
	}

	if opts.ServerURL == "" {
		return nil, errors.New("server address is empty")
	}
	address := fmt.Sprintf("%s%s", opts.ServerURL,
		fmt.Sprintf(_backendAPI,
			opts.Context.Project.ID,
			opts.Context.Environment.ID,
			opts.Context.Resource.ID,
			opts.ResourceRevision.ID))

	// Prepare terraform config files to be mounted to secret.
	requiredProviderNames := sets.NewString()
	for _, p := range providerRequirements {
		requiredProviderNames = requiredProviderNames.Insert(p.Name)
	}

	secretOptionMaps := map[string]config.CreateOptions{
		config.FileMain: {
			TerraformOptions: &config.TerraformOptions{
				Token:                opts.Token,
				Address:              address,
				SkipTLSVerify:        !apiconfig.TlsCertified.Get(),
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
				ResourcePrefix:    _resourcePrefix,
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

	// Save input plan to resource revision.
	opts.ResourceRevision.InputPlan = string(secretMaps[config.FileMain])
	// If resource revision does not inherit variables from cloned revision,
	// then save the parsed variables to resource revision.
	if len(opts.ResourceRevision.Variables) == 0 {
		variableMap := make(crypto.Map[string, string], len(variables))
		for _, s := range variables {
			variableMap[s.Name] = string(s.Value)
		}
		opts.ResourceRevision.Variables = variableMap
	}

	revision, err := d.modelClient.ResourceRevisions().UpdateOne(opts.ResourceRevision).
		Set(opts.ResourceRevision).
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
		if c.Type != types.ConnectorTypeKubernetes {
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
// get terraform module config block from resource revision.
func (d Deployer) getModuleConfig(
	ctx context.Context,
	opts createK8sSecretsOptions,
) (*config.ModuleConfig, []types.ProviderRequirement, error) {
	var (
		requiredProviders = make([]types.ProviderRequirement, 0)
		predicates        = make([]predicate.TemplateVersion, 0)
	)

	predicates = append(predicates, templateversion.And(
		templateversion.Version(opts.ResourceRevision.TemplateVersion),
		templateversion.TemplateID(opts.ResourceRevision.TemplateID),
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
			templateversion.FieldUiSchema,
		).
		Where(templateversion.Or(predicates...)).
		Only(ctx)
	if err != nil {
		return nil, nil, err
	}

	if len(templateVersion.Schema.RequiredProviders) != 0 {
		requiredProviders = append(requiredProviders, templateVersion.Schema.RequiredProviders...)
	}

	mc, err := getModuleConfig(opts.ResourceRevision, templateVersion, opts)
	if err != nil {
		return nil, nil, err
	}

	return mc, requiredProviders, err
}

// getPreviousRequiredProviders get previous succeed revision required providers.
// NB(alex): the previous revision may be failed, the failed revision may not contain required providers of states.
func (d Deployer) getPreviousRequiredProviders(
	ctx context.Context,
	resourceID object.ID,
) ([]types.ProviderRequirement, error) {
	prevRequiredProviders := make([]types.ProviderRequirement, 0)

	entity, err := d.modelClient.ResourceRevisions().Query().
		Where(resourcerevision.ResourceID(resourceID)).
		Order(model.Desc(resourcerevision.FieldCreateTime)).
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

	if len(templateVersion.Schema.RequiredProviders) != 0 {
		prevRequiredProviders = append(prevRequiredProviders, templateVersion.Schema.RequiredProviders...)
	}

	return prevRequiredProviders, nil
}
