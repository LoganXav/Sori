package helpers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// Handle http success response
//
//	param c *fiber.Ctx
//	param data interface{}
//	param message string
//	return error
func SuccessResponse(c *fiber.Ctx, data interface{}, message string) error {
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"error":   false,
		"data":    data,
		"message": message,
	})
}