package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	// cmd := exec.Command("docker run --rm -it --name myapp -v \"$PWD\":/app -w /app clearlinux/tesseract-ocr tesseract ../res/salary-info.png stdout --oem 1")
	// cmd := exec.Command("/usr/local/bin/docker", "run", "--rm", "grycap/cowsay")
	// cmd := exec.Command("docker", "run", "--rm", "grycap/cowsay", "/usr/games/cowsay", "Hello World")

	cmd := exec.Command(
		"docker",
		"run",
		"--rm",
		"--name",
		"myapp",
		"-v",
		// "$PWD:/app",
		"/Users/justincalleja/github/days-of-go/content/series/days-of-go/day-3:/app",
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

	// docker run --rm -it --name myapp \
	//   -v "$PWD":/app \
	//   -e "TESSDATA_PREFIX=/app/res" \
	//   -w /app \
	//   clearlinux/tesseract-ocr tesseract \
	//   ./res/salary-info.png stdout --oem 1

	// cmd := exec.Command("docker", "run", "--rm", "grycap/cowsay", "/usr/games/cowsay", "Hello World")
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	// cmdStr := "docker run --rm grycap/cowsay"
	// out, _ := exec.Command("/bin/sh", "-c", cmdStr).Output()
	// fmt.Printf("%s", out)

}
