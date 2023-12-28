package local

import (
	"os/exec"
	"strings"
)

func IsDockerInstalled() bool {
	cmd := exec.Command("docker", "--version")

	output, err := cmd.Output()
	if err != nil {
		return false
	}

	return strings.Contains(string(output), "Docker version")
}

func IsDockerExtensionAvailable() bool {
	cmd := exec.Command("docker", "extension")

	if err := cmd.Run(); err != nil {
		return false
	}

	return true
}

func IsLocalWalrusInstalled() bool {
	return isWalrusDockerExtensionInstalled() || isLocalWalrusContainerInstalled()
}

func isWalrusDockerExtensionInstalled() bool {
	cmd := exec.Command("docker", "extension", "ls")

	out, err := cmd.Output()
	if err != nil {
		return false
	}

	return strings.Contains(string(out), walrusDockerExtensionName)
}

func isLocalWalrusContainerInstalled() bool {
	cmd := exec.Command("docker", "ps", "-aqf", "name=local-walrus")

	out, err := cmd.Output()
	if err != nil {
		return false
	}

	return len(out) > 0
}
