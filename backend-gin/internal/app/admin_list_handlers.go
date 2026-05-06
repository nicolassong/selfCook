package app

import (
	"net/http"

	gin "github.com/gin-gonic/gin"
)

func (s Server) listAdminProducts(c *gin.Context) {
	var products []Product
	query := s.db.Preload("SKUs").Order("sort_order asc, id desc")
	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("name LIKE ? OR category_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Find(&products).Error; err != nil {
		writeError(c, http.StatusInternalServerError, 50000, "query admin products failed", nil)
		return
	}
	writeSuccess(c, gin.H{"list": products})
}

func (s Server) listAdminDailyMenus(c *gin.Context) {
	var menus []DailyMenu
	query := s.db.Preload("Items").Order("menu_date desc")
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Find(&menus).Error; err != nil {
		writeError(c, http.StatusInternalServerError, 50000, "query admin daily menus failed", nil)
		return
	}
	writeSuccess(c, gin.H{"list": menus})
}

func (s Server) listAdminGroups(c *gin.Context) {
	var groups []Group
	query := s.db.Preload("Items").Order("created_at desc")
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Find(&groups).Error; err != nil {
		writeError(c, http.StatusInternalServerError, 50000, "query admin groups failed", nil)
		return
	}
	writeSuccess(c, gin.H{"list": groups})
}

func (s Server) listAdminOrders(c *gin.Context) {
	var orders []Order
	query := s.db.Preload("Items").Order("created_at desc")
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if fulfillmentMode := c.Query("fulfillmentMode"); fulfillmentMode != "" {
		query = query.Where("fulfillment_mode = ?", fulfillmentMode)
	}
	if groupID := c.Query("groupId"); groupID != "" {
		query = query.Where("group_id = ?", groupID)
	}
	if err := query.Find(&orders).Error; err != nil {
		writeError(c, http.StatusInternalServerError, 50000, "query admin orders failed", nil)
		return
	}
	writeSuccess(c, gin.H{"list": orders})
}
