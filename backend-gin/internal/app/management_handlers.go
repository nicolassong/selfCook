package app

import (
	"net/http"
	"strconv"
	"time"

	gin "github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (s Server) createProduct(c *gin.Context) {
	var req createProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, 40001, "invalid request", gin.H{"detail": err.Error()})
		return
	}
	if len(req.SKUs) == 0 {
		writeError(c, http.StatusBadRequest, 40001, "skus is required", nil)
		return
	}

	var product Product
	err := s.db.Transaction(func(tx *gorm.DB) error {
		product = Product{
			Name:         req.Name,
			Subtitle:     req.Subtitle,
			CoverImage:   req.CoverImage,
			CategoryName: req.CategoryName,
			Description:  req.Description,
			Status:       "on_sale",
			SortOrder:    req.SortOrder,
		}
		if err := tx.Create(&product).Error; err != nil {
			return err
		}
		for _, sku := range req.SKUs {
			model := ProductSKU{
				ProductID:      product.ID,
				SKUName:        sku.SKUName,
				SKUCode:        sku.SKUCode,
				Price:          sku.Price,
				OriginalPrice:  sku.OriginalPrice,
				StockTotal:     sku.StockTotal,
				StockAvailable: sku.StockTotal,
				LimitPerUser:   sku.LimitPerUser,
				LimitPerOrder:  sku.LimitPerOrder,
				Status:         "active",
			}
			if err := tx.Create(&model).Error; err != nil {
				return err
			}
		}
		return tx.Preload("SKUs").First(&product, product.ID).Error
	})
	if err != nil {
		writeError(c, http.StatusInternalServerError, 50000, "create product failed", gin.H{"detail": err.Error()})
		return
	}
	writeSuccess(c, product)
}

func (s Server) createDailyMenu(c *gin.Context) {
	var req createDailyMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, 40001, "invalid request", gin.H{"detail": err.Error()})
		return
	}
	if len(req.Items) == 0 {
		writeError(c, http.StatusBadRequest, 40001, "items is required", nil)
		return
	}

	menuDate, err := time.Parse("2006-01-02", req.MenuDate)
	if err != nil {
		writeError(c, http.StatusBadRequest, 40001, "menuDate must be YYYY-MM-DD", nil)
		return
	}

	var menu DailyMenu
	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("menu_date = ?", menuDate).First(&menu).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				menu = DailyMenu{MenuDate: menuDate, Title: req.Title, Status: "active"}
				if menu.Title == "" {
					menu.Title = req.MenuDate + " 菜单"
				}
				if err := tx.Create(&menu).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			menu.Title = req.Title
			if menu.Title == "" {
				menu.Title = req.MenuDate + " 菜单"
			}
			if err := tx.Save(&menu).Error; err != nil {
				return err
			}
			if err := tx.Where("daily_menu_id = ?", menu.ID).Delete(&DailyMenuItem{}).Error; err != nil {
				return err
			}
		}

		for index, menuItem := range req.Items {
			var sku ProductSKU
			if err := tx.First(&sku, menuItem.ProductSKUID).Error; err != nil {
				return err
			}
			var product Product
			if err := tx.First(&product, sku.ProductID).Error; err != nil {
				return err
			}

			sortOrder := menuItem.SortOrder
			if sortOrder == 0 {
				sortOrder = index + 1
			}

			item := DailyMenuItem{
				DailyMenuID:    menu.ID,
				ProductID:      product.ID,
				ProductSKUID:   sku.ID,
				StockTotal:     menuItem.StockTotal,
				StockAvailable: menuItem.StockTotal,
				Price:          menuItem.Price,
				OriginalPrice:  menuItem.OriginalPrice,
				LimitPerUser:   menuItem.LimitPerUser,
				LimitPerOrder:  menuItem.LimitPerOrder,
				SortOrder:      sortOrder,
				Status:         "active",
			}
			if err := tx.Create(&item).Error; err != nil {
				return err
			}
		}

		return tx.Preload("Items").First(&menu, menu.ID).Error
	})
	if err != nil {
		writeError(c, http.StatusInternalServerError, 50000, "create daily menu failed", gin.H{"detail": err.Error()})
		return
	}
	writeSuccess(c, menu)
}

