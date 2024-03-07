package dashboard

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/utils/timex"
	"github.com/seal-io/walrus/utils/validation"
)

type RunStatusCount struct {
	Planned   int `json:"planned"`
	Canceled  int `json:"canceled"`
	Running   int `json:"running"`
	Failed    int `json:"failed"`
	Succeeded int `json:"succeeded"`
}

// RunStatusStats is the statistics of resource run status.
type RunStatusStats struct {
	RunStatusCount

	StartTime string `json:"startTime,omitempty"`
}

type (
	CollectionRouteGetLatestResourceRunsRequest struct {
		_ struct{} `route:"GET=/latest-resource-runs"`

		Context *gin.Context
	}

	CollectionRouteGetLatestResourceRunsResponse = []*model.ResourceRunOutput
)

func (r *CollectionRouteGetLatestResourceRunsRequest) SetGinContext(ctx *gin.Context) {
	r.Context = ctx
}

type (
	CollectionRouteGetBasicInformationRequest struct {
		_ struct{} `route:"GET=/basic-information"`

		WithResourceComponent bool `query:"withResourceComponent,omitempty"`
		WithResourceRun       bool `query:"withResourceRun,omitempty"`

		Context *gin.Context
	}

	CollectionRouteGetBasicInformationResponse struct {
		// Project number.
		Project int `json:"project"`
		// Environment number.
		Environment int `json:"environment"`
		// Connector number.
		Connector int `json:"connector"`
		// Resource number.
		Resource int `json:"resource"`
		// Resource component number.
		ResourceComponent int `json:"resourceComponent,omitempty"`
		// Resource run number.
		ResourceRun int `json:"resourceRun,omitempty"`
	}
)

func (r *CollectionRouteGetBasicInformationRequest) SetGinContext(ctx *gin.Context) {
	r.Context = ctx
}

type (
	CollectionRouteGetResourceRunStatisticsRequest struct {
		_ struct{} `route:"POST=/resource-run-statistics"`

		Step      string    `json:"step"`
		StartTime time.Time `json:"startTime"`
		EndTime   time.Time `json:"endTime"`

		Context *gin.Context
	}

	CollectionRouteGetResourceRunStatisticsResponse struct {
		StatusCount *RunStatusCount   `json:"statusCount"`
		StatusStats []*RunStatusStats `json:"statusStats"`
	}
)

func (r *CollectionRouteGetResourceRunStatisticsRequest) Validate() error {
	if err := validation.TimeRange(r.StartTime, r.EndTime); err != nil {
		return err
	}

	switch r.Step {
	default:
		return errors.New("invalid step: must be day, month or year")
	case timex.Day, timex.Month, timex.Year:
	}

	return nil
}

func (r *CollectionRouteGetResourceRunStatisticsRequest) SetGinContext(ctx *gin.Context) {
	r.Context = ctx
}
