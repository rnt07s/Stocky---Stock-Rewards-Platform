package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Fees     FeesConfig
	Service  ServiceConfig
}

type ServerConfig struct {
	Port    string
	GinMode string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type FeesConfig struct {
	BrokerageFeeBP int // Basis points (1 bp = 0.01%)
	STTFeeBC       int
	GSTFeeBC       int
	ExchangeFeeBC  int
	SEBIFeeBC      int
}

type ServiceConfig struct {
	PriceUpdateIntervalMinutes int
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if exists (ignore error if not found)
	_ = godotenv.Load()

	cfg := &Config{
		Server: ServerConfig{
			Port:    getEnv("PORT", "8080"),
			GinMode: getEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "stocky_user"),
			Password: getEnv("DB_PASSWORD", "stocky_password"),
			DBName:   getEnv("DB_NAME", "assignment"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Fees: FeesConfig{
			BrokerageFeeBC: getEnvAsInt("BROKERAGE_FEE_BP", 5),
			STTFeeBC:       getEnvAsInt("STT_FEE_BP", 25),
			GSTFeeBC:       getEnvAsInt("GST_FEE_BP", 18),
			ExchangeFeeBC:  getEnvAsInt("EXCHANGE_FEE_BP", 3),
			SEBIFeeBC:      getEnvAsInt("SEBI_FEE_BP", 1),
		},
		Service: ServiceConfig{
			PriceUpdateIntervalMinutes: getEnvAsInt("PRICE_UPDATE_INTERVAL_MINUTES", 60),
		},
	}

	return cfg, nil
}

// GetDSN returns PostgreSQL connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		logrus.Warnf("Invalid integer for %s, using default %d", key, defaultValue)
		return defaultValue
	}
	return value
}
