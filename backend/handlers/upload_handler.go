package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// UploadImage handles POST /api/upload
// Accepts multipart/form-data with a file field named "image"
func UploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "No file uploaded or invalid form field",
			"error":   err.Error(),
		})
		return
	}

	// Ensure the ./images directory exists
	uploadDir := "./images"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err = os.MkdirAll(uploadDir, os.ModePerm)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to create upload directory",
			})
			return
		}
	}

	// Generate a unique filename using timestamp
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("upload_%d%s", time.Now().UnixNano(), ext)
	savePath := filepath.Join(uploadDir, filename)

	// Save the file
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to save file",
			"error":   err.Error(),
		})
		return
	}

	// Return the accessible URL
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Image uploaded successfully",
		"url":     fmt.Sprintf("images/%s", filename),
	})
}
