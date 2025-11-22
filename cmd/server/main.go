package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stocky/assignment/internal/config"
	"github.com/stocky/assignment/internal/database"
	"github.com/stocky/assignment/internal/handlers"
	"github.com/stocky/assignment/internal/middleware"
	"github.com/stocky/assignment/internal/repository"
	"github.com/stocky/assignment/internal/services"
)

func main() {
	// Initialize logger
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)
	log.SetOutput(os.Stdout)

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set Gin mode
	gin.SetMode(cfg.Server.GinMode)

	// Connect to database
	db, err := database.Connect(cfg.Database.GetDSN(), log)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := database.RunMigrations(db, "migrations", log); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	rewardRepo := repository.NewRewardRepository(db)
	stockRepo := repository.NewStockRepository(db)
	ledgerRepo := repository.NewLedgerRepository(db)

	// Initialize services
	priceService := services.NewStockPriceService(stockRepo, log)
	rewardService := services.NewRewardService(
		rewardRepo,
		stockRepo,
		ledgerRepo,
		priceService,
		services.FeesConfig{
			BrokerageFeeBC: cfg.Fees.BrokerageFeeBC,
			STTFeeBC:       cfg.Fees.STTFeeBC,
			GSTFeeBC:       cfg.Fees.GSTFeeBC,
			ExchangeFeeBC:  cfg.Fees.ExchangeFeeBC,
			SEBIFeeBC:      cfg.Fees.SEBIFeeBC,
		},
		log,
	)

	// Start stock price updater
	priceService.StartPriceUpdater(cfg.Service.PriceUpdateIntervalMinutes)

	// Initialize handlers
	rewardHandler := handlers.NewRewardHandler(rewardService, log)

	// Setup router
	router := gin.New()
	router.Use(middleware.RecoveryMiddleware(log))
	router.Use(middleware.LoggingMiddleware(log))
	router.Use(middleware.CORSMiddleware())

	// Health check
	router.GET("/health", rewardHandler.HealthCheck)

	// API routes
	api := router.Group("/api/v1")
	{
		api.POST("/reward", rewardHandler.CreateReward)
		api.GET("/today-stocks/:userId", rewardHandler.GetTodayStocks)
		api.GET("/historical-inr/:userId", rewardHandler.GetHistoricalINR)
		api.GET("/stats/:userId", rewardHandler.GetStats)
		api.GET("/portfolio/:userId", rewardHandler.GetPortfolio)
	}

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Infof("Starting Stocky API server on %s", addr)
	
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
