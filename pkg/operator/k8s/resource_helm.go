package k8s

import (
	"context"

	"helm.sh/helm/v3/pkg/releaseutil"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/seal-io/walrus/pkg/operator/k8s/helm"
	"github.com/seal-io/walrus/pkg/operator/k8s/polymorphic"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/utils/strs"
)

// parseResourcesOfHelm parses the given `helm_release` model.ServiceResource,
// and keeps resource item which matches.
func parseResourcesOfHelm(
	ctx context.Context,
	op Operator,
	res *model.ServiceResource,
	match func(schema.GroupVersionKind) bool,
) ([]resource, error) {
	if match == nil {
		return nil, nil
	}

	opts := helm.GetReleaseOptions{
		RESTClientGetter: helm.IncompleteRestClientGetter(*op.RestConfig),
		Log:              op.Logger.Debugf,
	}

	hr, err := helm.GetRelease(ctx, res, opts)
	if err != nil {
		return nil, resourceParsingError("error parsing helm release: " + err.Error())
	}

	hms := releaseutil.SplitManifests(hr.Manifest)
	if len(hms) == 0 {
		return nil, nil
	}
	hs := polymorphic.YamlSerializer()
	// Parse helm release.
	rs := make([]resource, 0, len(hms))

	for k := range hms {
		var (
			obj unstructured.Unstructured
			gvk *schema.GroupVersionKind
		)
		ms := hms[k]

		_, gvk, err = hs.Decode(strs.ToBytes(&ms), nil, &obj)
		if err != nil {
			op.Logger.Warnf("error decoding helm release resource: %v", err)
			continue
		}

		if !match(*gvk) {
			continue
		}
		gvr, _ := meta.UnsafeGuessKindToResource(*gvk)

		ns := obj.GetNamespace()
		if ns == "" {
			ns = hr.Namespace
		}

		rs = append(rs, resource{
			GroupVersionResource: gvr,
			Namespace:            ns,
			Name:                 obj.GetName(),
		})
	}

	return rs, nil
}
