package patient

import (
	"simrs/src/libs/db/pg"
	"time"

	"github.com/google/uuid"
)

type PatientModel struct {
	*pg.Model
	MedicalRecordNumber            *string               `json:"medicalRecordNumber"`
	FamilyCardNumber               *string               `json:"familyCardNumber"`
	RelationshipInFamily           *relationshipInFamily `json:"relationshipInFamily"`
	PopulationIdentificationNumber *string               `json:"populationIdentificationNumber"`
	Salutation                     *salutation           `json:"salutation"`
	Name                           *string               `json:"name"`
	Gender                         *gender               `json:"gender"`
	PlaceOfBirth                   *string               `json:"placeOfBirth"`
	DateOfBirth                    *time.Time            `json:"dateOfBirth"`
	Address                        *string               `json:"address"`
	DistrictID                     *uuid.UUID            `json:"districtId"`
	Job                            *string               `json:"job"`
	Religion                       *string               `json:"religion"`
	BloodGroup                     *string               `json:"bloodGroup"`
	Insurence                      *string               `json:"insurence"`
	InsurenceNumber                *string               `json:"insurenceNumber"`
	Phone                          *string               `json:"phone"`
	LastHealthCheckTime            *time.Time            `json:"lastHealthCheckTime"`
}

func (PatientModel) TableName() string {
	return "patients"
}

func (*Module) autoMigrate() {
	if err := pg.GetDB().AutoMigrate(&PatientModel{}); err != nil {
		logger.Panic(err)
	}
}

type patientDB struct {
	*pg.Service[PatientModel]
}

var patientRepo *patientDB

func PatientRepository() *patientDB {
	if patientRepo == nil {
		patientRepo = &patientDB{
			Service: pg.NewService(
				&pg.Service[PatientModel]{
					DB: pg.GetDB(),
				},
			),
		}
	}
	return patientRepo
}
