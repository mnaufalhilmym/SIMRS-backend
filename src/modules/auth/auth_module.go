package auth

import (
	applogger "simrs/src/libs/logger"

	"github.com/gofiber/fiber/v2"
)

type Module struct {
	App *fiber.App
}

var authModule *Module
var logger = applogger.New("AuthModule")

func New(module *Module) {
	if authModule != nil {
		logger.Error("module has been initiated")
		return
	}
	authModule = module
}

func Load() {
	if authModule == nil {
		logger.Panic("module has not been initiated")
	}
	authModule.controller()
}
