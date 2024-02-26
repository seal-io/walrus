package types

import (
	"github.com/getkin/kin-openapi/openapi3"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/zclconf/go-cty/cty"

	"github.com/seal-io/walrus/pkg/dao/types/property"
	"github.com/seal-io/walrus/utils/json"
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
		p.ResourceComponentChanges[i] = &ResourceComponentChange{
			ResourceChange: rc,
			Change: &Change{
				Change: rc.Change,
			},
		}

		switch {
		case rc.Change.Actions.Create():
			p.ResourceComponentChanges[i].Change.Type = ResourceComponentChangeTypeCreate
		case rc.Change.Actions.Update():
			p.ResourceComponentChanges[i].Change.Type = ResourceComponentChangeTypeUpdate
		case rc.Change.Actions.Delete():
			p.ResourceComponentChanges[i].Change.Type = ResourceComponentChangeTypeDelete
		}
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
	ResourceComponentChangeTypeCreate = "create"
	ResourceComponentChangeTypeUpdate = "update"
	ResourceComponentChangeTypeDelete = "delete"
)

type Change struct {
	*tfjson.Change `json:",inline"`

	Type string `json:"type"`
}
