package app

import (
	"net/http"
	"strconv"

	gin "github.com/gin-gonic/gin"
)

func (s Server) listAdminUserCoupons(c *gin.Context) {
	userID, err := strconv.Atoi(c.Query("userId"))
	if err != nil || userID <= 0 {
		writeError(c, http.StatusBadRequest, 40001, "invalid userId", nil)
		return
	}

	var list []UserCoupon
	query := s.db.Preload("Coupon").Where("user_id = ?", userID).Order("id desc")
	if status := c.Query("status"); status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}
	if err := query.Find(&list).Error; err != nil {
		writeError(c, http.StatusInternalServerError, 50000, "query user coupons failed", nil)
		return
	}
	writeSuccess(c, gin.H{"list": list})
}
