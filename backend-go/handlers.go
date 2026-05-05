package main

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	isShieldActive  bool = false
	idempotencyKeys sync.Map
)

func ToggleShield(c *gin.Context) {
	isShieldActive = !isShieldActive
	status := "INACTIVE"
	if isShieldActive {
		status = "ACTIVE"
	}
	c.JSON(http.StatusOK, gin.H{"message": "Shield is now " + status, "active": isShieldActive})
}

func CreateData(c *gin.Context) {
	var reqData UserData
	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if isShieldActive {
		// Protected: Check Idempotency Key
		idemKey := c.GetHeader("X-Idempotency-Key")
		if idemKey == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "X-Idempotency-Key header is required"})
			return
		}

		if _, exists := idempotencyKeys.Load(idemKey); exists {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Duplicate request detected"})
			return
		}
		idempotencyKeys.Store(idemKey, true)
	}

	// Insert Data
	if err := DB.Create(&reqData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert data"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Data created successfully", "data": reqData})
}

func UpdateData(c *gin.Context) {
	id := c.Param("id")

	var reqData struct {
		Data string `json:"data"`
	}
	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if isShieldActive {
		// Protected: Check Authorization matches Username (Mitigate IDOR)
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		var existingRecord UserData
		if err := DB.First(&existingRecord, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
			return
		}

		if existingRecord.Username != authHeader {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: You do not own this record"})
			return
		}

		existingRecord.Data = reqData.Data
		if err := DB.Save(&existingRecord).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update record"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Record updated securely", "data": existingRecord})

	} else {
		// Vulnerable: IDOR update directly without ownership check
		if err := DB.Model(&UserData{}).Where("id = ?", id).Update("data", reqData.Data).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update record"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Record updated (Vulnerable to IDOR)"})
	}
}
