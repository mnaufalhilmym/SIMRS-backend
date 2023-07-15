package patientexamination

import (
	applogger "simrs/src/libs/logger"

	"github.com/gofiber/fiber/v2"
)

type Module struct {
	App *fiber.App
}

var patientExaminationModule *Module
var logger = applogger.New("PatientExaminationModule")

func New(module *Module) {
	if patientExaminationModule != nil {
		logger.Error("module has been initiated")
		return
	}
	patientExaminationModule = module
}

func Load() {
	if patientExaminationModule == nil {
		logger.Panic("module has not been initiated")
	}
	patientExaminationModule.autoMigrate()
	patientExaminationModule.controller()
}
