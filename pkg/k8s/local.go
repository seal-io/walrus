package k8s

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/seal-io/walrus/pkg/consts"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/utils/files"
	"github.com/seal-io/walrus/utils/log"
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

	// Enable nested cgroup v2,
	// ref to https://github.com/moby/moby/issues/43093.
	if files.Exists("/sys/fs/cgroup/cgroup.controllers") {
		// Move the processes from the root group to the /init group,
		// otherwise writing subtree_control fails with EBUSY.
		err := files.Copy(
			"/sys/fs/cgroup/cgroup.procs",
			"/sys/fs/cgroup/init/cgroup.procs")
		if err != nil {
			return fmt.Errorf("error moving processes to init group: %w", err)
		}
		// Enable controllers.
		err = files.Copy(
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
			return fmt.Errorf("error enabling group controllers: %w", err)
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

	// Reset server data.
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
		"--kubelet-arg=system-reserved=cpu=300m,memory=256Mi",
		"--kubelet-arg=kube-reserved=cpu=200m,memory=256Mi",
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

	err = Wait(ctx, cfg,
		prepareDeployWorkspace,
		prepareDeployPermission)
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

// prepareDeployWorkspace creates the kubernetes namespace to run deployer at next.
func prepareDeployWorkspace(ctx context.Context, cli *kubernetes.Clientset) error {
	ns := core.Namespace{
		ObjectMeta: meta.ObjectMeta{
			Name: types.WalrusSystemNamespace,
		},
	}

	_, err := cli.CoreV1().
		Namespaces().
		Create(ctx, &ns, meta.CreateOptions{})
	if err != nil && !kerrors.IsAlreadyExists(err) {
		return err
	}

	return nil
}

// prepareDeployPermission creates the kubernetes rbac resources for deployer running.
func prepareDeployPermission(ctx context.Context, cli *kubernetes.Clientset) error {
	// FIXME(thxCode): Remove this later,
	//  we should not use the same permission for deployer running,
	//  deployer should use the permission of the configured connectors.
	sa := core.ServiceAccount{
		ObjectMeta: meta.ObjectMeta{
			Namespace: types.WalrusSystemNamespace,
			Name:      types.DeployerServiceAccountName,
		},
	}

	_, err := cli.CoreV1().
		ServiceAccounts(sa.Namespace).
		Create(ctx, &sa, meta.CreateOptions{})
	if err != nil && !kerrors.IsAlreadyExists(err) {
		return err
	}

	r := rbac.Role{
		ObjectMeta: meta.ObjectMeta{
			Namespace: types.WalrusSystemNamespace,
			Name:      types.DeployerServiceAccountName,
		},
		Rules: []rbac.PolicyRule{
			{
				APIGroups: []string{batch.GroupName},
				Resources: []string{"jobs"},
				Verbs:     []string{rbac.VerbAll},
			},
			{
				APIGroups: []string{core.GroupName},
				Resources: []string{"secrets", "pods", "pods/log"},
				Verbs:     []string{rbac.VerbAll},
			},
		},
	}

	_, err = cli.RbacV1().
		Roles(r.Namespace).
		Create(ctx, &r, meta.CreateOptions{})
	if err != nil && !kerrors.IsAlreadyExists(err) {
		return err
	}

	rb := rbac.RoleBinding{
		ObjectMeta: meta.ObjectMeta{
			Namespace: types.WalrusSystemNamespace,
			Name:      types.DeployerServiceAccountName,
		},
		Subjects: []rbac.Subject{
			{
				Kind: rbac.ServiceAccountKind,
				Name: types.DeployerServiceAccountName,
			},
		},
		RoleRef: rbac.RoleRef{
			APIGroup: rbac.GroupName,
			Kind:     "Role",
			Name:     types.DeployerServiceAccountName,
		},
	}

	_, err = cli.RbacV1().
		RoleBindings(r.Namespace).
		Create(ctx, &rb, meta.CreateOptions{})
	if err != nil && !kerrors.IsAlreadyExists(err) {
		return err
	}

	return nil
}