func (s Server) deleteDailyMenu(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, 40001, "invalid menu id", nil)
		return
	}
	if err := s.db.Delete(&DailyMenu{}, id).Error; err != nil {
		writeError(c, http.StatusInternalServerError, 50000, "delete daily menu failed", gin.H{"detail": err.Error()})
		return
	}
	writeSuccess(c, gin.H{"id": id})
}

func (s Server) createGroup(c *gin.Context) {
	var req createGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, 40001, "invalid request", gin.H{"detail": err.Error()})
		return
	}
	if !req.CutoffAt.After(req.StartAt) {
		writeError(c, http.StatusBadRequest, 40001, "cutoffAt must be after startAt", nil)
		return
	}

	var group Group
	err := s.db.Transaction(func(tx *gorm.DB) error {
		group = Group{
			Title:                   req.Title,
			CoverImage:              req.CoverImage,
			LeaderUserID:            req.LeaderUserID,
			Status:                  "ongoing",
			StartAt:                 req.StartAt,
			CutoffAt:                req.CutoffAt,
			FulfillmentMode:         req.FulfillmentMode,
			AllowModifyBeforeCutoff: req.AllowModifyBeforeCutoff,
			ShowJoinList:            req.ShowJoinList,
			PickupRuleDesc:          req.PickupRuleDesc,
			DeliveryRuleDesc:        req.DeliveryRuleDesc,
			GroupNotice:             req.GroupNotice,
		}
		if group.LeaderUserID == 0 {
			group.LeaderUserID = 1
		}
		if err := tx.Create(&group).Error; err != nil {
			return err
		}

		items := req.Items
		if req.MenuDate != "" {
			menuDate, err := time.Parse("2006-01-02", req.MenuDate)
			if err != nil {
				return err
			}
			var menu DailyMenu
			if err := tx.Preload("Items").Where("menu_date = ?", menuDate).First(&menu).Error; err != nil {
				return err
			}
			if len(menu.Items) == 0 {
				return gorm.ErrRecordNotFound
			}
			menuItems := make([]createGroupItemRequest, 0, len(menu.Items))
			for _, menuItem := range menu.Items {
				menuItems = append(menuItems, createGroupItemRequest{
					ProductSKUID:  menuItem.ProductSKUID,
					StockTotal:    menuItem.StockTotal,
					Price:         menuItem.Price,
					OriginalPrice: menuItem.OriginalPrice,
					LimitPerUser:  menuItem.LimitPerUser,
					LimitPerOrder: menuItem.LimitPerOrder,
					SortOrder:     menuItem.SortOrder,
				})
			}
			items = menuItems
		}

		if len(items) == 0 {
			return gorm.ErrRecordNotFound
		}

		for _, item := range items {
			var sku ProductSKU
			if err := tx.First(&sku, item.ProductSKUID).Error; err != nil {
				return err
			}
			var product Product
			if err := tx.First(&product, sku.ProductID).Error; err != nil {
				return err
			}
			groupItem := GroupItem{
				GroupID:                group.ID,
				ProductID:              product.ID,
				ProductSKUID:           sku.ID,
				ProductNameSnapshot:    product.Name,
				SKUNameSnapshot:        sku.SKUName,
				CoverImageSnapshot:     product.CoverImage,
				PriceSnapshot:          item.Price,
				OriginalPriceSnapshot:  item.OriginalPrice,
				StockTotalSnapshot:     item.StockTotal,
				StockAvailableSnapshot: item.StockTotal,
				LimitPerUserSnapshot:   item.LimitPerUser,
				LimitPerOrderSnapshot:  item.LimitPerOrder,
				Status:                 "active",
				SortOrder:              item.SortOrder,
			}
			if err := tx.Create(&groupItem).Error; err != nil {
				return err
			}
		}
		return tx.Preload("Items").First(&group, group.ID).Error
	})
	if err != nil {
		writeError(c, http.StatusInternalServerError, 50000, "create group failed", gin.H{"detail": err.Error()})
		return
	}
	writeSuccess(c, group)
}
