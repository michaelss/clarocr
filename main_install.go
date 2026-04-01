package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

func installShortcut() error {
	const bindingPath = "/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/clarocr/"
	const schemaList = "org.gnome.settings-daemon.plugins.media-keys"
	const schemaBinding = "org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:" + bindingPath

	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("could not determine executable path: %w", err)
	}

	// Lê a lista atual de atalhos customizados
	out, err := exec.Command("gsettings", "get", schemaList, "custom-keybindings").Output()
	if err != nil {
		return fmt.Errorf("could not read gsettings: %w", err)
	}

	current := strings.TrimSpace(string(out))

	// Parseia a lista: "@as []" (vazia) ou "['path1', 'path2']"
	var paths []string
	if current != "@as []" {
		inner := strings.Trim(current, "[]")
		for _, p := range strings.Split(inner, ",") {
			p = strings.TrimSpace(strings.Trim(strings.TrimSpace(p), "'"))
			if p != "" {
				paths = append(paths, p)
			}
		}
	}

	// Adiciona nosso path se ainda não estiver na lista
	found := false
	for _, p := range paths {
		if p == bindingPath {
			found = true
			break
		}
	}
	if !found {
		paths = append(paths, bindingPath)
		quoted := make([]string, len(paths))
		for i, p := range paths {
			quoted[i] = "'" + p + "'"
		}
		newList := "[" + strings.Join(quoted, ", ") + "]"
		if err := exec.Command("gsettings", "set", schemaList, "custom-keybindings", newList).Run(); err != nil {
			return fmt.Errorf("could not update custom-keybindings list: %w", err)
		}
	}

	command := execPath + " capture"
	for _, args := range [][]string{
		{schemaBinding, "name", "ClaroOCR"},
		{schemaBinding, "command", command},
		{schemaBinding, "binding", "<Super><Shift>t"},
	} {
		if err := exec.Command("gsettings", append([]string{"set"}, args...)...).Run(); err != nil {
			return fmt.Errorf("could not set %s: %w", args[1], err)
		}
	}

	fmt.Println("Atalho Super+Shift+T configurado com sucesso.")
	return nil
}
