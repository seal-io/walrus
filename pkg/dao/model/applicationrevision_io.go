// SPDX-FileCopyrightText: 2023 Seal, Inc
// SPDX-License-Identifier: Apache-2.0

// GENERATED, DO NOT EDIT.

package model

import (
	"time"

	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/oid"
)

// ApplicationRevisionQueryInput is the input for the ApplicationRevision query.
type ApplicationRevisionQueryInput struct {
	// ID holds the value of the "id" field.
	ID oid.ID `uri:"id,omitempty" json:"id,omitempty"`
}

// Model converts the ApplicationRevisionQueryInput to ApplicationRevision.
func (in ApplicationRevisionQueryInput) Model() *ApplicationRevision {
	return &ApplicationRevision{
		ID: in.ID,
	}
}

// ApplicationRevisionCreateInput is the input for the ApplicationRevision creation.
type ApplicationRevisionCreateInput struct {
	// Status of the resource.
	Status string `json:"status,omitempty"`
	// Extra message for status, like error details.
	StatusMessage string `json:"statusMessage,omitempty"`
	// Application modules.
	Modules []types.ApplicationModule `json:"modules,omitempty"`
	// Input variables of the revision.
	InputVariables map[string]interface{} `json:"inputVariables,omitempty"`
	// Input plan of the revision.
	InputPlan string `json:"inputPlan,omitempty"`
	// Output of the revision.
	Output string `json:"output,omitempty"`
	// Type of deployer.
	DeployerType string `json:"deployerType,omitempty"`
	// Duration in seconds of the revision deploying.
	Duration int `json:"duration,omitempty"`
	// Application instance to which the revision belongs.
	Instance ApplicationInstanceQueryInput `json:"instance"`
	// Environment to which the revision deploys.
	Environment EnvironmentQueryInput `json:"environment"`
}

// Model converts the ApplicationRevisionCreateInput to ApplicationRevision.
func (in ApplicationRevisionCreateInput) Model() *ApplicationRevision {
	var entity = &ApplicationRevision{
		Status:         in.Status,
		StatusMessage:  in.StatusMessage,
		Modules:        in.Modules,
		InputVariables: in.InputVariables,
		InputPlan:      in.InputPlan,
		Output:         in.Output,
		DeployerType:   in.DeployerType,
		Duration:       in.Duration,
	}
	entity.InstanceID = in.Instance.ID
	entity.EnvironmentID = in.Environment.ID
	return entity
}

// ApplicationRevisionUpdateInput is the input for the ApplicationRevision modification.
type ApplicationRevisionUpdateInput struct {
	// ID holds the value of the "id" field.
	ID oid.ID `uri:"id" json:"-"`
	// Status of the resource.
	Status string `json:"status,omitempty"`
	// Extra message for status, like error details.
	StatusMessage string `json:"statusMessage,omitempty"`
	// Application modules.
	Modules []types.ApplicationModule `json:"modules,omitempty"`
	// Input variables of the revision.
	InputVariables map[string]interface{} `json:"inputVariables,omitempty"`
	// Input plan of the revision.
	InputPlan string `json:"inputPlan,omitempty"`
	// Output of the revision.
	Output string `json:"output,omitempty"`
	// Type of deployer.
	DeployerType string `json:"deployerType,omitempty"`
	// Duration in seconds of the revision deploying.
	Duration int `json:"duration,omitempty"`
}

// Model converts the ApplicationRevisionUpdateInput to ApplicationRevision.
func (in ApplicationRevisionUpdateInput) Model() *ApplicationRevision {
	var entity = &ApplicationRevision{
		ID:             in.ID,
		Status:         in.Status,
		StatusMessage:  in.StatusMessage,
		Modules:        in.Modules,
		InputVariables: in.InputVariables,
		InputPlan:      in.InputPlan,
		Output:         in.Output,
		DeployerType:   in.DeployerType,
		Duration:       in.Duration,
	}
	return entity
}

// ApplicationRevisionOutput is the output for the ApplicationRevision.
type ApplicationRevisionOutput struct {
	// ID holds the value of the "id" field.
	ID oid.ID `json:"id,omitempty"`
	// Status of the resource.
	Status string `json:"status,omitempty"`
	// Extra message for status, like error details.
	StatusMessage string `json:"statusMessage,omitempty"`
	// Describe creation time.
	CreateTime *time.Time `json:"createTime,omitempty"`
	// Application modules.
	Modules []types.ApplicationModule `json:"modules,omitempty"`
	// Type of deployer.
	DeployerType string `json:"deployerType,omitempty"`
	// Duration in seconds of the revision deploying.
	Duration int `json:"duration,omitempty"`
	// Application instance to which the revision belongs.
	Instance *ApplicationInstanceOutput `json:"instance,omitempty"`
	// Environment to which the revision deploys.
	Environment *EnvironmentOutput `json:"environment,omitempty"`
}

// ExposeApplicationRevision converts the ApplicationRevision to ApplicationRevisionOutput.
func ExposeApplicationRevision(in *ApplicationRevision) *ApplicationRevisionOutput {
	if in == nil {
		return nil
	}
	var entity = &ApplicationRevisionOutput{
		ID:            in.ID,
		Status:        in.Status,
		StatusMessage: in.StatusMessage,
		CreateTime:    in.CreateTime,
		Modules:       in.Modules,
		DeployerType:  in.DeployerType,
		Duration:      in.Duration,
		Instance:      ExposeApplicationInstance(in.Edges.Instance),
		Environment:   ExposeEnvironment(in.Edges.Environment),
	}
	if entity.Instance == nil {
		entity.Instance = &ApplicationInstanceOutput{}
	}
	entity.Instance.ID = in.InstanceID
	if entity.Environment == nil {
		entity.Environment = &EnvironmentOutput{}
	}
	entity.Environment.ID = in.EnvironmentID
	return entity
}

// ExposeApplicationRevisions converts the ApplicationRevision slice to ApplicationRevisionOutput pointer slice.
func ExposeApplicationRevisions(in []*ApplicationRevision) []*ApplicationRevisionOutput {
	var out = make([]*ApplicationRevisionOutput, 0, len(in))
	for i := 0; i < len(in); i++ {
		var o = ExposeApplicationRevision(in[i])
		if o == nil {
			continue
		}
		out = append(out, o)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}
