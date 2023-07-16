package account

import (
	accountrole "simrs/src/common/account_role"
	"simrs/src/libs/db/pg"
	"simrs/src/libs/env"
	"simrs/src/libs/hash/argon2"
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
	role := accountrole.ROLE_SUPERADMIN
	count, err := m.countAccount(&searchOption{byRole: &role})
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
