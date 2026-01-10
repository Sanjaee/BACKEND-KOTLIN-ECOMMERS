package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusSuccess   PaymentStatus = "success"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusCancelled PaymentStatus = "cancelled"
	PaymentStatusExpired   PaymentStatus = "expired"
)

type PaymentMethod string

const (
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
	PaymentMethodGopay        PaymentMethod = "gopay"
	PaymentMethodCreditCard   PaymentMethod = "credit_card"
	PaymentMethodQRIS         PaymentMethod = "qris"
	PaymentMethodAlfamart     PaymentMethod = "alfamart"
)

type Payment struct {
	ID                    string        `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrderID               string        `gorm:"type:varchar(50);uniqueIndex;not null;index" json:"order_id"` // order_number from orders table
	OrderUUID             string        `gorm:"type:uuid;not null;index" json:"order_uuid"`                  // UUID from orders table
	MidtransTransactionID *string       `gorm:"type:varchar(255);index" json:"midtrans_transaction_id,omitempty"`
	Amount                int           `gorm:"not null" json:"amount"`
	TotalAmount           int           `gorm:"not null" json:"total_amount"`
	Status                PaymentStatus `gorm:"type:varchar(50);not null;default:'pending';index" json:"status"`
	PaymentMethod         PaymentMethod `gorm:"type:varchar(50);not null" json:"payment_method"`
	PaymentType           string        `gorm:"type:varchar(50);default:'midtrans'" json:"payment_type"`
	FraudStatus           *string       `gorm:"type:varchar(50)" json:"fraud_status,omitempty"`
	VANumber              *string       `gorm:"type:varchar(50)" json:"va_number,omitempty"`
	BankType              *string       `gorm:"type:varchar(50)" json:"bank_type,omitempty"`
	QRCodeURL             *string       `gorm:"type:text" json:"qr_code_url,omitempty"`
	ExpiryTime            *time.Time    `gorm:"type:timestamp" json:"expiry_time,omitempty"`
	MidtransResponse      *string       `gorm:"type:text" json:"midtrans_response,omitempty"` // Raw JSON response from Midtrans
	CreatedAt             time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt             time.Time     `gorm:"autoUpdateTime" json:"updated_at"`

	Order Order `gorm:"foreignKey:OrderUUID" json:"order,omitempty"`
}

func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}

func (Payment) TableName() string {
	return "payments"
}
