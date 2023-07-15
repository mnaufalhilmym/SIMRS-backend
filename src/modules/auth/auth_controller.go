package auth

import (
	"errors"
	commonjwt "simrs/src/common/jwt"
	"simrs/src/helpers"
	"simrs/src/libs/hash/argon2"
	"simrs/src/libs/jwx/jwt"
	"simrs/src/libs/parser"
	"simrs/src/modules/account"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func (m *Module) controller() {
	m.App.Post("/api/v1/auth/signin", m.signIn)
	m.App.Get("/api/v1/auth/account", account.UpdateAccountLastActivityTime(), m.account)
	m.App.Post("/api/v1/auth/renew-token", account.UpdateAccountLastActivityTime(), m.renewToken)
}

func (m *Module) signIn(c *fiber.Ctx) error {
	req := new(signInReq)
	if err := parser.ParseReqBody(c, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrBadRequest.Error()),
		})
	}

	accountDetailData, err := m.getAccountDetailByUsernameService(req.Username)
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

	isAuthorized, err := argon2.CompareStringAndEncodedHash(req.Password, accountDetailData.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrInternalServerError.Error()),
		})
	}
	if !isAuthorized {
		return c.Status(fiber.StatusBadRequest).JSON(&response{
			Error: helpers.GetErrorMessage(errors.New("incorrect username or password").Error(), fiber.ErrBadRequest.Error()),
		})
	}

	jwtToken, err := jwt.Create(&commonjwt.JwtPayload{
		ID:   accountDetailData.ID,
		Role: accountDetailData.Role,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrInternalServerError.Error()),
		})
	}

	return c.Status(fiber.StatusOK).JSON(&response{
		Data: &signInRes{
			Token: jwtToken,
			Name:  accountDetailData.Name,
			Role:  accountDetailData.Role,
		},
	})
}

func (m *Module) account(c *fiber.Ctx) error {
	token := new(commonjwt.JwtPayload)
	if err := parser.ParseReqBearerToken(c, token); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrBadRequest.Error()),
		})
	}

	accountDetailData, err := m.getAccountDetailService(token.ID)
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
		Data: &accountRes{
			Name: accountDetailData.Name,
			Role: accountDetailData.Role,
		},
	})
}

func (m *Module) renewToken(c *fiber.Ctx) error {
	token := new(commonjwt.JwtPayload)
	if err := parser.ParseReqBearerToken(c, token); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrBadRequest.Error()),
		})
	}

	accountDetailData, err := m.getAccountDetailService(token.ID)
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

	jwtToken, err := jwt.Create(&commonjwt.JwtPayload{
		ID: accountDetailData.ID,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&response{
			Error: helpers.GetErrorMessage(err.Error(), fiber.ErrInternalServerError.Error()),
		})
	}

	return c.Status(fiber.StatusOK).JSON(&response{
		Data: &renewTokenRes{
			Token: jwtToken,
		},
	})
}
