package local

import (
	"errors"
	"fmt"
	"os/exec"
)

func UninstallLocalWalrus() error {
	if !IsDockerInstalled() {
		return errors.New("docker is not available")
	}

	if isWalrusDockerExtensionInstalled() {
		return uninstallWalrusDockerExtension()
	} else if isLocalWalrusContainerInstalled() {
		return uninstallLocalWalrusDockerContainer()
	}

	return errors.New("not installed")
}

func uninstallWalrusDockerExtension() error {
	cmd := exec.Command("docker", "extension", "rm", walrusDockerExtensionName)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%w: %s", err, string(output))
	}

	return nil
}

func uninstallLocalWalrusDockerContainer() error {
	cmd := exec.Command("docker", "rm", "--force", localWalrusContainerName)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%w: %s", err, string(output))
	}

	return nil
}
