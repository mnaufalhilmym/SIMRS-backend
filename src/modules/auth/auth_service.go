package auth

import (
	"sim-puskesmas/src/libs/db/pg"
	"sim-puskesmas/src/modules/account"

	"github.com/google/uuid"
)

func (m *Module) getAccountDetailService(id *uuid.UUID) (*account.AccountModel, error) {
	return account.AccountRepository().FindOne(&pg.FindOneOption{
		Where: &[]pg.Where{
			{
				Query: "id = ?",
				Args:  []interface{}{id},
			},
		},
	})
}

func (m *Module) getAccountDetailByUsernameService(username *string) (*account.AccountModel, error) {
	return account.AccountRepository().FindOne(&pg.FindOneOption{
		Where: &[]pg.Where{
			{
				Query: "username = ?",
				Args:  []interface{}{username},
			},
		},
	})
}
