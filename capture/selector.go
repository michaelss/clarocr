package capture

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// SelectRegion opens an interactive region selector and returns the selected Region.
// Uses slurp on Wayland or slop on X11.
func SelectRegion() (Region, error) {
	if IsWayland() {
		return selectWithSlurp()
	}
	return selectWithSlop()
}

func selectWithSlurp() (Region, error) {
	cmd := exec.Command("slurp")
	out, err := cmd.Output()
	if err != nil {
		return Region{}, selectorError("slurp", err)
	}
	return parseSlurpGeometry(strings.TrimSpace(string(out)))
}

func selectWithSlop() (Region, error) {
	cmd := exec.Command("slop", "--nodrag", "--bordersize=2", "--color=0.3,0.6,1,1", "--tolerance=0", "--format=%x %y %w %h")
	out, err := cmd.Output()
	if err != nil {
		return Region{}, selectorError("slop", err)
	}
	return parseSlopGeometry(strings.TrimSpace(string(out)))
}

func selectorError(tool string, err error) error {
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && len(exitErr.Stderr) > 0 {
		return fmt.Errorf("%s: %s", tool, strings.TrimSpace(string(exitErr.Stderr)))
	}
	return fmt.Errorf("%s failed: %w", tool, err)
}

// parseSlurpGeometry parses slurp output: "X,Y WxH"
func parseSlurpGeometry(s string) (Region, error) {
	var r Region
	_, err := fmt.Sscanf(s, "%d,%d %dx%d", &r.X, &r.Y, &r.W, &r.H)
	if err != nil {
		return Region{}, fmt.Errorf("unexpected slurp output %q: %w", s, err)
	}
	return r, nil
}

// parseSlopGeometry parses slop output: "X Y W H"
func parseSlopGeometry(s string) (Region, error) {
	var r Region
	_, err := fmt.Sscanf(s, "%d %d %d %d", &r.X, &r.Y, &r.W, &r.H)
	if err != nil {
		return Region{}, fmt.Errorf("unexpected slop output %q: %w", s, err)
	}
	return r, nil
}
