$ErrorActionPreference = 'Stop'
Set-Location (Join-Path $PSScriptRoot '..')
$version = [regex]::Match((Get-Content cmd/version.go -Raw), 'const version = "([0-9]+\.[0-9]+\.[0-9]+)"').Groups[1].Value
if (-not $version) { throw 'Missing numeric version' }
New-Item -ItemType Directory -Force dist/windows, dist/installers | Out-Null
$env:CGO_ENABLED = '0'
$env:GOOS = 'windows'
$env:GOARCH = 'amd64'
go build -trimpath -ldflags='-s -w' -o dist/windows/tunnexo.exe .
if ($LASTEXITCODE -ne 0) { throw 'Go build failed' }
wix build packaging/windows/Package.wxs -arch x64 -d "Version=$version" -ext WixToolset.UI.wixext -o "dist/installers/tunnexo_${version}_windows_amd64.msi"
if ($LASTEXITCODE -ne 0) { throw 'MSI build failed' }
