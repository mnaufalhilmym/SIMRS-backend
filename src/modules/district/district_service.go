package district

import (
	"simrs/src/libs/db/pg"

	"github.com/google/uuid"
)

type paginationOption struct {
	limit  *int
	lastID *uuid.UUID
}

func (m *Module) getDistrictListService(pagination *paginationOption, search *string) (*[]*DistrictModel, int, error) {
	where := []pg.Where{}
	limit := 0

	if search != nil && len(*search) > 0 {
		where = append(where, pg.Where{
			Query: "name ILIKE ?",
			Args:  []interface{}{"%" + *search + "%"},
		})
	}

	if pagination.limit != nil && *pagination.limit > 0 {
		limit = *pagination.limit
	}

	if pagination.lastID != nil && len(*pagination.lastID) > 0 {
		districtDetailData, err := m.getDistrictDetailService(pagination.lastID)
		if err == nil {
			where = append(where, pg.Where{
				Query: "updated_at < ?",
				Args:  []interface{}{districtDetailData.UpdatedAt},
			})
		}
	}

	return DistrictRepository().FindAll(&pg.FindAllOption{
		Where: &where,
		Limit: &limit,
		Order: &[]interface{}{"updated_at desc"},
	})
}

func (m *Module) getDistrictDetailService(id *uuid.UUID) (*DistrictModel, error) {
	return DistrictRepository().FindOne(&pg.FindOneOption{
		Where: &[]pg.Where{
			{
				Query: "id = ?",
				Args:  []interface{}{id},
			},
		},
	})
}

func (m *Module) addDistrictService(data *DistrictModel) (*DistrictModel, error) {
	return DistrictRepository().Create(data)
}

func (m *Module) updateDistrictService(data *DistrictModel) (*DistrictModel, error) {
	return DistrictRepository().Update(data)
}

func (m *Module) deleteDistrictService(id *uuid.UUID) error {
	return DistrictRepository().Destroy(&DistrictModel{
		Model: &pg.Model{
			ID: id,
		},
	})
}
