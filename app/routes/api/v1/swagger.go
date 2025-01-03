package v1

import (
	"github.com/gofiber/fiber/v2"
	swagger "github.com/gofiber/swagger"
)

func SwaggerRoute(router fiber.Router) {
	route := router.Group("/")

	route.Get("documentation/*", swagger.HandlerDefault)
}