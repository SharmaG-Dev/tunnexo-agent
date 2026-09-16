# Tunnexo

Share your local web application with a public URL.

Tunnexo lets you show work in progress, share a demo, or receive webhooks while your application runs on your computer.

**[Website](https://tunnexo.live)** · **[Downloads](https://github.com/SharmaG-Dev/Portune-Agent/releases)** · **[Complete installation guide](TUNNEXO_SETUP.md)**

## What you can do

- Share a local website with teammates or clients.
- Preview a development build on another device.
- Test webhooks against a local application.
- Share a temporary demo without deploying your application.

Your application, Tunnexo, and internet connection must stay running while someone uses the public URL.

## Install

Get the installer for your computer from [Releases](https://github.com/SharmaG-Dev/Portune-Agent/releases).

| Computer | Package |
| --- | --- |
| macOS Apple Silicon / Intel | `.pkg` for ARM64 / AMD64 |
| Windows x64 | `.msi` |
| Ubuntu / Debian | `.deb` for AMD64 / ARM64 |
| Fedora / RHEL-family | `.rpm` for AMD64 / ARM64 |

Double-click the macOS or Windows installer and follow the prompts. On Linux, install the downloaded package using your package manager. Open a new terminal after installation. No Go installation or manual PATH editing is required.

See the [installation guide](TUNNEXO_SETUP.md) for platform-specific steps and installer availability.

## Start sharing

### 1. Start your application

Start your app normally and check its local address in your browser.

### 2. Start Tunnexo in another terminal

```sh
tunnexo agent --target http://localhost:4000
```

Replace `4000` with your app's port. Normal public-service use needs no `.env` or manually configured token. Contact the maintainer if public guest access is unavailable.

### 3. Share your URL

Copy the **Public URL** displayed in the terminal. It will look like:

```text
https://example.tunnexo.live
```

This is an example; use the actual URL shown by your session. Keep both terminals running. Press **Ctrl+C** in the Tunnexo terminal to stop sharing.

## Update and uninstall

Install a newer package to update. To uninstall:

- Windows: Settings → Apps → Installed apps → Tunnexo → Uninstall.
- macOS: open the release's `tunnexo_<version>_macos_uninstall.pkg` and follow the prompts.
- Ubuntu/Debian: `sudo apt remove tunnexo`.
- Fedora/RHEL-family: `sudo dnf remove tunnexo`.

Stop running Tunnexo sessions before updating or uninstalling.

## Help

```sh
tunnexo --help
tunnexo agent --help
tunnexo version
```



## Common problems

| Problem | Try this |
| --- | --- |
| Command not found | Complete installation and reopen your terminal. |
| Permission denied on macOS/Linux | Reinstall using the official package. |
| Connection refused or a 502 response | Check that your local app is running and the target port is correct. |
| Guest access unavailable or connection fails | Check your connection, retry later, or contact the maintainer. |
| Link stops working | Keep your app and Tunnexo running, and prevent your computer from sleeping. Restart Tunnexo and check the new URL. |

More help: [Complete setup guide](TUNNEXO_SETUP.md).

## Support

[Report an issue](https://github.com/SharmaG-Dev/Portune-Agent/issues) with your operating system, Tunnexo version, and error message. Remove private information before posting.

## License

[MIT License](LICENSE) · Copyright (c) 2026 SharmaG-Dev.
