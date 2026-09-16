# Installer maintainer guide

User instructions live in README.md and TUNNEXO_SETUP.md. Installers use the numeric version in cmd/version.go; increase it before publishing an upgrade. Run the workflow manually or push a matching version tag. Workflow artifacts are build outputs, not automatically published GitHub Releases.

## Build

- macOS: `bash scripts/build-macos-installer.sh` using Apple's pkgbuild/productbuild.
- Windows: install WiX CLI 4.0.6 and WixToolset.UI.wixext 4.0.6, then run `./scripts/build-windows-installer.ps1` in PowerShell.
- Linux: install dpkg-deb and rpmbuild, then `bash scripts/build-linux-installers.sh`.

All require the Go version in go.mod. Outputs go to dist/installers. Windows and Linux packages must be built and tested on their native CI runners; cross-compiling Go alone does not validate installers.

macOS owns /usr/local/lib/tunnexo and /etc/paths.d/tunnexo. Separate uninstaller package removes only the named installed files and matching PATH file, preserving user files. Windows MSI owns Program Files/Tunnexo and its machine PATH entry; Windows Installer handles removal and upgrades. Linux packages own /usr/bin/tunnexo and /usr/share/doc/tunnexo.

No .env, secrets, or backend implementation document is included in installer payloads. Old manual installations are not removed automatically.

## Release checks

- Test fresh install, new-terminal PATH resolution, version/help, upgrade to a higher version, and uninstall on supported OS/architecture.
- Verify Windows PATH cleanup and DEB/RPM package removal.
- macOS installer and uninstaller must be tested in a disposable Mac/VM before publication; building does not install them.
- Sign/notarize macOS packages using your Developer ID Installer credentials; sign Windows MSI with your code-signing certificate. Current scripts produce unsigned artifacts. Credentials and release publishing are not configured.
- Publish installer packages and macOS uninstaller together. Do not advertise Windows/Linux artifacts before their CI build and checks pass.

Reference: https://docs.firegiant.com/wix/schema/wxs/environment/ (MSI PATH component lifecycle).
