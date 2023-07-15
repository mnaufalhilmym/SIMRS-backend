package patient

import (
	applogger "sim-puskesmas/src/libs/logger"

	"github.com/gofiber/fiber/v2"
)

type Module struct {
	App *fiber.App
}

var patientModule *Module
var logger = applogger.New("PatientModule")

func New(module *Module) {
	if patientModule != nil {
		logger.Error("module has been initiated")
		return
	}
	patientModule = module
}

func Load() {
	if patientModule == nil {
		logger.Panic("module has not been initiated")
	}
	patientModule.autoMigrate()
	patientModule.controller()
}
