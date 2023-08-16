package auths

import (
	"github.com/gin-gonic/gin"

	"github.com/seal-io/walrus/pkg/auths/builtin"
	"github.com/seal-io/walrus/pkg/auths/session"
	"github.com/seal-io/walrus/pkg/auths/token"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/subject"
	"github.com/seal-io/walrus/pkg/dao/types"
)

func authn(c *gin.Context, mc model.ClientSet, s session.Subject) (session.Subject, error) {
	var err error

	sid, tid, tv := decodeToken(c.Request)
	if tv != "" {
		s.ID = sid
		s.Domain, s.Groups, s.Name, err = token.Validate(c, mc, sid, tid, tv)

		return s, err
	}

	sid, d, sv := decodeSession(c.Request)
	if sv != "" {
		switch d {
		case "", types.SubjectDomainBuiltin:
			s.ID = sid
			s.Domain, s.Groups, s.Name, err = builtin.Validate(c, sid, sv)

			if err != nil {
				revertSession(c.Request, c.Writer)
				return s, err
			}

			_ = flushSession(c.Request, c.Writer)
		default:
			// Anonymous.
		}
	}

	return s, nil
}

func authnSkip(c *gin.Context, mc model.ClientSet, s session.Subject) (session.Subject, error) {
	sid, err := mc.Subjects().Query().
		Where(
			subject.Kind(types.SubjectKindUser),
			subject.Domain(types.SubjectDomainBuiltin),
			subject.Name("admin")).
		OnlyID(c)
	if err != nil {
		return s, err
	}

	s.ID = sid
	s.Domain = types.SubjectDomainBuiltin
	s.Groups = []string{}
	s.Name = "admin"

	return s, nil
}
