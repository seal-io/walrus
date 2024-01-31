package local

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/AlecAivazis/survey/v2"

	"github.com/seal-io/walrus/utils/version"
)

const (
	localWalrusContainerName  = "local-walrus"
	walrusDockerExtensionName = "sealio/walrus-docker-extension"
)

func InstallLocalWalrus(opts InstallOptions) error {
	if !IsDockerInstalled() {
		return errors.New("docker is not available")
	}

	if IsLocalWalrusInstalled() {
		return nil
	}

	// Check if we can use docker extension to install.
	// Installation using docker extension does not support bootstrap configurations.
	// When env is set, we fall back to use docker engine.
	if IsDockerExtensionAvailable() && len(opts.Env) == 0 {
		confirm := ""
		prompt := &survey.Input{
			Message: "Install Walrus docker extension to proceed [y/N]",
		}

		if err := survey.AskOne(prompt, &confirm); err != nil {
			return err
		}

		if confirm != "y" {
			fmt.Println("Installation aborted.")
			os.Exit(1)
		}

		fmt.Println("Installing...")

		return InstallWalrusDockerExtension()
	}

	fmt.Println("Installing...")

	return InstallLocalWalrusDockerContainer(opts)
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

func InstallLocalWalrusDockerContainer(opts InstallOptions) error {
	runArgs := []string{
		"run",
		"-d",
		"--name",
		localWalrusContainerName,
		"-p",
		"7080:80",
		"-p",
		"7443:443",
		"--privileged",
		"-v",
		"/var/run/docker.sock:/var/run/docker.sock",
	}

	for _, env := range opts.Env {
		runArgs = append(runArgs, "-e", env)
	}

	runArgs = append(runArgs,
		"-e",
		"SERVER_SETTING_LOCAL_ENVIRONMENT_MODE=docker",
		"-e",
		"SERVER_ENABLE_AUTHN=false",
		fmt.Sprintf("sealio/walrus:%s", getLocalWalrusTag()),
	)

	// #nosec G204
	cmd := exec.Command("docker", runArgs...)
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
