package parser

import (
	"errors"
	"simrs/src/libs/jwx/jwt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func ParseReqBearerToken[T any](c *fiber.Ctx, tokenData *T) error {
	authorizationHeader := c.Get("Authorization")
	if len(authorizationHeader) == 0 {
		err := errors.New("authorization header not found")
		logger.Error(err)
		return err
	}

	authorization := strings.Split(authorizationHeader, " ")
	if strings.ToLower(authorization[0]) != "bearer" {
		err := errors.New("not a bearer token")
		logger.Error(err)
		return err
	}

	if err := jwt.Parse(authorization[1], tokenData); err != nil {
		return err
	}

	return nil
}
