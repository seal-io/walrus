package platformk8s

import (
	"context"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	dynamicclient "k8s.io/client-go/dynamic"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/platformk8s/intercept"
	"github.com/seal-io/seal/pkg/platformk8s/kube"
)

// Label implements operator.Operator.
func (op Operator) Label(ctx context.Context, res *model.ApplicationResource, labels map[string]string) error {
	if res == nil || res.DeployerType != types.DeployerTypeTF {
		return nil
	}

	client, err := dynamicclient.NewForConfig(op.RestConfig)
	if err != nil {
		return fmt.Errorf("error creating kubernetes core client: %w", err)
	}

	gvr, ok := intercept.Terraform().GetGVR(res.Type)
	if !ok {
		return fmt.Errorf("error get resource %s's gvr", res.ID)
	}

	ns, name := kube.ParseNamespacedName(res.Name)
	obj, err := client.Resource(gvr).Namespace(ns).Get(ctx, name, meta.GetOptions{ResourceVersion: "0"})
	if err != nil {
		if apierrors.IsNotFound(err) {
			// skip apply labels while resource isn't existed
			return nil
		}
		return err
	}
	return kube.Label(ctx, client, obj, labels)
}
