package account

import (
	accountrole "simrs/src/common/account_role"

	"github.com/google/uuid"
)

type getAccountListReqQuery struct {
	Search *string    `query:"search"`
	Limit  *int       `query:"limit"`
	LastID *uuid.UUID `query:"lastId"`
}

type getAccountDetailReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}

type addAccountReq struct {
	Name     *string           `json:"name" validate:"required"`
	Username *string           `json:"username" validate:"required"`
	Password *string           `json:"password" validate:"required"`
	Role     *accountrole.Role `json:"role" validate:"required"`
}

type updateAccountReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}

type updateAccountReq struct {
	Name     *string           `json:"name"`
	Username *string           `json:"username"`
	Password *string           `json:"password"`
	Role     *accountrole.Role `json:"role"`
}

type deleteAccountReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}

type response struct {
	Error      *string     `json:"error,omitempty"`
	Pagination *pagination `json:"pagination,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}

type pagination struct {
	Total *int `json:"total,omitempty"`
}
