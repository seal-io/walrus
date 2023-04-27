// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package runtime

import (
	"time"

	"github.com/seal-io/seal/pkg/dao/model/allocationcost"
	"github.com/seal-io/seal/pkg/dao/model/application"
	"github.com/seal-io/seal/pkg/dao/model/applicationinstance"
	"github.com/seal-io/seal/pkg/dao/model/applicationmodulerelationship"
	"github.com/seal-io/seal/pkg/dao/model/applicationresource"
	"github.com/seal-io/seal/pkg/dao/model/applicationrevision"
	"github.com/seal-io/seal/pkg/dao/model/clustercost"
	"github.com/seal-io/seal/pkg/dao/model/connector"
	"github.com/seal-io/seal/pkg/dao/model/environment"
	"github.com/seal-io/seal/pkg/dao/model/environmentconnectorrelationship"
	"github.com/seal-io/seal/pkg/dao/model/module"
	"github.com/seal-io/seal/pkg/dao/model/moduleversion"
	"github.com/seal-io/seal/pkg/dao/model/perspective"
	"github.com/seal-io/seal/pkg/dao/model/project"
	"github.com/seal-io/seal/pkg/dao/model/role"
	"github.com/seal-io/seal/pkg/dao/model/secret"
	"github.com/seal-io/seal/pkg/dao/model/setting"
	"github.com/seal-io/seal/pkg/dao/model/subject"
	"github.com/seal-io/seal/pkg/dao/model/token"
	"github.com/seal-io/seal/pkg/dao/schema"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/crypto"
	"github.com/seal-io/seal/pkg/dao/types/property"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	allocationcostFields := schema.AllocationCost{}.Fields()
	_ = allocationcostFields
	// allocationcostDescConnectorID is the schema descriptor for connectorID field.
	allocationcostDescConnectorID := allocationcostFields[3].Descriptor()
	// allocationcost.ConnectorIDValidator is a validator for the "connectorID" field. It is called by the builders before save.
	allocationcost.ConnectorIDValidator = allocationcostDescConnectorID.Validators[0].(func(string) error)
	// allocationcostDescPvs is the schema descriptor for pvs field.
	allocationcostDescPvs := allocationcostFields[13].Descriptor()
	// allocationcost.DefaultPvs holds the default value on creation for the pvs field.
	allocationcost.DefaultPvs = allocationcostDescPvs.Default.(map[string]types.PVCost)
	// allocationcostDescLabels is the schema descriptor for labels field.
	allocationcostDescLabels := allocationcostFields[14].Descriptor()
	// allocationcost.DefaultLabels holds the default value on creation for the labels field.
	allocationcost.DefaultLabels = allocationcostDescLabels.Default.(map[string]string)
	// allocationcostDescTotalCost is the schema descriptor for totalCost field.
	allocationcostDescTotalCost := allocationcostFields[15].Descriptor()
	// allocationcost.DefaultTotalCost holds the default value on creation for the totalCost field.
	allocationcost.DefaultTotalCost = allocationcostDescTotalCost.Default.(float64)
	// allocationcost.TotalCostValidator is a validator for the "totalCost" field. It is called by the builders before save.
	allocationcost.TotalCostValidator = allocationcostDescTotalCost.Validators[0].(func(float64) error)
	// allocationcostDescCpuCost is the schema descriptor for cpuCost field.
	allocationcostDescCpuCost := allocationcostFields[17].Descriptor()
	// allocationcost.DefaultCpuCost holds the default value on creation for the cpuCost field.
	allocationcost.DefaultCpuCost = allocationcostDescCpuCost.Default.(float64)
	// allocationcost.CpuCostValidator is a validator for the "cpuCost" field. It is called by the builders before save.
	allocationcost.CpuCostValidator = allocationcostDescCpuCost.Validators[0].(func(float64) error)
	// allocationcostDescCpuCoreRequest is the schema descriptor for cpuCoreRequest field.
	allocationcostDescCpuCoreRequest := allocationcostFields[18].Descriptor()
	// allocationcost.DefaultCpuCoreRequest holds the default value on creation for the cpuCoreRequest field.
	allocationcost.DefaultCpuCoreRequest = allocationcostDescCpuCoreRequest.Default.(float64)
	// allocationcost.CpuCoreRequestValidator is a validator for the "cpuCoreRequest" field. It is called by the builders before save.
	allocationcost.CpuCoreRequestValidator = allocationcostDescCpuCoreRequest.Validators[0].(func(float64) error)
	// allocationcostDescGpuCost is the schema descriptor for gpuCost field.
	allocationcostDescGpuCost := allocationcostFields[19].Descriptor()
	// allocationcost.DefaultGpuCost holds the default value on creation for the gpuCost field.
	allocationcost.DefaultGpuCost = allocationcostDescGpuCost.Default.(float64)
	// allocationcost.GpuCostValidator is a validator for the "gpuCost" field. It is called by the builders before save.
	allocationcost.GpuCostValidator = allocationcostDescGpuCost.Validators[0].(func(float64) error)
	// allocationcostDescGpuCount is the schema descriptor for gpuCount field.
	allocationcostDescGpuCount := allocationcostFields[20].Descriptor()
	// allocationcost.DefaultGpuCount holds the default value on creation for the gpuCount field.
	allocationcost.DefaultGpuCount = allocationcostDescGpuCount.Default.(float64)
	// allocationcost.GpuCountValidator is a validator for the "gpuCount" field. It is called by the builders before save.
	allocationcost.GpuCountValidator = allocationcostDescGpuCount.Validators[0].(func(float64) error)
	// allocationcostDescRamCost is the schema descriptor for ramCost field.
	allocationcostDescRamCost := allocationcostFields[21].Descriptor()
	// allocationcost.DefaultRamCost holds the default value on creation for the ramCost field.
	allocationcost.DefaultRamCost = allocationcostDescRamCost.Default.(float64)
	// allocationcost.RamCostValidator is a validator for the "ramCost" field. It is called by the builders before save.
	allocationcost.RamCostValidator = allocationcostDescRamCost.Validators[0].(func(float64) error)
	// allocationcostDescRamByteRequest is the schema descriptor for ramByteRequest field.
	allocationcostDescRamByteRequest := allocationcostFields[22].Descriptor()
	// allocationcost.DefaultRamByteRequest holds the default value on creation for the ramByteRequest field.
	allocationcost.DefaultRamByteRequest = allocationcostDescRamByteRequest.Default.(float64)
	// allocationcost.RamByteRequestValidator is a validator for the "ramByteRequest" field. It is called by the builders before save.
	allocationcost.RamByteRequestValidator = allocationcostDescRamByteRequest.Validators[0].(func(float64) error)
	// allocationcostDescPvCost is the schema descriptor for pvCost field.
	allocationcostDescPvCost := allocationcostFields[23].Descriptor()
	// allocationcost.DefaultPvCost holds the default value on creation for the pvCost field.
	allocationcost.DefaultPvCost = allocationcostDescPvCost.Default.(float64)
	// allocationcost.PvCostValidator is a validator for the "pvCost" field. It is called by the builders before save.
	allocationcost.PvCostValidator = allocationcostDescPvCost.Validators[0].(func(float64) error)
	// allocationcostDescPvBytes is the schema descriptor for pvBytes field.
	allocationcostDescPvBytes := allocationcostFields[24].Descriptor()
	// allocationcost.DefaultPvBytes holds the default value on creation for the pvBytes field.
	allocationcost.DefaultPvBytes = allocationcostDescPvBytes.Default.(float64)
	// allocationcost.PvBytesValidator is a validator for the "pvBytes" field. It is called by the builders before save.
	allocationcost.PvBytesValidator = allocationcostDescPvBytes.Validators[0].(func(float64) error)
	// allocationcostDescLoadBalancerCost is the schema descriptor for loadBalancerCost field.
	allocationcostDescLoadBalancerCost := allocationcostFields[25].Descriptor()
	// allocationcost.DefaultLoadBalancerCost holds the default value on creation for the loadBalancerCost field.
	allocationcost.DefaultLoadBalancerCost = allocationcostDescLoadBalancerCost.Default.(float64)
	// allocationcost.LoadBalancerCostValidator is a validator for the "loadBalancerCost" field. It is called by the builders before save.
	allocationcost.LoadBalancerCostValidator = allocationcostDescLoadBalancerCost.Validators[0].(func(float64) error)
	// allocationcostDescCpuCoreUsageAverage is the schema descriptor for cpuCoreUsageAverage field.
	allocationcostDescCpuCoreUsageAverage := allocationcostFields[26].Descriptor()
	// allocationcost.DefaultCpuCoreUsageAverage holds the default value on creation for the cpuCoreUsageAverage field.
	allocationcost.DefaultCpuCoreUsageAverage = allocationcostDescCpuCoreUsageAverage.Default.(float64)
	// allocationcost.CpuCoreUsageAverageValidator is a validator for the "cpuCoreUsageAverage" field. It is called by the builders before save.
	allocationcost.CpuCoreUsageAverageValidator = allocationcostDescCpuCoreUsageAverage.Validators[0].(func(float64) error)
	// allocationcostDescCpuCoreUsageMax is the schema descriptor for cpuCoreUsageMax field.
	allocationcostDescCpuCoreUsageMax := allocationcostFields[27].Descriptor()
	// allocationcost.DefaultCpuCoreUsageMax holds the default value on creation for the cpuCoreUsageMax field.
	allocationcost.DefaultCpuCoreUsageMax = allocationcostDescCpuCoreUsageMax.Default.(float64)
	// allocationcost.CpuCoreUsageMaxValidator is a validator for the "cpuCoreUsageMax" field. It is called by the builders before save.
	allocationcost.CpuCoreUsageMaxValidator = allocationcostDescCpuCoreUsageMax.Validators[0].(func(float64) error)
	// allocationcostDescRamByteUsageAverage is the schema descriptor for ramByteUsageAverage field.
	allocationcostDescRamByteUsageAverage := allocationcostFields[28].Descriptor()
	// allocationcost.DefaultRamByteUsageAverage holds the default value on creation for the ramByteUsageAverage field.
	allocationcost.DefaultRamByteUsageAverage = allocationcostDescRamByteUsageAverage.Default.(float64)
	// allocationcost.RamByteUsageAverageValidator is a validator for the "ramByteUsageAverage" field. It is called by the builders before save.
	allocationcost.RamByteUsageAverageValidator = allocationcostDescRamByteUsageAverage.Validators[0].(func(float64) error)
	// allocationcostDescRamByteUsageMax is the schema descriptor for ramByteUsageMax field.
	allocationcostDescRamByteUsageMax := allocationcostFields[29].Descriptor()
	// allocationcost.DefaultRamByteUsageMax holds the default value on creation for the ramByteUsageMax field.
	allocationcost.DefaultRamByteUsageMax = allocationcostDescRamByteUsageMax.Default.(float64)
	// allocationcost.RamByteUsageMaxValidator is a validator for the "ramByteUsageMax" field. It is called by the builders before save.
	allocationcost.RamByteUsageMaxValidator = allocationcostDescRamByteUsageMax.Validators[0].(func(float64) error)
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
	applicationinstanceMixin := schema.ApplicationInstance{}.Mixin()
	applicationinstanceMixinHooks0 := applicationinstanceMixin[0].Hooks()
	applicationinstance.Hooks[0] = applicationinstanceMixinHooks0[0]
	applicationinstanceMixinFields1 := applicationinstanceMixin[1].Fields()
	_ = applicationinstanceMixinFields1
	applicationinstanceFields := schema.ApplicationInstance{}.Fields()
	_ = applicationinstanceFields
	// applicationinstanceDescCreateTime is the schema descriptor for createTime field.
	applicationinstanceDescCreateTime := applicationinstanceMixinFields1[0].Descriptor()
	// applicationinstance.DefaultCreateTime holds the default value on creation for the createTime field.
	applicationinstance.DefaultCreateTime = applicationinstanceDescCreateTime.Default.(func() time.Time)
	// applicationinstanceDescUpdateTime is the schema descriptor for updateTime field.
	applicationinstanceDescUpdateTime := applicationinstanceMixinFields1[1].Descriptor()
	// applicationinstance.DefaultUpdateTime holds the default value on creation for the updateTime field.
	applicationinstance.DefaultUpdateTime = applicationinstanceDescUpdateTime.Default.(func() time.Time)
	// applicationinstance.UpdateDefaultUpdateTime holds the default value on update for the updateTime field.
	applicationinstance.UpdateDefaultUpdateTime = applicationinstanceDescUpdateTime.UpdateDefault.(func() time.Time)
	// applicationinstanceDescApplicationID is the schema descriptor for applicationID field.
	applicationinstanceDescApplicationID := applicationinstanceFields[0].Descriptor()
	// applicationinstance.ApplicationIDValidator is a validator for the "applicationID" field. It is called by the builders before save.
	applicationinstance.ApplicationIDValidator = applicationinstanceDescApplicationID.Validators[0].(func(string) error)
	// applicationinstanceDescEnvironmentID is the schema descriptor for environmentID field.
	applicationinstanceDescEnvironmentID := applicationinstanceFields[1].Descriptor()
	// applicationinstance.EnvironmentIDValidator is a validator for the "environmentID" field. It is called by the builders before save.
	applicationinstance.EnvironmentIDValidator = applicationinstanceDescEnvironmentID.Validators[0].(func(string) error)
	// applicationinstanceDescName is the schema descriptor for name field.
	applicationinstanceDescName := applicationinstanceFields[2].Descriptor()
	// applicationinstance.NameValidator is a validator for the "name" field. It is called by the builders before save.
	applicationinstance.NameValidator = applicationinstanceDescName.Validators[0].(func(string) error)
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
	// applicationmodulerelationshipDescVersion is the schema descriptor for version field.
	applicationmodulerelationshipDescVersion := applicationmodulerelationshipFields[2].Descriptor()
	// applicationmodulerelationship.VersionValidator is a validator for the "version" field. It is called by the builders before save.
	applicationmodulerelationship.VersionValidator = applicationmodulerelationshipDescVersion.Validators[0].(func(string) error)
	// applicationmodulerelationshipDescName is the schema descriptor for name field.
	applicationmodulerelationshipDescName := applicationmodulerelationshipFields[3].Descriptor()
	// applicationmodulerelationship.NameValidator is a validator for the "name" field. It is called by the builders before save.
	applicationmodulerelationship.NameValidator = applicationmodulerelationshipDescName.Validators[0].(func(string) error)
	applicationresourceMixin := schema.ApplicationResource{}.Mixin()
	applicationresourceMixinHooks0 := applicationresourceMixin[0].Hooks()
	applicationresource.Hooks[0] = applicationresourceMixinHooks0[0]
	applicationresourceInters := schema.ApplicationResource{}.Interceptors()
	applicationresource.Interceptors[0] = applicationresourceInters[0]
	applicationresourceMixinFields1 := applicationresourceMixin[1].Fields()
	_ = applicationresourceMixinFields1
	applicationresourceFields := schema.ApplicationResource{}.Fields()
	_ = applicationresourceFields
	// applicationresourceDescCreateTime is the schema descriptor for createTime field.
	applicationresourceDescCreateTime := applicationresourceMixinFields1[0].Descriptor()
	// applicationresource.DefaultCreateTime holds the default value on creation for the createTime field.
	applicationresource.DefaultCreateTime = applicationresourceDescCreateTime.Default.(func() time.Time)
	// applicationresourceDescUpdateTime is the schema descriptor for updateTime field.
	applicationresourceDescUpdateTime := applicationresourceMixinFields1[1].Descriptor()
	// applicationresource.DefaultUpdateTime holds the default value on creation for the updateTime field.
	applicationresource.DefaultUpdateTime = applicationresourceDescUpdateTime.Default.(func() time.Time)
	// applicationresource.UpdateDefaultUpdateTime holds the default value on update for the updateTime field.
	applicationresource.UpdateDefaultUpdateTime = applicationresourceDescUpdateTime.UpdateDefault.(func() time.Time)
	// applicationresourceDescInstanceID is the schema descriptor for instanceID field.
	applicationresourceDescInstanceID := applicationresourceFields[0].Descriptor()
	// applicationresource.InstanceIDValidator is a validator for the "instanceID" field. It is called by the builders before save.
	applicationresource.InstanceIDValidator = applicationresourceDescInstanceID.Validators[0].(func(string) error)
	// applicationresourceDescConnectorID is the schema descriptor for connectorID field.
	applicationresourceDescConnectorID := applicationresourceFields[1].Descriptor()
	// applicationresource.ConnectorIDValidator is a validator for the "connectorID" field. It is called by the builders before save.
	applicationresource.ConnectorIDValidator = applicationresourceDescConnectorID.Validators[0].(func(string) error)
	// applicationresourceDescModule is the schema descriptor for module field.
	applicationresourceDescModule := applicationresourceFields[3].Descriptor()
	// applicationresource.ModuleValidator is a validator for the "module" field. It is called by the builders before save.
	applicationresource.ModuleValidator = applicationresourceDescModule.Validators[0].(func(string) error)
	// applicationresourceDescMode is the schema descriptor for mode field.
	applicationresourceDescMode := applicationresourceFields[4].Descriptor()
	// applicationresource.ModeValidator is a validator for the "mode" field. It is called by the builders before save.
	applicationresource.ModeValidator = applicationresourceDescMode.Validators[0].(func(string) error)
	// applicationresourceDescType is the schema descriptor for type field.
	applicationresourceDescType := applicationresourceFields[5].Descriptor()
	// applicationresource.TypeValidator is a validator for the "type" field. It is called by the builders before save.
	applicationresource.TypeValidator = applicationresourceDescType.Validators[0].(func(string) error)
	// applicationresourceDescName is the schema descriptor for name field.
	applicationresourceDescName := applicationresourceFields[6].Descriptor()
	// applicationresource.NameValidator is a validator for the "name" field. It is called by the builders before save.
	applicationresource.NameValidator = applicationresourceDescName.Validators[0].(func(string) error)
	// applicationresourceDescDeployerType is the schema descriptor for deployerType field.
	applicationresourceDescDeployerType := applicationresourceFields[7].Descriptor()
	// applicationresource.DeployerTypeValidator is a validator for the "deployerType" field. It is called by the builders before save.
	applicationresource.DeployerTypeValidator = applicationresourceDescDeployerType.Validators[0].(func(string) error)
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
	// applicationrevisionDescInstanceID is the schema descriptor for instanceID field.
	applicationrevisionDescInstanceID := applicationrevisionFields[0].Descriptor()
	// applicationrevision.InstanceIDValidator is a validator for the "instanceID" field. It is called by the builders before save.
	applicationrevision.InstanceIDValidator = applicationrevisionDescInstanceID.Validators[0].(func(string) error)
	// applicationrevisionDescEnvironmentID is the schema descriptor for environmentID field.
	applicationrevisionDescEnvironmentID := applicationrevisionFields[1].Descriptor()
	// applicationrevision.EnvironmentIDValidator is a validator for the "environmentID" field. It is called by the builders before save.
	applicationrevision.EnvironmentIDValidator = applicationrevisionDescEnvironmentID.Validators[0].(func(string) error)
	// applicationrevisionDescModules is the schema descriptor for modules field.
	applicationrevisionDescModules := applicationrevisionFields[2].Descriptor()
	// applicationrevision.DefaultModules holds the default value on creation for the modules field.
	applicationrevision.DefaultModules = applicationrevisionDescModules.Default.([]types.ApplicationModule)
	// applicationrevisionDescSecrets is the schema descriptor for secrets field.
	applicationrevisionDescSecrets := applicationrevisionFields[3].Descriptor()
	// applicationrevision.DefaultSecrets holds the default value on creation for the secrets field.
	applicationrevision.DefaultSecrets = applicationrevisionDescSecrets.Default.(crypto.Map[string, string])
	// applicationrevisionDescInputVariables is the schema descriptor for inputVariables field.
	applicationrevisionDescInputVariables := applicationrevisionFields[5].Descriptor()
	// applicationrevision.DefaultInputVariables holds the default value on creation for the inputVariables field.
	applicationrevision.DefaultInputVariables = applicationrevisionDescInputVariables.Default.(property.Values)
	// applicationrevisionDescDeployerType is the schema descriptor for deployerType field.
	applicationrevisionDescDeployerType := applicationrevisionFields[8].Descriptor()
	// applicationrevision.DefaultDeployerType holds the default value on creation for the deployerType field.
	applicationrevision.DefaultDeployerType = applicationrevisionDescDeployerType.Default.(string)
	// applicationrevisionDescDuration is the schema descriptor for duration field.
	applicationrevisionDescDuration := applicationrevisionFields[9].Descriptor()
	// applicationrevision.DefaultDuration holds the default value on creation for the duration field.
	applicationrevision.DefaultDuration = applicationrevisionDescDuration.Default.(int)
	// applicationrevisionDescPreviousRequiredProviders is the schema descriptor for previousRequiredProviders field.
	applicationrevisionDescPreviousRequiredProviders := applicationrevisionFields[10].Descriptor()
	// applicationrevision.DefaultPreviousRequiredProviders holds the default value on creation for the previousRequiredProviders field.
	applicationrevision.DefaultPreviousRequiredProviders = applicationrevisionDescPreviousRequiredProviders.Default.([]types.ProviderRequirement)
	clustercostFields := schema.ClusterCost{}.Fields()
	_ = clustercostFields
	// clustercostDescConnectorID is the schema descriptor for connectorID field.
	clustercostDescConnectorID := clustercostFields[3].Descriptor()
	// clustercost.ConnectorIDValidator is a validator for the "connectorID" field. It is called by the builders before save.
	clustercost.ConnectorIDValidator = clustercostDescConnectorID.Validators[0].(func(string) error)
	// clustercostDescClusterName is the schema descriptor for clusterName field.
	clustercostDescClusterName := clustercostFields[4].Descriptor()
	// clustercost.ClusterNameValidator is a validator for the "clusterName" field. It is called by the builders before save.
	clustercost.ClusterNameValidator = clustercostDescClusterName.Validators[0].(func(string) error)
	// clustercostDescTotalCost is the schema descriptor for totalCost field.
	clustercostDescTotalCost := clustercostFields[5].Descriptor()
	// clustercost.DefaultTotalCost holds the default value on creation for the totalCost field.
	clustercost.DefaultTotalCost = clustercostDescTotalCost.Default.(float64)
	// clustercost.TotalCostValidator is a validator for the "totalCost" field. It is called by the builders before save.
	clustercost.TotalCostValidator = clustercostDescTotalCost.Validators[0].(func(float64) error)
	// clustercostDescAllocationCost is the schema descriptor for allocationCost field.
	clustercostDescAllocationCost := clustercostFields[7].Descriptor()
	// clustercost.DefaultAllocationCost holds the default value on creation for the allocationCost field.
	clustercost.DefaultAllocationCost = clustercostDescAllocationCost.Default.(float64)
	// clustercost.AllocationCostValidator is a validator for the "allocationCost" field. It is called by the builders before save.
	clustercost.AllocationCostValidator = clustercostDescAllocationCost.Validators[0].(func(float64) error)
	// clustercostDescIdleCost is the schema descriptor for idleCost field.
	clustercostDescIdleCost := clustercostFields[8].Descriptor()
	// clustercost.DefaultIdleCost holds the default value on creation for the idleCost field.
	clustercost.DefaultIdleCost = clustercostDescIdleCost.Default.(float64)
	// clustercost.IdleCostValidator is a validator for the "idleCost" field. It is called by the builders before save.
	clustercost.IdleCostValidator = clustercostDescIdleCost.Validators[0].(func(float64) error)
	// clustercostDescManagementCost is the schema descriptor for managementCost field.
	clustercostDescManagementCost := clustercostFields[9].Descriptor()
	// clustercost.DefaultManagementCost holds the default value on creation for the managementCost field.
	clustercost.DefaultManagementCost = clustercostDescManagementCost.Default.(float64)
	// clustercost.ManagementCostValidator is a validator for the "managementCost" field. It is called by the builders before save.
	clustercost.ManagementCostValidator = clustercostDescManagementCost.Validators[0].(func(float64) error)
	connectorMixin := schema.Connector{}.Mixin()
	connectorMixinHooks0 := connectorMixin[0].Hooks()
	connector.Hooks[0] = connectorMixinHooks0[0]
	connectorMixinFields1 := connectorMixin[1].Fields()
	_ = connectorMixinFields1
	connectorMixinFields2 := connectorMixin[2].Fields()
	_ = connectorMixinFields2
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
	connectorDescCreateTime := connectorMixinFields2[0].Descriptor()
	// connector.DefaultCreateTime holds the default value on creation for the createTime field.
	connector.DefaultCreateTime = connectorDescCreateTime.Default.(func() time.Time)
	// connectorDescUpdateTime is the schema descriptor for updateTime field.
	connectorDescUpdateTime := connectorMixinFields2[1].Descriptor()
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
	connector.DefaultConfigData = connectorDescConfigData.Default.(crypto.Properties)
	// connectorDescCategory is the schema descriptor for category field.
	connectorDescCategory := connectorFields[5].Descriptor()
	// connector.CategoryValidator is a validator for the "category" field. It is called by the builders before save.
	connector.CategoryValidator = connectorDescCategory.Validators[0].(func(string) error)
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
	moduleDescLabels := moduleFields[3].Descriptor()
	// module.DefaultLabels holds the default value on creation for the labels field.
	module.DefaultLabels = moduleDescLabels.Default.(map[string]string)
	// moduleDescSource is the schema descriptor for source field.
	moduleDescSource := moduleFields[4].Descriptor()
	// module.SourceValidator is a validator for the "source" field. It is called by the builders before save.
	module.SourceValidator = moduleDescSource.Validators[0].(func(string) error)
	// moduleDescID is the schema descriptor for id field.
	moduleDescID := moduleFields[0].Descriptor()
	// module.IDValidator is a validator for the "id" field. It is called by the builders before save.
	module.IDValidator = moduleDescID.Validators[0].(func(string) error)
	moduleversionMixin := schema.ModuleVersion{}.Mixin()
	moduleversionMixinHooks0 := moduleversionMixin[0].Hooks()
	moduleversion.Hooks[0] = moduleversionMixinHooks0[0]
	moduleversionMixinFields1 := moduleversionMixin[1].Fields()
	_ = moduleversionMixinFields1
	moduleversionFields := schema.ModuleVersion{}.Fields()
	_ = moduleversionFields
	// moduleversionDescCreateTime is the schema descriptor for createTime field.
	moduleversionDescCreateTime := moduleversionMixinFields1[0].Descriptor()
	// moduleversion.DefaultCreateTime holds the default value on creation for the createTime field.
	moduleversion.DefaultCreateTime = moduleversionDescCreateTime.Default.(func() time.Time)
	// moduleversionDescUpdateTime is the schema descriptor for updateTime field.
	moduleversionDescUpdateTime := moduleversionMixinFields1[1].Descriptor()
	// moduleversion.DefaultUpdateTime holds the default value on creation for the updateTime field.
	moduleversion.DefaultUpdateTime = moduleversionDescUpdateTime.Default.(func() time.Time)
	// moduleversion.UpdateDefaultUpdateTime holds the default value on update for the updateTime field.
	moduleversion.UpdateDefaultUpdateTime = moduleversionDescUpdateTime.UpdateDefault.(func() time.Time)
	// moduleversionDescModuleID is the schema descriptor for moduleID field.
	moduleversionDescModuleID := moduleversionFields[0].Descriptor()
	// moduleversion.ModuleIDValidator is a validator for the "moduleID" field. It is called by the builders before save.
	moduleversion.ModuleIDValidator = moduleversionDescModuleID.Validators[0].(func(string) error)
	// moduleversionDescVersion is the schema descriptor for version field.
	moduleversionDescVersion := moduleversionFields[1].Descriptor()
	// moduleversion.VersionValidator is a validator for the "version" field. It is called by the builders before save.
	moduleversion.VersionValidator = moduleversionDescVersion.Validators[0].(func(string) error)
	// moduleversionDescSource is the schema descriptor for source field.
	moduleversionDescSource := moduleversionFields[2].Descriptor()
	// moduleversion.SourceValidator is a validator for the "source" field. It is called by the builders before save.
	moduleversion.SourceValidator = moduleversionDescSource.Validators[0].(func(string) error)
	// moduleversionDescSchema is the schema descriptor for schema field.
	moduleversionDescSchema := moduleversionFields[3].Descriptor()
	// moduleversion.DefaultSchema holds the default value on creation for the schema field.
	moduleversion.DefaultSchema = moduleversionDescSchema.Default.(*types.ModuleSchema)
	perspectiveMixin := schema.Perspective{}.Mixin()
	perspectiveMixinHooks0 := perspectiveMixin[0].Hooks()
	perspective.Hooks[0] = perspectiveMixinHooks0[0]
	perspectiveMixinFields1 := perspectiveMixin[1].Fields()
	_ = perspectiveMixinFields1
	perspectiveFields := schema.Perspective{}.Fields()
	_ = perspectiveFields
	// perspectiveDescCreateTime is the schema descriptor for createTime field.
	perspectiveDescCreateTime := perspectiveMixinFields1[0].Descriptor()
	// perspective.DefaultCreateTime holds the default value on creation for the createTime field.
	perspective.DefaultCreateTime = perspectiveDescCreateTime.Default.(func() time.Time)
	// perspectiveDescUpdateTime is the schema descriptor for updateTime field.
	perspectiveDescUpdateTime := perspectiveMixinFields1[1].Descriptor()
	// perspective.DefaultUpdateTime holds the default value on creation for the updateTime field.
	perspective.DefaultUpdateTime = perspectiveDescUpdateTime.Default.(func() time.Time)
	// perspective.UpdateDefaultUpdateTime holds the default value on update for the updateTime field.
	perspective.UpdateDefaultUpdateTime = perspectiveDescUpdateTime.UpdateDefault.(func() time.Time)
	// perspectiveDescName is the schema descriptor for name field.
	perspectiveDescName := perspectiveFields[0].Descriptor()
	// perspective.NameValidator is a validator for the "name" field. It is called by the builders before save.
	perspective.NameValidator = perspectiveDescName.Validators[0].(func(string) error)
	// perspectiveDescStartTime is the schema descriptor for startTime field.
	perspectiveDescStartTime := perspectiveFields[1].Descriptor()
	// perspective.StartTimeValidator is a validator for the "startTime" field. It is called by the builders before save.
	perspective.StartTimeValidator = perspectiveDescStartTime.Validators[0].(func(string) error)
	// perspectiveDescEndTime is the schema descriptor for endTime field.
	perspectiveDescEndTime := perspectiveFields[2].Descriptor()
	// perspective.EndTimeValidator is a validator for the "endTime" field. It is called by the builders before save.
	perspective.EndTimeValidator = perspectiveDescEndTime.Validators[0].(func(string) error)
	// perspectiveDescBuiltin is the schema descriptor for builtin field.
	perspectiveDescBuiltin := perspectiveFields[3].Descriptor()
	// perspective.DefaultBuiltin holds the default value on creation for the builtin field.
	perspective.DefaultBuiltin = perspectiveDescBuiltin.Default.(bool)
	// perspectiveDescAllocationQueries is the schema descriptor for allocationQueries field.
	perspectiveDescAllocationQueries := perspectiveFields[4].Descriptor()
	// perspective.DefaultAllocationQueries holds the default value on creation for the allocationQueries field.
	perspective.DefaultAllocationQueries = perspectiveDescAllocationQueries.Default.([]types.QueryCondition)
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
	role.DefaultPolicies = roleDescPolicies.Default.(types.RolePolicies)
	// roleDescBuiltin is the schema descriptor for builtin field.
	roleDescBuiltin := roleFields[4].Descriptor()
	// role.DefaultBuiltin holds the default value on creation for the builtin field.
	role.DefaultBuiltin = roleDescBuiltin.Default.(bool)
	// roleDescSession is the schema descriptor for session field.
	roleDescSession := roleFields[5].Descriptor()
	// role.DefaultSession holds the default value on creation for the session field.
	role.DefaultSession = roleDescSession.Default.(bool)
	secretMixin := schema.Secret{}.Mixin()
	secretMixinHooks0 := secretMixin[0].Hooks()
	secret.Hooks[0] = secretMixinHooks0[0]
	secretMixinFields1 := secretMixin[1].Fields()
	_ = secretMixinFields1
	secretFields := schema.Secret{}.Fields()
	_ = secretFields
	// secretDescCreateTime is the schema descriptor for createTime field.
	secretDescCreateTime := secretMixinFields1[0].Descriptor()
	// secret.DefaultCreateTime holds the default value on creation for the createTime field.
	secret.DefaultCreateTime = secretDescCreateTime.Default.(func() time.Time)
	// secretDescUpdateTime is the schema descriptor for updateTime field.
	secretDescUpdateTime := secretMixinFields1[1].Descriptor()
	// secret.DefaultUpdateTime holds the default value on creation for the updateTime field.
	secret.DefaultUpdateTime = secretDescUpdateTime.Default.(func() time.Time)
	// secret.UpdateDefaultUpdateTime holds the default value on update for the updateTime field.
	secret.UpdateDefaultUpdateTime = secretDescUpdateTime.UpdateDefault.(func() time.Time)
	// secretDescName is the schema descriptor for name field.
	secretDescName := secretFields[1].Descriptor()
	// secret.NameValidator is a validator for the "name" field. It is called by the builders before save.
	secret.NameValidator = secretDescName.Validators[0].(func(string) error)
	// secretDescValue is the schema descriptor for value field.
	secretDescValue := secretFields[2].Descriptor()
	// secret.ValueValidator is a validator for the "value" field. It is called by the builders before save.
	secret.ValueValidator = secretDescValue.Validators[0].(func(string) error)
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
	subject.DefaultRoles = subjectDescRoles.Default.(types.SubjectRoles)
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
	Version = "v0.11.8" // Version of ent codegen.
)
