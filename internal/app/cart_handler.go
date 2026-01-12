package app

import (
	"net/http"
	"yourapp/internal/service"
	"yourapp/internal/util"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	cartService service.CartService
}

func NewCartHandler(cartService service.CartService) *CartHandler {
	return &CartHandler{
		cartService: cartService,
	}
}

// GetCart handles getting user's cart
// GET /api/v1/carts
func (h *CartHandler) GetCart(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		util.Unauthorized(c, "User not authenticated")
		return
	}

	cart, err := h.cartService.GetCart(userID.(string))
	if err != nil {
		util.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	util.SuccessResponse(c, http.StatusOK, "Cart retrieved successfully", cart)
}

// AddItemToCart handles adding item to cart
// POST /api/v1/carts/items
func (h *CartHandler) AddItemToCart(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		util.Unauthorized(c, "User not authenticated")
		return
	}

	var req service.AddCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, err.Error())
		return
	}

	cartItem, err := h.cartService.AddItemToCart(userID.(string), &req)
	if err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	util.SuccessResponse(c, http.StatusCreated, "Item added to cart successfully", cartItem)
}

// UpdateCartItem handles updating cart item quantity
// PUT /api/v1/carts/items/:id
func (h *CartHandler) UpdateCartItem(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		util.Unauthorized(c, "User not authenticated")
		return
	}

	cartItemID := c.Param("id")
	if cartItemID == "" {
		util.BadRequest(c, "Cart item ID is required")
		return
	}

	var req service.UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, err.Error())
		return
	}

	cartItem, err := h.cartService.UpdateCartItem(userID.(string), cartItemID, &req)
	if err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	util.SuccessResponse(c, http.StatusOK, "Cart item updated successfully", cartItem)
}

// RemoveCartItem handles removing item from cart
// DELETE /api/v1/carts/items/:id
func (h *CartHandler) RemoveCartItem(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		util.Unauthorized(c, "User not authenticated")
		return
	}

	cartItemID := c.Param("id")
	if cartItemID == "" {
		util.BadRequest(c, "Cart item ID is required")
		return
	}

	err := h.cartService.RemoveCartItem(userID.(string), cartItemID)
	if err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	util.SuccessResponse(c, http.StatusOK, "Item removed from cart successfully", nil)
}

// ClearCart handles clearing all items from cart
// DELETE /api/v1/carts
func (h *CartHandler) ClearCart(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		util.Unauthorized(c, "User not authenticated")
		return
	}

	err := h.cartService.ClearCart(userID.(string))
	if err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	util.SuccessResponse(c, http.StatusOK, "Cart cleared successfully", nil)
}

// GetCartItems handles getting all cart items
// GET /api/v1/carts/items
func (h *CartHandler) GetCartItems(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		util.Unauthorized(c, "User not authenticated")
		return
	}

	cartItems, err := h.cartService.GetCartItems(userID.(string))
	if err != nil {
		util.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	util.SuccessResponse(c, http.StatusOK, "Cart items retrieved successfully", cartItems)
}
