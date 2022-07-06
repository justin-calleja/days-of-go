package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/justin-calleja/days-of-go/content/series/days-of-go/day-3/checkpoint-2/ocr"
	"github.com/justin-calleja/days-of-go/content/series/days-of-go/day-3/checkpoint-2/utils"
)

func main() {
	ocrInPath := "res/salary-info.png"
	ocrOutPath := "res/ocr-output.txt"

	tocr := ocr.NewTesseractOcr()

	if exists, err := utils.FileExists(ocrOutPath); exists {
		// TODO: read file
	} else if !exists {
		if err := tocr.RunOCR(ocrInPath); err != nil {
			log.Fatalf("OCR on \"%s\" failed with: %s\n", ocrInPath, err)
			os.Exit(1)
		}

		outBuff := tocr.OutBuffer()
		errBuff := tocr.ErrBuffer()
		if errBuff.Len() > 0 {
			fmt.Printf("OCR stderr: %s\n", errBuff)
		}

		file, err := os.Create(ocrOutPath)
		if err != nil {
			log.Fatalf("Creating / truncating %s failed with: %s\n", ocrOutPath, err)
			os.Exit(1)
		}
		defer file.Close()

		fmt.Printf("about to print the trim:\n")
		// for pos, char := range outBuff.String() {
		// 	fmt.Printf("character %d starts at byte position %d\n", char, pos)
		// }
		fmt.Printf("%U\n", outBuff)

		// fmt.Printf(strings.Trim(outBuff.String(), "\t"))
		// file.Write(bytes.Trim(outBuff.Bytes(), "\r\n"))
		if _, err := file.WriteString(strings.Trim(outBuff.String(), "\t")); err != nil {
			log.Fatalf("Writing to \"%s\" failed with: %s\n", ocrOutPath, err)
			os.Exit(1)
		}
		// if _, err := file.Write(bytes.TrimSpace(outBuff.Bytes())); err != nil {
		// 	log.Fatalf("Writing to \"%s\" failed with: %s\n", ocrOutPath, err)
		// 	os.Exit(1)
		// }
	} else {
		log.Fatalf("Failed to check existance of %s with error: %s\n", ocrOutPath, err)
		os.Exit(1)
	}
}

func Parse(s string) {
	fmt.Printf("TODO: need to parse:\n%s\n", s)
}
