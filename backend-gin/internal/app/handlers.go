package app

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	gin "github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type createOrderRequest struct {
	GroupID         uint                     `json:"groupId" binding:"required"`
	FulfillmentMode string                   `json:"fulfillmentMode" binding:"required"`
	PickupPointID   *uint                    `json:"pickupPointId"`
	AddressID       *uint                    `json:"addressId"`
	UserCouponID    *uint                    `json:"userCouponId"`
	ContactName     string                   `json:"contactName" binding:"required"`
	ContactPhone    string                   `json:"contactPhone" binding:"required"`
	Remark          string                   `json:"remark"`
	Items           []createOrderItemRequest `json:"items" binding:"required"`
}

type createOrderItemRequest struct {
	GroupItemID uint   `json:"groupItemId" binding:"required"`
	Quantity    int    `json:"quantity" binding:"required"`
	TasteRemark string `json:"tasteRemark"`
}

type cancelOrderRequest struct {
	Reason string `json:"reason"`
}

type cutoffGroupRequest struct {
	Reason string `json:"reason"`
}

type summaryRow struct {
	ProductName string `json:"productName"`
	SKUName     string `json:"skuName"`
	TotalQty    int    `json:"totalQty"`
}

func (s Server) listGroups(c *gin.Context) {
	var groups []Group
	query := s.db.Preload("Items", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order asc")
	}).Order("cutoff_at asc")

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if minCutoffAt := c.Query("minCutoffAt"); minCutoffAt != "" {
		parsed, err := time.Parse(time.RFC3339, minCutoffAt)
		if err != nil {
			parsed, err = time.Parse("2006-01-02", minCutoffAt)
		}
		if err != nil {
			writeError(c, http.StatusBadRequest, 40001, "invalid minCutoffAt, expect RFC3339 or YYYY-MM-DD", nil)
			return
		}
		query = query.Where("cutoff_at >= ?", parsed)
	}

	if err := query.Find(&groups).Error; err != nil {
		writeError(c, http.StatusInternalServerError, 50000, "query groups failed", nil)
		return
	}

	writeSuccess(c, gin.H{"list": groups})
}

func (s Server) getGroup(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, 40001, "invalid group id", nil)
		return
	}

	var group Group
	if err := s.db.Preload("Items", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order asc")
	}).First(&group, id).Error; err != nil {
		writeError(c, http.StatusNotFound, 40004, "group not found", nil)
		return
	}

	writeSuccess(c, group)
}

func (s Server) listProducts(c *gin.Context) {
	var products []Product
	query := s.db.Preload("SKUs").Order("sort_order asc, id desc")
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Find(&products).Error; err != nil {
		writeError(c, http.StatusInternalServerError, 50000, "query products failed", nil)
		return
	}
	writeSuccess(c, gin.H{"list": products})
}

func (s Server) listPickupPoints(c *gin.Context) {
	var points []PickupPoint
	if err := s.db.Order("id asc").Find(&points).Error; err != nil {
		writeError(c, http.StatusInternalServerError, 50000, "query pickup points failed", nil)
		return
	}
	writeSuccess(c, gin.H{"list": points})
}

func (s Server) getOrder(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, 40001, "invalid order id", nil)
		return
	}

	var order Order
	if err := s.db.Preload("Items").First(&order, id).Error; err != nil {
		writeError(c, http.StatusNotFound, 40004, "order not found", nil)
		return
	}

	writeSuccess(c, order)
}

func (s Server) listMyOrders(c *gin.Context) {
	var orders []Order
	query := s.db.Preload("Items").Where("user_id = ?", 1).Order("created_at desc")
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Find(&orders).Error; err != nil {
		writeError(c, http.StatusInternalServerError, 50000, "query my orders failed", nil)
		return
	}
	writeSuccess(c, gin.H{"list": orders})
}

