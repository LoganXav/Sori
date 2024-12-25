package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func DefaultMiddleware(app *fiber.App){
	app.Use(
		cors.New(),
		logger.New(),

		// Compresses HTTP responses to reduce payload size.
		compress.New(compress.Config{
			Level: compress.LevelBestSpeed,
		}),

		// Handles panics in the application to prevent crashes.
		recover.New(),
		func(c *fiber.Ctx) error {
			// Custom middleware here
			return c.Next()
		},
		

	)
}