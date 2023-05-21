package casdoor

import (
	"github.com/seal-io/seal/utils/vars"
)

var endpoint = &vars.SetOnce[string]{}

const (
	BuiltinOrg          = "built-in"
	BuiltinApp          = "app-built-in"
	BuiltinAdmin        = "admin"
	BuiltinAdminInitPwd = "Admin@123"
)
