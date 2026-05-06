package app

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	gin "github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type updateOrderStatusRequest struct {
	Remark string `json:"remark"`
}

type createProductRequest struct {
	Name         string                    `json:"name" binding:"required"`
	Subtitle     string                    `json:"subtitle"`
	CoverImage   string                    `json:"coverImage"`
	CategoryName string                    `json:"categoryName"`
	Description  string                    `json:"description"`
	SortOrder    int                       `json:"sortOrder"`
	SKUs         []createProductSKURequest `json:"skus" binding:"required"`
}

type createProductSKURequest struct {
	SKUName       string  `json:"skuName" binding:"required"`
	SKUCode       string  `json:"skuCode" binding:"required"`
	Price         float64 `json:"price" binding:"required"`
	OriginalPrice float64 `json:"originalPrice"`
	StockTotal    int     `json:"stockTotal"`
	LimitPerUser  int     `json:"limitPerUser"`
	LimitPerOrder int     `json:"limitPerOrder"`
}

type createGroupRequest struct {
	Title                   string                   `json:"title" binding:"required"`
	CoverImage              string                   `json:"coverImage"`
	LeaderUserID            uint                     `json:"leaderUserId"`
	StartAt                 time.Time                `json:"startAt" binding:"required"`
	CutoffAt                time.Time                `json:"cutoffAt" binding:"required"`
	FulfillmentMode         string                   `json:"fulfillmentMode" binding:"required"`
	AllowModifyBeforeCutoff bool                     `json:"allowModifyBeforeCutoff"`
	ShowJoinList            bool                     `json:"showJoinList"`
	PickupRuleDesc          string                   `json:"pickupRuleDesc"`
	DeliveryRuleDesc        string                   `json:"deliveryRuleDesc"`
	GroupNotice             string                   `json:"groupNotice"`
	MenuDate                string                   `json:"menuDate"`
	Items                   []createGroupItemRequest `json:"items"`
}

type createDailyMenuRequest struct {
	MenuDate string                       `json:"menuDate" binding:"required"`
	Title    string                       `json:"title"`
	Items    []createDailyMenuItemRequest `json:"items" binding:"required"`
}

type createDailyMenuItemRequest struct {
	ProductSKUID  uint    `json:"productSkuId" binding:"required"`
	StockTotal    int     `json:"stockTotal" binding:"required"`
	Price         float64 `json:"price" binding:"required"`
	OriginalPrice float64 `json:"originalPrice"`
	LimitPerUser  int     `json:"limitPerUser"`
	LimitPerOrder int     `json:"limitPerOrder"`
	SortOrder     int     `json:"sortOrder"`
}

type createGroupItemRequest struct {
	ProductSKUID  uint    `json:"productSkuId" binding:"required"`
	StockTotal    int     `json:"stockTotal" binding:"required"`
	Price         float64 `json:"price" binding:"required"`
	OriginalPrice float64 `json:"originalPrice"`
	LimitPerUser  int     `json:"limitPerUser"`
	LimitPerOrder int     `json:"limitPerOrder"`
	SortOrder     int     `json:"sortOrder"`
}

func (s Server) markOrderReadyForPickup(c *gin.Context) {
	s.updateOrderStatus(c, []string{"cutoff_locked"}, "ready_for_pickup")
}

func (s Server) markOrderDelivering(c *gin.Context) {
	s.updateOrderStatus(c, []string{"cutoff_locked"}, "delivering")
}

func (s Server) completeOrder(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, 40001, "invalid order id", nil)
		return
	}

	var req updateOrderStatusRequest
	_ = c.ShouldBindJSON(&req)

	err = s.db.Transaction(func(tx *gorm.DB) error {
		var order Order
		if err := tx.First(&order, id).Error; err != nil {
			return err
		}
		if order.Status != "ready_for_pickup" && order.Status != "delivering" {
			return fmt.Errorf("ORDER_STATUS_INVALID")
		}
		now := time.Now()
		order.Status = "completed"
		order.CompletedAt = &now
		if req.Remark != "" {
			order.Remark = req.Remark
		}
		if err := tx.Save(&order).Error; err != nil {
			return err
		}

		var user User
		if err := tx.First(&user, order.UserID).Error; err == nil {
			gain := int(order.PaidAmount)
			user.PointsBalance += gain
			if err := tx.Save(&user).Error; err != nil {
				return err
			}
			sourceID := order.ID
			ledger := PointsLedger{UserID: user.ID, ChangeValue: gain, BalanceAfter: user.PointsBalance, SourceType: "order_reward", SourceID: &sourceID, Remark: "订单完成赠送积分"}
			if err := tx.Create(&ledger).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		if err.Error() == "ORDER_STATUS_INVALID" {
			writeError(c, http.StatusBadRequest, 40015, "订单状态不可完成", nil)
			return
		}
		writeError(c, http.StatusInternalServerError, 50000, "complete order failed", gin.H{"detail": err.Error()})
		return
	}

	writeSuccess(c, gin.H{"orderId": id, "status": "completed"})
}

func (s Server) updateOrderStatus(c *gin.Context, allowed []string, target string) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, 40001, "invalid order id", nil)
		return
	}

	var req updateOrderStatusRequest
	_ = c.ShouldBindJSON(&req)

	err = s.db.Transaction(func(tx *gorm.DB) error {
		var order Order
		if err := tx.First(&order, id).Error; err != nil {
			return err
		}
		allowedMap := map[string]bool{}
		for _, item := range allowed {
			allowedMap[item] = true
		}
		if !allowedMap[order.Status] {
			return fmt.Errorf("ORDER_STATUS_INVALID")
		}
		order.Status = target
		if req.Remark != "" {
			order.Remark = req.Remark
		}
		return tx.Save(&order).Error
	})

	if err != nil {
		if err.Error() == "ORDER_STATUS_INVALID" {
			writeError(c, http.StatusBadRequest, 40015, "订单状态不可流转", nil)
			return
		}
		writeError(c, http.StatusInternalServerError, 50000, "update order status failed", gin.H{"detail": err.Error()})
		return
	}

	writeSuccess(c, gin.H{"orderId": id, "status": target})
}
