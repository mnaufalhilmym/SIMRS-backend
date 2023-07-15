package patientexamination

import (
	"time"

	"github.com/google/uuid"
)

type getPatientExaminationListQuery struct {
	PatientID *uuid.UUID `query:"patientId"`
	Search    *string    `query:"search"`
	Limit     *int       `query:"limit"`
	LastID    *uuid.UUID `query:"lastId"`
}

type getPatientExaminationDetailReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}

type addPatientExaminationReq struct {
	PatientID       *uuid.UUID `json:"patientId" validate:"required"`
	ExaminationTime *time.Time `json:"examinationTime" validate:"required"`
	Examination     *string    `json:"examination" validate:"required"`
	Treatment       *string    `json:"treatment" validate:"required"`
}

type updatePatientExaminationReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}

type updatePatientExaminationReq struct {
	ExaminationTime *time.Time `json:"examinationTime"`
	Examination     *string    `json:"examination"`
	Treatment       *string    `json:"treatment"`
}

type deletePatientExaminationReqParam struct {
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
