package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"openpeo-backend/config"
	"openpeo-backend/models"
)

// SalesData represents the aggregate response for the admin sales dashboard.
type SalesData struct {
	TotalRevenue    float64          `json:"total_revenue"`
	TotalPreOrders  int64            `json:"total_pre_orders"`
	TotalProducts   int64            `json:"total_products"`
	TotalCustomers  int64            `json:"total_customers"`
	PendingOrders   int64            `json:"pending_orders"`
	CompletedOrders int64            `json:"completed_orders"`
	RevenueByRegion []RegionRevenue  `json:"revenue_by_region"`
	RecentOrders    []models.Order   `json:"recent_orders"`
	OrdersByStatus  []StatusCount    `json:"orders_by_status"`
	DailySales      []DailySalePoint `json:"daily_sales"`
}

// RegionRevenue holds revenue data grouped by NTT region.
type RegionRevenue struct {
	Region  string  `json:"region"`
	Revenue float64 `json:"revenue"`
	Count   int64   `json:"count"`
}

// StatusCount holds order counts grouped by status.
type StatusCount struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

// DailySalePoint holds daily sales data for charting.
type DailySalePoint struct {
	Date    string  `json:"date"`
	Revenue float64 `json:"revenue"`
	Orders  int64   `json:"orders"`
}

// GetSalesData handles GET /api/admin/sales-data
// Protected endpoint for admins — runs aggregate queries on the orders table
// to fetch total revenue, pre-order counts, and dynamic sales statistics.
func GetSalesData(c *gin.Context) {
	var data SalesData

	// ── Total Revenue (sum of all order total_price) ──
	config.DB.Model(&models.Order{}).
		Select("COALESCE(SUM(total_price), 0)").
		Scan(&data.TotalRevenue)

	// ── Total Pre-Orders count ──
	config.DB.Model(&models.Order{}).Count(&data.TotalPreOrders)

	// ── Total Active Products ──
	config.DB.Model(&models.Product{}).
		Where("is_active = ?", true).
		Count(&data.TotalProducts)

	// ── Total Customers ──
	config.DB.Model(&models.User{}).
		Where("role = ?", "customer").
		Count(&data.TotalCustomers)

	// ── Pending Orders ──
	config.DB.Model(&models.Order{}).
		Where("status = ?", "pending").
		Count(&data.PendingOrders)

	// ── Completed Orders ──
	config.DB.Model(&models.Order{}).
		Where("status IN ('completed', 'shipped')").
		Count(&data.CompletedOrders)

	// ── Revenue by NTT Region ──
	config.DB.Model(&models.Order{}).
		Select("products.region as region, COALESCE(SUM(orders.total_price), 0) as revenue, COUNT(orders.id) as count").
		Joins("JOIN products ON orders.product_id = products.id").
		Group("products.region").
		Order("revenue DESC").
		Scan(&data.RevenueByRegion)

	// ── Orders by Status ──
	config.DB.Model(&models.Order{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&data.OrdersByStatus)

	// ── Daily Sales (last 7 days) ──
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	config.DB.Model(&models.Order{}).
		Select("DATE(created_at) as date, COALESCE(SUM(total_price), 0) as revenue, COUNT(*) as orders").
		Where("created_at >= ?", sevenDaysAgo).
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&data.DailySales)

	// ── Recent Orders (latest 10) ──
	config.DB.
		Preload("Customer").
		Preload("Product").
		Order("created_at DESC").
		Limit(10).
		Find(&data.RecentOrders)

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"data":       data,
		"fetched_at": time.Now().Format(time.RFC3339),
	})
}
