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
	return runClipboardCommand(cmd)
}

func copyX11(text string) error {
	cmd := exec.Command("xclip", "-selection", "clipboard")
	cmd.Stdin = strings.NewReader(text)
	return runClipboardCommand(cmd)
}

func runClipboardCommand(cmd *exec.Cmd) error {
	// Clipboard tools fork a background owner to keep the copied text available.
	// Use file descriptors directly: CombinedOutput creates pipes, and Run then
	// waits for EOF until that background owner also closes its inherited pipes.
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s failed: %w", cmd.Args[0], err)
	}
	return nil
}

func isWayland() bool {
	return os.Getenv("WAYLAND_DISPLAY") != "" || os.Getenv("XDG_SESSION_TYPE") == "wayland"
}
