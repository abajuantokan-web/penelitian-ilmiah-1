package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"openpeo-backend/config"
	"openpeo-backend/models"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

func CreateOrder(c *gin.Context) {
	userIDFloat, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}
	customerID := int32(userIDFloat.(float64))

	var req struct {
		Items	[]struct {
			ProductID	int32	`json:"product_id"`
			Quantity	int	`json:"quantity"`
			Price		float64	`json:"price"`
		}	`json:"items"`
		Note	string	`json:"note"`
	}
	c.ShouldBindJSON(&req)

	if len(req.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Payload items kosong"})
		return
	}

	itemsBySeller := make(map[int32][]models.OrderItem)
	var grandTotal float64

	for _, reqItem := range req.Items {
		var product models.Product
		if err := config.DB.First(&product, reqItem.ProductID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": fmt.Sprintf("Produk ID %d tidak ditemukan", reqItem.ProductID)})
			return
		}

		itemTotal := reqItem.Price * float64(reqItem.Quantity)
		grandTotal += itemTotal

		itemsBySeller[product.SellerID] = append(itemsBySeller[product.SellerID], models.OrderItem{
			ProductID:	reqItem.ProductID,
			Quantity:	reqItem.Quantity,
			Price:		reqItem.Price,
		})
	}

	paymentRef := fmt.Sprintf("INV-%d-%d", customerID, time.Now().Unix())

	var createdOrders []models.Order

	err := config.DB.Transaction(func(tx *gorm.DB) error {

		for sellerID, items := range itemsBySeller {
			var orderTotal float64
			for _, item := range items {
				orderTotal += item.Price * float64(item.Quantity)
			}

			order := models.Order{
				CustomerID:		customerID,
				SellerID:		sellerID,
				Quantity:		len(items),
				TotalPrice:		orderTotal,
				Status:			"Menunggu Pembayaran",
				Note:			req.Note,
				PaymentReference:	paymentRef,
			}

			if err := tx.Create(&order).Error; err != nil {
				return err
			}

			for i := range items {
				items[i].OrderID = order.ID
				if err := tx.Create(&items[i]).Error; err != nil {
					return err
				}
			}

			logEntry := models.ActivityLog{
				UserID:		customerID,
				UserRole:	"customer",
				ActionType:	"CHECKOUT",
				Description:	fmt.Sprintf("New Order Created: #%d (Seller %d)", order.ID, sellerID),
				CreatedAt:	time.Now(),
			}
			if err := tx.Create(&logEntry).Error; err != nil {
				return err
			}

			var sellerProfile models.SellerProfile
			if err := tx.Where("user_id = ?", sellerID).First(&sellerProfile).Error; err == nil {
				sellerProfile.PendingBalance += order.TotalPrice
				if err := tx.Save(&sellerProfile).Error; err != nil {
					return err
				}
			}

			createdOrders = append(createdOrders, order)
		}

		if err := tx.Where("user_id = ?", customerID).Delete(&models.CartItem{}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal membuat pesanan", "error": err.Error()})
		return
	}

	for _, order := range createdOrders {
		ChatHub.SendToUser(order.SellerID, gin.H{
			"type":	"NEW_ORDER_CREATED",
			"data":	order,
		})
	}

	midtrans.Environment = midtrans.Sandbox
	midtrans.ServerKey = "Mid-server-ZM2VH9KSHXRCMBsC_cdP_Xk7"

	var user models.User
	config.DB.First(&user, customerID)

	reqSnap := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:	paymentRef,
			GrossAmt:	int64(grandTotal),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName:	user.Name,
			Email:	user.Email,
			Phone:	user.Phone,
		},
		CustomField1:	req.Note,
	}

	var snapToken string
	snapResp, midErr := snap.CreateTransaction(reqSnap)
	if midErr == nil {
		snapToken = snapResp.Token
	} else {
		fmt.Printf("Midtrans Error: %v\n", midErr)
	}

	c.JSON(http.StatusCreated, gin.H{
		"success":		true,
		"message":		"Pesanan multi-vendor berhasil dibuat",
		"data":			createdOrders[0],
		"payment_reference":	paymentRef,
		"snap_token":		snapToken,
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
		Preload("OrderItems").
		Preload("OrderItems.Product").
		Order("created_at DESC").
		Find(&orders)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":	false,
			"message":	"Failed to fetch orders",
			"error":	result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":	true,
		"data":		orders,
		"total":	total,
	})
}

func DeleteOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":	false,
			"message":	"Invalid order ID",
		})
		return
	}

	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":	false,
			"message":	"Parameter user_id (Admin ID) diperlukan untuk menghapus log transaksi",
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
			"message":	"Hanya Admin yang memiliki hak akses untuk menghapus log transaksi",
		})
		return
	}

	var order models.Order
	if err := config.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success":	false,
			"message":	"Order tidak ditemukan",
		})
		return
	}

	if err := config.DB.Unscoped().Delete(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":	false,
			"message":	"Gagal menghapus order",
			"error":	err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":	true,
		"message":	"Transaksi berhasil dihapus",
	})
}

func CreateDirectOrder(c *gin.Context) {
	userIDFloat, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}
	customerID := int32(userIDFloat.(float64))

	var req struct {
		ProductID	int32	`json:"product_id" binding:"required"`
		Quantity	int	`json:"quantity" binding:"required,min=1"`
		Price		float64	`json:"price" binding:"required"`
		Note		string	`json:"note"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	var product models.Product
	if err := config.DB.First(&product, req.ProductID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Produk tidak ditemukan"})
		return
	}

	order := models.Order{
		CustomerID:	customerID,
		SellerID:	product.SellerID,
		Quantity:	req.Quantity,
		TotalPrice:	req.Price * float64(req.Quantity),
		Status:		"Menunggu Konfirmasi",
		Note:		req.Note,
	}

	if err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		orderItem := models.OrderItem{
			OrderID:	order.ID,
			ProductID:	req.ProductID,
			Quantity:	req.Quantity,
			Price:		req.Price,
		}
		if err := tx.Create(&orderItem).Error; err != nil {
			return err
		}

		var sellerProfile models.SellerProfile
		if err := tx.Where("user_id = ?", order.SellerID).First(&sellerProfile).Error; err != nil {
			return err
		}
		sellerProfile.PendingBalance += order.TotalPrice
		return tx.Save(&sellerProfile).Error
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal menyimpan order"})
		return
	}

	fmt.Printf("Order created. Assigned to SellerID: %d\n", order.SellerID)

	midtrans.Environment = midtrans.Sandbox
	midtrans.ServerKey = "Mid-server-ZM2VH9KSHXRCMBsC_cdP_Xk7"
	orderIDStr := fmt.Sprintf("ORDER-%d-%d", order.ID, time.Now().Unix())

	var user models.User
	config.DB.First(&user, customerID)

	reqSnap := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:	orderIDStr,
			GrossAmt:	int64(order.TotalPrice),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName:	user.Name,
			Email:	user.Email,
			Phone:	user.Phone,
		},
		CustomField1:	req.Note,
	}

	var snapToken string
	snapResp, midErr := snap.CreateTransaction(reqSnap)
	if midErr == nil {
		snapToken = snapResp.Token
	} else {
		fmt.Printf("Midtrans Error: %v\n", midErr)
	}

	c.JSON(http.StatusCreated, gin.H{
		"success":	true,
		"message":	"Pesanan langsung berhasil dibuat",
		"data":		order,
		"snap_token":	snapToken,
	})
}

func GetSellerOrders(c *gin.Context) {
	sellerID, ok := requireSeller(c)
	if !ok {
		return
	}

	var orders []models.Order

	err := config.DB.
		Where("seller_id = ?", sellerID).
		Preload("Customer").
		Preload("OrderItems").
		Preload("OrderItems.Product").
		Order("created_at DESC").
		Find(&orders).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengambil order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": orders})
}

func ProcessSellerOrder(c *gin.Context) {
	sellerID, ok := requireSeller(c)
	if !ok {
		return
	}
	orderID := c.Param("id")

	var order models.Order

	err := config.DB.
		Where("id = ? AND seller_id = ?", orderID, sellerID).
		First(&order).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Order tidak ditemukan"})
		return
	}

	if order.Status != "Menunggu Konfirmasi" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Order tidak dalam status Menunggu Konfirmasi"})
		return
	}

	order.Status = "Pesanan Sedang Diproses"
	if err := config.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal update status"})
		return
	}

	ChatHub.SendToUser(sellerID, gin.H{
		"type":	"ORDER_STATUS_UPDATED",
		"data":	order,
	})

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Status berhasil diubah ke Pesanan Sedang Diproses", "data": order})
}

func GetSellerDashboardStats(c *gin.Context) {
	sellerID, ok := requireSeller(c)
	if !ok {
		return
	}

	var pendingCount int64
	var processingCount int64
	var completedCount int64
	var totalRevenue float64
	var recentOrders []models.Order

	config.DB.Model(&models.Order{}).Where("seller_id = ? AND status = ?", sellerID, "Menunggu Konfirmasi").Count(&pendingCount)

	config.DB.Model(&models.Order{}).Where("seller_id = ? AND status = ?", sellerID, "Pesanan Sedang Diproses").Count(&processingCount)

	config.DB.Model(&models.Order{}).Where("seller_id = ? AND (status = ? OR status = ?)", sellerID, "Pesanan Selesai", "Selesai").Count(&completedCount)

	config.DB.Model(&models.Order{}).Where("seller_id = ? AND (status = ? OR status = ?)", sellerID, "Pesanan Selesai", "Selesai").Select("COALESCE(SUM(total_price), 0)").Scan(&totalRevenue)

	config.DB.Preload("Customer").Where("seller_id = ?", sellerID).Order("created_at desc").Limit(5).Find(&recentOrders)

	c.JSON(http.StatusOK, gin.H{
		"success":	true,
		"data": gin.H{
			"total_orders_pending":		pendingCount,
			"total_orders_processing":	processingCount,
			"total_completed_orders":	completedCount,
			"total_revenue":		totalRevenue,
			"recent_orders":		recentOrders,
		},
	})
}

func GetSellerDashboardChart(c *gin.Context) {
	sellerID, ok := requireSeller(c)
	if !ok {
		return
	}

	rangeParam := c.Query("range")
	days := 7
	switch rangeParam {
	case "30_days":
		days = 30
	case "12_months":
		days = 365
	}

	type ChartData struct {
		Date	string	`json:"date"`
		Total	float64	`json:"total"`
	}

	var results []ChartData

	query := `
		SELECT DATE(created_at) as date, SUM(total_price) as total
		FROM orders
		WHERE seller_id = ? AND (status = 'Pesanan Selesai' OR status = 'Selesai')
		AND created_at >= DATE_SUB(CURDATE(), INTERVAL ? DAY)
		GROUP BY DATE(created_at)
		ORDER BY DATE(created_at) ASC
	`

	if err := config.DB.Raw(query, sellerID, days).Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengambil data chart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":	true,
		"data":		results,
	})
}

func ConfirmPayment(c *gin.Context) {
	ref := c.Param("id")

	var orders []models.Order

	if err := config.DB.Where("payment_reference = ? OR id = ?", ref, ref).Find(&orders).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Pesanan tidak ditemukan"})
		return
	}

	if len(orders) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Pesanan tidak ditemukan"})
		return
	}

	for _, order := range orders {
		if order.Status == "Menunggu Pembayaran" {
			order.Status = "Menunggu Konfirmasi"
			if err := config.DB.Save(&order).Error; err != nil {
				continue
			}

			ChatHub.SendToUser(order.SellerID, gin.H{
				"type":	"NEW_ORDER_CREATED",
				"data":	order,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success":	true,
		"message":	"Pembayaran berhasil dikonfirmasi",
	})
}

func CompleteOrder(c *gin.Context) {
	userIDFloat, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}
	customerID := int32(userIDFloat.(float64))
	orderID := c.Param("id")

	var order models.Order
	if err := config.DB.Where("id = ? AND customer_id = ?", orderID, customerID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Order tidak ditemukan"})
		return
	}

	if order.Status == "Selesai" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Order sudah selesai"})
		return
	}

	err := config.DB.Transaction(func(tx *gorm.DB) error {

		order.Status = "Selesai"
		if err := tx.Save(&order).Error; err != nil {
			return err
		}

		var sellerProfile models.SellerProfile
		if err := tx.Where("user_id = ?", order.SellerID).First(&sellerProfile).Error; err != nil {
			return err
		}

		sellerProfile.PendingBalance -= order.TotalPrice
		if sellerProfile.PendingBalance < 0 {
			sellerProfile.PendingBalance = 0
		}
		sellerProfile.ActiveBalance += order.TotalPrice
		if err := tx.Save(&sellerProfile).Error; err != nil {
			return err
		}

		walletTx := models.WalletTransaction{
			SellerID:	sellerProfile.ID,
			Type:		"income",
			Amount:		order.TotalPrice,
			Status:		"completed",
		}
		if err := tx.Create(&walletTx).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal menyelesaikan order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Order berhasil diselesaikan, dana diteruskan ke penjual"})
}
