package view

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"k8s.io/apimachinery/pkg/util/sets"

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

	moduleNames := sets.Set[string]{}
	ps := make([]predicate.ModuleVersion, len(inputModules))
	for i, m := range inputModules {
		if m.ModuleID == "" {
			return errors.New("invalid module id: blank")
		}
		if err := validation.IsDNSSubdomainName(m.Name); err != nil {
			return fmt.Errorf("invalid module name %s: %w", m.Name, err)
		}
		if moduleNames.Has(m.Name) {
			return fmt.Errorf("invalid module name %s: duplicated", m.Name)
		}
		moduleNames.Insert(m.Name)
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

	var schemas = make(map[string]map[string]types.ModuleVariable)
	for _, m := range moduleVersions {
		if m.Schema == nil {
			continue
		}
		key := moduleVersionKey(m.ModuleID, m.Version)
		if _, ok := schemas[key]; !ok {
			schemas[key] = make(map[string]types.ModuleVariable)
		}
		for _, v := range m.Schema.Variables {
			schemas[key][v.Name] = v
		}
	}

	for _, v := range inputModules {
		schemaVariables := schemas[moduleVersionKey(v.ModuleID, v.Version)]
		// schema doesn't exist
		if schemaVariables == nil {
			return fmt.Errorf("invalid module %s: empty schema", v.Name)
		}

		// check input unsupported attributes
		for attrName := range v.Attributes {
			if _, supported := schemaVariables[attrName]; !supported {
				return fmt.Errorf("found unknown attribute %s in module %s, please refresh module %s to load the latest attributes", attrName, v.Name, v.ModuleID)
			}
		}

		// check attributes
		for _, schemaVariable := range schemaVariables {
			attrValue, existed := v.Attributes[schemaVariable.Name]
			// check required attribute existed
			if schemaVariable.Required && (!existed || attrValue == nil) {
				return fmt.Errorf("required attribute %s in module %s is empty", schemaVariable.Name, v.Name)
			}

			// omit it when the value is null
			if attrValue == nil {
				continue
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
				return fmt.Errorf("unexpected value type for attribute %s in module %s", schemaVariable.Name, v.Name)
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
	IDs        []types.ID                 `json:"ids,omitempty"`
	Collection []*model.ApplicationOutput `json:"collection,omitempty"`
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
