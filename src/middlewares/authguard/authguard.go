package authguard

import (
	"context"
	accountrole "sim-puskesmas/src/common/account_role"
	"sim-puskesmas/src/common/jwt"
	"sim-puskesmas/src/helpers"
	"sim-puskesmas/src/libs/parser"

	"github.com/gofiber/fiber/v2"
)

func AuthGuard(accessRole ...accountrole.Role) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if len(accessRole) == 0 {
			return c.Next()
		}

		token := new(jwt.JwtPayload)
		if err := parser.ParseReqBearerToken(c, token); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(&response{
				Error: helpers.GetErrorMessage(err.Error(), fiber.ErrUnauthorized.Error()),
			})
		}

		isAuthorized := false
		for index := range accessRole {
			if accessRole[index] == *token.Role {
				isAuthorized = true
			}
		}

		if !isAuthorized {
			return c.Status(fiber.StatusUnauthorized).JSON(&response{
				Error: helpers.GetErrorMessage(fiber.ErrUnauthorized.Error()),
			})
		}

		userCtx := context.WithValue(context.Background(), jwt.JwtCtxKey("jwtToken"), token)
		c.SetUserContext(userCtx)

		return c.Next()
	}
}
