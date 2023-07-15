package patient

import (
	"sim-puskesmas/src/libs/db/pg"

	"github.com/google/uuid"
)

type searchOption struct {
	byFamilyCardNumber *string
	byDistrictID       *uuid.UUID
	byAny              *string
}

type paginationOption struct {
	limit  *int
	lastID *uuid.UUID
}

func (m *Module) getPatientListService(pagination *paginationOption, search *searchOption) (*[]*PatientModel, int, error) {
	where := []pg.Where{}
	limit := 0

	if search != nil {
		if search.byFamilyCardNumber != nil && len(*search.byFamilyCardNumber) > 0 {
			where = append(where, pg.Where{
				Query: "family_card_number = ?",
				Args:  []interface{}{search.byFamilyCardNumber},
			})
		}
		if search.byDistrictID != nil && len(*search.byDistrictID) > 0 {
			where = append(where, pg.Where{
				Query: "district_id = ?",
				Args:  []interface{}{search.byDistrictID},
			})
		}
		if search.byAny != nil && len(*search.byAny) > 0 {
			where = append(where, pg.Where{
				Query: "medical_record_number ILIKE ? OR family_card_number ILIKE ? OR population_identification_number ILIKE ? OR name ILIKE ?",
				Args:  []interface{}{"%" + *search.byAny + "%", "%" + *search.byAny + "%", "%" + *search.byAny + "%", "%" + *search.byAny + "%"},
			})
		}
	}

	if pagination.limit != nil && *pagination.limit > 0 {
		limit = *pagination.limit
	}

	if pagination.lastID != nil && len(*pagination.lastID) > 0 {
		patientDetailData, err := m.getPatientDetailService(pagination.lastID)
		if err == nil {
			where = append(where, pg.Where{
				Query: "last_health_check_time < ?",
				Args:  []interface{}{patientDetailData.LastHealthCheckTime},
			})
		}
	}

	return PatientRepository().FindAll(&pg.FindAllOption{
		Where: &where,
		Limit: &limit,
		Order: &[]interface{}{"last_health_check_time desc"},
	})
}

func (m *Module) getPatientDetailService(id *uuid.UUID) (*PatientModel, error) {
	return PatientRepository().FindOne(&pg.FindOneOption{
		Where: &[]pg.Where{
			{
				Query: "id = ?",
				Args:  []interface{}{id},
			},
		},
	})
}

func (m *Module) getPatientCountService(districtId *uuid.UUID) (*int64, error) {
	where := []pg.Where{}

	if districtId != nil && len(*districtId) > 0 {
		where = append(where, pg.Where{
			Query: "district_id = ?",
			Args:  []interface{}{districtId},
		})
	}

	return PatientRepository().Count(&pg.CountOption{
		Where: &where,
	})
}

func (m *Module) addPatientService(data *PatientModel) (*PatientModel, error) {
	return PatientRepository().Create(data)
}

func (m *Module) updatePatientService(data *PatientModel) (*PatientModel, error) {
	return PatientRepository().Update(data)
}

func (m *Module) deletePatientService(id *uuid.UUID) error {
	return PatientRepository().Destroy(&PatientModel{
		Model: &pg.Model{
			ID: id,
		},
	})
}
