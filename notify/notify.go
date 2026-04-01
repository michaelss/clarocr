package notify

import "os/exec"

// Send displays a desktop notification via notify-send.
func Send(summary, body string) error {
	return exec.Command("notify-send", "--icon=edit-copy", "--expire-time=4000", summary, body).Run()
}
