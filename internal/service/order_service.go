package service

import (
	"errors"
	"yourapp/internal/model"
	"yourapp/internal/repository"
)

type OrderService interface {
	CreateOrder(userID string, req *CreateOrderRequest) (*model.Order, error)
	GetOrderByID(orderID string, userID string) (*model.Order, error)
	GetOrdersByUserID(userID string, page, limit int) ([]model.Order, int64, error)
	UpdateOrderStatus(orderID string, status string) error
}

type orderService struct {
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
	addressRepo repository.AddressRepository
}

type CreateOrderRequest struct {
	ShippingAddressID string                   `json:"shipping_address_id"` // Optional: will auto-create if not found
	Items             []CreateOrderItemRequest `json:"items" binding:"required,min=1"`
	Subtotal          int                      `json:"subtotal" binding:"required"`
	ShippingCost      int                      `json:"shipping_cost"`
	InsuranceCost     int                      `json:"insurance_cost"`
	WarrantyCost      int                      `json:"warranty_cost"`
	ServiceFee        int                      `json:"service_fee"`
	ApplicationFee    int                      `json:"application_fee"`
	TotalDiscount     int                      `json:"total_discount"`
	Bonus             int                      `json:"bonus"`
	Notes             *string                  `json:"notes,omitempty"`
}

type CreateOrderItemRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

func NewOrderService(
	orderRepo repository.OrderRepository,
	productRepo repository.ProductRepository,
	addressRepo repository.AddressRepository,
) OrderService {
	return &orderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		addressRepo: addressRepo,
	}
}

func (s *orderService) CreateOrder(userID string, req *CreateOrderRequest) (*model.Order, error) {
	// Validate or auto-create shipping address
	var address *model.Address
	var err error

	// If shipping_address_id is provided, try to find it
	if req.ShippingAddressID != "" && req.ShippingAddressID != "ADDR_1" {
		address, err = s.addressRepo.FindByID(req.ShippingAddressID)
		if err != nil {
			// Address ID not found, auto-create default address
			address = s.createDefaultAddress(userID)
			if err := s.addressRepo.Create(address); err != nil {
				return nil, errors.New("failed to create default address: " + err.Error())
			}
		} else if address.UserID != userID {
			return nil, errors.New("shipping address does not belong to user")
		}
		// If address found and belongs to user, use it
	} else {
		// No valid shipping_address_id provided, check if user has default address
		defaultAddr, err := s.addressRepo.FindDefaultByUserID(userID)
		if err == nil && defaultAddr != nil {
			address = defaultAddr
		} else {
			// No default address found, create one with static data
			address = s.createDefaultAddress(userID)
			if err := s.addressRepo.Create(address); err != nil {
				return nil, errors.New("failed to create default address: " + err.Error())
			}
		}
	}

	// Validate products and create order items
	var orderItems []model.OrderItem
	var totalAmount int

	for _, item := range req.Items {
		product, err := s.productRepo.FindByID(item.ProductID)
		if err != nil {
			return nil, errors.New("product not found: " + item.ProductID)
		}
		if !product.IsActive {
			return nil, errors.New("product is not active: " + item.ProductID)
		}
		if product.Stock < item.Quantity {
			return nil, errors.New("insufficient stock for product: " + product.Name)
		}

		subtotal := product.Price * item.Quantity
		totalAmount += subtotal

		orderItem := model.OrderItem{
			ProductID:   product.ID,
			ProductName: product.Name,
			Quantity:    item.Quantity,
			Price:       product.Price,
			Subtotal:    subtotal,
		}
		orderItems = append(orderItems, orderItem)
	}

	// Calculate total amount
	totalAmount = req.Subtotal + req.ShippingCost + req.InsuranceCost + req.WarrantyCost +
		req.ServiceFee + req.ApplicationFee - req.Bonus - req.TotalDiscount

	// Create order
	order := &model.Order{
		UserID:            userID,
		ShippingAddressID: address.ID,
		Subtotal:          req.Subtotal,
		ShippingCost:      req.ShippingCost,
		InsuranceCost:     req.InsuranceCost,
		WarrantyCost:      req.WarrantyCost,
		ServiceFee:        req.ServiceFee,
		ApplicationFee:    req.ApplicationFee,
		TotalDiscount:     req.TotalDiscount,
		Bonus:             req.Bonus,
		TotalAmount:       totalAmount,
		Status:            "pending",
		Notes:             req.Notes,
		OrderItems:        orderItems,
	}

	if err := s.orderRepo.Create(order); err != nil {
		return nil, err
	}

	// Update product stock
	for _, item := range req.Items {
		product, _ := s.productRepo.FindByID(item.ProductID)
		if product != nil {
			product.Stock -= item.Quantity
			s.productRepo.Update(product)
		}
	}

	return order, nil
}

func (s *orderService) GetOrderByID(orderID string, userID string) (*model.Order, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, errors.New("order not found")
	}
	if order.UserID != userID {
		return nil, errors.New("order does not belong to user")
	}
	return order, nil
}

func (s *orderService) GetOrdersByUserID(userID string, page, limit int) ([]model.Order, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	return s.orderRepo.FindByUserID(userID, page, limit)
}

func (s *orderService) UpdateOrderStatus(orderID string, status string) error {
	validStatuses := map[string]bool{
		"pending":    true,
		"processing": true,
		"shipped":    true,
		"delivered":  true,
		"cancelled":  true,
	}
	if !validStatuses[status] {
		return errors.New("invalid order status")
	}
	return s.orderRepo.UpdateStatus(orderID, status)
}

// createDefaultAddress creates a default static address for a user
// This uses static data matching the CheckoutViewModel in Android app
func (s *orderService) createDefaultAddress(userID string) *model.Address {
	return &model.Address{
		UserID:        userID,
		Label:         "Rumah",
		RecipientName: "Ahmad",
		Phone:         "+6281234567890",
		AddressLine1:  "JL.PELITA RT07/RW01 KONTRAKAN HJ.KEPOY",
		AddressLine2:  nil,
		City:          "Jakarta",
		Province:      "DKI Jakarta",
		PostalCode:    "12345",
		IsDefault:     true,
	}
}
