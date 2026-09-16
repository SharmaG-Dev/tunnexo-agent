#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
command -v dpkg-deb >/dev/null
command -v rpmbuild >/dev/null
command -v rpm >/dev/null
version=$(sed -n 's/^const version = "\([^"]*\)"/\1/p' cmd/version.go)
[[ "$version" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]
stage=$(mktemp -d)
trap 'rm -rf "$stage"' EXIT
mkdir -p dist/installers
for arch in amd64 arm64; do
  root="$stage/$arch"
  mkdir -p "$root/usr/bin" "$root/usr/share/doc/tunnexo" "$root/DEBIAN"
  CGO_ENABLED=0 GOOS=linux GOARCH="$arch" go build -trimpath -ldflags='-s -w' -o "$root/usr/bin/tunnexo" .
  cp LICENSE README.md TUNNEXO_SETUP.md "$root/usr/share/doc/tunnexo/"
  cat > "$root/DEBIAN/control" <<CONTROL
Package: tunnexo
Version: $version
Architecture: $arch
Maintainer: SharmaG-Dev <SharmaG-Dev@users.noreply.github.com>
Section: net
Priority: optional
Depends: ca-certificates
Homepage: https://tunnexo.live
Description: Share local web applications with a public URL
CONTROL
  dpkg-deb --root-owner-group --build "$root" "dist/installers/tunnexo_${version}_linux_${arch}.deb"
  rpmarch=x86_64
  if [[ "$arch" == arm64 ]]; then rpmarch=aarch64; fi
  top="$stage/rpm-$arch"
  mkdir -p "$top/BUILD" "$top/RPMS" "$top/SOURCES" "$top/SPECS" "$top/SRPMS"
  cat > "$top/SPECS/tunnexo.spec" <<SPEC
Name: tunnexo
Version: $version
Release: 1
Summary: Share local web applications with a public URL
License: MIT
URL: https://tunnexo.live
Requires: ca-certificates
AutoReqProv: no
%description
Share local web applications with a public URL.
%install
mkdir -p %{buildroot}/usr/bin %{buildroot}/usr/share/doc/tunnexo
cp "$root/usr/bin/tunnexo" %{buildroot}/usr/bin/tunnexo
cp "$root/usr/share/doc/tunnexo/"* %{buildroot}/usr/share/doc/tunnexo/
%files
/usr/bin/tunnexo
/usr/share/doc/tunnexo
SPEC
  # GOARCH selects the binary; --target selects RPM metadata. BuildArch would
  # additionally require host compatibility and reject ARM64 on an x86 runner.
  rpmbuild --define "_topdir $top" --define '__os_install_post %{nil}' --target "$rpmarch" -bb "$top/SPECS/tunnexo.spec"
  rpmfile="$top/RPMS/$rpmarch/tunnexo-$version-1.$rpmarch.rpm"
  actual_arch=$(rpm -qp --queryformat '%{ARCH}' "$rpmfile")
  if [[ "$actual_arch" != "$rpmarch" ]]; then
    echo "RPM architecture mismatch: expected $rpmarch, got $actual_arch" >&2
    exit 1
  fi
  cp "$rpmfile" "dist/installers/tunnexo_${version}_linux_${arch}.rpm"
done
