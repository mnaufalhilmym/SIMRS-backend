package argon2

import (
	applogger "simrs/src/libs/logger"
	"simrs/src/libs/validator"
)

type Config struct {
	Memory      uint32 `validate:"required"`
	Iterations  uint32 `validate:"required"`
	Parallelism uint8  `validate:"required"`
	SaltLength  int    `validate:"required"`
	KeyLength   uint32 `validate:"required"`
	IsDebug     bool
}

var argon2Config *Config
var logger = applogger.New("Argon2")

func New(config *Config) {
	if config.IsDebug {
		logger.Log("initializing argon2")
	}

	if err := validator.Struct(config); err != nil {
		logger.Panic(err)
	}

	argon2Config = config
}
