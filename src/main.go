package main

import (
	"runtime"
	"sim-puskesmas/src/constants"
	"sim-puskesmas/src/helpers"
	"sim-puskesmas/src/libs/env"
	"sim-puskesmas/src/libs/exithandler"
	applogger "sim-puskesmas/src/libs/logger"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var logger = applogger.New("App")

func main() {
	{
		appMode := env.Get(env.APP_MODE)
		if appMode != constants.APP_MODE_RELEASE {
			appName := env.Get(env.APP_NAME)
			logger.Log("starting " + appName + " in " + appMode + " on " + runtime.Version())
		}
	}

	app := fiber.New(fiber.Config{
		AppName:     env.Get(env.APP_NAME),
		Network:     fiber.NetworkTCP,
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
		ReadTimeout: 30 * time.Second,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			statusCode := fiber.StatusInternalServerError
			if fiberError, ok := err.(*fiber.Error); ok {
				statusCode = fiberError.Code
			}
			logger.Error(err)
			return c.Status(statusCode).JSON(&response{
				Error: helpers.GetErrorMessage(err.Error(), "an unexpected error occurred"),
			})
		},
	})

	app.Use(recover.New(recover.Config{
		EnableStackTrace: env.Get(env.APP_MODE) != constants.APP_MODE_RELEASE,
	}))

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins: func() string {
			if env.Get(env.APP_MODE) == constants.APP_MODE_RELEASE {
				return env.Get(env.WEB_ADDRESS)
			}
			return "*"
		}(),
	}))

	if env.Get(env.APP_MODE) != constants.APP_MODE_RELEASE {
		app.Use(fiberlogger.New())
	}

	module := &module{app: app}
	module.Load()

	exithandler.Add(exithandler.FnRunInExit{
		FnDescription: "shutting down app",
		Fn: func() {
			if err := app.Shutdown(); err != nil {
				logger.Error(err)
			}
		},
	})

	if err := app.Listen(env.Get(env.APP_ADDRESS)); err != nil {
		logger.Panic(err)
	}
}
