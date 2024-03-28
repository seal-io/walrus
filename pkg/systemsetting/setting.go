package systemsetting

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/seal-io/utils/stringx"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/seal-io/walrus/pkg/kubeclientset"
	"github.com/seal-io/walrus/pkg/system"
	"github.com/seal-io/walrus/pkg/systemkuberes"
	"github.com/seal-io/walrus/pkg/systemsetting/setting"
)

type _SettingProp uint8

const (
	_SettingPropPrivate _SettingProp = 1 << (iota)
	_SettingPropHidden
	_SettingPropEditable
	_SettingPropSensitive
)

type Setting struct {
	name        string
	description string
	defVal      string
	admit       setting.Admission
	props       _SettingProp
}

// Name returns the name of the setting.
func (s Setting) Name() string {
	return s.name
}

// Description returns the description of the setting.
func (s Setting) Description() string {
	return s.description
}

// Private returns true if the setting is private.
func (s Setting) Private() bool {
	return s.props == _SettingPropPrivate
}

// Hidden returns true if the setting is hidden.
func (s Setting) Hidden() bool {
	return s.props&_SettingPropHidden == _SettingPropHidden
}

// Editable returns true if the setting is editable.
func (s Setting) Editable() bool {
	return s.props&_SettingPropEditable == _SettingPropEditable
}

// Sensitive returns true if the setting is sensitive.
func (s Setting) Sensitive() bool {
	return s.props&_SettingPropSensitive == _SettingPropSensitive
}

// Configure configures the value of the setting.
func (s Setting) Configure(ctx context.Context, newVal string) error {
	loopbackKubeCli := system.LoopbackKubeClient.Get()

	// Update.
	secCli := loopbackKubeCli.CoreV1().Secrets(systemkuberes.SystemNamespaceName)
	eSec := &core.Secret{
		ObjectMeta: meta.ObjectMeta{
			Name: DelegatedSecretName,
		},
	}
	alignFn := func(aSec *core.Secret) (*core.Secret, bool, error) {
		var oldVal string
		if aSec.Data != nil {
			oldVal = string(aSec.Data[s.name])
		}
		if admitErr := s.admit(ctx, oldVal, newVal); admitErr != nil {
			// NB(thxCode): Skip update if the new value is invalid.
			return nil, false, admitErr
		}
		if oldVal == newVal {
			// Skip update if the new value is the same as the old value.
			return nil, true, nil
		}
		// Update the value of the setting.
		aSec.Data[s.name] = []byte(newVal)
		return aSec, false, nil
	}

	_, err := kubeclientset.Update(ctx, secCli, eSec,
		kubeclientset.WithUpdateAlign(alignFn))
	if err != nil {
		return fmt.Errorf("configure setting %s: %w", s.name, err)
	}

	return nil
}

// Value returns the value of the setting.
func (s Setting) Value(ctx context.Context) (string, error) {
	loopbackKubeCli := system.LoopbackKubeClient.Get()

	sec, err := loopbackKubeCli.CoreV1().
		Secrets(systemkuberes.SystemNamespaceName).
		Get(ctx, DelegatedSecretName, meta.GetOptions{ResourceVersion: "0"})
	if err != nil {
		return "", fmt.Errorf("get value of setting %s: %w", s.name, err)
	}

	if sec.Data == nil || sec.Data[s.name] == nil {
		return "", fmt.Errorf("get value of setting %s: not found", s.name)
	}

	return string(sec.Data[s.name]), nil
}

// ValueBool returns the bool value of the setting.
func (s Setting) ValueBool(ctx context.Context) (bool, error) {
	v, err := s.Value(ctx)
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(v)
}

// ValueInt64 returns the int64 value of the setting.
func (s Setting) ValueInt64(ctx context.Context) (int64, error) {
	v, err := s.Value(ctx)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(v, 10, 64)
}

// ValueUint64 returns the uint64 value of the setting.
func (s Setting) ValueUint64(ctx context.Context) (uint64, error) {
	v, err := s.Value(ctx)
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(v, 10, 64)
}

// ValueFloat64 returns the float64 value of the setting.
func (s Setting) ValueFloat64(ctx context.Context) (float64, error) {
	v, err := s.Value(ctx)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(v, 64)
}

// ValueURL returns the *url.URL value of the setting.
func (s Setting) ValueURL(ctx context.Context) (*url.URL, error) {
	v, err := s.Value(ctx)
	if err != nil {
		return nil, err
	}
	return url.Parse(v)
}

var settings = map[string]Setting{}

