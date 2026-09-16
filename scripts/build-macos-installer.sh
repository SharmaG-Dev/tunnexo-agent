#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
version=$(sed -n 's/^const version = "\([^"]*\)"/\1/p' cmd/version.go)
[[ "$version" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]
command -v pkgbuild >/dev/null
command -v productbuild >/dev/null
stage=$(mktemp -d)
trap 'rm -rf "$stage"' EXIT
mkdir -p dist/installers
for arch in arm64 amd64; do
  root="$stage/$arch/root"
  mkdir -p "$root/usr/local/lib/tunnexo" "$root/etc/paths.d"
  CGO_ENABLED=0 GOOS=darwin GOARCH="$arch" go build -trimpath -ldflags='-s -w' -o "$root/usr/local/lib/tunnexo/tunnexo" .
  cp LICENSE README.md TUNNEXO_SETUP.md "$root/usr/local/lib/tunnexo/"
  printf '%s\n' /usr/local/lib/tunnexo > "$root/etc/paths.d/tunnexo"
  pkgbuild --root "$root" --identifier live.tunnexo.agent --version "$version" --install-location / "$stage/$arch/component.pkg"
  hostarch="$arch"
  if [[ "$arch" == amd64 ]]; then hostarch=x86_64; fi
  cat > "$stage/$arch/distribution.xml" <<XML
<?xml version="1.0" encoding="utf-8"?>
<installer-gui-script minSpecVersion="1">
<title>Tunnexo</title>
<options customize="never" hostArchitectures="$hostarch" rootVolumeOnly="true"/>
<domains enable_localSystem="true" enable_currentUserHome="false" enable_anywhere="false"/>
<choices-outline><line choice="default"/></choices-outline>
<choice id="default" visible="false"><pkg-ref id="live.tunnexo.agent"/></choice>
<pkg-ref id="live.tunnexo.agent" version="$version">component.pkg</pkg-ref>
</installer-gui-script>
XML
  productbuild --distribution "$stage/$arch/distribution.xml" --package-path "$stage/$arch" "dist/installers/tunnexo_${version}_macos_${arch}.pkg"
done
pkgbuild --nopayload --scripts packaging/macos/uninstall --identifier live.tunnexo.uninstall --version "$version" "dist/installers/tunnexo_${version}_macos_uninstall.pkg"
