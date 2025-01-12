package controllers

import (
	"LoganXav/sori/app/models"
	repositoriesV1 "LoganXav/sori/app/repositories/v1"
	helpers "LoganXav/sori/helpers"
	"net/http"

	"github.com/gofiber/fiber/v2"
)


func JobsWorkflow(c *fiber.Ctx) error {
	workflow := map[string] interface{} {
		"message" : "Workflow Completed",
	}
	
	return helpers.SuccessResponse(c, workflow, "success")
}

func JobsQualityControl(c *fiber.Ctx) error {
	// Create a new job
	job, err := repositoriesV1.JobCreate(c); 
	
	if err != nil {
		return helpers.UnprocessableResponse(c, http.StatusInternalServerError, err.Error())
	}
	
	// Update job status to "running"
	err = repositoriesV1.JobUpdateStatus(job.ID, models.JobPending)
	if err != nil {
		return helpers.UnprocessableResponse(c, http.StatusInternalServerError, err.Error())
	}

	outputURL, err := helpers.RunFastQC(job.FileID)

	if err != nil {
		// Update job status to "failed" in case of error
		repositoriesV1.JobUpdateStatus(1, models.JobFailed)
		return helpers.UnprocessableResponse(c, http.StatusInternalServerError, err.Error())
	}

	// Update job status to "completed" and save result URL
	err = repositoriesV1.JobUpdateResult(job.ID, outputURL, models.JobCompleted)
	if err != nil {
		return helpers.UnprocessableResponse(c, http.StatusInternalServerError, "Failed to update job result")
	}


	result := map[string]interface{}{
		"message": "Quality Control Completed",
		"data":  job,
	}

	return helpers.SuccessResponse(c, result, "success")
}

func JobsResult(c *fiber.Ctx) error {

	workflow := map[string] interface{} {
		"data" : "Workflow Completed",
	}
	
	return helpers.SuccessResponse(c, workflow, "success")
}

func JobsStatus(c *fiber.Ctx) error {

	workflow := map[string] interface{} {
		"data" : "Workflow Completed",
	}
	
	return helpers.SuccessResponse(c, workflow, "success")
}