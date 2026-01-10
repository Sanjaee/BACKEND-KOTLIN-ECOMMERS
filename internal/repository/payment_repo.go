package repository

import (
	"time"
	"yourapp/internal/model"

	"gorm.io/gorm"
)

type PaymentRepository interface {
	Create(payment *model.Payment) error
	FindByID(id string) (*model.Payment, error)
	FindByOrderID(orderID string) (*model.Payment, error)
	FindByOrderNumber(orderNumber string) (*model.Payment, error)
	FindByMidtransTransactionID(transactionID string) (*model.Payment, error)
	FindPendingPayments() ([]*model.Payment, error) // Get all pending payments for background check
	Update(payment *model.Payment) error
	UpdateStatus(paymentID string, status model.PaymentStatus) error
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(payment *model.Payment) error {
	return r.db.Create(payment).Error
}

func (r *paymentRepository) FindByID(id string) (*model.Payment, error) {
	var payment model.Payment
	err := r.db.Preload("Order").
		Preload("Order.OrderItems").
		Preload("Order.OrderItems.Product").
		Where("id = ?", id).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) FindByOrderID(orderID string) (*model.Payment, error) {
	var payment model.Payment
	err := r.db.Preload("Order").
		Preload("Order.OrderItems").
		Preload("Order.OrderItems.Product").
		Where("order_uuid = ?", orderID).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) FindByOrderNumber(orderNumber string) (*model.Payment, error) {
	var payment model.Payment
	err := r.db.Preload("Order").
		Preload("Order.OrderItems").
		Preload("Order.OrderItems.Product").
		Where("order_id = ?", orderNumber).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) FindByMidtransTransactionID(transactionID string) (*model.Payment, error) {
	var payment model.Payment
	err := r.db.Preload("Order").
		Preload("Order.OrderItems").
		Preload("Order.OrderItems.Product").
		Where("midtrans_transaction_id = ?", transactionID).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) FindPendingPayments() ([]*model.Payment, error) {
	var payments []*model.Payment
	// Get all pending payments created in last 48 hours
	// We'll filter by transaction ID in Go code for reliability
	err := r.db.Where("status = ?", model.PaymentStatusPending).
		Where("created_at > ?", time.Now().Add(-48*time.Hour)). // Check payments created in last 48 hours
		Find(&payments).Error
	if err != nil {
		return nil, err
	}

	// Filter payments that have transaction ID
	var validPayments []*model.Payment
	for _, payment := range payments {
		if payment.MidtransTransactionID != nil && *payment.MidtransTransactionID != "" {
			validPayments = append(validPayments, payment)
		}
	}

	return validPayments, nil
}

func (r *paymentRepository) Update(payment *model.Payment) error {
	return r.db.Save(payment).Error
}

func (r *paymentRepository) UpdateStatus(paymentID string, status model.PaymentStatus) error {
	return r.db.Model(&model.Payment{}).
		Where("id = ?", paymentID).
		Update("status", status).Error
}
