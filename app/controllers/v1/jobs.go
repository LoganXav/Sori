package controllers

import (
	"LoganXav/sori/app/models"
	repositoriesV1 "LoganXav/sori/app/repositories/v1"
	helpers "LoganXav/sori/helpers"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

func JobsAlignment(c *fiber.Ctx) error {

	// Create a new job
	job, err := repositoriesV1.JobCreate(c)

	if err != nil {
		return helpers.UnprocessableResponse(c, http.StatusInternalServerError, err.Error())
	}

	// Update job status to "running"
	err = repositoriesV1.JobUpdateStatus(job.ID, models.JobPending)
	if err != nil {
		return helpers.UnprocessableResponse(c, http.StatusInternalServerError, err.Error())
	}

	// Define paths
	tempDir := os.TempDir()
	fastQFilePath := filepath.Join(tempDir, job.FileID+".fastq.gz")
	referenceGenomePath := filepath.Join(tempDir, "reference_genome.fasta")
	outputDir := filepath.Join(tempDir, "bwa_output")
	alignmentOutputFile := filepath.Join(outputDir, job.FileID+"_aligned.sam")

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return helpers.UnprocessableResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to create output directory: %v", err).Error())
	}
	defer os.RemoveAll(outputDir)

	// Download FASTQ file from S3
	fastqFileKey := "fastqcs/" + job.FileID + ".fastq.gz"
	if err := helpers.DownloadFromS3(fastqFileKey, fastQFilePath); err != nil {
		return helpers.UnprocessableResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to download FASTQ file from S3: %v", err).Error())
	}
	defer os.Remove(fastQFilePath)

	// Validate FASTQ file exists
	if _, err := os.Stat(fastQFilePath); err != nil {
		return helpers.UnprocessableResponse(c, http.StatusInternalServerError, fmt.Errorf("FASTQ file is missing or inaccessible: %v", err).Error())
	}

	// Download reference genome from S3
	referenceGenomeKey := "reference-genomes/" + job.ReferenceID + ".fasta"
	if err := helpers.DownloadFromS3(referenceGenomeKey, referenceGenomePath); err != nil {
		return helpers.UnprocessableResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to download reference genome from S3: %v", err).Error())
	}
	defer os.Remove(referenceGenomePath)

	// Validate reference genome file exists
	if _, err := os.Stat(referenceGenomePath); err != nil {
		return helpers.UnprocessableResponse(c, http.StatusInternalServerError, fmt.Errorf("reference genome file is missing or inaccessible: %v", err).Error())
	}

	// // Run BWA alignment
	// if err := helpers.RunBWA(fastQFilePath, referenceGenomePath, alignmentOutputFile); err != nil {
	// 	// Update job status to "failed"
	// 	repositoriesV1.JobUpdateStatus(job.ID, models.JobFailed)
	// 	return helpers.UnprocessableResponse(c, http.StatusInternalServerError, fmt.Errorf("BWA alignment failed: %v", err).Error())
	// }

	// Validate the output alignment file exists
	if _, err := os.Stat(alignmentOutputFile); err != nil {
		return helpers.UnprocessableResponse(c, http.StatusInternalServerError, fmt.Errorf("alignment output file is missing: %v", err).Error())
	}

	// Upload the alignment output file to S3
	alignmentFileKey := "alignments/" + job.FileID + "_aligned.sam"
	alignmentResultMap, err := helpers.UploadToS3(alignmentOutputFile, alignmentFileKey)

	if err != nil {
		return helpers.UnprocessableResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to upload alignment file to S3: %v", err).Error())
	}

	alignmentResultURL := alignmentResultMap.AWSUrl

	// if err != nil {
	// 	// Update job status to "failed" in case of error
	// 	repositoriesV1.JobUpdateStatus(job.ID, models.JobFailed)
	// 	return helpers.UnprocessableResponse(c, http.StatusInternalServerError, err.Error())
	// }

	// Update job status to "completed" and save result URL
	updatedJob, errUpdate := repositoriesV1.JobUpdateResult(job.ID, alignmentResultURL, models.JobCompleted)
	if errUpdate != nil {
		return helpers.UnprocessableResponse(c, http.StatusInternalServerError, "failed to update job result")
	}

	result := map[string]interface{}{
		"message": "Quality Control Analysis Completed.",
		"data":    updatedJob,
	}

	return helpers.SuccessResponse(c, result, "success")
}

func JobsQualityControl(c *fiber.Ctx) error {
	// Create a new job
	job, err := repositoriesV1.JobCreate(c)

	if err != nil {
		return helpers.UnprocessableResponse(c, http.StatusInternalServerError, err.Error())
	}

	// Update job status to "running"
	err = repositoriesV1.JobUpdateStatus(job.ID, models.JobPending)
	if err != nil {
		return helpers.UnprocessableResponse(c, http.StatusInternalServerError, err.Error())
	}

	// Define paths
	tempDir := os.TempDir()
	inputFilePath := filepath.Join(tempDir, job.FileID+".fastq.gz")
	outputDir := filepath.Join(tempDir, "fastqc_output")

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return helpers.UnprocessableResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to create output directory: %v", err).Error())
	}
	defer os.RemoveAll(outputDir)

	// Download FASTQ file from S3
	fileKey := "fastqcs/" + job.FileID + ".fastq.gz"

	if err := helpers.DownloadFromS3(fileKey, inputFilePath); err != nil {
		return helpers.UnprocessableResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to download file from S3: %v", err).Error())
	}
	defer os.Remove(inputFilePath)

	// Validate file exists
	if _, err := os.Stat(inputFilePath); err != nil {
		return helpers.UnprocessableResponse(c, http.StatusInternalServerError, fmt.Errorf("downloaded file is missing or inaccessible: %v", err).Error())
	}

	if err := helpers.RunFastQC(inputFilePath, outputDir); err != nil {
		return helpers.UnprocessableResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to save FastQ output to s3: %v", err).Error())
	}

	// Generate the output report URL
	outputFile := filepath.Join(outputDir, job.FileID+"_fastqc.html")

	reportKey := "fastq-reports/" + job.FileID + "_fastqc.html"

	uploadResultMap, errUpload := helpers.UploadToS3(outputFile, reportKey)

	if errUpload != nil {
		return helpers.UnprocessableResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to save FastQ output to s3: %v", errUpload).Error())
	}

	outputURL := uploadResultMap.AWSUrl

	// if err != nil {
	// 	// Update job status to "failed" in case of error
	// 	repositoriesV1.JobUpdateStatus(job.ID, models.JobFailed)
	// 	return helpers.UnprocessableResponse(c, http.StatusInternalServerError, err.Error())
	// }

	// Update job status to "completed" and save result URL
	updatedJob, errUpdate := repositoriesV1.JobUpdateResult(job.ID, outputURL, models.JobCompleted)
	if errUpdate != nil {
		return helpers.UnprocessableResponse(c, http.StatusInternalServerError, "failed to update job result")
	}

	result := map[string]interface{}{
		"message": "Quality Control Analysis Completed.",
		"data":    updatedJob,
	}

	return helpers.SuccessResponse(c, result, "success")
}

func JobsResult(c *fiber.Ctx) error {

	workflow := map[string]interface{}{
		"data": "Workflow Completed",
	}

	return helpers.SuccessResponse(c, workflow, "success")
}

func JobsStatus(c *fiber.Ctx) error {

	workflow := map[string]interface{}{
		"data": "Workflow Completed",
	}

	return helpers.SuccessResponse(c, workflow, "success")
}
