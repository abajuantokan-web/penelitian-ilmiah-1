package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"openpeo-backend/config"
	"openpeo-backend/models"
)

// requireSeller validates that the current user has the 'seller' role
func requireSeller(c *gin.Context) (int32, bool) {
	role, exists := c.Get("role")
	if !exists || role != "seller" {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Akses ditolak. Anda bukan penjual.",
		})
		return 0, false
	}

	userIDFloat, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Sesi tidak valid",
		})
		return 0, false
	}
	return int32(userIDFloat.(float64)), true
}

// GetSellerProducts handles GET /api/seller/products
func GetSellerProducts(c *gin.Context) {
	sellerID, ok := requireSeller(c)
	if !ok {
		return
	}

	var products []models.Product
	if err := config.DB.Where("seller_id = ?", sellerID).Order("created_at desc").Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengambil produk", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": products})
}

// CreateSellerProduct handles POST /api/seller/products
func CreateSellerProduct(c *gin.Context) {
	sellerID, ok := requireSeller(c)
	if !ok {
		return
	}

	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Payload tidak valid"})
		return
	}

	product.SellerID = sellerID
	product.IsActive = true

	if err := config.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal menambah produk"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Produk berhasil ditambahkan", "data": product})
}

// UpdateSellerProduct handles PUT /api/seller/products/:id
func UpdateSellerProduct(c *gin.Context) {
	sellerID, ok := requireSeller(c)
	if !ok {
		return
	}

	productID := c.Param("id")

	var existing models.Product
	if err := config.DB.First(&existing, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Produk tidak ditemukan"})
		return
	}

	// Security check: ensure the product belongs to this seller
	if existing.SellerID != sellerID {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Anda tidak memiliki akses untuk mengubah produk ini"})
		return
	}

	var updateData models.Product
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Payload tidak valid"})
		return
	}

	// Update specific fields
	existing.Name = updateData.Name
	existing.Description = updateData.Description
	existing.Price = updateData.Price
	existing.ImageURL = updateData.ImageURL
	existing.Category = updateData.Category
	existing.Region = updateData.Region
	existing.MinOrder = updateData.MinOrder
	existing.PreOrderDuration = updateData.PreOrderDuration
	existing.Stock = updateData.Stock

	if err := config.DB.Save(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengubah produk"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Produk berhasil diubah", "data": existing})
}

// DeleteSellerProduct handles DELETE /api/seller/products/:id
func DeleteSellerProduct(c *gin.Context) {
	sellerID, ok := requireSeller(c)
	if !ok {
		return
	}

	productID, _ := strconv.Atoi(c.Param("id"))

	var existing models.Product
	if err := config.DB.First(&existing, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Produk tidak ditemukan"})
		return
	}

	if existing.SellerID != sellerID {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Anda tidak memiliki akses untuk menghapus produk ini"})
		return
	}

	if err := config.DB.Delete(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal menghapus produk"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Produk berhasil dihapus"})
}
