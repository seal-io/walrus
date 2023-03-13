package view

import (
	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
)

// Basic APIs

// Batch APIs

type CollectionGetRequest struct {
	runtime.RequestCollection[predicate.Role] `query:",inline"`

	Domain string `query:"domain,omitempty"`
}

type CollectionGetResponse = []*model.RoleOutput

// Extensional APIs
