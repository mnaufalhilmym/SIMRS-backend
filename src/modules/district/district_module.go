package district

import (
	applogger "simrs/src/libs/logger"

	"github.com/gofiber/fiber/v2"
)

type Module struct {
	App *fiber.App
}

var districtModule *Module
var logger = applogger.New("DistrictModule")

func New(module *Module) {
	if districtModule != nil {
		logger.Error("module has been initiated")
		return
	}
	districtModule = module
}

func Load() {
	if districtModule == nil {
		logger.Panic("module has not been initiated")
	}
	districtModule.autoMigrate()
	districtModule.controller()
}
