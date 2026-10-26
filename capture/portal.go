package capture

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"time"

	"github.com/godbus/dbus/v5"
)

// CapturePortalScreenshot asks the desktop to present its interactive screenshot
// dialog. The returned file is owned by ClarOCR and must be removed by the caller.
func CapturePortalScreenshot() (string, error) {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return "", fmt.Errorf("connect to desktop portal: %w", err)
	}
	defer conn.Close()

	// Subscribe before calling Screenshot: Response can arrive before its reply.
	if err := conn.AddMatchSignal(
		dbus.WithMatchSender("org.freedesktop.portal.Desktop"),
		dbus.WithMatchInterface("org.freedesktop.portal.Request"),
		dbus.WithMatchMember("Response"),
	); err != nil {
		return "", fmt.Errorf("subscribe to screenshot response: %w", err)
	}
	signals := make(chan *dbus.Signal, 16)
	conn.Signal(signals)
	defer conn.RemoveSignal(signals)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	portal := conn.Object("org.freedesktop.portal.Desktop", "/org/freedesktop/portal/desktop")
	var handle dbus.ObjectPath
	err = portal.CallWithContext(ctx, "org.freedesktop.portal.Screenshot.Screenshot", 0,
		"", map[string]dbus.Variant{"interactive": dbus.MakeVariant(true)},
	).Store(&handle)
	if err != nil {
		return "", fmt.Errorf("desktop screenshot portal (check xdg-desktop-portal-gnome): %w", err)
	}
	for {
		select {
		case signal, ok := <-signals:
			if !ok {
				return "", fmt.Errorf("desktop portal disconnected")
			}
			if signal == nil || signal.Path != handle {
				continue
			}
			return screenshotFromResponse(signal.Body)
		case <-ctx.Done():
			closeCtx, closeCancel := context.WithTimeout(context.Background(), time.Second)
			conn.Object("org.freedesktop.portal.Desktop", handle).CallWithContext(closeCtx, "org.freedesktop.portal.Request.Close", 0)
			closeCancel()
			return "", fmt.Errorf("timed out waiting for screenshot selection: %w", ctx.Err())
		}
	}
}

func screenshotFromResponse(body []interface{}) (string, error) {
	var response uint32
	var results map[string]dbus.Variant
	if err := dbus.Store(body, &response, &results); err != nil {
		return "", fmt.Errorf("invalid screenshot portal response: %w", err)
	}
	if response == 1 {
		return "", fmt.Errorf("screenshot selection cancelled")
	}
	if response != 0 {
		return "", fmt.Errorf("desktop screenshot portal failed (response %d)", response)
	}
	uri, ok := results["uri"].Value().(string)
	if !ok {
		return "", fmt.Errorf("screenshot portal returned no image URI")
	}
	u, err := url.Parse(uri)
	if err != nil || u.Scheme != "file" || (u.Host != "" && u.Host != "localhost") || u.Path == "" {
		return "", fmt.Errorf("invalid screenshot file URI: %q", uri)
	}
	source, err := os.Open(u.Path)
	if err != nil {
		return "", fmt.Errorf("open portal screenshot: %w", err)
	}
	defer source.Close()
	// Copy to our own file: the portal's image may live in the user's Pictures.
	f, err := os.CreateTemp("", "clarocr_*.png")
	if err != nil {
		return "", err
	}
	_, copyErr := io.Copy(f, source)
	closeErr := f.Close()
	if copyErr != nil {
		err = copyErr
	} else if closeErr != nil {
		err = closeErr
	} else {
		err = validateScreenshot(f.Name())
	}
	if err != nil {
		os.Remove(f.Name())
		return "", fmt.Errorf("read portal screenshot: %w", err)
	}
	return f.Name(), nil
}
