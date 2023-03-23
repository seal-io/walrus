package view

import (
	"errors"
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
