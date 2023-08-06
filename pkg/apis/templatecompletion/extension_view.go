package templatecompletion

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/dao/types/object"
)

type PromptExample struct {
	Name   string `json:"name"`
	Prompt string `json:"prompt"`
}

type (
	CommonRequest struct {
		Text string `json:"text"`

		Context *gin.Context
	}

	CommonResponse struct {
		Text string `json:"text"`
	}
)

func (r *CommonRequest) Validate() error {
	if r.Text == "" {
		return errors.New("invalid input: blank")
	}

	return nil
}

func (r *CommonRequest) SetGinContext(ctx *gin.Context) {
	r.Context = ctx
}

type (
	CollectionRouteGetExampleRequest struct {
		_ struct{} `route:"GET=/examples"`

		Context *gin.Context
	}

	CollectionRouteGetExampleResponse = []PromptExample
)

func (r *CollectionRouteGetExampleRequest) SetGinContext(ctx *gin.Context) {
	r.Context = ctx
}

type (
	CollectionRouteGenerateRequest struct {
		_ struct{} `route:"POST=/generate"`

		CommonRequest `json:",inline"`
	}

	CollectionRouteGenerateResponse = CommonResponse
)

type (
	CollectionRouteExplainRequest struct {
		_ struct{} `route:"POST=/explain"`

		CommonRequest `json:",inline"`
	}

	CollectionRouteExplainResponse = CommonResponse
)

type (
	CollectionRouteCorrectRequest struct {
		_ struct{} `route:"POST=/correct"`

		CommonRequest `json:",inline"`
	}

	CollectionRouteCorrectResponse struct {
		Corrected   string `json:"corrected"`
		Explanation string `json:"explanation"`
	}
)

type (
	CollectionRouteCreatePrRequest struct {
		_ struct{} `route:"POST=/create-pr"`

		ConnectorID object.ID `json:"connectorID"`
		Repository  string    `json:"repository"`
		Branch      string    `json:"branch"`
		Path        string    `json:"path"`
		Content     string    `json:"content"`

		Context *gin.Context
	}

	CollectionRouteCreatePrResponse struct {
		Link string `json:"link"`
	}
)

func (r *CollectionRouteCreatePrRequest) Validate() error {
	if !r.ConnectorID.Valid() {
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

func (r *CollectionRouteCreatePrRequest) SetGinContext(ctx *gin.Context) {
	r.Context = ctx
}
