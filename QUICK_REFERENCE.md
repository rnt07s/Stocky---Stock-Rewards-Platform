# 🚀 Stocky - Quick Reference Card

## ⚡ 30-Second Overview
Stock rewards platform built with **Go + Gin + PostgreSQL**. Award fractional Indian stock shares to users as incentives. Features: double-entry ledger, idempotency, hourly price updates, comprehensive API.

---

## 📁 Key Files

| File | Purpose |
|------|---------|
| `README.md` | Complete API documentation |
| `ASSIGNMENT_SUMMARY.md` | Assignment checklist & deliverables |
| `ARCHITECTURE.md` | Design decisions & edge cases |
| `VISUAL_GUIDE.md` | Diagrams & flow charts |
| `WINDOWS_SETUP.md` | Installation guide for Windows |
| `GITHUB_SETUP.md` | Push to GitHub instructions |
| `cmd/server/main.go` | Application entry point |
| `internal/services/services.go` | Core business logic |
| `migrations/001_initial_schema.sql` | Database schema |

---

## 🎯 API Endpoints (6 Total)

```
BASE: http://localhost:8080/api/v1

POST   /reward                    # Award shares to user
GET    /today-stocks/:userId      # Today's rewards
GET    /historical-inr/:userId    # Past days INR value
GET    /stats/:userId             # User statistics
GET    /portfolio/:userId         # Complete portfolio (BONUS)
GET    /health                    # Health check
```

---

## 🏃 Quick Start (3 Commands)

### Option 1: Docker (Fastest)
```powershell
cd "c:\Users\RAUNEET SINGH\OneDrive\Desktop\backend Stocky"
docker-compose up --build
# Visit: http://localhost:8080/health
```

### Option 2: Manual
```powershell
# 1. Setup DB (see WINDOWS_SETUP.md)
# 2. Run server
go run cmd/server/main.go
```

---

## 🧪 Test API

```powershell
# PowerShell
.\test-api.ps1

# Or manually
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/portfolio/alice"
```

---

## 🗄️ Database Tables (7 Total)

| Table | Purpose |
|-------|---------|
| `users` | User profiles |
| `stocks` | Stock reference (TCS, INFY, etc.) |
| `reward_events` | Immutable reward log |
| `user_holdings` | Aggregated holdings |
| `stock_prices` | Hourly price snapshots |
| `ledger_entries` | Double-entry bookkeeping |
| `stock_events` | Corporate actions (splits, etc.) |

---

## 💡 Key Features

✅ **Idempotency** - Duplicate requests return same result  
✅ **Double-Entry Ledger** - Full accounting trail  
✅ **Fractional Shares** - Award 2.5 shares, not just whole numbers  
✅ **Hourly Price Updates** - Mock service with realistic prices  
✅ **Fee Tracking** - Brokerage, STT, GST, Exchange, SEBI  
✅ **Edge Cases** - Splits, mergers, delisting, refunds  
✅ **Precise Decimals** - No floating-point errors  
✅ **Production Ready** - Logging, error handling, validation  

---

## 🔧 Configuration (.env)

```env
PORT=8080
DB_HOST=localhost
DB_USER=stocky_user
DB_PASSWORD=stocky_password
DB_NAME=assignment
PRICE_UPDATE_INTERVAL_MINUTES=60
BROKERAGE_FEE_BP=5  # 0.05%
```

---

## 📊 Example Request/Response

### Create Reward
```bash
POST /api/v1/reward
{
  "idempotency_key": "reward-alice-001",
  "user_id": "alice",
  "stock_symbol": "TCS",
  "shares_quantity": 2.5,
  "reason": "onboarding_bonus"
}

Response (201):
{
  "success": true,
  "data": {
    "id": 1,
    "total_value": 8803.63,
    "total_fees": 30.72,
    "total_cost": 8834.35
  }
}
```

### Get Portfolio
```bash
GET /api/v1/portfolio/alice

Response (200):
{
  "success": true,
  "summary": {
    "total_value": 45678.90,
    "total_profit_loss": 3328.65
  },
  "holdings": [
    {
      "stock_symbol": "TCS",
      "total_shares": 10.5,
      "current_price": 3521.45,
      "profit_loss_pct": 1.18
    }
  ]
}
```

