package district

import (
	"errors"
	accountrole "sim-puskesmas/src/common/account_role"
	"sim-puskesmas/src/helpers"
	"sim-puskesmas/src/libs/db/pg"
	"sim-puskesmas/src/libs/parser"
	"sim-puskesmas/src/middlewares/authguard"
	"sim-puskesmas/src/modules/account"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func (m *Module) controller() {
	m.App.Get("/api/v1/districts", authguard.AuthGuard(accountrole.ROLE_ADMIN, accountrole.ROLE_SUPERADMIN), account.UpdateAccountLastActivityTime(), m.getDistrictList)
	m.App.Get("/api/v1/district/:id", authguard.AuthGuard(accountrole.ROLE_ADMIN, accountrole.ROLE_SUPERADMIN), account.UpdateAccountLastActivityTime(), m.getDistrictDetail)
	m.App.Post("/api/v1/district", authguard.AuthGuard(accountrole.ROLE_ADMIN, accountrole.ROLE_SUPERADMIN), account.UpdateAccountLastActivityTime(), m.addDistrict)
	m.App.Patch("/api/v1/district/:id", authguard.AuthGuard(accountrole.ROLE_ADMIN, accountrole.ROLE_SUPERADMIN), account.UpdateAccountLastActivityTime(), m.updateDistrict)
	m.App.Delete("/api/v1/district/:id", authguard.AuthGuard(accountrole.ROLE_ADMIN, accountrole.ROLE_SUPERADMIN), account.UpdateAccountLastActivityTime(), m.deleteDistrict)
}

func (m *Module) getDistrictList(c *fiber.Ctx) error {
	query := new(getDistrictListReqQuery)
	if err := parser.ParseReqQuery(c, query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrBadRequest.Error()),
		})
	}

	districtListData, total, err := m.getDistrictListService(&paginationOption{
		limit:  query.Limit,
		lastID: query.LastID,
	}, query.Search)
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
			Total: &total,
		},
		Data: districtListData,
	})
}

func (m *Module) getDistrictDetail(c *fiber.Ctx) error {
	param := new(getDistrictDetailReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrBadRequest.Error()),
		})
	}

	districtDetailData, err := m.getDistrictDetailService(param.ID)
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
		Data: districtDetailData,
	})
}

func (m *Module) addDistrict(c *fiber.Ctx) error {
	req := new(addDistrictReq)
	if err := parser.ParseReqBody(c, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrBadRequest.Error()),
		})
	}

	districtDetailData, err := m.addDistrictService(&DistrictModel{
		Name: req.Name,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrInternalServerError.Error()),
		})
	}

	return c.Status(fiber.StatusOK).JSON(&response{
		Data: districtDetailData,
	})
}

func (m *Module) updateDistrict(c *fiber.Ctx) error {
	param := new(updateDistrictReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrBadRequest.Error()),
		})
	}

	req := new(updateDistrictReq)
	if err := parser.ParseReqBody(c, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrBadRequest.Error()),
		})
	}

	districtDetailData, err := m.updateDistrictService(&DistrictModel{
		Model: &pg.Model{
			ID: param.ID,
		},
		Name: req.Name,
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
		Data: districtDetailData,
	})
}

func (m *Module) deleteDistrict(c *fiber.Ctx) error {
	param := new(deleteDistrictReqParam)
	if err := parser.ParseReqParam(c, param); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrBadRequest.Error()),
		})
	}

	if err := m.deleteDistrictService(param.ID); err != nil {
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
