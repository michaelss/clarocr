package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const desktopEntry = `[Desktop Entry]
Type=Application
Name=clarocr
Comment=Captura e extrai texto da tela
Exec=%s daemon
Icon=edit-copy
Terminal=false
Categories=Utility;
X-GNOME-Autostart-enabled=true
`

func installAutostart() error {
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("could not determine executable path: %w", err)
	}

	autostartDir := filepath.Join(os.Getenv("HOME"), ".config", "autostart")
	if err := os.MkdirAll(autostartDir, 0755); err != nil {
		return fmt.Errorf("could not create autostart dir: %w", err)
	}

	destPath := filepath.Join(autostartDir, "clarocr.desktop")
	content := fmt.Sprintf(desktopEntry, execPath)
	if err := os.WriteFile(destPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("could not write desktop entry: %w", err)
	}

	fmt.Printf("Autostart instalado em: %s\n", destPath)
	fmt.Println("O clarocr será iniciado automaticamente no próximo login.")
	return nil
}
