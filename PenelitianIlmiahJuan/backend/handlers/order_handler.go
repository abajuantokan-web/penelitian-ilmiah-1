package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"openpeo-backend/config"
	"openpeo-backend/models"
)




func CreateOrder(c *gin.Context) {
	var req models.OrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request payload",
			"error":   err.Error(),
		})
		return
	}

	
	var customer models.User
	if err := config.DB.First(&customer, req.CustomerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Customer not found",
		})
		return
	}

	
	var product models.Product
	if err := config.DB.First(&product, req.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Product not found",
		})
		return
	}

	if !product.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "This product is no longer available for pre-order",
		})
		return
	}

	
	if req.Quantity < product.MinOrder {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Quantity does not meet the minimum order requirement",
			"data": gin.H{
				"min_order":          product.MinOrder,
				"requested_quantity": req.Quantity,
			},
		})
		return
	}

	
	if product.Stock < req.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Stok produk tidak mencukupi untuk pesanan Anda",
			"data": gin.H{
				"stock":              product.Stock,
				"requested_quantity": req.Quantity,
			},
		})
		return
	}

	
	product.Stock -= req.Quantity
	if err := config.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengurangi stok produk",
			"error":   err.Error(),
		})
		return
	}

	
	totalPrice := product.Price * float64(req.Quantity)

	
	order := models.Order{
		CustomerID: req.CustomerID,
		ProductID:  req.ProductID,
		Quantity:   req.Quantity,
		TotalPrice: totalPrice,
		Status:     "pending",
		Note:       req.Note,
	}

	if err := config.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create order",
			"error":   err.Error(),
		})
		return
	}

	
	config.DB.Preload("Customer").Preload("Product").Preload("Product.Vendor").First(&order, order.ID)

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Order placed successfully! Your pre-order has been registered.",
		"data":    order,
	})
}



func GetOrders(c *gin.Context) {
	var orders []models.Order
	var total int64

	query := config.DB.Model(&models.Order{})

	
	if customerID := c.Query("customer_id"); customerID != "" {
		query = query.Where("customer_id = ?", customerID)
	}

	
	if productID := c.Query("product_id"); productID != "" {
		query = query.Where("product_id = ?", productID)
	}

	
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	result := query.
		Preload("Customer").
		Preload("Product").
		Order("created_at DESC").
		Find(&orders)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch orders",
			"error":   result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    orders,
		"total":   total,
	})
}



func DeleteOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid order ID",
		})
		return
	}

	
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Parameter user_id (Admin ID) diperlukan untuk menghapus log transaksi",
		})
		return
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "User ID tidak valid",
		})
		return
	}
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "User pengakses tidak ditemukan",
		})
		return
	}
	if user.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Hanya Admin yang memiliki hak akses untuk menghapus log transaksi",
		})
		return
	}

	var order models.Order
	if err := config.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Order tidak ditemukan",
		})
		return
	}

	
	if err := config.DB.Unscoped().Delete(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal menghapus order",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Transaksi berhasil dihapus",
	})
}
