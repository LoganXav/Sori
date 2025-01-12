package helpers

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func RunFastQC(fileID string) (string, error) {
	// Define paths
	tempDir := os.TempDir()
	inputFilePath := filepath.Join(tempDir, fileID+".fastq.gz")
	outputDir := filepath.Join(tempDir, "fastqc_output")

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %v", err)
	}
	defer os.RemoveAll(outputDir) 

	// Download FASTQ file from S3
	fileKey := "fastqcs/" + fileID + ".fastq.gz"

	if err := DownloadFromS3(fileKey, inputFilePath); err != nil {
		return "", fmt.Errorf("failed to download file from S3: %v", err)
	}
	defer os.Remove(inputFilePath) 

	// Validate file exists
	if _, err := os.Stat(inputFilePath); err != nil {
		return "", fmt.Errorf("downloaded file is missing or inaccessible: %v", err)
	}

	// Run FastQC
	cmd := exec.Command("fastqc", inputFilePath, "-o", outputDir)
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to run FastQC: %v", err)
	}

	// Generate the output report URL (e.g., store in S3 or local server)
	outputFile := filepath.Join(outputDir, fileID+"_fastqc.html")
	resultURL := fmt.Sprintf("https://your-storage-service.com/reports/%s", filepath.Base(outputFile))

	return resultURL, nil
}
