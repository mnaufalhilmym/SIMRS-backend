package parser

import (
	"sim-puskesmas/src/libs/validator"

	"github.com/gofiber/fiber/v2"
)

func ParseReqBody[T any](c *fiber.Ctx, req T) error {
	if err := c.BodyParser(req); err != nil {
		logger.Error(err)
		return err
	}
	if err := validator.Struct(req); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}
