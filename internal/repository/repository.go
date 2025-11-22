package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/stocky/assignment/internal/models"
)

type RewardRepository interface {
	CreateRewardEvent(event *models.RewardEvent) error
	GetRewardEventByIdempotencyKey(key string) (*models.RewardEvent, error)
	GetTodayRewards(userID string) ([]models.RewardEvent, error)
	GetHistoricalRewards(userID string) ([]models.RewardEvent, error)
	GetUserHolding(userID, stockSymbol string) (*models.UserHolding, error)
	UpsertUserHolding(holding *models.UserHolding) error
	GetUserPortfolio(userID string) ([]models.UserHolding, error)
}

type rewardRepository struct {
	db *sql.DB
}

func NewRewardRepository(db *sql.DB) RewardRepository {
	return &rewardRepository{db: db}
}

func (r *rewardRepository) CreateRewardEvent(event *models.RewardEvent) error {
	query := `
		INSERT INTO reward_events (
			idempotency_key, user_id, stock_symbol, shares_quantity, 
			price_per_share, total_value, brokerage_fee, stt_fee, 
			gst_fee, exchange_fee, sebi_fee, total_fees, total_cost,
			reason, metadata, rewarded_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING id, created_at
	`

	return r.db.QueryRow(
		query,
		event.IdempotencyKey, event.UserID, event.StockSymbol, event.SharesQuantity,
		event.PricePerShare, event.TotalValue, event.BrokerageFee, event.STTFee,
		event.GSTFee, event.ExchangeFee, event.SEBIFee, event.TotalFees, event.TotalCost,
		event.Reason, event.Metadata, event.RewardedAt,
	).Scan(&event.ID, &event.CreatedAt)
}

