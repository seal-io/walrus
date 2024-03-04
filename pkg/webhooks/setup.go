package webhooks

import (
	"context"
	"fmt"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	admreg "k8s.io/api/admissionregistration/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlwebhook "sigs.k8s.io/controller-runtime/pkg/webhook"
	ctrladmission "sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/seal-io/walrus/pkg/clients/clientset"
	"github.com/seal-io/walrus/pkg/kubeclientset"
	"github.com/seal-io/walrus/pkg/kubeclientset/review"
	"github.com/seal-io/walrus/pkg/webhook"
	"github.com/seal-io/walrus/pkg/webhooks/walruscore"
)

// NB(thxCode): Register webhooks below.
var (
	setupers = []webhook.Setup{
		new(walruscore.CatalogWebhook),
		new(walruscore.ConnectorWebhook),
		new(walruscore.ResourceDefinitionWebhook),
		new(walruscore.TemplateWebhook),
	}
	cfgGetters = []_WebhookConfigurationsGetter{
		walruscore.GetWebhookConfigurations,
	}
)

type (
	_DefaultWebhookHandler interface {
		ctrlwebhook.CustomDefaulter
		DefaultPath() string
	}
	_ValidatorWebhookHandler interface {
		ctrlwebhook.CustomValidator
		ValidatePath() string
	}
	_WebhookConfigurationsGetter func(string, admreg.WebhookClientConfig) (
		*admreg.ValidatingWebhookConfiguration,
		*admreg.MutatingWebhookConfiguration,
	)
	HTTPServeMux interface {
		Handle(string, http.Handler)
	}
)

func Setup(ctx context.Context, mgr ctrl.Manager, mux HTTPServeMux) error {
	scheme := mgr.GetScheme()
	for i := range setupers {
		switch setupers[i].(type) {
		default:
			continue
		case _DefaultWebhookHandler:
		case _ValidatorWebhookHandler:
		}

		opts := webhook.SetupOptions{Manager: mgr}
		obj, err := setupers[i].SetupWebhook(ctx, opts)
		if err != nil {
			return fmt.Errorf("webhook setup: %s: %w", spew.Sdump(setupers[i]), err)
		}

		if d, ok := setupers[i].(_DefaultWebhookHandler); ok {
			dh := ctrladmission.WithCustomDefaulter(scheme, obj, d).WithRecoverPanic(true)
			mux.Handle(d.DefaultPath(), dh)
		}
		if v, ok := setupers[i].(_ValidatorWebhookHandler); ok {
			vh := ctrladmission.WithCustomValidator(scheme, obj, v).WithRecoverPanic(true)
			mux.Handle(v.ValidatePath(), vh)
		}
	}

	return nil
}

// GetWebhookConfigurations returns the registered webhook configuration.
func GetWebhookConfigurations(cc admreg.WebhookClientConfig) (*admreg.ValidatingWebhookConfiguration, *admreg.MutatingWebhookConfiguration) {
	// NB(thxCode): add more webhook configurations getters here.
	// Merge all the webhook configurations from the getters.
	var (
		vret = make([]*admreg.ValidatingWebhookConfiguration, len(cfgGetters))
		vwsc int
		mret = make([]*admreg.MutatingWebhookConfiguration, len(cfgGetters))
		mwsc int
	)
	for i := range cfgGetters {
		vwc, mwc := cfgGetters[i]("walrus-webhook", cc)
		if vwc != nil {
			vret[i] = vwc
			vwsc += len(vwc.Webhooks)
		}
		if mwc != nil {
			mret[i] = mwc
			mwsc += len(mwc.Webhooks)
		}
	}

	var (
		vwc *admreg.ValidatingWebhookConfiguration
		mwc *admreg.MutatingWebhookConfiguration
	)
	if vwsc != 0 {
		vwc = &admreg.ValidatingWebhookConfiguration{
			ObjectMeta: meta.ObjectMeta{
				Name: "walrus-webhook-validation",
			},
		}
		for i := range vret {
			if vret[i] == nil {
				continue
			}
			vwc.Webhooks = append(vwc.Webhooks, vret[i].Webhooks...)
		}
	}
	if mwsc != 0 {
		mwc = &admreg.MutatingWebhookConfiguration{
			ObjectMeta: meta.ObjectMeta{
				Name: "walrus-webhook-mutation",
			},
		}
		for i := range mret {
			if mret[i] == nil {
				continue
			}
			mwc.Webhooks = append(mwc.Webhooks, mret[i].Webhooks...)
		}
	}

	return vwc, mwc
}

// InstallWebhookConfigurations installs the webhook configurations.
func InstallWebhookConfigurations(ctx context.Context, cli clientset.Interface, cc admreg.WebhookClientConfig) error {
	err := review.CanDoUpdate(ctx,
		cli.AuthorizationV1().SelfSubjectAccessReviews(),
		review.Simples{
			{
				Group:    admreg.SchemeGroupVersion.Group,
				Version:  admreg.SchemeGroupVersion.Version,
				Resource: "validatingwebhookconfigurations",
			},
			{
				Group:    admreg.SchemeGroupVersion.Group,
				Version:  admreg.SchemeGroupVersion.Version,
				Resource: "mutatingwebhookconfigurations",
			},
		},
		review.WithCreateIfNotExisted(),
	)
	if err != nil {
		return err
	}

	vwCli := cli.AdmissionregistrationV1().ValidatingWebhookConfigurations()
	mwCli := cli.AdmissionregistrationV1().MutatingWebhookConfigurations()

	vwc, mwc := GetWebhookConfigurations(cc)
	if vwc != nil {
		_, err := kubeclientset.Update(ctx, vwCli, vwc,
			kubeclientset.WithCreateIfNotExisted[*admreg.ValidatingWebhookConfiguration]())
		if err != nil {
			return fmt.Errorf("install validating webhook configuration %q: %w",
				vwc.GetName(), err)
		}
	}
	if mwc != nil {
		_, err := kubeclientset.Update(ctx, mwCli, mwc,
			kubeclientset.WithCreateIfNotExisted[*admreg.MutatingWebhookConfiguration]())
		if err != nil {
			return fmt.Errorf("install mutating webhook configuration %q: %w",
				mwc.GetName(), err)
		}
	}

	return nil
}
