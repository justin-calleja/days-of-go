package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	cmd := exec.Command(
		"docker",
		"run",
		"--rm",
		"--name",
		"myapp",
		"-v",
		path+":/app",
		"-e",
		"TESSDATA_PREFIX=/app/res",
		"-w",
		"/app",
		"clearlinux/tesseract-ocr",
		"tesseract",
		"./res/salary-info.png",
		"stdout",
		"--oem",
		"1",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("OCR (via docker) failed with %s\n", err)
	}
}
