package clipboard

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestCopyReturnsWhileClipboardOwnerIsRunning(t *testing.T) {
	for _, tool := range []string{"wl-copy", "xclip"} {
		t.Run(tool, func(t *testing.T) {
			dir := t.TempDir()
			t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
			t.Setenv("WAYLAND_DISPLAY", "")
			t.Setenv("XDG_SESSION_TYPE", "x11")
			if tool == "wl-copy" {
				t.Setenv("XDG_SESSION_TYPE", "wayland")
			}
			textPath := filepath.Join(dir, "copied.txt")
			t.Setenv("COPIED_TEXT", textPath)
			// Model a clipboard owner that inherits stdout/stderr after forking.
			script := "#!/bin/sh\ncat > \"$COPIED_TEXT\"\nsleep 2 &\nexit 0\n"
			if err := os.WriteFile(filepath.Join(dir, tool), []byte(script), 0700); err != nil {
				t.Fatal(err)
			}
			done := make(chan error, 1)
			go func() { done <- Copy("texto de teste\nsegunda linha") }()
			select {
			case err := <-done:
				if err != nil {
					t.Fatal(err)
				}
			case <-time.After(time.Second):
				// Wait for Copy before restoring the test's environment.
				<-done
				t.Fatal("Copy waited for the background clipboard owner")
			}
			got, err := os.ReadFile(textPath)
			if err != nil || string(got) != "texto de teste\nsegunda linha" {
				t.Fatalf("unexpected clipboard content %q: %v", got, err)
			}
		})
	}
}

func TestClipboardCommandFailure(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
	for _, tool := range []string{"wl-copy", "xclip"} {
		if err := os.WriteFile(filepath.Join(dir, tool), []byte("#!/bin/sh\nexit 7\n"), 0700); err != nil {
			t.Fatal(err)
		}
	}
	for _, copyText := range []func(string) error{copyWayland, copyX11} {
		if err := copyText("test"); err == nil || !strings.Contains(err.Error(), "exit status 7") {
			t.Fatalf("expected clipboard failure, got %v", err)
		}
	}
}
