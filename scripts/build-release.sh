#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
mkdir -p dist
stage=$(mktemp -d)
trap 'rm -rf "$stage"' EXIT
for platform in darwin/arm64 darwin/amd64 linux/amd64 linux/arm64 windows/amd64; do
  os=${platform%/*}
  arch=${platform#*/}
  name="tunnexo_${os}_${arch}"
  package="$stage/$name"
  mkdir -p "$package"
  executable=tunnexo
  if [[ "$os" == windows ]]; then executable=tunnexo.exe; fi
  echo "Building $name"
  CGO_ENABLED=0 GOOS="$os" GOARCH="$arch" go build -trimpath -ldflags="-s -w" -o "$package/$executable" .
  cp LICENSE README.md TUNNEXO_SETUP.md SERVER_GUEST_INTEGRATION.md .env.example "$package/"
  case "$platform" in
    darwin/arm64) folder=macos-arm64 ;;
    darwin/amd64) folder=macos-intel ;;
    linux/amd64) folder=linux-amd64 ;;
    linux/arm64) folder=linux-arm64 ;;
    windows/amd64) folder=windows ;;
  esac
  mkdir -p "dist/$folder"
  cp "$package/$executable" "dist/$folder/$executable"
  if [[ "$os" == linux ]]; then
    tar -czf "$stage/$name.tar.gz" -C "$package" .
    mv "$stage/$name.tar.gz" "dist/$name.tar.gz"
  else
    (cd "$package" && zip -q "$stage/$name.zip" "$executable" LICENSE README.md TUNNEXO_SETUP.md SERVER_GUEST_INTEGRATION.md .env.example)
    mv "$stage/$name.zip" "dist/$name.zip"
  fi
done
