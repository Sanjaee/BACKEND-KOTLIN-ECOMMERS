package service

import (
	"errors"
	"fmt"

	"yourapp/internal/model"
	"yourapp/internal/repository"
)

type ProductService interface {
	CreateProduct(userID string, req CreateProductRequest) (*model.Product, error)
	GetProductByID(id string) (*model.Product, error)
	GetProducts(page, limit int, categoryID, featured, activeOnly *string) (*ProductListResponse, error)
	UpdateProduct(id string, req UpdateProductRequest) (*model.Product, error)
	DeleteProduct(id string) error
	AddProductImage(productID string, req AddProductImageRequest) (*model.ProductImage, error)
	DeleteProductImage(imageID string) error
}

type productService struct {
	productRepo  repository.ProductRepository
	categoryRepo repository.CategoryRepository
	sellerRepo   repository.SellerRepository
}

type CreateProductRequest struct {
	CategoryID  string  `json:"category_id" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description,omitempty"`
	SKU         string  `json:"sku" binding:"required"`
	Price       int     `json:"price" binding:"required,min=0"`
	Stock       int     `json:"stock" binding:"min=0"`
	Weight      *int    `json:"weight,omitempty"`
	Thumbnail   *string `json:"thumbnail,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
	IsFeatured  *bool   `json:"is_featured,omitempty"`
}

type UpdateProductRequest struct {
	CategoryID  *string `json:"category_id,omitempty"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	SKU         *string `json:"sku,omitempty"`
	Price       *int    `json:"price,omitempty"`
	Stock       *int    `json:"stock,omitempty"`
	Weight      *int    `json:"weight,omitempty"`
	Thumbnail   *string `json:"thumbnail,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
	IsFeatured  *bool   `json:"is_featured,omitempty"`
}

type AddProductImageRequest struct {
	ImageURL  string `json:"image_url" binding:"required"`
	SortOrder *int   `json:"sort_order,omitempty"`
}

type ProductListResponse struct {
	Products []model.Product `json:"products"`
	Total    int64           `json:"total"`
	Page     int             `json:"page"`
	Limit    int             `json:"limit"`
}

func NewProductService(productRepo repository.ProductRepository, categoryRepo repository.CategoryRepository, sellerRepo repository.SellerRepository) ProductService {
	return &productService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
		sellerRepo:   sellerRepo,
	}
}

func (s *productService) CreateProduct(userID string, req CreateProductRequest) (*model.Product, error) {
	// Get seller by userID (1 user 1 toko)
	seller, err := s.sellerRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("seller not found. Please create a shop first")
	}

	// Validate category exists
	_, err = s.categoryRepo.FindByID(req.CategoryID)
	if err != nil {
		return nil, errors.New("category not found")
	}

	// Check SKU uniqueness
	existing, _ := s.productRepo.FindBySKU(req.SKU)
	if existing != nil {
		return nil, errors.New("SKU already exists")
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	isFeatured := false
	if req.IsFeatured != nil {
		isFeatured = *req.IsFeatured
	}

	product := &model.Product{
		SellerID:    seller.ID,
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		SKU:         req.SKU,
		Price:       req.Price,
		Stock:       req.Stock,
		Weight:      req.Weight,
		Thumbnail:   req.Thumbnail,
		IsActive:    isActive,
		IsFeatured:  isFeatured,
	}

	if err := s.productRepo.Create(product); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return s.productRepo.FindByID(product.ID)
}

func (s *productService) GetProductByID(id string) (*model.Product, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("product not found")
	}
	return product, nil
}

func (s *productService) GetProducts(page, limit int, categoryID, featured, activeOnly *string) (*ProductListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	var categoryIDPtr *string
	if categoryID != nil && *categoryID != "" {
		categoryIDPtr = categoryID
	}

	var featuredPtr *bool
	if featured != nil && *featured != "" {
		feat := *featured == "true"
		featuredPtr = &feat
	}

	activeOnlyBool := false
	if activeOnly != nil && *activeOnly == "true" {
		activeOnlyBool = true
	}

	products, total, err := s.productRepo.FindAll(page, limit, categoryIDPtr, featuredPtr, activeOnlyBool)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	return &ProductListResponse{
		Products: products,
		Total:    total,
		Page:     page,
		Limit:    limit,
	}, nil
}

func (s *productService) UpdateProduct(id string, req UpdateProductRequest) (*model.Product, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Validate category if provided
	if req.CategoryID != nil {
		_, err := s.categoryRepo.FindByID(*req.CategoryID)
		if err != nil {
			return nil, errors.New("category not found")
		}
		product.CategoryID = *req.CategoryID
	}

	// Check SKU uniqueness if provided
	if req.SKU != nil && *req.SKU != product.SKU {
		existing, _ := s.productRepo.FindBySKU(*req.SKU)
		if existing != nil && existing.ID != product.ID {
			return nil, errors.New("SKU already exists")
		}
		product.SKU = *req.SKU
	}

	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Description != nil {
		product.Description = req.Description
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Stock != nil {
		product.Stock = *req.Stock
	}
	if req.Weight != nil {
		product.Weight = req.Weight
	}
	if req.Thumbnail != nil {
		product.Thumbnail = req.Thumbnail
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}
	if req.IsFeatured != nil {
		product.IsFeatured = *req.IsFeatured
	}

	if err := s.productRepo.Update(product); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return s.productRepo.FindByID(product.ID)
}

func (s *productService) DeleteProduct(id string) error {
	_, err := s.productRepo.FindByID(id)
	if err != nil {
		return errors.New("product not found")
	}

	return s.productRepo.Delete(id)
}

func (s *productService) AddProductImage(productID string, req AddProductImageRequest) (*model.ProductImage, error) {
	// Validate product exists
	_, err := s.productRepo.FindByID(productID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	sortOrder := 0
	if req.SortOrder != nil {
		sortOrder = *req.SortOrder
	}

	image := &model.ProductImage{
		ProductID: productID,
		ImageURL:  req.ImageURL,
		SortOrder: sortOrder,
	}

	if err := s.productRepo.CreateImage(image); err != nil {
		return nil, fmt.Errorf("failed to add image: %w", err)
	}

	return image, nil
}

func (s *productService) DeleteProductImage(imageID string) error {
	return s.productRepo.DeleteImage(imageID)
}