func newSetting(
	name, description string, props _SettingProp,
	init setting.Initializer, admit setting.Admission,
) Setting {
	settings[name] = Setting{
		name:        name,
		defVal:      init(name),
		admit:       admit,
		description: description,
		props:       props,
	}
	return settings[name]
}

// Index returns the setting with the given name.
func Index(name string) (Setting, bool) {
	s, ok := settings[name]
	return s, ok
}

// Initials returns a map with all settings and their default values.
func Initials() map[string]string {
	m := map[string]string{}
	for k, v := range settings {
		m[k] = v.defVal
	}
	return m
}

// the built-in settings for deployer.
var (
	DeployerHttpProxy = newSetting(
		"deployer-http-proxy",
		"Indicates an address for proxying none SSL http outbound traffic used by deployer, "+
			"it's in form of [http|https]://[user:password@]address[:port].",
		_SettingPropEditable,
		setting.InitializeFromSpecifiedEnv("HTTP_PROXY"),
		setting.AdmitWith(setting.AllowBlank, setting.AllowUrlWithSchema("http", "https")),
	)
	DeployerHttpsProxy = newSetting(
		"deployer-https-proxy",
		"Indicates an address for proxying SSL http outbound traffic used by deployer, "+
			"it's in form of [http|https]://[user:password@]address[:port].",
		_SettingPropEditable,
		setting.InitializeFromSpecifiedEnv("HTTPS_PROXY"),
		setting.AdmitWith(setting.AllowBlank, setting.AllowUrlWithSchema("http", "https")),
	)
	DeployerAllProxy = newSetting(
		"deployer-all-proxy",
		"Indicates an address for proxying all outbound traffic used by deployer, "+
			"it's in form of [sock4|sock5]://[user:password@]address[:port].",
		_SettingPropEditable,
		setting.InitializeFromSpecifiedEnv("ALL_PROXY"),
		setting.AdmitWith(setting.AllowBlank, setting.AllowUrlWithSchema("sock4", "sock5")),
	)
	DeployerNoProxy = newSetting(
		"deployer-no-proxy",
		"Indicates addresses that should not to proxy, "+
			"it's in form of comma separated list of IPs or DNS names.",
		_SettingPropEditable,
		setting.InitializeFromSpecifiedEnv("NO_PROXY"),
		setting.Allow,
	)
	TerraformDeployerImage = newSetting(
		"terraform-deployer-image",
		"Indicates the image used by Terraform deployer.",
		_SettingPropEditable,
		setting.InitializeFromEnv("sealio/terraform-deployer:v1.5.7-seal.1"),
		setting.AdmitWith(setting.DisallowBlank, setting.AllowContainerImageReference),
	)
	TerraformDeployerNetworkMirrorUrl = newSetting(
		"terraform-deployer-network-mirror-url",
		"Indicates the URL to configure the network mirror for Terraform deployer, "+
			"it's in form of https://address[:port]/path/, must with a tail slash. ",
		_SettingPropEditable,
		setting.InitializeFromEnv(),
		setting.AdmitWith(setting.AllowBlank, setting.AllowUrlWithSchema("https")),
	)
)

