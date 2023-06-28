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

type BasicInformation struct {
	// Project number.
	Project int `json:"project"`
	// Service number.
	Service int `json:"service"`
	// Service resource number.
	Resource int `json:"resource"`
	// Service revision number.
	Revision int `json:"revision"`
	// Environment number.
	Environment int `json:"environment"`
	// Connector number.
	Connector int `json:"connector"`
}

// Basic APIs.

// Batch APIs.

type CollectionGetLatestServiceRevisionsRequest struct{}

type CollectionGetLatestServiceRevisionsResponse = []*model.ServiceRevisionOutput

// Extensional APIs.

type BasicInfoRequest struct {
	_ struct{} `route:"GET=/basic-information"`
}

type BasicInfoResponse = BasicInformation

type ServiceRevisionStatisticsRequest struct {
	_ struct{} `route:"POST=/service-revision-statistics"`

	Step      string    `json:"step"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

func (r *ServiceRevisionStatisticsRequest) Validate() error {
	if err := validation.TimeRange(r.StartTime, r.EndTime); err != nil {
		return err
	}

	if r.Step != timex.Day && r.Step != timex.Month && r.Step != timex.Year {
		return errors.New("step must be day, month or year")
	}

	return nil
}

type ServiceRevisionStatisticsResponse struct {
	StatusCount *RevisionStatusCount   `json:"statusCount"`
	StatusStats []*RevisionStatusStats `json:"statusStats"`
}
