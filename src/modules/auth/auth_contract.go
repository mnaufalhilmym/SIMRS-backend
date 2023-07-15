package auth

import (
	accountrole "sim-puskesmas/src/common/account_role"
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
	Name *string           `json:"name"`
	Role *accountrole.Role `json:"role"`
}

type renewTokenRes struct {
	Token *string `json:"token"`
}

type response struct {
	Error *string     `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}
