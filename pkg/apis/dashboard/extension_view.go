package dashboard

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/utils/timex"
	"github.com/seal-io/seal/utils/validation"
)

type RevisionStatusCount struct {
	Running int `json:"running"`
	Failed  int `json:"failed"`
	Succeed int `json:"succeed"`
}

// RevisionStatusStats is the statistics of service revision status.
type RevisionStatusStats struct {
	RevisionStatusCount

	StartTime string `json:"startTime,omitempty"`
}

type (
	CollectionRouteGetLatestServiceRevisionsRequest struct {
		_ struct{} `route:"GET=/latest-service-revisions"`

		Context *gin.Context
	}

	CollectionRouteGetLatestServiceRevisionsResponse = []*model.ServiceRevisionOutput
)

func (r *CollectionRouteGetLatestServiceRevisionsRequest) SetGinContext(ctx *gin.Context) {
	r.Context = ctx
}

type (
	CollectionRouteGetBasicInformationRequest struct {
		_ struct{} `route:"GET=/basic-information"`

		WithServiceResource bool `query:"withServiceResource,omitempty"`
		WithServiceRevision bool `query:"withServiceRevision,omitempty"`

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
		// Service resource number.
		ServiceResource int `json:"serviceResource,omitempty"`
		// Service revision number.
		ServiceRevision int `json:"serviceRevision,omitempty"`
	}
)

func (r *CollectionRouteGetBasicInformationRequest) SetGinContext(ctx *gin.Context) {
	r.Context = ctx
}

type (
	CollectionRouteGetServiceRevisionStatisticsRequest struct {
		_ struct{} `route:"POST=/service-revision-statistics"`

		Step      string    `json:"step"`
		StartTime time.Time `json:"startTime"`
		EndTime   time.Time `json:"endTime"`

		Context *gin.Context
	}

	CollectionRouteGetServiceRevisionStatisticsResponse struct {
		StatusCount *RevisionStatusCount   `json:"statusCount"`
		StatusStats []*RevisionStatusStats `json:"statusStats"`
	}
)

func (r *CollectionRouteGetServiceRevisionStatisticsRequest) Validate() error {
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

func (r *CollectionRouteGetServiceRevisionStatisticsRequest) SetGinContext(ctx *gin.Context) {
	r.Context = ctx
}
