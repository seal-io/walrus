package systemapp

import (
	"context"
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/seal-io/utils/funcx"
	"github.com/seal-io/utils/pools/gopool"
	"github.com/seal-io/utils/stringx"
	"golang.org/x/exp/maps"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/rest"

	"github.com/seal-io/walrus/pkg/systemapp/helm"
	"github.com/seal-io/walrus/pkg/systemkuberes"
	"github.com/seal-io/walrus/pkg/systemsetting"
)

// NB(thxCode): Register installer below.
var installers = []_Installer{
	installMinio,
	installHermitCrab,
}

type _Installer func(context.Context, *helm.Client, map[string]any, sets.Set[string]) error

// Install installs the system applications.
func Install(ctx context.Context, cliConfig rest.Config, disable sets.Set[string]) error {
	hc, err := helm.NewClient(&cliConfig, helm.WithNamespace(systemkuberes.SystemNamespaceName))
	if err != nil {
		return fmt.Errorf("create helm client: %w", err)
	}

	gvc := map[string]any{
		"ImageRegistry": funcx.NoError(systemsetting.ImageRegistry.Value(ctx)),
		"Env":           getCommonEnv(ctx),
		"ManagedLabel":  "walrus.seal.io/managed",
	}

	gp := gopool.Group()
	for i := range installers {
		gvc := maps.Clone(gvc)
		in := installers[i]
		gp.Go(func() error {
			err = in(ctx, hc, gvc, disable)
			if err != nil {
				return fmt.Errorf("%s: %w", loadInstallerName(in), err)
			}
			return nil
		})
	}
	return gp.Wait()
}

func getCommonEnv(ctx context.Context) (env []core.EnvVar) {
	if v := funcx.NoError(systemsetting.DeployerAllProxy.Value(ctx)); v != "" {
		env = append(env, core.EnvVar{
			Name:  "ALL_PROXY",
			Value: v,
		})
	}

	if v := funcx.NoError(systemsetting.DeployerHttpProxy.Value(ctx)); v != "" {
		env = append(env, core.EnvVar{
			Name:  "HTTP_PROXY",
			Value: v,
		})
	}

	if v := funcx.NoError(systemsetting.DeployerHttpsProxy.Value(ctx)); v != "" {
		env = append(env, core.EnvVar{
			Name:  "HTTPS_PROXY",
			Value: v,
		})
	}

	if v := funcx.NoError(systemsetting.DeployerNoProxy.Value(ctx)); v != "" {
		env = append(env, core.EnvVar{
			Name:  "NO_PROXY",
			Value: v,
		})
	}

	return env
}

func loadInstallerName(i _Installer) string {
	n := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	n = strings.TrimPrefix(strings.TrimSuffix(filepath.Ext(n), "-fm"), ".")
	return stringx.Decamelize(n, true)
}
