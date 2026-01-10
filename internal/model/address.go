package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Address struct {
	ID            string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID        string         `gorm:"type:uuid;not null;index" json:"user_id"`
	Label         string         `gorm:"type:varchar(100)" json:"label"` // e.g., Home, Office
	RecipientName string         `gorm:"type:varchar(255);not null" json:"recipient_name"`
	Phone         string         `gorm:"type:varchar(20);not null" json:"phone"`
	AddressLine1  string         `gorm:"type:text;not null" json:"address_line1"`
	AddressLine2  *string        `gorm:"type:text" json:"address_line2,omitempty"`
	City          string         `gorm:"type:varchar(100);not null" json:"city"`
	Province      string         `gorm:"type:varchar(100);not null" json:"province"`
	PostalCode    string         `gorm:"type:varchar(10);not null" json:"postal_code"`
	IsDefault     bool           `gorm:"default:false;index" json:"is_default"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (a *Address) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return nil
}

func (Address) TableName() string {
	return "addresses"
}
