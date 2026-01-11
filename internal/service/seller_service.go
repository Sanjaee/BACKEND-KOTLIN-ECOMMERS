package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"yourapp/internal/model"
	"yourapp/internal/repository"
)

type SellerService interface {
	CreateSeller(userID string, req CreateSellerRequest) (*model.Seller, error)
	GetSellerByID(sellerID string) (*model.Seller, error)
	GetSellerByUserID(userID string) (*model.Seller, error)
	UpdateSeller(userID string, req UpdateSellerRequest) (*model.Seller, error)
	DeleteSeller(userID string) error
}

type sellerService struct {
	sellerRepo repository.SellerRepository
	userRepo   repository.UserRepository
}

type CreateSellerRequest struct {
	ShopName       string  `json:"shop_name" binding:"required"`
	ShopDescription *string `json:"shop_description,omitempty"`
	ShopLogo       *string `json:"shop_logo,omitempty"`
	ShopBanner     *string `json:"shop_banner,omitempty"`
	ShopAddress    *string `json:"shop_address,omitempty"`
	ShopCity       *string `json:"shop_city,omitempty"`
	ShopProvince   *string `json:"shop_province,omitempty"`
	ShopPhone      *string `json:"shop_phone,omitempty"`
	ShopEmail      *string `json:"shop_email,omitempty"`
}

type UpdateSellerRequest struct {
	ShopName       *string `json:"shop_name,omitempty"`
	ShopDescription *string `json:"shop_description,omitempty"`
	ShopLogo       *string `json:"shop_logo,omitempty"`
	ShopBanner     *string `json:"shop_banner,omitempty"`
	ShopAddress    *string `json:"shop_address,omitempty"`
	ShopCity       *string `json:"shop_city,omitempty"`
	ShopProvince   *string `json:"shop_province,omitempty"`
	ShopPhone      *string `json:"shop_phone,omitempty"`
	ShopEmail      *string `json:"shop_email,omitempty"`
}

func NewSellerService(sellerRepo repository.SellerRepository, userRepo repository.UserRepository) SellerService {
	return &sellerService{
		sellerRepo: sellerRepo,
		userRepo:   userRepo,
	}
}

func (s *sellerService) CreateSeller(userID string, req CreateSellerRequest) (*model.Seller, error) {
	// Validasi user exists
	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Cek apakah user sudah punya toko (1 user 1 toko)
	existingSeller, _ := s.sellerRepo.FindByUserID(userID)
	if existingSeller != nil {
		return nil, errors.New("user already has a shop. One user can only have one shop")
	}

	// Generate slug dari shop_name
	shopSlug := generateSellerSlug(req.ShopName)

	// Validasi slug uniqueness
	existingBySlug, _ := s.sellerRepo.FindBySlug(shopSlug)
	if existingBySlug != nil {
		// Add timestamp or random string to make it unique
		shopSlug = shopSlug + "-" + strings.ToLower(generateUniqueSuffix())
	}

	seller := &model.Seller{
		UserID:         userID,
		ShopName:       req.ShopName,
		ShopSlug:       shopSlug,
		ShopDescription: req.ShopDescription,
		ShopLogo:       req.ShopLogo,
		ShopBanner:     req.ShopBanner,
		ShopAddress:    req.ShopAddress,
		ShopCity:       req.ShopCity,
		ShopProvince:   req.ShopProvince,
		ShopPhone:      req.ShopPhone,
		ShopEmail:      req.ShopEmail,
		IsActive:       true,
		IsVerified:     false,
		TotalProducts:  0,
		TotalSales:     0,
		RatingAverage:  0.00,
		TotalReviews:   0,
	}

	if err := s.sellerRepo.Create(seller); err != nil {
		// Check if error is due to duplicate shop_name
		if strings.Contains(err.Error(), "shop_name") || strings.Contains(err.Error(), "duplicate") {
			return nil, errors.New("shop name already exists")
		}
		return nil, fmt.Errorf("failed to create seller: %w", err)
	}

	return s.sellerRepo.FindByID(seller.ID)
}

func (s *sellerService) GetSellerByID(sellerID string) (*model.Seller, error) {
	seller, err := s.sellerRepo.FindByID(sellerID)
	if err != nil {
		return nil, errors.New("seller not found")
	}
	return seller, nil
}

func (s *sellerService) GetSellerByUserID(userID string) (*model.Seller, error) {
	seller, err := s.sellerRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("seller not found")
	}
	return seller, nil
}

func (s *sellerService) UpdateSeller(userID string, req UpdateSellerRequest) (*model.Seller, error) {
	// Get seller by user_id (hanya owner yang bisa update)
	seller, err := s.sellerRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("seller not found")
	}

	// Update shop_name dan generate slug baru jika shop_name berubah
	if req.ShopName != nil && *req.ShopName != seller.ShopName {
		// Generate slug dari shop_name baru
		newSlug := generateSellerSlug(*req.ShopName)
		// Validasi slug uniqueness
		existingBySlug, _ := s.sellerRepo.FindBySlug(newSlug)
		if existingBySlug != nil && existingBySlug.ID != seller.ID {
			return nil, errors.New("shop name already exists")
		}
		seller.ShopName = *req.ShopName
		// Slug akan diupdate otomatis oleh BeforeUpdate hook di model
	}

	if req.ShopDescription != nil {
		seller.ShopDescription = req.ShopDescription
	}
	if req.ShopLogo != nil {
		seller.ShopLogo = req.ShopLogo
	}
	if req.ShopBanner != nil {
		seller.ShopBanner = req.ShopBanner
	}
	if req.ShopAddress != nil {
		seller.ShopAddress = req.ShopAddress
	}
	if req.ShopCity != nil {
		seller.ShopCity = req.ShopCity
	}
	if req.ShopProvince != nil {
		seller.ShopProvince = req.ShopProvince
	}
	if req.ShopPhone != nil {
		seller.ShopPhone = req.ShopPhone
	}
	if req.ShopEmail != nil {
		seller.ShopEmail = req.ShopEmail
	}

	if err := s.sellerRepo.Update(seller); err != nil {
		// Check if error is due to duplicate shop_name
		if strings.Contains(err.Error(), "shop_name") || strings.Contains(err.Error(), "duplicate") {
			return nil, errors.New("shop name already exists")
		}
		return nil, fmt.Errorf("failed to update seller: %w", err)
	}

	return s.sellerRepo.FindByID(seller.ID)
}

func (s *sellerService) DeleteSeller(userID string) error {
	// Get seller by user_id (hanya owner yang bisa delete)
	seller, err := s.sellerRepo.FindByUserID(userID)
	if err != nil {
		return errors.New("seller not found")
	}

	// Soft delete
	return s.sellerRepo.Delete(seller.ID)
}

// generateSellerSlug generates a URL-friendly slug from a string
func generateSellerSlug(text string) string {
	slug := strings.ToLower(text)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")
	// Remove special characters (keep only alphanumeric and hyphens)
	var result strings.Builder
	for _, char := range slug {
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-' {
			result.WriteRune(char)
		}
	}
	return result.String()
}

// generateUniqueSuffix generates a short unique suffix
func generateUniqueSuffix() string {
	return fmt.Sprintf("%d", time.Now().Unix()%10000)
}
