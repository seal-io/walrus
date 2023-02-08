package casdoor

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"text/template"
	"time"

	"github.com/seal-io/seal/pkg/rds"
	"github.com/seal-io/seal/utils/bytespool"
	"github.com/seal-io/seal/utils/files"
	"github.com/seal-io/seal/utils/log"
)

const casdoorConfigPathEnvName = "BEEGO_CONFIG_PATH"

const embeddedCasdoorEndpointAddress = "http://127.0.0.1:8000"

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
sessionConfig = {"enableSetCookie":true,"cookieName":"casdoor_session_id","cookieLifeTime":3600,"providerConfig":"/var/lib/seal/casdoor","gclifetime":3600,"domain":"","secure":false,"disableHTTPOnly":false}
`

type Embedded struct{}

func (Embedded) Run(ctx context.Context, dataSourceAddress string) error {
	var configPath, err = writeConfig(dataSourceAddress)
	if err != nil {
		return err
	}

	const cmdName = "casdoor"
	var cmdArgs = []string{
		"-createDatabase=true",
	}
	var cmd = exec.CommandContext(ctx, cmdName, cmdArgs...)
	cmd.Env = append(os.Environ(), casdoorConfigPathEnvName+"="+configPath)
	var logger = log.WithName(cmdName).V(5)
	cmd.Stdout = logger
	cmd.Stderr = logger
	err = cmd.Run()
	if err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

func (Embedded) GetAddress(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 90*time.Second)
	defer cancel()

	var err = Wait(ctx, embeddedCasdoorEndpointAddress)
	if err != nil {
		return "", err
	}

	return embeddedCasdoorEndpointAddress, nil
}

func writeConfig(dataSourceAddress string) (string, error) {
	var dsd, dsn, err = rds.GetDriverAndName(dataSourceAddress)
	if err != nil {
		return "", err
	}

	tmpl, err := template.New("config").
		Parse(embeddedCasdoorConfigTmpl)
	if err != nil {
		return "", fmt.Errorf("error parsing casdoor config template: %w", err)
	}

	var buf = bytespool.GetBuffer()
	defer func() { bytespool.Put(buf) }()
	err = tmpl.Execute(buf, map[string]string{"DataSourceDriver": dsd, "DataSourceName": dsn})
	if err != nil {
		return "", fmt.Errorf("error rendering casdoor config: %w", err)
	}

	var configPath = os.Getenv(casdoorConfigPathEnvName)
	if configPath == "" {
		configPath = files.TempFile("")
	}
	return configPath, os.WriteFile(configPath, buf.Bytes(), 0600)
}
