package resourcerun

import (
	"errors"
	"fmt"
	"strings"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/predicate"
	"github.com/seal-io/walrus/pkg/dao/model/resourcerun"
	"github.com/seal-io/walrus/pkg/dao/types/object"
)

type (
	GetRequest = model.ResourceRunQueryInput

	GetResponse = *model.ResourceRunOutput
)

type (
	CollectionFieldQuery struct {
		QueryID      object.ID `query:"id"`
		QueryType    string    `query:"type"`
		QueryStatus  string    `query:"status"`
		QueryPreview *bool     `query:"preview"`
		QueryLabels  []string  `query:"labels"`
	}

	CollectionGetRequest struct {
		model.ResourceRunQueryInputs `path:",inline" query:",inline"`

		runtime.RequestCollection[
			predicate.ResourceRun, resourcerun.OrderOption,
		] `query:",inline"`

		Stream *runtime.RequestUnidiStream

		CollectionFieldQuery `query:",inline"`
		Rollback             bool `query:"rollback"`
	}

	CollectionGetResponse = []*model.ResourceRunOutput
)

func (r *CollectionGetRequest) SetStream(stream runtime.RequestUnidiStream) {
	r.Stream = &stream
}

func (q *CollectionFieldQuery) Queries() (queries []predicate.ResourceRun, ok bool) {
	if q.QueryID != "" {
		queries = append(queries, resourcerun.ID(q.QueryID))
		ok = true
	}

	if q.QueryType != "" {
		queries = append(queries, resourcerun.Type(q.QueryType))
		ok = true
	}

	if q.QueryStatus != "" {
		queries = append(queries, func(s *sql.Selector) {
			s.Where(sqljson.ValueEQ(
				resourcerun.FieldStatus,
				q.QueryStatus,
				sqljson.Path("summaryStatus"),
			))
		})

		ok = true
	}

	if q.QueryPreview != nil {
		queries = append(queries, resourcerun.Preview(*q.QueryPreview))
		ok = true
	}

	if len(q.QueryLabels) != 0 {
		labels := make(map[string]string)
		for _, ql := range q.QueryLabels {
			arr := strings.Split(ql, "=")
			if len(arr) != 2 {
				continue
			}

			labels[arr[0]] = arr[1]
		}

		var ps []*sql.Predicate
		for k, v := range labels {
			ps = append(ps, sqljson.ValueEQ(
				resourcerun.FieldLabels,
				v,
				sqljson.Path(k),
			))
		}

		queries = append(queries, func(s *sql.Selector) {
			s.Where(sql.And(ps...))
		})
		ok = true
	}

	return
}

type CollectionDeleteRequest struct {
	model.ResourceRunDeleteInputs `path:",inline" json:",inline"`
}

func (r *CollectionDeleteRequest) Validate() error {
	if err := r.ResourceRunDeleteInputs.Validate(); err != nil {
		return err
	}

	latestRun, err := r.Client.ResourceRuns().Query().
		Where(resourcerun.ResourceID(r.Resource.ID)).
		Order(model.Desc(resourcerun.FieldCreateTime)).
		Select(resourcerun.FieldID).
		First(r.Context)
	if err != nil {
		return fmt.Errorf("failed to get latest run: %w", err)
	}

	for i := range r.Items {
		// Prevent deleting the latest run.
		if r.Items[i].ID == latestRun.ID {
			return errors.New("invalid ids: can not delete latest run")
		}
	}

	return nil
}
