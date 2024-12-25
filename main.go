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

	// Database
	dbConnectionError := appDatabase.Connect()

	if dbConnectionError != nil{
		panic("Cannot connect database: " + dbConnectionError.Error())
	}

	db, err := appDatabase.DB.DB()

	if err != nil {
		errc := db.Close()

		if errc != nil {
			panic("Cannot connect database, closing connection")
		}

		panic("Cannot connect database: " + errc.Error())
	}

	// TODO: DB Migration
	defer db.Close()
	// TODO: Redis Setup
	// TODO: S3 Setup
	// TODO: Apply Middlewares
	// TODO: Setup Routes
	// TODO: Run Cron Jobs
	// TODO: Start Server

}