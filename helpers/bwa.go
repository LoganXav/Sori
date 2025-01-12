package helpers

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func RunBWA(fastaFile, referenceGenome, outputFile string) error {
	// Prepare the BWA command
	cmd := exec.Command("bwa", "mem", referenceGenome, fastaFile)

	// Capture stdout and stderr
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Run the command
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("bWA failed: %v, stderr: %s", err, stderr.String())
	}

	// Write output to file (aligned output in BAM format)
	err = os.WriteFile(outputFile, out.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("failed to write output file: %v", err)
	}

	return nil
}
