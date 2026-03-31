package tray

import (
	_ "embed"

	"github.com/getlantern/systray"
)

//go:embed icon.png
var iconData []byte

// Config holds runtime options for the tray daemon.
type Config struct {
	Lang           string
	OnCapture      func(lang string)
	AvailableLangs []LangOption
}

// LangOption represents a selectable OCR language.
type LangOption struct {
	Label string
	Value string
}

var defaultLangs = []LangOption{
	{Label: "Português + Inglês", Value: "por+eng"},
	{Label: "Português", Value: "por"},
	{Label: "Inglês", Value: "eng"},
}

// Run starts the system tray and blocks until the user selects Quit.
func Run(cfg *Config) {
	if cfg.AvailableLangs == nil {
		cfg.AvailableLangs = defaultLangs
	}
	systray.Run(func() { onReady(cfg) }, nil)
}

func onReady(cfg *Config) {
	systray.SetIcon(iconData)
	systray.SetTooltip("clarocr — Capturar texto da tela")

	mCapture := systray.AddMenuItem("Capturar texto", "Selecionar área e extrair texto")
	systray.AddSeparator()
	mLang := buildLangMenu(cfg)
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Sair", "Encerrar clarocr")

	go handleEvents(cfg, mCapture, mLang, mQuit)
}

func buildLangMenu(cfg *Config) []*systray.MenuItem {
	parent := systray.AddMenuItem("Idioma: "+currentLangLabel(cfg), "")
	items := make([]*systray.MenuItem, len(cfg.AvailableLangs))
	for i, l := range cfg.AvailableLangs {
		items[i] = parent.AddSubMenuItem(l.Label, l.Value)
		if l.Value == cfg.Lang {
			items[i].Check()
		}
	}
	return items
}

func handleEvents(cfg *Config, mCapture *systray.MenuItem, mLang []*systray.MenuItem, mQuit *systray.MenuItem) {
	for {
		select {
		case <-mCapture.ClickedCh:
			go cfg.OnCapture(cfg.Lang)
		case <-mQuit.ClickedCh:
			systray.Quit()
			return
		}
		checkLangSelection(cfg, mLang)
	}
}

func checkLangSelection(cfg *Config, items []*systray.MenuItem) {
	for i, item := range items {
		select {
		case <-item.ClickedCh:
			cfg.Lang = cfg.AvailableLangs[i].Value
			for _, it := range items {
				it.Uncheck()
			}
			item.Check()
		default:
		}
	}
}

func currentLangLabel(cfg *Config) string {
	for _, l := range cfg.AvailableLangs {
		if l.Value == cfg.Lang {
			return l.Label
		}
	}
	return cfg.Lang
}
