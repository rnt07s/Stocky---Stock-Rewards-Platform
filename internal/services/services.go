package services

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stocky/assignment/internal/models"
	"github.com/stocky/assignment/internal/repository"
)

// StockPriceService handles stock price updates
type StockPriceService interface {
	StartPriceUpdater(intervalMinutes int)
	GetCurrentPrice(symbol string) (float64, error)
	GetAllCurrentPrices() (map[string]float64, error)
}

type stockPriceService struct {
	stockRepo repository.StockRepository
	log       *logrus.Logger
}

func NewStockPriceService(stockRepo repository.StockRepository, log *logrus.Logger) StockPriceService {
	return &stockPriceService{
		stockRepo: stockRepo,
		log:       log,
	}
}

// StartPriceUpdater starts hourly price update goroutine
func (s *stockPriceService) StartPriceUpdater(intervalMinutes int) {
	ticker := time.NewTicker(time.Duration(intervalMinutes) * time.Minute)
	
	// Run immediately on start
	s.updateAllPrices()
	
	go func() {
		for range ticker.C {
			s.updateAllPrices()
		}
	}()
	
	s.log.Infof("Stock price updater started (interval: %d minutes)", intervalMinutes)
}

func (s *stockPriceService) updateAllPrices() {
	s.log.Info("Updating stock prices...")
	
	// Hardcoded list of stocks (in production, fetch from database)
	stocks := []string{
		"RELIANCE", "TCS", "INFY", "HDFCBANK", "ICICIBANK",
		"HINDUNILVR", "ITC", "BHARTIARTL", "KOTAKBANK", "WIPRO",
	}
	
	now := time.Now()
	updatedCount := 0
	
	for _, symbol := range stocks {
		price := s.generateMockPrice(symbol)
		
		stockPrice := &models.StockPrice{
			StockSymbol: symbol,
			Price:       price,
			Timestamp:   now,
			Source:      "mock",
		}
		
		if err := s.stockRepo.CreateStockPrice(stockPrice); err != nil {
			s.log.Errorf("Failed to save price for %s: %v", symbol, err)
		} else {
			updatedCount++
		}
	}
	
	s.log.Infof("Updated %d stock prices at %s", updatedCount, now.Format(time.RFC3339))
}

// generateMockPrice generates realistic stock prices
func (s *stockPriceService) generateMockPrice(symbol string) float64 {
	// Base prices for Indian stocks (in INR)
	basePrices := map[string]float64{
		"RELIANCE":    2500.0,
		"TCS":         3500.0,
		"INFY":        1500.0,
		"HDFCBANK":    1600.0,
		"ICICIBANK":   950.0,
		"HINDUNILVR":  2400.0,
		"ITC":         450.0,
		"BHARTIARTL":  900.0,
		"KOTAKBANK":   1750.0,
		"WIPRO":       420.0,
	}
	
	basePrice := basePrices[symbol]
	if basePrice == 0 {
		basePrice = 1000.0 // Default
	}
	
	// Add random variation (+/- 5%)
	variation := (rand.Float64() - 0.5) * 0.10 // -5% to +5%
	price := basePrice * (1 + variation)
	
	// Round to 2 decimal places
	return float64(int(price*100)) / 100
}

func (s *stockPriceService) GetCurrentPrice(symbol string) (float64, error) {
	price, err := s.stockRepo.GetLatestStockPrice(symbol)
	if err != nil {
		return 0, err
	}
	return price.Price, nil
}

func (s *stockPriceService) GetAllCurrentPrices() (map[string]float64, error) {
	return s.stockRepo.GetLatestStockPrices()
}

// RewardService handles reward business logic
type RewardService interface {
	CreateReward(req *RewardRequest) (*models.RewardEvent, error)
	GetTodayStocks(userID string) ([]models.RewardEvent, error)
	GetHistoricalINR(userID string) (*HistoricalINRResponse, error)
	GetUserStats(userID string) (*UserStatsResponse, error)
	GetUserPortfolio(userID string) ([]models.PortfolioItem, error)
}

type rewardService struct {
	rewardRepo    repository.RewardRepository
	stockRepo     repository.StockRepository
	ledgerRepo    repository.LedgerRepository
	priceService  StockPriceService
	feesConfig    FeesConfig
	log           *logrus.Logger
}

type FeesConfig struct {
	BrokerageFeeBC int
	STTFeeBC       int
	GSTFeeBC       int
	ExchangeFeeBC  int
	SEBIFeeBC      int
}

