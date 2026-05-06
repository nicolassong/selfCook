package app

import (
	"net/http"
	"strconv"
	"time"

	gin "github.com/gin-gonic/gin"
)

type createAddressRequest struct {
	ContactName   string  `json:"contactName" binding:"required"`
	ContactPhone  string  `json:"contactPhone" binding:"required"`
	Province      string  `json:"province"`
	City          string  `json:"city"`
	District      string  `json:"district"`
	DetailAddress string  `json:"detailAddress" binding:"required"`
	CommunityName string  `json:"communityName"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	IsDefault     bool    `json:"isDefault"`
}

func (s Server) listMyAddresses(c *gin.Context) {
	var addresses []Address
	if err := s.db.Where("user_id = ?", 1).Order("is_default desc, id desc").Find(&addresses).Error; err != nil {
		writeError(c, http.StatusInternalServerError, 50000, "query addresses failed", nil)
		return
	}
	writeSuccess(c, gin.H{"list": addresses})
}

func (s Server) createAddress(c *gin.Context) {
	var req createAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, 40001, "invalid request", gin.H{"detail": err.Error()})
		return
	}

	address := Address{
		UserID:        1,
		ContactName:   req.ContactName,
		ContactPhone:  req.ContactPhone,
		Province:      req.Province,
		City:          req.City,
		District:      req.District,
		DetailAddress: req.DetailAddress,
		CommunityName: req.CommunityName,
		Latitude:      req.Latitude,
		Longitude:     req.Longitude,
		IsDefault:     req.IsDefault,
	}

	if req.IsDefault {
		_ = s.db.Model(&Address{}).Where("user_id = ?", 1).Update("is_default", false).Error
	}
	if err := s.db.Create(&address).Error; err != nil {
		writeError(c, http.StatusInternalServerError, 50000, "create address failed", nil)
		return
	}
	writeSuccess(c, address)
}

func (s Server) updateAddress(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, 40001, "invalid address id", nil)
		return
	}
	var req createAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, 40001, "invalid request", gin.H{"detail": err.Error()})
		return
	}
	var address Address
	if err := s.db.Where("id = ? AND user_id = ?", id, 1).First(&address).Error; err != nil {
		writeError(c, http.StatusNotFound, 40004, "address not found", nil)
		return
	}
	if req.IsDefault {
		_ = s.db.Model(&Address{}).Where("user_id = ?", 1).Update("is_default", false).Error
	}
	address.ContactName = req.ContactName
	address.ContactPhone = req.ContactPhone
	address.Province = req.Province
	address.City = req.City
	address.District = req.District
	address.DetailAddress = req.DetailAddress
	address.CommunityName = req.CommunityName
	address.Latitude = req.Latitude
	address.Longitude = req.Longitude
	address.IsDefault = req.IsDefault
	if err := s.db.Save(&address).Error; err != nil {
		writeError(c, http.StatusInternalServerError, 50000, "update address failed", nil)
		return
	}
	writeSuccess(c, address)
}

func (s Server) deleteAddress(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, 40001, "invalid address id", nil)
		return
	}
	if err := s.db.Where("id = ? AND user_id = ?", id, 1).Delete(&Address{}).Error; err != nil {
		writeError(c, http.StatusInternalServerError, 50000, "delete address failed", nil)
		return
	}
	writeSuccess(c, gin.H{"id": id})
}

func (s Server) setDefaultAddress(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, 40001, "invalid address id", nil)
		return
	}
	_ = s.db.Model(&Address{}).Where("user_id = ?", 1).Update("is_default", false).Error
	if err := s.db.Model(&Address{}).Where("id = ? AND user_id = ?", id, 1).Update("is_default", true).Error; err != nil {
		writeError(c, http.StatusInternalServerError, 50000, "set default address failed", nil)
		return
	}
	writeSuccess(c, gin.H{"id": id, "isDefault": true})
}

func (s Server) listMyCoupons(c *gin.Context) {
	var list []UserCoupon
	query := s.db.Preload("Coupon").Where("user_id = ?", 1).Order("id desc")
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Find(&list).Error; err != nil {
		writeError(c, http.StatusInternalServerError, 50000, "query coupons failed", nil)
		return
	}
	writeSuccess(c, gin.H{"list": list})
}

func (s Server) listMyPoints(c *gin.Context) {
	var list []PointsLedger
	if err := s.db.Where("user_id = ?", 1).Order("id desc").Find(&list).Error; err != nil {
		writeError(c, http.StatusInternalServerError, 50000, "query points failed", nil)
		return
	}
	writeSuccess(c, gin.H{"list": list})
}

func (s Server) listMyNotifications(c *gin.Context) {
	var list []Notification
	if err := s.db.Where("user_id = ?", 1).Order("id desc").Find(&list).Error; err != nil {
		writeError(c, http.StatusInternalServerError, 50000, "query notifications failed", nil)
		return
	}
	writeSuccess(c, gin.H{"list": list})
}

func (s Server) getCurrentUser(c *gin.Context) {
	var user User
	if err := s.db.First(&user, 1).Error; err != nil {
		writeError(c, http.StatusNotFound, 40004, "user not found", nil)
		return
	}
	writeSuccess(c, user)
}

func (s Server) getOrderByNo(c *gin.Context) {
	orderNo := c.Param("orderNo")
	var order Order
	if err := s.db.Preload("Items").Where("order_no = ?", orderNo).First(&order).Error; err != nil {
		writeError(c, http.StatusNotFound, 40004, "order not found", nil)
		return
	}
	writeSuccess(c, order)
}

func (s Server) seedNotificationForOrder(orderID uint) {
	now := time.Now()
	notice := Notification{UserID: 1, OrderID: &orderID, SceneCode: "order_created", TemplateID: "tmpl_order_created", SendStatus: "mock_sent", RequestPayload: "{}", ResponsePayload: "{\"ok\":true}", SentAt: &now}
	_ = s.db.Create(&notice).Error
}
