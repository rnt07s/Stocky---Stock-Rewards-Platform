package models

import (
	"database/sql"
	"time"
)

// User represents a user in the system
type User struct {
	ID        int64     `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Stock represents a stock symbol
type Stock struct {
	ID          int64     `json:"id"`
	Symbol      string    `json:"symbol"`
	CompanyName string    `json:"company_name"`
	Exchange    string    `json:"exchange"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// RewardEvent represents a single reward transaction
type RewardEvent struct {
	ID              int64          `json:"id"`
	IdempotencyKey  string         `json:"idempotency_key"`
	UserID          string         `json:"user_id"`
	StockSymbol     string         `json:"stock_symbol"`
	SharesQuantity  float64        `json:"shares_quantity"`
	PricePerShare   float64        `json:"price_per_share"`
	TotalValue      float64        `json:"total_value"`
	BrokerageFee    float64        `json:"brokerage_fee"`
	STTFee          float64        `json:"stt_fee"`
	GSTFee          float64        `json:"gst_fee"`
	ExchangeFee     float64        `json:"exchange_fee"`
	SEBIFee         float64        `json:"sebi_fee"`
	TotalFees       float64        `json:"total_fees"`
	TotalCost       float64        `json:"total_cost"`
	Reason          string         `json:"reason"`
	Metadata        sql.NullString `json:"metadata,omitempty"`
	RewardedAt      time.Time      `json:"rewarded_at"`
	CreatedAt       time.Time      `json:"created_at"`
}

// UserHolding represents aggregated holdings for a user
type UserHolding struct {
	ID           int64     `json:"id"`
	UserID       string    `json:"user_id"`
	StockSymbol  string    `json:"stock_symbol"`
	TotalShares  float64   `json:"total_shares"`
	AveragePrice float64   `json:"average_price"`
	LastUpdated  time.Time `json:"last_updated"`
}

// StockPrice represents a stock price snapshot
type StockPrice struct {
	ID          int64     `json:"id"`
	StockSymbol string    `json:"stock_symbol"`
	Price       float64   `json:"price"`
	Timestamp   time.Time `json:"timestamp"`
	Source      string    `json:"source"`
}

// LedgerEntry represents a double-entry accounting record
type LedgerEntry struct {
	ID             int64          `json:"id"`
	EntryGroupID   string         `json:"entry_group_id"`
	RewardEventID  sql.NullInt64  `json:"reward_event_id,omitempty"`
	AccountType    string         `json:"account_type"`
	StockSymbol    sql.NullString `json:"stock_symbol,omitempty"`
	DebitAmount    float64        `json:"debit_amount"`
	CreditAmount   float64        `json:"credit_amount"`
	Description    string         `json:"description"`
	CreatedAt      time.Time      `json:"created_at"`
}

// StockEvent represents corporate actions like splits, mergers
type StockEvent struct {
	ID                 int64          `json:"id"`
	StockSymbol        string         `json:"stock_symbol"`
	EventType          string         `json:"event_type"`
	EventDate          time.Time      `json:"event_date"`
	SplitRatioOld      sql.NullInt64  `json:"split_ratio_old,omitempty"`
	SplitRatioNew      sql.NullInt64  `json:"split_ratio_new,omitempty"`
	MergedIntoSymbol   sql.NullString `json:"merged_into_symbol,omitempty"`
	ConversionRatio    sql.NullFloat64 `json:"conversion_ratio,omitempty"`
	Description        string         `json:"description"`
	Processed          bool           `json:"processed"`
	ProcessedAt        sql.NullTime   `json:"processed_at,omitempty"`
	CreatedAt          time.Time      `json:"created_at"`
}

// Portfolio represents a user's complete portfolio
type PortfolioItem struct {
	StockSymbol  string  `json:"stock_symbol"`
	CompanyName  string  `json:"company_name"`
	TotalShares  float64 `json:"total_shares"`
	AveragePrice float64 `json:"average_price"`
	CurrentPrice float64 `json:"current_price"`
	CurrentValue float64 `json:"current_value"`
	TotalCost    float64 `json:"total_cost"`
	ProfitLoss   float64 `json:"profit_loss"`
	ProfitLossPct float64 `json:"profit_loss_pct"`
}
