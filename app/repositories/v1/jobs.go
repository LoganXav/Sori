package repositories

import (
	"LoganXav/sori/app/models"
	"LoganXav/sori/app/structs"
	appDatabase "LoganXav/sori/database"
	"LoganXav/sori/helpers"
	"encoding/json"
	"fmt"

	// structs "LoganXav/sori/app/structs"

	"github.com/gofiber/fiber/v2"
)

func JobCreate(c *fiber.Ctx) (models.Job, error) {
	db := appDatabase.DB

	var job models.Job
	var jobsCreateStructure structs.JobsCreate

	// Parse JSON body into the struct
	body := c.Body()
	if err := json.Unmarshal(body, &jobsCreateStructure); err != nil {
		return job, fmt.Errorf("invalid payload: %v", err)
	}

	fileID := helpers.SanitiseText(jobsCreateStructure.FileID)
	jobType := models.JobType(helpers.SanitiseText(jobsCreateStructure.Type))
	// status := models.JobStatus(helpers.SanitiseText(jobsCreateStructure.Status))

	// Validate jobType
	if jobType != models.JobTypeQC && jobType != models.JobTypeAlignment && jobType != models.JobTypeDownstream {
		return job, fmt.Errorf("invalid job type: %v", jobType)
	}


	// Set default status
	status := models.JobPending

	
	// Check if the Job already exists (based on FileID)
	exists := db.Model(&models.Job{}).Where("file_id = ?", fileID).First(&job).RowsAffected > 0
	if exists {
		return job, fmt.Errorf("job with file_id '%s' already exists", fileID)
	}


	// Map the struct to the Job model
	newJob := models.Job{
		FileID: fileID,
		Type:   jobType,
		Status: status,
	}

	// Save the Job to the database
	if err := db.Create(&newJob).Error; err != nil {
		return job, fmt.Errorf("failed to create job: %v", err)
	}

	return newJob, nil

}

func JobUpdateStatus(jobId uint, jobStatus models.JobStatus)  error {
	db := appDatabase.DB
	return db.Model(&models.Job{}).Where("id = ?", jobId).Update("status", jobStatus).Error
}

func JobUpdateResult(jobID uint, resultURL string, status models.JobStatus) (models.Job, error) {
	var job models.Job
	db := appDatabase.DB

	if err := db.Model(&job).Where("id = ?", jobID).Updates(map[string]interface{}{
		"result_url": resultURL,
		"status":     status,
	}).Error; err != nil {
		return job, fmt.Errorf("failed to update job: %v", err)
	}

	if err := db.Where("id = ?", jobID).First(&job).Error; err != nil {
		return job, fmt.Errorf("failed to retrieve updated job: %v", err)
	}

	return job, nil
}


