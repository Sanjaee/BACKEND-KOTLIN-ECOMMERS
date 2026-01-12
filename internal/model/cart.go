package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Cart struct {
	ID        string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    string    `gorm:"type:uuid;not null;uniqueIndex;index" json:"user_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	User       User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CartItems  []CartItem  `gorm:"foreignKey:CartID" json:"cart_items,omitempty"`
}

func (c *Cart) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}

func (Cart) TableName() string {
	return "carts"
}

type CartItem struct {
	ID        string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CartID    string    `gorm:"type:uuid;not null;index" json:"cart_id"`
	ProductID string    `gorm:"type:uuid;not null;index" json:"product_id"`
	Quantity  int       `gorm:"not null;default:1" json:"quantity"`
	Price     int       `gorm:"not null" json:"price"` // Price at time of adding to cart
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Cart    Cart    `gorm:"foreignKey:CartID" json:"cart,omitempty"`
	Product Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

func (ci *CartItem) BeforeCreate(tx *gorm.DB) error {
	if ci.ID == "" {
		ci.ID = uuid.New().String()
	}
	return nil
}

func (CartItem) TableName() string {
	return "cart_items"
}
