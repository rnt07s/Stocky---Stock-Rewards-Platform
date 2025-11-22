package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stocky/assignment/internal/services"
)

type RewardHandler struct {
	rewardService services.RewardService
	log           *logrus.Logger
}

func NewRewardHandler(rewardService services.RewardService, log *logrus.Logger) *RewardHandler {
	return &RewardHandler{
		rewardService: rewardService,
		log:           log,
	}
}

// CreateReward handles POST /reward
func (h *RewardHandler) CreateReward(c *gin.Context) {
	var req services.RewardRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Errorf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}
	
	// If rewarded_at is not provided, set to current time
	if req.RewardedAt.IsZero() {
		req.RewardedAt = time.Now()
	}
	
	event, err := h.rewardService.CreateReward(&req)
	if err != nil {
		h.log.Errorf("Failed to create reward: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create reward",
			"details": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    event,
	})
}

// GetTodayStocks handles GET /today-stocks/:userId
func (h *RewardHandler) GetTodayStocks(c *gin.Context) {
	userID := c.Param("userId")
	
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_id is required",
		})
		return
	}
	
	rewards, err := h.rewardService.GetTodayStocks(userID)
	if err != nil {
		h.log.Errorf("Failed to get today's stocks: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch today's stocks",
			"details": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user_id": userID,
		"date":    time.Now().Format("2006-01-02"),
		"count":   len(rewards),
		"data":    rewards,
	})
}

// GetHistoricalINR handles GET /historical-inr/:userId
func (h *RewardHandler) GetHistoricalINR(c *gin.Context) {
	userID := c.Param("userId")
	
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_id is required",
		})
		return
	}
	
	historical, err := h.rewardService.GetHistoricalINR(userID)
	if err != nil {
		h.log.Errorf("Failed to get historical INR: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch historical INR data",
			"details": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    historical,
	})
}

// GetStats handles GET /stats/:userId
func (h *RewardHandler) GetStats(c *gin.Context) {
	userID := c.Param("userId")
	
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_id is required",
		})
		return
	}
	
	stats, err := h.rewardService.GetUserStats(userID)
	if err != nil {
		h.log.Errorf("Failed to get user stats: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch user stats",
			"details": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

// GetPortfolio handles GET /portfolio/:userId (bonus endpoint)
func (h *RewardHandler) GetPortfolio(c *gin.Context) {
	userID := c.Param("userId")
	
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_id is required",
		})
		return
	}
	
	portfolio, err := h.rewardService.GetUserPortfolio(userID)
	if err != nil {
		h.log.Errorf("Failed to get portfolio: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch portfolio",
			"details": err.Error(),
		})
		return
	}
	
	// Calculate total portfolio value
	totalValue := 0.0
	totalCost := 0.0
	for _, item := range portfolio {
		totalValue += item.CurrentValue
		totalCost += item.TotalCost
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user_id": userID,
		"summary": gin.H{
			"total_value":      totalValue,
			"total_cost":       totalCost,
			"total_profit_loss": totalValue - totalCost,
			"holdings_count":   len(portfolio),
		},
		"holdings": portfolio,
	})
}

// HealthCheck handles GET /health
func (h *RewardHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
