$ErrorActionPreference = "Stop"

# ============================================================
# AgentBoard project bootstrap
# Run from the AgentBoard repository root:
#
#   powershell -ExecutionPolicy Bypass -File .\bootstrap-project.ps1
#
# or in PowerShell 7:
#
#   pwsh -File .\bootstrap-project.ps1
# ============================================================

function Write-Step {
    param([string]$Message)

    Write-Host ""
    Write-Host "============================================================"
    Write-Host $Message
    Write-Host "============================================================"
}

function Assert-Command {
    param([string]$Name)

    if (-not (Get-Command $Name -ErrorAction SilentlyContinue)) {
        throw "Required command '$Name' was not found in PATH."
    }
}

function Assert-PathExists {
    param(
        [string]$Path,
        [string]$Description
    )

    if (-not (Test-Path $Path)) {
        throw "$Description not found: $Path"
    }
}

# ------------------------------------------------------------
# 1. Validate environment
# ------------------------------------------------------------

Write-Step "1/10 Validating AgentBoard project"

$Root = (Get-Location).Path

Write-Host "Project root: $Root"

Assert-Command "go"
Assert-Command "git"
Assert-Command "node"
Assert-Command "npm"

Assert-PathExists ".\go.mod" "go.mod"
Assert-PathExists ".\cmd\agentboard\main.go" "Go entry point"
Assert-PathExists ".\web\package.json" "Frontend package.json"
Assert-PathExists ".\web\vite.config.ts" "Vite config"

Write-Host "Go:"
go version

Write-Host "Node:"
node --version

Write-Host "npm:"
npm --version

Write-Host "Git:"
git --version


# ------------------------------------------------------------
# 2. Ensure project directories exist
# ------------------------------------------------------------

Write-Step "2/10 Creating project directories"

$dirs = @(
    "cmd/agentboard"
    "internal/app"
    "internal/domain"
    "internal/workflow"
    "internal/storage"
    "internal/httpapi"
    "internal/mcpserver"
    "internal/history"
    "internal/testing"
    "internal/telemetry"
    "internal/properties"
    "internal/git"
    "internal/webui"
    "internal/webui/dist"
    "migrations"
    "scripts"
    "docs"
)

foreach ($dir in $dirs) {
    New-Item -ItemType Directory -Force -Path $dir | Out-Null
}

Write-Host "Directories ready."


# ------------------------------------------------------------
# 3. Configure embedded frontend
# ------------------------------------------------------------

Write-Step "3/10 Configuring embedded frontend"

@'
package webui

import "embed"

//go:embed all:dist
var FS embed.FS
'@ | Set-Content -Encoding UTF8 .\internal\webui\embed.go

