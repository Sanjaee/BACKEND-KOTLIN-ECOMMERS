package repository

import (
	"yourapp/internal/model"

	"gorm.io/gorm"
)

type AddressRepository interface {
	Create(address *model.Address) error
	FindByID(id string) (*model.Address, error)
	FindByUserID(userID string) ([]model.Address, error)
	FindDefaultByUserID(userID string) (*model.Address, error)
	Update(address *model.Address) error
	Delete(id string) error
}

type addressRepository struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) AddressRepository {
	return &addressRepository{db: db}
}

func (r *addressRepository) Create(address *model.Address) error {
	return r.db.Create(address).Error
}

func (r *addressRepository) FindByID(id string) (*model.Address, error) {
	var address model.Address
	err := r.db.Where("id = ?", id).First(&address).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}

func (r *addressRepository) FindByUserID(userID string) ([]model.Address, error) {
	var addresses []model.Address
	err := r.db.Where("user_id = ?", userID).Order("is_default DESC, created_at DESC").Find(&addresses).Error
	return addresses, err
}

func (r *addressRepository) FindDefaultByUserID(userID string) (*model.Address, error) {
	var address model.Address
	err := r.db.Where("user_id = ? AND is_default = ?", userID, true).First(&address).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}

func (r *addressRepository) Update(address *model.Address) error {
	return r.db.Save(address).Error
}

func (r *addressRepository) Delete(id string) error {
	return r.db.Delete(&model.Address{}, "id = ?", id).Error
}
