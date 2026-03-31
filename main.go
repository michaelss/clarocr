package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"clarocr/capture"
	"clarocr/clipboard"
	"clarocr/notify"
	"clarocr/ocr"
)

func main() {
	captureCmd := flag.NewFlagSet("capture", flag.ExitOnError)
	lang := captureCmd.String("lang", "por+eng", "Tesseract language string (e.g. por+eng)")

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: clarocr capture [--lang por+eng]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "capture":
		captureCmd.Parse(os.Args[2:])
		if err := runCapture(*lang); err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
		fmt.Fprintln(os.Stderr, "Usage: clarocr capture [--lang por+eng]")
		os.Exit(1)
	}
}

func runCapture(lang string) error {
	geometry, err := capture.SelectRegion()
	if err != nil {
		return fmt.Errorf("region selection cancelled or failed: %w", err)
	}

	imagePath, err := capture.CaptureRegion(geometry)
	if err != nil {
		return fmt.Errorf("screen capture failed: %w", err)
	}
	defer os.Remove(imagePath)

	text, err := ocr.ExtractText(imagePath, lang)
	if err != nil {
		return fmt.Errorf("OCR failed: %w", err)
	}

	if text == "" {
		notify.Send("clarocr", "Nenhum texto encontrado na região selecionada.")
		return nil
	}

	if err := clipboard.Copy(text); err != nil {
		return fmt.Errorf("clipboard copy failed: %w", err)
	}

	notify.Send("Texto copiado!", notify.TextPreview(text))
	return nil
}
