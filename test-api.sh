#!/bin/bash

# Stocky API Test Script

BASE_URL="http://localhost:8080/api/v1"

echo "========================================="
echo "Stocky API Testing Script"
echo "========================================="
echo ""

# Health check
echo "1. Testing Health Check..."
curl -s "$BASE_URL/../health" | jq '.'
echo -e "\n"

# Create reward for user alice
echo "2. Creating reward for user 'alice' (TCS shares)..."
curl -s -X POST "$BASE_URL/reward" \
  -H "Content-Type: application/json" \
  -d '{
    "idempotency_key": "reward-alice-20250122-001",
    "user_id": "alice",
    "stock_symbol": "TCS",
    "shares_quantity": 2.5,
    "reason": "onboarding_bonus"
  }' | jq '.'
echo -e "\n"

# Create another reward for alice
echo "3. Creating another reward for user 'alice' (INFY shares)..."
curl -s -X POST "$BASE_URL/reward" \
  -H "Content-Type: application/json" \
  -d '{
    "idempotency_key": "reward-alice-20250122-002",
    "user_id": "alice",
    "stock_symbol": "INFY",
    "shares_quantity": 5.0,
    "reason": "referral_bonus"
  }' | jq '.'
echo -e "\n"

# Test idempotency
echo "4. Testing Idempotency (duplicate request)..."
curl -s -X POST "$BASE_URL/reward" \
  -H "Content-Type: application/json" \
  -d '{
    "idempotency_key": "reward-alice-20250122-001",
    "user_id": "alice",
    "stock_symbol": "TCS",
    "shares_quantity": 2.5,
    "reason": "onboarding_bonus"
  }' | jq '.'
echo -e "\n"

# Get today's stocks
echo "5. Getting today's stocks for 'alice'..."
curl -s "$BASE_URL/today-stocks/alice" | jq '.'
echo -e "\n"

# Get user stats
echo "6. Getting stats for 'alice'..."
curl -s "$BASE_URL/stats/alice" | jq '.'
echo -e "\n"

# Get portfolio
echo "7. Getting portfolio for 'alice'..."
curl -s "$BASE_URL/portfolio/alice" | jq '.'
echo -e "\n"

# Create reward for user bob
echo "8. Creating reward for user 'bob' (RELIANCE shares)..."
curl -s -X POST "$BASE_URL/reward" \
  -H "Content-Type: application/json" \
  -d '{
    "idempotency_key": "reward-bob-20250122-001",
    "user_id": "bob",
    "stock_symbol": "RELIANCE",
    "shares_quantity": 1.75,
    "reason": "milestone_achieved"
  }' | jq '.'
echo -e "\n"

# Get portfolio for bob
echo "9. Getting portfolio for 'bob'..."
curl -s "$BASE_URL/portfolio/bob" | jq '.'
echo -e "\n"

echo "========================================="
echo "Testing Complete!"
echo "========================================="
