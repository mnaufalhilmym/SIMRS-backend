package patient

import (
	"time"

	"github.com/google/uuid"
)

type getPatientListReqQuery struct {
	SearchByDistrictID *uuid.UUID `query:"searchByDistrictId"`
	Search             *string    `query:"search"`
	Limit              *int       `query:"limit"`
	LastID             *uuid.UUID `query:"lastId"`
}

type getPatientDetailReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}

type getPatientCountReqQuery struct {
	DistrictID *uuid.UUID `query:"districtId"`
}

type addPatientReq struct {
	MedicalRecordNumber            *string    `json:"medicalRecordNumber" validate:"required"`
	FamilyCardNumber               *string    `json:"familyCardNumber" validate:"required"`
	PopulationIdentificationNumber *string    `json:"populationIdentificationNumber" validate:"required"`
	Name                           *string    `json:"name" validate:"required"`
	Gender                         *gender    `json:"gender" validate:"required"`
	PlaceOfBirth                   *string    `json:"placeOfBirth" validate:"required"`
	DateOfBirth                    *time.Time `json:"dateOfBirth" validate:"required"`
	Address                        *string    `json:"address" validate:"required"`
	DistrictID                     *uuid.UUID `json:"districtId" validate:"required"`
	Job                            *string    `json:"job" validate:"required"`
	Religion                       *string    `json:"religion" validate:"required"`
	BloodGroup                     *string    `json:"bloodGroup" validate:"required"`
	Insurence                      *string    `json:"insurence" validate:"required"`
	InsurenceNumber                *string    `json:"insurenceNumber" validate:"required"`
	Phone                          *string    `json:"phone" validate:"required"`
}

type updatePatientReqParam struct {
	ID *uuid.UUID `params:"id" validate:"required"`
}

type updatePatientReq struct {
	MedicalRecordNumber            *string    `json:"medicalRecordNumber"`
	FamilyCardNumber               *string    `json:"familyCardNumber"`
	PopulationIdentificationNumber *string    `json:"populationIdentificationNumber"`
	Name                           *string    `json:"name"`
	Gender                         *gender    `json:"gender"`
	PlaceOfBirth                   *string    `json:"placeOfBirth"`
	DateOfBirth                    *time.Time `json:"dateOfBirth"`
	Address                        *string    `json:"address"`
	DistrictID                     *uuid.UUID `json:"districtId"`
	Job                            *string    `json:"job"`
	Religion                       *string    `json:"religion"`
	BloodGroup                     *string    `json:"bloodGroup"`
	Insurence                      *string    `json:"insurence"`
	InsurenceNumber                *string    `json:"insurenceNumber"`
	Phone                          *string    `json:"phone"`
}

type deletePatientReqParam struct {
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