// the built-in settings for server.
var (
	ServeIdentify = newSetting(
		"serve-identify",
		"Indicates the UUID after server installation, "+
			"it's used for telemetry.",
		_SettingPropPrivate,
		setting.InitializeFrom(stringx.Hex(16)),
		setting.Disallow,
	)
	ServeUiUrl = newSetting(
		"serve-ui-url",
		"Indicates a URL to provide the server UI, "+
			"it's in form of [https|file]://address[:port]/path.",
		_SettingPropPrivate,
		setting.InitializeFromEnv("https://walrus-ui-1303613262.cos.ap-guangzhou.myqcloud.com/latest/index.html"), // nolint: lll
		setting.AdmitWith(setting.DisallowBlank, setting.AllowUrlWithSchema("https", "file")),
	)
	ServeWalrusFilesUrl = newSetting(
		"serve-walrus-files-url",
		"Indicates a URL to provide the WalrusFiles, "+
			"it's in form of [http|https|file]://address[:port]/path.",
		_SettingPropEditable,
		setting.InitializeFromEnv("https://github.com/seal-io/walrus-file-hub"),
		setting.AdmitWith(setting.DisallowBlank, setting.AllowUrlWithSchema("http", "https", "file")),
	)
	ServeUrl = newSetting(
		"serve-url",
		"Indicates the URL to access server, "+
			"it's in form of https://address[:port].",
		_SettingPropEditable,
		setting.InitializeFromEnv(),
		setting.AdmitWith(setting.DisallowBlank, setting.AllowUrlWithSchema("https")),
	)
	ServeObjectStorageUrl = newSetting(
		"serve-object-storage-url",
		"Indicates the URL to provide the Object Storage, "+
			"it's in form of s3://[accessKey[:secretKey]@]endpoint[:port]/bucketName[?param1=value1&...&paramN=valueN].",
		_SettingPropEditable|_SettingPropSensitive,
		setting.InitializeFromEnv(),
		setting.AdmitWith(setting.DisallowBlank, setting.AllowUrlWithSchema("s3")),
	)
	EnableTelemetry = newSetting(
		"enable-telemetry",
		"Indicates whether to enable telemetry.",
		_SettingPropEditable,
		setting.InitializeFromEnv("true"),
		setting.AllowBoolean,
	)
	EnableSyncCatalog = newSetting(
		"enable-sync-catalog",
		"Indicates whether to enable the catalog synchronization, "+
			"If enabled, the server will synchronize all versioned templates under the catalog.",
		_SettingPropEditable,
		setting.InitializeFromEnv("true"),
		setting.AllowBoolean,
	)
	EnableBuiltInCatalog = newSetting(
		"enable-builtin-catalog",
		"Indicates whether to enable the builtin catalog.",
		_SettingPropEditable,
		setting.InitializeFromEnv("true"),
		setting.AllowBoolean,
	)
	EnableRemoteTlsVerify = newSetting(
		"enable-remote-tls-verify",
		"Indicates whether to enable the remote TLS verification.",
		_SettingPropEditable,
		setting.InitializeFromEnv("true"),
		setting.AllowBoolean,
	)
	ImageRegistry = newSetting(
		"image-registry",
		"Indicates the registry used by the server to pull images.",
		_SettingPropEditable,
		setting.InitializeFromEnv("docker.io"),
		setting.AdmitWith(setting.AllowBlank, setting.AllowContainerRegistry),
	)
	DefaultEnvironmentMode = newSetting(
		"default-environment-mode",
		"Indicates the default environment mode.",
		_SettingPropPrivate,
		setting.InitializeFromEnv("kubernetes"),
		setting.Disallow,
	)
)

// the built-in settings for server cron jobs.
var (
	ConnectorStatusSyncCron = newSetting(
		"connector-status-sync-cron",
		"Indicates the Cron Expression of sync connector status, "+
			"default Cron Expression means executing check every 5 minutes. "+
			"The Cron Expression to sync is at least 1 minute.",
		_SettingPropEditable,
		setting.InitializeFromEnv("0 */5 * ? * *" /* every 5 minutes */),
		setting.AllowCronExpressionAtLeast(1*time.Minute),
	)
	CatalogSyncCron = newSetting(
		"catalog-sync-cron",
		"Indicates the Cron Expression of sync catalog, "+
			"default Cron Expression means executing check every day at 1 o'clock. "+
			"The Cron Expression to sync is at least 30 minutes.",
		_SettingPropEditable,
		setting.InitializeFromEnv("0 0 1 * * *" /* every day at 1 o'clock */),
		setting.AllowCronExpressionAtLeast(30*time.Minute),
	)
	ResourceRelationshipCheckCron = newSetting(
		"resource-relationship-check-cron",
		"Indicates the Cron Expression of check resource relationship, "+
			"default Cron Expression means executing check every 30 minutes. "+
			"The Cron Expression to sync is at least 1 minute.",
		_SettingPropEditable,
		setting.InitializeFromEnv("*/30 * * ? * *" /* every 30 minutes */),
		setting.AllowCronExpressionAtLeast(1*time.Minute),
	)
	ResourceComponentStatusSyncCron = newSetting(
		"resource-component-status-sync-cron",
		"Indicates the Cron Expression of sync resource component status, "+
			"default Cron Expression means executing check every 1 minutes. "+
			"The Cron Expression to sync is at least 30 seconds",
		_SettingPropEditable,
		setting.InitializeFromEnv("0 */1 * ? * *" /* every 1 minutes */),
		setting.AllowCronExpressionAtLeast(30*time.Second),
	)
	TelemetryReportCron = newSetting(
		"telemetry-report-cron",
		"Indicates the Cron Expression of report telemetry, "+
			"default Cron Expression means executing check every day at 2 o'clock. "+
			"The Cron Expression to sync is at least 1 hour.",
		_SettingPropEditable,
		setting.InitializeFromEnv("0 0 2 * * *" /* every day at 2 o'clock */),
		setting.AllowCronExpressionAtLeast(1*time.Hour),
	)
)
