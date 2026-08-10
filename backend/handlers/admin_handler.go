package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"openpeo-backend/config"
	"openpeo-backend/models"
)

// GetAdminStats handles GET /api/admin/stats
func GetAdminStats(c *gin.Context) {
	var userCount, sellerCount, productCount int64
	var totalRevenue struct {
		Total float64
	}

	config.DB.Model(&models.User{}).Where("role = ?", "customer").Count(&userCount)
	config.DB.Model(&models.User{}).Where("role = ?", "seller").Count(&sellerCount)
	config.DB.Model(&models.Product{}).Count(&productCount)
	
	// Sum total_price of successful or pending orders (excluding cancelled)
	config.DB.Model(&models.Order{}).
		Select("sum(total_price) as total").
		Where("status != ?", "Dibatalkan Admin").
		Scan(&totalRevenue)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"total_users":    userCount,
			"total_sellers":  sellerCount,
			"total_products": productCount,
			"total_revenue":  totalRevenue.Total,
		},
	})
}

// GetAdminUsers handles GET /api/admin/users
func GetAdminUsers(c *gin.Context) {
	var users []models.User
	// Get all users (customers and sellers, excluding admins for this list typically)
	if err := config.DB.Preload("SellerProfile").Order("created_at DESC").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch users", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": users})
}



// GetAdminProducts handles GET /api/admin/products
func GetAdminProducts(c *gin.Context) {
	var products []models.Product
	if err := config.DB.Preload("SellerProfile").Order("created_at DESC").Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch products", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": products})
}

// ForceDeleteProduct handles DELETE /api/admin/products/:id
func ForceDeleteProduct(c *gin.Context) {
	productID := c.Param("id")

	if err := config.DB.Unscoped().Delete(&models.Product{}, productID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to delete product", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Produk berhasil dihapus permanen"})
}

// GetAdminOrders handles GET /api/admin/orders
func GetAdminOrders(c *gin.Context) {
	var orders []models.Order
	// Preload Customer, Product, and Product.Seller for comprehensive table view
	if err := config.DB.Preload("Customer").Preload("SellerProfile").Preload("OrderItems").Preload("OrderItems.Product").Order("created_at DESC").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch orders", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": orders})
}

// CancelOrder handles PUT /api/admin/orders/:id/cancel
func CancelOrder(c *gin.Context) {
	orderID := c.Param("id")

	if err := config.DB.Model(&models.Order{}).Where("id = ?", orderID).Update("status", "Dibatalkan Admin").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal membatalkan pesanan", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Pesanan berhasil dibatalkan"})
}

// GetActivityLogs handles GET /api/admin/activity-logs
func GetActivityLogs(c *gin.Context) {
	var logs []models.ActivityLog
	
	if err := config.DB.Order("created_at DESC").Limit(20).Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch logs", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": logs})
}
