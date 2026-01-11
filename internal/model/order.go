package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	ID                string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrderNumber       string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"order_number"`
	UserID            string         `gorm:"type:uuid;not null;index" json:"user_id"`
	ShippingAddressID string         `gorm:"type:uuid;not null" json:"shipping_address_id"`
	Subtotal          int            `gorm:"not null" json:"subtotal"`
	ShippingCost      int            `gorm:"default:0" json:"shipping_cost"`
	InsuranceCost     int            `gorm:"default:0" json:"insurance_cost"`
	WarrantyCost      int            `gorm:"default:0" json:"warranty_cost"`
	ServiceFee        int            `gorm:"default:0" json:"service_fee"`
	ApplicationFee    int            `gorm:"default:0" json:"application_fee"`
	TotalDiscount     int            `gorm:"default:0" json:"total_discount"`
	Bonus             int            `gorm:"default:0" json:"bonus"`
	TotalAmount       int            `gorm:"not null" json:"total_amount"`
	Status            string         `gorm:"type:varchar(50);not null;default:'pending';index" json:"status"` // pending, processing, shipped, delivered, cancelled
	Notes             *string        `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt         time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	User            User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ShippingAddress Address     `gorm:"foreignKey:ShippingAddressID" json:"shipping_address,omitempty"`
	OrderItems      []OrderItem `gorm:"foreignKey:OrderID" json:"order_items,omitempty"`
	Payment         *Payment    `gorm:"foreignKey:OrderUUID" json:"payment,omitempty"`
}

func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.ID == "" {
		o.ID = uuid.New().String()
	}
	if o.OrderNumber == "" {
		o.OrderNumber = generateOrderNumber()
	}
	return nil
}

func (Order) TableName() string {
	return "orders"
}

type OrderItem struct {
	ID          string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrderID     string    `gorm:"type:uuid;not null;index" json:"order_id"`
	ProductID   string    `gorm:"type:uuid;not null;index" json:"product_id"`
	SellerID    string    `gorm:"type:uuid;not null;index" json:"seller_id"`
	ProductName string    `gorm:"type:varchar(255);not null" json:"product_name"`
	Quantity    int       `gorm:"not null" json:"quantity"`
	Price       int       `gorm:"not null" json:"price"` // Price at time of order
	Subtotal    int       `gorm:"not null" json:"subtotal"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`

	Order   Order  `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Product Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Seller  Seller  `gorm:"foreignKey:SellerID" json:"seller,omitempty"`
}

func (oi *OrderItem) BeforeCreate(tx *gorm.DB) error {
	if oi.ID == "" {
		oi.ID = uuid.New().String()
	}
	return nil
}

func (OrderItem) TableName() string {
	return "order_items"
}

// generateOrderNumber generates a unique order number
func generateOrderNumber() string {
	// Format: ORD-YYYYMMDD-HHMMSS-XXXX
	now := time.Now()
	return "ORD-" + now.Format("20060102") + "-" + now.Format("150405") + "-" + uuid.New().String()[:4]
}
