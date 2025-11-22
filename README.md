# Stocky - Stock Rewards Platform

A robust REST API system that enables users to earn shares of Indian stocks (NSE/BSE) as rewards for various actions. Built with Go, Gin, and PostgreSQL, featuring double-entry bookkeeping, hourly price updates, and comprehensive edge case handling.

## 📋 Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Database Schema](#database-schema)
- [API Endpoints](#api-endpoints)
- [Setup Instructions](#setup-instructions)
- [Edge Cases & Solutions](#edge-cases--solutions)
- [Scaling Considerations](#scaling-considerations)

---

## ✨ Features

- **Stock Rewards System**: Award fractional shares to users with full transparency
- **Double-Entry Ledger**: Track all financial transactions (cash, fees, stock inventory)
- **Hourly Price Updates**: Mock stock price service with configurable intervals
- **Idempotency**: Prevent duplicate reward events using idempotency keys
- **Portfolio Management**: Real-time INR valuation of user holdings
- **Fee Tracking**: Brokerage, STT, GST, Exchange, and SEBI fees
- **Corporate Actions**: Support for stock splits, mergers, and delisting

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────┐
│                     API Layer (Gin)                      │
│  POST /reward  │  GET /today-stocks/:userId              │
│  GET /historical-inr/:userId  │  GET /stats/:userId      │
│  GET /portfolio/:userId                                  │
└─────────────────┬───────────────────────────────────────┘
                  │
┌─────────────────┴───────────────────────────────────────┐
│                  Service Layer                           │
│  • RewardService (business logic)                        │
│  • StockPriceService (hourly updates)                    │
└─────────────────┬───────────────────────────────────────┘
                  │
┌─────────────────┴───────────────────────────────────────┐
│                Repository Layer                          │
│  • RewardRepository  • StockRepository                   │
│  • LedgerRepository                                      │
└─────────────────┬───────────────────────────────────────┘
                  │
┌─────────────────┴───────────────────────────────────────┐
│              PostgreSQL Database                         │
│  Tables: reward_events, user_holdings, stock_prices,    │
│          ledger_entries, stock_events, users, stocks    │
└──────────────────────────────────────────────────────────┘
```

### Package Structure

```
stocky/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── database/
│   │   └── database.go          # DB connection & migrations
│   ├── models/
│   │   └── models.go            # Data models
│   ├── repository/
│   │   └── repository.go        # Data access layer
│   ├── services/
│   │   └── services.go          # Business logic
│   ├── handlers/
│   │   └── handlers.go          # HTTP handlers
│   └── middleware/
│       └── middleware.go        # Logging, CORS, recovery
├── migrations/
│   └── 001_initial_schema.sql   # Database schema
├── go.mod
├── .env.example
└── README.md
```

---

## 🗄️ Database Schema

### 1. **users**
Stores user information.

| Column     | Type         | Description          |
|------------|--------------|----------------------|
| id         | SERIAL       | Primary key          |
| user_id    | VARCHAR(100) | Unique user ID       |
| name       | VARCHAR(255) | User name            |
| email      | VARCHAR(255) | User email           |
| created_at | TIMESTAMP    | Creation timestamp   |
| updated_at | TIMESTAMP    | Update timestamp     |

---

### 2. **stocks**
Reference table for stock symbols.

| Column       | Type         | Description               |
|--------------|--------------|---------------------------|
| id           | SERIAL       | Primary key               |
| symbol       | VARCHAR(20)  | Stock symbol (e.g., TCS)  |
| company_name | VARCHAR(255) | Full company name         |
| exchange     | VARCHAR(10)  | NSE or BSE                |
| is_active    | BOOLEAN      | Active status             |
| created_at   | TIMESTAMP    | Creation timestamp        |
| updated_at   | TIMESTAMP    | Update timestamp          |

---

### 3. **reward_events** (Immutable Log)
Records every reward transaction.

| Column          | Type            | Description                      |
|-----------------|-----------------|----------------------------------|
| id              | SERIAL          | Primary key                      |
| idempotency_key | VARCHAR(255)    | Unique key (prevents duplicates) |
| user_id         | VARCHAR(100)    | User receiving reward            |
| stock_symbol    | VARCHAR(20)     | Stock symbol                     |
| shares_quantity | NUMERIC(18,6)   | Fractional shares awarded        |
| price_per_share | NUMERIC(18,4)   | Price at reward time (INR)       |
| total_value     | NUMERIC(18,4)   | shares × price                   |
| brokerage_fee   | NUMERIC(18,4)   | Brokerage fee paid by Stocky     |
| stt_fee         | NUMERIC(18,4)   | Securities Transaction Tax       |
| gst_fee         | NUMERIC(18,4)   | GST on brokerage                 |
| exchange_fee    | NUMERIC(18,4)   | Exchange transaction fee         |
| sebi_fee        | NUMERIC(18,4)   | SEBI turnover charges            |
| total_fees      | NUMERIC(18,4)   | Sum of all fees                  |
| total_cost      | NUMERIC(18,4)   | total_value + total_fees         |
| reason          | VARCHAR(255)    | Reward reason                    |
| metadata        | JSONB           | Additional context               |
| rewarded_at     | TIMESTAMP       | When reward was given            |
| created_at      | TIMESTAMP       | Record creation time             |

**Indexes**: `user_id`, `stock_symbol`, `rewarded_at`, `idempotency_key`

---

### 4. **user_holdings** (Aggregated)
Current holdings per user per stock.

| Column        | Type            | Description                   |
|---------------|-----------------|-------------------------------|
| id            | SERIAL          | Primary key                   |
| user_id       | VARCHAR(100)    | User ID                       |
| stock_symbol  | VARCHAR(20)     | Stock symbol                  |
| total_shares  | NUMERIC(18,6)   | Total shares owned            |
| average_price | NUMERIC(18,4)   | Weighted average cost         |
| last_updated  | TIMESTAMP       | Last update timestamp         |

**Unique constraint**: `(user_id, stock_symbol)`

---

### 5. **stock_prices**
Hourly price snapshots.

| Column       | Type            | Description                |
|--------------|-----------------|----------------------------|
| id           | SERIAL          | Primary key                |
| stock_symbol | VARCHAR(20)     | Stock symbol               |
| price        | NUMERIC(18,4)   | Price in INR               |
| timestamp    | TIMESTAMP       | Price snapshot time        |
| source       | VARCHAR(50)     | Data source (mock/nse/bse) |

**Unique constraint**: `(stock_symbol, timestamp)`

---

### 6. **ledger_entries** (Double-Entry Bookkeeping)
Tracks all financial transactions.

| Column         | Type            | Description                          |
|----------------|-----------------|--------------------------------------|
| id             | SERIAL          | Primary key                          |
| entry_group_id | UUID            | Groups related debit/credit entries  |
| reward_event_id| INTEGER         | FK to reward_events                  |
| account_type   | VARCHAR(50)     | stock_inventory, cash_outflow, fees_expense |
| stock_symbol   | VARCHAR(20)     | Stock symbol (NULL for cash accounts)|
| debit_amount   | NUMERIC(18,4)   | Debit amount                         |
| credit_amount  | NUMERIC(18,4)   | Credit amount                        |
| description    | TEXT            | Entry description                    |
| created_at     | TIMESTAMP       | Creation timestamp                   |

**Constraint**: Either debit or credit must be > 0 (not both)

---

### 7. **stock_events**
Tracks corporate actions (splits, mergers, delisting).

| Column              | Type            | Description                     |
|---------------------|-----------------|---------------------------------|
| id                  | SERIAL          | Primary key                     |
| stock_symbol        | VARCHAR(20)     | Affected stock                  |
| event_type          | VARCHAR(50)     | split, merger, delisting, bonus |
| event_date          | DATE            | Event date                      |
| split_ratio_old     | INTEGER         | Old ratio (e.g., 1 in 1:10)     |
| split_ratio_new     | INTEGER         | New ratio (e.g., 10 in 1:10)    |
| merged_into_symbol  | VARCHAR(20)     | Target stock for mergers        |
| conversion_ratio    | NUMERIC(18,6)   | Conversion ratio                |
| description         | TEXT            | Event details                   |
| processed           | BOOLEAN         | Processing status               |
| processed_at        | TIMESTAMP       | Processing timestamp            |
| created_at          | TIMESTAMP       | Creation timestamp              |

---

## 🌐 API Endpoints

### Base URL: `http://localhost:8080/api/v1`

---

### 1. **POST /reward**
Award shares to a user.

**Request Body:**
```json
{
  "idempotency_key": "reward-user123-20250122-001",
  "user_id": "user123",
  "stock_symbol": "TCS",
  "shares_quantity": 2.5,
  "reason": "referral_bonus",
  "metadata": "{\"referral_code\": \"REF123\"}",
  "rewarded_at": "2025-01-22T10:30:00Z"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "idempotency_key": "reward-user123-20250122-001",
    "user_id": "user123",
    "stock_symbol": "TCS",
    "shares_quantity": 2.5,
    "price_per_share": 3521.45,
    "total_value": 8803.63,
    "brokerage_fee": 4.40,
    "stt_fee": 22.01,
    "gst_fee": 0.79,
    "exchange_fee": 2.64,
    "sebi_fee": 0.88,
    "total_fees": 30.72,
    "total_cost": 8834.35,
    "reason": "referral_bonus",
    "rewarded_at": "2025-01-22T10:30:00Z",
    "created_at": "2025-01-22T10:30:05Z"
  }
}
```

**Idempotency**: Sending the same `idempotency_key` returns the original event without creating duplicates.

---

### 2. **GET /today-stocks/:userId**
Get all stock rewards for a user today.

**Example:** `GET /today-stocks/user123`

**Response (200 OK):**
```json
{
  "success": true,
  "user_id": "user123",
  "date": "2025-01-22",
  "count": 3,
  "data": [
    {
      "id": 1,
      "stock_symbol": "TCS",
      "shares_quantity": 2.5,
      "price_per_share": 3521.45,
      "total_value": 8803.63,
      "rewarded_at": "2025-01-22T10:30:00Z"
    },
    {
      "id": 2,
      "stock_symbol": "INFY",
      "shares_quantity": 5.0,
      "price_per_share": 1487.20,
      "total_value": 7436.00,
      "rewarded_at": "2025-01-22T14:15:00Z"
    }
  ]
}
```

---

### 3. **GET /historical-inr/:userId**
Get INR value of rewards for all past days (excluding today).

**Example:** `GET /historical-inr/user123`

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "user_id": "user123",
    "daily_inr": [
      {
        "date": "2025-01-21",
        "total_value": 15240.50
      },
      {
        "date": "2025-01-20",
        "total_value": 8920.00
      }
    ],
    "total_value": 24160.50
  }
}
```

---

### 4. **GET /stats/:userId**
Get user statistics including today's rewards and portfolio value.

**Example:** `GET /stats/user123`

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "user_id": "user123",
    "today_rewards": [
      {
        "stock_symbol": "TCS",
        "total_shares": 2.5,
        "reward_count": 1
      },
      {
        "stock_symbol": "INFY",
        "total_shares": 5.0,
        "reward_count": 1
      }
    ],
    "portfolio_value_inr": 45678.90,
    "total_shares_rewarded": 7.5
  }
}
```

---

### 5. **GET /portfolio/:userId** (Bonus)
Get complete portfolio with profit/loss analysis.

**Example:** `GET /portfolio/user123`

**Response (200 OK):**
```json
{
  "success": true,
  "user_id": "user123",
  "summary": {
    "total_value": 45678.90,
    "total_cost": 42350.25,
    "total_profit_loss": 3328.65,
    "holdings_count": 3
  },
  "holdings": [
    {
      "stock_symbol": "TCS",
      "company_name": "Tata Consultancy Services Limited",
      "total_shares": 10.5,
      "average_price": 3480.25,
      "current_price": 3521.45,
      "current_value": 36975.23,
      "total_cost": 36542.63,
      "profit_loss": 432.60,
      "profit_loss_pct": 1.18
    },
    {
      "stock_symbol": "INFY",
      "company_name": "Infosys Limited",
      "total_shares": 5.8,
      "average_price": 1475.00,
      "current_price": 1487.20,
      "current_value": 8625.76,
      "total_cost": 8555.00,
      "profit_loss": 70.76,
      "profit_loss_pct": 0.83
    }
  ]
}
```

---

### 6. **GET /health**
Health check endpoint.

**Response (200 OK):**
```json
{
  "status": "healthy",
  "timestamp": "2025-01-22T10:30:00Z"
}
```

---

## 🚀 Setup Instructions

### Prerequisites

- Go 1.21+
- PostgreSQL 13+
- Git

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/stocky-assignment.git
cd stocky-assignment
```

### 2. Setup PostgreSQL Database

```bash
# Create database and user
psql -U postgres

CREATE DATABASE assignment;
CREATE USER stocky_user WITH PASSWORD 'stocky_password';
GRANT ALL PRIVILEGES ON DATABASE assignment TO stocky_user;
\q
```

### 3. Configure Environment

```bash
cp .env.example .env
# Edit .env with your database credentials
```

**`.env` file:**
```env
PORT=8080
GIN_MODE=debug

DB_HOST=localhost
DB_PORT=5432
DB_USER=stocky_user
DB_PASSWORD=stocky_password
DB_NAME=assignment
DB_SSLMODE=disable

PRICE_UPDATE_INTERVAL_MINUTES=60

BROKERAGE_FEE_BP=5
STT_FEE_BP=25
GST_FEE_BP=18
EXCHANGE_FEE_BP=3
SEBI_FEE_BP=1
```

### 4. Install Dependencies

```bash
go mod download
```

### 5. Run Migrations

Migrations run automatically on server start, or manually:

```bash
psql -U stocky_user -d assignment -f migrations/001_initial_schema.sql
```

### 6. Start the Server

```bash
go run cmd/server/main.go
```

Server starts on `http://localhost:8080`

---

## 🧪 Testing the API

### Using cURL

**Create a reward:**
```bash
curl -X POST http://localhost:8080/api/v1/reward \
  -H "Content-Type: application/json" \
  -d '{
    "idempotency_key": "reward-alice-20250122-001",
    "user_id": "alice",
    "stock_symbol": "TCS",
    "shares_quantity": 2.5,
    "reason": "onboarding_bonus"
  }'
```

**Get today's stocks:**
```bash
curl http://localhost:8080/api/v1/today-stocks/alice
```

**Get portfolio:**
```bash
curl http://localhost:8080/api/v1/portfolio/alice
```

---

## 🛡️ Edge Cases & Solutions

### 1. **Duplicate Reward Events / Replay Attacks**

**Problem**: API might receive duplicate requests due to network retries.

**Solution**: 
- Mandatory `idempotency_key` field in POST /reward
- Database unique constraint on `idempotency_key`
- Returns existing event if key already exists (HTTP 201)

**Implementation**:
```go
existingEvent, _ := repo.GetRewardEventByIdempotencyKey(req.IdempotencyKey)
if existingEvent != nil {
    return existingEvent, nil // Return cached result
}
```

---

### 2. **Stock Splits**

**Problem**: 1:10 split means 1 share becomes 10 shares.

**Solution**:
- `stock_events` table tracks splits with `split_ratio_old` and `split_ratio_new`
- Background job processes unprocessed events
- Updates `user_holdings.total_shares` proportionally
- Adjusts `average_price` inversely

**Example**:
```sql
-- Before split: 10 shares @ ₹3500 = ₹35,000
-- After 1:2 split: 20 shares @ ₹1750 = ₹35,000
UPDATE user_holdings 
SET total_shares = total_shares * 2,
    average_price = average_price / 2
WHERE stock_symbol = 'TCS';
```

---

### 3. **Stock Mergers**

**Problem**: Company A merges into Company B with conversion ratio.

**Solution**:
- Record merger event in `stock_events`
- Convert holdings: `holding_B = holding_A × conversion_ratio`
- Mark original stock inactive

---

### 4. **Delisting**

**Problem**: Stock is delisted from exchange.

**Solution**:
- Set `stocks.is_active = false`
- Prevent new rewards for delisted stocks
- Retain historical data for compliance
- Portfolio shows last known price with warning

---

### 5. **Rounding Errors in INR Valuation**

**Problem**: Fractional shares × fractional prices cause precision errors.

**Solution**:
- Use `NUMERIC(18, 6)` for shares (6 decimals)
- Use `NUMERIC(18, 4)` for INR (4 decimals = paise level)
- Round calculations explicitly:

```go
func roundToDecimal(value float64, decimals int) float64 {
    multiplier := math.Pow(10, float64(decimals))
    return math.Round(value*multiplier) / multiplier
}
```

---

### 6. **Price API Downtime / Stale Data**

**Problem**: External price API fails or returns stale data.

**Solution**:
- **Fallback**: Use last known price from `stock_prices` table
- **Staleness check**: Flag prices older than 2 hours
- **Circuit breaker**: Retry with exponential backoff
- **Manual override**: Admin can set prices manually

```go
price, err := priceService.GetCurrentPrice(symbol)
if err != nil {
    // Fallback to last known price
    price = holding.AveragePrice
    log.Warn("Using fallback price")
}
```

---

### 7. **Adjustments/Refunds of Rewards**

**Problem**: Need to revoke incorrectly given rewards.

**Solution**:
- Create negative reward event (shares_quantity < 0 allowed for adjustments)
- Update `user_holdings` accordingly
- Ledger entries reflect reversal (debit/credit swapped)

**Example**:
```json
{
  "idempotency_key": "refund-user123-20250122-001",
  "user_id": "user123",
  "stock_symbol": "TCS",
  "shares_quantity": -2.5,
  "reason": "correction"
}
```

---

## 📈 Scaling Considerations

### 1. **Database Optimization**

- **Partitioning**: Partition `reward_events` by `rewarded_at` (monthly)
- **Read replicas**: Route reads to replicas, writes to primary
- **Connection pooling**: Already configured (max 25 open connections)
- **Indexes**: Critical columns indexed (see schema)

### 2. **Caching**

- **Redis**: Cache current stock prices (TTL: 5 minutes)
- **User portfolios**: Cache with invalidation on new rewards
- **Leaderboards**: Pre-compute daily/weekly

```go
// Pseudocode
price := redis.Get("price:TCS")
if price == nil {
    price = db.GetLatestPrice("TCS")
    redis.Set("price:TCS", price, 5*time.Minute)
}
```

### 3. **Async Processing**

- **Message queue**: Use RabbitMQ/Kafka for reward processing
- **Worker pool**: Dedicated workers for fee calculations, ledger entries
- **Batch updates**: Update stock prices in bulk

### 4. **Horizontal Scaling**

- **Stateless servers**: Multiple API instances behind load balancer
- **Database sharding**: Shard by `user_id` for large user base
- **Microservices**: Split into services (rewards, portfolio, analytics)

### 5. **Monitoring & Observability**

- **Metrics**: Prometheus + Grafana
  - Request latency (p50, p95, p99)
  - Error rates
  - Database query times
- **Logging**: Structured logs (Logrus → ELK stack)
- **Alerting**: PagerDuty for critical failures

### 6. **Rate Limiting**

- **Per-user limits**: 100 requests/minute
- **Per-IP limits**: 1000 requests/minute
- **Token bucket algorithm**

```go
// Pseudocode
if !rateLimiter.Allow(userID) {
    return http.StatusTooManyRequests
}
```

### 7. **Load Testing Benchmarks**

Expected performance on moderate hardware (4 CPU, 8GB RAM):

- **POST /reward**: 500 req/s
- **GET /portfolio**: 2000 req/s (with caching)
- **Price updates**: 10,000 stocks in <5 seconds

---

## 🔒 Security Considerations

1. **Authentication**: Add JWT-based auth middleware
2. **Authorization**: Ensure users can only access their own data
3. **Input validation**: Validate all inputs (already using Gin binding)
4. **SQL injection**: Using parameterized queries throughout
5. **HTTPS**: Deploy behind reverse proxy (Nginx) with TLS
6. **Secrets management**: Use vault (e.g., HashiCorp Vault)

---

## 📊 Double-Entry Ledger Example

When user `alice` receives 2.5 shares of TCS @ ₹3500:

| Account Type       | Stock | Debit (₹) | Credit (₹) | Description                  |
|--------------------|-------|-----------|------------|------------------------------|
| `stock_inventory`  | TCS   | 8,750.00  | 0          | Stock purchased for user     |
| `cash_outflow`     | -     | 0         | 8,750.00   | Cash paid for stocks         |
| `fees_expense`     | -     | 30.72     | 0          | Brokerage, STT, GST, etc.    |
| `cash_outflow`     | -     | 0         | 30.72      | Cash paid for fees           |

**Total Debit = Total Credit = ₹8,780.72** ✅

---

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

---

## 📝 License

This project is created as an assignment and is available for educational purposes.

---

## 👤 Author

**Your Name**
- GitHub: [@yourusername](https://github.com/yourusername)
- Email: your.email@example.com

---

## 🙏 Acknowledgments

- Indian stock exchanges (NSE/BSE) for reference data
- Gin framework for excellent HTTP routing
- PostgreSQL for robust database features

---

**Built with ❤️ using Go, Gin, and PostgreSQL**
