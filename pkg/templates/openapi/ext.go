package openapi

import (
	"github.com/getkin/kin-openapi/openapi3"

	"github.com/seal-io/walrus/utils/json"
)

const OpenAPIVersion = "3.0.3"

const (
	// Extension for walrus.
	ExtWalrusKey = "x-walrus"
	// ExtWalrusVersionKey is a string, for walrus version constraint.
	ExtWalrusVersionKey = "version"
	/* Version Constraint is under info extension.
	   Example:
	   ```yaml
	   openapi: 3.0.3
	   info:
	     title: OpenAPI schema for template webservice
	     x-walrus:
	   	version: '>=0.4.0-rc1'
	    ```
	*/
)

const (
	// Extension for UI.
	ExtUIKey = "x-walrus-ui"
	/* UI is under schema extension.
		   Example:
		   ```yaml
		   components:
		     schemas:
		       variables:
		         type: object
		         properties:
	                   image:
		             title: Image Name
		   	     type: string
		             description: Docker image name
		             x-walrus-ui:
		   	       group: Basic
		               showIf:
		               hidden:
		               immutable:
		   	       widget:
	                       order:
	                       colSpan:
		   ```
	*/
)

// ExtUI is a struct wrap the UI extension.
type ExtUI struct {
	// Group is a string, for grouping the properties.
	Group string `json:"group,omitempty" yaml:"group,omitempty"`
	// GroupOrder is a list, for ordering the group in the UI.
	GroupOrder []string `json:"groupOrder,omitempty" yaml:"groupOrder,omitempty"`
	// ShowIf is a string, for showing the property.
	ShowIf string `json:"showIf,omitempty" yaml:"showIf,omitempty"`
	// Hidden is a boolean, for hiding the property.
	Hidden bool `json:"hidden,omitempty" yaml:"hidden,omitempty"`
	// Immutable is a boolean, for making the property immutable.
	Immutable bool `json:"immutable,omitempty" yaml:"immutable,omitempty"`
	// Widget is a string, for customizing the UI widget.
	Widget string `json:"widget,omitempty" yaml:"widget,omitempty"`
	// Order is a number, for ordering the properties in the UI.
	Order int `json:"order,omitempty" yaml:"order,omitempty"`
	// ColSpan is a number between 1 and 12, for typical 12-column grid systems.
	ColSpan int `json:"colSpan,omitempty" yaml:"colSpan,omitempty"`
}

// IsEmpty reports if the extension is empty.
func (e ExtUI) IsEmpty() bool {
	return e.Group == "" &&
		e.ShowIf == "" &&
		!e.Hidden &&
		!e.Immutable &&
		e.Widget == "" &&
		e.Order <= 0 &&
		e.ColSpan <= 0 &&
		len(e.GroupOrder) == 0
}

const (
	// ExtOriginalKey for original value.
	ExtOriginalKey = "x-walrus-original"
	/* Origin is under schema extension.
	   Example:
	   ```yaml
	   components:
	     schemas:
	       variables:
	         type: object
	         properties:
	           image:
	             title: Image Name
	             type: string
	             description: Docker image name
	             x-walrus-original:
	               type: list
	   ```
	*/
)

// ExtOriginal is a struct wrap the original extension.
type ExtOriginal struct {
	// Type is a string, for original type.
	Type any `json:"type,omitempty" yaml:"type,omitempty"`
	// ValueExpression is a string, for original value expression.
	ValueExpression []byte `json:"value-expression,omitempty" yaml:"value-expression,omitempty"`
	// VariablesSequence is a list, for original variables sequence.
	VariablesSequence []string `json:"sequence,omitempty" yaml:"sequence,omitempty"`
}

// IsEmpty reports if the extension is empty.
func (e ExtOriginal) IsEmpty() bool {
	return e.Type == nil && len(e.ValueExpression) == 0 && len(e.VariablesSequence) == 0
}

// Ext is a struct wrap the extension.
type Ext struct {
	ExtUI
	ExtOriginal
}

// NewExt creates a new Ext.
func NewExt() *Ext {
	return &Ext{}
}

// NewExtFromMap creates a new Ext from extension map.
func NewExtFromMap(m map[string]any) *Ext {
	e := NewExt()
	if len(m) == 0 {
		return e
	}

	ui, ok := m[ExtUIKey]
	if ok {
		e.ExtUI, ok = ui.(ExtUI)
		if !ok {
			b, err := json.Marshal(m[ExtUIKey])
			if err == nil {
				_ = json.Unmarshal(b, &e.ExtUI)
			}
		}
	}

	origin, ok := m[ExtOriginalKey]
	if ok {
		e.ExtOriginal, ok = origin.(ExtOriginal)
		if !ok {
			b, err := json.Marshal(m[ExtOriginalKey])
			if err == nil {
				_ = json.Unmarshal(b, &e.ExtOriginal)
			}
		}
	}

	return e
}

func (e *Ext) WithOriginal(origin any) *Ext {
	o, ok := origin.(ExtOriginal)
	if ok {
		e.ExtOriginal = o
	}

	return e
}

