package main

import (
	"simrs/src/constants"
	"simrs/src/libs/db/pg"
	"simrs/src/libs/env"
	"simrs/src/libs/exithandler"
	"simrs/src/libs/hash/argon2"
	"simrs/src/libs/jwx/jwt"
	"simrs/src/modules/account"
	"simrs/src/modules/auth"
	"simrs/src/modules/district"
	"simrs/src/modules/patient"
	patientexamination "simrs/src/modules/patient_examination"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type module struct {
	app *fiber.App
}

func (*module) loadRootModule() {
	// PostgreSQL database
	pg.New(&pg.Config{
		Address:  env.Get(env.POSTGRES_ADDRESS),
		User:     env.Get(env.POSTGRES_USER),
		Password: env.Get(env.POSTGRES_PASSWORD),
		DbName:   env.Get(env.POSTGRES_DBNAME),
		IsDebug:  env.Get(env.APP_MODE) != constants.APP_MODE_RELEASE,
	})

	// JWT
	jwt.New(&jwt.Config{
		Duration: func() *time.Duration {
			duration, err := time.ParseDuration(env.Get(env.JWT_DURATION))
			if err != nil {
				logger.Panic(err)
			}
			return &duration
		}(),
		IsDebug: env.Get(env.APP_MODE) != constants.APP_MODE_RELEASE,
	})

	// Argon2
	argon2.New(&argon2.Config{
		Memory: func() uint32 {
			hashMemory, _ := strconv.ParseUint(env.Get(env.HASH_MEMORY), 10, 32)
			return uint32(hashMemory)
		}(),
		Iterations: func() uint32 {
			hashIterations, _ := strconv.ParseUint(env.Get(env.HASH_ITERATIONS), 10, 32)
			return uint32(hashIterations)
		}(),
		Parallelism: func() uint8 {
			hashParallelism, _ := strconv.ParseUint(env.Get(env.HASH_PARALLELISM), 10, 8)
			return uint8(hashParallelism)
		}(),
		SaltLength: func() int {
			hashSaltLength, _ := strconv.Atoi(env.Get(env.HASH_SALTLENGTH))
			return hashSaltLength
		}(),
		KeyLength: func() uint32 {
			hashKeyLength, _ := strconv.ParseUint(env.Get(env.HASH_KEYLENGTH), 10, 32)
			return uint32(hashKeyLength)
		}(),
		IsDebug: env.Get(env.APP_MODE) != constants.APP_MODE_RELEASE,
	})

	// ExitHandler
	exithandler.New(&exithandler.Config{
		IsDebug: env.Get(env.APP_MODE) != constants.APP_MODE_RELEASE,
	})
}

func (m *module) Load() {
	m.loadRootModule()

	m.controller()

	account.New(&account.Module{App: m.app})
	account.Load()

	auth.New(&auth.Module{App: m.app})
	auth.Load()

	district.New(&district.Module{App: m.app})
	district.Load()

	patient.New(&patient.Module{App: m.app})
	patient.Load()

	patientexamination.New(&patientexamination.Module{App: m.app})
	patientexamination.Load()
}
