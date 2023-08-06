package auths

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/apis/runtime/bind"
	"github.com/seal-io/seal/pkg/auths/builtin"
	"github.com/seal-io/seal/pkg/auths/session"
	tokenbus "github.com/seal-io/seal/pkg/bus/token"
	"github.com/seal-io/seal/pkg/casdoor"
	"github.com/seal-io/seal/pkg/dao/model"
	"github.com/seal-io/seal/pkg/dao/model/predicate"
	"github.com/seal-io/seal/pkg/dao/model/subject"
	"github.com/seal-io/seal/pkg/dao/model/token"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/settings"
	"github.com/seal-io/seal/utils/log"
)

func RequestAccount(mc model.ClientSet, withAuthn bool) Account {
	a := Account{
		modelClient: mc,
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
		modelClient model.ClientSet
		filters     []filter
	}
)

func (a Account) Filter(c *gin.Context) {
	var (
		sj  session.Subject
		err error
	)

	for i := range a.filters {
		sj, err = a.filters[i](c, a.modelClient, sj)
		if err != nil {
			_ = c.Error(err)
			c.Abort()

			return
		}
	}

	session.SetSubject(c, sj)

	c.Next()
}

// Authorize implements the runtime.ResourceRouteAuthorizer interface.
func (a Account) Authorize(c *gin.Context, p runtime.ResourceRouteProfile) int {
	sj, err := session.GetSubject(c)
	if err != nil {
		return http.StatusUnauthorized
	}

	if len(p.Resources) == 0 {
		return http.StatusOK
	}

	allow := sj.Enforce(
		c.Param("project"),
		p.Method,
		p.Resources[len(p.Resources)-1],
		c.Param(p.ResourcePathRefers[len(p.ResourcePathRefers)-1]),
		p.Path)

	if allow {
		return http.StatusOK
	}

	if sj.IsAnonymous() {
		return http.StatusUnauthorized
	}

	return http.StatusForbidden
}

func (a Account) Login(c *gin.Context) error {
	sj, err := session.GetSubject(c)
	if err == nil && !sj.IsAnonymous() {
		// Return account info if already login.
		c.JSON(http.StatusOK, sj)

		return nil
	}

	// Bind and validate.
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if !bind.WithJSON(c, &req) {
		return runtime.Errorc(http.StatusBadRequest)
	}

	if req.Username == "" {
		return runtime.Error(http.StatusBadRequest, "invalid username: blank")
	}

	if req.Password == "" {
		return runtime.Error(http.StatusBadRequest, "invalid password: blank")
	}

	// Get domain form username.
	var (
		sv string

		d string
		u = req.Username
	)

	if ss := strings.SplitN(req.Username, "/", 2); len(ss) == 2 {
		d, u = ss[0], ss[1]
	}

	// Authenticate.
	switch d {
	default:
		// TODO(thxCode): support other authentication system.
		return runtime.Errorf(http.StatusBadRequest, "invalid domain: %s", d)
	case "", types.SubjectDomainBuiltin:
		sv, err = builtin.Login(c, u, req.Password)
		if err != nil {
			return err
		}
		d = types.SubjectDomainBuiltin
	}

	// Get subject id for session.
	sid, err := a.modelClient.Subjects().Query().
		Where(
			subject.Kind(types.SubjectKindUser),
			subject.Domain(d),
			subject.Name(u)).
		OnlyID(c)
	if err != nil {
		// TODO(thxCode): support creating new subject if not found.
		return err
	}

	err = flushSession(c.Request, c.Writer, string(sid), d, sv)
	if err != nil {
		return err
	}

	// Return account info at login.
	sj, err = authz(c, a.modelClient, session.Subject{
		Ctx:    c,
		ID:     sid,
		Domain: d,
		Name:   u,
	})
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, sj)

	return nil
}

func (a Account) Logout(c *gin.Context) error {
	sj, err := session.GetSubject(c)
	if err != nil || sj.IsAnonymous() {
		return runtime.Errorc(http.StatusUnauthorized)
	}

	_, d, sv := decodeSession(c.Request)
	if sv == "" {
		return nil
	}

	switch d {
	default:
		// TODO(thxCode): support other authentication system.
		return nil
	case "", types.SubjectDomainBuiltin:
		builtin.Logout(c, sv)
	}

	revertSession(c.Request, c.Writer)

	return nil
}

func (a Account) GetInfo(c *gin.Context) error {
	sj, err := session.GetSubject(c)
	if err != nil || sj.IsAnonymous() {
		return runtime.Errorc(http.StatusUnauthorized)
	}

	c.JSON(http.StatusOK, sj)

	return nil
}

