package capture

import "os"

// Region represents a screen area defined by its top-left corner and dimensions.
type Region struct {
	X, Y, W, H int
}

// IsWayland returns true when running under a Wayland compositor.
func IsWayland() bool {
	return os.Getenv("WAYLAND_DISPLAY") != "" || os.Getenv("XDG_SESSION_TYPE") == "wayland"
}
