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

    # The frontend is embedded at Go build time. Refresh the executable only
    # after Vite has written internal/webui/dist.
    go build -o ".\agentboard.exe" ./cmd/agentboard

    if ($LASTEXITCODE -ne 0) {
        throw "Could not build agentboard.exe. Stop a running AgentBoard server and retry."
    }

    Write-Host ""
    Write-Host "============================================================"
    Write-Host "AgentBoard build completed successfully."
    Write-Host "============================================================"
}
finally {
    Pop-Location
}
