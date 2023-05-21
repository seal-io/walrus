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

// RevisionStatusStats is the statistics of application revision status.
type RevisionStatusStats struct {
	RevisionStatusCount

	StartTime string `json:"startTime,omitempty"`
}

type BasicInformation struct {
	// Application number.
	Application int `json:"application"`
	//  module number
	Module int `json:"module"`
	// Instance number.
	Instance int `json:"instance"`
	// Application resource number.
	Resource int `json:"resource"`
	// Application revision number.
	Revision int `json:"revision"`
	// Environment number.
	Environment int `json:"environment"`
	// Connector number.
	Connector int `json:"connector"`
}

// Basic APIs.

// Batch APIs.

type CollectionGetLatestApplicationRevisionsRequest struct{}

type CollectionGetLatestApplicationRevisionsResponse = []*model.ApplicationRevisionOutput

// Extensional APIs.

type BasicInfoRequest struct {
	_ struct{} `route:"GET=/basic-information"`
}

type BasicInfoResponse = BasicInformation

type ApplicationRevisionStatisticsRequest struct {
	_ struct{} `route:"POST=/application-revision-statistics"`

	Step      string    `json:"step"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

func (r *ApplicationRevisionStatisticsRequest) Validate() error {
	if err := validation.TimeRange(r.StartTime, r.EndTime); err != nil {
		return err
	}

	if r.Step != timex.Day && r.Step != timex.Month && r.Step != timex.Year {
		return errors.New("step must be day, month or year")
	}

	return nil
}

type ApplicationRevisionStatisticsResponse struct {
	StatusCount *RevisionStatusCount   `json:"statusCount"`
	StatusStats []*RevisionStatusStats `json:"statusStats"`
}
