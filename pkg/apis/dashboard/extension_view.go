package dashboard

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/utils/timex"
	"github.com/seal-io/walrus/utils/validation"
)

type RevisionStatusCount struct {
	Running   int `json:"running"`
	Failed    int `json:"failed"`
	Succeeded int `json:"succeeded"`
}

// RevisionStatusStats is the statistics of resource revision status.
type RevisionStatusStats struct {
	RevisionStatusCount

	StartTime string `json:"startTime,omitempty"`
}

type (
	CollectionRouteGetLatestResourceRevisionsRequest struct {
		_ struct{} `route:"GET=/latest-resource-revisions"`

		IsService *bool `query:"isService,omitempty"`

		Context *gin.Context
	}

	CollectionRouteGetLatestResourceRevisionsResponse = []*model.ResourceRevisionOutput
)

func (r *CollectionRouteGetLatestResourceRevisionsRequest) SetGinContext(ctx *gin.Context) {
	r.Context = ctx
}

type (
	CollectionRouteGetBasicInformationRequest struct {
		_ struct{} `route:"GET=/basic-information"`

		WithResourceComponent bool `query:"withResourceComponent,omitempty"`
		WithResourceRevision  bool `query:"withResourceRevision,omitempty"`

		Context *gin.Context
	}

	CollectionRouteGetBasicInformationResponse struct {
		// Project number.
		Project int `json:"project"`
		// Environment number.
		Environment int `json:"environment"`
		// Connector number.
		Connector int `json:"connector"`
		// Service number.
		Service int `json:"service"`
		// Resource number.
		Resource int `json:"resource"`
		// Resource component number.
		ResourceComponent int `json:"resourceComponent,omitempty"`
		// Resource revision number.
		ResourceRevision int `json:"resourceRevision,omitempty"`
	}
)

func (r *CollectionRouteGetBasicInformationRequest) SetGinContext(ctx *gin.Context) {
	r.Context = ctx
}

type (
	CollectionRouteGetResourceRevisionStatisticsRequest struct {
		_ struct{} `route:"POST=/resource-revision-statistics"`

		Step      string    `json:"step"`
		StartTime time.Time `json:"startTime"`
		EndTime   time.Time `json:"endTime"`

		Context *gin.Context
	}

	CollectionRouteGetResourceRevisionStatisticsResponse struct {
		StatusCount *RevisionStatusCount   `json:"statusCount"`
		StatusStats []*RevisionStatusStats `json:"statusStats"`
	}
)

func (r *CollectionRouteGetResourceRevisionStatisticsRequest) Validate() error {
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

func (r *CollectionRouteGetResourceRevisionStatisticsRequest) SetGinContext(ctx *gin.Context) {
	r.Context = ctx
}
