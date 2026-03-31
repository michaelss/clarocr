# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Compilar o binário
make build          # ou: go build -o clarocr .

# Gerar pacote .deb
make deb
make deb VERSION=1.2.0

# Limpar artefatos
make clean
```

O binário requer CGO (via `gosseract`). Antes de compilar, instale os headers:

```bash
sudo apt install libtesseract-dev libleptonica-dev libayatana-appindicator3-dev
```

## Arquitetura

O ponto de entrada (`main.go` + `main_install.go`) expõe três subcomandos: `capture`, `daemon` e `install`. O fluxo principal é orquestrado por `runCapture`, que encadeia as chamadas aos pacotes internos em sequência:

```
capture.SelectRegion() → capture.CaptureRegion() → ocr.ExtractText() → clipboard.Copy() → notify.Send()
```

**Detecção de ambiente (X11 vs Wayland)** — `capture.IsWayland()` verifica `WAYLAND_DISPLAY` e `XDG_SESSION_TYPE`. Cada pacote que precisa ramificar por ambiente chama essa função (ou replica a lógica localmente, como `clipboard`). As ferramentas externas usadas são:

| Função | X11 | Wayland |
|---|---|---|
| Seleção de região | `slop` | `slurp` |
| Captura | `maim` | `grim` |
| Clipboard | `xclip` | `wl-copy` |

**Modo daemon** — `runDaemon` passa um `tray.Config` com um callback `OnCapture` que aponta para `runCaptureOrLog`. O `tray` nunca importa `main`; a inversão de dependência é feita via `Config.OnCapture func(lang string)`. O idioma selecionado no menu da bandeja é mantido em `cfg.Lang` (string mutável no struct).

**Captura de região** — `capture.Region` é o tipo de transferência entre `selector.go` e `screenshot.go`. O `selector` parseia o formato de saída de cada ferramenta (`slurp`: `X,Y WxH`; `slop`: `X Y W H`) e entrega um `Region` uniforme para o `screenshot`.

**Ícone da bandeja** — embutido via `//go:embed icon.png` em `tray/tray.go`. Para trocar o ícone, substitua `tray/icon.png` (PNG 22×22).

## Dependências de runtime

O `.deb` gerado declara as dependências em `packaging/debian/control`. Ao adicionar novas ferramentas externas, atualizar o campo `Depends` nesse arquivo.
