package pg

import (
	"simrs/src/libs/exithandler"
	applogger "simrs/src/libs/logger"
	"simrs/src/libs/validator"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var db *gorm.DB
var logger = applogger.New("DB")

type Config struct {
	Address  string `validate:"required"`
	User     string `validate:"required"`
	Password string `validate:"required"`
	DbName   string `validate:"required"`
	IsDebug  bool
}

func New(config *Config) {
	if db != nil {
		logger.Error("database has been initiated")
		return
	}
	if err := validator.Struct(config); err != nil {
		logger.Panic(err)
	}

	if config.IsDebug {
		logger.Log("initializing database")
	}

	_db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: "postgresql://" + config.User + ":" + config.Password + "@" + config.Address + "/" + config.DbName,
	}), &gorm.Config{
		Logger: func() gormlogger.Interface {
			if config.IsDebug {
				return gormlogger.Default.LogMode(gormlogger.Info)
			}
			return nil
		}(),
		PrepareStmt: true,
	})
	if err != nil {
		logger.Panic(err)
	}

	db = _db

	exithandler.Add(exithandler.FnRunInExit{
		FnDescription: "closing database",
		Fn: func() {
			db, err := db.DB()
			if err != nil {
				logger.Error(err)
				return
			}
			if err := db.Close(); err != nil {
				logger.Error(err)
			}
		},
	})
}

func GetDB() *gorm.DB {
	if db == nil {
		logger.Panic("database has not been initiated")
	}
	return db
}