func (e *Ext) WithOriginalType(ty any) *Ext {
	e.Type = ty
	return e
}

func (e *Ext) WithOriginalValueExpression(ve []byte) *Ext {
	e.ValueExpression = ve
	return e
}

func (e *Ext) WithOriginalVariablesSequence(vq []string) *Ext {
	e.VariablesSequence = vq
	return e
}

func (e *Ext) WithUIGroup(gp string) *Ext {
	e.Group = gp
	return e
}

func (e *Ext) WithUIGroupOrder(grd ...string) *Ext {
	if len(grd) == 0 {
		return e
	}
	e.GroupOrder = grd

	return e
}

func (e *Ext) WithUIShowIf(showIf string) *Ext {
	e.ShowIf = showIf
	return e
}

func (e *Ext) WithUIHidden() *Ext {
	e.Hidden = true
	return e
}

func (e *Ext) WithUIImmutable() *Ext {
	e.Immutable = true
	return e
}

func (e *Ext) WithUIWidget(widget string) *Ext {
	e.Widget = widget
	return e
}

func (e *Ext) WithUIOrder(order int) *Ext {
	if order > 0 {
		e.Order = order
	}

	return e
}

func (e *Ext) WithUIColSpan(cs int) *Ext {
	if cs > 0 {
		e.ColSpan = cs
	}

	return e
}

func (e *Ext) Export() map[string]any {
	if e.ExtUI.IsEmpty() && e.ExtOriginal.IsEmpty() {
		return nil
	}

	result := make(map[string]any)
	if !e.ExtUI.IsEmpty() {
		result[ExtUIKey] = e.ExtUI
	}

	if !e.ExtOriginal.IsEmpty() {
		result[ExtOriginalKey] = e.ExtOriginal
	}

	return result
}

func GetExtOriginal(e map[string]any) ExtOriginal {
	if e[ExtOriginalKey] == nil {
		return ExtOriginal{}
	}

	eo, ok := e[ExtOriginalKey].(ExtOriginal)
	if ok {
		return eo
	}

	b, err := json.Marshal(e[ExtOriginalKey])
	if err == nil {
		_ = json.Unmarshal(b, &eo)
		return eo
	}

	return ExtOriginal{}
}

func GetExtUI(e map[string]any) ExtUI {
	if e[ExtUIKey] == nil {
		return ExtUI{}
	}

	eo, ok := e[ExtUIKey].(ExtUI)
	if ok {
		return eo
	}

	b, err := json.Marshal(e[ExtUIKey])
	if err == nil {
		_ = json.Unmarshal(b, &eo)
		return eo
	}

	return ExtUI{}
}

func GetExtWalrusVersion(e map[string]any) string {
	if e[ExtWalrusKey] == nil {
		return ""
	}

	eo, ok := e[ExtWalrusKey].(map[string]any)
	if !ok {
		return ""
	}

	val, ok := eo[ExtWalrusVersionKey]
	if !ok {
		return ""
	}

	vb, _ := val.(string)

	return vb
}

func IsSchemaRefEmpty(s *openapi3.SchemaRef) bool {
	return s == nil || s.Value == nil || s.Value.IsEmpty() || len(s.Value.Extensions) == 0
}

func RemoveExt(key string, s *openapi3.Schema) *openapi3.Schema {
	if s == nil || s.Extensions == nil {
		return nil
	}

	// Self.
	delete(s.Extensions, key)

	// Properties.
	if s.Properties != nil {
		for pk := range s.Properties {
			s.Properties[pk].Value = RemoveExt(key, s.Properties[pk].Value)
		}
	}

	// AdditionalProperties.
	if !IsSchemaRefEmpty(s.AdditionalProperties.Schema) {
		s.AdditionalProperties.Schema.Value = RemoveExt(key, s.AdditionalProperties.Schema.Value)
	}

	// Items.
	if s.Items != nil && !IsSchemaRefEmpty(s.Items) {
		// Self.
		s.Items.Value = RemoveExt(key, s.Items.Value)

		// Items Properties.
		if len(s.Items.Value.Properties) != 0 {
			for pk := range s.Items.Value.Properties {
				s.Items.Value.Properties[pk].Value = RemoveExt(key, s.Items.Value.Properties[pk].Value)
			}
		}

		// AdditionalProperties.
		if s.Items.Value.AdditionalProperties.Has != nil && *s.Items.Value.AdditionalProperties.Has &&
			!IsSchemaRefEmpty(s.Items.Value.AdditionalProperties.Schema) {
			s.Items.Value.AdditionalProperties.Schema.Value = RemoveExt(
				key,
				s.Items.Value.AdditionalProperties.Schema.Value,
			)
		}
	}

	return s
}

func RemoveExtOriginal(s *openapi3.Schema) *openapi3.Schema {
	return RemoveExt(ExtOriginalKey, s)
}

func RemoveExtUI(s *openapi3.Schema) *openapi3.Schema {
	return RemoveExt(ExtUIKey, s)
}
