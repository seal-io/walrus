package casdoor

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
	"time"

	"github.com/seal-io/seal/pkg/consts"
	"github.com/seal-io/seal/pkg/rds"
	"github.com/seal-io/seal/utils/bytespool"
	"github.com/seal-io/seal/utils/files"
	"github.com/seal-io/seal/utils/log"
	"github.com/seal-io/seal/utils/strs"
)

const casdoorConfigPathEnvName = "BEEGO_CONFIG_PATH"

const embeddedCasdoorEndpointAddress = "http://127.0.0.1:8000"

//nolint:lll
const embeddedCasdoorConfigTmpl = `
appname = casdoor
httpport = 8000
runmode = dev
copyrequestbody = true
driverName = {{ .DataSourceDriver }}
dataSourceName = {{ .DataSourceName }}
dbName =
tableNamePrefix = casdoor_
showSql = false
redisEndpoint =
defaultStorageProvider =
isCloudIntranet = false
authState = "casdoor"
sock5Proxy = "127.0.0.1:10808"
verificationCodeTimeout = 10
initScore = 2000
logPostOnly = true
origin =
staticBaseUrl = "https://cdn.casbin.org"
isDemoMode = false
batchSize = 100
ldapServerPort = 389
languages = en,zh,es,fr,de,ja,ko,ru
quota = {"organization": -1, "user": -1, "application": -1, "provider": -1}
logConfig = {"filename": "logs/casdoor.log", "maxdays":99999, "perm":"0770"}
initDataFile = "./init_data.json"
sessionConfig = {"enableSetCookie":true,"cookieName":"casdoor_session_id","cookieLifeTime":3600,"providerConfig":"{{ .DataDir }}","gclifetime":3600,"domain":"","secure":false,"disableHTTPOnly":false}
`

type Embedded struct{}

func (Embedded) Run(ctx context.Context, dataSourceAddress string) error {
	runDataPath := filepath.Join(consts.DataDir, "casdoor")

	configPath, err := writeConfig(dataSourceAddress, runDataPath)
	if err != nil {
		return err
	}

	const cmdName = "casdoor"
	logger := log.WithName(cmdName)
	cmdArgs := []string{
		"-createDatabase=true",
	}
	logger.Infof("run: %s %s", cmdName, strs.Join(" ", cmdArgs...))
	cmd := exec.CommandContext(ctx, cmdName, cmdArgs...)

	cmd.Env = append(os.Environ(), casdoorConfigPathEnvName+"="+configPath)
	cmd.Stdout = logger.V(5)
	cmd.Stderr = logger.V(5)

	err = cmd.Run()
	if err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}

func (Embedded) GetAddress(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 90*time.Second)
	defer cancel()

	err := Wait(ctx, embeddedCasdoorEndpointAddress)
	if err != nil {
		return "", err
	}

	return embeddedCasdoorEndpointAddress, nil
}

func writeConfig(dataSourceAddress, dataDir string) (string, error) {
	dsd, dsn, err := rds.GetDriverAndName(dataSourceAddress)
	if err != nil {
		return "", err
	}

	tmpl, err := template.New("config").
		Parse(embeddedCasdoorConfigTmpl)
	if err != nil {
		return "", fmt.Errorf("error parsing casdoor config template: %w", err)
	}

	buf := bytespool.GetBuffer()
	defer func() { bytespool.Put(buf) }()

	err = tmpl.Execute(buf, map[string]string{
		"DataSourceDriver": dsd,
		"DataSourceName":   dsn,
		"DataDir":          dataDir,
	})
	if err != nil {
		return "", fmt.Errorf("error rendering casdoor config: %w", err)
	}

	configPath := os.Getenv(casdoorConfigPathEnvName)
	if configPath == "" {
		configPath = files.TempFile("")
	}

	return configPath, os.WriteFile(configPath, buf.Bytes(), 0o600)
}