func (a Account) UpdateInfo(c *gin.Context) error {
	sj, err := session.GetSubject(c)
	if err != nil || sj.IsAnonymous() {
		return runtime.Errorc(http.StatusUnauthorized)
	}

	if sj.Domain != types.SubjectDomainBuiltin {
		return runtime.Error(http.StatusForbidden, "invalid user: not builtin")
	}

	// Bind and validate.
	var req struct {
		Password    string `json:"password,omitempty"`
		OldPassword string `json:"oldPassword,omitempty"`
	}

	if !bind.WithJSON(c, &req) {
		return runtime.Errorc(http.StatusBadRequest)
	}

	if req.Password == "" {
		return runtime.Error(http.StatusBadRequest, "invalid password: blank")
	}

	if req.OldPassword == "" {
		return runtime.Error(http.StatusBadRequest, "invalid old password: blank")
	}

	if req.Password == req.OldPassword {
		return runtime.Error(http.StatusBadRequest, "invalid password: the same")
	}

	// Get casdoor application credential.
	var ac casdoor.ApplicationCredential

	err = settings.CasdoorCred.ValueJSONUnmarshal(c, a.modelClient, &ac)
	if err != nil {
		return err
	}

	// Update casdoor password.
	err = casdoor.UpdateUserPassword(c, ac.ClientID, ac.ClientSecret,
		casdoor.BuiltinOrg, sj.Name, req.OldPassword, req.Password)
	if err != nil {
		if strings.HasSuffix(err.Error(), "not found") {
			return runtime.Error(http.StatusNotFound,
				"invalid user: not found")
		}

		return runtime.Error(http.StatusBadRequest, err)
	}

	if sj.IsAdmin() {
		// Nullify the bootstrap password gain source.
		if settings.BootPwdGainSource.ShouldValue(c, a.modelClient) != "Invalid" {
			_, err = settings.BootPwdGainSource.Set(c, a.modelClient, "Invalid")
			return err
		}
	}

	return nil
}

func (a Account) CreateToken(c *gin.Context) error {
	sj, err := session.GetSubject(c)
	if err != nil || sj.IsAnonymous() {
		return runtime.Errorc(http.StatusUnauthorized)
	}

	// Bind and validate.
	var req struct {
		Name              string `json:"name"`
		ExpirationSeconds *int   `json:"expirationSeconds,omitempty"`
	}

	if !bind.WithJSON(c, &req) {
		return runtime.Errorc(http.StatusBadRequest)
	}

	if req.Name == "" {
		return runtime.Error(http.StatusBadRequest, "invalid name: blank")
	}

	if req.ExpirationSeconds != nil {
		if *req.ExpirationSeconds < 0 {
			return runtime.Error(http.StatusBadRequest, "invalid expiration seconds: negative")
		}
	}

	// Create API token.
	entity, err := CreateAccessToken(c,
		a.modelClient, sj.ID, types.TokenKindAPI, req.Name, req.ExpirationSeconds)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, model.ExposeToken(entity))

	return nil
}

func (a Account) DeleteToken(c *gin.Context) error {
	sj, err := session.GetSubject(c)
	if err != nil || sj.IsAnonymous() {
		return runtime.Errorc(http.StatusUnauthorized)
	}

	// Bind and validate.
	var r model.TokenDeleteInput

	if !bind.WithPath(c, &r) {
		return runtime.Errorc(http.StatusBadRequest)
	}

	if err = r.ValidateWith(c, a.modelClient); err != nil {
		return runtime.Error(http.StatusBadRequest, err)
	}

	// Delete token.
	entity, err := a.modelClient.Tokens().Query().
		Where(
			token.ID(r.ID),
			token.SubjectID(sj.ID)).
		Select(
			token.FieldID,
			token.FieldValue).
		Only(c)
	if err != nil {
		return err
	}

	err = a.modelClient.Tokens().DeleteOneID(r.ID).
		Exec(c)
	if err != nil {
		return err
	}

	if err = tokenbus.Notify(c, model.Tokens{entity}); err != nil {
		// Proceed on clean up failure.
		log.WithName("account").
			Warnf("token post deletion hook failed: %v", err)
	}

	return nil
}

var (
	queryTokenFields = []string{
		token.FieldName,
	}
	getTokenFields = token.WithoutFields(
		token.FieldValue)
	sortTokenFields = []string{
		token.FieldName,
		token.FieldCreateTime,
	}
)

func (a Account) GetTokens(c *gin.Context) error {
	sj, err := session.GetSubject(c)
	if err != nil || sj.IsAnonymous() {
		return runtime.Errorc(http.StatusUnauthorized)
	}

	// Bind and validate.
	var req struct {
		model.TokenQueryInputs `path:",inline" query:",inline"`

		runtime.RequestCollection[
			predicate.Token, token.OrderOption,
		] `query:",inline"`
	}

	if !bind.WithPath(c, &req) || !bind.WithQuery(c, &req) {
		return runtime.Errorc(http.StatusBadRequest)
	}

	if err = req.ValidateWith(c, a.modelClient); err != nil {
		return runtime.Error(http.StatusBadRequest, err)
	}

	query := a.modelClient.Tokens().Query().
		Where(
			token.SubjectID(sj.ID),
			token.Kind(types.TokenKindAPI))

	if queries, ok := req.Querying(queryTokenFields); ok {
		query.Where(queries)
	}

	// Get count.
	cnt, err := query.Clone().Count(c)
	if err != nil {
		return err
	}

	// Get entities.
	if limit, offset, ok := req.Paging(); ok {
		query.Limit(limit).Offset(offset)
	}

	if fields, ok := req.Extracting(getTokenFields, getTokenFields...); ok {
		query.Select(fields...)
	}

	if orders, ok := req.Sorting(sortTokenFields, model.Desc(token.FieldCreateTime)); ok {
		query.Order(orders...)
	}

	entities, err := query.
		Unique(false).
		All(c)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK,
		runtime.PageResponse(
			req.Page, req.PerPage,
			model.ExposeTokens(entities), cnt))

	return nil
}
