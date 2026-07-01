#Requires -Version 5.1

$ErrorActionPreference = "Stop"

$Repo    = "tosterabgx/marten"
$BinName = "marten.exe"

$arch = (Get-CimInstance Win32_Processor).Architecture
if ($arch -ne 9) {
    Write-Error "Unsupported architecture. Only Windows x64 (amd64) is supported. Download manually from: https://github.com/$Repo/releases/latest"
    exit 1
}

$InstallDir = Join-Path $env:LOCALAPPDATA "marten"

if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
}

$DownloadUrl = "https://github.com/$Repo/releases/latest/download/marten-windows-amd64.exe"
$TempFile    = Join-Path $env:TEMP "marten-install-$([System.Guid]::NewGuid()).exe"
 
Write-Host "Downloading marten (windows/amd64)..."
Write-Host "  -> $DownloadUrl"

try {
    Invoke-WebRequest -Uri $DownloadUrl -OutFile $TempFile -UseBasicParsing
} catch {
    Write-Error "Download failed: $_"
    exit 1
}

$Dest = Join-Path $InstallDir $BinName
 
Move-Item -Path $TempFile -Destination $Dest -Force
 
Write-Host "marten installed to $Dest"

$UserPath = [Environment]::GetEnvironmentVariable("PATH", "User")
 
if ($UserPath -notlike "*$InstallDir*") {
    [Environment]::SetEnvironmentVariable(
        "PATH",
        "$UserPath;$InstallDir",
        "User"
    )
    Write-Host ""
    Write-Host "Added $InstallDir to your PATH."
    Write-Host "Restart your terminal for the change to take effect."
} else {
    Write-Host "$InstallDir is already in your PATH."
}
 

Write-Host ""
Write-Host "Try it: marten tcp 3000"
