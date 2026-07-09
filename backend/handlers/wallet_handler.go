package handlers

import (
	"net/http"
	"openpeo-backend/config"
	"openpeo-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetWalletDetails fetches wallet balance and transactions
func GetWalletDetails(c *gin.Context) {
	userIDFloat, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}
	userID := int32(userIDFloat.(float64))

	var profile models.SellerProfile
	if err := config.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Profile tidak ditemukan"})
		return
	}

	var transactions []models.WalletTransaction
	config.DB.Where("seller_id = ?", profile.ID).Order("created_at desc").Find(&transactions)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"active_balance":      profile.ActiveBalance,
			"pending_balance":     profile.PendingBalance,
			"bank_name":           profile.BankName,
			"bank_account_number": profile.BankAccountNumber,
			"bank_account_name":   profile.BankAccountName,
			"transactions":        transactions,
		},
	})
}

// WithdrawBalance handles the withdrawal request
func WithdrawBalance(c *gin.Context) {
	userIDFloat, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}
	userID := int32(userIDFloat.(float64))

	var input struct {
		Amount float64 `json:"amount" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Nominal tidak valid"})
		return
	}

	var profile models.SellerProfile
	if err := config.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Profile tidak ditemukan"})
		return
	}

	// VALIDATION: Bank Setup
	if profile.BankAccountNumber == "" || profile.BankName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Silakan atur rekening bank di Pengaturan Toko terlebih dahulu."})
		return
	}

	// VALIDATION: Minimum Withdrawal
	if input.Amount < 50000 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Minimal penarikan adalah Rp 50.000."})
		return
	}

	// VALIDATION: Sufficient Balance
	if input.Amount > profile.ActiveBalance {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Saldo tidak mencukupi."})
		return
	}

	// DB TRANSACTION
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		// Deduct Balance
		profile.ActiveBalance -= input.Amount
		if err := tx.Save(&profile).Error; err != nil {
			return err
		}

		// Insert Transaction Log
		transaction := models.WalletTransaction{
			SellerID: profile.ID,
			Type:     "withdrawal",
			Amount:   input.Amount,
			Status:   "processing",
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		// TODO: Call Xendit/Midtrans Disbursement API here

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal memproses penarikan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Penarikan berhasil diajukan dan sedang diproses.",
	})
}
