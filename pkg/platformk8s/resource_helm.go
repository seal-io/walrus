package platformk8s

import (
	"context"

	"helm.sh/helm/v3/pkg/releaseutil"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/platformk8s/helm"
	"github.com/seal-io/seal/pkg/platformk8s/polymorphic"
	"github.com/seal-io/seal/utils/strs"
)

// parseResourcesOfHelm parses the given `helm_release` model.ApplicationResource,
// and keeps resource item which matches.
func parseResourcesOfHelm(ctx context.Context, op Operator, res *model.ApplicationResource, match func(schema.GroupVersionKind) bool) ([]resource, error) {
	if match == nil {
		return nil, nil
	}

	var opts = helm.GetReleaseOptions{
		RESTClientGetter: helm.IncompleteRestClientGetter(*op.RestConfig),
		Log:              op.Logger.Debugf,
	}
	var hr, err = helm.GetRelease(ctx, res, opts)
	if err != nil {
		return nil, resourceParsingError("error parsing helm release: " + err.Error())
	}

	// parse helm release.
	var rs []resource
	var hms = releaseutil.SplitManifests(hr.Manifest)
	if len(hms) == 0 {
		return nil, nil
	}
	var hs = polymorphic.YamlSerializer()
	for k := range hms {
		var (
			obj unstructured.Unstructured
			gvk *schema.GroupVersionKind
		)
		var ms = hms[k]
		_, gvk, err = hs.Decode(strs.ToBytes(&ms), nil, &obj)
		if err != nil {
			op.Logger.Warnf("error decoding helm release resource: %v", err)
			continue
		}
		if !match(*gvk) {
			continue
		}
		var gvr, _ = meta.UnsafeGuessKindToResource(*gvk)
		rs = append(rs, resource{
			GroupVersionResource: gvr,
			Namespace:            obj.GetNamespace(),
			Name:                 obj.GetName(),
		})
	}
	return rs, nil
}
