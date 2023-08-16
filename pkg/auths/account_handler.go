package auths

import (
	"net/http"
	"strings"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/pkg/auths/builtin"
	"github.com/seal-io/walrus/pkg/auths/session"
	tokenbus "github.com/seal-io/walrus/pkg/bus/token"
	"github.com/seal-io/walrus/pkg/casdoor"
	"github.com/seal-io/walrus/pkg/dao/model"
	"github.com/seal-io/walrus/pkg/dao/model/subject"
	"github.com/seal-io/walrus/pkg/dao/model/token"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/settings"
	"github.com/seal-io/walrus/utils/log"
)

// Login logins a subject with the given username and password.
func (a Account) Login(req LoginRequest) (LoginResponse, error) {
	sj, err := session.GetSubject(req.Context)
	if err == nil && !sj.IsAnonymous() {
		return &sj, nil
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
		return nil, runtime.Errorf(http.StatusBadRequest, "invalid domain: %s", d)
	case "", types.SubjectDomainBuiltin:
		sv, err = builtin.Login(req.Context, u, req.Password)
		if err != nil {
			return nil, err
		}
		d = types.SubjectDomainBuiltin
	}

	// Get subject id for session.
	sid, err := a.modelClient.Subjects().Query().
		Where(
			subject.Kind(types.SubjectKindUser),
			subject.Domain(d),
			subject.Name(u)).
		OnlyID(req.Context)
	if err != nil {
		// TODO(thxCode): support creating new subject if not found.
		return nil, err
	}

	err = flushSession(req.Context.Request, req.Context.Writer, string(sid), d, sv)
	if err != nil {
		return nil, err
	}

	// Return account info at login.
	sj, err = authz(req.Context, a.modelClient, session.Subject{
		Ctx:    req.Context,
		ID:     sid,
		Domain: d,
		Name:   u,
	})
	if err != nil {
		return nil, err
	}

	return &sj, nil
}

// Logout logouts the session subject.
func (a Account) Logout(req LogoutRequest) error {
	_, d, sv := decodeSession(req.Context.Request)
	if sv == "" {
		return nil
	}

	switch d {
	default:
		// TODO(thxCode): support other authentication system.
		return nil
	case "", types.SubjectDomainBuiltin:
		builtin.Logout(req.Context, sv)
	}

	revertSession(req.Context.Request, req.Context.Writer)

	return nil
}

// GetInfo returns the session subject.
func (a Account) GetInfo(req GetInfoRequest) (GetInfoResponse, error) {
	sj := session.MustGetSubject(req.Context)
	return &sj, nil
}

// UpdateInfo updates the session subject.
func (a Account) UpdateInfo(req UpdateInfoRequest) error {
	sj := session.MustGetSubject(req.Context)

	if sj.Domain != types.SubjectDomainBuiltin {
		// Nothing to do with other authentication system.
		return nil
	}

	// Get casdoor application credential.
	var ac casdoor.ApplicationCredential

	err := settings.CasdoorCred.ValueJSONUnmarshal(req.Context, a.modelClient, &ac)
	if err != nil {
		return err
	}

	// Update casdoor password.
	err = casdoor.UpdateUserPassword(req.Context, ac.ClientID, ac.ClientSecret,
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
		if settings.BootPwdGainSource.ShouldValue(req.Context, a.modelClient) != "Invalid" {
			_, err = settings.BootPwdGainSource.Set(req.Context, a.modelClient, "Invalid")
			return err
		}
	}

	return nil
}

// CreateToken creates a new API token for the session subject.
func (a Account) CreateToken(req CreateTokenRequest) (CreateTokenResponse, error) {
	sj := session.MustGetSubject(req.Context)

	// Create API token.
	entity, err := CreateAccessToken(req.Context,
		a.modelClient, sj.ID, types.TokenKindAPI, req.Name, req.ExpirationSeconds)
	if err != nil {
		return nil, err
	}

	return model.ExposeToken(entity), nil
}

// DeleteToken deletes a API token of the session subject.
func (a Account) DeleteToken(req DeleteTokenRequest) error {
	sj := session.MustGetSubject(req.Context)

	entity, err := a.modelClient.Tokens().Query().
		Where(
			token.ID(req.ID),
			token.SubjectID(sj.ID)).
		Select(
			token.FieldID,
			token.FieldValue).
		Only(req.Context)
	if err != nil {
		return err
	}

	err = a.modelClient.Tokens().DeleteOneID(req.ID).
		Exec(req.Context)
	if err != nil {
		return err
	}

	if err = tokenbus.Notify(req.Context, model.Tokens{entity}); err != nil {
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

func (a Account) GetTokens(req GetTokensRequest) (GetTokensResponse, int, error) {
	sj := session.MustGetSubject(req.Context)

	query := a.modelClient.Tokens().Query().
		Where(
			token.SubjectID(sj.ID),
			token.Kind(types.TokenKindAPI))

	if queries, ok := req.Querying(queryTokenFields); ok {
		query.Where(queries)
	}

	// Get count.
	cnt, err := query.Clone().Count(req.Context)
	if err != nil {
		return nil, 0, err
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
		All(req.Context)
	if err != nil {
		return nil, 0, err
	}

	return model.ExposeTokens(entities), cnt, nil
}