---

## 🏗️ Architecture (3 Layers)

```
Handlers (HTTP)
    ↓
Services (Business Logic)
    ↓
Repositories (Data Access)
    ↓
PostgreSQL Database
```

---

## 🛡️ Edge Cases Handled

| Problem | Solution |
|---------|----------|
| Duplicate requests | Idempotency keys |
| Stock splits (1→10) | stock_events table |
| Rounding errors | NUMERIC types |
| Price API down | Fallback to last price |
| Delisted stocks | is_active flag |
| Refunds | Negative reward events |

---

## 📈 Performance

| Metric | Target |
|--------|--------|
| POST /reward | 500 req/s |
| GET /portfolio | 2000 req/s |
| Latency (p95) | < 200ms |

---

## 🔍 Troubleshooting

**Go not found?**
```powershell
# Install from https://go.dev/dl/
# Restart PowerShell
go version
```

**PostgreSQL error?**
```powershell
# Check service is running
Get-Service postgresql*
```

**Port 8080 in use?**
```powershell
# Change port in .env
PORT=8081
```

---

## 📦 Project Structure

```
stocky/
├── cmd/server/main.go           # Entry point
├── internal/
│   ├── handlers/                # HTTP handlers
│   ├── services/                # Business logic
│   ├── repository/              # Data access
│   ├── models/                  # Data structures
│   ├── config/                  # Configuration
│   └── middleware/              # Logging, CORS
├── migrations/                  # SQL schema
├── README.md                    # Documentation
└── docker-compose.yml           # Container setup
```

---

## 🎓 Technologies Used

- **Go 1.21** - Backend language
- **Gin** - Web framework
- **PostgreSQL** - Database
- **Logrus** - Logging
- **Docker** - Containerization

---

## 📝 Assignment Checklist

✅ POST /reward endpoint  
✅ GET /today-stocks endpoint  
✅ GET /historical-inr endpoint  
✅ GET /stats endpoint  
✅ GET /portfolio endpoint (BONUS)  
✅ Database schema with relationships  
✅ Double-entry ledger system  
✅ Fractional shares support  
✅ Fee tracking (5 types)  
✅ Idempotency handling  
✅ Edge cases handled (8 types)  
✅ API documentation  
✅ Setup instructions  
✅ Docker support  
✅ Test scripts  

---

## 🔗 Quick Links

| Resource | Location |
|----------|----------|
| Full API Docs | `README.md` |
| Edge Cases | `ARCHITECTURE.md` |
| Diagrams | `VISUAL_GUIDE.md` |
| Windows Setup | `WINDOWS_SETUP.md` |
| GitHub Guide | `GITHUB_SETUP.md` |
| Postman | `Stocky-API.postman_collection.json` |

---

## 📞 Need Help?

1. Check `README.md` for detailed documentation
2. Review `WINDOWS_SETUP.md` for installation issues
3. See `ARCHITECTURE.md` for design questions
4. View `VISUAL_GUIDE.md` for flow diagrams

---

## 🎉 Ready to Submit!

1. Test locally: `.\test-api.ps1`
2. Push to GitHub: See `GITHUB_SETUP.md`
3. Submit repo URL: `https://github.com/YOUR_USERNAME/stocky-assignment`

---

**Built with ❤️ for the Stocky Assignment**

---

## 💰 Fee Breakdown (Quick Ref)

On ₹10,000 stock purchase:
- Brokerage: ₹5.00 (0.05%)
- STT: ₹25.00 (0.25%)
- GST: ₹0.90 (18% of brokerage)
- Exchange: ₹3.00 (0.03%)
- SEBI: ₹1.00 (0.01%)
- **Total: ₹34.90**

---

## 🗂️ Database Quick Reference

**Precision:**
- Shares: `NUMERIC(18, 6)` → 123456.789012
- INR: `NUMERIC(18, 4)` → 123456.7890

**Key Indexes:**
- `reward_events(user_id, rewarded_at)` - Fast user queries
- `reward_events(idempotency_key)` - Duplicate detection
- `user_holdings(user_id, stock_symbol)` - Portfolio lookups

---

**⭐ Star this repo if you found it helpful!**
