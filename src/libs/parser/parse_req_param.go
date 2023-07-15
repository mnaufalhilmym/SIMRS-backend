package parser

import (
	"sim-puskesmas/src/libs/validator"

	"github.com/gofiber/fiber/v2"
)

func ParseReqParam[T any](c *fiber.Ctx, param T) error {
	if err := c.ParamsParser(param); err != nil {
		logger.Error(err)
		return err
	}
	if err := validator.Struct(param); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}
