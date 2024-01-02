package connector

import (
	"context"
	"errors"
	"fmt"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/connector"
	"github.com/seal-io/walrus/pkg/dao/model/predicate"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/operator"
	optypes "github.com/seal-io/walrus/pkg/operator/types"
	"github.com/seal-io/walrus/pkg/vcs"
	"github.com/seal-io/walrus/utils/validation"
)

type (
	CreateRequest struct {
		model.ConnectorCreateInput `path:",inline" json:",inline"`
	}

	CreateResponse = *model.ConnectorOutput
)

func (r *CreateRequest) Validate() error {
	if err := r.ConnectorCreateInput.Validate(); err != nil {
		return err
	}

	if err := validation.IsDNSLabel(r.Name); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	if r.Type == "" {
		return errors.New("invalid type: blank")
	}

	if r.ConfigData == nil {
		return errors.New("invalid config data: empty")
	}

	if !types.IsEnvironmentType(r.ApplicableEnvironmentType) {
		return fmt.Errorf("invalid applicable environment type: %s", r.ApplicableEnvironmentType)
	}

	if err := validateConnector(r.Context, r.Client, r.Model()); err != nil {
		return err
	}

	return nil
}

type (
	GetRequest struct {
		model.ConnectorQueryInput `path:",inline"`
	}

	GetResponse = *model.ConnectorOutput
)

type UpdateRequest struct {
	model.ConnectorUpdateInput `path:",inline" json:",inline"`
}

func (r *UpdateRequest) Validate() error {
	if err := r.ConnectorUpdateInput.Validate(); err != nil {
		return err
	}

	if err := validation.IsDNSLabel(r.Name); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	if r.Type == "" {
		return errors.New("invalid type: blank")
	}

	if r.ConfigData != nil {
		if err := validateConnector(r.Context, r.Client, r.Model()); err != nil {
			return err
		}
	}

	return nil
}

type DeleteRequest = model.ConnectorDeleteInput

type (
	CollectionGetRequest struct {
		model.ConnectorQueryInputs `path:",inline" query:",inline"`

		runtime.RequestCollection[
			predicate.Connector, connector.OrderOption,
		] `query:",inline"`

		WithGlobal bool `query:"withGlobal,omitempty"`

		Stream *runtime.RequestUnidiStream
	}

	CollectionGetResponse = []*model.ConnectorOutput
)

func (r *CollectionGetRequest) SetStream(stream runtime.RequestUnidiStream) {
	r.Stream = &stream
}

type CollectionDeleteRequest = model.ConnectorDeleteInputs

func validateConnector(ctx context.Context, mc model.ClientSet, entity *model.Connector) error {
	ops := optypes.CreateOptions{
		Connector:   *entity,
		ModelClient: mc,
	}

	switch entity.Category {
	case types.ConnectorCategoryKubernetes, types.ConnectorCategoryCloudProvider:
		op, err := operator.Get(ctx, ops)
		if err != nil {
			return fmt.Errorf("invalid connector config: %w", err)
		}

		if err = op.IsConnected(ctx); err != nil {
			return fmt.Errorf("unreachable connector: %w", err)
		}
	case types.ConnectorCategoryVersionControl:
		vcsClient, err := vcs.NewClient(entity)
		if err != nil {
			return fmt.Errorf("invalid connector config: %w", err)
		}

		_, _, err = vcsClient.Users.Find(ctx)
		if err != nil {
			return fmt.Errorf("invalid connector: %w", err)
		}
	case types.ConnectorCategoryCustom:

	default:
		return errors.New("invalid connector category")
	}

	if entity.Category != types.ConnectorCategoryKubernetes && entity.EnableFinOps {
		return errors.New("invalid connector: finOps not supported")
	}

	return nil
}
