package account

import (
	applogger "simrs/src/libs/logger"

	"github.com/gofiber/fiber/v2"
)

type Module struct {
	App *fiber.App
}

var accountModule *Module
var logger = applogger.New("AccountModule")

func New(module *Module) {
	if accountModule != nil {
		logger.Error("module has been initiated")
		return
	}
	accountModule = module
}

func Load() {
	if accountModule == nil {
		logger.Panic("module has not been initiated")
	}
	accountModule.autoMigrate()
	accountModule.createInitialAccount()
	accountModule.controller()
}
