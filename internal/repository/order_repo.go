package repository

import (
	"yourapp/internal/model"

	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *model.Order) error
	FindByID(id string) (*model.Order, error)
	FindByOrderNumber(orderNumber string) (*model.Order, error)
	FindByUserID(userID string, page, limit int, status, paymentStatus string) ([]model.Order, int64, error)
	Update(order *model.Order) error
	UpdateStatus(orderID string, status string) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *model.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) FindByID(id string) (*model.Order, error) {
	var order model.Order
	err := r.db.Preload("User").
		Preload("ShippingAddress").
		Preload("OrderItems").
		Preload("OrderItems.Product").
		Preload("Payment").
		Where("id = ?", id).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) FindByOrderNumber(orderNumber string) (*model.Order, error) {
	var order model.Order
	err := r.db.Preload("User").
		Preload("ShippingAddress").
		Preload("OrderItems").
		Preload("OrderItems.Product").
		Preload("Payment").
		Where("order_number = ?", orderNumber).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) FindByUserID(userID string, page, limit int, status, paymentStatus string) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	offset := (page - 1) * limit

	// Base query with user_id filter
	query := r.db.Where("orders.user_id = ?", userID)

	// Filter by order status if provided
	if status != "" {
		validStatuses := map[string]bool{
			"pending":    true,
			"processing": true,
			"shipped":    true,
			"delivered":  true,
			"cancelled":  true,
		}
		if validStatuses[status] {
			query = query.Where("orders.status = ?", status)
		}
	}

	// Filter by payment status if provided
	if paymentStatus != "" {
		validPaymentStatuses := map[string]bool{
			"pending":   true,
			"success":   true,
			"failed":    true,
			"cancelled": true,
			"expired":   true,
		}
		if validPaymentStatuses[paymentStatus] {
			// Join with payments table to filter by payment status
			query = query.Joins("LEFT JOIN payments ON payments.order_uuid = orders.id").
				Where("payments.status = ?", paymentStatus)
		}
	}

	// Count total
	if err := query.Model(&model.Order{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch orders with preloads
	err := query.Preload("ShippingAddress").
		Preload("OrderItems").
		Preload("OrderItems.Product").
		Preload("Payment").
		Order("orders.created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&orders).Error

	return orders, total, err
}

func (r *orderRepository) Update(order *model.Order) error {
	return r.db.Save(order).Error
}

func (r *orderRepository) UpdateStatus(orderID string, status string) error {
	return r.db.Model(&model.Order{}).
		Where("id = ?", orderID).
		Update("status", status).Error
}
