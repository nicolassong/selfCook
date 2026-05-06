package app

import "time"

type User struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	OpenID        string    `gorm:"column:open_id;size:64;uniqueIndex;not null" json:"openId"`
	UnionID       string    `gorm:"column:union_id;size:64" json:"unionId"`
	Nickname      string    `gorm:"size:64;not null" json:"nickname"`
	AvatarURL     string    `gorm:"column:avatar_url;size:255" json:"avatarUrl"`
	Phone         string    `gorm:"size:20" json:"phone"`
	Role          string    `gorm:"size:20;not null;default:user" json:"role"`
	Status        string    `gorm:"size:20;not null;default:active" json:"status"`
	PointsBalance int       `gorm:"column:points_balance;not null;default:0" json:"pointsBalance"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type Address struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"index;not null" json:"userId"`
	ContactName   string    `gorm:"size:50;not null" json:"contactName"`
	ContactPhone  string    `gorm:"size:20;not null" json:"contactPhone"`
	Province      string    `gorm:"size:50" json:"province"`
	City          string    `gorm:"size:50" json:"city"`
	District      string    `gorm:"size:50" json:"district"`
	DetailAddress string    `gorm:"size:255;not null" json:"detailAddress"`
	CommunityName string    `gorm:"size:100" json:"communityName"`
	Latitude      float64   `json:"latitude"`
	Longitude     float64   `json:"longitude"`
	IsDefault     bool      `gorm:"not null;default:false" json:"isDefault"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type Product struct {
	ID           uint         `gorm:"primaryKey" json:"id"`
	Name         string       `gorm:"size:100;not null" json:"name"`
	Subtitle     string       `gorm:"size:255" json:"subtitle"`
	CoverImage   string       `gorm:"size:255" json:"coverImage"`
	CategoryName string       `gorm:"size:50" json:"categoryName"`
	Description  string       `gorm:"type:text" json:"description"`
	Status       string       `gorm:"size:20;not null;default:on_sale" json:"status"`
	SortOrder    int          `gorm:"not null;default:0" json:"sortOrder"`
	CreatedAt    time.Time    `json:"createdAt"`
	UpdatedAt    time.Time    `json:"updatedAt"`
	SKUs         []ProductSKU `gorm:"foreignKey:ProductID" json:"skus,omitempty"`
}

type ProductSKU struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ProductID      uint      `gorm:"index;not null" json:"productId"`
	SKUName        string    `gorm:"column:sku_name;size:100;not null" json:"skuName"`
	SKUCode        string    `gorm:"column:sku_code;size:50;not null" json:"skuCode"`
	Price          float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	OriginalPrice  float64   `gorm:"type:decimal(10,2);not null;default:0" json:"originalPrice"`
	StockTotal     int       `gorm:"not null;default:0" json:"stockTotal"`
	StockAvailable int       `gorm:"not null;default:0" json:"stockAvailable"`
	LimitPerUser   int       `gorm:"not null;default:0" json:"limitPerUser"`
	LimitPerOrder  int       `gorm:"not null;default:0" json:"limitPerOrder"`
	Status         string    `gorm:"size:20;not null;default:active" json:"status"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type DailyMenu struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	MenuDate  time.Time       `gorm:"column:menu_date;type:date;uniqueIndex;not null" json:"menuDate"`
	Title     string          `gorm:"size:100;not null;default:''" json:"title"`
	Status    string          `gorm:"size:20;not null;default:active" json:"status"`
	Remark    string          `gorm:"size:255" json:"remark"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
	Items     []DailyMenuItem `gorm:"foreignKey:DailyMenuID" json:"items,omitempty"`
}

type DailyMenuItem struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	DailyMenuID    uint      `gorm:"index;not null" json:"dailyMenuId"`
	ProductID      uint      `gorm:"not null" json:"productId"`
	ProductSKUID   uint      `gorm:"column:product_sku_id;not null" json:"productSkuId"`
	StockTotal     int       `gorm:"not null;default:0" json:"stockTotal"`
	StockAvailable int       `gorm:"not null;default:0" json:"stockAvailable"`
	Price          float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	OriginalPrice  float64   `gorm:"type:decimal(10,2);not null;default:0" json:"originalPrice"`
	LimitPerUser   int       `gorm:"not null;default:0" json:"limitPerUser"`
	LimitPerOrder  int       `gorm:"not null;default:0" json:"limitPerOrder"`
	SortOrder      int       `gorm:"not null;default:0" json:"sortOrder"`
	Status         string    `gorm:"size:20;not null;default:active" json:"status"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type PickupPoint struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Name          string    `gorm:"size:100;not null" json:"name"`
	ContactName   string    `gorm:"size:50" json:"contactName"`
	ContactPhone  string    `gorm:"size:20" json:"contactPhone"`
	Address       string    `gorm:"size:255;not null" json:"address"`
	BusinessHours string    `gorm:"size:100" json:"businessHours"`
	Status        string    `gorm:"size:20;not null;default:active" json:"status"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type Group struct {
	ID                      uint        `gorm:"primaryKey" json:"id"`
	Title                   string      `gorm:"size:150;not null" json:"title"`
	CoverImage              string      `gorm:"size:255" json:"coverImage"`
	LeaderUserID            uint        `gorm:"index;not null" json:"leaderUserId"`
	Status                  string      `gorm:"size:20;not null;default:ongoing" json:"status"`
	StartAt                 time.Time   `json:"startAt"`
	CutoffAt                time.Time   `gorm:"index" json:"cutoffAt"`
	FulfillmentMode         string      `gorm:"size:20;not null;default:mixed" json:"fulfillmentMode"`
	AllowModifyBeforeCutoff bool        `gorm:"not null;default:false" json:"allowModifyBeforeCutoff"`
	ShowJoinList            bool        `gorm:"not null;default:false" json:"showJoinList"`
	PickupRuleDesc          string      `gorm:"size:255" json:"pickupRuleDesc"`
	DeliveryRuleDesc        string      `gorm:"size:255" json:"deliveryRuleDesc"`
	GroupNotice             string      `gorm:"type:text" json:"groupNotice"`
	CreatedAt               time.Time   `json:"createdAt"`
	UpdatedAt               time.Time   `json:"updatedAt"`
	Items                   []GroupItem `gorm:"foreignKey:GroupID" json:"items,omitempty"`
}

type GroupItem struct {
	ID                     uint      `gorm:"primaryKey" json:"id"`
	GroupID                uint      `gorm:"index;not null" json:"groupId"`
	ProductID              uint      `gorm:"not null" json:"productId"`
	ProductSKUID           uint      `gorm:"column:product_sku_id;not null" json:"productSkuId"`
	ProductNameSnapshot    string    `gorm:"size:100;not null" json:"productName"`
	SKUNameSnapshot        string    `gorm:"size:100;not null" json:"skuName"`
	CoverImageSnapshot     string    `gorm:"size:255" json:"coverImage"`
	PriceSnapshot          float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	OriginalPriceSnapshot  float64   `gorm:"type:decimal(10,2);not null;default:0" json:"originalPrice"`
	StockTotalSnapshot     int       `gorm:"not null;default:0" json:"stockTotal"`
	StockAvailableSnapshot int       `gorm:"not null;default:0" json:"stockAvailable"`
	LimitPerUserSnapshot   int       `gorm:"not null;default:0" json:"limitPerUser"`
	LimitPerOrderSnapshot  int       `gorm:"not null;default:0" json:"limitPerOrder"`
	Status                 string    `gorm:"size:20;not null;default:active" json:"status"`
	SortOrder              int       `gorm:"not null;default:0" json:"sortOrder"`
	CreatedAt              time.Time `json:"createdAt"`
	UpdatedAt              time.Time `json:"updatedAt"`
}

type Order struct {
	ID                      uint        `gorm:"primaryKey" json:"id"`
	OrderNo                 string      `gorm:"size:40;uniqueIndex;not null" json:"orderNo"`
	GroupID                 uint        `gorm:"index;not null" json:"groupId"`
	UserID                  uint        `gorm:"index;not null;default:1" json:"userId"`
	Status                  string      `gorm:"size:30;not null" json:"status"`
	FulfillmentMode         string      `gorm:"size:20;not null" json:"fulfillmentMode"`
	ContactName             string      `gorm:"size:50;not null" json:"contactName"`
	ContactPhone            string      `gorm:"size:20;not null" json:"contactPhone"`
	PickupPointID           *uint       `json:"pickupPointId"`
	AddressID               *uint       `json:"addressId"`
	DeliveryAddressSnapshot string      `gorm:"size:255" json:"deliveryAddressSnapshot"`
	GoodsAmount             float64     `gorm:"type:decimal(10,2);not null" json:"goodsAmount"`
	DiscountAmount          float64     `gorm:"type:decimal(10,2);not null;default:0" json:"discountAmount"`
	DeliveryFee             float64     `gorm:"type:decimal(10,2);not null;default:0" json:"deliveryFee"`
	PayableAmount           float64     `gorm:"type:decimal(10,2);not null" json:"payableAmount"`
	PaidAmount              float64     `gorm:"type:decimal(10,2);not null" json:"paidAmount"`
	Remark                  string      `gorm:"size:255" json:"remark"`
	CutoffAtSnapshot        time.Time   `json:"cutoffAtSnapshot"`
	CancelledAt             *time.Time  `json:"cancelledAt"`
	CompletedAt             *time.Time  `json:"completedAt"`
	CreatedAt               time.Time   `json:"createdAt"`
	UpdatedAt               time.Time   `json:"updatedAt"`
	Items                   []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
}

type OrderItem struct {
	ID                  uint      `gorm:"primaryKey" json:"id"`
	OrderID             uint      `gorm:"index;not null" json:"orderId"`
	GroupItemID         uint      `gorm:"not null" json:"groupItemId"`
	ProductID           uint      `gorm:"not null" json:"productId"`
	ProductSKUID        uint      `gorm:"column:product_sku_id;not null" json:"productSkuId"`
	ProductNameSnapshot string    `gorm:"size:100;not null" json:"productName"`
	SKUNameSnapshot     string    `gorm:"size:100;not null" json:"skuName"`
	UnitPriceSnapshot   float64   `gorm:"type:decimal(10,2);not null" json:"unitPrice"`
	Quantity            int       `gorm:"not null" json:"quantity"`
	SubtotalAmount      float64   `gorm:"type:decimal(10,2);not null" json:"subtotalAmount"`
	TasteRemark         string    `gorm:"size:100" json:"tasteRemark"`
	ItemStatus          string    `gorm:"size:20;not null;default:normal" json:"itemStatus"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
}

