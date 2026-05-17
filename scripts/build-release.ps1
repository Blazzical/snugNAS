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

Write-Host "`nBuilt snugnas.exe ($((Get-Item snugnas.exe).Length) bytes)"
Write-Host "Reported version: $(& .\snugnas.exe --version)"

$nsis = Get-Command makensis -ErrorAction SilentlyContinue
if ($null -eq $nsis) {
    Write-Warning "makensis not on PATH — skipping installer build."
    Write-Host "Install NSIS with: winget install --id NSIS.NSIS --silent"
    exit 0
}

Write-Host "`nBuilding NSIS installer..."
& makensis ".\installer\snugnas.nsi"
if ($LASTEXITCODE -ne 0) { throw "makensis failed" }

$out = Join-Path $repo "installer\snugnas-setup.exe"
if (Test-Path $out) {
    Write-Host "`nInstaller: $out ($((Get-Item $out).Length) bytes)"
} else {
    throw "Installer not found at expected path"
}
