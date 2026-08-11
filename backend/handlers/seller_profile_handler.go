package handlers

import (
	"net/http"
	"openpeo-backend/config"
	"openpeo-backend/models"

	"github.com/gin-gonic/gin"
)

func GetSellerProfile(c *gin.Context) {
	userIDFloat, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := int32(userIDFloat.(float64))

	var profile models.SellerProfile
	err := config.DB.Where("user_id = ?", userID).First(&profile).Error
	if err != nil {

		c.JSON(http.StatusOK, models.SellerProfile{
			UserID: userID,
		})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func UpsertSellerProfile(c *gin.Context) {
	userIDFloat, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := int32(userIDFloat.(float64))

	var input models.SellerProfile
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var profile models.SellerProfile
	err := config.DB.Where("user_id = ?", userID).First(&profile).Error

	if err != nil {

		input.UserID = userID
		if err := config.DB.Create(&input).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan profil"})
			return
		}
		c.JSON(http.StatusOK, input)
	} else {

		profile.StoreName = input.StoreName
		profile.Description = input.Description
		profile.StoreLogo = input.StoreLogo
		profile.Phone = input.Phone
		profile.Address = input.Address
		profile.Region = input.Region
		profile.BankName = input.BankName
		profile.BankAccountNumber = input.BankAccountNumber
		profile.BankAccountName = input.BankAccountName

		if err := config.DB.Save(&profile).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui profil"})
			return
		}
		c.JSON(http.StatusOK, profile)
	}
}