type InventoryLog struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	GroupID      uint      `gorm:"index" json:"groupId"`
	GroupItemID  uint      `gorm:"index" json:"groupItemId"`
	ProductSKUID uint      `gorm:"column:product_sku_id" json:"productSkuId"`
	OrderID      *uint     `gorm:"index" json:"orderId"`
	ChangeType   string    `gorm:"size:30;not null" json:"changeType"`
	ChangeQty    int       `json:"changeQty"`
	BeforeStock  int       `json:"beforeStock"`
	AfterStock   int       `json:"afterStock"`
	OperatorID   *uint     `json:"operatorId"`
	OperatorRole string    `gorm:"size:20" json:"operatorRole"`
	Remark       string    `gorm:"size:255" json:"remark"`
	CreatedAt    time.Time `json:"createdAt"`
}

type Coupon struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	Name            string    `gorm:"size:100;not null" json:"name"`
	CouponType      string    `gorm:"size:20;not null" json:"couponType"`
	Amount          float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
	ThresholdAmount float64   `gorm:"type:decimal(10,2);not null;default:0" json:"thresholdAmount"`
	ApplicableScope string    `gorm:"size:20;not null;default:all" json:"applicableScope"`
	Status          string    `gorm:"size:20;not null;default:active" json:"status"`
	ValidFrom       time.Time `json:"validFrom"`
	ValidTo         time.Time `json:"validTo"`
	TotalCount      int       `json:"totalCount"`
	PerUserLimit    int       `json:"perUserLimit"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type UserCoupon struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	CouponID   uint       `gorm:"index;not null" json:"couponId"`
	UserID     uint       `gorm:"index;not null" json:"userId"`
	Status     string     `gorm:"size:20;not null;default:unused" json:"status"`
	AcquiredAt time.Time  `json:"acquiredAt"`
	UsedAt     *time.Time `json:"usedAt"`
	OrderID    *uint      `json:"orderId"`
	ValidFrom  time.Time  `json:"validFrom"`
	ValidTo    time.Time  `json:"validTo"`
	Coupon     Coupon     `gorm:"foreignKey:CouponID" json:"coupon"`
}

type PointsLedger struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"index;not null" json:"userId"`
	ChangeValue  int       `json:"changeValue"`
	BalanceAfter int       `json:"balanceAfter"`
	SourceType   string    `gorm:"size:30;not null" json:"sourceType"`
	SourceID     *uint     `json:"sourceId"`
	Remark       string    `gorm:"size:255" json:"remark"`
	CreatedAt    time.Time `json:"createdAt"`
}

type Notification struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	UserID          uint      `gorm:"index;not null" json:"userId"`
	OrderID         *uint     `gorm:"index" json:"orderId"`
	GroupID         *uint     `gorm:"index" json:"groupId"`
	SceneCode       string    `gorm:"size:50;not null" json:"sceneCode"`
	TemplateID      string    `gorm:"size:100;not null" json:"templateId"`
	SendStatus      string    `gorm:"size:20;not null;default:pending" json:"sendStatus"`
	RequestPayload  string    `gorm:"type:text" json:"requestPayload"`
	ResponsePayload string    `gorm:"type:text" json:"responsePayload"`
	FailReason      string    `gorm:"size:255" json:"failReason"`
	SentAt          *time.Time `json:"sentAt"`
	CreatedAt       time.Time `json:"createdAt"`
}
