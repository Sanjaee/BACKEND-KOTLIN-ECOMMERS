package app

import (
	"net/http"

	"yourapp/internal/service"
	"yourapp/internal/util"

	"github.com/gin-gonic/gin"
)

type SellerHandler struct {
	sellerService service.SellerService
}

func NewSellerHandler(sellerService service.SellerService) *SellerHandler {
	return &SellerHandler{
		sellerService: sellerService,
	}
}

// CreateSeller handles shop creation
// POST /api/v1/sellers
func (h *SellerHandler) CreateSeller(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		util.Unauthorized(c, "User not authenticated")
		return
	}

	var req service.CreateSellerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, err.Error())
		return
	}

	seller, err := h.sellerService.CreateSeller(userID.(string), req)
	if err != nil {
		if err.Error() == "user already has a shop. One user can only have one shop" {
			util.ErrorResponse(c, http.StatusConflict, err.Error(), nil)
			return
		}
		util.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	util.SuccessResponse(c, http.StatusCreated, "Shop created successfully", seller)
}

// GetSeller handles getting shop by ID
// GET /api/v1/sellers/:id
func (h *SellerHandler) GetSeller(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		util.BadRequest(c, "Seller ID is required")
		return
	}

	seller, err := h.sellerService.GetSellerByID(id)
	if err != nil {
		util.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	util.SuccessResponse(c, http.StatusOK, "Shop retrieved successfully", seller)
}

// GetMySeller handles getting current user's shop
// GET /api/v1/sellers/me
func (h *SellerHandler) GetMySeller(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		util.Unauthorized(c, "User not authenticated")
		return
	}

	seller, err := h.sellerService.GetSellerByUserID(userID.(string))
	if err != nil {
		util.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	util.SuccessResponse(c, http.StatusOK, "Shop retrieved successfully", seller)
}

// UpdateSeller handles shop update
// PUT /api/v1/sellers
func (h *SellerHandler) UpdateSeller(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		util.Unauthorized(c, "User not authenticated")
		return
	}

	var req service.UpdateSellerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, err.Error())
		return
	}

	seller, err := h.sellerService.UpdateSeller(userID.(string), req)
	if err != nil {
		if err.Error() == "shop name already exists" {
			util.ErrorResponse(c, http.StatusConflict, err.Error(), nil)
			return
		}
		util.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	util.SuccessResponse(c, http.StatusOK, "Shop updated successfully", seller)
}

// DeleteSeller handles shop deletion
// DELETE /api/v1/sellers
func (h *SellerHandler) DeleteSeller(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		util.Unauthorized(c, "User not authenticated")
		return
	}

	err := h.sellerService.DeleteSeller(userID.(string))
	if err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	util.SuccessResponse(c, http.StatusOK, "Shop deleted successfully", nil)
}
