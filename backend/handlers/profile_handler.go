package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"openpeo-backend/config"
	"openpeo-backend/models"
)

// UpdateProfileRequest is the expected JSON payload for profile updates.
type UpdateProfileRequest struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

// ChangePasswordRequest is the expected JSON payload for password changes.
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}

// GetProfile handles GET /api/user/profile
// Returns the authenticated user's profile from JWT context.
func GetProfile(c *gin.Context) {
	userIDFloat, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}
	userID := int32(userIDFloat.(float64))

	var user models.User
	if err := config.DB.Preload("SellerProfile").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "User tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"id":       user.ID,
			"name":     user.Name,
			"email":    user.Email,
			"phone":    user.Phone,
			"address":  user.Address,
			"role":     user.Role,
			"avatar":   user.Avatar,
		"store_name": func() string {
				if user.Role == "seller" && user.SellerProfile != nil { return user.SellerProfile.StoreName }
				return ""
			}(),
			"store_logo": func() string {
				if user.Role == "seller" && user.SellerProfile != nil { return user.SellerProfile.StoreLogo }
				return ""
			}(),
		},
	})
}

// UpdateProfile handles PUT /api/user/profile
// Updates name, phone, and address for the authenticated user.
func UpdateProfile(c *gin.Context) {
	userIDFloat, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}
	userID := int32(userIDFloat.(float64))

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data tidak valid",
			"error":   err.Error(),
		})
		return
	}

	var user models.User
	if err := config.DB.Preload("SellerProfile").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "User tidak ditemukan",
		})
		return
	}

	// Update fields
	if req.Name != "" {
		user.Name = req.Name
	}
	user.Phone = req.Phone
	user.Address = req.Address

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal memperbarui profil",
			"error":   err.Error(),
		})
		return
	}

	// Update localStorage data on frontend by returning fresh user object
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profil berhasil diperbarui",
		"data": gin.H{
			"id":       user.ID,
			"name":     user.Name,
			"email":    user.Email,
			"phone":    user.Phone,
			"address":  user.Address,
			"role":     user.Role,
			"avatar":   user.Avatar,
		},
	})
}

// ChangePassword handles PUT /api/user/password
// Validates current password and updates to new password using bcrypt.
func ChangePassword(c *gin.Context) {
	userIDFloat, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}
	userID := int32(userIDFloat.(float64))

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data tidak valid. Password baru minimal 6 karakter.",
			"error":   err.Error(),
		})
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "User tidak ditemukan",
		})
		return
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Password saat ini salah",
		})
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal memproses password baru",
		})
		return
	}

	user.Password = string(hashedPassword)
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengubah password",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Password berhasil diubah",
	})
}

// GetMyOrders handles GET /api/user/orders
// Returns the order history for the authenticated user, with product preloading.
func GetMyOrders(c *gin.Context) {
	userIDFloat, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}
	userID := int32(userIDFloat.(float64))

	var orders []models.Order
	if err := config.DB.
		Where("customer_id = ?", userID).
		Preload("SellerProfile").
		Preload("OrderItems").
		Preload("OrderItems.Product").
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil riwayat pesanan",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    orders,
		"total":   len(orders),
	})
}
