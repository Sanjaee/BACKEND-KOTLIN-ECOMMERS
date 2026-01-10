package service

import (
	"errors"
	"fmt"
	"strings"

	"yourapp/internal/model"
	"yourapp/internal/repository"
)

type CategoryService interface {
	CreateCategory(req CreateCategoryRequest) (*model.Category, error)
	GetCategoryByID(id string) (*model.Category, error)
	GetCategoryBySlug(slug string) (*model.Category, error)
	GetCategories(activeOnly bool) ([]model.Category, error)
	UpdateCategory(id string, req UpdateCategoryRequest) (*model.Category, error)
	DeleteCategory(id string) error
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
}

type CreateCategoryRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description,omitempty"`
	Slug        string  `json:"slug" binding:"required"`
	ImageURL    *string `json:"image_url,omitempty"`
	ParentID    *string `json:"parent_id,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

type UpdateCategoryRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Slug        *string `json:"slug,omitempty"`
	ImageURL    *string `json:"image_url,omitempty"`
	ParentID    *string `json:"parent_id,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *categoryService) CreateCategory(req CreateCategoryRequest) (*model.Category, error) {
	// Generate slug from name if not provided
	slug := req.Slug
	if slug == "" {
		slug = generateSlug(req.Name)
	}

	// Validate slug uniqueness
	existing, _ := s.categoryRepo.FindBySlug(slug)
	if existing != nil {
		return nil, errors.New("slug already exists")
	}

	// Validate parent category if provided
	if req.ParentID != nil && *req.ParentID != "" {
		parent, err := s.categoryRepo.FindByID(*req.ParentID)
		if err != nil {
			return nil, errors.New("parent category not found")
		}
		// Prevent circular reference (parent can't be itself)
		if parent.ID == *req.ParentID {
			return nil, errors.New("category cannot be its own parent")
		}
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	category := &model.Category{
		Name:        req.Name,
		Description: req.Description,
		Slug:        slug,
		ImageURL:    req.ImageURL,
		ParentID:    req.ParentID,
		IsActive:    isActive,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return s.categoryRepo.FindByID(category.ID)
}

func (s *categoryService) GetCategoryByID(id string) (*model.Category, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("category not found")
	}
	return category, nil
}

func (s *categoryService) GetCategoryBySlug(slug string) (*model.Category, error) {
	category, err := s.categoryRepo.FindBySlug(slug)
	if err != nil {
		return nil, errors.New("category not found")
	}
	return category, nil
}

func (s *categoryService) GetCategories(activeOnly bool) ([]model.Category, error) {
	categories, err := s.categoryRepo.FindAll(activeOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}
	return categories, nil
}

func (s *categoryService) UpdateCategory(id string, req UpdateCategoryRequest) (*model.Category, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("category not found")
	}

	// Validate slug uniqueness if provided
	if req.Slug != nil && *req.Slug != category.Slug {
		existing, _ := s.categoryRepo.FindBySlug(*req.Slug)
		if existing != nil && existing.ID != category.ID {
			return nil, errors.New("slug already exists")
		}
		category.Slug = *req.Slug
	}

	// Validate parent category if provided
	if req.ParentID != nil {
		if *req.ParentID == "" {
			// Remove parent (set to null)
			category.ParentID = nil
		} else {
			// Check if parent exists
			parent, err := s.categoryRepo.FindByID(*req.ParentID)
			if err != nil {
				return nil, errors.New("parent category not found")
			}
			// Prevent circular reference (can't set parent to itself or its children)
			if parent.ID == category.ID {
				return nil, errors.New("category cannot be its own parent")
			}
			category.ParentID = req.ParentID
		}
	}

	if req.Name != nil {
		category.Name = *req.Name
	}
	if req.Description != nil {
		category.Description = req.Description
	}
	if req.ImageURL != nil {
		category.ImageURL = req.ImageURL
	}
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}

	if err := s.categoryRepo.Update(category); err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return s.categoryRepo.FindByID(category.ID)
}

func (s *categoryService) DeleteCategory(id string) error {
	_, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return errors.New("category not found")
	}

	// Check if category has children (optional validation)
	// This can be implemented if needed to prevent deletion of parent categories

	return s.categoryRepo.Delete(id)
}

// generateSlug generates a URL-friendly slug from a string
func generateSlug(text string) string {
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
