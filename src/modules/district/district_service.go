package district

import (
	"simrs/src/libs/db/pg"

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

func (m *Module) getDistrictListService(pagination *paginationOption, search *string) (*[]*DistrictModel, *paginationQuery, error) {
	where := []pg.FindAllWhere{}
	limit := 0

	if search != nil && len(*search) > 0 {
		where = append(where, pg.FindAllWhere{
			Query:          "name ILIKE ?",
			Args:           []interface{}{"%" + *search + "%"},
			IncludeInCount: true,
		})
	}

	if pagination.limit != nil && *pagination.limit > 0 {
		limit = *pagination.limit
	}

	if pagination.lastID != nil && len(*pagination.lastID) > 0 {
		districtDetailData, err := m.getDistrictDetailService(pagination.lastID)
		if err == nil {
			where = append(where, pg.FindAllWhere{
				Query:          "updated_at < ?",
				Args:           []interface{}{districtDetailData.UpdatedAt},
				IncludeInCount: false,
			})
		}
	}

	data, page, err := DistrictRepository().FindAll(&pg.FindAllOption{
		Where: &where,
		Limit: &limit,
		Order: &[]interface{}{"updated_at desc"},
	})

	return data, &paginationQuery{count: &page.Count, limit: &page.Limit, total: &page.Total}, err
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
