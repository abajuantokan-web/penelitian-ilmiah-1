package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"openpeo-backend/config"
	"openpeo-backend/models"
)

// LoginRequest is the expected JSON payload for the login endpoint.
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest represents the JSON payload for registration.
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
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
			"message": "Username dan password harus diisi",
			"error":   err.Error(),
		})
		return
	}

	// Look up user by username
	var user models.User
	result := config.DB.Where("username = ?", req.Username).First(&user)

	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Username tidak ditemukan",
		})
		return
	}

	// Validate password (plain text comparison for demo — use bcrypt in production)
	if user.Password != req.Password {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Password salah",
		})
		return
	}

	// Return user session data
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login berhasil!",
		"data": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"name":     user.Name,
			"role":     user.Role,
			"email":    user.Email,
			"avatar":   user.Avatar,
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

	// Verify if username is already taken
	var existing models.User
	if err := config.DB.Where("username = ?", req.Username).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Username sudah terdaftar",
		})
		return
	}

	// Create user record (role is hardcoded as 'customer' for safety)
	newUser := models.User{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Role:     "customer",
		Phone:    req.Phone,
		Address:  req.Address,
	}

	if err := config.DB.Create(&newUser).Error; err != nil {
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
			"username": newUser.Username,
			"name":     newUser.Name,
			"role":     newUser.Role,
		},
	})
}

// GetCurrentUser handles GET /api/user/:id
// Returns profile data for a specific user (used for session validation).
func GetCurrentUser(c *gin.Context) {
	userID := c.Param("id")

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
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
			"username": user.Username,
			"name":     user.Name,
			"role":     user.Role,
			"email":    user.Email,
			"avatar":   user.Avatar,
		},
	})
}

// GetChatContacts handles GET /api/chat/contacts
// For admins: returns a list of unique customers who have sent messages.
// For customers: returns a list of admins/vendors they can chat with.
func GetChatContacts(c *gin.Context) {
	userID := c.Query("user_id")
	role := c.Query("role")

	if userID == "" || role == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "user_id and role are required",
		})
		return
	}

	var contacts []models.User

	if role == "admin" {
		// For admins: show all customers to allow messaging any registered customer
		config.DB.Where("role = ?", "customer").Order("name ASC").Find(&contacts)
	} else {
		// For customers: show all admins and vendors they can contact
		config.DB.Where("role IN ('admin', 'vendor')").Find(&contacts)
	}

	// Calculate unread count for each contact
	type ContactInfo struct {
		ID          int32  `json:"id"`
		Username    string `json:"username"`
		Name        string `json:"name"`
		Role        string `json:"role"`
		Avatar      string `json:"avatar"`
		UnreadCount int64  `json:"unread_count"`
	}

	contactInfos := make([]ContactInfo, 0, len(contacts))
	for _, contact := range contacts {
		var unreadCount int64
		config.DB.Model(&models.Message{}).
			Where("sender_id = ? AND receiver_id = ? AND is_read = ?", contact.ID, userID, false).
			Count(&unreadCount)

		contactInfos = append(contactInfos, ContactInfo{
			ID:          contact.ID,
			Username:    contact.Username,
			Name:        contact.Name,
			Role:        contact.Role,
			Avatar:      contact.Avatar,
			UnreadCount: unreadCount,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    contactInfos,
	})
}
