package auths

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/auths/builtin"
	"github.com/seal-io/seal/pkg/auths/session"
	"github.com/seal-io/seal/pkg/casdoor"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/subject"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/settings"
)

func RequestAccount(mc model.ClientSet, withAuthn bool) Account {
	a := Account{
		mc: mc,
		filters: []filter{
			authn,
			authz,
		},
	}
	if !withAuthn {
		a.filters[0] = authnSkip
	}

	return a
}

type (
	filter  func(*gin.Context, model.ClientSet, session.Subject) (session.Subject, error)
	Account struct {
		mc      model.ClientSet
		filters []filter
	}
)

func (a Account) Filter(c *gin.Context) {
	var (
		s   session.Subject
		err error
	)

	for i := range a.filters {
		s, err = a.filters[i](c, a.mc, s)
		if err != nil {
			_ = c.Error(err)
			c.Abort()

			return
		}
	}

	session.SetSubject(c, s)

	c.Next()
}

func (a Account) Login(c *gin.Context) (err error) {
	var r struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err = c.ShouldBindJSON(&r); err != nil {
		return err
	}

	var (
		sv string

		d string
		u = r.Username
	)

	if ss := strings.SplitN(r.Username, "/", 2); len(ss) == 2 {
		d, u = ss[0], ss[1]
	}

	switch d {
	default:
		// TODO(thxCode): support other authentication system.
		return runtime.Errorf(http.StatusBadRequest, "invalid domain: %s", d)
	case "", types.SubjectDomainBuiltin:
		sv, err = builtin.Login(c, u, r.Password)
		if err != nil {
			return err
		}
		d = types.SubjectDomainBuiltin
	}

	sid, err := a.mc.Subjects().Query().
		Where(
			subject.Kind(types.SubjectKindUser),
			subject.Domain(d),
			subject.Name(u)).
		OnlyID(c)
	if err != nil {
		// TODO(thxCode): support creating new subject if not found.
		return err
	}

	return flushSession(c.Request, c.Writer, string(sid), d, sv)
}

func (a Account) Logout(c *gin.Context) {
	_, d, sv := decodeSession(c.Request)
	if sv == "" {
		return
	}

	switch d {
	default:
		// TODO(thxCode): support other authentication system.
		return
	case "", types.SubjectDomainBuiltin:
		builtin.Logout(c, sv)
	}

	revertSession(c.Request, c.Writer)
}

func (a Account) GetInfo(c *gin.Context) error {
	s := session.MustGetSubject(c)
	if s.IsAnonymous() {
		return runtime.Errorc(http.StatusUnauthorized)
	}

	c.JSON(http.StatusOK, s)

	return nil
}

func (a Account) UpdateInfo(c *gin.Context) error {
	s := session.MustGetSubject(c)
	if s.IsAnonymous() {
		return runtime.Errorc(http.StatusUnauthorized)
	}

	if s.Domain != types.SubjectDomainBuiltin {
		return runtime.Error(http.StatusForbidden, "invalid user: not builtin")
	}

	var r struct {
		Password    string `json:"password,omitempty"`
		OldPassword string `json:"oldPassword,omitempty"`
	}

	if err := c.ShouldBindJSON(&r); err != nil {
		return err
	}

	if r.Password == "" {
		return runtime.Error(http.StatusBadRequest, "invalid password: blank")
	}

	if r.OldPassword == "" {
		return runtime.Error(http.StatusBadRequest, "invalid old password: blank")
	}

	if r.Password == r.OldPassword {
		return runtime.Error(http.StatusBadRequest, "invalid password: the same")
	}

	var ac casdoor.ApplicationCredential

	err := settings.CasdoorCred.ValueJSONUnmarshal(c, a.mc, &ac)
	if err != nil {
		return err
	}

	err = casdoor.UpdateUserPassword(c, ac.ClientID, ac.ClientSecret,
		casdoor.BuiltinOrg, s.Name, r.OldPassword, r.Password)
	if err != nil {
		if strings.HasSuffix(err.Error(), "not found") {
			return runtime.Error(http.StatusNotFound,
				"invalid user: not found")
		}

		return runtime.Error(http.StatusBadRequest, err)
	}

	if s.IsAdmin() {
		// Nullify the bootstrap password gain source.
		if settings.BootPwdGainSource.ShouldValue(c, a.mc) != "Invalid" {
			_, err = settings.BootPwdGainSource.Set(c, a.mc, "Invalid")
			return err
		}
	}

	return nil
}
