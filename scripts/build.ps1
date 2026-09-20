$ErrorActionPreference = "Stop"

$Root = Split-Path -Parent $PSScriptRoot

Push-Location $Root

try {
    Write-Host ""
    Write-Host "==> Building frontend"

    Push-Location ".\web"

    try {
        npm run build

        if ($LASTEXITCODE -ne 0) {
            throw "Frontend build failed."
        }
    }
    finally {
        Pop-Location
    }

    # Vite uses emptyOutDir, therefore recreate the tracked placeholder
    # after every production frontend build.
    New-Item `
        -ItemType File `
        -Force `
        ".\internal\webui\dist\.placeholder" |
        Out-Null

    Write-Host ""
    Write-Host "==> Running Go tests"

    go test ./...

    if ($LASTEXITCODE -ne 0) {
        throw "Go tests failed."
    }

    Write-Host ""
    Write-Host "==> Building Go application"

    go build ./...

    if ($LASTEXITCODE -ne 0) {
        throw "Go build failed."
    }

    Write-Host ""
    Write-Host "============================================================"
    Write-Host "AgentBoard build completed successfully."
    Write-Host "============================================================"
}
finally {
    Pop-Location
}
