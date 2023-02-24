// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package runtime

import (
	"time"

	"github.com/seal-io/seal/pkg/dao/model/application"
	"github.com/seal-io/seal/pkg/dao/model/applicationmodulerelationship"
	"github.com/seal-io/seal/pkg/dao/model/applicationresource"
	"github.com/seal-io/seal/pkg/dao/model/applicationrevision"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/model/environmentconnectorrelationship"
	"github.com/seal-io/seal/pkg/dao/model/module"
	"github.com/seal-io/seal/pkg/dao/model/project"
	"github.com/seal-io/seal/pkg/dao/model/role"
	"github.com/seal-io/seal/pkg/dao/model/setting"
	"github.com/seal-io/seal/pkg/dao/model/subject"
	"github.com/seal-io/seal/pkg/dao/model/token"
	"github.com/seal-io/seal/pkg/dao/schema"
	"github.com/seal-io/seal/pkg/dao/types"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	applicationMixin := schema.Application{}.Mixin()
	applicationMixinHooks0 := applicationMixin[0].Hooks()
	application.Hooks[0] = applicationMixinHooks0[0]
	applicationMixinFields1 := applicationMixin[1].Fields()
	_ = applicationMixinFields1
	applicationMixinFields2 := applicationMixin[2].Fields()
	_ = applicationMixinFields2
	applicationFields := schema.Application{}.Fields()
	_ = applicationFields
	// applicationDescName is the schema descriptor for name field.
	applicationDescName := applicationMixinFields1[0].Descriptor()
	// application.NameValidator is a validator for the "name" field. It is called by the builders before save.
	application.NameValidator = applicationDescName.Validators[0].(func(string) error)
	// applicationDescLabels is the schema descriptor for labels field.
	applicationDescLabels := applicationMixinFields1[2].Descriptor()
	// application.DefaultLabels holds the default value on creation for the labels field.
	application.DefaultLabels = applicationDescLabels.Default.(map[string]string)
	// applicationDescCreateTime is the schema descriptor for createTime field.
	applicationDescCreateTime := applicationMixinFields2[0].Descriptor()
	// application.DefaultCreateTime holds the default value on creation for the createTime field.
	application.DefaultCreateTime = applicationDescCreateTime.Default.(func() time.Time)
	// applicationDescUpdateTime is the schema descriptor for updateTime field.
	applicationDescUpdateTime := applicationMixinFields2[1].Descriptor()
	// application.DefaultUpdateTime holds the default value on creation for the updateTime field.
	application.DefaultUpdateTime = applicationDescUpdateTime.Default.(func() time.Time)
	// application.UpdateDefaultUpdateTime holds the default value on update for the updateTime field.
	application.UpdateDefaultUpdateTime = applicationDescUpdateTime.UpdateDefault.(func() time.Time)
	// applicationDescProjectID is the schema descriptor for projectID field.
	applicationDescProjectID := applicationFields[0].Descriptor()
	// application.ProjectIDValidator is a validator for the "projectID" field. It is called by the builders before save.
	application.ProjectIDValidator = applicationDescProjectID.Validators[0].(func(string) error)
	// applicationDescEnvironmentID is the schema descriptor for environmentID field.
	applicationDescEnvironmentID := applicationFields[1].Descriptor()
	// application.EnvironmentIDValidator is a validator for the "environmentID" field. It is called by the builders before save.
	application.EnvironmentIDValidator = applicationDescEnvironmentID.Validators[0].(func(string) error)
	applicationmodulerelationshipMixin := schema.ApplicationModuleRelationship{}.Mixin()
	applicationmodulerelationshipMixinFields0 := applicationmodulerelationshipMixin[0].Fields()
	_ = applicationmodulerelationshipMixinFields0
	applicationmodulerelationshipFields := schema.ApplicationModuleRelationship{}.Fields()
	_ = applicationmodulerelationshipFields
	// applicationmodulerelationshipDescCreateTime is the schema descriptor for createTime field.
	applicationmodulerelationshipDescCreateTime := applicationmodulerelationshipMixinFields0[0].Descriptor()
	// applicationmodulerelationship.DefaultCreateTime holds the default value on creation for the createTime field.
	applicationmodulerelationship.DefaultCreateTime = applicationmodulerelationshipDescCreateTime.Default.(func() time.Time)
	// applicationmodulerelationshipDescUpdateTime is the schema descriptor for updateTime field.
	applicationmodulerelationshipDescUpdateTime := applicationmodulerelationshipMixinFields0[1].Descriptor()
	// applicationmodulerelationship.DefaultUpdateTime holds the default value on creation for the updateTime field.
	applicationmodulerelationship.DefaultUpdateTime = applicationmodulerelationshipDescUpdateTime.Default.(func() time.Time)
	// applicationmodulerelationship.UpdateDefaultUpdateTime holds the default value on update for the updateTime field.
	applicationmodulerelationship.UpdateDefaultUpdateTime = applicationmodulerelationshipDescUpdateTime.UpdateDefault.(func() time.Time)
	// applicationmodulerelationshipDescApplicationID is the schema descriptor for application_id field.
	applicationmodulerelationshipDescApplicationID := applicationmodulerelationshipFields[0].Descriptor()
	// applicationmodulerelationship.ApplicationIDValidator is a validator for the "application_id" field. It is called by the builders before save.
	applicationmodulerelationship.ApplicationIDValidator = applicationmodulerelationshipDescApplicationID.Validators[0].(func(string) error)
	// applicationmodulerelationshipDescModuleID is the schema descriptor for module_id field.
	applicationmodulerelationshipDescModuleID := applicationmodulerelationshipFields[1].Descriptor()
	// applicationmodulerelationship.ModuleIDValidator is a validator for the "module_id" field. It is called by the builders before save.
	applicationmodulerelationship.ModuleIDValidator = applicationmodulerelationshipDescModuleID.Validators[0].(func(string) error)
	// applicationmodulerelationshipDescName is the schema descriptor for name field.
	applicationmodulerelationshipDescName := applicationmodulerelationshipFields[2].Descriptor()
	// applicationmodulerelationship.NameValidator is a validator for the "name" field. It is called by the builders before save.
	applicationmodulerelationship.NameValidator = applicationmodulerelationshipDescName.Validators[0].(func(string) error)
	applicationresourceMixin := schema.ApplicationResource{}.Mixin()
	applicationresourceMixinHooks0 := applicationresourceMixin[0].Hooks()
	applicationresource.Hooks[0] = applicationresourceMixinHooks0[0]
	applicationresourceMixinFields2 := applicationresourceMixin[2].Fields()
	_ = applicationresourceMixinFields2
	applicationresourceFields := schema.ApplicationResource{}.Fields()
	_ = applicationresourceFields
	// applicationresourceDescCreateTime is the schema descriptor for createTime field.
	applicationresourceDescCreateTime := applicationresourceMixinFields2[0].Descriptor()
	// applicationresource.DefaultCreateTime holds the default value on creation for the createTime field.
	applicationresource.DefaultCreateTime = applicationresourceDescCreateTime.Default.(func() time.Time)
	// applicationresourceDescUpdateTime is the schema descriptor for updateTime field.
	applicationresourceDescUpdateTime := applicationresourceMixinFields2[1].Descriptor()
	// applicationresource.DefaultUpdateTime holds the default value on creation for the updateTime field.
	applicationresource.DefaultUpdateTime = applicationresourceDescUpdateTime.Default.(func() time.Time)
	// applicationresource.UpdateDefaultUpdateTime holds the default value on update for the updateTime field.
	applicationresource.UpdateDefaultUpdateTime = applicationresourceDescUpdateTime.UpdateDefault.(func() time.Time)
	// applicationresourceDescApplicationID is the schema descriptor for applicationID field.
	applicationresourceDescApplicationID := applicationresourceFields[0].Descriptor()
	// applicationresource.ApplicationIDValidator is a validator for the "applicationID" field. It is called by the builders before save.
	applicationresource.ApplicationIDValidator = applicationresourceDescApplicationID.Validators[0].(func(string) error)
	// applicationresourceDescConnectorID is the schema descriptor for connectorID field.
	applicationresourceDescConnectorID := applicationresourceFields[1].Descriptor()
	// applicationresource.ConnectorIDValidator is a validator for the "connectorID" field. It is called by the builders before save.
	applicationresource.ConnectorIDValidator = applicationresourceDescConnectorID.Validators[0].(func(string) error)
	// applicationresourceDescModule is the schema descriptor for module field.
	applicationresourceDescModule := applicationresourceFields[2].Descriptor()
	// applicationresource.ModuleValidator is a validator for the "module" field. It is called by the builders before save.
	applicationresource.ModuleValidator = applicationresourceDescModule.Validators[0].(func(string) error)
	// applicationresourceDescMode is the schema descriptor for mode field.
	applicationresourceDescMode := applicationresourceFields[3].Descriptor()
	// applicationresource.ModeValidator is a validator for the "mode" field. It is called by the builders before save.
	applicationresource.ModeValidator = applicationresourceDescMode.Validators[0].(func(string) error)
	// applicationresourceDescType is the schema descriptor for type field.
	applicationresourceDescType := applicationresourceFields[4].Descriptor()
	// applicationresource.TypeValidator is a validator for the "type" field. It is called by the builders before save.
	applicationresource.TypeValidator = applicationresourceDescType.Validators[0].(func(string) error)
	// applicationresourceDescName is the schema descriptor for name field.
	applicationresourceDescName := applicationresourceFields[5].Descriptor()
	// applicationresource.NameValidator is a validator for the "name" field. It is called by the builders before save.
	applicationresource.NameValidator = applicationresourceDescName.Validators[0].(func(string) error)
	applicationrevisionMixin := schema.ApplicationRevision{}.Mixin()
	applicationrevisionMixinHooks0 := applicationrevisionMixin[0].Hooks()
	applicationrevision.Hooks[0] = applicationrevisionMixinHooks0[0]
	applicationrevisionMixinFields2 := applicationrevisionMixin[2].Fields()
	_ = applicationrevisionMixinFields2
	applicationrevisionFields := schema.ApplicationRevision{}.Fields()
	_ = applicationrevisionFields
	// applicationrevisionDescCreateTime is the schema descriptor for createTime field.
	applicationrevisionDescCreateTime := applicationrevisionMixinFields2[0].Descriptor()
	// applicationrevision.DefaultCreateTime holds the default value on creation for the createTime field.
	applicationrevision.DefaultCreateTime = applicationrevisionDescCreateTime.Default.(func() time.Time)
	// applicationrevisionDescUpdateTime is the schema descriptor for updateTime field.
	applicationrevisionDescUpdateTime := applicationrevisionMixinFields2[1].Descriptor()
	// applicationrevision.DefaultUpdateTime holds the default value on creation for the updateTime field.
	applicationrevision.DefaultUpdateTime = applicationrevisionDescUpdateTime.Default.(func() time.Time)
	// applicationrevision.UpdateDefaultUpdateTime holds the default value on update for the updateTime field.
	applicationrevision.UpdateDefaultUpdateTime = applicationrevisionDescUpdateTime.UpdateDefault.(func() time.Time)
	// applicationrevisionDescApplicationID is the schema descriptor for applicationID field.
	applicationrevisionDescApplicationID := applicationrevisionFields[0].Descriptor()
	// applicationrevision.ApplicationIDValidator is a validator for the "applicationID" field. It is called by the builders before save.
	applicationrevision.ApplicationIDValidator = applicationrevisionDescApplicationID.Validators[0].(func(string) error)
	// applicationrevisionDescEnvironmentID is the schema descriptor for environmentID field.
	applicationrevisionDescEnvironmentID := applicationrevisionFields[1].Descriptor()
	// applicationrevision.EnvironmentIDValidator is a validator for the "environmentID" field. It is called by the builders before save.
	applicationrevision.EnvironmentIDValidator = applicationrevisionDescEnvironmentID.Validators[0].(func(string) error)
	// applicationrevisionDescModules is the schema descriptor for modules field.
	applicationrevisionDescModules := applicationrevisionFields[2].Descriptor()
	// applicationrevision.DefaultModules holds the default value on creation for the modules field.
	applicationrevision.DefaultModules = applicationrevisionDescModules.Default.([]types.ApplicationModule)
	// applicationrevisionDescInputVariables is the schema descriptor for inputVariables field.
	applicationrevisionDescInputVariables := applicationrevisionFields[3].Descriptor()
	// applicationrevision.DefaultInputVariables holds the default value on creation for the inputVariables field.
	applicationrevision.DefaultInputVariables = applicationrevisionDescInputVariables.Default.(map[string]interface{})
	connectorMixin := schema.Connector{}.Mixin()
	connectorMixinHooks0 := connectorMixin[0].Hooks()
	connector.Hooks[0] = connectorMixinHooks0[0]
	connectorMixinFields1 := connectorMixin[1].Fields()
	_ = connectorMixinFields1
	connectorMixinFields3 := connectorMixin[3].Fields()
	_ = connectorMixinFields3
	connectorFields := schema.Connector{}.Fields()
	_ = connectorFields
	// connectorDescName is the schema descriptor for name field.
	connectorDescName := connectorMixinFields1[0].Descriptor()
	// connector.NameValidator is a validator for the "name" field. It is called by the builders before save.
	connector.NameValidator = connectorDescName.Validators[0].(func(string) error)
	// connectorDescLabels is the schema descriptor for labels field.
	connectorDescLabels := connectorMixinFields1[2].Descriptor()
	// connector.DefaultLabels holds the default value on creation for the labels field.
	connector.DefaultLabels = connectorDescLabels.Default.(map[string]string)
	// connectorDescCreateTime is the schema descriptor for createTime field.
	connectorDescCreateTime := connectorMixinFields3[0].Descriptor()
	// connector.DefaultCreateTime holds the default value on creation for the createTime field.
	connector.DefaultCreateTime = connectorDescCreateTime.Default.(func() time.Time)
	// connectorDescUpdateTime is the schema descriptor for updateTime field.
	connectorDescUpdateTime := connectorMixinFields3[1].Descriptor()
	// connector.DefaultUpdateTime holds the default value on creation for the updateTime field.
	connector.DefaultUpdateTime = connectorDescUpdateTime.Default.(func() time.Time)
	// connector.UpdateDefaultUpdateTime holds the default value on update for the updateTime field.
	connector.UpdateDefaultUpdateTime = connectorDescUpdateTime.UpdateDefault.(func() time.Time)
	// connectorDescType is the schema descriptor for type field.
	connectorDescType := connectorFields[0].Descriptor()
	// connector.TypeValidator is a validator for the "type" field. It is called by the builders before save.
	connector.TypeValidator = connectorDescType.Validators[0].(func(string) error)
	// connectorDescConfigVersion is the schema descriptor for configVersion field.
	connectorDescConfigVersion := connectorFields[1].Descriptor()
	// connector.ConfigVersionValidator is a validator for the "configVersion" field. It is called by the builders before save.
	connector.ConfigVersionValidator = connectorDescConfigVersion.Validators[0].(func(string) error)
	// connectorDescConfigData is the schema descriptor for configData field.
	connectorDescConfigData := connectorFields[2].Descriptor()
	// connector.DefaultConfigData holds the default value on creation for the configData field.
	connector.DefaultConfigData = connectorDescConfigData.Default.(map[string]interface{})
	environmentMixin := schema.Environment{}.Mixin()
	environmentMixinHooks0 := environmentMixin[0].Hooks()
	environment.Hooks[0] = environmentMixinHooks0[0]
	environmentMixinFields1 := environmentMixin[1].Fields()
	_ = environmentMixinFields1
	environmentMixinFields2 := environmentMixin[2].Fields()
	_ = environmentMixinFields2
	environmentFields := schema.Environment{}.Fields()
	_ = environmentFields
	// environmentDescName is the schema descriptor for name field.
	environmentDescName := environmentMixinFields1[0].Descriptor()
	// environment.NameValidator is a validator for the "name" field. It is called by the builders before save.
	environment.NameValidator = environmentDescName.Validators[0].(func(string) error)
	// environmentDescLabels is the schema descriptor for labels field.
	environmentDescLabels := environmentMixinFields1[2].Descriptor()
	// environment.DefaultLabels holds the default value on creation for the labels field.
	environment.DefaultLabels = environmentDescLabels.Default.(map[string]string)
	// environmentDescCreateTime is the schema descriptor for createTime field.
	environmentDescCreateTime := environmentMixinFields2[0].Descriptor()
	// environment.DefaultCreateTime holds the default value on creation for the createTime field.
	environment.DefaultCreateTime = environmentDescCreateTime.Default.(func() time.Time)
	// environmentDescUpdateTime is the schema descriptor for updateTime field.
	environmentDescUpdateTime := environmentMixinFields2[1].Descriptor()
	// environment.DefaultUpdateTime holds the default value on creation for the updateTime field.
	environment.DefaultUpdateTime = environmentDescUpdateTime.Default.(func() time.Time)
	// environment.UpdateDefaultUpdateTime holds the default value on update for the updateTime field.
	environment.UpdateDefaultUpdateTime = environmentDescUpdateTime.UpdateDefault.(func() time.Time)
	environmentconnectorrelationshipMixin := schema.EnvironmentConnectorRelationship{}.Mixin()
	environmentconnectorrelationshipMixinFields0 := environmentconnectorrelationshipMixin[0].Fields()
	_ = environmentconnectorrelationshipMixinFields0
	environmentconnectorrelationshipFields := schema.EnvironmentConnectorRelationship{}.Fields()
	_ = environmentconnectorrelationshipFields
	// environmentconnectorrelationshipDescCreateTime is the schema descriptor for createTime field.
	environmentconnectorrelationshipDescCreateTime := environmentconnectorrelationshipMixinFields0[0].Descriptor()
	// environmentconnectorrelationship.DefaultCreateTime holds the default value on creation for the createTime field.
	environmentconnectorrelationship.DefaultCreateTime = environmentconnectorrelationshipDescCreateTime.Default.(func() time.Time)
	// environmentconnectorrelationshipDescEnvironmentID is the schema descriptor for environment_id field.
	environmentconnectorrelationshipDescEnvironmentID := environmentconnectorrelationshipFields[0].Descriptor()
	// environmentconnectorrelationship.EnvironmentIDValidator is a validator for the "environment_id" field. It is called by the builders before save.
	environmentconnectorrelationship.EnvironmentIDValidator = environmentconnectorrelationshipDescEnvironmentID.Validators[0].(func(string) error)
	// environmentconnectorrelationshipDescConnectorID is the schema descriptor for connector_id field.
	environmentconnectorrelationshipDescConnectorID := environmentconnectorrelationshipFields[1].Descriptor()
	// environmentconnectorrelationship.ConnectorIDValidator is a validator for the "connector_id" field. It is called by the builders before save.
	environmentconnectorrelationship.ConnectorIDValidator = environmentconnectorrelationshipDescConnectorID.Validators[0].(func(string) error)
	moduleMixin := schema.Module{}.Mixin()
	moduleMixinFields1 := moduleMixin[1].Fields()
	_ = moduleMixinFields1
	moduleFields := schema.Module{}.Fields()
	_ = moduleFields
	// moduleDescCreateTime is the schema descriptor for createTime field.
	moduleDescCreateTime := moduleMixinFields1[0].Descriptor()
	// module.DefaultCreateTime holds the default value on creation for the createTime field.
	module.DefaultCreateTime = moduleDescCreateTime.Default.(func() time.Time)
	// moduleDescUpdateTime is the schema descriptor for updateTime field.
	moduleDescUpdateTime := moduleMixinFields1[1].Descriptor()
	// module.DefaultUpdateTime holds the default value on creation for the updateTime field.
	module.DefaultUpdateTime = moduleDescUpdateTime.Default.(func() time.Time)
	// module.UpdateDefaultUpdateTime holds the default value on update for the updateTime field.
	module.UpdateDefaultUpdateTime = moduleDescUpdateTime.UpdateDefault.(func() time.Time)
	// moduleDescLabels is the schema descriptor for labels field.
	moduleDescLabels := moduleFields[2].Descriptor()
	// module.DefaultLabels holds the default value on creation for the labels field.
	module.DefaultLabels = moduleDescLabels.Default.(map[string]string)
	// moduleDescSource is the schema descriptor for source field.
	moduleDescSource := moduleFields[3].Descriptor()
	// module.SourceValidator is a validator for the "source" field. It is called by the builders before save.
	module.SourceValidator = moduleDescSource.Validators[0].(func(string) error)
	// moduleDescSchema is the schema descriptor for schema field.
	moduleDescSchema := moduleFields[5].Descriptor()
	// module.DefaultSchema holds the default value on creation for the schema field.
	module.DefaultSchema = moduleDescSchema.Default.(*types.ModuleSchema)
	// moduleDescID is the schema descriptor for id field.
	moduleDescID := moduleFields[0].Descriptor()
	// module.IDValidator is a validator for the "id" field. It is called by the builders before save.
	module.IDValidator = moduleDescID.Validators[0].(func(string) error)
	projectMixin := schema.Project{}.Mixin()
	projectMixinHooks0 := projectMixin[0].Hooks()
	project.Hooks[0] = projectMixinHooks0[0]
	projectMixinFields1 := projectMixin[1].Fields()
	_ = projectMixinFields1
	projectMixinFields2 := projectMixin[2].Fields()
	_ = projectMixinFields2
	projectFields := schema.Project{}.Fields()
	_ = projectFields
	// projectDescName is the schema descriptor for name field.
	projectDescName := projectMixinFields1[0].Descriptor()
	// project.NameValidator is a validator for the "name" field. It is called by the builders before save.
	project.NameValidator = projectDescName.Validators[0].(func(string) error)
	// projectDescLabels is the schema descriptor for labels field.
	projectDescLabels := projectMixinFields1[2].Descriptor()
	// project.DefaultLabels holds the default value on creation for the labels field.
	project.DefaultLabels = projectDescLabels.Default.(map[string]string)
	// projectDescCreateTime is the schema descriptor for createTime field.
	projectDescCreateTime := projectMixinFields2[0].Descriptor()
	// project.DefaultCreateTime holds the default value on creation for the createTime field.
	project.DefaultCreateTime = projectDescCreateTime.Default.(func() time.Time)
	// projectDescUpdateTime is the schema descriptor for updateTime field.
	projectDescUpdateTime := projectMixinFields2[1].Descriptor()
	// project.DefaultUpdateTime holds the default value on creation for the updateTime field.
	project.DefaultUpdateTime = projectDescUpdateTime.Default.(func() time.Time)
	// project.UpdateDefaultUpdateTime holds the default value on update for the updateTime field.
	project.UpdateDefaultUpdateTime = projectDescUpdateTime.UpdateDefault.(func() time.Time)
	roleMixin := schema.Role{}.Mixin()
	roleMixinHooks0 := roleMixin[0].Hooks()
	role.Hooks[0] = roleMixinHooks0[0]
	roleMixinFields1 := roleMixin[1].Fields()
	_ = roleMixinFields1
	roleFields := schema.Role{}.Fields()
	_ = roleFields
	// roleDescCreateTime is the schema descriptor for createTime field.
	roleDescCreateTime := roleMixinFields1[0].Descriptor()
	// role.DefaultCreateTime holds the default value on creation for the createTime field.
	role.DefaultCreateTime = roleDescCreateTime.Default.(func() time.Time)
	// roleDescUpdateTime is the schema descriptor for updateTime field.
	roleDescUpdateTime := roleMixinFields1[1].Descriptor()
	// role.DefaultUpdateTime holds the default value on creation for the updateTime field.
	role.DefaultUpdateTime = roleDescUpdateTime.Default.(func() time.Time)
	// role.UpdateDefaultUpdateTime holds the default value on update for the updateTime field.
	role.UpdateDefaultUpdateTime = roleDescUpdateTime.UpdateDefault.(func() time.Time)
	// roleDescDomain is the schema descriptor for domain field.
	roleDescDomain := roleFields[0].Descriptor()
	// role.DefaultDomain holds the default value on creation for the domain field.
	role.DefaultDomain = roleDescDomain.Default.(string)
	// roleDescName is the schema descriptor for name field.
	roleDescName := roleFields[1].Descriptor()
	// role.NameValidator is a validator for the "name" field. It is called by the builders before save.
	role.NameValidator = roleDescName.Validators[0].(func(string) error)
	// roleDescPolicies is the schema descriptor for policies field.
	roleDescPolicies := roleFields[3].Descriptor()
	// role.DefaultPolicies holds the default value on creation for the policies field.
	role.DefaultPolicies = roleDescPolicies.Default.(schema.RolePolicies)
	// roleDescBuiltin is the schema descriptor for builtin field.
	roleDescBuiltin := roleFields[4].Descriptor()
	// role.DefaultBuiltin holds the default value on creation for the builtin field.
	role.DefaultBuiltin = roleDescBuiltin.Default.(bool)
	// roleDescSession is the schema descriptor for session field.
	roleDescSession := roleFields[5].Descriptor()
	// role.DefaultSession holds the default value on creation for the session field.
	role.DefaultSession = roleDescSession.Default.(bool)
	settingMixin := schema.Setting{}.Mixin()
	settingMixinHooks0 := settingMixin[0].Hooks()
	setting.Hooks[0] = settingMixinHooks0[0]
	settingMixinFields1 := settingMixin[1].Fields()
	_ = settingMixinFields1
	settingFields := schema.Setting{}.Fields()
	_ = settingFields
	// settingDescCreateTime is the schema descriptor for createTime field.
	settingDescCreateTime := settingMixinFields1[0].Descriptor()
	// setting.DefaultCreateTime holds the default value on creation for the createTime field.
	setting.DefaultCreateTime = settingDescCreateTime.Default.(func() time.Time)
	// settingDescUpdateTime is the schema descriptor for updateTime field.
	settingDescUpdateTime := settingMixinFields1[1].Descriptor()
	// setting.DefaultUpdateTime holds the default value on creation for the updateTime field.
	setting.DefaultUpdateTime = settingDescUpdateTime.Default.(func() time.Time)
	// setting.UpdateDefaultUpdateTime holds the default value on update for the updateTime field.
	setting.UpdateDefaultUpdateTime = settingDescUpdateTime.UpdateDefault.(func() time.Time)
	// settingDescName is the schema descriptor for name field.
	settingDescName := settingFields[0].Descriptor()
	// setting.NameValidator is a validator for the "name" field. It is called by the builders before save.
	setting.NameValidator = settingDescName.Validators[0].(func(string) error)
	// settingDescHidden is the schema descriptor for hidden field.
	settingDescHidden := settingFields[2].Descriptor()
	// setting.DefaultHidden holds the default value on creation for the hidden field.
	setting.DefaultHidden = settingDescHidden.Default.(bool)
	// settingDescEditable is the schema descriptor for editable field.
	settingDescEditable := settingFields[3].Descriptor()
	// setting.DefaultEditable holds the default value on creation for the editable field.
	setting.DefaultEditable = settingDescEditable.Default.(bool)
	// settingDescPrivate is the schema descriptor for private field.
	settingDescPrivate := settingFields[4].Descriptor()
	// setting.DefaultPrivate holds the default value on creation for the private field.
	setting.DefaultPrivate = settingDescPrivate.Default.(bool)
	subjectMixin := schema.Subject{}.Mixin()
	subjectMixinHooks0 := subjectMixin[0].Hooks()
	subject.Hooks[0] = subjectMixinHooks0[0]
	subjectMixinFields1 := subjectMixin[1].Fields()
	_ = subjectMixinFields1
	subjectFields := schema.Subject{}.Fields()
	_ = subjectFields
	// subjectDescCreateTime is the schema descriptor for createTime field.
	subjectDescCreateTime := subjectMixinFields1[0].Descriptor()
	// subject.DefaultCreateTime holds the default value on creation for the createTime field.
	subject.DefaultCreateTime = subjectDescCreateTime.Default.(func() time.Time)
	// subjectDescUpdateTime is the schema descriptor for updateTime field.
	subjectDescUpdateTime := subjectMixinFields1[1].Descriptor()
	// subject.DefaultUpdateTime holds the default value on creation for the updateTime field.
	subject.DefaultUpdateTime = subjectDescUpdateTime.Default.(func() time.Time)
	// subject.UpdateDefaultUpdateTime holds the default value on update for the updateTime field.
	subject.UpdateDefaultUpdateTime = subjectDescUpdateTime.UpdateDefault.(func() time.Time)
	// subjectDescKind is the schema descriptor for kind field.
	subjectDescKind := subjectFields[0].Descriptor()
	// subject.DefaultKind holds the default value on creation for the kind field.
	subject.DefaultKind = subjectDescKind.Default.(string)
	// subjectDescGroup is the schema descriptor for group field.
	subjectDescGroup := subjectFields[1].Descriptor()
	// subject.DefaultGroup holds the default value on creation for the group field.
	subject.DefaultGroup = subjectDescGroup.Default.(string)
	// subjectDescName is the schema descriptor for name field.
	subjectDescName := subjectFields[2].Descriptor()
	// subject.NameValidator is a validator for the "name" field. It is called by the builders before save.
	subject.NameValidator = subjectDescName.Validators[0].(func(string) error)
	// subjectDescMountTo is the schema descriptor for mountTo field.
	subjectDescMountTo := subjectFields[4].Descriptor()
	// subject.DefaultMountTo holds the default value on creation for the mountTo field.
	subject.DefaultMountTo = subjectDescMountTo.Default.(bool)
	// subjectDescLoginTo is the schema descriptor for loginTo field.
	subjectDescLoginTo := subjectFields[5].Descriptor()
	// subject.DefaultLoginTo holds the default value on creation for the loginTo field.
	subject.DefaultLoginTo = subjectDescLoginTo.Default.(bool)
	// subjectDescRoles is the schema descriptor for roles field.
	subjectDescRoles := subjectFields[6].Descriptor()
	// subject.DefaultRoles holds the default value on creation for the roles field.
	subject.DefaultRoles = subjectDescRoles.Default.(schema.SubjectRoles)
	// subjectDescPaths is the schema descriptor for paths field.
	subjectDescPaths := subjectFields[7].Descriptor()
	// subject.DefaultPaths holds the default value on creation for the paths field.
	subject.DefaultPaths = subjectDescPaths.Default.([]string)
	// subjectDescBuiltin is the schema descriptor for builtin field.
	subjectDescBuiltin := subjectFields[8].Descriptor()
	// subject.DefaultBuiltin holds the default value on creation for the builtin field.
	subject.DefaultBuiltin = subjectDescBuiltin.Default.(bool)
	tokenMixin := schema.Token{}.Mixin()
	tokenMixinHooks0 := tokenMixin[0].Hooks()
	token.Hooks[0] = tokenMixinHooks0[0]
	tokenMixinFields1 := tokenMixin[1].Fields()
	_ = tokenMixinFields1
	tokenFields := schema.Token{}.Fields()
	_ = tokenFields
	// tokenDescCreateTime is the schema descriptor for createTime field.
	tokenDescCreateTime := tokenMixinFields1[0].Descriptor()
	// token.DefaultCreateTime holds the default value on creation for the createTime field.
	token.DefaultCreateTime = tokenDescCreateTime.Default.(func() time.Time)
	// tokenDescUpdateTime is the schema descriptor for updateTime field.
	tokenDescUpdateTime := tokenMixinFields1[1].Descriptor()
	// token.DefaultUpdateTime holds the default value on creation for the updateTime field.
	token.DefaultUpdateTime = tokenDescUpdateTime.Default.(func() time.Time)
	// token.UpdateDefaultUpdateTime holds the default value on update for the updateTime field.
	token.UpdateDefaultUpdateTime = tokenDescUpdateTime.UpdateDefault.(func() time.Time)
	// tokenDescCasdoorTokenName is the schema descriptor for casdoorTokenName field.
	tokenDescCasdoorTokenName := tokenFields[0].Descriptor()
	// token.CasdoorTokenNameValidator is a validator for the "casdoorTokenName" field. It is called by the builders before save.
	token.CasdoorTokenNameValidator = tokenDescCasdoorTokenName.Validators[0].(func(string) error)
	// tokenDescCasdoorTokenOwner is the schema descriptor for casdoorTokenOwner field.
	tokenDescCasdoorTokenOwner := tokenFields[1].Descriptor()
	// token.CasdoorTokenOwnerValidator is a validator for the "casdoorTokenOwner" field. It is called by the builders before save.
	token.CasdoorTokenOwnerValidator = tokenDescCasdoorTokenOwner.Validators[0].(func(string) error)
	// tokenDescName is the schema descriptor for name field.
	tokenDescName := tokenFields[2].Descriptor()
	// token.NameValidator is a validator for the "name" field. It is called by the builders before save.
	token.NameValidator = tokenDescName.Validators[0].(func(string) error)
}

const (
	Version = "v0.11.7"                                         // Version of ent codegen.
	Sum     = "h1:V+wKFh0jhAbY/FoU+PPbdMOf2Ma5vh07R/IdF+N/nFg=" // Sum of ent codegen.
)
