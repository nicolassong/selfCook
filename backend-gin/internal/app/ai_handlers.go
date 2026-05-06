package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	gin "github.com/gin-gonic/gin"
)

type aiRecommendRequest struct {
	Message   string             `json:"message"`
	History   []aiMessageRequest `json:"history"`
	MenuItems []aiMenuItem       `json:"menuItems,omitempty"`
}

type aiMessageRequest struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type aiRecommendResponse struct {
	Reply string `json:"reply"`
}

type aiMenuItem struct {
	GroupID         uint    `json:"groupId"`
	GroupTitle      string  `json:"groupTitle"`
	CutoffAt        string  `json:"cutoffAt"`
	FulfillmentMode string  `json:"fulfillmentMode"`
	GroupNotice     string  `json:"groupNotice"`
	ProductID       uint    `json:"productId"`
	GroupItemID     uint    `json:"groupItemId"`
	Name            string  `json:"name"`
	SKUName         string  `json:"skuName"`
	Price           float64 `json:"price"`
	OriginalPrice   float64 `json:"originalPrice"`
	StockAvailable  int     `json:"stockAvailable"`
	LimitPerOrder   int     `json:"limitPerOrder"`
	CategoryName    string  `json:"categoryName"`
	Description     string  `json:"description"`
}

func (s Server) aiRecommend(c *gin.Context) {
	var req aiRecommendRequest
	body, err := io.ReadAll(io.LimitReader(c.Request.Body, 1<<20))
	if err != nil {
		writeError(c, http.StatusBadRequest, 400, "read request body failed", nil)
		return
	}
	body = bytes.TrimPrefix(body, []byte{0xEF, 0xBB, 0xBF})
	if err := json.Unmarshal(body, &req); err != nil {
		writeError(c, http.StatusBadRequest, 400, "invalid request body", nil)
		return
	}

	req.Message = strings.TrimSpace(req.Message)
	if req.Message == "" {
		writeError(c, http.StatusBadRequest, 400, "message is required", nil)
		return
	}

	menuItems, err := s.loadAIMenuItems()
	if err != nil {
		writeError(c, http.StatusInternalServerError, 500, "query menu items failed", nil)
		return
	}
	req.MenuItems = menuItems

	payload, err := json.Marshal(req)
	if err != nil {
		writeError(c, http.StatusInternalServerError, 500, "marshal ai request failed", nil)
		return
	}

	target := strings.TrimRight(s.cfg.AIServiceURL, "/") + "/crew/recommend"
	httpClient := &http.Client{Timeout: 60 * time.Second}
	httpReq, err := http.NewRequestWithContext(c.Request.Context(), http.MethodPost, target, bytes.NewReader(payload))
	if err != nil {
		writeError(c, http.StatusInternalServerError, 500, "create ai request failed", nil)
		return
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		writeError(c, http.StatusBadGateway, 502, "ai service unavailable", nil)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		writeError(c, http.StatusBadGateway, 502, "read ai response failed", nil)
		return
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		writeError(c, http.StatusBadGateway, 502, "ai service returned error", string(respBody))
		return
	}

	var aiResp aiRecommendResponse
	if err := json.Unmarshal(respBody, &aiResp); err != nil {
		writeError(c, http.StatusBadGateway, 502, "invalid ai response", nil)
		return
	}

	writeSuccess(c, aiResp)
}

func (s Server) loadAIMenuItems() ([]aiMenuItem, error) {
	type menuRow struct {
		GroupID         uint
		GroupTitle      string
		CutoffAt        time.Time
		FulfillmentMode string
		GroupNotice     string
		ProductID       uint
		GroupItemID     uint
		Name            string
		SKUName         string
		Price           float64
		OriginalPrice   float64
		StockAvailable  int
		LimitPerOrder   int
		CategoryName    string
		Description     string
	}

	var rows []menuRow
	err := s.db.Table("group_items gi").
		Select(`
			g.id AS group_id,
			g.title AS group_title,
			g.cutoff_at AS cutoff_at,
			g.fulfillment_mode AS fulfillment_mode,
			g.group_notice AS group_notice,
			gi.product_id AS product_id,
			gi.id AS group_item_id,
			gi.product_name_snapshot AS name,
			gi.sku_name_snapshot AS sku_name,
			gi.price_snapshot AS price,
			gi.original_price_snapshot AS original_price,
			gi.stock_available_snapshot AS stock_available,
			gi.limit_per_order_snapshot AS limit_per_order,
			p.category_name AS category_name,
			p.description AS description
		`).
		Joins("JOIN `groups` g ON g.id = gi.group_id").
		Joins("LEFT JOIN products p ON p.id = gi.product_id").
		Where("g.status = ? AND gi.status = ? AND gi.stock_available_snapshot > 0 AND g.cutoff_at >= ?", "ongoing", "active", time.Now()).
		Order("g.cutoff_at ASC, gi.sort_order ASC, gi.id ASC").
		Limit(80).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	items := make([]aiMenuItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, aiMenuItem{
			GroupID:         row.GroupID,
			GroupTitle:      row.GroupTitle,
			CutoffAt:        row.CutoffAt.Format(time.RFC3339),
			FulfillmentMode: row.FulfillmentMode,
			GroupNotice:     row.GroupNotice,
			ProductID:       row.ProductID,
			GroupItemID:     row.GroupItemID,
			Name:            row.Name,
			SKUName:         row.SKUName,
			Price:           row.Price,
			OriginalPrice:   row.OriginalPrice,
			StockAvailable:  row.StockAvailable,
			LimitPerOrder:   row.LimitPerOrder,
			CategoryName:    row.CategoryName,
			Description:     row.Description,
		})
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("no available menu items")
	}
	return items, nil
}
