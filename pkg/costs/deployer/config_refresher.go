package deployer

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/status"
	"github.com/seal-io/seal/pkg/platformk8s"
	"github.com/seal-io/seal/utils/log"
)

func UpdateCustomPricing(ctx context.Context, conn *model.Connector) error {
	log.WithName("cost").Debugf("updating cost custom pricing for connector %s", conn.Name)

	restCfg, err := platformk8s.GetConfig(*conn)
	if err != nil {
		return err
	}

	if err = updateCustomPricingConfigMap(ctx, conn, restCfg); err != nil {
		return err
	}

	if !status.ConnectorStatusReady.IsTrue(conn) {
		return nil
	}

	if err = refreshCustomPricing(conn, restCfg); err != nil {
		return err
	}
	return nil
}

func updateCustomPricingConfigMap(ctx context.Context, conn *model.Connector, restCfg *rest.Config) error {
	corev1Client, err := corev1.NewForConfig(restCfg)
	if err != nil {
		return fmt.Errorf("error creating kubernetes core client: %w", err)
	}

	configMaps := corev1Client.ConfigMaps(types.SealSystemNamespace)
	current := opencostCustomPricingConfigMap(conn)
	existed, err := configMaps.Get(ctx, ConfigMapNameOpencost, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			_, err = configMaps.Create(ctx, current, metav1.CreateOptions{})
			if err != nil && apierrors.IsAlreadyExists(err) {
				return fmt.Errorf("error create configmap %s:%s, %w", types.SealSystemNamespace, ConfigMapNameOpencost, err)
			}
			return nil
		}

		return fmt.Errorf("error get configmap %s:%s, %w", types.SealSystemNamespace, ConfigMapNameOpencost, err)
	}

	if reflect.DeepEqual(existed.Data, current.Data) {
		return nil
	}

	existed.Data = current.Data
	_, err = configMaps.Update(ctx, existed, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("error update configmap %s:%s, %w", types.SealSystemNamespace, ConfigMapNameOpencost, err)
	}

	return nil
}

func refreshCustomPricing(conn *model.Connector, restCfg *rest.Config) error {
	// call opencost API to refresh pricing
	url, err := opencostRefreshPricingURL(restCfg)
	if err != nil {
		return fmt.Errorf("error get refresh princing url for connector: %s: %w", conn.Name, err)
	}

	clusterClient, err := rest.HTTPClientFor(restCfg)
	if err != nil {
		return err
	}

	resp, err := clusterClient.Post(url, "application/json", nil)
	if err != nil {
		return fmt.Errorf("error request to %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error response from %s, expected %d but get code: %d", url, http.StatusOK, resp.StatusCode)
	}
	return nil
}
