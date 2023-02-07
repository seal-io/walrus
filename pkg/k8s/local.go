package k8s

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/seal-io/seal/utils/files"
	"github.com/seal-io/seal/utils/log"
)

var (
	home                   string
	embeddedKubeConfigPath string
)

func init() {
	var err error
	home, err = os.UserHomeDir()
	if err != nil {
		panic(fmt.Errorf("failed to get home dir: %w", err))
	}
	embeddedKubeConfigPath = filepath.Join(home, clientcmd.RecommendedHomeDir, "k3s.yaml")
}

type Embedded struct{}

func (Embedded) Run(ctx context.Context) error {
	const cmdName = "k3s"
	var cmdArgs = []string{
		"server",
		"--cluster-init",
		"--disable=traefik,servicelb,metrics-server",
		"--node-name=local-node",
		"--data-dir=/var/lib/seal/k3s",
		fmt.Sprintf("--write-kubeconfig=%s", embeddedKubeConfigPath),
		"--kubelet-arg=system-reserved=cpu=300m,memory=256Mi",
		"--kubelet-arg=kube-reserved=cpu=200m,memory=256Mi",
	}
	var cmd = exec.CommandContext(ctx, cmdName, cmdArgs...)
	cmd.SysProcAttr = getSysProcAttr()
	var logger = log.WithName(cmdName).V(5)
	cmd.Stdout = logger
	cmd.Stderr = logger
	var err = cmd.Run()
	if err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

func (Embedded) GetConfig(ctx context.Context) (string, *rest.Config, error) {
	ctx, cancel := context.WithTimeout(ctx, 90*time.Second)
	defer cancel()

	var err = wait.PollUntilWithContext(ctx, time.Second, func(ctx context.Context) (bool, error) {
		return files.Exists(embeddedKubeConfigPath), nil
	})
	if err != nil {
		return "", nil, err
	}

	cfg, err := LoadConfig(embeddedKubeConfigPath)
	if err != nil {
		return "", nil, err
	}

	err = Wait(ctx, cfg)
	if err != nil {
		return "", nil, err
	}
	return embeddedKubeConfigPath, cfg, err
}
