package database

import (
	appConfig "LoganXav/sori/configs"
	appHelper "LoganXav/sori/helpers"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() error {
	var err error
	dbConnectionUrl, _ := appHelper.ConnectionUrlBuilder(appConfig.GetEnv("DB_DRIVER"))

	var loggerInfo logger.Interface

	if appConfig.GetEnv("ENV") == "development" {
		loggerInfo = logger.Default.LogMode(logger.Info)
	}

	switch appConfig.GetEnv("DB_DRIVER") {
	case "mysql":
		DB, err = gorm.Open(mysql.Open(dbConnectionUrl), &gorm.Config{
			SkipDefaultTransaction:                   true,
			PrepareStmt:                              true,
			DisableForeignKeyConstraintWhenMigrating: true,
			Logger:                                   loggerInfo,
		})
	case "postgres":
		DB, err = gorm.Open(postgres.Open(dbConnectionUrl), &gorm.Config{
			SkipDefaultTransaction: true,
			PrepareStmt:            true,
			Logger:                 loggerInfo,
		})
	}

	if err != nil {
		panic(err)
	}

	return nil

}