func (s Server) createOrder(c *gin.Context) {
	var req createOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, 40001, "invalid request", gin.H{"detail": err.Error()})
		return
	}

	if len(req.Items) == 0 {
		writeError(c, http.StatusBadRequest, 40001, "items is required", nil)
		return
	}

	if req.FulfillmentMode == "pickup" && req.PickupPointID == nil {
		writeError(c, http.StatusBadRequest, 40014, "pickupPointId is required", nil)
		return
	}
	if req.FulfillmentMode == "delivery" && req.AddressID == nil {
		writeError(c, http.StatusBadRequest, 40014, "addressId is required", nil)
		return
	}

	var created Order
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var group Group
		if err := tx.First(&group, req.GroupID).Error; err != nil {
			return err
		}

		if group.Status != "ongoing" {
			return fmt.Errorf("GROUP_NOT_ONGOING")
		}
		if time.Now().After(group.CutoffAt) {
			return fmt.Errorf("GROUP_CUTOFF")
		}

		created = Order{
			OrderNo:          fmt.Sprintf("GM%d", time.Now().UnixNano()),
			GroupID:          req.GroupID,
			UserID:           1,
			Status:           "joined",
			FulfillmentMode:  req.FulfillmentMode,
			ContactName:      req.ContactName,
			ContactPhone:     req.ContactPhone,
			PickupPointID:    req.PickupPointID,
			AddressID:        req.AddressID,
			Remark:           req.Remark,
			CutoffAtSnapshot: group.CutoffAt,
		}

		if req.AddressID != nil {
			var addr Address
			if err := tx.First(&addr, *req.AddressID).Error; err == nil {
				created.DeliveryAddressSnapshot = strings.TrimSpace(addr.Province + addr.City + addr.District + addr.DetailAddress)
			}
		}

		var goodsAmount float64
		for _, item := range req.Items {
			if item.Quantity <= 0 {
				return fmt.Errorf("INVALID_QUANTITY")
			}

			var groupItem GroupItem
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&groupItem, item.GroupItemID).Error; err != nil {
				return err
			}

			if groupItem.GroupID != req.GroupID {
				return fmt.Errorf("ITEM_GROUP_MISMATCH")
			}
			if groupItem.Status != "active" {
				return fmt.Errorf("ITEM_NOT_ACTIVE")
			}
			if item.Quantity > groupItem.StockAvailableSnapshot {
				return fmt.Errorf("INSUFFICIENT_STOCK:%d", groupItem.ID)
			}
			if groupItem.LimitPerUserSnapshot > 0 && item.Quantity > groupItem.LimitPerUserSnapshot {
				return fmt.Errorf("LIMIT_EXCEEDED")
			}

			goodsAmount += groupItem.PriceSnapshot * float64(item.Quantity)
		}

		created.GoodsAmount = goodsAmount
		created.PayableAmount = goodsAmount
		created.PaidAmount = goodsAmount

		if req.UserCouponID != nil {
			var userCoupon UserCoupon
			if err := tx.Where("id = ? AND user_id = ?", *req.UserCouponID, 1).First(&userCoupon).Error; err == nil {
				var coupon Coupon
				if err := tx.First(&coupon, userCoupon.CouponID).Error; err == nil && userCoupon.Status == "unused" && coupon.Status == "active" {
					now := time.Now()
					if now.After(userCoupon.ValidFrom) && now.Before(userCoupon.ValidTo) && goodsAmount >= coupon.ThresholdAmount {
						discount := coupon.Amount
						if discount > goodsAmount {
							discount = goodsAmount
						}
						created.DiscountAmount = discount
						created.PayableAmount = goodsAmount - discount
						created.PaidAmount = created.PayableAmount
					}
				}
			}
		}

		if err := tx.Create(&created).Error; err != nil {
			return err
		}

		if req.UserCouponID != nil && created.DiscountAmount > 0 {
			now := time.Now()
			if err := tx.Model(&UserCoupon{}).Where("id = ? AND user_id = ? AND status = ?", *req.UserCouponID, 1, "unused").Updates(map[string]any{
				"status":  "used",
				"used_at": &now,
				"order_id": created.ID,
			}).Error; err != nil {
				return err
			}
		}

		for _, item := range req.Items {
			var groupItem GroupItem
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&groupItem, item.GroupItemID).Error; err != nil {
				return err
			}

			before := groupItem.StockAvailableSnapshot
			groupItem.StockAvailableSnapshot -= item.Quantity
			if err := tx.Save(&groupItem).Error; err != nil {
				return err
			}

			orderItem := OrderItem{
				OrderID:             created.ID,
				GroupItemID:         groupItem.ID,
				ProductID:           groupItem.ProductID,
				ProductSKUID:        groupItem.ProductSKUID,
				ProductNameSnapshot: groupItem.ProductNameSnapshot,
				SKUNameSnapshot:     groupItem.SKUNameSnapshot,
				UnitPriceSnapshot:   groupItem.PriceSnapshot,
				Quantity:            item.Quantity,
				SubtotalAmount:      groupItem.PriceSnapshot * float64(item.Quantity),
				TasteRemark:         item.TasteRemark,
				ItemStatus:          "normal",
			}
			if err := tx.Create(&orderItem).Error; err != nil {
				return err
			}

			inv := InventoryLog{
				GroupID:      groupItem.GroupID,
				GroupItemID:  groupItem.ID,
				ProductSKUID: groupItem.ProductSKUID,
				OrderID:      &created.ID,
				ChangeType:   "reserve",
				ChangeQty:    -item.Quantity,
				BeforeStock:  before,
				AfterStock:   groupItem.StockAvailableSnapshot,
				OperatorID:   &created.UserID,
				OperatorRole: "user",
				Remark:       "create order reserve stock",
			}
			if err := tx.Create(&inv).Error; err != nil {
				return err
			}
		}

		if err := tx.Preload("Items").First(&created, created.ID).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		switch {
		case err.Error() == "GROUP_NOT_ONGOING" || err.Error() == "GROUP_CUTOFF":
			writeError(c, http.StatusBadRequest, 40010, "活动已截单或不可下单", nil)
		case err.Error() == "LIMIT_EXCEEDED":
			writeError(c, http.StatusBadRequest, 40012, "超出限购", nil)
		case err.Error() == "INVALID_QUANTITY":
			writeError(c, http.StatusBadRequest, 40001, "invalid quantity", nil)
		case len(err.Error()) >= 18 && err.Error()[:18] == "INSUFFICIENT_STOCK":
			writeError(c, http.StatusBadRequest, 40011, "库存不足", nil)
		default:
			writeError(c, http.StatusInternalServerError, 50000, "create order failed", gin.H{"detail": err.Error()})
		}
		return
	}

	writeSuccess(c, created)
	s.seedNotificationForOrder(created.ID)
}

