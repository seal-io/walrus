package view

import (
	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
)

// Basic APIs

// Batch APIs

type CollectionGetRequest struct {
	runtime.RequestPagination `query:",inline"`

	Domain string `query:"domain,omitempty"`
}

type CollectionGetResponse = []*model.RoleOutput

// Extensional APIs
