package v1

import (
	controllerV1 "LoganXav/sori/app/controllers/v1"

	"github.com/gofiber/fiber/v2"
)

func AnalysisRoute(router fiber.Router) {
	
	analysis := router.Group("/analysis")

	// analysis.Post("/", validators.CreateAnalysis, controllerV1.AnalysisWorkflow)
	analysis.Post("/", controllerV1.AnalysisWorkflow)
}