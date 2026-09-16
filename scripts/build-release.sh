#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
case "$(uname -s)" in
  Darwin) exec bash scripts/build-macos-installer.sh ;;
  Linux) exec bash scripts/build-linux-installers.sh ;;
  *) echo 'On Windows, run scripts/build-windows-installer.ps1 in PowerShell.' >&2; exit 1 ;;
esac
