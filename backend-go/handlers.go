package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var ShieldActive bool = false
var idempotencyKeys sync.Map
var limiter = rate.NewLimiter(rate.Every(1*time.Second), 5) // 5 requests per second

type DataRequest struct {
	Username string `json:"username"`
	Data     string `json:"data"`
}

func ToggleShield(c *gin.Context) {
	ShieldActive = !ShieldActive
	status := "INACTIVE"
	if ShieldActive {
		status = "ACTIVE"
	}
	c.JSON(http.StatusOK, gin.H{"message": "Shield is now " + status, "active": ShieldActive})
}

func CreateData(c *gin.Context) {
	var req DataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if ShieldActive {
		// SHIELD ON: Enforce rate limiting and Idempotency Key
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded. Try again later."})
			return
		}

		idempKey := c.GetHeader("Idempotency-Key")
		if idempKey == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Idempotency-Key header is required when shield is active"})
			return
		}

		if _, exists := idempotencyKeys.Load(idempKey); exists {
			c.JSON(http.StatusConflict, gin.H{"error": "Duplicate request detected based on Idempotency-Key"})
			return
		}
		idempotencyKeys.Store(idempKey, true)

		// Secure Insert via GORM
		newData := UserData{Username: req.Username, Data: req.Data}
		if DB == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database not initialized"})
			return
		}
		if err := DB.Create(&newData).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Data inserted securely", "data": newData})

	} else {
		// SHIELD OFF: Vulnerable to DoS (no limits) and Duplication (no idempotency check)
		if SqlDB == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database not initialized"})
			return
		}
		// Raw SQL execution without rate limits
		query := fmt.Sprintf("INSERT INTO user_data (username, data) VALUES ('%s', '%s')", req.Username, req.Data)
		_, err := SqlDB.Exec(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Data inserted vulnerable to DoS/Duplication"})
	}
}

func UpdateData(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Data string `json:"data"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if ShieldActive {
		// SHIELD ON: Check Authorization (Mitigate IDOR)
		authUsername := c.GetHeader("Authorization")
		if authUsername == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required (username)"})
			return
		}

		if DB == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database not initialized"})
			return
		}

		var existingData UserData
		if err := DB.First(&existingData, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
			return
		}

		// Ensure the logged-in user owns the record they are trying to update
		if existingData.Username != authUsername {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: You do not own this record"})
			return
		}

		// Update securely via GORM
		existingData.Data = req.Data
		if err := DB.Save(&existingData).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update record"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Record updated securely", "data": existingData})

	} else {
		// SHIELD OFF: IDOR Vulnerability - updates record purely based on URL parameter ID
		if SqlDB == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database not initialized"})
			return
		}
		query := fmt.Sprintf("UPDATE user_data SET data = '%s' WHERE id = %s", req.Data, id)
		_, err := SqlDB.Exec(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Record updated (Vulnerable to IDOR)"})
	}
}
