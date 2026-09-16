# Tunnexo — Install aur Use Guide

Apni local website ko public link se share karein. Yeh guide macOS, Linux aur Windows users ke liye hai.

**[Website](https://tunnexo.live)** · **[Download](https://github.com/SharmaG-Dev/Portune-Agent/releases)** · **[Support](https://github.com/SharmaG-Dev/Portune-Agent/issues)**

## Tunnexo kab helpful hai?

- Client ya teammate ko apni website ka live demo dikhana.
- Development build ko doosre device par preview karna.
- Local application par webhooks test karna.
- Application deploy kiye bina temporary preview share karna.

Example: aapki app `http://localhost:4000` par chalti hai. Tunnexo chalane par terminal mein share karne ke liye `https://example.tunnexo.live` jaisa URL milta hai.

Actual URL terminal se copy karein. App, Tunnexo aur internet connection sharing ke dauraan running rehne chahiye.

## Setup se pehle

1. Internet-connected computer.
2. Aapke computer ke liye Tunnexo download.
3. Running local application aur uska address/port.
4. macOS/Linux par Terminal, Windows par PowerShell.

Downloaded executable ke liye Go install nahi karna hai. Normal public-service use mein `.env` ya manual token setup nahi chahiye. Service unavailable ho to maintainer se contact karein.

## 1. Sahi download choose karein

[GitHub Releases](https://github.com/SharmaG-Dev/Portune-Agent/releases) kholein aur matching archive download karein:

| Computer | File |
| --- | --- |
| macOS Apple Silicon — M-series | `tunnexo_darwin_arm64.zip` |
| macOS Intel | `tunnexo_darwin_amd64.zip` |
| Linux x86-64 / AMD64 | `tunnexo_linux_amd64.tar.gz` |
| Linux ARM64 | `tunnexo_linux_arm64.tar.gz` |
| Windows x64 | `tunnexo_windows_amd64.zip` |

Archive Releases mein available nahi hai to maintainer se matching file lein. Windows ARM64 aur 32-bit ke liye native download is list mein nahi hai.

Neeche sirf apne operating system ka installation section follow karein, phir section 5 se continue karein.

## 2. macOS installation

### Processor check

Apple menu → About This Mac mein chip dekhein. M-series ke liye Apple Silicon download aur Intel ke liye Intel download use karein.

Terminal se bhi check kar sakte hain:

```sh
uname -m
```

`arm64` Apple Silicon aur `x86_64` Intel architecture indicate karta hai. Translated terminal use ho raha ho to About This Mac se confirm karein.

### Extract karein

ZIP ko `Downloads` mein rakhein. Apple Silicon ke liye:

```sh
mkdir -p "$HOME/Downloads/tunnexo"
unzip "$HOME/Downloads/tunnexo_darwin_arm64.zip" -d "$HOME/Downloads/tunnexo"
cd "$HOME/Downloads/tunnexo"
```

Intel ke liye upar ke commands ki jagah:

```sh
mkdir -p "$HOME/Downloads/tunnexo"
unzip "$HOME/Downloads/tunnexo_darwin_amd64.zip" -d "$HOME/Downloads/tunnexo"
cd "$HOME/Downloads/tunnexo"
```

`mkdir` folder banata hai, `unzip` archive extract karta hai aur `cd` terminal ko us folder mein le jaata hai. Browser ne pehle hi extract kiya ho to actual extracted folder mein terminal kholein.

### Install karein

```sh
chmod +x ./tunnexo
./tunnexo --help
mkdir -p "$HOME/.local/bin"
cp ./tunnexo "$HOME/.local/bin/tunnexo"
export PATH="$HOME/.local/bin:$PATH"
tunnexo version
```

| Command | Kaam |
| --- | --- |
| `chmod +x` | File ko executable permission deta hai. |
| `./tunnexo --help` | Downloaded executable ki help dikhata hai. |
| `mkdir -p` | User installation folder banata hai. |
| `cp` | Executable ko installation folder mein copy karta hai. |
| `export PATH=...` | Current terminal mein `tunnexo` command available karta hai. |
| `tunnexo version` | Installed version dikhata hai. |

Future terminals ke liye default zsh configuration kholein:

```sh
nano "$HOME/.zshrc"
```

Yeh line ek baar add karein:

```sh
export PATH="$HOME/.local/bin:$PATH"
```

`Ctrl+O`, Enter se save aur `Ctrl+X` se exit karein. Phir:

```sh
source "$HOME/.zshrc"
command -v tunnexo
```

`source` settings reload karta hai. `command -v` installed command ka path dikhata hai. Bash use karte hain to apni Bash startup file mein PATH line add karein.

macOS downloaded app block kare to source verify karke System Settings → Privacy & Security mein app ka approval option check karein.

## 3. Linux installation

### Architecture check

```sh
uname -m
```

`x86_64` ke liye AMD64 archive; `aarch64`/`arm64` ke liye ARM64 archive choose karein.

### Extract karein

Archive ko `Downloads` folder mein rakhein. AMD64 ke liye:

```sh
mkdir -p "$HOME/Downloads/tunnexo"
tar -xzf "$HOME/Downloads/tunnexo_linux_amd64.tar.gz" -C "$HOME/Downloads/tunnexo"
cd "$HOME/Downloads/tunnexo"
```

ARM64 ke liye upar ke commands ki jagah:

```sh
mkdir -p "$HOME/Downloads/tunnexo"
tar -xzf "$HOME/Downloads/tunnexo_linux_arm64.tar.gz" -C "$HOME/Downloads/tunnexo"
cd "$HOME/Downloads/tunnexo"
```

`tar -xzf` archive extract karta hai; `-C` destination folder select karta hai. Download path alag ho to command mein actual path use karein.

### Install karein

```sh
chmod +x ./tunnexo
./tunnexo --help
mkdir -p "$HOME/.local/bin"
cp ./tunnexo "$HOME/.local/bin/tunnexo"
export PATH="$HOME/.local/bin:$PATH"
tunnexo version
```

Yeh executable permission set karta hai, binary user folder mein install karta hai aur current terminal ka PATH update karta hai. `sudo` ki zaroorat nahi hai.

Future Bash terminals ke liye:

```sh
nano "$HOME/.bashrc"
```

Yeh line ek baar add karein:

```sh
export PATH="$HOME/.local/bin:$PATH"
```

Save/exit ke baad:

```sh
source "$HOME/.bashrc"
command -v tunnexo
```

Zsh use karte hain to `.bashrc` ki jagah `.zshrc` use karein.

## 4. Windows installation

Windows Settings → System → About mein x64 system confirm karein. `tunnexo_windows_amd64.zip` download karke `Downloads` mein rakhein.

**PowerShell** kholein. Yeh commands Command Prompt ke liye nahi hain.

### Extract aur install karein

```powershell
$tunnexoZip = Join-Path $HOME 'Downloads\tunnexo_windows_amd64.zip'
$tunnexoExtract = Join-Path $HOME 'Downloads\tunnexo'
Expand-Archive -LiteralPath $tunnexoZip -DestinationPath $tunnexoExtract -Force
Set-Location $tunnexoExtract
.\tunnexo.exe --help

$tunnexoInstall = Join-Path $env:LOCALAPPDATA 'Tunnexo'
New-Item -ItemType Directory -Path $tunnexoInstall -Force | Out-Null
Copy-Item '.\tunnexo.exe' -Destination $tunnexoInstall -Force
```

| Command | Kaam |
| --- | --- |
| `Join-Path` | User folder ka correct file path banata hai. |
| `Expand-Archive` | ZIP extract karta hai; matching extracted files replace ho sakti hain. |
| `Set-Location` | Terminal ka current folder change karta hai. |
| `.\tunnexo.exe --help` | Extracted executable ki help dikhata hai. |
| `New-Item` | Installation folder banata hai. |
| `Copy-Item` | Binary install karta hai; previous binary replace ho sakti hai. |

### Command ko PATH mein add karein

Isi PowerShell window mein:

```powershell
$tunnexoUserPath = [Environment]::GetEnvironmentVariable('Path', 'User')
if (($tunnexoUserPath -split ';') -notcontains $tunnexoInstall) {
    $tunnexoNewPath = (@($tunnexoUserPath, $tunnexoInstall) | Where-Object { $_ }) -join ';'
    [Environment]::SetEnvironmentVariable('Path', $tunnexoNewPath, 'User')
}
$env:Path = "$tunnexoInstall;$env:Path"
Get-Command tunnexo
tunnexo version
```

Yeh existing user PATH preserve karke installation folder add karta hai. Current terminal mein command turant available hogi. `Get-Command` executable ka location aur `tunnexo version` version dikhata hai.

Future sessions ke liye terminal application close karke dobara kholein. Zaroorat pade to sign out/in karein. Administrator PowerShell ki zaroorat nahi hai.

## 5. Apni local application start karein

Pehle terminal mein app ka normal start command chalayein. Browser mein app ka local address kholkar confirm karein ki woh kaam kar rahi hai.

Example:

```text
http://localhost:4000
```

Aapki app `3000`, `5173`, `8080` ya kisi aur port par ho sakti hai. Aage ke commands mein wahi port use karein. Tunnexo aapki local app start nahi karta.

Optional check — macOS/Linux:

```sh
curl -i http://localhost:4000
```

Windows PowerShell:

```powershell
Test-NetConnection -ComputerName localhost -Port 4000
Invoke-WebRequest -Uri 'http://localhost:4000'
```

Yeh local app ki reachability/response check karte hain. Connection refused aaye to app start karein ya port correct karein. Root page 404 de to app ki known working route check karein.

## 6. Public link banayein

Doosra terminal kholein aur run karein. Yeh command macOS, Linux aur Windows PowerShell par same hai:

```sh
tunnexo agent --target http://localhost:4000
```

| Part | Meaning |
| --- | --- |
| `tunnexo` | Installed application chalata hai. |
| `agent` | Sharing session start karta hai. |
| `--target` | Aapki local app ka address specify karta hai. |
| `http://localhost:4000` | Aapki running app ka example address. |

Successful connection ke baad terminal mein **Public URL** milega, jaise:

```text
https://example.tunnexo.live
```

Wahi actual URL copy karke share karein. Example URL ko use na karein. Local app aur Tunnexo dono running rakhein; computer sleep hone par link unavailable ho sakta hai.

Aapki app public link se accessible hogi, isliye private app ke required login/access controls enabled rakhein.

### Doosre ports ke examples

Apni app ke liye inmein se matching command use karein:

```sh
tunnexo agent --target http://localhost:3000
tunnexo agent --target http://localhost:5173
tunnexo agent --target http://127.0.0.1:8080
```

Address mein `http://` ya `https://` include karein. `--target` mein username/password, query string ya fragment na dein.

### PATH installation skip ki ho to

Extracted folder se directly run karein.

macOS/Linux:

```sh
./tunnexo agent --target http://localhost:4000
```

Windows PowerShell:

```powershell
.\tunnexo.exe agent --target http://localhost:4000
```

`./` aur `.\` current folder ki executable run karte hain.

## 7. Stop aur restart

Tunnexo terminal mein **Ctrl+C** press karein. Sharing stop hogi. Local application ko stop karna ho to uske alag terminal mein bhi Ctrl+C karein.

Restart:

```sh
tunnexo agent --target http://localhost:4000
```

Restart ke baad public URL dobara check karein; same URL milna guaranteed nahi hai.

## 8. Help aur version commands

| Command | Result |
| --- | --- |
| `tunnexo --help` | Available commands. |
| `tunnexo agent --help` | Sharing command ka usage. |
| `tunnexo help agent` | Agent help ka alternate command. |
| `tunnexo version` | Installed version. |
| `tunnexo completion --help` | Shell completion help. |

PATH installation se pehle command mein `tunnexo` ki jagah `./tunnexo` ya `.\tunnexo.exe` use karein.

## 9. Daily use

1. Local app start karein.
2. Doosra terminal kholein.
3. `tunnexo agent --target http://localhost:4000` run karein, apne port ke saath.
4. Terminal ka public URL share karein.
5. Kaam complete hone par Ctrl+C karein.

Installation roz repeat nahi karni hai.

## 10. Update kaise karein?

Running Tunnexo ko stop karein. [Releases](https://github.com/SharmaG-Dev/Portune-Agent/releases) se apne platform ka naya archive download/extract karein. **Naye extracted folder** mein terminal kholkar executable replace karein.

macOS/Linux:

```sh
chmod +x ./tunnexo
cp ./tunnexo "$HOME/.local/bin/tunnexo"
tunnexo version
```

Windows PowerShell:

```powershell
Copy-Item '.\tunnexo.exe' -Destination "$env:LOCALAPPDATA\Tunnexo\tunnexo.exe" -Force
tunnexo version
```

Phir normal run command se Tunnexo start karein. Automatic update command available nahi hai.

## 11. Troubleshooting

| Problem | Solution |
| --- | --- |
| Command not found / not recognized | PATH installation complete karein aur terminal reopen karein, ya extracted folder se executable run karein. |
| Permission denied | macOS/Linux par `chmod +x ./tunnexo` chalayein. |
| Incompatible executable / exec format error | Apne OS aur processor ka correct archive download karein. |
| Download missing / 404 | Releases page check karein ya maintainer se archive lein. |
| Guest registration unavailable | Public service ke maintainer se contact karein; manual setup ki zaroorat assume na karein. |
| Guest registration limit reached | Thodi der baad retry karein. |
| Connection / authentication failure | Internet check karein, latest release use karein aur error persist ho to support ko message bhejein. |
| Public URL rejected | Latest release check karein; exact error support ko bhejein. |
| Cannot detect local IP | Active network connection check karein. |
| 502 / connection refused | Local app running ho aur `--target` ka host/port correct ho. |
| App public hostname reject karti hai | App ki allowed-host settings mein assigned hostname allow karein. |
| Link sleep/disconnection ke baad band | Computer awake aur connected rakhein; Tunnexo restart karke fresh URL check karein. |
| HTTPS certificate error | Error maintainer ko report karein. |

## 12. Support

[GitHub Issues](https://github.com/SharmaG-Dev/Portune-Agent/issues) par yeh details bhejein:

- Operating system aur processor.
- `tunnexo version` ka output.
- Run command aur error message.
- Local application browser mein chal rahi hai ya nahi.

Private credentials ya personal information share na karein.

## License

Tunnexo [MIT License](LICENSE) ke under available hai.
