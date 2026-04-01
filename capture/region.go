package capture

import (
	"os"
	"strings"
)

// Region represents a screen area defined by its top-left corner and dimensions.
type Region struct {
	X, Y, W, H int
}

// IsWayland returns true when running under a Wayland compositor.
func IsWayland() bool {
	return os.Getenv("WAYLAND_DISPLAY") != "" || os.Getenv("XDG_SESSION_TYPE") == "wayland"
}

// IsGnomeWayland returns true when running under GNOME on Wayland.
// GNOME's Mutter compositor does not support zwlr_layer_shell_v1, so
// slurp/grim cannot be used; gnome-screenshot must be used instead.
func IsGnomeWayland() bool {
	if !IsWayland() {
		return false
	}
	desktop := strings.ToLower(os.Getenv("XDG_CURRENT_DESKTOP"))
	return strings.Contains(desktop, "gnome")
}
