package ocr

import (
	"fmt"
	"strings"

	"github.com/otiai10/gosseract/v2"
)

// ExtractText runs Tesseract OCR on the given image file and returns the extracted text.
// langs is a Tesseract language string, e.g. "por+eng".
func ExtractText(imagePath string, langs string) (string, error) {
	client := gosseract.NewClient()
	defer client.Close()

	if err := client.SetImage(imagePath); err != nil {
		return "", fmt.Errorf("failed to load image: %w", err)
	}

	if langs != "" {
		if err := client.SetLanguage(strings.Split(langs, "+")...); err != nil {
			return "", fmt.Errorf("failed to set language: %w", err)
		}
	}

	text, err := client.Text()
	if err != nil {
		return "", fmt.Errorf("OCR failed: %w", err)
	}

	return strings.TrimSpace(text), nil
}
