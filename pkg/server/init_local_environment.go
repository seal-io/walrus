package server

import (
	"context"
	"os"
	"path/filepath"

	"k8s.io/client-go/tools/clientcmd"

	envbus "github.com/seal-io/walrus/pkg/bus/environment"
	"github.com/seal-io/walrus/pkg/dao"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/connector"
	"github.com/seal-io/walrus/pkg/dao/model/environment"
	"github.com/seal-io/walrus/pkg/dao/model/project"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/crypto"
	"github.com/seal-io/walrus/pkg/settings"
)

const (
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
		kubeConfig, err := r.readKubeConfig()
		if err != nil {
			return err
		}

		defaultProject, err := tx.Projects().Query().
			Where(project.Name("default")).
			Only(ctx)
		if err != nil {
			return err
		}

		conn, err := tx.Connectors().Query().
			Where(
				connector.Name("local"),
				connector.ProjectID(defaultProject.ID),
			).
			Only(ctx)

		switch {
		case model.IsNotFound(err):
			conn = &model.Connector{
				Name:                      "local",
				Type:                      types.ConnectorTypeKubernetes,
				ProjectID:                 defaultProject.ID,
				ApplicableEnvironmentType: types.EnvironmentDevelopment,
				Category:                  types.ConnectorCategoryKubernetes,
				ConfigVersion:             "v1",
				ConfigData: crypto.Properties{
					"kubeconfig": crypto.StringProperty(kubeConfig),
				},
			}

			if os.Getenv("KUBERNETES_SERVICE_HOST") == "" && os.Getenv("_RUNNING_INSIDE_CONTAINER_") != "" {
				// Set label for embedded k3s.
				conn.Labels = map[string]string{
					types.LabelEmbeddedKubernetes: "true",
				}
			}

			conn, err = tx.Connectors().Create().
				Set(conn).
				Save(ctx)
			if err != nil {
				return err
			}
		case err != nil:
			return err
		default:
			// Update config data of existing local kubernetes connector as it may change on restart.
			conn.ConfigData = crypto.Properties{
				"kubeconfig": crypto.StringProperty(kubeConfig),
			}

			conn, err = tx.Connectors().UpdateOne(conn).
				Set(conn).
				Save(ctx)
			if err != nil {
				return err
			}
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

func (r *Server) readKubeConfig() (string, error) {
	kubeConfig := r.KubeConfig
	if kubeConfig == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		kubeConfig = filepath.Join(home, clientcmd.RecommendedHomeDir, clientcmd.RecommendedFileName)
	}

	kubeConfigData, err := os.ReadFile(kubeConfig)
	if err != nil {
		return "", err
	}

	return string(kubeConfigData), nil
}
