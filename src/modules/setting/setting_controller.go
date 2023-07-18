package setting

import (
	"errors"
	accountrole "simrs/src/common/account_role"
	"simrs/src/helpers"
	"simrs/src/libs/parser"
	"simrs/src/middlewares/authguard"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func (m *Module) controller() {
	m.App.Get("/api/v1/setting", authguard.AuthGuard(accountrole.ROLE_ADMIN, accountrole.ROLE_SUPERADMIN), m.getSetting)
	m.App.Post("/api/v1/setting", authguard.AuthGuard(accountrole.ROLE_ADMIN, accountrole.ROLE_SUPERADMIN), m.addSetting)
	m.App.Put("/api/v1/setting", authguard.AuthGuard(accountrole.ROLE_ADMIN, accountrole.ROLE_SUPERADMIN), m.updateSetting)
}

func (m *Module) getSetting(c *fiber.Ctx) error {
	settingData, err := m.getSettingService()
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
		Data: settingData,
	})
}

func (m *Module) addSetting(c *fiber.Ctx) error {
	req := new(addSettingReq)
	if err := parser.ParseReqBody(c, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrBadRequest.Error()),
		})
	}

	workersBytes, err := sonic.ConfigFastest.Marshal(req.Workers)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrInternalServerError.Error()),
		})
	}
	workers := datatypes.JSON(workersBytes)

	settingData, err := m.addSettingService(&SettingModel{
		CoverImg: req.CoverImg,
		Workers:  &workers,
		Vision:   req.Vision,
		Mission:  req.Mission,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrInternalServerError.Error()),
		})
	}

	return c.Status(fiber.StatusOK).JSON(&response{
		Data: settingData,
	})
}

func (m *Module) updateSetting(c *fiber.Ctx) error {
	req := new(updateSettingReq)
	if err := parser.ParseReqBody(c, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrBadRequest.Error()),
		})
	}

	settingData := &SettingModel{
		CoverImg: req.CoverImg,
		Vision:   req.Vision,
		Mission:  req.Mission,
	}

	if req.Workers != nil {
		workersBytes, err := sonic.ConfigFastest.Marshal(req.Workers)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&response{
				Error: helpers.GetErrorMessage(err.Error(), fiber.ErrInternalServerError.Error()),
			})
		}
		workers := datatypes.JSON(workersBytes)
		settingData.Workers = &workers
	}

	settingData, err := m.updateSettingService(settingData)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrBadRequest.Error()),
		})
	}

	return c.Status(fiber.StatusOK).JSON(&response{
		Data: settingData,
	})
}
