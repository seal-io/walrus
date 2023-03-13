package view

import (
	"context"
	"errors"
	"net/http"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/application"
	"github.com/seal-io/seal/pkg/dao/model/applicationinstance"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/types"
)

// Basic APIs

type DeleteRequest struct {
	*model.ApplicationInstanceQueryInput `uri:",inline"`
}

func (r *DeleteRequest) Validate() error {
	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	return nil
}

// Batch APIs

type CollectionDeleteRequest []*model.ApplicationInstanceQueryInput

func (r CollectionDeleteRequest) Validate() error {
	if len(r) == 0 {
		return errors.New("invalid input: empty")
	}
	for _, i := range r {
		if !i.ID.Valid(0) {
			return errors.New("invalid id: blank")
		}
	}
	return nil
}

type CollectionGetRequest struct {
	runtime.RequestCollection[predicate.ApplicationInstance] `query:",inline"`

	ApplicationID types.ID `query:"applicationID"`
}

func (r *CollectionGetRequest) ValidateWith(ctx context.Context, input any) error {
	var modelClient = input.(model.ClientSet)

	if !r.ApplicationID.Valid(0) {
		return errors.New("invalid application id: blank")
	}
	_, err := modelClient.Applications().Query().
		Where(application.ID(r.ApplicationID)).
		OnlyID(ctx)
	if err != nil {
		return runtime.Error(http.StatusNotFound, "invalid application id: not found")
	}
	return nil
}

type CollectionGetResponse = []*model.ApplicationInstanceOutput

// Extensional APIs

type RouteUpgradeRequest struct {
	_ struct{} `route:"PUT=/upgrade"`

	*model.ApplicationInstanceUpdateInput `uri:",inline" json:",inline"`
}

func (r *RouteUpgradeRequest) ValidateWith(ctx context.Context, input any) error {
	var modelClient = input.(model.ClientSet)

	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	_, err := modelClient.ApplicationInstances().Query().
		Where(applicationinstance.ID(r.ID)).
		OnlyID(ctx)
	if err != nil {
		return runtime.Error(http.StatusNotFound, "invalid id: not found")
	}

	return nil
}
