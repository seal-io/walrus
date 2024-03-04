package extensionapi

import (
	"context"
	"fmt"
	"strings"

	autoscaling "k8s.io/api/autoscaling/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/validation/path"
	"k8s.io/apimachinery/pkg/apis/meta/internalversion"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/utils/ptr"
	ctrlcli "sigs.k8s.io/controller-runtime/pkg/client"
)

func newCreateOptionsFromUpdateOptions(in *v1.UpdateOptions) *v1.CreateOptions {
	co := &v1.CreateOptions{
		DryRun:          in.DryRun,
		FieldManager:    in.FieldManager,
		FieldValidation: in.FieldValidation,
	}
	co.TypeMeta.SetGroupVersionKind(v1.SchemeGroupVersion.WithKind("CreateOptions"))
	return co
}

func convertCtrlCreateOptionsFromMeta(in *v1.CreateOptions) ctrlcli.CreateOptions {
	return ctrlcli.CreateOptions{
		DryRun:       in.DryRun,
		FieldManager: in.FieldManager,
		Raw:          in,
	}
}

func convertCtrlUpdateOptionsFromMeta(in *v1.UpdateOptions) ctrlcli.UpdateOptions {
	return ctrlcli.UpdateOptions{
		DryRun:       in.DryRun,
		FieldManager: in.FieldManager,
		Raw:          in,
	}
}

func convertCtrlDeleteOptionsFromMeta(in *v1.DeleteOptions) ctrlcli.DeleteOptions {
	return ctrlcli.DeleteOptions{
		GracePeriodSeconds: in.GracePeriodSeconds,
		Preconditions:      in.Preconditions,
		PropagationPolicy:  in.PropagationPolicy,
		Raw:                in,
		DryRun:             in.DryRun,
	}
}

func convertCtrlGetOptionsFromMeta(in *v1.GetOptions) ctrlcli.GetOptions {
	return ctrlcli.GetOptions{
		Raw: in,
	}
}

func convertInternalListOptions(in *internalversion.ListOptions) v1.ListOptions {
	var out v1.ListOptions
	_ = internalversion.Convert_internalversion_ListOptions_To_v1_ListOptions(in, &out, nil)
	out.TypeMeta.SetGroupVersionKind(v1.SchemeGroupVersion.WithKind("ListOptions"))
	return out
}

func convertCtrlListOptionsFromMeta(ctx context.Context, in *internalversion.ListOptions) ctrlcli.ListOptions {
	ns, _ := request.NamespaceFrom(ctx)

	opts := ctrlcli.ListOptions{
		Namespace: ns,
		Limit:     in.Limit,
		Continue:  in.Continue,
		Raw:       ptr.To(convertInternalListOptions(in)),
	}
	if in.FieldSelector != nil && !in.FieldSelector.Empty() {
		opts.FieldSelector = in.FieldSelector
	}
	if in.LabelSelector != nil && !in.LabelSelector.Empty() {
		opts.LabelSelector = in.LabelSelector
	}
	return opts
}

func keyFuncForNamespacedScope(ctx context.Context, name string) (types.NamespacedName, error) {
	ns, ok := request.NamespaceFrom(ctx)
	if !ok || len(ns) == 0 {
		return types.NamespacedName{},
			kerrors.NewBadRequest("Namespace parameter required.")
	}
	if len(name) == 0 {
		return types.NamespacedName{},
			kerrors.NewBadRequest("Name parameter required.")
	}
	if msgs := path.IsValidPathSegmentName(name); len(msgs) != 0 {
		return types.NamespacedName{},
			kerrors.NewBadRequest(fmt.Sprintf("Name parameter invalid: %q: %s", name, strings.Join(msgs, ";")))
	}
	return types.NamespacedName{Name: name, Namespace: ns}, nil
}

func keyFuncForClusterScope(ctx context.Context, name string) (types.NamespacedName, error) {
	if len(name) == 0 {
		return types.NamespacedName{},
			kerrors.NewBadRequest("Name parameter required.")
	}
	if msgs := path.IsValidPathSegmentName(name); len(msgs) != 0 {
		return types.NamespacedName{},
			kerrors.NewBadRequest(fmt.Sprintf("Name parameter invalid: %q: %s", name, strings.Join(msgs, ";")))
	}
	return types.NamespacedName{Name: name}, nil
}

func beforeUpdateFuncForPreventStatusModify(ctx context.Context, newObj, oldObj runtime.Object) (runtime.Object, error) {
	if oldWithStatus, ok := oldObj.(ObjectWithStatusSubResource); ok {
		oldWithStatus.CopyStatusTo(newObj)
		return newObj, nil
	}
	return nil, kerrors.NewBadRequest("do not support status subresource")
}

func beforeUpdateFuncForStatusModifyOnly(ctx context.Context, newObj, oldObj runtime.Object) (runtime.Object, error) {
	if newWithStatus, ok := newObj.(ObjectWithStatusSubResource); ok {
		newWithStatus.CopyStatusTo(oldObj)
		return oldObj, nil
	}
	return nil, kerrors.NewBadRequest("do not support status subresource")
}

type _ScaleUpdatedObjectInfo struct {
	rest.UpdatedObjectInfo
}

func toScaleUpdatedObjectInfo(info rest.UpdatedObjectInfo) rest.UpdatedObjectInfo {
	return _ScaleUpdatedObjectInfo{info}
}

func (s _ScaleUpdatedObjectInfo) UpdatedObject(ctx context.Context, oldObj runtime.Object) (runtime.Object, error) {
	withScale, ok := oldObj.(ObjectWithScaleSubResource)
	if !ok {
		return nil, kerrors.NewBadRequest("do not support scale subresource")
	}

	oldScale := withScale.GetScale()
	obj, err := s.UpdatedObjectInfo.UpdatedObject(ctx, oldScale)
	if err != nil {
		return nil, err
	}
	if obj == nil {
		return nil, kerrors.NewBadRequest("returned nil updated object")
	}

	newScale, ok := obj.(*autoscaling.Scale)
	if !ok {
		return nil, kerrors.NewBadRequest(fmt.Sprintf("wrong object passed to Scale update: %v", obj))
	}
	withScale.SetScale(newScale)
	if len(newScale.ResourceVersion) != 0 {
		withScale.GetObjectMeta().SetResourceVersion(newScale.ResourceVersion)
	}

	return withScale, nil
}

func toScaleCreateValidation(fn rest.ValidateObjectFunc) rest.ValidateObjectFunc {
	return func(ctx context.Context, obj runtime.Object) error {
		withScale, ok := obj.(ObjectWithScaleSubResource)
		if !ok {
			return kerrors.NewBadRequest("do not support scale subresource")
		}
		return fn(ctx, withScale.GetScale())
	}
}
