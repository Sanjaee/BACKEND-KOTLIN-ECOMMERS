package repository

import (
	"yourapp/internal/model"

	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *model.Product) error
	FindByID(id string) (*model.Product, error)
	FindBySKU(sku string) (*model.Product, error)
	FindAll(page, limit int, categoryID *string, featured *bool, activeOnly bool) ([]model.Product, int64, error)
	Update(product *model.Product) error
	Delete(id string) error
	CreateImage(image *model.ProductImage) error
	DeleteImage(id string) error
	FindImagesByProductID(productID string) ([]model.ProductImage, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *model.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) FindByID(id string) (*model.Product, error) {
	var product model.Product
	err := r.db.Preload("Category").Preload("ProductImages", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC")
	}).Where("id = ?", id).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) FindBySKU(sku string) (*model.Product, error) {
	var product model.Product
	err := r.db.Where("sku = ?", sku).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) FindAll(page, limit int, categoryID *string, featured *bool, activeOnly bool) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	query := r.db.Model(&model.Product{}).Preload("Category").Preload("ProductImages", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC")
	})

	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}

	if featured != nil {
		query = query.Where("is_featured = ?", *featured)
	}

	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&products).Error
	return products, total, err
}

func (r *productRepository) Update(product *model.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id string) error {
	return r.db.Delete(&model.Product{}, "id = ?", id).Error
}

func (r *productRepository) CreateImage(image *model.ProductImage) error {
	return r.db.Create(image).Error
}

func (r *productRepository) DeleteImage(id string) error {
	return r.db.Delete(&model.ProductImage{}, "id = ?", id).Error
}

func (r *productRepository) FindImagesByProductID(productID string) ([]model.ProductImage, error) {
	var images []model.ProductImage
	err := r.db.Where("product_id = ?", productID).Order("sort_order ASC").Find(&images).Error
	return images, err
}
