package view

import (
	"errors"
	"time"

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

// Basic APIs.

// Batch APIs.

type CollectionGetLatestServiceRevisionsRequest struct{}

type CollectionGetLatestServiceRevisionsResponse = []*model.ServiceRevisionOutput

// Extensional APIs.

type CollectionRouteBasicInformationRequest struct {
	_ struct{} `route:"GET=/basic-information"`

	WithServiceResource bool `query:"withServiceResource,omitempty"`
	WithServiceRevision bool `query:"withServiceRevision,omitempty"`
}

type CollectionRouteBasicInformationResponse struct {
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

type CollectionRouteServiceRevisionStatisticsRequest struct {
	_ struct{} `route:"POST=/service-revision-statistics"`

	Step      string    `json:"step"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

func (r *CollectionRouteServiceRevisionStatisticsRequest) Validate() error {
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

type CollectionRouteServiceRevisionStatisticsResponse struct {
	StatusCount *RevisionStatusCount   `json:"statusCount"`
	StatusStats []*RevisionStatusStats `json:"statusStats"`
}
