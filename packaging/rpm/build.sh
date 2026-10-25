#!/bin/sh
set -eu

version=${1:-0.0.1}
release=${2:-1}

# Restrict values used in paths and RPM macros. Release tags use vX.Y.Z.
valid_chars=true
case "$version$release" in
    *[!0-9.]*) valid_chars=false ;;
esac
if [ "$valid_chars" != true ] ||
   ! printf '%s\n' "$version" | LC_ALL=C grep -Eq '^[0-9]+(\.[0-9]+)*$' ||
   ! printf '%s\n' "$release" | LC_ALL=C grep -Eq '^[1-9][0-9]*$'; then
    echo 'Use VERSION numérica (ex.: 1.2.0) e RELEASE inteiro positivo.' >&2
    exit 1
fi

for tool in go rpmbuild rpmspec tar; do
    if ! command -v "$tool" >/dev/null 2>&1; then
        echo "Pré-requisito ausente: $tool. No Fedora, execute make rpm-deps." >&2
        exit 1
    fi
done

root=$(pwd)
topdir="$root/dist/rpmbuild"
spec="$root/packaging/rpm/clarocr.spec"
arch=$(rpm --eval '%{_arch}')
platform=$(. /etc/os-release; printf '%s:%s' "$ID" "$VERSION_ID")
if [ "$arch" != x86_64 ] || [ "$platform" != fedora:44 ]; then
    echo 'O RPM é suportado inicialmente apenas em Fedora 44 x86_64.' >&2
    exit 1
fi

# Check declared BuildRequires before downloading Go modules.
if ! rpmbuild -bp --nobuild --define "_topdir $topdir" \
    --define "clarocr_version $version" --define "clarocr_release $release" "$spec"; then
    echo 'Verifique os pré-requisitos com make rpm-deps.' >&2
    exit 1
fi

mkdir -p "$topdir/SOURCES" "$topdir/RPMS"
work=$(mktemp -d "$root/dist/rpm-source.XXXXXX")
trap 'rm -rf "$work"' EXIT
trap 'exit 1' HUP INT TERM
source_dir="$work/clarocr-$version"
mkdir -p "$source_dir"

# Copy only application/build sources; never include .git, secrets or old artifacts.
cp "$root/"*.go "$root/go.mod" "$root/go.sum" "$root/README.md" "$source_dir/"
cp -R "$root/capture" "$root/clipboard" "$root/notify" "$root/ocr" "$root/tray" "$source_dir/"
(
    cd "$source_dir"
    GOTOOLCHAIN=local go mod vendor
    mkdir THIRD_PARTY_LICENSES
    cp vendor/modules.txt THIRD_PARTY_LICENSES/modules.txt
    # Fedora stores the Go license outside GOROOT; upstream toolchains keep it there.
    go_license="$(go env GOROOT)/LICENSE"
    if [ ! -f "$go_license" ]; then
        go_license=/usr/share/licenses/golang/LICENSE
    fi
    cp "$go_license" THIRD_PARTY_LICENSES/Go-LICENSE
    find vendor -type f \( -iname 'license*' -o -iname 'copying*' -o -iname 'notice*' -o -iname 'copyright*' \) \
        -exec cp --parents '{}' THIRD_PARTY_LICENSES/ \;
)
tar -czf "$topdir/SOURCES/clarocr-$version.tar.gz" -C "$work" "clarocr-$version"

rpmbuild -bb --define "_topdir $topdir" \
    --define "clarocr_version $version" --define "clarocr_release $release" "$spec"

filename=$(rpmspec -q --rpms --qf '%{NAME}-%{VERSION}-%{RELEASE}.%{ARCH}.rpm' \
    --define "clarocr_version $version" --define "clarocr_release $release" "$spec")
cp "$topdir/RPMS/$arch/$filename" "$root/dist/$filename"
printf '\nPacote gerado: dist/%s\nPara instalar: sudo dnf install ./dist/%s\n' "$filename" "$filename"
