package settings

import (
	"time"

	"github.com/seal-io/walrus/pkg/casdoor"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/types/crypto"
	"github.com/seal-io/walrus/utils/pointer"
	"github.com/seal-io/walrus/utils/strs"
)

// the built-in settings for deployer.
var (
	// DeployerHttpProxy indicates the address for proxying none SSL http outbound traffic used by deployer,
	// it's in form of http(s)://[user:password@]address[:port].
	DeployerHttpProxy = newValue(
		"DeployerHttpProxy",
		editable,
		initializeFromSpecifiedEnv("HTTP_PROXY", ""),
		modifyWith(httpUrl),
	)
	// DeployerHttpsProxy indicates the address for proxying SSL http outbound traffic used by deployer,
	// it's in form of http(s)://[user:password@]address[:port].
	DeployerHttpsProxy = newValue(
		"DeployerHttpsProxy",
		editable,
		initializeFromSpecifiedEnv("HTTPS_PROXY", ""),
		modifyWith(httpUrl),
	)
	// DeployerAllProxy indicates the address for proxying outbound traffic used by deployer,
	// it's in form of scheme://[user:password@]address[:port].
	DeployerAllProxy = newValue(
		"DeployerAllProxy",
		editable,
		initializeFromSpecifiedEnv("ALL_PROXY", ""),
		modifyWith(sockUrl),
	)
	// DeployerNoProxy indicates the host exclusion list when proxying outbound traffic used by deployer,
	// it's a comma-separated string.
	DeployerNoProxy = newValue(
		"DeployerNoProxy",
		editable,
		initializeFromSpecifiedEnv("NO_PROXY", ""),
		nil)
	// DeployerImage indicates the image used by deployer.
	DeployerImage = newValue(
		"DeployerImage",
		editable,
		// When the image is updated, sync the one in server Dockerfile.
		initializeFrom("sealio/terraform-deployer:v0.1.4"),
		modifyWith(notBlank, containerImageReference),
	)
)

// the built-in settings for server.
var (
	// BootPwdGainSource indicates the bootstrap password provision mode.
	BootPwdGainSource = newValue(
		"BootPwdGainSource",
		hidden,
		initializeFrom("Specified"),
		nil)
	// CasdoorCred keeps the AK/SK for accessing Casdoor server.
	CasdoorCred = newValue(
		"CasdoorCred",
		private,
		initializeFromJSON(casdoor.ApplicationCredential{}),
		modifyWith(once),
	)
	// CasdoorApiToken keeps the token for accessing Casdoor server.
	CasdoorApiToken = newValue(
		"CasdoorApiToken",
		private,
		nil,
		nil)
	// ServeUrl keeps the URL for accessing server.
	ServeUrl = newValue(
		"ServeUrl",
		editable,
		nil,
		modifyWith(notBlank, httpUrl))
	// ServeUiIndex keeps the address for serving UI.
	ServeUiIndex = newValue(
		"ServeUiIndex",
		editable|hidden,
		initializeFromEnv("https://walrus-ui-1303613262.cos.ap-guangzhou.myqcloud.com/latest/index.html"),
		modifyWith(notBlank, anyUrl),
	)
	// DataEncryptionSentry keeps the sentry for indicating whether enables data encryption.
	DataEncryptionSentry = newValue(
		"DataEncryptionSentry",
		private,
		nil,
		modifyWith(notBlank))
	// AuthsEncryptionAesGcmKey keeps the key for encrypting public token value with AES-GCM algorithm.
	AuthsEncryptionAesGcmKey = newValue(
		"AuthsEncryptionAesGcmKey",
		private,
		initializeFrom(strs.Hex(32)),
		modifyWith(never))
	// OpenAiApiToken keeps the openAI API token for generating module completions.
	// TODO protect the stored token.
	OpenAiApiToken = newValue(
		"OpenAiApiToken",
		editable|sensitive,
		nil,
		nil)
	// InstallationUUID keeps the uuid for installation.
	InstallationUUID = newValue(
		"InstallationUUID",
		private,
		initializeFrom(strs.Hex(16)),
		modifyWith(never))
	// EnableTelemetry keeps the user config for enable telemetry or not.
	EnableTelemetry = newValue(
		"EnableTelemetry",
		editable,
		initializeFrom("true"),
		modifyWith(notBlank))
	// EnableSyncCatalog keeps the user config for enable sync catalog or not.
	EnableSyncCatalog = newValue(
		"EnableSyncCatalog",
		editable,
		initializeFrom("true"),
		modifyWith(notBlank))
	// ImageRegistry config the image registry for seal tools, like finOps tools.
	ImageRegistry = newValue(
		"ImageRegistry",
		editable,
		initializeFrom("docker.io"),
		modifyWith(notBlank))
	EnableBuiltinCatalog = newValue(
		"EnableBuiltinCatalog",
		editable,
		initializeFrom("true"),
		modifyWith(notBlank))
	// SkipRemoteTLSVerify indicates whether skip SSL verification when accessing remote server.
	SkipRemoteTLSVerify = newValue(
		"SkipRemoteTLSVerify",
		editable,
		initializeFrom("false"),
		modifyWith(notBlank))
)

