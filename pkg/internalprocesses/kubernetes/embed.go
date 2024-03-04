package kubernetes

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/seal-io/utils/netx"
	"github.com/seal-io/utils/osx"
	"github.com/seal-io/utils/stringx"
	"github.com/seal-io/utils/waitx"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"

	"github.com/seal-io/walrus/pkg/kubeconfig"
	"github.com/seal-io/walrus/pkg/system"
)

var embeddedKubeConfigPath = filepath.Join(clientcmd.RecommendedConfigDir, "k3s.yaml")

// Embedded represents the embedded Kubernetes cluster,
// which driven by k3s.
type Embedded struct{}

func (Embedded) Start(ctx context.Context) error {
	// Validate if run with privileged.
	if !osx.ExistsDevice("/dev/kmsg") {
		if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
			return errors.New(`require "securityContext.privileged" feature of walrus pod`)
		}

		return errors.New(`require "--privileged" flag to run walrus container`)
	}

	// NB(thxCode): With cgroup v2,
	// in order to nest walrus multiple times,
	// we need to create specific cgroup for k3s.
	cgroupRoot := "/"

	// Enable nested cgroup v2,
	// ref to https://github.com/moby/moby/issues/43093.
	if osx.ExistsFile("/sys/fs/cgroup/cgroup.controllers") {
		// Move the processes from the root group to the /init group,
		// otherwise writing subtree_control fails with EBUSY.
		if !osx.ExistsFile("/sys/fs/cgroup/init/cgroup.procs") {
			err := osx.Copy(
				"/sys/fs/cgroup/cgroup.procs",
				"/sys/fs/cgroup/init/cgroup.procs")
			if err != nil {
				return fmt.Errorf("error moving processes to root init group: %w", err)
			}
		}
		// Enable controllers to allow nesting.
		err := osx.Copy(
			"/sys/fs/cgroup/cgroup.controllers",
			"/sys/fs/cgroup/cgroup.subtree_control",
			osx.CopyWithModifier(func(data []byte) ([]byte, error) {
				if len(data) == 0 {
					return data, nil
				}
				cs := strings.Split(stringx.FromBytes(&data), " ")
				s := "+" + stringx.Join(" +", cs...)

				return stringx.ToBytes(&s), nil
			}))
		if err != nil {
			return fmt.Errorf("error enabling root group controllers: %w", err)
		}

		// Calculate specific cgroup root.
		cgroupRoot = "/k3s-" + stringx.SumByFNV64a(system.PrimarySubnet.Get())

		// Create specific cgroup root.
		err = os.Mkdir(path.Join("/sys/fs/cgroup", cgroupRoot), 0o755)
		if err != nil && !os.IsExist(err) {
			return fmt.Errorf("error creating cgroup: %w", err)
		}

		// Enable controllers to run k3s.
		err = osx.Copy(
			"/sys/fs/cgroup/cgroup.subtree_control",
			path.Join("/sys/fs/cgroup", cgroupRoot, "cgroup.subtree_control"),
			osx.CopyWithModifier(func(data []byte) ([]byte, error) {
				if len(data) == 0 {
					return data, nil
				}
				cs := strings.Split(stringx.FromBytes(&data), " ")
				s := "+" + stringx.Join(" +", cs...)

				return stringx.ToBytes(&s), nil
			}))
		if err != nil {
			return fmt.Errorf("error enabling %q group controllers: %w", cgroupRoot, err)
		}
	}

	var (
		k3sDataDir       = osx.Getenv("K3S_DATA_DIR", "/var/lib/k3s")
		k3sServerDataDir = filepath.Join(k3sDataDir, "server")
		runDataPath      = filepath.Join(system.DataDir, "k3s")
	)

	// Link run data directory.
	err := osx.Link(
		runDataPath,
		k3sServerDataDir,
		osx.LinkEvenIfNotFound(false, 0o766),
		osx.LinkInReplace())
	if err != nil {
		return fmt.Errorf("error link server data: %w", err)
	}

	// Reset server database.
	if osx.ExistsDir(filepath.Join(k3sServerDataDir, "db", "etcd")) {
		_ = os.Remove(filepath.Join(k3sServerDataDir, "db", "reset-flag")) // Clean reset flag.

		cmdArgs := []string{
			"server",
			"--cluster-reset",
			"--data-dir=" + k3sDataDir,
		}
		if err = runK3sWith(ctx, cmdArgs); err != nil {
			return err
		}
	}

	cmdArgs := []string{
		"server",
		"--cluster-init",
		"--etcd-disable-snapshots",
		"--disable=traefik,servicelb,metrics-server",
		"--disable-cloud-controller",
		"--disable-network-policy",
		"--disable-helm-controller",
		"--flannel-backend=host-gw",
		"--node-name=local",
		"--data-dir=" + k3sDataDir,
		"--write-kubeconfig=" + embeddedKubeConfigPath,
		"--kubelet-arg=cgroup-root=" + cgroupRoot,
		"--kubelet-arg=system-reserved=cpu=300m,memory=256Mi",
		"--kubelet-arg=kube-reserved=cpu=200m,memory=256Mi",
		"--kube-apiserver-arg=service-node-port-range=30000-30100",
	}

	// NB(thxCode): With embedded cluster,
	// in order to avoid the conflict with the default cluster CIDR when running in Kubernetes cluster,
	// we need to adjust the embedded cluster CIDR.
	if v := system.PrimarySubnet.Get(); v != "" {
		sn := netx.MustIPv4FromCIDR(v)
		cls := sn.Next().Next() // Offset 2 positions to get cluster CIDR.
		svc := cls.Next()       // Offset 1 position of cluster CIDR to get service CIDR.
		cmdArgs = append(cmdArgs,
			"--cluster-cidr="+cls.String(),
			"--service-cidr="+svc.String())

		// NB(thxCode): At present, it's impossible to reset the PodCIDR of a joined Kubernetes Node,
		// even the `--cluster-reset` introduced by k3s.
		// Embedded Kubernetes cluster is for demonstration purpose,
		// we never guarantee it for production using.
		// In order to restart from the mutable networking env,
		// we compare with the previous cluster CIDR here to determine erasing the data directory or not.
		clsFp := filepath.Join(k3sServerDataDir, "cluster.cidr")
		if osx.ExistsFile(clsFp) {
			bs, _ := os.ReadFile(clsFp)
			if string(bs) != cls.String() {
				_ = os.RemoveAll(filepath.Join(k3sServerDataDir, "db", "etcd"))
			}
		}
		_ = os.WriteFile(clsFp, []byte(cls.String()), 0o600)
	}

	return runK3sWith(ctx, cmdArgs)
}

func (Embedded) GetConfig(ctx context.Context) (string, *rest.Config, error) {
	err := waitx.PollUntilContextTimeout(ctx, time.Second, 300*time.Second, true,
		func(ctx context.Context) error {
			if osx.ExistsFile(embeddedKubeConfigPath) {
				return nil
			}
			return errors.New("wait for embedded kubeconfig")
		})
	if err != nil {
		return "", nil, err
	}

	cfg, err := kubeconfig.LoadRestConfig(embeddedKubeConfigPath)
	if err != nil {
		return "", nil, err
	}

	return embeddedKubeConfigPath, cfg, err
}

func runK3sWith(ctx context.Context, cmdArgs []string) error {
	const cmdName = "k3s"

	logger := klog.Background().WithName(cmdName)
	logger.Infof("run: %s %s", cmdName, stringx.Join(" ", cmdArgs...))
	cmd := exec.CommandContext(ctx, cmdName, cmdArgs...)
	cmd.SysProcAttr = getSysProcAttr()
	cmd.Stdout = logger.V(6)
	cmd.Stderr = logger.V(6)

	err := cmd.Run()
	if err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}
