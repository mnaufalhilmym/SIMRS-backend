package env

import (
	"os"
	applogger "simrs/src/libs/logger"

	_ "github.com/joho/godotenv/autoload"
)

type Env string

const (
	APP_NAME    Env = "APP_NAME"
	APP_MODE    Env = "APP_MODE"
	APP_ADDRESS Env = "APP_ADDRESS"

	WEB_ADDRESS Env = "WEB_ADDRESS"

	POSTGRES_ADDRESS  Env = "POSTGRES_ADDRESS"
	POSTGRES_USER     Env = "POSTGRES_USER"
	POSTGRES_PASSWORD Env = "POSTGRES_PASSWORD"
	POSTGRES_DBNAME   Env = "POSTGRES_DB"

	JWT_DURATION Env = "JWT_DURATION"

	HASH_MEMORY      Env = "HASH_MEMORY"
	HASH_ITERATIONS  Env = "HASH_ITERATIONS"
	HASH_PARALLELISM Env = "HASH_PARALLELISM"
	HASH_SALTLENGTH  Env = "HASH_SALTLENGTH"
	HASH_KEYLENGTH   Env = "HASH_KEYLENGTH"

	INITIAL_ACCOUNT_NAME     Env = "INITIAL_ACCOUNT_NAME"
	INITIAL_ACCOUNT_USERNAME Env = "INITIAL_ACCOUNT_USERNAME"
	INITIAL_ACCOUNT_PASSWORD Env = "INITIAL_ACCOUNT_PASSWORD"
	INITIAL_ACCOUNT_ROLE     Env = "INITIAL_ACCOUNT_ROLE"
)

var logger = applogger.New("ENV")

func Get(env Env) string {
	envVal, isExist := os.LookupEnv(string(env))
	if !isExist {
		logger.Error("unknown env variable: " + env)
	}

	return envVal
}
