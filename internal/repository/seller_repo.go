package repository

import (
	"errors"

	"yourapp/internal/model"

	"gorm.io/gorm"
)

type SellerRepository interface {
	Create(seller *model.Seller) error
	FindByID(id string) (*model.Seller, error)
	FindByUserID(userID string) (*model.Seller, error)
	FindBySlug(slug string) (*model.Seller, error)
	Update(seller *model.Seller) error
	Delete(sellerID string) error
}

type sellerRepository struct {
	db *gorm.DB
}

func NewSellerRepository(db *gorm.DB) SellerRepository {
	return &sellerRepository{db: db}
}

func (r *sellerRepository) Create(seller *model.Seller) error {
	return r.db.Create(seller).Error
}

func (r *sellerRepository) FindByID(id string) (*model.Seller, error) {
	var seller model.Seller
	err := r.db.Where("id = ?", id).Preload("User").First(&seller).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("seller not found")
		}
		return nil, err
	}
	return &seller, nil
}

func (r *sellerRepository) FindByUserID(userID string) (*model.Seller, error) {
	var seller model.Seller
	err := r.db.Where("user_id = ?", userID).Preload("User").First(&seller).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("seller not found")
		}
		return nil, err
	}
	return &seller, nil
}

func (r *sellerRepository) FindBySlug(slug string) (*model.Seller, error) {
	var seller model.Seller
	err := r.db.Where("shop_slug = ?", slug).Preload("User").First(&seller).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("seller not found")
		}
		return nil, err
	}
	return &seller, nil
}

func (r *sellerRepository) Update(seller *model.Seller) error {
	return r.db.Save(seller).Error
}

func (r *sellerRepository) Delete(sellerID string) error {
	// Soft delete
	result := r.db.Where("id = ?", sellerID).Delete(&model.Seller{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("seller not found")
	}
	return nil
}
