package repository

import (
	"yourapp/internal/model"

	"gorm.io/gorm"
)

type CartRepository interface {
	GetOrCreateByUserID(userID string) (*model.Cart, error)
	GetByUserID(userID string) (*model.Cart, error)
	GetCartItemByID(cartItemID string) (*model.CartItem, error)
	GetCartItemByProductID(cartID, productID string) (*model.CartItem, error)
	AddCartItem(cartItem *model.CartItem) error
	UpdateCartItem(cartItem *model.CartItem) error
	DeleteCartItem(cartItemID string) error
	ClearCart(cartID string) error
	GetCartItems(cartID string) ([]model.CartItem, error)
}

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) GetOrCreateByUserID(userID string) (*model.Cart, error) {
	var cart model.Cart
	err := r.db.Where("user_id = ?", userID).First(&cart).Error
	if err == gorm.ErrRecordNotFound {
		// Create new cart if not exists
		cart.UserID = userID
		if err := r.db.Create(&cart).Error; err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	
	// Preload cart items with product details
	err = r.db.Preload("CartItems").Preload("CartItems.Product").Preload("CartItems.Product.Seller").Preload("CartItems.Product.Category").Preload("CartItems.Product.ProductImages").Where("id = ?", cart.ID).First(&cart).Error
	return &cart, err
}

func (r *cartRepository) GetByUserID(userID string) (*model.Cart, error) {
	var cart model.Cart
	err := r.db.Preload("CartItems").Preload("CartItems.Product").Preload("CartItems.Product.Seller").Preload("CartItems.Product.Category").Preload("CartItems.Product.ProductImages").Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepository) GetCartItemByID(cartItemID string) (*model.CartItem, error) {
	var cartItem model.CartItem
	err := r.db.Preload("Product").Preload("Product.Seller").Preload("Product.Category").Preload("Product.ProductImages").Where("id = ?", cartItemID).First(&cartItem).Error
	if err != nil {
		return nil, err
	}
	return &cartItem, nil
}

func (r *cartRepository) GetCartItemByProductID(cartID, productID string) (*model.CartItem, error) {
	var cartItem model.CartItem
	err := r.db.Where("cart_id = ? AND product_id = ?", cartID, productID).First(&cartItem).Error
	if err != nil {
		return nil, err
	}
	return &cartItem, nil
}

func (r *cartRepository) AddCartItem(cartItem *model.CartItem) error {
	return r.db.Create(cartItem).Error
}

func (r *cartRepository) UpdateCartItem(cartItem *model.CartItem) error {
	return r.db.Save(cartItem).Error
}

func (r *cartRepository) DeleteCartItem(cartItemID string) error {
	return r.db.Delete(&model.CartItem{}, "id = ?", cartItemID).Error
}

func (r *cartRepository) ClearCart(cartID string) error {
	return r.db.Where("cart_id = ?", cartID).Delete(&model.CartItem{}).Error
}

func (r *cartRepository) GetCartItems(cartID string) ([]model.CartItem, error) {
	var cartItems []model.CartItem
	err := r.db.Preload("Product").Preload("Product.Seller").Preload("Product.Category").Preload("Product.ProductImages").Where("cart_id = ?", cartID).Find(&cartItems).Error
	return cartItems, err
}