type RewardRequest struct {
	IdempotencyKey string    `json:"idempotency_key" binding:"required"`
	UserID         string    `json:"user_id" binding:"required"`
	StockSymbol    string    `json:"stock_symbol" binding:"required"`
	SharesQuantity float64   `json:"shares_quantity" binding:"required,gt=0"`
	Reason         string    `json:"reason"`
	Metadata       string    `json:"metadata"`
	RewardedAt     time.Time `json:"rewarded_at"`
}

type HistoricalINRResponse struct {
	UserID     string                  `json:"user_id"`
	DailyINR   []DailyINR              `json:"daily_inr"`
	TotalValue float64                 `json:"total_value"`
}

type DailyINR struct {
	Date       string  `json:"date"`
	TotalValue float64 `json:"total_value"`
}

type UserStatsResponse struct {
	UserID              string                 `json:"user_id"`
	TodayRewards        []StockRewardSummary   `json:"today_rewards"`
	PortfolioValueINR   float64                `json:"portfolio_value_inr"`
	TotalSharesRewarded float64                `json:"total_shares_rewarded"`
}

type StockRewardSummary struct {
	StockSymbol string  `json:"stock_symbol"`
	TotalShares float64 `json:"total_shares"`
	RewardCount int     `json:"reward_count"`
}

func NewRewardService(
	rewardRepo repository.RewardRepository,
	stockRepo repository.StockRepository,
	ledgerRepo repository.LedgerRepository,
	priceService StockPriceService,
	feesConfig FeesConfig,
	log *logrus.Logger,
) RewardService {
	return &rewardService{
		rewardRepo:   rewardRepo,
		stockRepo:    stockRepo,
		ledgerRepo:   ledgerRepo,
		priceService: priceService,
		feesConfig:   feesConfig,
		log:          log,
	}
}

func (s *rewardService) CreateReward(req *RewardRequest) (*models.RewardEvent, error) {
	// Check idempotency
	existingEvent, err := s.rewardRepo.GetRewardEventByIdempotencyKey(req.IdempotencyKey)
	if err != nil {
		return nil, fmt.Errorf("idempotency check failed: %w", err)
	}
	if existingEvent != nil {
		s.log.Infof("Duplicate reward request detected (key: %s), returning existing event", req.IdempotencyKey)
		return existingEvent, nil
	}
	
	// Validate stock exists
	stock, err := s.stockRepo.GetStockBySymbol(req.StockSymbol)
	if err != nil {
		return nil, fmt.Errorf("invalid stock symbol: %w", err)
	}
	if !stock.IsActive {
		return nil, fmt.Errorf("stock %s is not active (possibly delisted)", req.StockSymbol)
	}
	
	// Get current stock price
	currentPrice, err := s.priceService.GetCurrentPrice(req.StockSymbol)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock price: %w", err)
	}
	
	// Calculate fees
	totalValue := req.SharesQuantity * currentPrice
	fees := s.calculateFees(totalValue)
	
	// Set rewarded_at to now if not provided
	rewardedAt := req.RewardedAt
	if rewardedAt.IsZero() {
		rewardedAt = time.Now()
	}
	
	// Create reward event
	event := &models.RewardEvent{
		IdempotencyKey: req.IdempotencyKey,
		UserID:         req.UserID,
		StockSymbol:    req.StockSymbol,
		SharesQuantity: req.SharesQuantity,
		PricePerShare:  currentPrice,
		TotalValue:     totalValue,
		BrokerageFee:   fees.Brokerage,
		STTFee:         fees.STT,
		GSTFee:         fees.GST,
		ExchangeFee:    fees.Exchange,
		SEBIFee:        fees.SEBI,
		TotalFees:      fees.Total,
		TotalCost:      totalValue + fees.Total,
		Reason:         req.Reason,
		Metadata:       sql.NullString{String: req.Metadata, Valid: req.Metadata != ""},
		RewardedAt:     rewardedAt,
	}
	
	// Save reward event
	if err := s.rewardRepo.CreateRewardEvent(event); err != nil {
		return nil, fmt.Errorf("failed to create reward event: %w", err)
	}
	
	// Update user holdings
	if err := s.updateUserHoldings(event); err != nil {
		s.log.Errorf("Failed to update user holdings: %v", err)
		// Don't fail the entire operation, but log the error
	}
	
	// Create ledger entries
	if err := s.createLedgerEntries(event); err != nil {
		s.log.Errorf("Failed to create ledger entries: %v", err)
		// Don't fail the entire operation, but log the error
	}
	
	s.log.Infof("Reward created: user=%s, stock=%s, shares=%.6f, price=%.2f, total_cost=%.2f",
		event.UserID, event.StockSymbol, event.SharesQuantity, event.PricePerShare, event.TotalCost)
	
	return event, nil
}

