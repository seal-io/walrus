package token

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/subject"
	"github.com/seal-io/seal/pkg/dao/model/token"
	"github.com/seal-io/seal/pkg/dao/types/object"
)

func Validate(
	c *gin.Context,
	mc model.ClientSet,
	sid, tid object.ID,
	tv string,
) (domain string, groups []string, name string, err error) {
	domain, groups, name, exist := getCached(c, tv)
	if exist {
		return
	}

	t, err := mc.Tokens().Query().
		Where(
			token.ID(tid),
			token.SubjectID(sid)).
		WithSubject(func(sq *model.SubjectQuery) {
			sq.Select(
				subject.FieldDomain,
				subject.FieldName)
		}).
		Only(c)
	if err != nil {
		if model.IsNotFound(err) || model.IsNotLoaded(err) {
			err = nil // Anonymous.
		}

		return
	}

	if string(t.Value) != tv {
		return // Anonymous.
	}

	if t.Expiration != nil && t.Expiration.Before(time.Now()) {
		return // Anonymous.
	}

	domain = t.Edges.Subject.Domain
	groups = []string{}
	name = t.Edges.Subject.Name

	cache(c, tv, domain, groups, name)

	return
}
