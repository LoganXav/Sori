package v1

import (
	controllerV1 "LoganXav/sori/app/controllers/v1"
	validators "LoganXav/sori/app/validators"

	"github.com/gofiber/fiber/v2"
)

func JobsRoute(router fiber.Router) {
	
	jobs := router.Group("/jobs")

	jobs.Post("/", controllerV1.JobsWorkflow)
	jobs.Post("/qc", validators.JobsQualityControl, controllerV1.JobsQualityControl)
	jobs.Get("/:id/results", controllerV1.JobsResult)
	jobs.Get("/:id", controllerV1.JobsStatus)
}