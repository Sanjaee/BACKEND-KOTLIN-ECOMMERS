package app

import (
	"log"
	"net/http"
	"yourapp/internal/model"
	"yourapp/internal/service"
	"yourapp/internal/util"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentService service.PaymentService
}

func NewPaymentHandler(paymentService service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

// CreatePayment handles payment creation for an order
// POST /api/v1/payments
func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var req struct {
		OrderID       string  `json:"order_id" binding:"required"`
		PaymentMethod string  `json:"payment_method" binding:"required"`
		Bank          *string `json:"bank,omitempty"` // bca, bni, mandiri, etc (for bank_transfer)
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, err.Error())
		return
	}

	// Validate payment method
	paymentMethod := model.PaymentMethod(req.PaymentMethod)
	validMethods := map[model.PaymentMethod]bool{
		model.PaymentMethodBankTransfer: true,
		model.PaymentMethodGopay:        true,
		model.PaymentMethodCreditCard:   true,
		model.PaymentMethodQRIS:         true,
		model.PaymentMethodAlfamart:     true,
	}
	if !validMethods[paymentMethod] {
		util.BadRequest(c, "Invalid payment method")
		return
	}

	payment, err := h.paymentService.CreatePayment(req.OrderID, paymentMethod, req.Bank)
	if err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	util.SuccessResponse(c, http.StatusCreated, "Payment created successfully", payment)
}

// GetPayment handles getting payment by ID
// GET /api/v1/payments/:id
func (h *PaymentHandler) GetPayment(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		util.BadRequest(c, "Payment ID is required")
		return
	}

	payment, err := h.paymentService.GetPaymentByID(id)
	if err != nil {
		util.ErrorResponse(c, http.StatusNotFound, "Payment not found", nil)
		return
	}

	util.SuccessResponse(c, http.StatusOK, "Payment retrieved successfully", payment)
}

// GetPaymentByOrder handles getting payment by order ID
// GET /api/v1/payments/order/:order_id
func (h *PaymentHandler) GetPaymentByOrder(c *gin.Context) {
	orderID := c.Param("order_id")
	if orderID == "" {
		util.BadRequest(c, "Order ID is required")
		return
	}

	payment, err := h.paymentService.GetPaymentByOrderID(orderID)
	if err != nil {
		util.ErrorResponse(c, http.StatusNotFound, "Payment not found", nil)
		return
	}

	util.SuccessResponse(c, http.StatusOK, "Payment retrieved successfully", payment)
}

// CheckPaymentStatus handles checking payment status
// GET /api/v1/payments/:id/status
// This endpoint always checks latest status from Midtrans API if payment is still pending
func (h *PaymentHandler) CheckPaymentStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		util.BadRequest(c, "Payment ID is required")
		return
	}

	// Force check from Midtrans API if payment is pending
	payment, err := h.paymentService.CheckPaymentStatus(id)
	if err != nil {
		util.ErrorResponse(c, http.StatusNotFound, "Payment not found", nil)
		return
	}

	util.SuccessResponse(c, http.StatusOK, "Payment status retrieved successfully", payment)
}

// MidtransCallback handles Midtrans payment callback
// POST /api/v1/payments/midtrans/callback
// This is a PUBLIC endpoint - Midtrans will POST webhook notifications here
// Note: In production, you should verify the signature for security
func (h *PaymentHandler) MidtransCallback(c *gin.Context) {
	var notification map[string]interface{}
	if err := c.ShouldBindJSON(&notification); err != nil {
		log.Printf("❌ Invalid Midtrans callback JSON: %v", err)
		util.BadRequest(c, "Invalid notification format")
		return
	}

	// Log raw notification for debugging
	log.Printf("📥 Received Midtrans callback: %+v", notification)

	// Process callback asynchronously to respond quickly to Midtrans
	// Midtrans expects fast response (< 10 seconds)
	go func() {
		if err := h.paymentService.HandleMidtransCallback(notification); err != nil {
			log.Printf("❌ Failed to process Midtrans callback: %v", err)
			// Note: We still return 200 OK to Midtrans even if processing fails
			// This prevents Midtrans from retrying immediately
			// Error will be logged and can be retried manually or via background job
		}
	}()

	// Respond immediately to Midtrans (within 10 seconds requirement)
	// Status will be updated in background
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Callback received",
	})
}
