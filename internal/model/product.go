package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID          string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	SellerID    string         `gorm:"type:uuid;not null;index" json:"seller_id"`
	CategoryID  string         `gorm:"type:uuid;not null;index" json:"category_id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Description *string        `gorm:"type:text" json:"description,omitempty"`
	SKU         string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"sku"`
	Price       int            `gorm:"not null" json:"price"`
	Stock       int            `gorm:"default:0" json:"stock"`
	Weight      *int           `gorm:"type:int" json:"weight,omitempty"`
	Thumbnail   *string        `gorm:"type:text" json:"thumbnail,omitempty"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	IsFeatured  bool           `gorm:"default:false" json:"is_featured"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Seller        Seller         `gorm:"foreignKey:SellerID" json:"seller,omitempty"`
	Category      Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	ProductImages []ProductImage `gorm:"foreignKey:ProductID" json:"images,omitempty"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}

func (Product) TableName() string {
	return "products"
}

type ProductImage struct {
	ID        string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ProductID string    `gorm:"type:uuid;not null;index" json:"product_id"`
	ImageURL  string    `gorm:"type:text;not null" json:"image_url"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (pi *ProductImage) BeforeCreate(tx *gorm.DB) error {
	if pi.ID == "" {
		pi.ID = uuid.New().String()
	}
	return nil
}

func (ProductImage) TableName() string {
	return "product_images"
}
