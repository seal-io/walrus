package builtin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/seal-io/walrus/pkg/casdoor"
	"github.com/seal-io/walrus/pkg/dao/types"
	"github.com/seal-io/walrus/pkg/dao/types/object"
	"github.com/seal-io/walrus/utils/errorx"
)

func Login(c *gin.Context, username, password string) (sessionValue string, err error) {
	// Login to casdoor.
	cs, err := casdoor.SignInUser(c, casdoor.BuiltinApp, casdoor.BuiltinOrg,
		username, password)
	if err != nil {
		return "", errorx.NewHttpError(http.StatusUnauthorized, err.Error())
	}

	// Extract session value.
	sessionValue = casdoor.UnwrapSession(cs)
	if sessionValue != "" {
		return
	}

	return "", errorx.NewHttpError(http.StatusInternalServerError, "not found login succeeded token")
}

func Logout(c *gin.Context, sessionValue string) {
	delCached(c, sessionValue)

	// Wrap session with value.
	s := casdoor.WrapSession(sessionValue)

	// Logout from casdoor.
	_ = casdoor.SignOutUser(c, s)
}

func Validate(c *gin.Context, sid object.ID, sv string) (domain string, groups []string, name string, err error) {
	domain, groups, name, exist := getCached(c, sv)
	if exist {
		return
	}

	// Wrap session with value.
	s := casdoor.WrapSession(sv)

	// Get info from casdoor.
	r, err := casdoor.GetUserInfo(c, s)
	if err != nil {
		return
	}

	domain = types.SubjectDomainBuiltin
	groups = []string{}
	name = r.Name

	cache(c, sv, domain, groups, name)

	return
}
