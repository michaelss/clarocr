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
	"clarocr/tray"
)

func main() {
	captureCmd := flag.NewFlagSet("capture", flag.ExitOnError)
	captureLang := captureCmd.String("lang", "por+eng", "Tesseract language string (e.g. por+eng)")

	daemonCmd := flag.NewFlagSet("daemon", flag.ExitOnError)
	daemonLang := daemonCmd.String("lang", "por+eng", "Default language for OCR")

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "capture":
		captureCmd.Parse(os.Args[2:])
		if err := runCapture(*captureLang); err != nil {
			log.Fatal(err)
		}
	case "daemon":
		daemonCmd.Parse(os.Args[2:])
		runDaemon(*daemonLang)
	case "install":
		if err := installAutostart(); err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr, "  clarocr capture [--lang por+eng]")
	fmt.Fprintln(os.Stderr, "  clarocr daemon  [--lang por+eng]")
	fmt.Fprintln(os.Stderr, "  clarocr install")
}

func runDaemon(lang string) {
	cfg := &tray.Config{
		Lang:      lang,
		OnCapture: func(l string) { runCaptureOrLog(l) },
	}
	tray.Run(cfg)
}

func runCaptureOrLog(lang string) {
	if err := runCapture(lang); err != nil {
		log.Println("capture error:", err)
		notify.Send("clarocr — Erro", err.Error())
	}
}

func runCapture(lang string) error {
	region, err := capture.SelectRegion()
	if err != nil {
		return fmt.Errorf("region selection cancelled or failed: %w", err)
	}

	imagePath, err := capture.CaptureRegion(region)
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
