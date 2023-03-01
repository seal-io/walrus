package view

import (
	"errors"

	"github.com/seal-io/seal/pkg/dao/schema"
)

type UpdateInfoRequest struct {
	LoginGroup  *string `json:"loginGroup,omitempty"`
	Password    *string `json:"password,omitempty"`
	OldPassword *string `json:"oldPassword,omitempty"`
}

func (r *UpdateInfoRequest) Validate() error {
	var needUpdate bool
	if r.LoginGroup != nil {
		if *r.LoginGroup == "" {
			return errors.New("invalid group: blank")
		}
		needUpdate = true
	}
	if r.Password != nil {
		if *r.Password == "" {
			return errors.New("invalid password: blank")
		}
		if r.OldPassword == nil || *r.OldPassword == "" {
			return errors.New("invalid old password: blank")
		}
		if *r.OldPassword == *r.Password {
			return errors.New("invalid password: the same")
		}
		needUpdate = true
	}
	if !needUpdate {
		return errors.New("invalid input: nothing update")
	}
	return nil
}

type GetInfoResponse struct {
	Name       string              `json:"name"`
	Roles      schema.SubjectRoles `json:"roles"`
	Policies   schema.RolePolicies `json:"policies"`
	Groups     []GetInfoGroup      `json:"groups"`
	LoginGroup string              `json:"loginGroup"`
}

type GetInfoGroup struct {
	Name  string   `json:"name"`
	Paths []string `json:"paths"`
}
