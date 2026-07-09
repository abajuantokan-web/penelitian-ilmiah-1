package handlers

import (
	"net/http"
	"openpeo-backend/config"
	"openpeo-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DevResetData handles POST /api/dev/reset-data
func DevResetData(c *gin.Context) {
	// 1. Delete all order items and orders
	if err := config.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.OrderItem{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to delete order_items"})
		return
	}
	if err := config.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Order{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to delete orders"})
		return
	}

	// 2. Delete all wallet transactions
	if err := config.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.WalletTransaction{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to delete wallet transactions"})
		return
	}

	// 3. Reset balances
	if err := config.DB.Model(&models.SellerProfile{}).Where("1 = 1").Updates(map[string]interface{}{
		"active_balance":  0,
		"pending_balance": 0,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to reset balances"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Database reset successfully"})
}
