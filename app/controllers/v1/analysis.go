package controllers

import (
	"LoganXav/sori/helpers"

	"github.com/gofiber/fiber/v2"
)


func AnalysisWorkflow(c *fiber.Ctx) error {

	workflow := map[string] interface{} {
		"data" : "Workflow Completed",
	}
	
	return helpers.SuccessResponse(c, workflow, "success")
}