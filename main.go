package main

import (
	"github.com/gofiber/fiber/v2"

	appConfig "LoganXav/sori/configs"
	appDatabase "LoganXav/sori/database"
)

// @title Sori
// @version 1.0
// @Description This is an genomics API using Golang
// @contact.name Sogbesan Segun
// @contact.url https://github.com/LoganXav/
// @contact.email sogbesansegun22@gmail.com
// @BasePath /
// @schemas http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {

  config := appConfig.FiberConfig()

	// Create a new Fiber app
	app := fiber.New(config)

	// GORM connect
	dbConnectionError := appDatabase.Connect()

	if dbConnectionError != nil{
		panic("Cannot connect database: " + dbConnectionError.Error())
	}

	// Underlying SQL connect
	db, err := appDatabase.DB.DB()

	if err != nil {
		errc := db.Close()

		if errc != nil {
			panic("Cannot connect database, closing connection")
		}

		panic("Cannot connect database: " + errc.Error())
	}

	// Generate Migrations
	errMigrate := appDatabase.MigrateDatabase()
	if errMigrate != nil {
		panic("migration error: " + errMigrate.Error())
	}

	defer db.Close()
	
	// Redis Setup
	if appConfig.GetEnv("REDIS_ACTIVATE") == "true"{
		errRedis := appDatabase.RedisConnect()

		if errRedis != nil {
			panic("Cannot start redis connection: " + errRedis.Error())
		}
	}
	// TODO: S3 Setup
	// TODO: Apply Middlewares
	// TODO: Setup Routes
	// TODO: Run Cron Jobs
	// TODO: Start Server

}