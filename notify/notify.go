package notify

import (
	"os/exec"
	"unicode/utf8"
)

const previewMaxRunes = 80

// Send displays a desktop notification via notify-send.
func Send(summary, body string) error {
	cmd := exec.Command("notify-send", "--icon=edit-copy", "--expire-time=4000", summary, body)
	return cmd.Run()
}

// TextPreview truncates text to previewMaxRunes for use in notifications.
func TextPreview(text string) string {
	if utf8.RuneCountInString(text) <= previewMaxRunes {
		return text
	}
	runes := []rune(text)
	return string(runes[:previewMaxRunes]) + "…"
}
