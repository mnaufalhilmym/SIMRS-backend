package account

import (
	"errors"
	accountrole "simrs/src/common/account_role"
	"simrs/src/helpers"
	"simrs/src/libs/db/pg"
	"simrs/src/libs/hash/argon2"
	"simrs/src/libs/parser"
	"simrs/src/middlewares/authguard"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func (m *Module) controller() {
	m.App.Get("/api/v1/accounts", authguard.AuthGuard(accountrole.ROLE_SUPERADMIN), UpdateAccountLastActivityTime(), m.getAccountList)
	m.App.Get("/api/v1/account/:id", authguard.AuthGuard(accountrole.ROLE_SUPERADMIN), UpdateAccountLastActivityTime(), m.getAccountDetail)
	m.App.Post("/api/v1/account", authguard.AuthGuard(accountrole.ROLE_SUPERADMIN), UpdateAccountLastActivityTime(), m.addAccount)
	m.App.Patch("/api/v1/account/:id", authguard.AuthGuard(accountrole.ROLE_SUPERADMIN), UpdateAccountLastActivityTime(), m.updateAccount)
	m.App.Delete("/api/v1/account/:id", authguard.AuthGuard(accountrole.ROLE_SUPERADMIN), UpdateAccountLastActivityTime(), m.deleteAccount)
}

func (m *Module) getAccountList(c *fiber.Ctx) error {
	query := new(getAccountListReqQuery)
	if err := parser.ParseReqQuery(c, query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrBadRequest.Error()),
		})
	}

	accountListData, page, err := m.getAccountListService(&paginationOption{
		limit:  query.Limit,
		lastID: query.LastID,
	}, &searchOption{
		byRole: query.SearchByAccountRole,
		byAny:  query.Search,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(&response{
				Error: helpers.GetErrorMessage(err.Error(), fiber.ErrNotFound.Error()),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrInternalServerError.Error()),
		})
	}

	return c.Status(fiber.StatusOK).JSON(&response{
		Pagination: &pagination{
			Count: page.count,
			Limit: page.limit,
			Total: page.total,
		},
		Data: accountListData,
	})
}

func (m *Module) getAccountDetail(c *fiber.Ctx) error {
	param := new(getAccountDetailReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrBadRequest.Error()),
		})
	}

	accountDetailData, err := m.getAccountDetailService(param.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(&response{
				Error: helpers.GetErrorMessage(err.Error(), fiber.ErrNotFound.Error()),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrInternalServerError.Error()),
		})
	}

	return c.Status(fiber.StatusOK).JSON(&response{
		Data: accountDetailData,
	})
}

func (m *Module) addAccount(c *fiber.Ctx) error {
	req := new(addAccountReq)
	if err := parser.ParseReqBody(c, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrBadRequest.Error()),
		})
	}

	encodedHash, err := argon2.GetEncodedHash(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrInternalServerError.Error()),
		})
	}

	initialActivityTime := time.Unix(0, 0)
	accountDetailData, err := m.addAccountService(&AccountModel{
		Name:             req.Name,
		Username:         req.Username,
		Password:         encodedHash,
		Role:             req.Role,
		LastActivityTime: &initialActivityTime,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrInternalServerError.Error()),
		})
	}

	return c.Status(fiber.StatusOK).JSON(&response{
		Data: accountDetailData,
	})
}

func (m *Module) updateAccount(c *fiber.Ctx) error {
	param := new(updateAccountReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrBadRequest.Error()),
		})
	}

	req := new(updateAccountReq)
	if err := parser.ParseReqBody(c, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrBadRequest.Error()),
		})
	}

	accountDetailData, err := m.getAccountDetailService(param.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(&response{
				Error: helpers.GetErrorMessage(err.Error(), fiber.ErrNotFound.Error()),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrInternalServerError.Error()),
		})
	}

	if *accountDetailData.Role == accountrole.ROLE_SUPERADMIN && *req.Role == accountrole.ROLE_ADMIN {
		role := accountrole.ROLE_SUPERADMIN
		accountCount, err := m.countAccount(&searchOption{
			byRole: &role,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&response{
				Error: helpers.GetErrorMessage(err.Error(), fiber.ErrInternalServerError.Error()),
			})
		}

		if *accountCount <= 1 {
			err := errors.New("unable to change the role of the last superadmin account")
			return c.Status(fiber.StatusBadRequest).JSON(&response{
				Error: helpers.GetErrorMessage(err.Error(), fiber.ErrBadRequest.Error()),
			})
		}
	}

	accountDetailData = &AccountModel{
		Model: &pg.Model{
			ID: param.ID,
		},
		Name:     req.Name,
		Username: req.Username,
		Role:     req.Role,
	}

	if req.Password != nil && len(*req.Password) > 0 {
		encodedHash, err := argon2.GetEncodedHash(req.Password)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&response{
				Error: helpers.GetErrorMessage(err.Error(), fiber.ErrInternalServerError.Error()),
			})
		}
		accountDetailData.Password = encodedHash
	}

	accountDetailData, err = m.updateAccountService(accountDetailData)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(&response{
				Error: helpers.GetErrorMessage(err.Error(), fiber.ErrNotFound.Error()),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrInternalServerError.Error()),
		})
	}

	return c.Status(fiber.StatusOK).JSON(&response{
		Data: accountDetailData,
	})
}

func (m *Module) deleteAccount(c *fiber.Ctx) error {
	param := new(deleteAccountReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrBadRequest.Error()),
		})
	}

	accountDetailData, err := m.getAccountDetailService(param.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(&response{
				Error: helpers.GetErrorMessage(err.Error(), fiber.ErrNotFound.Error()),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrInternalServerError.Error()),
		})
	}

	if *accountDetailData.Role == accountrole.ROLE_SUPERADMIN {
		role := accountrole.ROLE_SUPERADMIN
		accountCount, err := m.countAccount(&searchOption{
			byRole: &role,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&response{
				Error: helpers.GetErrorMessage(err.Error(), fiber.ErrInternalServerError.Error()),
			})
		}

		if *accountCount <= 1 {
			err := errors.New("unable to delete the last superadmin account")
			return c.Status(fiber.StatusBadRequest).JSON(&response{
				Error: helpers.GetErrorMessage(err.Error(), fiber.ErrBadRequest.Error()),
			})
		}
	}

	if err := m.deleteAccountService(param.ID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(&response{
				Error: helpers.GetErrorMessage(err.Error(), fiber.ErrNotFound.Error()),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrInternalServerError.Error()),
		})
	}

	return c.Status(fiber.StatusOK).JSON(&response{
		Data: param.ID,
	})
}