type Fees struct {
	Brokerage float64
	STT       float64
	GST       float64
	Exchange  float64
	SEBI      float64
	Total     float64
}

func (s *rewardService) calculateFees(totalValue float64) Fees {
	brokerage := totalValue * float64(s.feesConfig.BrokerageFeeBC) / 10000
	stt := totalValue * float64(s.feesConfig.STTFeeBC) / 10000
	exchange := totalValue * float64(s.feesConfig.ExchangeFeeBC) / 10000
	sebi := totalValue * float64(s.feesConfig.SEBIFeeBC) / 10000
	gst := brokerage * float64(s.feesConfig.GSTFeeBC) / 100
	
	total := brokerage + stt + gst + exchange + sebi
	
	return Fees{
		Brokerage: roundToDecimal(brokerage, 4),
		STT:       roundToDecimal(stt, 4),
		GST:       roundToDecimal(gst, 4),
		Exchange:  roundToDecimal(exchange, 4),
		SEBI:      roundToDecimal(sebi, 4),
		Total:     roundToDecimal(total, 4),
	}
}

func (s *rewardService) updateUserHoldings(event *models.RewardEvent) error {
	// Get existing holding
	holding, err := s.rewardRepo.GetUserHolding(event.UserID, event.StockSymbol)
	if err != nil {
		return err
	}
	
	if holding == nil {
		// Create new holding
		holding = &models.UserHolding{
			UserID:       event.UserID,
			StockSymbol:  event.StockSymbol,
			TotalShares:  event.SharesQuantity,
			AveragePrice: event.PricePerShare,
			LastUpdated:  time.Now(),
		}
	} else {
		// Update existing holding with weighted average price
		totalCost := (holding.TotalShares * holding.AveragePrice) + (event.SharesQuantity * event.PricePerShare)
		holding.TotalShares += event.SharesQuantity
		holding.AveragePrice = totalCost / holding.TotalShares
		holding.LastUpdated = time.Now()
	}
	
	return s.rewardRepo.UpsertUserHolding(holding)
}

func (s *rewardService) createLedgerEntries(event *models.RewardEvent) error {
	entryGroupID := fmt.Sprintf("reward-%d", event.ID)
	
	entries := []models.LedgerEntry{
		// Debit: Stock Inventory (asset increases)
		{
			EntryGroupID:  entryGroupID,
			RewardEventID: sql.NullInt64{Int64: event.ID, Valid: true},
			AccountType:   "stock_inventory",
			StockSymbol:   sql.NullString{String: event.StockSymbol, Valid: true},
			DebitAmount:   event.TotalValue,
			CreditAmount:  0,
			Description:   fmt.Sprintf("Stock reward: %s x %.6f shares to user %s", event.StockSymbol, event.SharesQuantity, event.UserID),
		},
		// Credit: Cash (asset decreases)
		{
			EntryGroupID:  entryGroupID,
			RewardEventID: sql.NullInt64{Int64: event.ID, Valid: true},
			AccountType:   "cash_outflow",
			StockSymbol:   sql.NullString{},
			DebitAmount:   0,
			CreditAmount:  event.TotalValue,
			Description:   fmt.Sprintf("Cash paid for stock purchase"),
		},
		// Debit: Fees Expense
		{
			EntryGroupID:  entryGroupID,
			RewardEventID: sql.NullInt64{Int64: event.ID, Valid: true},
			AccountType:   "fees_expense",
			StockSymbol:   sql.NullString{},
			DebitAmount:   event.TotalFees,
			CreditAmount:  0,
			Description:   fmt.Sprintf("Fees: brokerage=%.2f, STT=%.2f, GST=%.2f, exchange=%.2f, SEBI=%.2f", event.BrokerageFee, event.STTFee, event.GSTFee, event.ExchangeFee, event.SEBIFee),
		},
		// Credit: Cash (for fees)
		{
			EntryGroupID:  entryGroupID,
			RewardEventID: sql.NullInt64{Int64: event.ID, Valid: true},
			AccountType:   "cash_outflow",
			StockSymbol:   sql.NullString{},
			DebitAmount:   0,
			CreditAmount:  event.TotalFees,
			Description:   "Cash paid for transaction fees",
		},
	}
	
	return s.ledgerRepo.CreateLedgerEntries(entries)
}

