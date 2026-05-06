package app

import (
	"net/http"
	"os"
	"strings"
	"time"

	gin "github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	db  *gorm.DB
	cfg Config
}

func RegisterRoutes(router *gin.Engine, db *gorm.DB, cfg Config) {
	_ = ConfigurePool(db)

	server := Server{db: db, cfg: cfg}
	_ = os.MkdirAll(cfg.UploadDir, 0755)

	router.Use(corsMiddleware(cfg.CORSAllowOrigin))
	router.Static("/uploads", cfg.UploadDir)

	router.GET("/api/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	api := router.Group("/api/v1")
	{
		api.GET("/health", server.health)
		api.GET("/groups", server.listGroups)
		api.GET("/groups/:id", server.getGroup)
		api.GET("/products", server.listProducts)
		api.GET("/pickup-points", server.listPickupPoints)
		api.POST("/orders", server.createOrder)
		api.GET("/orders/:id", server.getOrder)
		api.GET("/orders/no/:orderNo", server.getOrderByNo)
		api.POST("/orders/:id/cancel", server.cancelOrder)
		api.GET("/me", server.getCurrentUser)
		api.GET("/me/orders", server.listMyOrders)
		api.GET("/me/addresses", server.listMyAddresses)
		api.POST("/me/addresses", server.createAddress)
		api.PUT("/me/addresses/:id", server.updateAddress)
		api.DELETE("/me/addresses/:id", server.deleteAddress)
		api.POST("/me/addresses/:id/default", server.setDefaultAddress)
		api.GET("/me/coupons", server.listMyCoupons)
		api.GET("/me/points", server.listMyPoints)
		api.GET("/me/notifications", server.listMyNotifications)
		api.GET("/leader/groups/:id/summary", server.getGroupSummary)
		api.POST("/ai/recommend", server.aiRecommend)
		api.POST("/leader/groups/:id/cutoff", server.cutoffGroup)
		api.POST("/admin/orders/:id/ready-for-pickup", server.markOrderReadyForPickup)
		api.POST("/admin/orders/:id/start-delivery", server.markOrderDelivering)
		api.POST("/admin/orders/:id/complete", server.completeOrder)
		api.GET("/admin/orders", server.listAdminOrders)
		api.GET("/admin/coupons", server.listAdminCoupons)
		api.GET("/admin/user-coupons", server.listAdminUserCoupons)
		api.POST("/admin/coupons", server.createCoupon)
		api.POST("/admin/coupons/:id/disable", server.disableCoupon)
		api.POST("/admin/coupons/grant", server.grantCoupon)
		api.GET("/admin/products", server.listAdminProducts)
		api.GET("/admin/daily-menus", server.listAdminDailyMenus)
		api.POST("/admin/daily-menus", server.createDailyMenu)
		api.DELETE("/admin/daily-menus/:id", server.deleteDailyMenu)
		api.POST("/admin/products", server.createProduct)
		api.POST("/admin/uploads/image", server.uploadImage)
		api.GET("/admin/groups", server.listAdminGroups)
		api.POST("/admin/groups", server.createGroup)
	}
}

func corsMiddleware(allowOrigin string) gin.HandlerFunc {
	allowed := strings.Split(allowOrigin, ",")
	for i := range allowed {
		allowed[i] = strings.TrimSpace(allowed[i])
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		for _, item := range allowed {
			if item != "" && origin == item {
				c.Header("Access-Control-Allow-Origin", origin)
				break
			}
		}
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Idempotency-Key")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Max-Age", "86400")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func writeSuccess(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"data":    data,
	})
}

func writeError(c *gin.Context, status int, code int, message string, data any) {
	c.JSON(status, gin.H{
		"code":    code,
		"message": message,
		"data":    data,
	})
}

type healthResponse struct {
	ServerTime time.Time `json:"serverTime"`
}

func (s Server) health(c *gin.Context) {
	writeSuccess(c, healthResponse{ServerTime: time.Now()})
}
