# Stocky API Test Script (PowerShell)

$baseUrl = "http://localhost:8080/api/v1"

Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "Stocky API Testing Script" -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host ""

# Health check
Write-Host "1. Testing Health Check..." -ForegroundColor Yellow
Invoke-RestMethod -Uri "$baseUrl/../health" -Method Get | ConvertTo-Json
Write-Host ""

# Create reward for user alice
Write-Host "2. Creating reward for user 'alice' (TCS shares)..." -ForegroundColor Yellow
$body1 = @{
    idempotency_key = "reward-alice-20250122-001"
    user_id = "alice"
    stock_symbol = "TCS"
    shares_quantity = 2.5
    reason = "onboarding_bonus"
} | ConvertTo-Json

Invoke-RestMethod -Uri "$baseUrl/reward" -Method Post -Body $body1 -ContentType "application/json" | ConvertTo-Json -Depth 10
Write-Host ""

# Create another reward for alice
Write-Host "3. Creating another reward for user 'alice' (INFY shares)..." -ForegroundColor Yellow
$body2 = @{
    idempotency_key = "reward-alice-20250122-002"
    user_id = "alice"
    stock_symbol = "INFY"
    shares_quantity = 5.0
    reason = "referral_bonus"
} | ConvertTo-Json

Invoke-RestMethod -Uri "$baseUrl/reward" -Method Post -Body $body2 -ContentType "application/json" | ConvertTo-Json -Depth 10
Write-Host ""

# Test idempotency
Write-Host "4. Testing Idempotency (duplicate request)..." -ForegroundColor Yellow
Invoke-RestMethod -Uri "$baseUrl/reward" -Method Post -Body $body1 -ContentType "application/json" | ConvertTo-Json -Depth 10
Write-Host ""

# Get today's stocks
Write-Host "5. Getting today's stocks for 'alice'..." -ForegroundColor Yellow
Invoke-RestMethod -Uri "$baseUrl/today-stocks/alice" -Method Get | ConvertTo-Json -Depth 10
Write-Host ""

# Get user stats
Write-Host "6. Getting stats for 'alice'..." -ForegroundColor Yellow
Invoke-RestMethod -Uri "$baseUrl/stats/alice" -Method Get | ConvertTo-Json -Depth 10
Write-Host ""

# Get portfolio
Write-Host "7. Getting portfolio for 'alice'..." -ForegroundColor Yellow
Invoke-RestMethod -Uri "$baseUrl/portfolio/alice" -Method Get | ConvertTo-Json -Depth 10
Write-Host ""

# Create reward for user bob
Write-Host "8. Creating reward for user 'bob' (RELIANCE shares)..." -ForegroundColor Yellow
$body3 = @{
    idempotency_key = "reward-bob-20250122-001"
    user_id = "bob"
    stock_symbol = "RELIANCE"
    shares_quantity = 1.75
    reason = "milestone_achieved"
} | ConvertTo-Json

Invoke-RestMethod -Uri "$baseUrl/reward" -Method Post -Body $body3 -ContentType "application/json" | ConvertTo-Json -Depth 10
Write-Host ""

# Get portfolio for bob
Write-Host "9. Getting portfolio for 'bob'..." -ForegroundColor Yellow
Invoke-RestMethod -Uri "$baseUrl/portfolio/bob" -Method Get | ConvertTo-Json -Depth 10
Write-Host ""

Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "Testing Complete!" -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan
