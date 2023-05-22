package view

import (
	"errors"

	"github.com/seal-io/seal/pkg/dao/types/oid"
)

type CompletionResponse struct {
	Text string `json:"text"`
}

type GenerateRequest struct {
	_ struct{} `route:"POST=/generate"`

	Text string `json:"text"`
}

func (r *GenerateRequest) Validate() error {
	if r.Text == "" {
		return errors.New("invalid input: blank")
	}

	return nil
}

type GenerateResponse = CompletionResponse

type ExplainRequest struct {
	_ struct{} `route:"POST=/explain"`

	Text string `json:"text"`
}

func (r *ExplainRequest) Validate() error {
	if r.Text == "" {
		return errors.New("invalid input: blank")
	}

	return nil
}

type ExplainResponse = CompletionResponse

type CorrectRequest struct {
	_ struct{} `route:"POST=/correct"`

	Text string `json:"text"`
}

func (r *CorrectRequest) Validate() error {
	if r.Text == "" {
		return errors.New("invalid input: blank")
	}

	return nil
}

type CorrectResponse struct {
	Corrected   string `json:"corrected"`
	Explanation string `json:"explanation"`
}

type ExampleRequest struct {
	_ struct{} `route:"GET=/examples"`
}

type ExampleResponse []ModuleCompletionPromptExample

type ModuleCompletionPromptExample struct {
	Name   string `json:"name"`
	Prompt string `json:"prompt"`
}

type CreatePrRequest struct {
	_ struct{} `route:"POST=/create-pr"`

	ConnectorID oid.ID `json:"connectorID"`
	Repository  string `json:"repository"`
	Branch      string `json:"branch"`
	Path        string `json:"path"`
	Content     string `json:"content"`
}

func (r *CreatePrRequest) Validate() error {
	if !r.ConnectorID.Valid(0) {
		return errors.New("invalid connector id: blank")
	}

	if r.Repository == "" {
		return errors.New("invalid repository: blank")
	}

	if r.Branch == "" {
		return errors.New("invalid branch: blank")
	}

	if r.Path == "" {
		return errors.New("invalid path: blank")
	}

	if r.Content == "" {
		return errors.New("invalid content: blank")
	}

	return nil
}

type CreatePrResponse struct {
	Link string `json:"link"`
}
