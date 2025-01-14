package helpers

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func RunBWA(fastQFilePath, referenceGenomeFilePath, outputDirPath string) error {
	cmd := exec.Command("bwa", "mem", referenceGenomeFilePath, fastQFilePath)

	outputDir, err := os.Create(outputDirPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer outputDir.Close()

	// Redirect stdout to the output file
	cmd.Stdout = outputDir

	// Capture stderr
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Run the command
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("BWA failed: %v, stderr: %s", err, stderr.String())
	}

	return nil
}
