package routes

import (
	routeV1 "LoganXav/sori/app/routes/api/v1"

	appConfig "LoganXav/sori/configs"

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

func ApiRoutes(app *fiber.App) {
	apiV1 := app.Group("/api/v1", logger.New())

	// Sample of protected route
	// routeV1.IndexProtectedRoute(apiV1)

	routeV1.IndexRoute(apiV1)
	routeV1.JobsRoute(apiV1)

	if appConfig.GetEnv("ENV") == "dev" {
		routeV1.SwaggerRoute(apiV1)
	}
}