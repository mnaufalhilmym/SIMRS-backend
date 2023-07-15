package district

import "github.com/google/uuid"

type getDistrictListReqQuery struct {
	Search *string    `query:"search"`
	Limit  *int       `query:"limiy"`
	LastID *uuid.UUID `query:"lastId"`
}

type getDistrictDetailReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}

type addDistrictReq struct {
	Name *string `json:"name" validate:"required"`
}

type updateDistrictReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}

type updateDistrictReq struct {
	Name *string `json:"name"`
}

type deleteDistrictReqParam struct {
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
