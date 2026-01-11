package model

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Seller struct {
	ID              string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID          string         `gorm:"type:uuid;uniqueIndex;not null;index" json:"user_id"`
	ShopName        string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"shop_name"`
	ShopSlug        string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"shop_slug"`
	ShopDescription *string        `gorm:"type:text" json:"shop_description,omitempty"`
	ShopLogo        *string        `gorm:"type:text" json:"shop_logo,omitempty"`
	ShopBanner      *string        `gorm:"type:text" json:"shop_banner,omitempty"`
	ShopAddress     *string        `gorm:"type:text" json:"shop_address,omitempty"`
	ShopCity        *string        `gorm:"type:varchar(100)" json:"shop_city,omitempty"`
	ShopProvince    *string        `gorm:"type:varchar(100)" json:"shop_province,omitempty"`
	ShopPhone       *string        `gorm:"type:varchar(20)" json:"shop_phone,omitempty"`
	ShopEmail       *string        `gorm:"type:varchar(255)" json:"shop_email,omitempty"`
	IsVerified      bool           `gorm:"default:false" json:"is_verified"`
	IsActive        bool           `gorm:"default:true" json:"is_active"`
	TotalProducts   int            `gorm:"default:0" json:"total_products"`
	TotalSales      int            `gorm:"default:0" json:"total_sales"`
	RatingAverage   float64        `gorm:"type:decimal(3,2);default:0.00" json:"rating_average"`
	TotalReviews    int            `gorm:"default:0" json:"total_reviews"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (s *Seller) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	if s.ShopSlug == "" && s.ShopName != "" {
		s.ShopSlug = generateSlug(s.ShopName)
	}
	return nil
}

func (s *Seller) BeforeUpdate(tx *gorm.DB) error {
	// Update slug jika shop_name berubah
	if tx.Statement.Changed("ShopName") {
		s.ShopSlug = generateSlug(s.ShopName)
	}
	return nil
}

func (Seller) TableName() string {
	return "sellers"
}

// generateSlug creates URL-friendly slug from shop name
func generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")
	// Remove special characters, keep only alphanumeric and hyphens
	var result strings.Builder
	for _, char := range slug {
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-' {
			result.WriteRune(char)
		}
	}
	return result.String()
}
