package patient

import (
	"simrs/src/libs/db/pg"

	"github.com/google/uuid"
)

type searchOption struct {
	byMedicalRecordNumber  *string
	byFamilyCardNumber     *string
	byRelationshipInFamily *relationshipInFamily
	byDistrictID           *uuid.UUID
	byAny                  *string
}

type paginationOption struct {
	limit  *int
	lastID *uuid.UUID
}

func (m *Module) getPatientListService(pagination *paginationOption, search *searchOption) (*[]*PatientModel, int, error) {
	where := []pg.Where{}
	limit := 0

	if search != nil {
		if search.byMedicalRecordNumber != nil && len(*search.byMedicalRecordNumber) > 0 {
			where = append(where, pg.Where{
				Query: "medical_record_number ILIKE ?",
				Args:  []interface{}{*search.byMedicalRecordNumber + "%"},
			})
		}
		if search.byFamilyCardNumber != nil && len(*search.byFamilyCardNumber) > 0 {
			where = append(where, pg.Where{
				Query: "family_card_number = ?",
				Args:  []interface{}{search.byFamilyCardNumber},
			})
		}
		if search.byRelationshipInFamily != nil && len(*search.byRelationshipInFamily) > 0 {
			where = append(where, pg.Where{
				Query: "byRelationship_in_family = ?",
				Args:  []interface{}{search.byRelationshipInFamily},
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
				Query: "name > ?",
				Args:  []interface{}{patientDetailData.Name},
			})
		}
	}

	return PatientRepository().FindAll(&pg.FindAllOption{
		Where: &where,
		Limit: &limit,
		Order: &[]interface{}{"name asc"},
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

func (m *Module) getPatientCountService(search *searchOption) (*int64, error) {
	where := []pg.Where{}

	if search != nil {
		if search.byFamilyCardNumber != nil && len(*search.byFamilyCardNumber) > 0 {
			where = append(where, pg.Where{
				Query: "family_card_number = ?",
				Args:  []interface{}{search.byFamilyCardNumber},
			})
		}
		if search.byRelationshipInFamily != nil && len(*search.byRelationshipInFamily) > 0 {
			where = append(where, pg.Where{
				Query: "byRelationship_in_family = ?",
				Args:  []interface{}{search.byRelationshipInFamily},
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
