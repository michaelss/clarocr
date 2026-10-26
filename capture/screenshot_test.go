package capture

import (
	"image"
	"image/png"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/godbus/dbus/v5"
)

func TestScreenshotFromResponse(t *testing.T) {
	for _, scenario := range []string{"valid", "empty", "truncated", "missing", "cancelled", "failed", "malformed", "no-uri", "remote-uri"} {
		t.Run(scenario, func(t *testing.T) {
			dir := t.TempDir()
			t.Setenv("TMPDIR", dir)
			fixture := filepath.Join(t.TempDir(), "screenshot with spaces.png")
			f, err := os.Create(fixture)
			if err != nil {
				t.Fatal(err)
			}
			if err := png.Encode(f, image.NewRGBA(image.Rect(0, 0, 2, 2))); err != nil {
				t.Fatal(err)
			}
			if err := f.Close(); err != nil {
				t.Fatal(err)
			}
			uri := (&url.URL{Scheme: "file", Path: fixture}).String()
			body := []interface{}{uint32(0), map[string]dbus.Variant{"uri": dbus.MakeVariant(uri)}}
			switch scenario {
			case "empty":
				err = os.Truncate(fixture, 0)
			case "truncated":
				err = os.Truncate(fixture, 40)
			case "missing":
				err = os.Remove(fixture)
			case "cancelled":
				body[0] = uint32(1)
			case "failed":
				body[0] = uint32(2)
			case "malformed":
				body = []interface{}{"bad"}
			case "no-uri":
				body[1] = map[string]dbus.Variant{}
			case "remote-uri":
				body[1] = map[string]dbus.Variant{"uri": dbus.MakeVariant("file://remote/image.png")}
			}
			if err != nil {
				t.Fatal(err)
			}
			path, err := screenshotFromResponse(body)
			if scenario == "valid" {
				if err != nil {
					t.Fatal(err)
				}
				if path == fixture {
					t.Fatal("must return a copy, not the portal's original")
				}
				if err := validateScreenshot(path); err != nil {
					t.Fatal(err)
				}
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
				if _, err := os.Stat(fixture); err != nil {
					t.Fatalf("portal's original must be preserved: %v", err)
				}
			} else if err == nil || path != "" {
				t.Fatalf("expected error and no path, got %q, %v", path, err)
			} else if scenario == "cancelled" && !strings.Contains(err.Error(), "cancelled") {
				t.Fatalf("expected cancellation diagnostic: %v", err)
			}
			entries, err := os.ReadDir(dir)
			if err != nil || len(entries) != 0 {
				t.Fatalf("temporary screenshot leak: %v, %v", entries, err)
			}
		})
	}
}
