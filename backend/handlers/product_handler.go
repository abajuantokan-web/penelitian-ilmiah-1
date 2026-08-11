package handlers

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"openpeo-backend/config"
	"openpeo-backend/models"
)

func CreateProduct(c *gin.Context) {
	var product models.Product

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":	false,
			"message":	"Invalid request payload",
			"error":	err.Error(),
		})
		return
	}

	if product.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":	false,
			"message":	"Product name is required",
		})
		return
	}

	if product.SellerID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":	false,
			"message":	"Seller ID is required",
		})
		return
	}

	if product.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":	false,
			"message":	"Price must be greater than zero",
		})
		return
	}

	if product.Region == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":	false,
			"message":	"Region is required (e.g., Sumba, Manggarai, Kupang, Flores, Timor)",
		})
		return
	}

	if product.MinOrder < 1 {
		product.MinOrder = 1
	}
	if product.PreOrderDuration < 1 {
		product.PreOrderDuration = 7
	}
	product.IsActive = true

	var seller models.User
	if err := config.DB.First(&seller, product.SellerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success":	false,
			"message":	"Seller not found",
		})
		return
	}
	if seller.Role != "admin" && seller.Role != "seller" {
		c.JSON(http.StatusForbidden, gin.H{
			"success":	false,
			"message":	"Hanya Seller/Admin yang memiliki hak akses untuk mengelola produk",
		})
		return
	}

	if err := config.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":	false,
			"message":	"Failed to create product",
			"error":	err.Error(),
		})
		return
	}

	config.DB.Preload("SellerProfile").First(&product, product.ID)

	c.JSON(http.StatusCreated, gin.H{
		"success":	true,
		"message":	"Product created successfully",
		"data":		product,
	})
}

func GetProducts(c *gin.Context) {
	var products []models.Product
	var total int64

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "12"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 12
	}
	offset := (page - 1) * limit

	query := config.DB.Model(&models.Product{}).Where("is_active = ?", true)

	if region := c.Query("region"); region != "" {
		query = query.Where("region = ?", region)
	}

	if category := c.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}

	if search := c.Query("search"); search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	query.Count(&total)

	result := query.
		Preload("SellerProfile").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&products)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":	false,
			"message":	"Failed to fetch products",
			"error":	result.Error.Error(),
		})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	c.JSON(http.StatusOK, gin.H{
		"success":	true,
		"data":		products,
		"meta": gin.H{
			"page":		page,
			"limit":	limit,
			"total":	total,
			"total_pages":	totalPages,
		},
	})
}

func GetProductByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":	false,
			"message":	"Invalid product ID",
		})
		return
	}

	var product models.Product
	if err := config.DB.Preload("SellerProfile").First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success":	false,
			"message":	"Product not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":	true,
		"data":		product,
	})
}

func UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":	false,
			"message":	"Invalid product ID",
		})
		return
	}

	var product models.Product
	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success":	false,
			"message":	"Product not found",
		})
		return
	}

	var req models.Product
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":	false,
			"message":	"Invalid request payload",
			"error":	err.Error(),
		})
		return
	}

	var adminUser models.User
	userIDToCheck := req.SellerID
	if userIDToCheck == 0 {
		userIDToCheck = product.SellerID
	}
	if err := config.DB.First(&adminUser, userIDToCheck).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"success":	false,
			"message":	"User pengakses tidak ditemukan",
		})
		return
	}
	if adminUser.Role != "admin" && adminUser.Role != "seller" {
		c.JSON(http.StatusForbidden, gin.H{
			"success":	false,
			"message":	"Hanya Admin/Seller yang memiliki hak akses untuk mengelola produk",
		})
		return
	}

	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	product.Description = req.Description
	if req.Category != "" {
		product.Category = req.Category
	}
	if req.Region != "" {
		product.Region = req.Region
	}
	if req.PreOrderDuration > 0 {
		product.PreOrderDuration = req.PreOrderDuration
	}
	if req.MinOrder > 0 {
		product.MinOrder = req.MinOrder
	}
	product.Stock = req.Stock
	if req.ImageURL != "" {
		product.ImageURL = req.ImageURL
	}
	product.IsActive = req.IsActive

	if err := config.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":	false,
			"message":	"Failed to update product",
			"error":	err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":	true,
		"message":	"Product updated successfully",
		"data":		product,
	})
}

func DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":	false,
			"message":	"Invalid product ID",
		})
		return
	}

	var product models.Product
	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success":	false,
			"message":	"Product not found",
		})
		return
	}

	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":	false,
			"message":	"Parameter user_id (Admin ID) diperlukan untuk menghapus produk",
		})
		return
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":	false,
			"message":	"User ID tidak valid",
		})
		return
	}
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"success":	false,
			"message":	"User pengakses tidak ditemukan",
		})
		return
	}
	if user.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"success":	false,
			"message":	"Hanya Admin yang memiliki hak akses untuk mengelola produk",
		})
		return
	}

	if err := config.DB.Unscoped().Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":	false,
			"message":	"Failed to delete product",
			"error":	err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":	true,
		"message":	"Product deleted successfully",
	})
}
