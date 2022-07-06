package ocr

import "bytes"

type OCRer interface {
	RunOCR(inPath string) error
	OutBuffer() *bytes.Buffer
	ErrBuffer() *bytes.Buffer
}
