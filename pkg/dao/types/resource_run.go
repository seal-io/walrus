package types

import (
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/mohae/deepcopy"
	"github.com/zclconf/go-cty/cty"

	"github.com/seal-io/walrus/pkg/dao/types/property"
	"github.com/seal-io/walrus/utils/json"
	"github.com/seal-io/walrus/utils/log"
)

type OutputValue struct {
	Name      string          `json:"name,omitempty"`
	Value     property.Value  `json:"value,omitempty,cli-table-column"`
	Type      cty.Type        `json:"type,omitempty"`
	Sensitive bool            `json:"sensitive,omitempty"`
	Schema    openapi3.Schema `json:"schema,omitempty"`
}

// Task type defines the type of the task to be performed with deployer.
const (
	RunTaskTypeApply   RunJobType = "apply"
	RunTaskTypePlan    RunJobType = "plan"
	RunTaskTypeDestroy RunJobType = "destroy"
)

const (
	RunTypeCreate   RunType = "create"
	RunTypeUpdate   RunType = "update"
	RunTypeDelete   RunType = "delete"
	RunTypeStart    RunType = "start"
	RunTypeStop     RunType = "stop"
	RunTypeRollback RunType = "rollback"
)

type RunType string

func (t RunType) String() string {
	return string(t)
}

type RunJobType string

func (t RunJobType) String() string {
	return string(t)
}

type ResourceRunConfigData = []byte

type Plan struct {
	tfjson.Plan `json:",inline"`

	// ResourceComponentChanges is the changes of the resource components.
	ResourceComponentChanges []*ResourceComponentChange `json:"resource_changes,omitempty"`
}

// UnmarshalJSON customizes the JSON decoding of the Plan.
func (p *Plan) UnmarshalJSON(data []byte) error {
	// Unmarshal data into the embedded tfjson.Plan.
	if err := json.Unmarshal(data, &p.Plan); err != nil {
		return err
	}

	p.ResourceComponentChanges = make([]*ResourceComponentChange, len(p.Plan.ResourceChanges))

	for i, rc := range p.Plan.ResourceChanges {
		change := (&Change{Change: rc.Change}).Process()

		rcc := &ResourceComponentChange{
			ResourceChange: rc,
			Change:         change,
		}

		rcc.ResourceChange.Change = nil
		p.ResourceComponentChanges[i] = rcc
	}

	return nil
}

func (p *Plan) GetResourceChangeSummary() ResourceComponentChangeSummary {
	summary := ResourceComponentChangeSummary{}

	for _, change := range p.ResourceComponentChanges {
		switch change.Change.Type {
		case ResourceComponentChangeTypeCreate:
			summary.Created++
		case ResourceComponentChangeTypeUpdate:
			summary.Updated++
		case ResourceComponentChangeTypeDelete:
			summary.Deleted++
		}
	}

	return summary
}

// ResourceComponentChangeSummary is the summary of the resource component changes.
type ResourceComponentChangeSummary struct {
	Created int `json:"created"`
	Updated int `json:"updated"`
	Deleted int `json:"deleted"`
}

// ResourceComponentChange is the change of the resource component.
type ResourceComponentChange struct {
	*tfjson.ResourceChange `json:",inline"`

	Change *Change `json:"change"`
}

const (
	ResourceComponentChangeTypeCreate   = "create"
	ResourceComponentChangeTypeUpdate   = "update"
	ResourceComponentChangeTypeDelete   = "delete"
	ResourceComponentChangeTypeNoChange = "no-change"
)

type Change struct {
	*tfjson.Change `json:",inline"`

	Type string `json:"type"`
}

// Process parses the change type from the actions and sets the type to the change.
func (c *Change) Process() *Change {
	logger := log.WithName("component-change")

	switch {
	case c.Actions.Create():
		c.Type = ResourceComponentChangeTypeCreate
	case c.Actions.Update(),
		c.Actions.Replace():
		c.Type = ResourceComponentChangeTypeUpdate
	case c.Actions.Delete():
		c.Type = ResourceComponentChangeTypeDelete
	case c.Actions.NoOp(), c.Actions.Read():
		c.Type = ResourceComponentChangeTypeNoChange
	}

	var (
		diff any
		// SensitivePatchedDiff is the diff with the sensitive values patched.
		sensitivePatchedDiff any
		err                  error
	)
	if c.Change.Before != nil && c.Change.After != nil {
		// Ignore error as the diff is not critical.
		diff, err = json.CreateMergePatch(c.Change.Before, c.Change.After)
		if err != nil {
			logger.Warnf("failed to create the merge diff: %v", err)
		}

		afterSensitive := deepcopy.Copy(c.Change.AfterSensitive)

		// Patch the sensitive values changes.
		sensitivePatchedDiff = patchLeaf(diff, afterSensitive, "<sensitive value(changed)>", false)
	}

	c.Change.Before = patchLeaf(c.Change.Before, c.Change.BeforeSensitive, "<sensitive value>", false)
	c.Change.After = patchLeaf(c.Change.After, c.Change.AfterSensitive, "<sensitive value>", false)

	if sensitivePatchedDiff != nil && c.Change.After != nil {
		object, err := json.PatchObject(c.Change.After, sensitivePatchedDiff)
		if err == nil {
			c.Change.After = object
		} else {
			logger.Warnf("failed to patch the sensitive values: %v", err)
		}
	}

	c.Change.After = patchLeaf(c.Change.After, c.Change.AfterUnknown, "<known after apply>", true)

	c.Change.BeforeSensitive = nil
	c.Change.AfterSensitive = nil
	c.Change.AfterUnknown = nil
	c.Change.GeneratedConfig = ""
	c.Change.ReplacePaths = nil

	return c
}

