#!/bin/sh
set -eu

destination=$1
mkdir -p "$destination"
work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT
trap 'exit 1' HUP INT TERM

go mod vendor -o "$work/vendor"
cp "$work/vendor/modules.txt" "$destination/modules.txt"
go_license="$(go env GOROOT)/LICENSE"
if [ ! -f "$go_license" ]; then
    go_license=/usr/share/licenses/golang/LICENSE
fi
if [ ! -f "$go_license" ]; then
    go_license="$(dirname "$0")/../Go-LICENSE"
fi
cp "$go_license" "$destination/Go-LICENSE"
(
    cd "$work"
    find vendor -type f \( -iname 'license*' -o -iname 'copying*' -o -iname 'notice*' -o -iname 'copyright*' \) \
        -exec cp --parents '{}' "$destination/" \;
)
