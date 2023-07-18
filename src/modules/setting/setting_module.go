package setting

import (
	applogger "simrs/src/libs/logger"

	"github.com/gofiber/fiber/v2"
)

type Module struct {
	App *fiber.App
}

var settingModule *Module
var logger = applogger.New("SettingModule")

func New(module *Module) {
	if settingModule != nil {
		logger.Error("module has been initiated")
		return
	}
	settingModule = module
}

func Load() {
	if settingModule == nil {
		logger.Panic("module has not been initiated")
	}
	settingModule.autoMigrate()
	settingModule.createInitialSetting()
	settingModule.controller()
}
