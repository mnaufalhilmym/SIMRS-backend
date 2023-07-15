package account

import (
	"simrs/src/libs/db/pg"

	"github.com/google/uuid"
)

type paginationOption struct {
	limit  *int
	lastID *uuid.UUID
}

func (m *Module) countAccount() (*int64, error) {
	return AccountRepository().Count(&pg.CountOption{})
}

func (m *Module) getAccountListService(pagination *paginationOption, search *string) (*[]*AccountModel, int, error) {
	where := []pg.Where{}
	limit := 0

	if search != nil && len(*search) > 0 {
		where = append(where, pg.Where{
			Query: "name ILIKE ? OR username ILIKE ? OR role ILIKE ?",
			Args:  []interface{}{"%" + *search + "%", "%" + *search + "%", "%" + *search + "%"},
		})
	}

	if pagination.limit != nil && *pagination.limit > 0 {
		limit = *pagination.limit
	}

	if pagination.lastID != nil && len(*pagination.lastID) > 0 {
		accountDetailData, err := m.getAccountDetailService(pagination.lastID)
		if err == nil {
			where = append(where, pg.Where{
				Query: "last_activity_time < ?",
				Args:  []interface{}{accountDetailData.LastActivityTime},
			})
		}
	}

	return AccountRepository().FindAll(&pg.FindAllOption{
		Where: &where,
		Limit: &limit,
		Order: &[]interface{}{"last_activity_time desc"},
	})
}

func (m *Module) getAccountDetailService(id *uuid.UUID) (*AccountModel, error) {
	return AccountRepository().FindOne(&pg.FindOneOption{
		Where: &[]pg.Where{
			{
				Query: "id = ?",
				Args:  []interface{}{id},
			},
		},
	})
}

func (m *Module) addAccountService(data *AccountModel) (*AccountModel, error) {
	return AccountRepository().Create(data)
}

func (m *Module) updateAccountService(data *AccountModel) (*AccountModel, error) {
	return AccountRepository().Update(data)
}

func (m *Module) deleteAccountService(id *uuid.UUID) error {
	return AccountRepository().Destroy(&AccountModel{
		Model: &pg.Model{
			ID: id,
		},
	})
}
