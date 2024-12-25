package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func MainRoutes(app *fiber.App) {
	mainRoute := app.Group("/", logger.New())
	mainRoute.Get("/", func(c *fiber.Ctx) error {
		err := c.SendString("Hello Welcome to Sori API!")
		return err
	})
}