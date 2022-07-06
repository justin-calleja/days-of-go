package ocr

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

type tesseractOcr struct {
	stdout *bytes.Buffer
	stderr *bytes.Buffer
}

func NewTesseractOcr() OCRer {
	return &tesseractOcr{
		stdout: &bytes.Buffer{},
		stderr: &bytes.Buffer{},
	}
}

func (tocr *tesseractOcr) RunOCR(inPath string) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory failed with: %w", err)
	}

	cmd := exec.Command(
		"docker",
		"run",
		"--rm",
		"--name",
		"myapp",
		"-v",
		wd+":/app",
		"-e",
		"TESSDATA_PREFIX=/app/res",
		"-w",
		"/app",
		"clearlinux/tesseract-ocr",
		"tesseract",
		inPath,
		"stdout",
		"--oem",
		"1",
	)

	cmd.Stdout = tocr.stdout
	cmd.Stderr = tocr.stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("OCR failed with: %w", err)
	}

	return nil
}

func (tocr *tesseractOcr) OutBuffer() *bytes.Buffer {
	return tocr.stdout
}

func (tocr *tesseractOcr) ErrBuffer() *bytes.Buffer {
	return tocr.stderr
}
