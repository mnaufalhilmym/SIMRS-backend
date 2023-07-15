package patientexamination

import (
	"simrs/src/libs/db/pg"
	"time"

	"github.com/google/uuid"
)

type PatientExaminationModel struct {
	*pg.Model
	PatientID       *uuid.UUID `json:"patientId"`
	ExaminationTime *time.Time `json:"examinationTime"`
	Examination     *string    `json:"examination"`
	Treatment       *string    `json:"treatment"`
}

func (PatientExaminationModel) TableName() string {
	return "patient_examinations"
}

func (*Module) autoMigrate() {
	if err := pg.GetDB().AutoMigrate(&PatientExaminationModel{}); err != nil {
		logger.Panic(err)
	}
}

type patientExaminationDB struct {
	*pg.Service[PatientExaminationModel]
}

var patientExaminationRepo *patientExaminationDB

func PatientExaminationRepository() *patientExaminationDB {
	if patientExaminationRepo == nil {
		patientExaminationRepo = &patientExaminationDB{
			Service: pg.NewService(
				&pg.Service[PatientExaminationModel]{
					DB: pg.GetDB(),
				},
			),
		}
	}
	return patientExaminationRepo
}
