package clipboard

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Copy writes text to the system clipboard.
// Uses wl-copy on Wayland or xclip on X11.
func Copy(text string) error {
	if isWayland() {
		return copyWayland(text)
	}
	return copyX11(text)
}

func copyWayland(text string) error {
	cmd := exec.Command("wl-copy")
	cmd.Stdin = strings.NewReader(text)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("wl-copy failed: %w\n%s", err, out)
	}
	return nil
}

func copyX11(text string) error {
	cmd := exec.Command("xclip", "-selection", "clipboard")
	cmd.Stdin = strings.NewReader(text)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("xclip failed: %w\n%s", err, out)
	}
	return nil
}

func isWayland() bool {
	return os.Getenv("WAYLAND_DISPLAY") != "" || os.Getenv("XDG_SESSION_TYPE") == "wayland"
}
