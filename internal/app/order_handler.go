package app

import (
	"net/http"
	"strconv"
	"yourapp/internal/service"
	"yourapp/internal/util"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// CreateOrder handles order creation from checkout
// POST /api/v1/orders
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		util.Unauthorized(c, "User not authenticated")
		return
	}

	var req service.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, err.Error())
		return
	}

	order, err := h.orderService.CreateOrder(userID.(string), &req)
	if err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	util.SuccessResponse(c, http.StatusCreated, "Order created successfully", order)
}

// GetOrder handles getting order by ID
// GET /api/v1/orders/:id
func (h *OrderHandler) GetOrder(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		util.Unauthorized(c, "User not authenticated")
		return
	}

	id := c.Param("id")
	if id == "" {
		util.BadRequest(c, "Order ID is required")
		return
	}

	order, err := h.orderService.GetOrderByID(id, userID.(string))
	if err != nil {
		util.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	util.SuccessResponse(c, http.StatusOK, "Order retrieved successfully", order)
}

// GetOrders handles getting list of orders for authenticated user
// GET /api/v1/orders
func (h *OrderHandler) GetOrders(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		util.Unauthorized(c, "User not authenticated")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	orders, total, err := h.orderService.GetOrdersByUserID(userID.(string), page, limit)
	if err != nil {
		util.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	util.SuccessResponse(c, http.StatusOK, "Orders retrieved successfully", gin.H{
		"orders": orders,
		"total":  total,
		"page":   page,
		"limit":  limit,
	})
}
