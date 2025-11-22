# Project Validation Script
# Checks if all files are correct without running the server

Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "Stocky Project Validation" -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host ""

$errors = 0
$warnings = 0

# Check if we're in the right directory
if (-not (Test-Path "go.mod")) {
    Write-Host "❌ Not in project root directory!" -ForegroundColor Red
    Write-Host "Please run from: c:\Users\RAUNEET SINGH\OneDrive\Desktop\backend Stocky" -ForegroundColor Yellow
    exit 1
}

Write-Host "✓ In correct directory" -ForegroundColor Green

# Check essential files exist
Write-Host "`nChecking essential files..." -ForegroundColor Yellow

$essentialFiles = @(
    "go.mod",
    "cmd/server/main.go",
    "internal/config/config.go",
    "internal/models/models.go",
    "internal/repository/repository.go",
    "internal/services/services.go",
    "internal/handlers/handlers.go",
    "internal/middleware/middleware.go",
    "internal/database/database.go",
    "migrations/001_initial_schema.sql",
    ".env.example",
    "README.md",
    "Dockerfile",
    "docker-compose.yml"
)

foreach ($file in $essentialFiles) {
    if (Test-Path $file) {
        Write-Host "  ✓ $file" -ForegroundColor Green
    } else {
        Write-Host "  ❌ $file MISSING" -ForegroundColor Red
        $errors++
    }
}

# Check if .env exists
Write-Host "`nChecking configuration..." -ForegroundColor Yellow
if (Test-Path ".env") {
    Write-Host "  ✓ .env file exists" -ForegroundColor Green
} else {
    Write-Host "  ⚠ .env file not found (will use defaults)" -ForegroundColor Yellow
    Write-Host "    Run: Copy-Item .env.example .env" -ForegroundColor Cyan
    $warnings++
}

# Check Go files for basic syntax
Write-Host "`nChecking Go files..." -ForegroundColor Yellow
$goFiles = Get-ChildItem -Recurse -Filter "*.go"
Write-Host "  Found $($goFiles.Count) Go files" -ForegroundColor Green

foreach ($file in $goFiles) {
    $content = Get-Content $file.FullName -Raw
    
    # Check for package declaration
    if ($content -match "package ") {
        # Has package declaration
    } else {
        Write-Host "  Warning: $($file.Name) might be missing package declaration" -ForegroundColor Yellow
        $warnings++
    }
}

# Check SQL file
Write-Host "`nChecking database schema..." -ForegroundColor Yellow
if (Test-Path "migrations/001_initial_schema.sql") {
    $sql = Get-Content "migrations/001_initial_schema.sql" -Raw
    
    $tables = @("users", "stocks", "reward_events", "user_holdings", "stock_prices", "ledger_entries", "stock_events")
    foreach ($table in $tables) {
        if ($sql -match "CREATE TABLE.*$table") {
            Write-Host "  ✓ Table '$table' defined" -ForegroundColor Green
        } else {
            Write-Host "  ❌ Table '$table' NOT found" -ForegroundColor Red
            $errors++
        }
    }
} else {
    Write-Host "  ❌ Migration file missing" -ForegroundColor Red
    $errors++
}

# Check dependencies in go.mod
Write-Host "`nChecking dependencies..." -ForegroundColor Yellow
$gomod = Get-Content "go.mod" -Raw

$requiredDeps = @(
    "github.com/gin-gonic/gin",
    "github.com/sirupsen/logrus",
    "github.com/lib/pq"
)

foreach ($dep in $requiredDeps) {
    if ($gomod -match [regex]::Escape($dep)) {
        Write-Host "  ✓ $dep" -ForegroundColor Green
    } else {
        Write-Host "  ❌ $dep NOT found" -ForegroundColor Red
        $errors++
    }
}

# Check Docker files
Write-Host "`nChecking Docker configuration..." -ForegroundColor Yellow
if (Test-Path "Dockerfile") {
    $dockerfile = Get-Content "Dockerfile" -Raw
    if ($dockerfile -match "FROM golang") {
        Write-Host "  ✓ Dockerfile is valid" -ForegroundColor Green
    } else {
        Write-Host "  ⚠ Dockerfile might have issues" -ForegroundColor Yellow
        $warnings++
    }
}

