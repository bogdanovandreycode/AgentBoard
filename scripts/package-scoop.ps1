param(
    [Parameter(Mandatory = $true)][ValidatePattern('^\d+\.\d+\.\d+$')][string]$Version,
    [string]$OutputDir = "release"
)
$ErrorActionPreference = "Stop"
$root = Split-Path -Parent $PSScriptRoot
$output = Join-Path $root $OutputDir
$stage = Join-Path $output "stage"
$outputFull = [System.IO.Path]::GetFullPath($output)
$stageFull = [System.IO.Path]::GetFullPath($stage)
if (-not $stageFull.StartsWith($outputFull.TrimEnd([System.IO.Path]::DirectorySeparatorChar) + [System.IO.Path]::DirectorySeparatorChar, [System.StringComparison]::OrdinalIgnoreCase)) {
    throw "Staging directory is outside the release directory: $stageFull"
}
$zipName = "agentboard-$Version-windows-amd64.zip"
$zipPath = Join-Path $output $zipName
Push-Location $root
try {
    New-Item -ItemType Directory -Force -Path $output | Out-Null
    if (Test-Path -LiteralPath $stageFull) { Remove-Item -LiteralPath $stageFull -Recurse -Force }
    New-Item -ItemType Directory -Force -Path $stage | Out-Null
    Push-Location (Join-Path $root "web")
    try {
        npm ci
        if ($LASTEXITCODE -ne 0) { throw "npm ci failed" }
        npm run build
        if ($LASTEXITCODE -ne 0) { throw "Frontend build failed" }
    } finally { Pop-Location }
    go test ./...
    if ($LASTEXITCODE -ne 0) { throw "Go tests failed" }
    go build -ldflags "-X main.version=$Version" -o (Join-Path $stage "agentboard.exe") ./cmd/agentboard
    if ($LASTEXITCODE -ne 0) { throw "Go build failed" }
    $reported = & (Join-Path $stage "agentboard.exe") version
    if ($reported -ne "AgentBoard $Version") { throw "Unexpected binary version: $reported" }
    Copy-Item -LiteralPath (Join-Path $root "README.md") -Destination $stage
    Copy-Item -LiteralPath (Join-Path $root "doc") -Destination $stage -Recurse
    if (Test-Path -LiteralPath $zipPath) { Remove-Item -LiteralPath $zipPath }
    Compress-Archive -Path (Join-Path $stage "*") -DestinationPath $zipPath
    $hash = (Get-FileHash -LiteralPath $zipPath -Algorithm SHA256).Hash.ToLowerInvariant()
    $url = "https://github.com/bogdanovandreycode/AgentBoard/releases/download/v$Version/$zipName"
    $manifest = [ordered]@{
        version = $Version
        description = "Local task board for human and AI workers with a scoped MCP server"
        homepage = "https://github.com/bogdanovandreycode/AgentBoard"
        license = "Unknown"
        architecture = @{ "64bit" = @{ url = $url; hash = $hash } }
        bin = "agentboard.exe"
        shortcuts = @(, @("agentboard.exe", "AgentBoard", "open"))
        notes = "Run 'agentboard init' in a project folder, then 'agentboard open'. Data is stored in your user configuration directory."
        checkver = "github"
        autoupdate = @{ architecture = @{ "64bit" = @{ url = "https://github.com/bogdanovandreycode/AgentBoard/releases/download/v`$version/agentboard-`$version-windows-amd64.zip" } } }
    }
    $manifestPath = Join-Path $output "agentboard.json"
    $manifest | ConvertTo-Json -Depth 10 | Set-Content -LiteralPath $manifestPath -Encoding utf8
    Write-Host "Package: $zipPath"
    Write-Host "SHA256: $hash"
    Write-Host "Manifest: $manifestPath"
} finally {
    Pop-Location
}
