package variable

import (
	"fmt"

	"entgo.io/ent/dialect/sql"
	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/variable/view"
	"github.com/seal-io/seal/pkg/auths/session"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/service"
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

func (h Handler) Create(ctx *gin.Context, req view.CreateRequest) (view.CreateResponse, error) {
	entity := req.Model()

	entity, err := h.modelClient.Variables().Create().
		Set(entity).
		Save(ctx)
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

	return h.modelClient.Variables().UpdateOne(entity).
		Set(entity).
		Exec(ctx)
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

var sortFields = []string{
	variable.FieldName,
	variable.FieldCreateTime,
}

var getFields = variable.WithoutFields(variable.FieldUpdateTime)

func (h Handler) CollectionGet(
	ctx *gin.Context,
	req view.CollectionGetRequest,
) (view.CollectionGetResponse, int, error) {
	s := session.MustGetSubject(ctx)

	s.IncognitoOn()
	defer s.IncognitoOff()

	/* SQL generate for count.

	SELECT
	  COUNT("variables"."id")
	FROM
	  "variables"
	  JOIN (
	    SELECT
	      "alias"."name",
	      MAX(
	        CASE WHEN project_id IS NOT NULL
	        AND environment_id IS NOT NULL THEN 3 WHEN project_id IS NOT NULL
	        AND environment_id IS NULL THEN 2 ELSE 1 END
	      ) AS scope
	    FROM
	      "variables" AS "alias"
	    WHERE
	      (
	        "project_id" = $1
	        AND "environment_id" = $2
	      )
	      OR (
	        "project_id" = $3
	        AND "environment_id" IS NULL
	      )
	      OR "project_id" IS NULL
	    GROUP BY
	      "alias"."name"
	  ) AS "alias" ON (
	    (
	      "project_id" = $4
	      AND "environment_id" = $5
	    )
	    OR (
	      "project_id" = $6
	      AND "environment_id" IS NULL
	    )
	    OR "project_id" IS NULL
	  )
	  AND "variables"."name" = "alias"."name"
	  AND (
	    CASE WHEN project_id IS NOT NULL
	    AND environment_id IS NOT NULL THEN 3 WHEN project_id IS NOT NULL
	    AND environment_id IS NULL THEN 2 ELSE 1 END
	  ) = alias.scope
	*/

	/* SQL generate for query.

	SELECT
	  "variables"."id",
	  "variables"."create_time",
	  "variables"."project_id",
	  "variables"."name",
	  "variables"."value",
	  "variables"."sensitive",
	  "variables"."description",
	  "variables"."environment_id"
	FROM
	  "variables"
	  JOIN (
	    SELECT
	      "alias"."name",
	      MAX(
	        CASE WHEN project_id IS NOT NULL
	        AND environment_id IS NOT NULL THEN 3 WHEN project_id IS NOT NULL
	        AND environment_id IS NULL THEN 2 ELSE 1 END
	      ) AS scope
	    FROM
	      "variables" AS "alias"
	    WHERE
	      (
	        (
	          "project_id" = $1
	          AND "environment_id" = $2
	        )
	        OR (
	          "project_id" = $3
	          AND "environment_id" IS NULL
	        )
	        OR "project_id" IS NULL
	      )
	    GROUP BY
	      "alias"."name"
	  ) AS "alias" ON (
	    (
	      "project_id" = $4
	      AND "environment_id" = $5
	    )
	    OR (
	      "project_id" = $6
	      AND "environment_id" IS NULL
	    )
	    OR "project_id" IS NULL
	  )
	  AND "variables"."name" = "alias"."name"
	  AND (
	    CASE WHEN project_id IS NOT NULL
	    AND environment_id IS NOT NULL THEN 3 WHEN project_id IS NOT NULL
	    AND environment_id IS NULL THEN 2 ELSE 1 END
	  ) = alias.scope
	ORDER BY
	  "variables"."create_time" DESC
	LIMIT
	  100
	*/

	var (
		// EnvPs is the predicate for query only env scope variables.
		envPs = sql.And(
			sql.EQ(variable.FieldProjectID, req.ProjectID),
			sql.EQ(variable.FieldEnvironmentID, req.EnvironmentID),
		)

		// ProjectPs is the predicate for query only project scope variables.
		projPs = sql.And(
			sql.EQ(variable.FieldProjectID, req.ProjectID),
			sql.IsNull(variable.FieldEnvironmentID),
		)

		// GlobalPs is the predicate for query only global scope variables.
		globalPs = sql.IsNull(variable.FieldProjectID)
	)

	// Ps is the generated query condition base on input.
	var ps *sql.Predicate

	switch {
	case req.EnvironmentID != "":
		// Environment scope.
		// Environment scope only.
		ps = envPs

		// With Project scope and global scope.
		if req.IncludeInherited {
			ps = sql.Or(
				envPs,
				projPs,
				globalPs,
			)
		}
	case req.ProjectID != "":
		// Project scope.
		// Project scope only.
		ps = projPs

		// With global scope.
		if req.IncludeInherited {
			ps = sql.Or(
				projPs,
				globalPs,
			)
		}
	default:
		// Global scope.
		ps = globalPs
	}

	scopeExpr := `CASE WHEN project_id IS NOT NULL AND environment_id IS NOT NULL THEN 3
                                          WHEN project_id IS NOT NULL AND environment_id IS NULL THEN 2
                                          ELSE 1 END`

	modifier := func(s *sql.Selector) {
		alias := sql.Table(variable.Table).
			As("alias")

		subQuery := sql.Select(alias.C(variable.FieldName)).
			AppendSelectExpr(sql.Expr(fmt.Sprintf("MAX(%s) AS scope", scopeExpr))).
			From(alias).
			Where(ps).
			GroupBy(alias.C(variable.FieldName)).
			As("alias")

		s.Join(subQuery).
			OnP(
				sql.And(
					ps,
					sql.ColumnsEQ(
						s.C(variable.FieldName),
						alias.C(variable.FieldName),
					),
					sql.ExprP(fmt.Sprintf("(%s) = alias.scope", scopeExpr)),
				),
			)
	}

	// Generate query.
	query := h.modelClient.Variables().Query().
		Modify(modifier)

	// Search query.
	if req.Query != nil {
		query.Where(variable.NameContainsFold(*req.Query))
	}

	// Get count.
	cnt, err := query.Clone().
		Modify(modifier).
		Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Get entities.
	if fields, ok := req.Extracting(getFields, getFields...); ok {
		query.Select(fields...)
	}

	if limit, offset, ok := req.Paging(); ok {
		query.Limit(limit).Offset(offset)
	}

	if orders, ok := req.Sorting(sortFields, model.Desc(service.FieldCreateTime)); ok {
		query.Order(orders...).
			Unique(false)
	}

	vars, err := query.All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return view.ExposeVariables(vars), cnt, err
}

// Extensional APIs.
