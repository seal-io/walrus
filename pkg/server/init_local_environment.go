package server

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	dtypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	envbus "github.com/seal-io/walrus/pkg/bus/environment"
	"github.com/seal-io/walrus/pkg/dao"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/connector"
	"github.com/seal-io/walrus/pkg/dao/model/environment"
	"github.com/seal-io/walrus/pkg/dao/model/project"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/crypto"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/pkg/settings"
)

const (
	// LocalEnvironmentKubernetes indicates creating local environment with management kubernetes.
	localEnvironmentModeKubernetes = "kubernetes"
	// LocalEnvironmentDocker indicates creating local environment with docker.
	localEnvironmentModeDocker = "docker"
	// LocalEnvironmentModeDisabled disables local environment creation.
	localEnvironmentModeDisabled = "disabled"
)

func (r *Server) createLocalEnvironment(ctx context.Context, opts initOptions) error {
	localEnvironmentMode, err := settings.LocalEnvironmentMode.Value(ctx, opts.ModelClient)
	if err != nil {
		return err
	}

	if localEnvironmentMode == localEnvironmentModeDisabled {
		return nil
	}

	return opts.ModelClient.WithTx(ctx, func(tx *model.Tx) error {
		defaultProject, err := tx.Projects().Query().
			Where(project.Name("default")).
			Only(ctx)
		if err != nil {
			return err
		}

		var conn *model.Connector

		switch localEnvironmentMode {
		case localEnvironmentModeDocker:
			conn, err = r.applyDockerConnector(ctx, tx, defaultProject.ID)
			if err != nil {
				return err
			}

			if err = applyLocalDockerNetwork(ctx); err != nil {
				return err
			}
		case localEnvironmentModeKubernetes:
			conn, err = r.applyKubernetesConnector(ctx, tx, defaultProject.ID)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("invalid local environment mode %q", localEnvironmentMode)
		}

		_, err = tx.Environments().Query().
			Where(
				environment.Name("local"),
				environment.ProjectID(defaultProject.ID),
			).Only(ctx)
		if err == nil {
			return nil
		} else if !model.IsNotFound(err) {
			return err
		}

		env := &model.Environment{
			Name:      "local",
			ProjectID: defaultProject.ID,
			Type:      types.EnvironmentDevelopment,
			Edges: model.EnvironmentEdges{
				Connectors: []*model.EnvironmentConnectorRelationship{
					{
						ConnectorID: conn.ID,
					},
				},
			},
		}

		env, err = tx.Environments().Create().
			Set(env).
			SaveE(ctx, dao.EnvironmentConnectorsEdgeSave)
		if err != nil {
			return err
		}

		return envbus.NotifyIDs(ctx, tx, envbus.EventCreate, env.ID)
	})
}

func (r *Server) applyKubernetesConnector(
	ctx context.Context,
	mc model.ClientSet,
	projectID object.ID,
) (*model.Connector, error) {
	kubeConfig, err := r.readKubeConfig()
	if err != nil {
		return nil, err
	}

	conn, err := mc.Connectors().Query().
		Where(
			connector.Name("local"),
			connector.ProjectID(projectID),
		).
		Only(ctx)

	switch {
	case model.IsNotFound(err):
		conn = &model.Connector{
			Name:                      "local",
			Type:                      types.ConnectorTypeKubernetes,
			ProjectID:                 projectID,
			ApplicableEnvironmentType: types.EnvironmentDevelopment,
			Category:                  types.ConnectorCategoryKubernetes,
			ConfigVersion:             "v1",
			ConfigData: crypto.Properties{
				"kubeconfig": crypto.StringProperty(kubeConfig),
			},
		}

		conn, err = mc.Connectors().Create().
			Set(conn).
			Save(ctx)
		if err != nil {
			return nil, err
		}
	case err != nil:
		return nil, err
	default:
		// Update config data of existing local kubernetes connector as it may change on restart.
		conn.ConfigData = crypto.Properties{
			"kubeconfig": crypto.StringProperty(kubeConfig),
		}

		conn, err = mc.Connectors().UpdateOne(conn).
			Set(conn).
			Save(ctx)
		if err != nil {
			return nil, err
		}
	}

	return conn, nil
}

func (r *Server) readKubeConfig() (string, error) {
	inClusterConfig, err := rest.InClusterConfig()
	if err == nil {
		caData, err := os.ReadFile(inClusterConfig.TLSClientConfig.CAFile)
		if err != nil {
			return "", err
		}

		kc := clientcmdapi.Config{
			Clusters: map[string]*clientcmdapi.Cluster{
				"default": {
					Server:                   inClusterConfig.Host,
					CertificateAuthorityData: caData,
				},
			},
			Contexts: map[string]*clientcmdapi.Context{
				"default": {
					Cluster:  "default",
					AuthInfo: "default",
				},
			},
			CurrentContext: "default",
			AuthInfos: map[string]*clientcmdapi.AuthInfo{
				"default": {
					Token: inClusterConfig.BearerToken,
				},
			},
		}
		kcData, err := clientcmd.Write(kc)

		return string(kcData), err
	}

	kubeConfigPath := r.KubeConfig
	if kubeConfigPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		kubeConfigPath = filepath.Join(home, clientcmd.RecommendedHomeDir, clientcmd.RecommendedFileName)
	}

	kubeConfigData, err := os.ReadFile(kubeConfigPath)
	if err != nil {
		return "", err
	}

	return string(kubeConfigData), nil
}

func (r *Server) applyDockerConnector(
	ctx context.Context,
	mc model.ClientSet,
	projectID object.ID,
) (*model.Connector, error) {
	conn, err := mc.Connectors().Query().
		Where(
			connector.Name("local"),
			connector.ProjectID(projectID),
		).
		Only(ctx)

	if model.IsNotFound(err) {
		conn = &model.Connector{
			Name:                      "local",
			Type:                      types.ConnectorTypeDocker,
			ProjectID:                 projectID,
			ApplicableEnvironmentType: types.EnvironmentDevelopment,
			Category:                  types.ConnectorCategoryDocker,
			ConfigVersion:             "v1",
			ConfigData:                crypto.Properties{},
		}

		conn, err = mc.Connectors().Create().
			Set(conn).
			Save(ctx)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return conn, nil
}

func applyLocalDockerNetwork(ctx context.Context) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	networkName := "local-walrus"

	networks, err := cli.NetworkList(ctx, dtypes.NetworkListOptions{})
	if err != nil {
		return err
	}

	exists := false

	for _, n := range networks {
		if n.Name == networkName {
			exists = true
			break
		}
	}

	if !exists {
		_, err = cli.NetworkCreate(ctx, networkName, dtypes.NetworkCreate{
			Driver: "bridge",
		})
		if err != nil {
			return err
		}
	}

	return nil
}
