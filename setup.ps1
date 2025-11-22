# Quick Start Script for Stocky (Windows PowerShell)

Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "Stocky Setup & Installation" -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host ""

# Check if Go is installed
Write-Host "Checking for Go installation..." -ForegroundColor Yellow
try {
    $goVersion = go version
    Write-Host "✓ Go is installed: $goVersion" -ForegroundColor Green
} catch {
    Write-Host "✗ Go is not installed!" -ForegroundColor Red
    Write-Host "Please install Go from: https://go.dev/dl/" -ForegroundColor Yellow
    Write-Host "After installation, restart PowerShell and run this script again." -ForegroundColor Yellow
    exit 1
}

# Check if PostgreSQL is installed
Write-Host "Checking for PostgreSQL..." -ForegroundColor Yellow
try {
    $pgVersion = psql --version
    Write-Host "✓ PostgreSQL is installed: $pgVersion" -ForegroundColor Green
} catch {
    Write-Host "✗ PostgreSQL is not installed!" -ForegroundColor Red
    Write-Host "Please install PostgreSQL from: https://www.postgresql.org/download/windows/" -ForegroundColor Yellow
    Write-Host "After installation, restart PowerShell and run this script again." -ForegroundColor Yellow
    exit 1
}

Write-Host ""
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "Installing Dependencies" -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan

# Install Go dependencies
Write-Host "Downloading Go modules..." -ForegroundColor Yellow
go mod download
if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Go dependencies installed successfully" -ForegroundColor Green
} else {
    Write-Host "✗ Failed to install Go dependencies" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "Database Setup" -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Yellow
Write-Host "1. Create PostgreSQL database and user (if not already done):" -ForegroundColor White
Write-Host "   psql -U postgres" -ForegroundColor Cyan
Write-Host "   CREATE DATABASE assignment;" -ForegroundColor Cyan
Write-Host "   CREATE USER stocky_user WITH PASSWORD 'stocky_password';" -ForegroundColor Cyan
Write-Host "   GRANT ALL PRIVILEGES ON DATABASE assignment TO stocky_user;" -ForegroundColor Cyan
Write-Host "   \q" -ForegroundColor Cyan
Write-Host ""
Write-Host "2. Copy .env.example to .env and update credentials:" -ForegroundColor White
Write-Host "   Copy-Item .env.example .env" -ForegroundColor Cyan
Write-Host "   notepad .env" -ForegroundColor Cyan
Write-Host ""
Write-Host "3. Run migrations:" -ForegroundColor White
Write-Host "   Get-Content migrations\001_initial_schema.sql | psql -U stocky_user -d assignment" -ForegroundColor Cyan
Write-Host ""
Write-Host "4. Start the server:" -ForegroundColor White
Write-Host "   go run cmd/server/main.go" -ForegroundColor Cyan
Write-Host ""
Write-Host "5. Test the API:" -ForegroundColor White
Write-Host "   .\test-api.ps1" -ForegroundColor Cyan
Write-Host ""
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "Setup Complete!" -ForegroundColor Green
Write-Host "=========================================" -ForegroundColor Cyan