New-Item `
    -ItemType File `
    -Force `
    .\internal\webui\dist\.placeholder |
    Out-Null

Write-Host "internal/webui/embed.go created."


# ------------------------------------------------------------
# 4. Create .gitignore
# ------------------------------------------------------------

Write-Step "4/10 Creating .gitignore"

@'
# Go
/bin/
/dist/
*.exe

# Runtime
*.db
*.db-shm
*.db-wal
*.log

# Local config
.env
.env.*

# Frontend
web/node_modules/
web/dist/

# Embedded generated frontend
internal/webui/dist/*
!internal/webui/dist/.placeholder

# IDE
.idea/
.vscode/

# OS
.DS_Store
Thumbs.db
'@ | Set-Content -Encoding UTF8 .\.gitignore

Write-Host ".gitignore created."


# ------------------------------------------------------------
# 5. Create AGENTS.md
# ------------------------------------------------------------

Write-Step "5/10 Creating AGENTS.md"

@'
# AgentBoard development rules

AgentBoard is a local control plane for human + AI task execution.

## Technology

Backend:
- Go
- database/sql
- modernc.org/sqlite
- chi
- official MCP Go SDK

Frontend:
- React
- TypeScript
- Vite
- TanStack Query
- dnd-kit
- Tailwind CSS

Do not replace the technology stack without explicit instruction.

## Architecture

Business rules belong in the Go core.

HTTP API and MCP are adapters over the same core services.

Do not duplicate workflow/state-transition rules in:
- frontend;
- HTTP handlers;
- MCP handlers.

SQLite is the application state source of truth.

AI access and Human access are separate capability surfaces.

AI must never receive a generic `move_task(target_state)` operation.

Agent workflow transitions must be exposed only through explicitly allowed MCP tools.

## Task workflow

Main workflow:

backlog -> features -> in_progress -> testing -> verification -> complete

Backlog and Complete are human-owned states.

## AI workflow permissions

AI may only perform:

features -> in_progress
in_progress -> testing
testing -> verification

AI must never:
- move tasks backwards;
- move tasks to backlog;
- move tasks to complete;
- read backlog;
- read complete;
- impersonate Human or System in task history.

These restrictions must be enforced by Go Core, not only by prompts.

## Human workflow

Human operates through the Web UI and Human HTTP API.

Human has broader workflow permissions and owns final task acceptance.

AI MCP and Human HTTP API are separate capability surfaces over the same Core.

## Testing

Tasks may use:

- ai
- human
- hybrid

AI must respect the configured testing policy.

A task requiring unresolved human testing must not be moved by AI to verification.

## History

Task history is shared between Human, AI and System.

History is append-oriented.

System transitions must be generated by AgentBoard itself.

AI must not fabricate System events.

## Quality

Do not implement placeholders for requested functionality.

Do not silently skip acceptance criteria.

Before each milestone commit:
- run Go tests;
- run Go build;
- run frontend build.

Keep commits atomic by top-level implementation milestone.

If the entire implementation cannot be completed in one session:
- leave the repository buildable whenever possible;
- commit completed milestones;
- document unfinished work in IMPLEMENTATION_STATUS.md;
- do not claim unfinished functionality is complete.
'@ | Set-Content -Encoding UTF8 .\AGENTS.md

Write-Host "AGENTS.md created."


# ------------------------------------------------------------
# 6. Create README.md
# ------------------------------------------------------------

Write-Step "6/10 Creating README.md"

@'
# AgentBoard

AgentBoard is a local-first task control plane designed for collaboration between a human project owner and AI agents.

## Core idea

Human and AI work with the same tasks and shared task history, but have different workflow capabilities.

Task workflow:

backlog -> features -> in_progress -> testing -> verification -> complete

AI can only move forward through:

features -> in_progress -> testing -> verification

Final acceptance belongs to the human project owner.

## Features planned for MVP

- Persistent local projects
- Kanban task board
- Shared Human / AI / System task history
- Strict AI workflow capabilities
- Human-only final acceptance
- AI / Human / Hybrid testing modes
- Structured test results
- MCP integration
- Human HTTP API
- AI usage telemetry
- Custom task properties
- Task dependencies
- Git metadata integration
- Embedded React Web UI

## Stack

Backend:
- Go
- SQLite
- database/sql
- chi
- official Model Context Protocol Go SDK

Frontend:
- React
- TypeScript
- Vite
- TanStack Query
- dnd-kit
- Tailwind CSS

## Architecture

AgentBoard Core owns:
- workflow rules;
- persistence;
- permissions;
- history;
- testing rules;
- telemetry.

Adapters:
- HTTP API for Human interaction;
- MCP for AI interaction;
- CLI for local administration.

SQLite is stored outside managed workspaces.

Managed projects contain only lightweight AgentBoard project metadata.

## Status

MVP development.
'@ | Set-Content -Encoding UTF8 .\README.md

Write-Host "README.md created."


# ------------------------------------------------------------
# 7. Create build/dev scripts
# ------------------------------------------------------------

Write-Step "7/10 Creating development scripts"

@'
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
'@ | Set-Content -Encoding UTF8 .\scripts\build.ps1


@'
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
'@ | Set-Content -Encoding UTF8 .\scripts\dev.ps1

Write-Host "scripts/build.ps1 created."
Write-Host "scripts/dev.ps1 created."


# ------------------------------------------------------------
# 8. Run full baseline build
# ------------------------------------------------------------

Write-Step "8/10 Running baseline build"

& .\scripts\build.ps1

if ($LASTEXITCODE -ne 0) {
    throw "Baseline build failed."
}

Write-Host "Baseline build successful."


# ------------------------------------------------------------
# 9. Initialize Git repository if necessary
# ------------------------------------------------------------

Write-Step "9/10 Preparing Git repository"

if (-not (Test-Path ".\.git")) {
    Write-Host "Initializing Git repository..."

    git init

    if ($LASTEXITCODE -ne 0) {
        throw "git init failed."
    }
}

git branch -M main

if ($LASTEXITCODE -ne 0) {
    throw "Could not set main branch."
}

Write-Host ""
Write-Host "Git status before baseline commit:"
git status --short

Write-Host ""
Write-Host "Adding project files..."

git add .

if ($LASTEXITCODE -ne 0) {
    throw "git add failed."
}

Write-Host ""
Write-Host "Files staged:"
git status --short


# ------------------------------------------------------------
# Check whether there is anything to commit.
# ------------------------------------------------------------

git diff --cached --quiet
$HasChanges = ($LASTEXITCODE -ne 0)

if ($HasChanges) {
    Write-Host ""
    Write-Host "Creating baseline commit..."

    git commit -m "chore: bootstrap AgentBoard project"

    if ($LASTEXITCODE -ne 0) {
        Write-Host ""
        Write-Host "Git commit failed."
        Write-Host ""
        Write-Host "If Git reports that user.name or user.email is missing,"
        Write-Host "configure them and run:"
        Write-Host ""
        Write-Host '  git config --global user.name "Your Name"'
        Write-Host '  git config --global user.email "you@example.com"'
        Write-Host ""
        Write-Host "Then:"
        Write-Host ""
        Write-Host '  git commit -m "chore: bootstrap AgentBoard project"'
        Write-Host ""
        throw "Baseline commit could not be created."
    }
}
else {
    Write-Host "Nothing new to commit. Baseline commit already exists or working tree is unchanged."
}


# ------------------------------------------------------------
# 10. Create / switch to feat/mvp branch
# ------------------------------------------------------------

Write-Step "10/10 Preparing feat/mvp development branch"

$ExistingBranch = git branch --list "feat/mvp"

if ($ExistingBranch) {
    Write-Host "Branch feat/mvp already exists. Switching to it..."

    git checkout feat/mvp

    if ($LASTEXITCODE -ne 0) {
        throw "Could not switch to feat/mvp."
    }
}
else {
    Write-Host "Creating feat/mvp..."

    git checkout -b feat/mvp

    if ($LASTEXITCODE -ne 0) {
        throw "Could not create feat/mvp."
    }
}


# ------------------------------------------------------------
# Final report
# ------------------------------------------------------------

Write-Host ""
Write-Host "============================================================"
Write-Host "AgentBoard bootstrap completed successfully"
Write-Host "============================================================"
Write-Host ""

Write-Host "Project:"
Write-Host "  $Root"
Write-Host ""

Write-Host "Current branch:"
git branch --show-current

Write-Host ""

Write-Host "Git status:"
git status --short

Write-Host ""

Write-Host "Recent commits:"
git --no-pager log --oneline -5

Write-Host ""

Write-Host "Next step:"
Write-Host "  Give the AgentBoard MVP specification to Codex."
Write-Host ""

Write-Host "Do NOT run go mod tidy manually yet."
Write-Host "Codex should run it after the Go dependencies are actually imported."
Write-Host ""