func (r *rewardRepository) GetRewardEventByIdempotencyKey(key string) (*models.RewardEvent, error) {
	query := `
		SELECT id, idempotency_key, user_id, stock_symbol, shares_quantity,
			   price_per_share, total_value, brokerage_fee, stt_fee, gst_fee,
			   exchange_fee, sebi_fee, total_fees, total_cost, reason, metadata,
			   rewarded_at, created_at
		FROM reward_events
		WHERE idempotency_key = $1
	`

	event := &models.RewardEvent{}
	err := r.db.QueryRow(query, key).Scan(
		&event.ID, &event.IdempotencyKey, &event.UserID, &event.StockSymbol,
		&event.SharesQuantity, &event.PricePerShare, &event.TotalValue,
		&event.BrokerageFee, &event.STTFee, &event.GSTFee, &event.ExchangeFee,
		&event.SEBIFee, &event.TotalFees, &event.TotalCost, &event.Reason,
		&event.Metadata, &event.RewardedAt, &event.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return event, err
}

func (r *rewardRepository) GetTodayRewards(userID string) ([]models.RewardEvent, error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	query := `
		SELECT id, idempotency_key, user_id, stock_symbol, shares_quantity,
			   price_per_share, total_value, brokerage_fee, stt_fee, gst_fee,
			   exchange_fee, sebi_fee, total_fees, total_cost, reason, metadata,
			   rewarded_at, created_at
		FROM reward_events
		WHERE user_id = $1 AND rewarded_at >= $2 AND rewarded_at < $3
		ORDER BY rewarded_at DESC
	`

	rows, err := r.db.Query(query, userID, startOfDay, endOfDay)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.RewardEvent
	for rows.Next() {
		var event models.RewardEvent
		err := rows.Scan(
			&event.ID, &event.IdempotencyKey, &event.UserID, &event.StockSymbol,
			&event.SharesQuantity, &event.PricePerShare, &event.TotalValue,
			&event.BrokerageFee, &event.STTFee, &event.GSTFee, &event.ExchangeFee,
			&event.SEBIFee, &event.TotalFees, &event.TotalCost, &event.Reason,
			&event.Metadata, &event.RewardedAt, &event.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, rows.Err()
}

func (r *rewardRepository) GetHistoricalRewards(userID string) ([]models.RewardEvent, error) {
	now := time.Now()
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	query := `
		SELECT id, idempotency_key, user_id, stock_symbol, shares_quantity,
			   price_per_share, total_value, brokerage_fee, stt_fee, gst_fee,
			   exchange_fee, sebi_fee, total_fees, total_cost, reason, metadata,
			   rewarded_at, created_at
		FROM reward_events
		WHERE user_id = $1 AND rewarded_at < $2
		ORDER BY rewarded_at DESC
	`

	rows, err := r.db.Query(query, userID, startOfToday)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.RewardEvent
	for rows.Next() {
		var event models.RewardEvent
		err := rows.Scan(
			&event.ID, &event.IdempotencyKey, &event.UserID, &event.StockSymbol,
			&event.SharesQuantity, &event.PricePerShare, &event.TotalValue,
			&event.BrokerageFee, &event.STTFee, &event.GSTFee, &event.ExchangeFee,
			&event.SEBIFee, &event.TotalFees, &event.TotalCost, &event.Reason,
			&event.Metadata, &event.RewardedAt, &event.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, rows.Err()
}

func (r *rewardRepository) GetUserHolding(userID, stockSymbol string) (*models.UserHolding, error) {
	query := `
		SELECT id, user_id, stock_symbol, total_shares, average_price, last_updated
		FROM user_holdings
		WHERE user_id = $1 AND stock_symbol = $2
	`

	holding := &models.UserHolding{}
	err := r.db.QueryRow(query, userID, stockSymbol).Scan(
		&holding.ID, &holding.UserID, &holding.StockSymbol,
		&holding.TotalShares, &holding.AveragePrice, &holding.LastUpdated,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return holding, err
}

func (r *rewardRepository) UpsertUserHolding(holding *models.UserHolding) error {
	query := `
		INSERT INTO user_holdings (user_id, stock_symbol, total_shares, average_price, last_updated)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, stock_symbol)
		DO UPDATE SET
			total_shares = EXCLUDED.total_shares,
			average_price = EXCLUDED.average_price,
			last_updated = EXCLUDED.last_updated
		RETURNING id
	`

	return r.db.QueryRow(
		query,
		holding.UserID, holding.StockSymbol, holding.TotalShares,
		holding.AveragePrice, holding.LastUpdated,
	).Scan(&holding.ID)
}

func (r *rewardRepository) GetUserPortfolio(userID string) ([]models.UserHolding, error) {
	query := `
		SELECT id, user_id, stock_symbol, total_shares, average_price, last_updated
		FROM user_holdings
		WHERE user_id = $1 AND total_shares > 0
		ORDER BY stock_symbol
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var holdings []models.UserHolding
	for rows.Next() {
		var holding models.UserHolding
		err := rows.Scan(
			&holding.ID, &holding.UserID, &holding.StockSymbol,
			&holding.TotalShares, &holding.AveragePrice, &holding.LastUpdated,
		)
		if err != nil {
			return nil, err
		}
		holdings = append(holdings, holding)
	}

	return holdings, rows.Err()
}

// StockRepository handles stock-related database operations
type StockRepository interface {
	GetStockBySymbol(symbol string) (*models.Stock, error)
	CreateStockPrice(price *models.StockPrice) error
	GetLatestStockPrice(symbol string) (*models.StockPrice, error)
	GetLatestStockPrices() (map[string]float64, error)
}

type stockRepository struct {
	db *sql.DB
}

func NewStockRepository(db *sql.DB) StockRepository {
	return &stockRepository{db: db}
}

func (r *stockRepository) GetStockBySymbol(symbol string) (*models.Stock, error) {
	query := `
		SELECT id, symbol, company_name, exchange, is_active, created_at, updated_at
		FROM stocks
		WHERE symbol = $1
	`

	stock := &models.Stock{}
	err := r.db.QueryRow(query, symbol).Scan(
		&stock.ID, &stock.Symbol, &stock.CompanyName,
		&stock.Exchange, &stock.IsActive, &stock.CreatedAt, &stock.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("stock not found: %s", symbol)
	}

	return stock, err
}

func (r *stockRepository) CreateStockPrice(price *models.StockPrice) error {
	query := `
		INSERT INTO stock_prices (stock_symbol, price, timestamp, source)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (stock_symbol, timestamp) DO UPDATE
		SET price = EXCLUDED.price, source = EXCLUDED.source
		RETURNING id
	`

	return r.db.QueryRow(
		query,
		price.StockSymbol, price.Price, price.Timestamp, price.Source,
	).Scan(&price.ID)
}

func (r *stockRepository) GetLatestStockPrice(symbol string) (*models.StockPrice, error) {
	query := `
		SELECT id, stock_symbol, price, timestamp, source
		FROM stock_prices
		WHERE stock_symbol = $1
		ORDER BY timestamp DESC
		LIMIT 1
	`

	price := &models.StockPrice{}
	err := r.db.QueryRow(query, symbol).Scan(
		&price.ID, &price.StockSymbol, &price.Price, &price.Timestamp, &price.Source,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no price found for stock: %s", symbol)
	}

	return price, err
}

func (r *stockRepository) GetLatestStockPrices() (map[string]float64, error) {
	query := `
		SELECT DISTINCT ON (stock_symbol) stock_symbol, price
		FROM stock_prices
		ORDER BY stock_symbol, timestamp DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	prices := make(map[string]float64)
	for rows.Next() {
		var symbol string
		var price float64
		if err := rows.Scan(&symbol, &price); err != nil {
			return nil, err
		}
		prices[symbol] = price
	}

	return prices, rows.Err()
}

// LedgerRepository handles ledger operations
type LedgerRepository interface {
	CreateLedgerEntries(entries []models.LedgerEntry) error
}

type ledgerRepository struct {
	db *sql.DB
}

func NewLedgerRepository(db *sql.DB) LedgerRepository {
	return &ledgerRepository{db: db}
}

func (r *ledgerRepository) CreateLedgerEntries(entries []models.LedgerEntry) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO ledger_entries (
			entry_group_id, reward_event_id, account_type, stock_symbol,
			debit_amount, credit_amount, description
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	for _, entry := range entries {
		_, err := tx.Exec(
			query,
			entry.EntryGroupID, entry.RewardEventID, entry.AccountType, entry.StockSymbol,
			entry.DebitAmount, entry.CreditAmount, entry.Description,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