// patchLeaf patch the raw value with the masked leaf value with the mask.
func patchLeaf(value, toMaskLeaf any, mask string, merge bool) any {
	logger := log.WithName("component-change")
	if value == nil || toMaskLeaf == nil {
		return value
	}

	maskedLeaf := maskLeafValues(value, toMaskLeaf, mask, merge)
	if maskedLeaf == nil {
		return value
	}

	patched, err := json.PatchObject(value, maskedLeaf)
	if err != nil {
		logger.Warnf("failed to patch the leaf value: %v", err)

		return value
	}

	if ptr := reflect.ValueOf(patched); ptr.Kind() == reflect.Ptr {
		patched = ptr.Elem().Interface()
	}

	return patched
}

// maskLeafValues masks the leaf values of the raw value with the mask.
// The toMaskLeafs record the leaf key to be masked.
// If merge is true, the leaf values will be merged with raw value.
func maskLeafValues(rawValue, toMaskLeafs any, mask string, merge bool) any {
	if isEmptyValueLeaf(toMaskLeafs) {
		return nil
	}

	// If the mask value is true, replace the raw value with the mask.
	if boolVal, isBool := toMaskLeafs.(bool); isBool {
		if boolVal {
			return mask
		}
		return rawValue
	}

	if ptr := reflect.ValueOf(rawValue); ptr.Kind() == reflect.Ptr {
		rawValue = ptr.Elem().Interface()
	}

	if ptr := reflect.ValueOf(toMaskLeafs); ptr.Kind() == reflect.Ptr {
		toMaskLeafs = ptr.Elem().Interface()
	}

	switch leafVal := toMaskLeafs.(type) {
	case map[string]any:
		val, ok := rawValue.(map[string]any)
		if !ok {
			return rawValue
		}

		for k := range leafVal {
			if merge && leafVal[k] == true {
				leafVal[k] = mask
			} else {
				if _, ok := val[k]; !ok {
					if merge {
						continue
					}

					delete(leafVal, k)
					continue
				}

				leafVal[k] = maskLeafValues(val[k], leafVal[k], mask, merge)

				if leafVal[k] == nil {
					delete(leafVal, k)
				}
			}
		}

		if len(leafVal) == 0 {
			return nil
		}

		return leafVal
	case []any:
		val, ok := rawValue.([]any)
		if !ok {
			return rawValue
		}

		maskLen := len(val)
		if merge && len(leafVal) > maskLen {
			maskLen = len(leafVal)
		}

		masked := make([]any, maskLen)

		for i := range masked {
			if merge && i >= len(val) && !isEmptyValueLeaf(leafVal[i]) {
				if leafVal[i] == true {
					masked[i] = mask
					continue
				} else {
					masked[i] = leafVal[i]
					continue
				}
			}

			masked[i] = val[i]
			if i >= len(leafVal) {
				continue
			}

			m := maskLeafValues(val[i], leafVal[i], mask, merge)
			if m != nil {
				masked[i] = m
			}
		}

		return masked
	case bool:
		if leafVal {
			return mask
		}
	}

	return rawValue
}

func isEmptyValueLeaf(v any) bool {
	reflectValue := reflect.ValueOf(v)
	switch reflectValue.Kind() {
	case reflect.String, reflect.Array:
		return reflectValue.Len() == 0
	case reflect.Map, reflect.Slice:
		return reflectValue.Len() == 0 || reflectValue.IsNil()
	case reflect.Bool:
		// As the value is a boolean, it is not an empty value.
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return reflectValue.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return reflectValue.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return reflectValue.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return reflectValue.IsNil()
	default:
	}

	return false
}
