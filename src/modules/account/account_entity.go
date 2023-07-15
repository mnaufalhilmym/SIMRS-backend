package account

import (
	accountrole "sim-puskesmas/src/common/account_role"
	"sim-puskesmas/src/libs/db/pg"
	"sim-puskesmas/src/libs/env"
	"sim-puskesmas/src/libs/hash/argon2"
	"time"
)

type AccountModel struct {
	*pg.Model
	Name             *string           `json:"name"`
	Username         *string           `gorm:"unique" json:"username"`
	Password         *string           `json:"-"`
	Role             *accountrole.Role `json:"role"`
	LastActivityTime *time.Time        `json:"lastActivityTime"`
}

func (AccountModel) TableName() string {
	return "accounts"
}

func (*Module) autoMigrate() {
	if err := pg.GetDB().AutoMigrate(&AccountModel{}); err != nil {
		logger.Panic(err)
	}
}

func (m *Module) createInitialAccount() {
	count, err := m.countAccount()
	if err != nil {
		logger.Panic(err)
	}
	if *count > 0 {
		return
	}

	name := env.Get(env.INITIAL_ACCOUNT_NAME)
	username := env.Get(env.INITIAL_ACCOUNT_USERNAME)
	password := env.Get(env.INITIAL_ACCOUNT_PASSWORD)
	encodedHash, err := argon2.GetEncodedHash(&password)
	if err != nil {
		logger.Panic(err)
	}
	role := accountrole.ROLE_SUPERADMIN

	if _, err := m.addAccountService(&AccountModel{
		Name:     &name,
		Username: &username,
		Password: encodedHash,
		Role:     &role,
	}); err != nil {
		logger.Panic(err)
	}
}

type accountDB struct {
	*pg.Service[AccountModel]
}

var accountRepo *accountDB

func AccountRepository() *accountDB {
	if accountRepo == nil {
		accountRepo = &accountDB{
			Service: pg.NewService(
				&pg.Service[AccountModel]{
					DB: pg.GetDB(),
				},
			),
		}
	}
	return accountRepo
}
