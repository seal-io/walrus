package walruscore

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrlwebhook "sigs.k8s.io/controller-runtime/pkg/webhook"

	walruscore "github.com/seal-io/walrus/pkg/apis/walruscore/v1"
	"github.com/seal-io/walrus/pkg/webhook"
)

// ResourceDefinitionWebhook hooks a v1.ResourceDefinition object.
//
// +k8s:webhook-gen:mutating:group="walruscore.seal.io",version="v1",resource="resourcedefinitions",scope="Namespaced"
// +k8s:webhook-gen:mutating:operations=["CREATE","UPDATE"],failurePolicy="Fail",sideEffects="None",matchPolicy="Equivalent",timeoutSeconds=10
type ResourceDefinitionWebhook struct{}

func (r *ResourceDefinitionWebhook) SetupWebhook(_ context.Context, opts webhook.SetupOptions) (runtime.Object, error) {
	return &walruscore.ResourceDefinition{}, nil
}

var _ ctrlwebhook.CustomDefaulter = (*ResourceDefinitionWebhook)(nil)

func (r *ResourceDefinitionWebhook) Default(ctx context.Context, obj runtime.Object) error {
	// TODO: your logic here

	return nil
}
