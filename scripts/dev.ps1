$ErrorActionPreference = "Stop"

$Root = Split-Path -Parent $PSScriptRoot

Write-Host ""
Write-Host "AgentBoard development"
Write-Host "Project: $Root"
Write-Host ""
Write-Host "Backend:"
Write-Host "  go run ./cmd/agentboard serve"
Write-Host ""
Write-Host "Frontend:"
Write-Host "  cd web"
Write-Host "  npm run dev"
Write-Host ""
Write-Host "Full verification:"
Write-Host "  .\scripts\build.ps1"
Write-Host ""
