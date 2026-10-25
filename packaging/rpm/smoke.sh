#!/bin/sh
# Run as root only inside a disposable Fedora 44 container.
set -eu

dnf install -y --setopt=install_weak_deps=False "$1"
rpm -q clarocr
test -x /usr/bin/clarocr
ldd /usr/bin/clarocr > /tmp/clarocr-ldd.txt
cat /tmp/clarocr-ldd.txt
if grep -q 'not found' /tmp/clarocr-ldd.txt; then
    exit 1
fi

# Flag help loads the executable and CGO libraries without opening a display.
clarocr capture -h
clarocr daemon -h
for tool in tesseract slop maim xclip slurp grim wl-copy gnome-screenshot notify-send gsettings; do
    command -v "$tool"
done
tesseract --list-langs > /tmp/clarocr-langs.txt
grep -qx eng /tmp/clarocr-langs.txt
grep -qx por /tmp/clarocr-langs.txt

dnf remove -y clarocr
test ! -e /usr/bin/clarocr
echo 'RPM: instalação, bibliotecas, idiomas, comandos e remoção OK.'
