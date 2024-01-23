package k8s

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

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/seal-io/walrus/pkg/consts"
	"github.com/seal-io/walrus/pkg/servervars"
	"github.com/seal-io/walrus/utils/files"
	"github.com/seal-io/walrus/utils/hash"
	"github.com/seal-io/walrus/utils/log"
	"github.com/seal-io/walrus/utils/netx"
	"github.com/seal-io/walrus/utils/osx"
	"github.com/seal-io/walrus/utils/strs"
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
	// Validate if run with privileged.
	if !files.Exists("/dev/kmsg") {
		if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
			return errors.New(`require "securityContext.privileged" feature of seal pod`)
		}

		return errors.New(`require "--privileged" flag to run seal container`)
	}

	// NB(thxCode): With cgroup v2,
	// in order to nest walrus multiple times,
	// we need to create specific cgroup for k3s.
	cgroupRoot := "/"

	// Enable nested cgroup v2,
	// ref to https://github.com/moby/moby/issues/43093.
	if files.Exists("/sys/fs/cgroup/cgroup.controllers") {
		// Move the processes from the root group to the /init group,
		// otherwise writing subtree_control fails with EBUSY.
		if !files.Exists("/sys/fs/cgroup/init/cgroup.procs") {
			err := files.Copy(
				"/sys/fs/cgroup/cgroup.procs",
				"/sys/fs/cgroup/init/cgroup.procs")
			if err != nil {
				return fmt.Errorf("error moving processes to root init group: %w", err)
			}
		}
		// Enable controllers to allow nesting.
		err := files.Copy(
			"/sys/fs/cgroup/cgroup.controllers",
			"/sys/fs/cgroup/cgroup.subtree_control",
			files.CopyWithModifier(func(data []byte) ([]byte, error) {
				if len(data) == 0 {
					return data, nil
				}
				cs := strings.Split(strs.FromBytes(&data), " ")
				s := "+" + strs.Join(" +", cs...)

				return strs.ToBytes(&s), nil
			}))
		if err != nil {
			return fmt.Errorf("error enabling root group controllers: %w", err)
		}

		// Calculate specific cgroup root.
		cgroupRoot = "/k3s-" + hash.SumStrings(servervars.Subnet.Get())

		// Create specific cgroup root.
		err = os.Mkdir(path.Join("/sys/fs/cgroup", cgroupRoot), 0o755)
		if err != nil && !os.IsExist(err) {
			return fmt.Errorf("error creating cgroup: %w", err)
		}

		// Enable controllers to run k3s.
		err = files.Copy(
			"/sys/fs/cgroup/cgroup.subtree_control",
			path.Join("/sys/fs/cgroup", cgroupRoot, "cgroup.subtree_control"),
			files.CopyWithModifier(func(data []byte) ([]byte, error) {
				if len(data) == 0 {
					return data, nil
				}
				cs := strings.Split(strs.FromBytes(&data), " ")
				s := "+" + strs.Join(" +", cs...)

				return strs.ToBytes(&s), nil
			}))
		if err != nil {
			return fmt.Errorf("error enabling %q group controllers: %w", cgroupRoot, err)
		}
	}

	var (
		k3sDataDir       = osx.Getenv("K3S_DATA_DIR", "/var/lib/k3s")
		k3sServerDataDir = filepath.Join(k3sDataDir, "server")
		runDataPath      = filepath.Join(consts.DataDir, "k3s")
	)

	// Link run data directory.
	err := files.Link(
		runDataPath,
		k3sServerDataDir,
		files.LinkEvenIfNotFound(false, 0o766),
		files.LinkInReplace())
	if err != nil {
		return fmt.Errorf("error link server data: %w", err)
	}

	// Reset server database.
	if files.Exists(filepath.Join(k3sServerDataDir, "db", "etcd")) {
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
	if v := servervars.Subnet.Get(); v != "" {
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
		if files.ExistsFile(clsFp) {
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
	ctx, cancel := context.WithTimeout(ctx, 300*time.Second)
	defer cancel()

	err := wait.PollUntilContextCancel(ctx, time.Second, true,
		func(ctx context.Context) (bool, error) {
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

func runK3sWith(ctx context.Context, cmdArgs []string) error {
	const cmdName = "k3s"
	logger := log.WithName(cmdName)
	logger.Infof("run: %s %s", cmdName, strs.Join(" ", cmdArgs...))
	cmd := exec.CommandContext(ctx, cmdName, cmdArgs...)
	cmd.SysProcAttr = getSysProcAttr()
	cmd.Stdout = logger.V(5)
	cmd.Stderr = logger.V(5)

	err := cmd.Run()
	if err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}
