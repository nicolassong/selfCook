package app

import (
	"net/http"
	"strconv"
	"time"

	gin "github.com/gin-gonic/gin"
)

type createCouponRequest struct {
	Name            string  `json:"name" binding:"required"`
	CouponType      string  `json:"couponType" binding:"required"`
	Amount          float64 `json:"amount" binding:"required"`
	ThresholdAmount float64 `json:"thresholdAmount"`
	ApplicableScope string  `json:"applicableScope"`
	ValidDays       int     `json:"validDays"`
	TotalCount      int     `json:"totalCount"`
	PerUserLimit    int     `json:"perUserLimit"`
}

type grantCouponRequest struct {
	CouponID uint `json:"couponId" binding:"required"`
	UserID   uint `json:"userId" binding:"required"`
}

type adminCouponRow struct {
	Coupon
	GrantedCount int `json:"grantedCount"`
	UsedCount    int `json:"usedCount"`
}

func (s Server) listAdminCoupons(c *gin.Context) {
	var coupons []Coupon
	if err := s.db.Order("id desc").Find(&coupons).Error; err != nil {
		writeError(c, http.StatusInternalServerError, 50000, "query coupons failed", nil)
		return
	}
	rows := make([]adminCouponRow, 0, len(coupons))
	for _, coupon := range coupons {
		var grantedCount int64
		var usedCount int64
		_ = s.db.Model(&UserCoupon{}).Where("coupon_id = ?", coupon.ID).Count(&grantedCount).Error
		_ = s.db.Model(&UserCoupon{}).Where("coupon_id = ? AND status = ?", coupon.ID, "used").Count(&usedCount).Error
		rows = append(rows, adminCouponRow{Coupon: coupon, GrantedCount: int(grantedCount), UsedCount: int(usedCount)})
	}
	writeSuccess(c, gin.H{"list": rows})
}

func (s Server) disableCoupon(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, 40001, "invalid coupon id", nil)
		return
	}
	if err := s.db.Model(&Coupon{}).Where("id = ?", id).Update("status", "inactive").Error; err != nil {
		writeError(c, http.StatusInternalServerError, 50000, "disable coupon failed", nil)
		return
	}
	writeSuccess(c, gin.H{"id": id, "status": "inactive"})
}

func (s Server) createCoupon(c *gin.Context) {
	var req createCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, 40001, "invalid request", gin.H{"detail": err.Error()})
		return
	}
	if req.ValidDays <= 0 {
		req.ValidDays = 30
	}
	if req.TotalCount <= 0 {
		req.TotalCount = 1000
	}
	if req.PerUserLimit <= 0 {
		req.PerUserLimit = 1
	}
	coupon := Coupon{
		Name:            req.Name,
		CouponType:      req.CouponType,
		Amount:          req.Amount,
		ThresholdAmount: req.ThresholdAmount,
		ApplicableScope: fallbackString(req.ApplicableScope, "all"),
		Status:          "active",
		ValidFrom:       time.Now(),
		ValidTo:         time.Now().AddDate(0, 0, req.ValidDays),
		TotalCount:      req.TotalCount,
		PerUserLimit:    req.PerUserLimit,
	}
	if err := s.db.Create(&coupon).Error; err != nil {
		writeError(c, http.StatusInternalServerError, 50000, "create coupon failed", nil)
		return
	}
	writeSuccess(c, coupon)
}

func (s Server) grantCoupon(c *gin.Context) {
	var req grantCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, 40001, "invalid request", gin.H{"detail": err.Error()})
		return
	}
	var coupon Coupon
	if err := s.db.First(&coupon, req.CouponID).Error; err != nil {
		writeError(c, http.StatusNotFound, 40004, "coupon not found", nil)
		return
	}
	if coupon.Status != "active" {
		writeError(c, http.StatusBadRequest, 40001, "coupon is inactive", nil)
		return
	}
	var grantedCount int64
	_ = s.db.Model(&UserCoupon{}).Where("coupon_id = ?", coupon.ID).Count(&grantedCount).Error
	if coupon.TotalCount > 0 && grantedCount >= int64(coupon.TotalCount) {
		writeError(c, http.StatusBadRequest, 40001, "coupon out of stock", nil)
		return
	}
	var userOwnCount int64
	_ = s.db.Model(&UserCoupon{}).Where("coupon_id = ? AND user_id = ?", coupon.ID, req.UserID).Count(&userOwnCount).Error
	if coupon.PerUserLimit > 0 && userOwnCount >= int64(coupon.PerUserLimit) {
		writeError(c, http.StatusBadRequest, 40001, "user coupon limit exceeded", nil)
		return
	}
	userCoupon := UserCoupon{
		CouponID:   coupon.ID,
		UserID:     req.UserID,
		Status:     "unused",
		AcquiredAt: time.Now(),
		ValidFrom:  coupon.ValidFrom,
		ValidTo:    coupon.ValidTo,
	}
	if err := s.db.Create(&userCoupon).Error; err != nil {
		writeError(c, http.StatusInternalServerError, 50000, "grant coupon failed", nil)
		return
	}
	writeSuccess(c, userCoupon)
}

func fallbackString(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
