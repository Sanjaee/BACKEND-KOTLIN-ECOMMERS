package app

import (
	"net/http"

	"yourapp/internal/service"
	"yourapp/internal/util"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryService service.CategoryService
}

func NewCategoryHandler(categoryService service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// CreateCategory handles category creation
// POST /api/v1/categories
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req service.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, err.Error())
		return
	}

	category, err := h.categoryService.CreateCategory(req)
	if err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	util.SuccessResponse(c, http.StatusCreated, "Category created successfully", category)
}

// GetCategory handles getting category by ID
// GET /api/v1/categories/:id
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		util.BadRequest(c, "Category ID is required")
		return
	}

	category, err := h.categoryService.GetCategoryByID(id)
	if err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	util.SuccessResponse(c, http.StatusOK, "Category retrieved successfully", category)
}

// GetCategoryBySlug handles getting category by slug
// GET /api/v1/categories/slug/:slug
func (h *CategoryHandler) GetCategoryBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		util.BadRequest(c, "Category slug is required")
		return
	}

	category, err := h.categoryService.GetCategoryBySlug(slug)
	if err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	util.SuccessResponse(c, http.StatusOK, "Category retrieved successfully", category)
}

// GetCategories handles getting list of categories
// GET /api/v1/categories
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	activeOnly := c.Query("active_only") == "true"

	categories, err := h.categoryService.GetCategories(activeOnly)
	if err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	util.SuccessResponse(c, http.StatusOK, "Categories retrieved successfully", categories)
}

// UpdateCategory handles category update
// PUT /api/v1/categories/:id
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		util.BadRequest(c, "Category ID is required")
		return
	}

	var req service.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, err.Error())
		return
	}

	category, err := h.categoryService.UpdateCategory(id, req)
	if err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	util.SuccessResponse(c, http.StatusOK, "Category updated successfully", category)
}

// DeleteCategory handles category deletion
// DELETE /api/v1/categories/:id
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		util.BadRequest(c, "Category ID is required")
		return
	}

	if err := h.categoryService.DeleteCategory(id); err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	util.SuccessResponse(c, http.StatusOK, "Category deleted successfully", nil)
}
