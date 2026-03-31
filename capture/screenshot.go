package capture

import (
	"fmt"
	"os"
	"os/exec"
)

// CaptureRegion captures the given screen region to a temporary PNG file.
// The caller is responsible for removing the file after use.
// Uses grim on Wayland or maim on X11.
func CaptureRegion(r Region) (string, error) {
	f, err := os.CreateTemp("", "clarocr_*.png")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	f.Close()

	if IsWayland() {
		err = captureWithGrim(r, f.Name())
	} else {
		err = captureWithMaim(r, f.Name())
	}
	if err != nil {
		os.Remove(f.Name())
		return "", err
	}
	return f.Name(), nil
}

func captureWithGrim(r Region, dest string) error {
	geometry := fmt.Sprintf("%d,%d %dx%d", r.X, r.Y, r.W, r.H)
	cmd := exec.Command("grim", "-g", geometry, dest)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("grim failed: %w\n%s", err, out)
	}
	return nil
}

func captureWithMaim(r Region, dest string) error {
	geometry := fmt.Sprintf("%dx%d+%d+%d", r.W, r.H, r.X, r.Y)
	cmd := exec.Command("maim", "--geometry="+geometry, dest)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("maim failed: %w\n%s", err, out)
	}
	return nil
}
