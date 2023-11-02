package property

import (
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/seal-io/walrus/utils/errorx"
	"github.com/seal-io/walrus/utils/json"
)

type (
	// Value indicates the value of property.
	Value = json.RawMessage
)

// Values holds the Value collection in map,
// the key of map is the name of Property,
// stores into json.
type Values map[string]Value

// Value implements driver.Valuer.
func (i Values) Value() (driver.Value, error) {
	return json.Marshal(i)
}

// Scan implements sql.Scanner.
func (i *Values) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		return nil
	case []byte:
		return json.Unmarshal(v, i)
	}

	return errors.New("not a valid property values")
}

// ValidateWith validates the property value with the given schemas.
func (i Values) ValidateWith(schema *openapi3.Schema) error {
	if schema == nil {
		return nil
	}

	// Check required and undefined.
	l := sets.StringKeySet(i)
	r := sets.NewString(schema.Required...)
	a := sets.StringKeySet(schema.Properties)

	if diff := r.Difference(l).UnsortedList(); len(diff) != 0 {
		return fmt.Errorf("not found required values %v", diff)
	}

	if diff := l.Difference(a).UnsortedList(); len(diff) != 0 {
		return fmt.Errorf("found undefiend values %v", diff)
	}

	// Validate.
	for n, v := range i {
		if schema.Properties[n] == nil || schema.Properties[n].Value == nil {
			continue
		}

		var (
			s           = schema.Properties[n].Value
			errTypeFunc = func(name string, ok bool, err error, expectedType string, actualValue any) error {
				if !ok {
					return fmt.Errorf("%s is not type %s, actual value: %v", name, expectedType, actualValue)
				}

				if err != nil {
					return fmt.Errorf("failed to convert %s to %s: %w", name, expectedType, err)
				}
				return nil
			}
			validateSchemaFunc = func(name string, val any) error {
				err := s.VisitJSON(val)
				if err != nil {
					var e *openapi3.SchemaError
					if errors.As(err, &e) {
						return errorx.Errorf("invalid %s: %v", name, e.Reason)
					}
					return err
				}
				return nil
			}
		)

		switch s.Type {
		case openapi3.TypeString:
			val, ok, err := GetString(v)
			if !ok || err != nil {
				return errTypeFunc(n, ok, err, s.Type, v)
			}

			return validateSchemaFunc(n, val)
		case openapi3.TypeBoolean:
			val, ok, err := GetBool(v)
			if !ok || err != nil {
				return errTypeFunc(n, ok, err, s.Type, v)
			}

			return validateSchemaFunc(n, val)
		case openapi3.TypeInteger:
			val, ok, err := GetInt(v)
			if !ok || err != nil {
				return errTypeFunc(n, ok, err, s.Type, v)
			}

			return validateSchemaFunc(n, val)
		case openapi3.TypeNumber:
			val, ok, err := GetNumber(v)
			if !ok || err != nil {
				return errTypeFunc(n, ok, err, s.Type, v)
			}

			return validateSchemaFunc(n, val)
		case openapi3.TypeArray:
			val, ok, err := GetSlice[any](v)
			if !ok || err != nil {
				return errTypeFunc(n, ok, err, s.Type, v)
			}

			return validateSchemaFunc(n, val)
		case openapi3.TypeObject:
			val, ok, err := GetMap[any](v)
			if !ok || err != nil {
				return errTypeFunc(n, ok, err, s.Type, v)
			}

			return validateSchemaFunc(n, val)
		default:
			_, ok, err := GetAny[any](v)
			if !ok || err != nil {
				return errTypeFunc(n, ok, err, s.Type, v)
			}
		}
	}

	return nil
}
