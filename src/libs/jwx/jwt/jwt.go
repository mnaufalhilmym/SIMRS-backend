package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	applogger "sim-puskesmas/src/libs/logger"
	"sim-puskesmas/src/libs/validator"
	"time"
)

var conf *jwtConf
var logger = applogger.New("JWT")

type jwtConf struct {
	duration   *time.Duration
	privateKey *rsa.PrivateKey
}

type Config struct {
	Bits     int
	Duration *time.Duration `validate:"required"`
	IsDebug  bool
}

func New(config *Config) {
	if err := validator.Struct(config); err != nil {
		logger.Panic(err)
	}

	if config.IsDebug {
		logger.Log("initializing private key for jwt token")
	}

	if config.Bits == 0 {
		config.Bits = 2048
	}

	privKey, err := rsa.GenerateKey(rand.Reader, config.Bits)
	if err != nil {
		logger.Panic(err)
	}

	conf = &jwtConf{
		duration:   config.Duration,
		privateKey: privKey,
	}
}
