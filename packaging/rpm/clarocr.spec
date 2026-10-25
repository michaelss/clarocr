%{!?clarocr_version:%global clarocr_version 0.0.1}
%{!?clarocr_release:%global clarocr_release 1}
%global debug_package %{nil}

Name:           clarocr
Version:        %{clarocr_version}
Release:        %{clarocr_release}%{?dist}
Summary:        Captura e extrai texto da tela via OCR
License:        LicenseRef-Proprietary AND Apache-2.0 AND MIT AND BSD-3-Clause
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
Requires:       gnome-screenshot
Requires:       libnotify
Requires:       /usr/bin/gsettings

%description
Ferramenta para Linux/GNOME que permite selecionar uma região da tela e
extrair seu texto com Tesseract OCR, copiando o resultado para o clipboard.
Suporta X11 e Wayland. Após instalar, execute clarocr install como usuário
da sessão gráfica para configurar autostart e o atalho Super+Shift+T.

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

%files
%{_bindir}/clarocr
%doc README.md
%license THIRD_PARTY_LICENSES

%changelog
* Sun Oct 25 2026 ClarOCR contributors - 0.0.1-1
- Empacotamento inicial para Fedora 44 x86_64.
