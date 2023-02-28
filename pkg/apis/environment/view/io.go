package view

import (
	"errors"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/utils/json"
)

// Basic APIs

type EnvironmentVO struct {
	*model.Environment `json:",inline"`

	ConnectorIDs []types.ID `json:"connectorIDs,omitempty"`
}

func (v EnvironmentVO) MarshalJSON() ([]byte, error) {
	type Alias EnvironmentVO

	// mutate `.Edges.Connectors` to `.ConnectorIDs`.
	if len(v.Edges.Connectors) != 0 {
		for _, c := range v.Edges.Connectors {
			if c == nil {
				continue
			}
			v.ConnectorIDs = append(v.ConnectorIDs, c.ConnectorID)
		}
		v.Edges.Connectors = nil // release
	}

	return json.Marshal(&struct {
		*Alias `json:",inline"`
	}{
		Alias: (*Alias)(&v),
	})
}

func (v EnvironmentVO) Model() *model.Environment {
	// mutate `.ConnectorIDs` to `.Edges.Connectors`.
	for _, id := range v.ConnectorIDs {
		v.Environment.Edges.Connectors = append(v.Environment.Edges.Connectors, &model.EnvironmentConnectorRelationship{
			ConnectorID: id,
		})
	}

	if v.ConnectorIDs != nil && len(v.ConnectorIDs) == 0 {
		v.Environment.Edges.Connectors = make([]*model.EnvironmentConnectorRelationship, 0)
	}

	return v.Environment
}

type EnvironmentCreateRequest struct {
	*EnvironmentVO `json:",inline"`
}

func (r *EnvironmentCreateRequest) Validate() error {
	if r.Name == "" {
		return errors.New("invalid name: blank")
	}

	return nil
}

type EnvironmentCreateResponse = EnvironmentVO

type EnvironmentUpdateRequest struct {
	UriID types.ID `uri:"id"`

	*EnvironmentVO `json:",inline"`
}

func (r *EnvironmentUpdateRequest) Validate() error {
	if !r.UriID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	r.ID = r.UriID
	return nil
}

type EnvironmentGetRequest struct {
	ID types.ID `uri:"id"`
}

func (r *EnvironmentGetRequest) Validate() error {
	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}
	return nil
}

type EnvironmentGetResponse = EnvironmentVO

type EnvironmentDeleteRequest = EnvironmentGetRequest

// Batch APIs

type CollectionGetRequest struct {
	runtime.RequestPagination `query:",inline"`
	runtime.RequestExtracting `query:",inline"`
	runtime.RequestSorting    `query:",inline"`
}

type CollectionGetResponse = []*EnvironmentVO

// Extensional APIs
