package walruscore

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrlwebhook "sigs.k8s.io/controller-runtime/pkg/webhook"
	ctrladmission "sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	walruscore "github.com/seal-io/walrus/pkg/apis/walruscore/v1"
	"github.com/seal-io/walrus/pkg/webhook"
)

// TemplateWebhook hooks a v1.Template object.
//
// +k8s:webhook-gen:validating:group="walruscore.seal.io",version="v1",resource="templates",scope="Namespaced"
// +k8s:webhook-gen:validating:operations=["UPDATE","DELETE"],failurePolicy="Fail",sideEffects="None",matchPolicy="Equivalent",timeoutSeconds=10
type TemplateWebhook struct {
	webhook.DefaultCustomValidator
}

func (r *TemplateWebhook) SetupWebhook(_ context.Context, opts webhook.SetupOptions) (runtime.Object, error) {
	return &walruscore.Template{}, nil
}

var _ ctrlwebhook.CustomValidator = (*TemplateWebhook)(nil)

func (r *TemplateWebhook) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (ctrladmission.Warnings, error) {
	// TODO: your logic here

	return nil, nil
}

func (r *TemplateWebhook) ValidateDelete(ctx context.Context, obj runtime.Object) (ctrladmission.Warnings, error) {
	// TODO: your logic here

	return nil, nil
}
