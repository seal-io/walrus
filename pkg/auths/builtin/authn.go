package builtin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/pkg/casdoor"
	"github.com/seal-io/seal/pkg/dao/types"
	"github.com/seal-io/seal/pkg/dao/types/oid"
)

func Login(c *gin.Context, username, password string) (sessionValue string, err error) {
	// Login to casdoor.
	cs, err := casdoor.SignInUser(c, casdoor.BuiltinApp, casdoor.BuiltinOrg,
		username, password)
	if err != nil {
		return "", runtime.Error(http.StatusUnauthorized, err)
	}

	// Extract session value.
	sessionValue = casdoor.UnwrapSession(cs)
	if sessionValue != "" {
		return
	}

	return "", runtime.Error(http.StatusInternalServerError, "not found login succeeded token")
}

func Logout(c *gin.Context, sessionValue string) {
	// Wrap session with value.
	s := casdoor.WrapSession(sessionValue)

	// Logout from casdoor.
	_ = casdoor.SignOutUser(c, s)
}

func Validate(c *gin.Context, sid oid.ID, sv string) (domain string, groups []string, name string, err error) {
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

	return
}
