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
	admit       Admission
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
		kubeclientset.UpdateAfterAlign(alignFn))
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
	init Initializer, admit Admission,
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
			"it's in form of http(s)://[user:password@]address[:port].",
		_SettingPropEditable,
		InitializeFromSpecifiedEnv("HTTP_PROXY"),
		AdmitWith(AllowBlank, AllowHttpUrl),
	)
	DeployerHttpsProxy = newSetting(
		"deployer-https-proxy",
		"Indicates an address for proxying SSL http outbound traffic used by deployer, "+
			"it's in form of http(s)://[user:password@]address[:port].",
		_SettingPropEditable,
		InitializeFromSpecifiedEnv("HTTPS_PROXY"),
		AdmitWith(AllowBlank, AllowHttpUrl),
	)
	DeployerAllProxy = newSetting(
		"deployer-all-proxy",
		"Indicates an address for proxying all outbound traffic used by deployer, "+
			"it's in form of scheme://[user:password@]address[:port].",
		_SettingPropEditable,
		InitializeFromSpecifiedEnv("ALL_PROXY"),
		AdmitWith(AllowBlank, AllowSockUrl),
	)
	DeployerNoProxy = newSetting(
		"deployer-no-proxy",
		"Indicates addresses that should not to proxy, "+
			"it's in form of comma separated list of IPs or DNS names.",
		_SettingPropEditable,
		InitializeFromSpecifiedEnv("NO_PROXY"),
		Allow,
	)
	DeployerImage = newSetting(
		"deployer-image",
		"Indicates the image used by deployer.",
		_SettingPropEditable,
		InitializeFromEnv("sealio/terraform-deployer:v1.5.7-seal.1"),
		AdmitWith(DisallowBlank, AllowContainerImageReference),
	)
	DeployerNetworkMirrorUrl = newSetting(
		"deployer-network-mirror-url",
		"Indicates the URL to configure the network mirror for deployer.",
		_SettingPropEditable,
		InitializeFromEnv(),
		AdmitWith(AllowBlank, AllowHttpUrl),
	)
)

// the built-in settings for server.
var (
	ServeIdentify = newSetting(
		"serve-identify",
		"Indicates the UUID after server installation, "+
			"it's used for telemetry.",
		_SettingPropPrivate,
		InitializeFrom(stringx.Hex(16)),
		Disallow,
	)
	ServeUiUrl = newSetting(
		"serve-ui-url",
		"Indicates a URL to provide the server UI, "+
			"it's in form of scheme://address[:port]/path.",
		_SettingPropPrivate,
		InitializeFromEnv("https://walrus-ui-1303613262.cos.ap-guangzhou.myqcloud.com/latest/index.html"), // nolint: lll
		AdmitWith(DisallowBlank, AllowUrl),
	)
	ServeWalrusFilesUrl = newSetting(
		"serve-walrus-files-url",
		"Indicates a URL to provide the WalrusFiles.",
		_SettingPropEditable,
		InitializeFromEnv("https://github.com/seal-io/walrus-file-hub"),
		AdmitWith(DisallowBlank, AllowUrl),
	)
	ServeUrl = newSetting(
		"serve-url",
		"Indicates the URL to access server, "+
			"it's in form of https://address[:port].",
		_SettingPropEditable,
		InitializeFromEnv(),
		AdmitWith(DisallowBlank, AllowHttpUrl),
	)
	EnableTelemetry = newSetting(
		"enable-telemetry",
		"Indicates whether to enable telemetry.",
		_SettingPropEditable,
		InitializeFromEnv("true"),
		AllowBoolean,
	)
	EnableSyncCatalog = newSetting(
		"enable-sync-catalog",
		"Indicates whether to enable the catalog synchronization, "+
			"If enabled, the server will synchronize all versioned templates under the catalog.",
		_SettingPropEditable,
		InitializeFromEnv("true"),
		AllowBoolean,
	)
	EnableBuiltInCatalog = newSetting(
		"enable-builtin-catalog",
		"Indicates whether to enable the builtin catalog.",
		_SettingPropEditable,
		InitializeFromEnv("true"),
		AllowBoolean,
	)
	EnableRemoteTlsVerify = newSetting(
		"enable-remote-tls-verify",
		"Indicates whether to enable the remote TLS verification.",
		_SettingPropEditable,
		InitializeFromEnv("true"),
		AllowBoolean,
	)
	OpenAiApiToken = newSetting(
		"openai-api-token",
		"Indicates the OpenAI API token for completing template generation.",
		_SettingPropEditable|_SettingPropSensitive,
		InitializeFromEnv(),
		Allow,
	)
	ImageRegistry = newSetting(
		"image-registry",
		"Indicates the registry used by the server to pull images.",
		_SettingPropEditable,
		InitializeFromEnv("docker.io"),
		AdmitWith(AllowBlank, AllowContainerRegistry),
	)
	DefaultEnvironmentMode = newSetting(
		"default-environment-mode",
		"Indicates the default environment mode.",
		_SettingPropPrivate,
		InitializeFromEnv("kubernetes"),
		Disallow,
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
		InitializeFromEnv("0 */5 * ? * *" /* every 5 minutes */),
		AllowCronExpressionAtLeast(1*time.Minute),
	)
	CatalogSyncCron = newSetting(
		"catalog-sync-cron",
		"Indicates the Cron Expression of sync catalog, "+
			"default Cron Expression means executing check every day at 1 o'clock. "+
			"The Cron Expression to sync is at least 30 minutes.",
		_SettingPropEditable,
		InitializeFromEnv("0 0 1 * * *" /* every day at 1 o'clock */),
		AllowCronExpressionAtLeast(30*time.Minute),
	)
	ResourceRelationshipCheckCron = newSetting(
		"resource-relationship-check-cron",
		"Indicates the Cron Expression of check resource relationship, "+
			"default Cron Expression means executing check every 30 minutes. "+
			"The Cron Expression to sync is at least 1 minute.",
		_SettingPropEditable,
		InitializeFromEnv("*/30 * * ? * *" /* every 30 minutes */),
		AllowCronExpressionAtLeast(1*time.Minute),
	)
	ResourceComponentStatusSyncCron = newSetting(
		"resource-component-status-sync-cron",
		"Indicates the Cron Expression of sync resource component status, "+
			"default Cron Expression means executing check every 1 minutes. "+
			"The Cron Expression to sync is at least 30 seconds",
		_SettingPropEditable,
		InitializeFromEnv("0 */1 * ? * *" /* every 1 minutes */),
		AllowCronExpressionAtLeast(30*time.Second),
	)
	TelemetryReportCron = newSetting(
		"telemetry-report-cron",
		"Indicates the Cron Expression of report telemetry, "+
			"default Cron Expression means executing check every day at 2 o'clock. "+
			"The Cron Expression to sync is at least 1 hour.",
		_SettingPropEditable,
		InitializeFromEnv("0 0 2 * * *" /* every day at 2 o'clock */),
		AllowCronExpressionAtLeast(1*time.Hour),
	)
)
