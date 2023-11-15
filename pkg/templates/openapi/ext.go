package openapi

import (
	"github.com/getkin/kin-openapi/openapi3"
)

const OpenAPIVersion = "3.0.3"

const (
	// Extension for UI.
	ExtUI          = "x-walrus-ui"
	ExtUIGroup     = "group"
	ExtUIShowIf    = "show-if"
	ExtUIHidden    = "hidden"
	ExtUIImmutable = "immutable"
	ExtUIWidget    = "widget"
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
		   ```
	*/
)

const (
	// Extension for original value.
	ExtOriginal                  = "x-walrus-original"
	ExtOriginalType              = "type"
	ExtOriginalValueExpression   = "value-expression"
	ExtOriginalVariablesSequence = "sequence"
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

const (
	// Extension for walrus.
	ExtWalrus        = "x-walrus"
	ExtWalrusVersion = "version"

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

type Ext map[string]map[string]any

func NewExt(c map[string]any) Ext {
	if c == nil {
		return Ext{}
	}

	e := Ext{}
	for k, v := range c {
		e[k] = v.(map[string]any)
	}

	return e
}

func (e Ext) SetOriginalType(t any) Ext {
	if e[ExtOriginal] == nil {
		e[ExtOriginal] = map[string]any{}
	}

	e[ExtOriginal][ExtOriginalType] = t

	return e
}

func (e Ext) SetOriginalValueExpression(ve []byte) Ext {
	if e[ExtOriginal] == nil {
		e[ExtOriginal] = map[string]any{}
	}

	e[ExtOriginal][ExtOriginalValueExpression] = string(ve)

	return e
}

func (e Ext) SetOriginalVariablesSequence(s []string) Ext {
	if e[ExtOriginal] == nil {
		e[ExtOriginal] = map[string]any{}
	}

	e[ExtOriginal][ExtOriginalVariablesSequence] = s

	return e
}

func (e Ext) SetUIGroup(gp string) Ext {
	if e[ExtUI] == nil {
		e[ExtUI] = map[string]any{}
	}

	e[ExtUI][ExtUIGroup] = gp

	return e
}

func (e Ext) SetUIWidget(w string) Ext {
	if e[ExtUI] == nil {
		e[ExtUI] = map[string]any{}
	}

	e[ExtUI][ExtUIWidget] = w

	return e
}

func (e Ext) SetUIHidden() Ext {
	if e[ExtUI] == nil {
		e[ExtUI] = map[string]any{}
	}

	e[ExtUI][ExtUIHidden] = true

	return e
}

func (e Ext) SetUIImmutable() Ext {
	if e[ExtUI] == nil {
		e[ExtUI] = map[string]any{}
	}

	e[ExtUI][ExtUIImmutable] = true

	return e
}

func (e Ext) SetUIShowIf(showIf string) Ext {
	if e[ExtUI] == nil {
		e[ExtUI] = map[string]any{}
	}

	e[ExtUI][ExtUIShowIf] = showIf

	return e
}

func (e Ext) SetWalrusVersion(s string) Ext {
	if e[ExtWalrus] == nil {
		e[ExtWalrus] = map[string]any{}
	}

	e[ExtWalrus][ExtWalrusVersion] = s

	return e
}

func (e Ext) Export() map[string]any {
	if len(e) == 0 {
		return nil
	}

	result := make(map[string]any)
	for k := range e {
		result[k] = e[k]
	}

	return result
}

func GetOriginalType(e map[string]any) any {
	if e[ExtOriginal] == nil {
		return nil
	}

	eo, ok := e[ExtOriginal].(map[string]any)
	if !ok {
		return nil
	}

	val, ok := eo[ExtOriginalType]
	if !ok {
		return nil
	}

	return val
}

func GetOriginalValueExpression(e map[string]any) []byte {
	if e[ExtOriginal] == nil {
		return nil
	}

	eo, ok := e[ExtOriginal].(map[string]any)
	if !ok {
		return nil
	}

	val, ok := eo[ExtOriginalValueExpression]
	if !ok {
		return nil
	}

	vb, _ := val.([]byte)

	return vb
}

func GetOriginalVariablesSequence(e map[string]any) []string {
	if e[ExtOriginal] == nil {
		return nil
	}

	eo, ok := e[ExtOriginal].(map[string]any)
	if !ok {
		return nil
	}

	val, ok := eo[ExtOriginalVariablesSequence]
	if !ok {
		return nil
	}

	switch v := val.(type) {
	case []string:
		return v
	case []any:
		var vb []string

		for _, vq := range v {
			if s, ok := vq.(string); ok {
				vb = append(vb, s)
			}
		}

		return vb
	}

	return nil
}

func GetUIGroup(e map[string]any) string {
	if e[ExtUI] == nil {
		return ""
	}

	eo, ok := e[ExtUI].(map[string]any)
	if !ok {
		return ""
	}

	val, ok := eo[ExtUIGroup]
	if !ok {
		return ""
	}

	vb, _ := val.(string)

	return vb
}

func GetUIShowIf(e map[string]any) string {
	if e[ExtUI] == nil {
		return ""
	}

	eo, ok := e[ExtUI].(map[string]any)
	if !ok {
		return ""
	}

	val, ok := eo[ExtUIShowIf]
	if !ok {
		return ""
	}

	vb, _ := val.(string)

	return vb
}

func GetUIHidden(e map[string]any) bool {
	if e[ExtUI] == nil {
		return false
	}

	eo, ok := e[ExtUI].(map[string]any)
	if !ok {
		return false
	}

	val, ok := eo[ExtUIHidden]
	if !ok {
		return false
	}

	vb, _ := val.(bool)

	return vb
}

func GetUIImmutable(e map[string]any) bool {
	if e[ExtUI] == nil {
		return false
	}

	eo, ok := e[ExtUI].(map[string]any)
	if !ok {
		return false
	}

	val, ok := eo[ExtUIImmutable]
	if !ok {
		return false
	}

	vb, _ := val.(bool)

	return vb
}

func GetUIWidget(e map[string]any) string {
	if e[ExtUI] == nil {
		return ""
	}

	eo, ok := e[ExtUI].(map[string]any)
	if !ok {
		return ""
	}

	val, ok := eo[ExtUIWidget]
	if !ok {
		return ""
	}

	vb, _ := val.(string)

	return vb
}

func GetWalrusVersion(e map[string]any) string {
	if e[ExtWalrus] == nil {
		return ""
	}

	eo, ok := e[ExtWalrus].(map[string]any)
	if !ok {
		return ""
	}

	val, ok := eo[ExtWalrusVersion]
	if !ok {
		return ""
	}

	vb, _ := val.(string)

	return vb
}

func IsSchemaRefEmpty(s *openapi3.SchemaRef) bool {
	return s == nil || s.Value == nil || s.Value.IsEmpty()
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
	return RemoveExt(ExtOriginal, s)
}

func RemoveExtUI(s *openapi3.Schema) *openapi3.Schema {
	return RemoveExt(ExtUI, s)
}