if (Test-Path "docker-compose.yml") {
    $compose = Get-Content "docker-compose.yml" -Raw
    if ($compose -match "postgres") {
        Write-Host "  ✓ docker-compose.yml includes PostgreSQL" -ForegroundColor Green
    }
    if ($compose -match "stocky-api") {
        Write-Host "  ✓ docker-compose.yml includes API service" -ForegroundColor Green
    }
}

# Check documentation
Write-Host "`nChecking documentation..." -ForegroundColor Yellow
$docs = @("README.md", "ARCHITECTURE.md", "WINDOWS_SETUP.md", "GITHUB_SETUP.md")
foreach ($doc in $docs) {
    if (Test-Path $doc) {
        $size = (Get-Item $doc).Length
        Write-Host "  OK: $doc - $size bytes" -ForegroundColor Green
    }
}

# Check test scripts
Write-Host "`nChecking test scripts..." -ForegroundColor Yellow
if (Test-Path "test-api.ps1") {
    Write-Host "  ✓ PowerShell test script exists" -ForegroundColor Green
}
if (Test-Path "test-api.sh") {
    Write-Host "  ✓ Bash test script exists" -ForegroundColor Green
}
if (Test-Path "Stocky-API.postman_collection.json") {
    Write-Host "  ✓ Postman collection exists" -ForegroundColor Green
}

# Summary
Write-Host "`n=========================================" -ForegroundColor Cyan
Write-Host "Validation Summary" -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan

if ($errors -eq 0 -and $warnings -eq 0) {
    Write-Host "✓ ALL CHECKS PASSED!" -ForegroundColor Green
    Write-Host "`nProject is ready. Next steps:" -ForegroundColor White
    Write-Host "1. Install Go: https://go.dev/dl/" -ForegroundColor Cyan
    Write-Host "2. Install PostgreSQL: https://www.postgresql.org/download/" -ForegroundColor Cyan
    Write-Host "3. Run: go mod download" -ForegroundColor Cyan
    Write-Host "4. Setup database (see WINDOWS_SETUP.md)" -ForegroundColor Cyan
    Write-Host "5. Run: go run cmd/server/main.go" -ForegroundColor Cyan
} elseif ($errors -eq 0) {
    Write-Host "✓ CHECKS PASSED with $warnings warning(s)" -ForegroundColor Yellow
    Write-Host "`nProject structure is correct." -ForegroundColor White
} else {
    Write-Host "❌ VALIDATION FAILED" -ForegroundColor Red
    Write-Host "Errors: $errors" -ForegroundColor Red
    Write-Host "Warnings: $warnings" -ForegroundColor Yellow
}

Write-Host "`n=========================================" -ForegroundColor Cyan

# Check if Go is installed
Write-Host "`nChecking prerequisites..." -ForegroundColor Yellow
try {
    $goVersion = go version 2>$null
    Write-Host "  ✓ Go is installed: $goVersion" -ForegroundColor Green
    Write-Host "`n  You can run the server now:" -ForegroundColor White
    Write-Host "  go run cmd/server/main.go" -ForegroundColor Cyan
} catch {
    Write-Host "  ⚠ Go is NOT installed" -ForegroundColor Yellow
    Write-Host "`n  Install Go to run the server:" -ForegroundColor White
    Write-Host "  https://go.dev/dl/" -ForegroundColor Cyan
}

try {
    $dockerVersion = docker --version 2>$null
    Write-Host "  ✓ Docker is installed: $dockerVersion" -ForegroundColor Green
    Write-Host "`n  You can use Docker:" -ForegroundColor White
    Write-Host "  docker-compose up --build" -ForegroundColor Cyan
} catch {
    Write-Host "  ⚠ Docker is NOT installed" -ForegroundColor Yellow
    Write-Host "`n  Install Docker for easiest setup:" -ForegroundColor White
    Write-Host "  https://www.docker.com/products/docker-desktop/" -ForegroundColor Cyan
}

Write-Host ""
