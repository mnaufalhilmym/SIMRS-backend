package auth

import (
	accountrole "simrs/src/common/account_role"
)

type signInReq struct {
	Username *string `json:"username" validate:"required"`
	Password *string `json:"password" validate:"required"`
}

type signInRes struct {
	Token *string           `json:"token"`
	Name  *string           `json:"name"`
	Role  *accountrole.Role `json:"role"`
}

type accountRes struct {
	Token *string           `json:"token"`
	Name  *string           `json:"name"`
	Role  *accountrole.Role `json:"role"`
}

type response struct {
	Error *string     `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}
