#!/bin/sh
# Run as root only inside a disposable Ubuntu container.
set -eu

apt-get update
# Minimal Ubuntu images may discard /usr/share/doc during installation.
# Verify license inclusion in the archive itself, independent of that policy.
dpkg-deb --contents "$1" > /tmp/clarocr-files.txt
grep -q '/THIRD_PARTY_LICENSES/vendor/github.com/godbus/dbus/v5/LICENSE$' /tmp/clarocr-files.txt
DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends "$1"
dpkg-query -W clarocr xdg-desktop-portal xdg-desktop-portal-gnome
test -x /usr/local/bin/clarocr
ldd /usr/local/bin/clarocr > /tmp/clarocr-ldd.txt
cat /tmp/clarocr-ldd.txt
if grep -q 'not found' /tmp/clarocr-ldd.txt; then
    exit 1
fi
clarocr capture -h
clarocr daemon -h
for tool in tesseract slop maim xclip slurp grim wl-copy notify-send gsettings; do
    command -v "$tool"
done
tesseract --list-langs > /tmp/clarocr-langs.txt
grep -qx eng /tmp/clarocr-langs.txt
grep -qx por /tmp/clarocr-langs.txt
apt-get remove -y clarocr
test ! -e /usr/local/bin/clarocr
echo 'DEB: instalação, bibliotecas, idiomas, comandos e remoção OK.'
