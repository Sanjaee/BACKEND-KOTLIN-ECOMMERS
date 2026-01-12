package service

import (
	"errors"
	"yourapp/internal/model"
	"yourapp/internal/repository"
)

type CartService interface {
	GetCart(userID string) (*model.Cart, error)
	AddItemToCart(userID string, req *AddCartItemRequest) (*model.CartItem, error)
	UpdateCartItem(userID string, cartItemID string, req *UpdateCartItemRequest) (*model.CartItem, error)
	RemoveCartItem(userID string, cartItemID string) error
	ClearCart(userID string) error
	GetCartItems(userID string) ([]model.CartItem, error)
}

type cartService struct {
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
}

type AddCartItemRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}

func NewCartService(
	cartRepo repository.CartRepository,
	productRepo repository.ProductRepository,
) CartService {
	return &cartService{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

func (s *cartService) GetCart(userID string) (*model.Cart, error) {
	cart, err := s.cartRepo.GetOrCreateByUserID(userID)
	if err != nil {
		return nil, err
	}
	return cart, nil
}

func (s *cartService) AddItemToCart(userID string, req *AddCartItemRequest) (*model.CartItem, error) {
	// Get or create cart
	cart, err := s.cartRepo.GetOrCreateByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Get product
	product, err := s.productRepo.FindByID(req.ProductID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Check if product is active
	if !product.IsActive {
		return nil, errors.New("product is not available")
	}

	// Check stock
	if product.Stock < req.Quantity {
		return nil, errors.New("insufficient stock")
	}

	// Check if item already exists in cart
	existingItem, err := s.cartRepo.GetCartItemByProductID(cart.ID, req.ProductID)
	if err == nil {
		// Update quantity if item exists
		newQuantity := existingItem.Quantity + req.Quantity
		if product.Stock < newQuantity {
			return nil, errors.New("insufficient stock")
		}
		existingItem.Quantity = newQuantity
		existingItem.Price = product.Price // Update price to current price
		if err := s.cartRepo.UpdateCartItem(existingItem); err != nil {
			return nil, err
		}
		return existingItem, nil
	}

	// Create new cart item
	cartItem := &model.CartItem{
		CartID:    cart.ID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		Price:     product.Price,
	}

	if err := s.cartRepo.AddCartItem(cartItem); err != nil {
		return nil, err
	}

	// Load product details
	cartItem, err = s.cartRepo.GetCartItemByID(cartItem.ID)
	if err != nil {
		return nil, err
	}

	return cartItem, nil
}

func (s *cartService) UpdateCartItem(userID string, cartItemID string, req *UpdateCartItemRequest) (*model.CartItem, error) {
	// Get cart to verify ownership
	cart, err := s.cartRepo.GetByUserID(userID)
	if err != nil {
		return nil, errors.New("cart not found")
	}

	// Get cart item
	cartItem, err := s.cartRepo.GetCartItemByID(cartItemID)
	if err != nil {
		return nil, errors.New("cart item not found")
	}

	// Verify cart item belongs to user's cart
	if cartItem.CartID != cart.ID {
		return nil, errors.New("unauthorized")
	}

	// Get product to check stock
	product, err := s.productRepo.FindByID(cartItem.ProductID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Check stock
	if product.Stock < req.Quantity {
		return nil, errors.New("insufficient stock")
	}

	// Update cart item
	cartItem.Quantity = req.Quantity
	cartItem.Price = product.Price // Update price to current price

	if err := s.cartRepo.UpdateCartItem(cartItem); err != nil {
		return nil, err
	}

	// Reload with product details
	cartItem, err = s.cartRepo.GetCartItemByID(cartItemID)
	if err != nil {
		return nil, err
	}

	return cartItem, nil
}

func (s *cartService) RemoveCartItem(userID string, cartItemID string) error {
	// Get cart to verify ownership
	cart, err := s.cartRepo.GetByUserID(userID)
	if err != nil {
		return errors.New("cart not found")
	}

	// Get cart item
	cartItem, err := s.cartRepo.GetCartItemByID(cartItemID)
	if err != nil {
		return errors.New("cart item not found")
	}

	// Verify cart item belongs to user's cart
	if cartItem.CartID != cart.ID {
		return errors.New("unauthorized")
	}

	return s.cartRepo.DeleteCartItem(cartItemID)
}

func (s *cartService) ClearCart(userID string) error {
	cart, err := s.cartRepo.GetByUserID(userID)
	if err != nil {
		return errors.New("cart not found")
	}

	return s.cartRepo.ClearCart(cart.ID)
}

func (s *cartService) GetCartItems(userID string) ([]model.CartItem, error) {
	cart, err := s.cartRepo.GetByUserID(userID)
	if err != nil {
		return nil, errors.New("cart not found")
	}

	return s.cartRepo.GetCartItems(cart.ID)
}
