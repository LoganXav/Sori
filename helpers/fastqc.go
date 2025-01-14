package helpers

import (
	"fmt"
	"os"
	"os/exec"
)

func RunFastQC(inputFilePath, outputDir string) error {
	cmd := exec.Command("fastqc", inputFilePath, "-o", outputDir)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run FastQC: %v", err)
	}

	return nil
}
