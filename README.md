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

## Download

Download and extract the archive for your computer from [GitHub Releases](https://github.com/SharmaG-Dev/Portune-Agent/releases).

| Computer | Archive |
| --- | --- |
| macOS — Apple Silicon (M-series) | `tunnexo_darwin_arm64.zip` |
| macOS — Intel | `tunnexo_darwin_amd64.zip` |
| Linux — x86-64 | `tunnexo_linux_amd64.tar.gz` |
| Linux — ARM64 | `tunnexo_linux_arm64.tar.gz` |
| Windows — x64 | `tunnexo_windows_amd64.zip` |

Go is not required to use a downloaded executable. If your archive is not listed in Releases, request it from the maintainer.

## Start in a few steps

### 1. Start your application

Start your app normally and open its local address in your browser. The examples below use `http://localhost:4000`; replace this with your app's actual address.

### 2. Open a second terminal

Open a terminal in the extracted Tunnexo folder.

**macOS / Linux:**

```sh
chmod +x ./tunnexo
./tunnexo agent --target http://localhost:4000
```

**Windows PowerShell:**

```powershell
.\tunnexo.exe agent --target http://localhost:4000
```

Normal public-service use does not require a `.env` file or a manually configured token. The public service must be available; if guest access is unavailable, contact the maintainer.

### 3. Share your URL

Copy the **Public URL** displayed in the terminal. It will look like:

```text
https://example.tunnexo.live
```

This is an example; use the actual URL shown by your session. Keep both terminals running. Press **Ctrl+C** in the Tunnexo terminal to stop sharing.

## Run as `tunnexo`

Install the executable on your PATH to use it from any folder:

```sh
tunnexo agent --target http://localhost:4000
```

Follow the [macOS, Linux, or Windows installation steps](TUNNEXO_SETUP.md) for your computer.

## Help

```sh
tunnexo --help
tunnexo agent --help
tunnexo version
```

Before installing on PATH, use `./tunnexo` on macOS/Linux or `.\tunnexo.exe` in Windows PowerShell.

## Common problems

| Problem | Try this |
| --- | --- |
| Command not found | Follow the PATH installation steps, or run the executable from its extracted folder. |
| Permission denied on macOS/Linux | Run `chmod +x ./tunnexo`. |
| Connection refused or a 502 response | Check that your local app is running and the target port is correct. |
| Guest access unavailable or connection fails | Check your connection, retry later, or contact the maintainer. |
| Link stops working | Keep your app and Tunnexo running, and prevent your computer from sleeping. Restart Tunnexo and check the new URL. |

More help: [Complete setup guide](TUNNEXO_SETUP.md).

## Support

[Report an issue](https://github.com/SharmaG-Dev/Portune-Agent/issues) with your operating system, Tunnexo version, and error message. Remove private information before posting.

## License

[MIT License](LICENSE) · Copyright (c) 2026 SharmaG-Dev.
