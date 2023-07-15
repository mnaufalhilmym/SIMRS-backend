package account

import (
	"errors"
	"simrs/src/common/jwt"
	"simrs/src/libs/db/pg"
	"time"

	"github.com/gofiber/fiber/v2"
)

func UpdateAccountLastActivityTime() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.UserContext().Value(jwt.JwtCtxKey("jwtToken")).(*jwt.JwtPayload)
		if !ok {
			err := errors.New("invalid token")
			logger.Log(err)
			return c.Next()
		}

		lastActivityTime := time.Now()
		accountDetailData := &AccountModel{
			Model: &pg.Model{
				ID: token.ID,
			},
			LastActivityTime: &lastActivityTime,
		}

		_, err := AccountRepository().Update(accountDetailData)
		if err != nil {
			logger.Log(err)
			return c.Next()
		}

		return c.Next()
	}
}