// the built-in settings for server cron jobs.
var (
	// ConnectorCostCollectCronExpr indicates the cron expression of collect cost data,
	// default cron expression means executing collection per hour,
	// the cron expression is in form of `Seconds Minutes Hours DayOfMonth Month DayOfWeek`.
	ConnectorCostCollectCronExpr = newValue(
		"ConnectorCostCollectCronExpr",
		editable,
		initializeFrom("0 0 * ? * *"),
		modifyWith(notBlank, cronExpression),
	)
	// ConnectorStatusSyncCronExpr indicates the cron expression of sync connector status,
	// default cron expression means executing check every 5 minutes.
	ConnectorStatusSyncCronExpr = newValue(
		"ConnectorStatusSyncCronExpr",
		editable,
		initializeFrom("0 */5 * ? * *"),
		modifyWith(notBlank, cronExpression),
	)
	// ResourceStatusSyncCronExpr indicates the cron expression of sync service resource status,
	// default cron expression means stating every 1 minute.
	ResourceStatusSyncCronExpr = newValue(
		"ResourceStatusSyncCronExpr",
		editable,
		initializeFrom("0 */1 * ? * *"),
		modifyWith(notBlank, cronExpression),
	)
	// ResourceLabelApplyCronExpr indicates the cron expression of set labels to service resource,
	// default cron expression means setting every 2 minutes.
	ResourceLabelApplyCronExpr = newValue(
		"ResourceLabelApplyCronExpr",
		editable,
		initializeFrom("0 */2 * ? * *"),
		modifyWith(notBlank, cronExpression),
	)
	// ResourceComponentsDiscoverCronExpr indicates the cron expression of discover service resource basics,
	// default cron expression means discovering every 1 minute.
	ResourceComponentsDiscoverCronExpr = newValue(
		"ResourceComponentsDiscoverCronExpr",
		editable,
		initializeFrom("0 */1 * ? * *"),
		modifyWith(notBlank, cronExpression),
	)
	// TokenDeploymentExpiredCleanCronExpr indicates the cron expression of clean expired deployment token,
	// default cron expression means cleaning up every 30 minutes.
	TokenDeploymentExpiredCleanCronExpr = newValue(
		"TokenDeploymentExpiredCleanCronExpr",
		hidden,
		initializeFrom("0 */30 * ? * *"),
		modifyWith(notBlank, cronExpression),
	)
	// ServiceRelationshipCheckCronExpr indicates the cron expression of deploy scheduled service,
	// default cron expression means deploying every 30 seconds.
	ServiceRelationshipCheckCronExpr = newValue(
		"ServiceRelationshipCheckCronExpr",
		editable,
		initializeFrom("*/30 * * ? * *"),
		modifyWith(notBlank, cronExpression),
	)
	// TelemetryPeriodicReportCronExpr indicates the cron expression of telemetry synchronization event,
	// default cron expression means sync at 2 o'clock evey day.
	TelemetryPeriodicReportCronExpr = newValue(
		"TelemetryPeriodicReportCronExpr",
		private,
		initializeFrom("0 0 2 * * *"),
		modifyWith(notBlank, cronExpression),
	)
	// CatalogTemplateSyncCronExpr indicates the cron expression of catalog template synchronization event,
	// default cron expression means sync at 1 o'clock evey day, and new cron must be at least 30 minutes.
	CatalogTemplateSyncCronExpr = newValue(
		"CatalogTemplateSyncCronExpr",
		private,
		initializeFrom("0 0 1 * * *"),
		modifyWith(notBlank, cronExpression, cronAtLeast(30*time.Minute)),
	)
)

// setting property list.
const (
	hidden uint8 = 1 << (iota)
	editable
	sensitive
	private
)

var (
	valuesOrder []string
	valuesIndex = map[string]value{}
)

// newValue creates a value with the given name and modifier,
// then indexes the new value by its name.
func newValue(name string, property uint8, initialize initializer, modify modifier) value {
	if modify == nil {
		modify = modifyWith(many)
	}

	if initialize == nil {
		initialize = initializeFromEnv("")
		if property&private == private {
			initialize = initializeFrom("")
		}
	}
	v := value{
		refer: model.Setting{
			Name:      name,
			Value:     crypto.String(initialize(name)),
			Hidden:    pointer.Bool(property&hidden == hidden),
			Editable:  pointer.Bool(property&editable == editable),
			Sensitive: pointer.Bool(property&sensitive == sensitive),
			Private:   pointer.Bool(property&private == private),
		},
		modify: modify,
	}

	valuesOrder = append(valuesOrder, name)
	valuesIndex[name] = v

	return v
}

// ForEach iterates each setting in input.
func ForEach(input func(setting model.Setting) error) error {
	if input == nil {
		return nil
	}

	for _, n := range valuesOrder {
		err := input(valuesIndex[n].refer)
		if err != nil {
			return err
		}
	}

	return nil
}

// All returns all settings.
func All() (r model.Settings) {
	_ = ForEach(func(s model.Setting) error {
		r = append(r, &s)
		return nil
	})

	return
}

// Index returns the setting with the given name.
func Index(name string) Value {
	return valuesIndex[name]
}
