# Backend Startup Script for Windows PowerShell
# Run this to start the backend server

Write-Host "🚀 Starting Yayasan ERP Backend..." -ForegroundColor Green
Write-Host ""

# Check if Go is installed
Write-Host "🔍 Checking Go installation..." -ForegroundColor Cyan
$goVersion = go version 2>$null
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Go is not installed!" -ForegroundColor Red
    Write-Host "   Please install Go 1.21+ from https://go.dev/dl/" -ForegroundColor Yellow
    exit 1
}
Write-Host "✅ Go installed: $goVersion" -ForegroundColor Green
Write-Host ""

# Check if .env exists
Write-Host "🔍 Checking .env file..." -ForegroundColor Cyan
if (-not (Test-Path ".env")) {
    Write-Host "⚠️  .env file not found. Creating from template..." -ForegroundColor Yellow
    Copy-Item ".env.example" ".env" -ErrorAction SilentlyContinue
    if (Test-Path ".env") {
        Write-Host "✅ .env created from template" -ForegroundColor Green
        Write-Host "   Please update .env with your database credentials" -ForegroundColor Yellow
    } else {
        Write-Host "❌ Could not create .env file" -ForegroundColor Red
        Write-Host "   Please create .env manually" -ForegroundColor Yellow
    }
    Write-Host ""
}

# Download dependencies
Write-Host "📦 Downloading dependencies..." -ForegroundColor Cyan
go mod tidy
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Failed to download dependencies!" -ForegroundColor Red
    Write-Host "   Try running: go mod tidy" -ForegroundColor Yellow
    exit 1
}
Write-Host "✅ Dependencies downloaded" -ForegroundColor Green
Write-Host ""

# Check if go.sum exists
if (-not (Test-Path "go.sum")) {
    Write-Host "⚠️  go.sum not found. Running go mod download..." -ForegroundColor Yellow
    go mod download
    Write-Host ""
}

# Start the server
Write-Host "🚀 Starting server..." -ForegroundColor Cyan
Write-Host "   Server will be available at: http://localhost:8080" -ForegroundColor Yellow
Write-Host "   API documentation: http://localhost:8080/api/v1/health" -ForegroundColor Yellow
Write-Host ""
Write-Host "   Press Ctrl+C to stop the server" -ForegroundColor Gray
Write-Host ""

go run cmd/api/main.go