func (s *rewardService) GetTodayStocks(userID string) ([]models.RewardEvent, error) {
	return s.rewardRepo.GetTodayRewards(userID)
}

func (s *rewardService) GetHistoricalINR(userID string) (*HistoricalINRResponse, error) {
	rewards, err := s.rewardRepo.GetHistoricalRewards(userID)
	if err != nil {
		return nil, err
	}
	
	// Group by date and sum values
	dailyMap := make(map[string]float64)
	for _, reward := range rewards {
		dateStr := reward.RewardedAt.Format("2006-01-02")
		dailyMap[dateStr] += reward.TotalValue
	}
	
	// Convert to slice and sort
	var dailyINR []DailyINR
	totalValue := 0.0
	for date, value := range dailyMap {
		dailyINR = append(dailyINR, DailyINR{
			Date:       date,
			TotalValue: roundToDecimal(value, 2),
		})
		totalValue += value
	}
	
	return &HistoricalINRResponse{
		UserID:     userID,
		DailyINR:   dailyINR,
		TotalValue: roundToDecimal(totalValue, 2),
	}, nil
}

func (s *rewardService) GetUserStats(userID string) (*UserStatsResponse, error) {
	// Get today's rewards
	todayRewards, err := s.rewardRepo.GetTodayRewards(userID)
	if err != nil {
		return nil, err
	}
	
	// Group by stock symbol
	stockMap := make(map[string]*StockRewardSummary)
	totalShares := 0.0
	
	for _, reward := range todayRewards {
		if summary, exists := stockMap[reward.StockSymbol]; exists {
			summary.TotalShares += reward.SharesQuantity
			summary.RewardCount++
		} else {
			stockMap[reward.StockSymbol] = &StockRewardSummary{
				StockSymbol: reward.StockSymbol,
				TotalShares: reward.SharesQuantity,
				RewardCount: 1,
			}
		}
		totalShares += reward.SharesQuantity
	}
	
	var summaries []StockRewardSummary
	for _, summary := range stockMap {
		summaries = append(summaries, *summary)
	}
	
	// Calculate portfolio value
	portfolio, err := s.GetUserPortfolio(userID)
	if err != nil {
		return nil, err
	}
	
	portfolioValue := 0.0
	for _, item := range portfolio {
		portfolioValue += item.CurrentValue
	}
	
	return &UserStatsResponse{
		UserID:              userID,
		TodayRewards:        summaries,
		PortfolioValueINR:   roundToDecimal(portfolioValue, 2),
		TotalSharesRewarded: roundToDecimal(totalShares, 6),
	}, nil
}

func (s *rewardService) GetUserPortfolio(userID string) ([]models.PortfolioItem, error) {
	holdings, err := s.rewardRepo.GetUserPortfolio(userID)
	if err != nil {
		return nil, err
	}
	
	// Get current prices
	prices, err := s.priceService.GetAllCurrentPrices()
	if err != nil {
		return nil, err
	}
	
	var portfolio []models.PortfolioItem
	for _, holding := range holdings {
		currentPrice := prices[holding.StockSymbol]
		if currentPrice == 0 {
			// Fallback if price not found
			currentPrice = holding.AveragePrice
		}
		
		currentValue := holding.TotalShares * currentPrice
		totalCost := holding.TotalShares * holding.AveragePrice
		profitLoss := currentValue - totalCost
		profitLossPct := 0.0
		if totalCost > 0 {
			profitLossPct = (profitLoss / totalCost) * 100
		}
		
		// Get company name
		stock, _ := s.stockRepo.GetStockBySymbol(holding.StockSymbol)
		companyName := holding.StockSymbol
		if stock != nil {
			companyName = stock.CompanyName
		}
		
		portfolio = append(portfolio, models.PortfolioItem{
			StockSymbol:   holding.StockSymbol,
			CompanyName:   companyName,
			TotalShares:   roundToDecimal(holding.TotalShares, 6),
			AveragePrice:  roundToDecimal(holding.AveragePrice, 2),
			CurrentPrice:  roundToDecimal(currentPrice, 2),
			CurrentValue:  roundToDecimal(currentValue, 2),
			TotalCost:     roundToDecimal(totalCost, 2),
			ProfitLoss:    roundToDecimal(profitLoss, 2),
			ProfitLossPct: roundToDecimal(profitLossPct, 2),
		})
	}
	
	return portfolio, nil
}

func roundToDecimal(value float64, decimals int) float64 {
	multiplier := 1.0
	for i := 0; i < decimals; i++ {
		multiplier *= 10
	}
	return float64(int(value*multiplier+0.5)) / multiplier
}
