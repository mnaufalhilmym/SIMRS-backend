package patientexamination

import (
	"simrs/src/libs/db/pg"
	"simrs/src/modules/patient"

	"github.com/google/uuid"
)

type paginationOption struct {
	limit  *int
	lastID *uuid.UUID
}

type paginationQuery struct {
	count *int
	limit *int
	total *int
}

func (m *Module) getPatientExaminationListService(pagination *paginationOption, patientID *uuid.UUID, search *string) (*[]*PatientExaminationModel, *paginationQuery, error) {
	where := []pg.Where{}
	limit := 0

	if patientID != nil && len(patientID) > 0 {
		where = append(where, pg.Where{
			Query: "patient_id = ?",
			Args:  []interface{}{patientID},
		})
	}

	if search != nil && len(*search) > 0 {
		where = append(where, pg.Where{
			Query: "examination_time ILIKE ? OR examination ILIKE ? OR treatment ILIKE ?",
			Args:  []interface{}{"%" + *search + "%", "%" + *search + "%", "%" + *search + "%"},
		})
	}

	if pagination.limit != nil && *pagination.limit > 0 {
		limit = *pagination.limit
	}

	if pagination.lastID != nil && len(*pagination.lastID) > 0 {
		patientExaminationDetailData, err := m.getPatientExaminationDetailService(pagination.lastID)
		if err == nil {
			where = append(where, pg.Where{
				Query: "examination_time < ?",
				Args:  []interface{}{patientExaminationDetailData.ExaminationTime},
			})
		}
	}

	data, page, err := PatientExaminationRepository().FindAll(&pg.FindAllOption{
		Where: &where,
		Limit: &limit,
		Order: &[]interface{}{"examination_time desc"},
	})

	return data, &paginationQuery{count: &page.Count, limit: &page.Limit, total: &page.Total}, err
}

func (m *Module) getPatientExaminationDetailService(id *uuid.UUID) (*PatientExaminationModel, error) {
	return PatientExaminationRepository().FindOne(&pg.FindOneOption{
		Where: &[]pg.Where{
			{
				Query: "id = ?",
				Args:  []interface{}{id},
			},
		},
	})
}

func (m *Module) addPatientExaminationService(data *PatientExaminationModel) (*PatientExaminationModel, error) {
	return PatientExaminationRepository().Create(data)
}

func (m *Module) updatePatientExaminationService(data *PatientExaminationModel) (*PatientExaminationModel, error) {
	return PatientExaminationRepository().Update(data)
}

func (m *Module) deletePatientExaminationService(id *uuid.UUID) error {
	return PatientExaminationRepository().Destroy(&PatientExaminationModel{
		Model: &pg.Model{
			ID: id,
		},
	})
}

func (m *Module) getPatientDetailService(id *uuid.UUID) (*patient.PatientModel, error) {
	return patient.PatientRepository().FindOne(&pg.FindOneOption{
		Where: &[]pg.Where{
			{
				Query: "id = ?",
				Args:  []interface{}{id},
			},
		},
	})
}

func (m *Module) updatePatientService(data *patient.PatientModel) (*patient.PatientModel, error) {
	return patient.PatientRepository().Update(data)
}