func (s Server) cancelOrder(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, 40001, "invalid order id", nil)
		return
	}

	var req cancelOrderRequest
	_ = c.ShouldBindJSON(&req)

	err = s.db.Transaction(func(tx *gorm.DB) error {
		var order Order
		if err := tx.Preload("Items").First(&order, id).Error; err != nil {
			return err
		}
		if order.UserID != 1 {
			return fmt.Errorf("FORBIDDEN")
		}
		if order.Status != "joined" {
			return fmt.Errorf("ORDER_NOT_CANCELABLE")
		}

		now := time.Now()
		order.Status = "cancelled"
		order.CancelledAt = &now
		if req.Reason != "" {
			order.Remark = req.Reason
		}
		if err := tx.Save(&order).Error; err != nil {
			return err
		}

		if order.DiscountAmount > 0 {
			_ = tx.Model(&UserCoupon{}).Where("order_id = ? AND user_id = ?", order.ID, 1).Updates(map[string]any{
				"status":  "unused",
				"used_at": nil,
				"order_id": nil,
			}).Error
		}

		for _, item := range order.Items {
			var groupItem GroupItem
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&groupItem, item.GroupItemID).Error; err != nil {
				return err
			}
			before := groupItem.StockAvailableSnapshot
			groupItem.StockAvailableSnapshot += item.Quantity
			if err := tx.Save(&groupItem).Error; err != nil {
				return err
			}
			inv := InventoryLog{
				GroupID:      groupItem.GroupID,
				GroupItemID:  groupItem.ID,
				ProductSKUID: groupItem.ProductSKUID,
				OrderID:      &order.ID,
				ChangeType:   "release",
				ChangeQty:    item.Quantity,
				BeforeStock:  before,
				AfterStock:   groupItem.StockAvailableSnapshot,
				OperatorID:   &order.UserID,
				OperatorRole: "user",
				Remark:       "cancel order release stock",
			}
			if err := tx.Create(&inv).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		switch err.Error() {
		case "FORBIDDEN":
			writeError(c, http.StatusForbidden, 40003, "无权限", nil)
		case "ORDER_NOT_CANCELABLE":
			writeError(c, http.StatusBadRequest, 40015, "订单状态不可取消", nil)
		default:
			writeError(c, http.StatusInternalServerError, 50000, "cancel order failed", gin.H{"detail": err.Error()})
		}
		return
	}

	writeSuccess(c, gin.H{"orderId": id, "status": "cancelled"})
}

func (s Server) cutoffGroup(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, 40001, "invalid group id", nil)
		return
	}

	var req cutoffGroupRequest
	_ = c.ShouldBindJSON(&req)

	err = s.db.Transaction(func(tx *gorm.DB) error {
		var group Group
		if err := tx.First(&group, id).Error; err != nil {
			return err
		}
		if group.Status != "ongoing" {
			return fmt.Errorf("GROUP_NOT_ONGOING")
		}
		group.Status = "cutoff"
		if err := tx.Save(&group).Error; err != nil {
			return err
		}
		if err := tx.Model(&Order{}).Where("group_id = ? AND status = ?", id, "joined").Updates(map[string]any{"status": "cutoff_locked"}).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		if err.Error() == "GROUP_NOT_ONGOING" {
			writeError(c, http.StatusBadRequest, 40015, "团状态不可截单", nil)
			return
		}
		writeError(c, http.StatusInternalServerError, 50000, "cutoff group failed", gin.H{"detail": err.Error()})
		return
	}

	writeSuccess(c, gin.H{"groupId": id, "status": "cutoff", "remark": req.Reason})
}

func (s Server) getGroupSummary(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, 40001, "invalid group id", nil)
		return
	}

	var summary []summaryRow
	if err := s.db.Table("order_items oi").
		Select("oi.product_name_snapshot as product_name, oi.sku_name_snapshot as sku_name, SUM(oi.quantity) as total_qty").
		Joins("JOIN orders o ON o.id = oi.order_id").
		Where("o.group_id = ? AND o.status IN ?", id, []string{"joined", "cutoff_locked", "ready_for_pickup", "delivering", "completed"}).
		Group("oi.product_name_snapshot, oi.sku_name_snapshot").
		Scan(&summary).Error; err != nil {
		writeError(c, http.StatusInternalServerError, 50000, "query summary failed", nil)
		return
	}

	writeSuccess(c, gin.H{"groupId": id, "bySku": summary})
}
