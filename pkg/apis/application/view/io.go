package view

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/application"
	"github.com/seal-io/seal/pkg/dao/model/moduleversion"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/topic/datamessage"
	"github.com/seal-io/seal/utils/validation"
)

// Basic APIs

type CreateRequest struct {
	*model.ApplicationCreateInput `json:",inline"`
}

func (r *CreateRequest) ValidateWith(ctx context.Context, input any) error {
	var modelClient = input.(model.ClientSet)

	if r.Project.ID == "" {
		return errors.New("invalid project id: blank")
	}
	if err := validation.IsDNSSubdomainName(r.Name); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}
	if len(r.Modules) != 0 {
		return validateModules(ctx, modelClient, r.Model().Edges.Modules)
	}
	return nil
}

func validateModules(ctx context.Context, modelClient model.ClientSet, inputModules []*model.ApplicationModuleRelationship) error {
	moduleVersionKey := func(moduleID, version string) string {
		return fmt.Sprintf("%s/%s", moduleID, version)
	}

	ps := make([]predicate.ModuleVersion, len(inputModules))
	for i, m := range inputModules {
		if inputModules[i].ModuleID == "" {
			return errors.New("invalid module id: blank")
		}
		if err := validation.IsDNSSubdomainName(inputModules[i].Name); err != nil {
			return fmt.Errorf("invalid module name: %w", err)
		}
		ps[i] = moduleversion.And(moduleversion.ModuleID(m.ModuleID), moduleversion.Version(m.Version))
	}

	moduleVersions, err := modelClient.ModuleVersions().
		Query().
		Select(
			moduleversion.FieldModuleID,
			moduleversion.FieldVersion,
			moduleversion.FieldSchema,
		).
		Where(moduleversion.Or(ps...)).
		All(ctx)
	if err != nil {
		return err
	}

	var attrs = make(map[string]map[string]types.ModuleVariable)
	for _, m := range moduleVersions {
		if m.Schema == nil {
			continue
		}
		key := moduleVersionKey(m.ModuleID, m.Version)
		if _, ok := attrs[key]; !ok {
			attrs[key] = make(map[string]types.ModuleVariable)
		}
		for _, v := range m.Schema.Variables {
			attrs[key][v.Name] = v
		}
	}

	for _, v := range inputModules {
		key := moduleVersionKey(v.ModuleID, v.Version)
		for attrName, attrValue := range v.Attributes {
			// check attribute existed.
			schemaVariable, ok := attrs[key][attrName]
			if !ok {
				return fmt.Errorf("invalid attribute %s in module %s: not supported", attrName, key)
			}

			// check attribute type,
			// only check primitive types, skip complex types now, since we currently use string to represent the these values.
			var (
				valueTypeExpected = true
				actualValueType   = reflect.TypeOf(attrValue).Kind()
			)
			switch schemaVariable.Type {
			case "string":
				valueTypeExpected = actualValueType == reflect.String
			case "bool":
				valueTypeExpected = actualValueType == reflect.Bool
			case "number":
				switch actualValueType {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				case reflect.Float32, reflect.Float64:
				default:
					valueTypeExpected = false
				}
			}

			if !valueTypeExpected {
				return fmt.Errorf("unexpected value type for attribute %s in module %s", attrName, key)
			}
		}
	}
	return nil
}

type CreateResponse = *model.ApplicationOutput

type DeleteRequest = GetRequest

type UpdateRequest struct {
	*model.ApplicationUpdateInput `uri:",inline" json:",inline"`
}

func (r *UpdateRequest) ValidateWith(ctx context.Context, input any) error {
	var modelClient = input.(model.ClientSet)

	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}
	if err := validation.IsDNSSubdomainName(r.Name); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}
	if len(r.Modules) != 0 {
		return validateModules(ctx, modelClient, r.Model().Edges.Modules)
	}
	return nil
}

type GetRequest struct {
	*model.ApplicationQueryInput `uri:",inline"`
}

func (r *GetRequest) Validate() error {
	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}

	return nil
}

type GetResponse = *model.ApplicationOutput

type StreamResponse struct {
	Type       datamessage.EventType      `json:"type"`
	IDs        []types.ID                 `json:"ids"`
	Collection []*model.ApplicationOutput `json:"collection"`
}

type StreamRequest struct {
	ID types.ID `uri:"id"`
}

func (r *StreamRequest) ValidateWith(ctx context.Context, input any) error {
	if !r.ID.Valid(0) {
		return errors.New("invalid id: blank")
	}
	var client = input.(model.ClientSet)
	exist, err := client.Applications().Query().
		Where(application.ID(r.ID)).
		Exist(ctx)
	if err != nil || !exist {
		return runtime.Errorw(err, "invalid id: not found")
	}

	return nil
}

// Batch APIs

type CollectionDeleteRequest []*model.ApplicationQueryInput

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
	runtime.RequestCollection[predicate.Application] `query:",inline"`

	ProjectIDs []types.ID `query:"projectID"`
}

func (r *CollectionGetRequest) Validate() error {
	if len(r.ProjectIDs) == 0 {
		return errors.New("invalid input: missing project id")
	}

	for i := range r.ProjectIDs {
		if !r.ProjectIDs[i].Valid(0) {
			return errors.New("invalid project id: blank")
		}
	}
	return nil
}

type CollectionGetResponse = []*model.ApplicationOutput

type CollectionStreamRequest struct {
	runtime.RequestExtracting `query:",inline"`

	ProjectIDs []types.ID `query:"projectID,omitempty"`
}

func (r *CollectionStreamRequest) Validate() error {
	for i := range r.ProjectIDs {
		if !r.ProjectIDs[i].Valid(0) {
			return errors.New("invalid project id: blank")
		}
	}
	return nil
}

// Extensional APIs
