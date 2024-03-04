package extensionapi

import (
	"context"

	autoscaling "k8s.io/api/autoscaling/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/rest"
)

type (
	StatusSubResourceParentStore interface {
		rest.Storage
		rest.Getter
		rest.Updater
		// GetStatusSubResourceUpdater returns an Updater for the status subresource.
		GetStatusSubResourceUpdater() UpdateOperation
	}

	StatusSubResourceStorage struct {
		StatusSubResourceParentStore
	}
)

var (
	_ rest.Storage = (*StatusSubResourceStorage)(nil)
	_ rest.Getter  = (*StatusSubResourceStorage)(nil)
	_ rest.Updater = (*StatusSubResourceStorage)(nil)
)

// AsStatusSubResourceStorage returns a default status SubResourceStorage for the given SubResourceObject.
func AsStatusSubResourceStorage(s StatusSubResourceParentStore) StatusSubResourceStorage {
	return StatusSubResourceStorage{StatusSubResourceParentStore: s}
}

func (s StatusSubResourceStorage) Update(
	ctx context.Context,
	name string,
	objInfo rest.UpdatedObjectInfo,
	createValidation rest.ValidateObjectFunc,
	updateValidation rest.ValidateObjectUpdateFunc,
	forceAllowCreate bool,
	options *meta.UpdateOptions,
) (runtime.Object, bool, error) {
	return s.GetStatusSubResourceUpdater().Update(ctx, name, objInfo, createValidation, updateValidation, forceAllowCreate, options)
}

type (
	ScaleSubResourceParentStore interface {
		rest.Getter
		rest.Updater
	}

	ScaleSubResourceStorage struct {
		ScaleSubResourceParentStore
	}
)

var (
	_ rest.Storage = (*ScaleSubResourceStorage)(nil)
	_ rest.Getter  = (*ScaleSubResourceStorage)(nil)
	_ rest.Updater = (*ScaleSubResourceStorage)(nil)
)

// AsScaleSubResourceStorage returns a default scale SubResourceStorage for the given SubResourceObject.
func AsScaleSubResourceStorage(s ScaleSubResourceParentStore) ScaleSubResourceStorage {
	return ScaleSubResourceStorage{ScaleSubResourceParentStore: s}
}

func (s ScaleSubResourceStorage) New() runtime.Object {
	return &autoscaling.Scale{}
}

func (s ScaleSubResourceStorage) Destroy() {
}

func (s ScaleSubResourceStorage) Get(
	ctx context.Context,
	name string,
	options *meta.GetOptions,
) (runtime.Object, error) {
	obj, err := s.ScaleSubResourceParentStore.Get(ctx, name, options)
	if err != nil {
		return nil, err
	}

	withScale, ok := obj.(ObjectWithScaleSubResource)
	if !ok {
		return nil, kerrors.NewBadRequest("do not support scale subresource")
	}
	return withScale.GetScale(), nil
}

func (s ScaleSubResourceStorage) Update(
	ctx context.Context,
	name string,
	objInfo rest.UpdatedObjectInfo,
	createValidation rest.ValidateObjectFunc,
	updateValidation rest.ValidateObjectUpdateFunc,
	forceAllowCreate bool,
	options *meta.UpdateOptions,
) (runtime.Object, bool, error) {
	obj, updated, err := s.ScaleSubResourceParentStore.Update(
		ctx,
		name,
		toScaleUpdatedObjectInfo(objInfo),
		toScaleCreateValidation(createValidation),
		updateValidation,
		forceAllowCreate,
		options,
	)
	if err != nil {
		return nil, false, err
	}

	withScale, ok := obj.(ObjectWithScaleSubResource)
	if !ok {
		return nil, false, kerrors.NewBadRequest("do not support scale subresource")
	}
	return withScale.GetScale(), updated, nil
}
