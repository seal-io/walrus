package local

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/AlecAivazis/survey/v2"

	"github.com/seal-io/walrus/utils/version"
)

const (
	localWalrusContainerName  = "local-walrus"
	walrusDockerExtensionName = "sealio/walrus-docker-extension"
)

func InstallLocalWalrus() error {
	if !IsDockerInstalled() {
		return errors.New("docker is not available")
	}

	if IsLocalWalrusInstalled() {
		return nil
	}

	if IsDockerExtensionAvailable() {
		confirm := ""
		prompt := &survey.Input{
			Message: "Install Walrus docker extension to proceed [y/N]",
		}

		if err := survey.AskOne(prompt, &confirm); err != nil {
			return err
		}

		if confirm != "y" {
			return nil
		}

		fmt.Println("Installing...")

		return InstallWalrusDockerExtension()
	}

	fmt.Println("Docker extension is not available, fall back to use docker engine...")

	return InstallLocalWalrusDockerContainer()
}

func InstallWalrusDockerExtension() error {
	// #nosec G204
	cmd := exec.Command(
		"docker",
		"extension",
		"install",
		fmt.Sprintf("%s:%s", walrusDockerExtensionName, getLocalWalrusTag()),
		"--force",
	)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%w: %s", err, string(output))
	}

	return nil
}

func InstallLocalWalrusDockerContainer() error {
	// #nosec G204
	cmd := exec.Command("docker", "run",
		"-d",
		"--name",
		localWalrusContainerName,
		"-p",
		"7080:80",
		"-p",
		"7443:443",
		"--privileged",
		"-e",
		"SERVER_SETTING_LOCAL_ENVIRONMENT_MODE=docker",
		"-e",
		"SERVER_BUILTIN_CATALOG_PROVIDER",
		"-v",
		"/var/run/docker.sock:/var/run/docker.sock",
		fmt.Sprintf("sealio/walrus:%s", getLocalWalrusTag()),
		"walrus",
		"--enable-authn=false",
	)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%w: %s", err, string(output))
	}

	return nil
}

func getLocalWalrusTag() string {
	if version.IsValid() {
		return version.Version
	}

	return "main"
}
