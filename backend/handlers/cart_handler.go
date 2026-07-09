package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"openpeo-backend/config"
	"openpeo-backend/models"
)

type AddCartRequest struct {
	ProductID int32 `json:"product_id" binding:"required"`
	Quantity  int   `json:"quantity" binding:"required,min=1"`
}

// AddToCart handles POST /api/cart
func AddToCart(c *gin.Context) {
	// User ID from JWT middleware
	userIDFloat, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}
	userID := int32(userIDFloat.(float64))

	var req AddCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request", "error": err.Error()})
		return
	}

	// Verify product exists
	var product models.Product
	if err := config.DB.First(&product, req.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Product not found"})
		return
	}

	// Check if item already in cart
	var cartItem models.CartItem
	err := config.DB.Where("user_id = ? AND product_id = ?", userID, req.ProductID).First(&cartItem).Error
	if err == nil {
		// Update quantity
		cartItem.Quantity += req.Quantity
		config.DB.Save(&cartItem)
	} else {
		// Add new item
		cartItem = models.CartItem{
			UserID:    userID,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
		}
		config.DB.Create(&cartItem)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Product added to cart",
		"data":    cartItem,
	})
}

// GetCart handles GET /api/cart
func GetCart(c *gin.Context) {
	userIDFloat, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}
	userID := int32(userIDFloat.(float64))

	var cartItems []models.CartItem
	if err := config.DB.Preload("Product").Where("user_id = ?", userID).Find(&cartItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    cartItems,
	})
}

// RemoveFromCart handles DELETE /api/cart/:id
func RemoveFromCart(c *gin.Context) {
	userIDFloat, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}
	userID := int32(userIDFloat.(float64))

	cartItemIDStr := c.Param("id")
	cartItemID, err := strconv.ParseInt(cartItemIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid cart item ID"})
		return
	}

	var cartItem models.CartItem
	if err := config.DB.Where("id = ? AND user_id = ?", cartItemID, userID).First(&cartItem).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Cart item not found"})
		return
	}

	if err := config.DB.Delete(&cartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to delete cart item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Cart item removed",
	})
}
