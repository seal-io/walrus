package variable

import (
	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/variable/view"
	"github.com/seal-io/seal/pkg/auths/session"
	"github.com/seal-io/seal/pkg/dao"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/variable"
)

func Handle(mc model.ClientSet) Handler {
	return Handler{
		modelClient: mc,
	}
}

type Handler struct {
	modelClient model.ClientSet
}

func (h Handler) Kind() string {
	return "Variable"
}

func (h Handler) Validating() any {
	return h.modelClient
}

// Basic APIs.

var queryFields = []string{
	variable.FieldName,
}

func (h Handler) Create(ctx *gin.Context, req view.CreateRequest) (view.CreateResponse, error) {
	entity := req.Model()

	err := h.modelClient.WithTx(ctx, func(tx *model.Tx) error {
		creates, err := dao.VariableCreates(tx, entity)
		if err != nil {
			return err
		}
		entity, err = creates[0].Save(ctx)

		return err
	})
	if err != nil {
		return nil, err
	}

	return view.ExposeVariable(entity), nil
}

func (h Handler) Delete(ctx *gin.Context, req view.DeleteRequest) error {
	return h.modelClient.Variables().DeleteOne(req.Model()).Exec(ctx)
}

func (h Handler) Update(ctx *gin.Context, req view.UpdateRequest) error {
	entity := req.Model()

	return h.modelClient.WithTx(ctx, func(tx *model.Tx) error {
		updates, err := dao.VariableUpdates(tx, entity)
		if err != nil {
			return err
		}

		return updates[0].Exec(ctx)
	})
}

// Batch APIs.

func (h Handler) CollectionDelete(ctx *gin.Context, req view.CollectionDeleteRequest) error {
	return h.modelClient.WithTx(ctx, func(tx *model.Tx) (err error) {
		for i := range req.Items {
			err = tx.Variables().DeleteOne(req.Items[i].Model()).
				Exec(ctx)
			if err != nil {
				return err
			}
		}

		return
	})
}

var getFields = variable.WithoutFields(
	variable.FieldUpdateTime)

func (h Handler) CollectionGet(
	ctx *gin.Context,
	req view.CollectionGetRequest,
) (view.CollectionGetResponse, int, error) {
	s := session.MustGetSubject(ctx)

	s.IncognitoOn()
	defer s.IncognitoOff()

	var (
		// EnvPs is the predicate for query only env scope variables.
		envPs = variable.And(
			variable.ProjectID(req.ProjectID),
			variable.EnvironmentID(req.EnvironmentID),
		)

		// ProjectPs is the predicate for query only project scope variables.
		projPs = variable.And(
			variable.ProjectID(req.ProjectID),
			variable.EnvironmentIDIsNil(),
		)

		// GlobalPs is the predicate for query only global scope variables.
		globalPs = variable.ProjectIDIsNil()
	)

	// Ps is the generated query condition base on input.
	var ps predicate.Variable

	switch {
	case req.EnvironmentID != "":
		// Environment scope.
		switch {
		case req.WithGlobal && req.WithProject:
			// With Project scope and global scope.
			ps = variable.Or(
				envPs,
				projPs,
				globalPs,
			)
		case req.WithGlobal:
			// With Global scope only.
			ps = variable.Or(
				envPs,
				globalPs,
			)
		case req.WithProject:
			// With Project scope only.
			ps = variable.Or(
				envPs,
				projPs,
			)
		default:
			// Environment scope only.
			ps = envPs
		}

	case req.ProjectID != "":
		// Project scope.
		switch {
		case req.WithGlobal:
			// With global scope.
			ps = variable.Or(
				projPs,
				globalPs,
			)
		default:
			// Project scope only.
			ps = projPs
		}
	default:
		// Global scope.
		ps = globalPs
	}

	// Generate query.
	query := h.modelClient.Variables().Query().Where(ps)

	if queries, ok := req.Querying(queryFields); ok {
		query.Where(queries)
	}

	// Get count.
	cnt, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Get entities.
	if limit, offset, ok := req.Paging(); !ok {
		query.Limit(limit).Offset(offset)
	}

	if fields, ok := req.Extracting(getFields, getFields...); ok {
		query.Select(fields...)
	}

	// Allow returning without sorting keys.
	query.Order(model.Desc(variable.FieldCreateTime)).
		Unique(false)

	entities, err := query.All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return view.ExposeVariables(entities), cnt, nil
}

// Extensional APIs.
