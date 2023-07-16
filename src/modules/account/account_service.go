package account

import (
	accountrole "simrs/src/common/account_role"
	"simrs/src/libs/db/pg"

	"github.com/google/uuid"
)

type searchOption struct {
	byRole *accountrole.Role
	byAny  *string
}

type paginationOption struct {
	limit  *int
	lastID *uuid.UUID
}

func (m *Module) countAccount(search *searchOption) (*int64, error) {
	where := []pg.Where{}

	if search != nil {
		if search.byRole != nil && len(*search.byRole) > 0 {
			where = append(where, pg.Where{
				Query: "role = ?",
				Args:  []interface{}{search.byRole},
			})
		}
		if search.byAny != nil && len(*search.byAny) > 0 {
			where = append(where, pg.Where{
				Query: "name ILIKE ? OR username ILIKE ? OR role ILIKE ?",
				Args:  []interface{}{"%" + *search.byAny + "%", "%" + *search.byAny + "%", "%" + *search.byAny + "%"},
			})
		}
	}
	return AccountRepository().Count(&pg.CountOption{
		Where: &where,
	})
}

func (m *Module) getAccountListService(pagination *paginationOption, search *searchOption) (*[]*AccountModel, int, error) {
	where := []pg.Where{}
	limit := 0

	if search != nil {
		if search.byRole != nil && len(*search.byRole) > 0 {
			where = append(where, pg.Where{
				Query: "role = ?",
				Args:  []interface{}{search.byRole},
			})
		}
		if search.byAny != nil && len(*search.byAny) > 0 {
			where = append(where, pg.Where{
				Query: "name ILIKE ? OR username ILIKE ? OR role ILIKE ?",
				Args:  []interface{}{"%" + *search.byAny + "%", "%" + *search.byAny + "%", "%" + *search.byAny + "%"},
			})
		}
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
