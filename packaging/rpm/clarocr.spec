%{!?clarocr_version:%global clarocr_version 0.0.1}
%{!?clarocr_release:%global clarocr_release 1}
%global debug_package %{nil}

Name:           clarocr
Version:        %{clarocr_version}
Release:        %{clarocr_release}%{?dist}
Summary:        Captura e extrai texto da tela via OCR
License:        LicenseRef-Proprietary AND Apache-2.0 AND MIT AND BSD-3-Clause AND BSD-2-Clause
URL:            https://github.com/michaelss/clarocr
Source0:        %{name}-%{version}.tar.gz
ExclusiveArch:  x86_64

BuildRequires:  golang >= 1.19
BuildRequires:  gcc
BuildRequires:  gcc-c++
BuildRequires:  make
BuildRequires:  pkgconf-pkg-config
BuildRequires:  tesseract-devel
BuildRequires:  leptonica-devel
BuildRequires:  libayatana-appindicator-gtk3-devel

Requires:       tesseract
Requires:       tesseract-langpack-por
Requires:       tesseract-langpack-eng
Requires:       slop
Requires:       maim
Requires:       xclip
Requires:       slurp
Requires:       grim
Requires:       wl-clipboard
Requires:       xdg-desktop-portal
Requires:       xdg-desktop-portal-gnome
Requires:       libnotify
Requires:       /usr/bin/gsettings
Requires(post): /usr/sbin/runuser
Requires(post): /usr/bin/getent
Requires(post): /usr/bin/id

%description
Ferramenta para Linux/GNOME que permite selecionar uma região da tela e
extrair seu texto com Tesseract OCR, copiando o resultado para o clipboard.
Suporta X11 e Wayland. A instalação tenta configurar autostart e o atalho
Super+Shift+T para o usuário que iniciou a instalação na sessão gráfica.
Se não houver sessão disponível, execute clarocr install como usuário.

%prep
%setup -q

%build
CGO_ENABLED=1 GOPROXY=off GOSUMDB=off GOTOOLCHAIN=local \
    go build -mod=vendor -trimpath -buildvcs=false -ldflags='-s -w' -o clarocr .

%install
install -Dpm 0755 clarocr %{buildroot}%{_bindir}/clarocr

%check
CGO_ENABLED=1 GOPROXY=off GOSUMDB=off GOTOOLCHAIN=local \
    go test -mod=vendor ./...

%post
# Configure only on first installation; upgrades must preserve user choices.
if [ "$1" -eq 1 ]; then
    setup_user=${SUDO_USER:-}
    if [ -z "$setup_user" ] && [ -n "${PKEXEC_UID:-}" ]; then
        setup_user=$(getent passwd "$PKEXEC_UID" | cut -d: -f1)
    fi
    setup_uid=
    if [ -n "$setup_user" ] && [ "$setup_user" != root ]; then
        setup_uid=$(id -u "$setup_user" 2>/dev/null) || setup_uid=
    fi
    if [ -n "$setup_uid" ] && [ "$setup_uid" != 0 ] && [ -S "/run/user/$setup_uid/bus" ]; then
        if ! /usr/sbin/runuser -u "$setup_user" -- env \
            XDG_RUNTIME_DIR="/run/user/$setup_uid" \
            DBUS_SESSION_BUS_ADDRESS="unix:path=/run/user/$setup_uid/bus" \
            %{_bindir}/clarocr install; then
            echo 'ClaroOCR: execute clarocr install na sessão gráfica para concluir a configuração.'
        fi
    else
        echo 'ClaroOCR: execute clarocr install na sessão gráfica para configurar autostart e Super+Shift+T.'
    fi
fi
exit 0

%files
%{_bindir}/clarocr
%doc README.md
%license THIRD_PARTY_LICENSES

%changelog
* Sun Oct 25 2026 ClarOCR contributors - 0.0.1-1
- Empacotamento inicial para Fedora 44 x86_64.
