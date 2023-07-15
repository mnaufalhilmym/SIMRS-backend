package main

import (
	"sim-puskesmas/src/libs/env"

	"github.com/gofiber/fiber/v2"
)

func (m *module) controller() {
	m.app.Get("/", m.rootController)
}

func (m *module) rootController(c *fiber.Ctx) error {
	return c.JSON(&response{
		Data: env.Get(env.APP_NAME) + " is running",
	})
}
