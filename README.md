# clarocr

Ferramenta para Linux/GNOME que permite selecionar uma região da tela e extrair o texto presente nela via OCR, copiando o resultado automaticamente para a área de transferência.

![Demonstração](https://github.com/michaelss/clarocr/blob/master/ClarOCR.gif)

Casos de uso típicos:
- Extrair texto de vídeos pausados no navegador
- Capturar texto de imagens ou PDFs exibidos na tela
- Copiar trechos de aplicações que não permitem seleção de texto

## Ambiente

- **Sistema operacional:** Linux
- **Desktop:** GNOME
- **Display server:** X11 e Wayland (detecção automática)

## Instalação

1. Baixe o pacote na área "Releases", ao lado: `.deb` para Debian/Ubuntu ou `.rpm` para Fedora 44 x86_64.
2. Instale com `sudo apt install ./<arquivo>.deb` ou `sudo dnf install ./<arquivo>.rpm`.
3. Na sessão gráfica, execute `clarocr install` para configurar o autostart e o atalho Super+Shift+T.


## Tecnologias

| Componente | Tecnologia |
|---|---|
| Linguagem | [Go](https://go.dev) |
| OCR engine | [Tesseract OCR](https://github.com/tesseract-ocr/tesseract) via [gosseract](https://github.com/otiai10/gosseract) |
| Seleção de região (X11) | [slop](https://github.com/naelstrof/slop) |
| Seleção de região (Wayland) | [slurp](https://github.com/emersion/slurp) |
| Captura de tela (X11) | [maim](https://github.com/naelstrof/maim) |
| Captura de tela (Wayland) | [grim](https://github.com/emersion/grim) |
| Clipboard (X11) | [xclip](https://github.com/astrand/xclip) |
| Clipboard (Wayland) | [wl-clipboard](https://github.com/bugaevc/wl-clipboard) |
| Ícone na bandeja | [systray](https://github.com/getlantern/systray) via libayatana-appindicator3 |
| Notificações | [libnotify](https://gitlab.gnome.org/GNOME/libnotify) (`notify-send`) |

## Dependências do sistema

### Para compilação

```bash
sudo apt install \
  libtesseract-dev \
  libleptonica-dev \
  libayatana-appindicator3-dev
```

### Para execução

```bash
# X11
sudo apt install \
  libtesseract5 tesseract-ocr tesseract-ocr-por tesseract-ocr-eng \
  slop maim xclip libayatana-appindicator3-1 libnotify-bin

# Wayland (substituir slop/maim/xclip por:)
sudo apt install slurp grim wl-clipboard
```

> **Nota:** No GNOME, o ícone na bandeja requer a extensão [AppIndicator](https://extensions.gnome.org/extension/615/appindicator-support/).
> ```bash
> sudo apt install gnome-shell-extension-appindicator
> ```

## Compilação

```bash
# Compilar o binário
make build

# Gerar pacote .deb
make deb

# Gerar pacote .deb com versão específica
make deb VERSION=1.2.0
```

O pacote gerado estará em `dist/clarocr_<versão>_amd64.deb`.

### Pacote RPM (Fedora 44 x86_64)

Os comandos têm funções diferentes:

- `make rpm-deps` instala as ferramentas e bibliotecas necessárias para gerar
  o RPM, usando `sudo dnf`. Execute no primeiro uso ou quando faltarem
  dependências de compilação. Esse comando não gera o pacote.
- `make rpm` compila a aplicação, executa as verificações e gera o pacote em
  `dist/`. Não instala dependências nem executa `sudo`.

No primeiro uso:

```bash
make rpm-deps
make rpm
```

O pacote padrão será `dist/clarocr-0.0.1-1.fc44.x86_64.rpm`.
Nas próximas compilações, basta executar `make rpm`.

Para definir a versão da aplicação e a revisão do pacote:

```bash
make rpm VERSION=1.2.0 RELEASE=1
```

O resultado estará em `dist/clarocr-1.2.0-1.fc44.x86_64.rpm`. `VERSION` aceita
números separados por pontos, com padrão `0.0.1`; `RELEASE` é um inteiro
positivo, com padrão `1`. Incremente `RELEASE` ao revisar o empacotamento
sem alterar a versão da aplicação.

Compile no Fedora 44 x86_64: o binário usa CGO e depende das bibliotecas dessa
distribuição. O empacotamento prepara os módulos Go em `dist/` (com acesso à
internet ou cache Go preenchido); a compilação e as verificações do RPM usam
esses módulos vendorizados, sem baixar dependências Go. Não há garantia de
compatibilidade com outras versões do Fedora ou outras distribuições RPM.

```bash
sudo dnf install ./dist/clarocr-1.2.0-1.fc44.x86_64.rpm
clarocr install
```

O RPM instala o executável em `/usr/bin/clarocr`. A configuração de autostart
e atalho é feita pelo próprio usuário, sem scripts de instalação que alterem
a sessão gráfica. Para o ícone na bandeja do GNOME, instale também a extensão
AppIndicator (`sudo dnf install gnome-shell-extension-appindicator`).

O projeto não declara uma licença aberta. O RPM usa `LicenseRef-Proprietary`
para o código da aplicação, declara também as licenças dos componentes
incorporados e inclui os avisos de licença do Go e de suas dependências.
Os pacotes são distribuídos sem assinatura RPM e sem repositório DNF próprio.

## Instalação via .deb

```bash
sudo dpkg -i dist/clarocr_1.0.0_amd64.deb
```

O apt resolverá as dependências automaticamente. Após instalar, execute `clarocr install` para configurar o autostart.

## Execução

### Captura avulsa

Abre o seletor de região, extrai o texto e copia para o clipboard:

```bash
clarocr capture
clarocr capture --lang por+eng   # idioma padrão
clarocr capture --lang eng        # somente inglês
```

Os códigos de idioma seguem o padrão Tesseract. Combine idiomas com `+` (ex: `por+eng+spa`).

### Daemon com ícone na bandeja

Inicia em background com ícone na bandeja do sistema:

```bash
clarocr daemon
clarocr daemon --lang eng
```

O menu da bandeja permite acionar a captura e trocar o idioma sem reiniciar o daemon.

### Configurar autostart

Registra o daemon para iniciar automaticamente no login:

```bash
clarocr install
```

### Atalho de teclado

Para acionar a captura por atalho, adicione um atalho personalizado nas configurações do GNOME apontando para:

```
clarocr capture --lang por+eng
```

Ou configure via linha de comando (substitua `<Super><Shift>t` pelo atalho desejado):

```bash
BINARY=$(which clarocr)

gsettings set org.gnome.settings-daemon.plugins.media-keys custom-keybindings \
  "['/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/clarocr/']"

gsettings set org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/clarocr/ \
  name 'Capturar texto'

gsettings set org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/clarocr/ \
  command "$BINARY capture --lang por+eng"

gsettings set org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/clarocr/ \
  binding '<Super><Shift>t'
```

## Release via GitHub Actions

Ao criar uma tag numérica `v*` (por exemplo, `v1.0.0`), o workflow em
`.github/workflows/release.yml` compila o `.deb` no Ubuntu e o `.rpm` em um
contêiner Fedora 44. Após os dois builds passarem, publica ambos na mesma
release do GitHub:

```bash
git tag v1.0.0
git push origin v1.0.0
```
