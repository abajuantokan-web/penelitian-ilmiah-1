package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"openpeo-backend/config"
	"openpeo-backend/middleware"
	"openpeo-backend/models"
)

// LoginRequest is the expected JSON payload for the login endpoint.
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest represents the JSON payload for registration.
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
}

// Login handles POST /api/login
// Authenticates a user by matching username and password against the users table.
// Returns the user's id, username, name, and role on success.
//
// In a production system this should use bcrypt password hashing and JWT tokens.
// For this academic/demo project, passwords are compared as plain text to match
// the dummy data in db.sql.
func Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Email dan password harus diisi",
			"error":   err.Error(),
		})
		return
	}

	// Look up user by email
	var user models.User
	result := config.DB.Preload("SellerProfile").Where("email = ?", req.Email).First(&user)

	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Email tidak ditemukan",
		})
		return
	}

	// Validate password using bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Password salah",
		})
		return
	}

	// Generate JWT
	token, err := middleware.GenerateJWT(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal membuat token",
		})
		return
	}

	logEntry := models.ActivityLog{
		UserID:      user.ID,
		UserRole:    user.Role,
		ActionType:  "LOGIN",
		Description: "User Logged In: " + user.Email,
		CreatedAt:   time.Now(),
	}
	config.DB.Create(&logEntry)

	// Return user session data
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login berhasil!",
		"token":   token,
		"data": gin.H{
			"id":       user.ID,
			"name":     user.Name,
			"role":     user.Role,
			"email":    user.Email,
			"avatar":   user.Avatar,
			"store_name": func() string { if user.SellerProfile != nil && user.SellerProfile.StoreName != "" { return user.SellerProfile.StoreName }; return user.StoreName }(),
			"store_logo": func() string { if user.SellerProfile != nil { return user.SellerProfile.StoreLogo }; return "" }(),
		},
	})
}

// RegisterUser handles POST /api/register
// Registers a new customer into the system.
func RegisterUser(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data pendaftaran tidak valid",
			"error":   err.Error(),
		})
		return
	}

	// Verify if email is already taken
	var existing models.User
	if err := config.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Email sudah terdaftar",
		})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal memproses kata sandi",
		})
		return
	}

	// Create user record (role is hardcoded as 'customer' for safety)
	newUser := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     "customer",
		Phone:    req.Phone,
		Address:  req.Address,
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&newUser).Error; err != nil {
			return err
		}
		logEntry := models.ActivityLog{
			UserID:      newUser.ID,
			UserRole:    newUser.Role,
			ActionType:  "REGISTER_CUSTOMER",
			Description: "New User Registered: " + newUser.Email,
			CreatedAt:   time.Now(),
		}
		return tx.Create(&logEntry).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mendaftarkan pengguna baru",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Pendaftaran berhasil! Silakan masuk.",
		"data": gin.H{
			"id":       newUser.ID,
			"name":     newUser.Name,
			"email":    newUser.Email,
			"role":     newUser.Role,
		},
	})
}

// RegisterSellerRequest represents the JSON payload for seller registration.
type RegisterSellerRequest struct {
	Name      string `json:"name" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
	StoreName string `json:"store_name" binding:"required"`
}

// RegisterSeller handles POST /api/register-seller
func RegisterSeller(c *gin.Context) {
	var req RegisterSellerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Data pendaftaran tidak valid", "error": err.Error()})
		return
	}

	var existing models.User
	if err := config.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Email sudah terdaftar"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal memproses kata sandi"})
		return
	}

	newUser := models.User{
		Name:      req.Name,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Role:      "seller",
		StoreName: req.StoreName,
	}

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&newUser).Error; err != nil {
			return err
		}
		logEntry := models.ActivityLog{
			UserID:      newUser.ID,
			UserRole:    newUser.Role,
			ActionType:  "REGISTER_SELLER",
			Description: "New Seller Registered: " + newUser.Email,
			CreatedAt:   time.Now(),
		}
		return tx.Create(&logEntry).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mendaftarkan toko", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Pendaftaran Toko berhasil! Silakan masuk.",
		"data": gin.H{"id": newUser.ID, "email": newUser.Email, "role": newUser.Role},
	})
}

// GetCurrentUser handles GET /api/user/:id
// Returns profile data for a specific user (used for session validation).
func GetCurrentUser(c *gin.Context) {
	userID := c.Param("id")

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
			"role":     user.Role,
			"email":    user.Email,
			"avatar":   user.Avatar,
		},
	})
}


