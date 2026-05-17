#!/usr/bin/env pwsh
# Build a release: tagged snugnas.exe + NSIS installer.
#
#   .\scripts\build-release.ps1 -Version 0.6.2

param(
    [Parameter(Mandatory = $true)]
    [string]$Version
)

$ErrorActionPreference = 'Stop'

$repo = Resolve-Path (Join-Path $PSScriptRoot '..')
Set-Location $repo

Write-Host "Building snugnas.exe v$Version..."
$ldflags = "-s -w -X github.com/Blazzical/snugNAS/internal/buildinfo.Version=$Version"
go build -ldflags $ldflags -o snugnas.exe .\cmd\snugnas
if ($LASTEXITCODE -ne 0) { throw "go build failed" }

$exeSize = (Get-Item snugnas.exe).Length
Write-Host "`nBuilt snugnas.exe ($exeSize bytes)"
$reported = & .\snugnas.exe --version
Write-Host "Reported version: $reported"

$nsis = Get-Command makensis -ErrorAction SilentlyContinue
if ($null -eq $nsis) {
    Write-Warning "makensis not on PATH - skipping installer build."
    Write-Host "Install NSIS with: winget install --id NSIS.NSIS --silent"
    Write-Host "Then add 'C:\Program Files (x86)\NSIS' to your PATH."
    exit 0
}

Write-Host "`nBuilding NSIS installer..."
# /D defines pass APP_VERSION (semver) and APP_VERSION_NUMERIC (4-part) into
# the .nsi so the installer metadata matches the binary's --version output.
& makensis "/DAPP_VERSION=$Version" "/DAPP_VERSION_NUMERIC=$Version.0" ".\installer\snugnas.nsi"
if ($LASTEXITCODE -ne 0) { throw "makensis failed" }

$out = Join-Path $repo "installer\snugnas-setup.exe"
if (Test-Path $out) {
    $instSize = (Get-Item $out).Length
    Write-Host "`nInstaller: $out ($instSize bytes)"
} else {
    throw "Installer not found at expected path"
}
