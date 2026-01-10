package service

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	"yourapp/internal/config"
	"yourapp/internal/model"
	"yourapp/internal/repository"
)

type PaymentService interface {
	CreatePayment(orderID string, paymentMethod model.PaymentMethod, bankType *string) (*model.Payment, error)
	GetPaymentByID(paymentID string) (*model.Payment, error)
	GetPaymentByOrderID(orderID string) (*model.Payment, error)
	HandleMidtransCallback(notification map[string]interface{}) error
	CheckPaymentStatus(paymentID string) (*model.Payment, error)
	CheckPaymentStatusFromMidtrans(orderID string) error
	UpdatePaymentStatus(orderID string, status string, transactionID string, vaNumber string, bankType string, qrCodeURL string, expiryTime *time.Time, midtransResponse string) error
}

type paymentService struct {
	paymentRepo repository.PaymentRepository
	orderRepo   repository.OrderRepository
	cfg         *config.Config
}

// Midtrans API request/response structures
type MidtransChargeRequest struct {
	PaymentType        string                     `json:"payment_type"`
	TransactionDetails MidtransTransactionDetails `json:"transaction_details"`
	CustomerDetails    MidtransCustomerDetails    `json:"customer_details"`
	ItemDetails        []MidtransItemDetail       `json:"item_details"`
	BankTransfer       *MidtransBankTransfer      `json:"bank_transfer,omitempty"`
	Gopay              *MidtransGopay             `json:"gopay,omitempty"`
	CreditCard         *MidtransCreditCard        `json:"credit_card,omitempty"`
}

type MidtransTransactionDetails struct {
	OrderID     string `json:"order_id"`
	GrossAmount int    `json:"gross_amount"`
}

type MidtransCustomerDetails struct {
	FirstName string `json:"first_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone,omitempty"`
}

type MidtransItemDetail struct {
	ID       string `json:"id"`
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
	Name     string `json:"name"`
	Category string `json:"category,omitempty"`
}

type MidtransBankTransfer struct {
	Bank string `json:"bank"`
}

type MidtransGopay struct {
	EnableCallback bool   `json:"enable_callback"`
	CallbackURL    string `json:"callback_url"`
}

type MidtransCreditCard struct {
	Secure         bool `json:"secure"`
	Authentication bool `json:"authentication"`
}

type MidtransChargeResponse struct {
	TransactionID     string             `json:"transaction_id"`
	OrderID           string             `json:"order_id"`
	GrossAmount       string             `json:"gross_amount"`
	PaymentType       string             `json:"payment_type"`
	TransactionTime   string             `json:"transaction_time"`
	TransactionStatus string             `json:"transaction_status"`
	FraudStatus       string             `json:"fraud_status"`
	StatusMessage     string             `json:"status_message"`
	VANumbers         []MidtransVANumber `json:"va_numbers,omitempty"`
	Actions           []MidtransAction   `json:"actions,omitempty"`
	ExpiryTime        string             `json:"expiry_time,omitempty"`
	QRCodeURL         string             `json:"qr_code_url,omitempty"`
}

type MidtransVANumber struct {
	Bank     string `json:"bank"`
	VANumber string `json:"va_number"`
}

type MidtransAction struct {
	Name   string `json:"name"`
	Method string `json:"method"`
	URL    string `json:"url"`
}

func NewPaymentService(
	paymentRepo repository.PaymentRepository,
	orderRepo repository.OrderRepository,
	cfg *config.Config,
) PaymentService {
	return &paymentService{
		paymentRepo: paymentRepo,
		orderRepo:   orderRepo,
		cfg:         cfg,
	}
}

// mapMidtransStatusToPaymentStatus maps Midtrans status to PaymentStatus
func mapMidtransStatusToPaymentStatus(status string) model.PaymentStatus {
	switch status {
	case "pending":
		return model.PaymentStatusPending
	case "settlement", "capture":
		return model.PaymentStatusSuccess
	case "deny":
		return model.PaymentStatusFailed
	case "cancel":
		return model.PaymentStatusCancelled
	case "expire":
		return model.PaymentStatusExpired
	default:
		return model.PaymentStatusPending
	}
}

// getMidtransBaseURL returns Midtrans API base URL based on environment
func (s *paymentService) getMidtransBaseURL() string {
	if s.cfg.MidtransServerKey != "" {
		// Check if it's production key (starts with Mid-server) or sandbox (starts with SB-Mid-server)
		if strings.HasPrefix(s.cfg.MidtransServerKey, "Mid-server") {
			return "https://api.midtrans.com/v2"
		}
	}
	return "https://api.sandbox.midtrans.com/v2"
}

// getAuthHeader returns base64 encoded authorization header
func (s *paymentService) getAuthHeader() string {
	auth := base64.StdEncoding.EncodeToString([]byte(s.cfg.MidtransServerKey + ":"))
	return "Basic " + auth
}

func (s *paymentService) CreatePayment(orderID string, paymentMethod model.PaymentMethod, bankType *string) (*model.Payment, error) {
	// Get order with preloaded data
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	// Check if payment already exists
	existingPayment, _ := s.paymentRepo.FindByOrderID(orderID)
	if existingPayment != nil {
		return existingPayment, nil
	}

	// Create payment record first
	payment := &model.Payment{
		OrderID:       order.OrderNumber,
		OrderUUID:     order.ID,
		Amount:        order.TotalAmount,
		TotalAmount:   order.TotalAmount,
		Status:        model.PaymentStatusPending,
		PaymentMethod: paymentMethod,
		PaymentType:   "midtrans",
	}

	if err := s.paymentRepo.Create(payment); err != nil {
		log.Printf("❌ Failed to create payment: %v", err)
		return nil, fmt.Errorf("failed to create payment: %v", err)
	}

	// If Midtrans is not configured, return payment without transaction
	if s.cfg.MidtransServerKey == "" {
		log.Printf("⚠️  Midtrans not configured, returning payment without transaction")
		return payment, nil
	}

	// Prepare customer details
	customerPhone := ""
	if order.User.Phone != nil {
		customerPhone = *order.User.Phone
	}

	customerDetails := MidtransCustomerDetails{
		FirstName: order.User.FullName,
		Email:     order.User.Email,
		Phone:     customerPhone,
	}

	// Prepare item details
	var itemDetails []MidtransItemDetail
	for _, item := range order.OrderItems {
		itemDetails = append(itemDetails, MidtransItemDetail{
			ID:       item.ProductID,
			Price:    item.Price,
			Quantity: item.Quantity,
			Name:     item.ProductName,
			Category: "product",
		})
	}

	// Add shipping cost, insurance, warranty as separate items
	if order.ShippingCost > 0 {
		itemDetails = append(itemDetails, MidtransItemDetail{
			ID:       "shipping",
			Price:    order.ShippingCost,
			Quantity: 1,
			Name:     "Shipping Cost",
			Category: "shipping",
		})
	}

	if order.InsuranceCost > 0 {
		itemDetails = append(itemDetails, MidtransItemDetail{
			ID:       "insurance",
			Price:    order.InsuranceCost,
			Quantity: 1,
			Name:     "Shipping Insurance",
			Category: "insurance",
		})
	}

	if order.WarrantyCost > 0 {
		itemDetails = append(itemDetails, MidtransItemDetail{
			ID:       "warranty",
			Price:    order.WarrantyCost,
			Quantity: 1,
			Name:     "Warranty Protection",
			Category: "warranty",
		})
	}

	if order.ServiceFee > 0 {
		itemDetails = append(itemDetails, MidtransItemDetail{
			ID:       "service_fee",
			Price:    order.ServiceFee,
			Quantity: 1,
			Name:     "Service Fee",
			Category: "fee",
		})
	}

	// Add discount as negative item (Midtrans requires item_details sum to equal gross_amount)
	if order.TotalDiscount > 0 {
		itemDetails = append(itemDetails, MidtransItemDetail{
			ID:       "discount",
			Price:    -order.TotalDiscount, // Negative price for discount
			Quantity: 1,
			Name:     "Discount",
			Category: "discount",
		})
	}

	// Add bonus as negative item (cashback/promotion)
	if order.Bonus > 0 {
		itemDetails = append(itemDetails, MidtransItemDetail{
			ID:       "bonus",
			Price:    -order.Bonus, // Negative price for bonus/cashback
			Quantity: 1,
			Name:     "Bonus Cashback",
			Category: "bonus",
		})
	}

	// Calculate gross_amount as sum of all item_details to ensure it matches Midtrans requirement
	// This ensures: gross_amount = sum(item_details[i].price * item_details[i].quantity)
	var grossAmount int
	for _, item := range itemDetails {
		grossAmount += item.Price * item.Quantity
	}

	// Verify that calculated gross_amount matches order.TotalAmount (they should be equal)
	if grossAmount != order.TotalAmount {
		log.Printf("⚠️  Warning: Calculated gross_amount (%d) does not match order.TotalAmount (%d). Using calculated value.", grossAmount, order.TotalAmount)
	}

	// Prepare charge request
	chargeData := MidtransChargeRequest{
		PaymentType: string(paymentMethod),
		TransactionDetails: MidtransTransactionDetails{
			OrderID:     order.OrderNumber,
			GrossAmount: grossAmount, // Use calculated sum to ensure it matches item_details
		},
		CustomerDetails: customerDetails,
		ItemDetails:     itemDetails,
	}

	// Add payment method specific config
	frontendURL := s.cfg.ClientURL
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	switch paymentMethod {
	case model.PaymentMethodBankTransfer:
		bank := "bca" // Default to BCA
		if bankType != nil && *bankType != "" {
			bank = strings.ToLower(*bankType)
		}
		chargeData.BankTransfer = &MidtransBankTransfer{Bank: bank}

	case model.PaymentMethodGopay:
		chargeData.Gopay = &MidtransGopay{
			EnableCallback: true,
			CallbackURL:    fmt.Sprintf("%s/api/v1/payments/midtrans/callback", frontendURL),
		}

	case model.PaymentMethodQRIS:
		// QRIS uses qris payment type
		chargeData.PaymentType = "qris"
		chargeData.Gopay = &MidtransGopay{
			EnableCallback: true,
			CallbackURL:    fmt.Sprintf("%s/api/v1/payments/midtrans/callback", frontendURL),
		}

	case model.PaymentMethodCreditCard:
		chargeData.CreditCard = &MidtransCreditCard{
			Secure:         true,
			Authentication: true,
		}

	case model.PaymentMethodAlfamart:
		// Alfamart uses cstore payment type
		chargeData.PaymentType = "cstore"
	}

	// Charge to Midtrans
	chargeJSON, err := json.Marshal(chargeData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal charge data: %v", err)
	}

	baseURL := s.getMidtransBaseURL()
	authHeader := s.getAuthHeader()

	// Make HTTP request to Midtrans
	reqHTTP, err := http.NewRequest("POST", baseURL+"/charge", bytes.NewBuffer(chargeJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	reqHTTP.Header.Set("Authorization", authHeader)
	reqHTTP.Header.Set("Content-Type", "application/json")
	reqHTTP.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(reqHTTP)
	if err != nil {
		log.Printf("⚠️  Failed to charge Midtrans: %v", err)
		return payment, nil // Return payment even if Midtrans fails
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("⚠️  Failed to read Midtrans response: %v", err)
		return payment, nil
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		log.Printf("⚠️  Midtrans API returned status %d: %s", resp.StatusCode, string(body))
		// Store error response but don't fail
		errorResp := string(body)
		payment.MidtransResponse = &errorResp
		s.paymentRepo.Update(payment)
		return payment, nil
	}

	var midtransResp MidtransChargeResponse
	if err := json.Unmarshal(body, &midtransResp); err != nil {
		log.Printf("⚠️  Failed to parse Midtrans response: %v", err)
		return payment, nil
	}

	// Extract payment details from response
	var vaNumber, bankTypeStr, qrCodeURL string
	if len(midtransResp.VANumbers) > 0 {
		vaNumber = midtransResp.VANumbers[0].VANumber
		bankTypeStr = midtransResp.VANumbers[0].Bank
	}

	// Extract QR code URL from actions (for Gopay/QRIS)
	for _, action := range midtransResp.Actions {
		if action.Name == "generate-qr-code" || action.Name == "generate-qr-code-v2" || action.Name == "qr-code" {
			qrCodeURL = action.URL
			break
		}
	}
	// If not found by name, try by method GET
	if qrCodeURL == "" {
		for _, action := range midtransResp.Actions {
			if action.Method == "GET" && action.URL != "" && strings.Contains(strings.ToLower(action.URL), "qr") {
				qrCodeURL = action.URL
				break
			}
		}
	}

	// Use QRCodeURL directly from response if available
	if qrCodeURL == "" && midtransResp.QRCodeURL != "" {
		qrCodeURL = midtransResp.QRCodeURL
	}

	// Parse expiry time
	var expiryTime *time.Time
	if midtransResp.ExpiryTime != "" {
		formats := []string{
			time.RFC3339,
			"2006-01-02 15:04:05",
			"2006-01-02T15:04:05",
		}
		for _, format := range formats {
			exp, err := time.Parse(format, midtransResp.ExpiryTime)
			if err == nil {
				expiryTime = &exp
				break
			}
		}
	}

	// Update payment with Midtrans response
	updateData := map[string]interface{}{
		"midtrans_transaction_id": midtransResp.TransactionID,
		"status":                  mapMidtransStatusToPaymentStatus(midtransResp.TransactionStatus),
		"fraud_status":            midtransResp.FraudStatus,
		"midtrans_response":       string(body),
		"va_number":               vaNumber,
		"bank_type":               bankTypeStr,
		"qr_code_url":             qrCodeURL,
		"expiry_time":             expiryTime,
		"updated_at":              time.Now(),
	}

	// Update payment using repository
	if err := s.updatePaymentFields(payment.ID, updateData); err != nil {
		log.Printf("⚠️  Failed to update payment: %v", err)
	}

	// Reload payment with updated data
	updatedPayment, err := s.paymentRepo.FindByID(payment.ID)
	if err != nil {
		return payment, nil
	}

	return updatedPayment, nil
}

// updatePaymentFields updates payment fields using repository
func (s *paymentService) updatePaymentFields(paymentID string, updateData map[string]interface{}) error {
	payment, err := s.paymentRepo.FindByID(paymentID)
	if err != nil {
		return err
	}

	// Update fields manually since we're using map[string]interface{}
	if transactionID, ok := updateData["midtrans_transaction_id"].(string); ok {
		payment.MidtransTransactionID = &transactionID
	}
	if status, ok := updateData["status"].(model.PaymentStatus); ok {
		payment.Status = status
	}
	if fraudStatus, ok := updateData["fraud_status"].(string); ok && fraudStatus != "" {
		payment.FraudStatus = &fraudStatus
	}
	if midtransResponse, ok := updateData["midtrans_response"].(string); ok {
		payment.MidtransResponse = &midtransResponse
	}
	if vaNumber, ok := updateData["va_number"].(string); ok && vaNumber != "" {
		payment.VANumber = &vaNumber
	}
	if bankType, ok := updateData["bank_type"].(string); ok && bankType != "" {
		payment.BankType = &bankType
	}
	if qrCodeURL, ok := updateData["qr_code_url"].(string); ok && qrCodeURL != "" {
		payment.QRCodeURL = &qrCodeURL
	}
	if expiryTime, ok := updateData["expiry_time"].(*time.Time); ok && expiryTime != nil {
		payment.ExpiryTime = expiryTime
	}

	return s.paymentRepo.Update(payment)
}

func (s *paymentService) GetPaymentByID(paymentID string) (*model.Payment, error) {
	return s.paymentRepo.FindByID(paymentID)
}

func (s *paymentService) GetPaymentByOrderID(orderID string) (*model.Payment, error) {
	return s.paymentRepo.FindByOrderID(orderID)
}

func (s *paymentService) HandleMidtransCallback(notification map[string]interface{}) error {
	orderID, ok := notification["order_id"].(string)
	if !ok {
		return errors.New("invalid notification: missing order_id")
	}

	transactionID, ok := notification["transaction_id"].(string)
	if !ok {
		return errors.New("invalid notification: missing transaction_id")
	}

	transactionStatus, _ := notification["transaction_status"].(string)

	var vaNumber, bankType, qrCodeURL string

	// Extract VA numbers
	if vaNumbers, ok := notification["va_numbers"].([]interface{}); ok && len(vaNumbers) > 0 {
		if vaNum, ok := vaNumbers[0].(map[string]interface{}); ok {
			vaNumber, _ = vaNum["va_number"].(string)
			bankType, _ = vaNum["bank"].(string)
		}
	}

	// Extract QR code URL
	if qrCode, ok := notification["qr_code_url"].(string); ok {
		qrCodeURL = qrCode
	} else if actions, ok := notification["actions"].([]interface{}); ok && len(actions) > 0 {
		for _, action := range actions {
			if act, ok := action.(map[string]interface{}); ok {
				name, _ := act["name"].(string)
				url, _ := act["url"].(string)
				if (name == "generate-qr-code" || name == "generate-qr-code-v2" || name == "qr-code") && url != "" {
					qrCodeURL = url
					break
				}
			}
		}
		// If not found by name, try by method GET
		if qrCodeURL == "" {
			for _, action := range actions {
				if act, ok := action.(map[string]interface{}); ok {
					method, _ := act["method"].(string)
					url, _ := act["url"].(string)
					if method == "GET" && url != "" && strings.Contains(strings.ToLower(url), "qr") {
						qrCodeURL = url
						break
					}
				}
			}
		}
	}

	var expiryTime *time.Time
	if expiry, ok := notification["expiry_time"].(string); ok && expiry != "" {
		formats := []string{
			time.RFC3339,
			"2006-01-02 15:04:05",
			"2006-01-02T15:04:05",
		}
		for _, format := range formats {
			exp, err := time.Parse(format, expiry)
			if err == nil {
				expiryTime = &exp
				break
			}
		}
	}

	webhookJSON, _ := json.Marshal(notification)

	// Update payment status with fraud status included in midtransResponse
	return s.UpdatePaymentStatus(orderID, transactionStatus, transactionID, vaNumber, bankType, qrCodeURL, expiryTime, string(webhookJSON))
}

func (s *paymentService) CheckPaymentStatus(paymentID string) (*model.Payment, error) {
	payment, err := s.paymentRepo.FindByID(paymentID)
	if err != nil {
		return nil, err
	}

	// Check status from Midtrans if transaction ID exists
	if payment.MidtransTransactionID != nil && s.cfg.MidtransServerKey != "" {
		if err := s.CheckPaymentStatusFromMidtrans(payment.OrderID); err != nil {
			log.Printf("⚠️  Failed to check payment status from Midtrans: %v", err)
		}
		// Reload payment after status check
		payment, _ = s.paymentRepo.FindByID(paymentID)
	}

	return payment, nil
}

// CheckPaymentStatusFromMidtrans checks payment status from Midtrans API
func (s *paymentService) CheckPaymentStatusFromMidtrans(orderID string) error {
	// Get payment from database first
	payment, err := s.paymentRepo.FindByOrderNumber(orderID)
	if err != nil {
		return fmt.Errorf("payment not found: %v", err)
	}

	// If already successful, skip check
	if payment.Status == model.PaymentStatusSuccess {
		return nil
	}

	// If no transaction ID, cannot check
	if payment.MidtransTransactionID == nil || *payment.MidtransTransactionID == "" {
		return fmt.Errorf("no transaction ID for payment")
	}

	// Call Midtrans status API
	baseURL := s.getMidtransBaseURL()
	authHeader := s.getAuthHeader()
	url := fmt.Sprintf("%s/%s/status", baseURL, *payment.MidtransTransactionID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call Midtrans API: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("⚠️  Midtrans API returned status %d: %s", resp.StatusCode, string(body))
		return fmt.Errorf("Midtrans API error: %s", string(body))
	}

	var midtransResp map[string]interface{}
	if err := json.Unmarshal(body, &midtransResp); err != nil {
		return fmt.Errorf("failed to parse response: %v", err)
	}

	// Extract status information
	transactionStatus, _ := midtransResp["transaction_status"].(string)
	transactionID, _ := midtransResp["transaction_id"].(string)

	var vaNumber, bankType, qrCodeURL string
	if vaNumbers, ok := midtransResp["va_numbers"].([]interface{}); ok && len(vaNumbers) > 0 {
		if va, ok := vaNumbers[0].(map[string]interface{}); ok {
			vaNumber, _ = va["va_number"].(string)
			bankType, _ = va["bank"].(string)
		}
	}

	// Extract QR code URL from actions
	if actions, ok := midtransResp["actions"].([]interface{}); ok && len(actions) > 0 {
		for _, action := range actions {
			if act, ok := action.(map[string]interface{}); ok {
				name, _ := act["name"].(string)
				url, _ := act["url"].(string)
				if (name == "generate-qr-code" || name == "generate-qr-code-v2" || name == "qr-code") && url != "" {
					qrCodeURL = url
					log.Printf("✅ Found QR code URL from action '%s': %s", name, qrCodeURL)
					break
				}
			}
		}
		// If not found by name, try by method GET
		if qrCodeURL == "" {
			for _, action := range actions {
				if act, ok := action.(map[string]interface{}); ok {
					method, _ := act["method"].(string)
					url, _ := act["url"].(string)
					if method == "GET" && url != "" && strings.Contains(strings.ToLower(url), "qr") {
						qrCodeURL = url
						log.Printf("✅ Found QR code URL from GET method: %s", qrCodeURL)
						break
					}
				}
			}
		}
	}

	// If QR code URL not found in response but payment already has one, preserve it
	if qrCodeURL == "" && payment.QRCodeURL != nil && *payment.QRCodeURL != "" {
		log.Printf("⚠️  QR code URL not in response, preserving existing: %s", *payment.QRCodeURL)
		qrCodeURL = *payment.QRCodeURL
	}

	// Extract expiry time
	var expiryTime *time.Time
	if expiry, ok := midtransResp["expiry_time"].(string); ok && expiry != "" {
		formats := []string{
			time.RFC3339,
			"2006-01-02 15:04:05",
			"2006-01-02T15:04:05",
		}
		for _, format := range formats {
			exp, err := time.Parse(format, expiry)
			if err == nil {
				expiryTime = &exp
				break
			}
		}
	}

	webhookJSON, _ := json.Marshal(midtransResp)

	return s.UpdatePaymentStatus(orderID, transactionStatus, transactionID, vaNumber, bankType, qrCodeURL, expiryTime, string(webhookJSON))
}

// UpdatePaymentStatus updates payment status from Midtrans webhook or status check
func (s *paymentService) UpdatePaymentStatus(orderID string, status string, transactionID string, vaNumber string, bankType string, qrCodeURL string, expiryTime *time.Time, midtransResponse string) error {
	paymentStatus := mapMidtransStatusToPaymentStatus(status)

	// Get payment by order number
	payment, err := s.paymentRepo.FindByOrderNumber(orderID)
	if err != nil {
		log.Printf("❌ Payment not found for order %s: %v", orderID, err)
		return fmt.Errorf("payment not found for order: %s", orderID)
	}

	// Preserve existing values if new ones are empty
	if qrCodeURL == "" && payment.QRCodeURL != nil && *payment.QRCodeURL != "" {
		qrCodeURL = *payment.QRCodeURL
	}
	if vaNumber == "" && payment.VANumber != nil && *payment.VANumber != "" {
		vaNumber = *payment.VANumber
	}
	if bankType == "" && payment.BankType != nil && *payment.BankType != "" {
		bankType = *payment.BankType
	}

	// Update payment fields
	payment.Status = paymentStatus
	if transactionID != "" {
		payment.MidtransTransactionID = &transactionID
	}
	if vaNumber != "" {
		payment.VANumber = &vaNumber
	}
	if bankType != "" {
		payment.BankType = &bankType
	}
	if qrCodeURL != "" {
		payment.QRCodeURL = &qrCodeURL
	}
	if expiryTime != nil {
		payment.ExpiryTime = expiryTime
	}
	if midtransResponse != "" {
		payment.MidtransResponse = &midtransResponse
		// Extract fraud_status from midtransResponse if available
		var responseMap map[string]interface{}
		if err := json.Unmarshal([]byte(midtransResponse), &responseMap); err == nil {
			if fraudStatus, ok := responseMap["fraud_status"].(string); ok && fraudStatus != "" {
				payment.FraudStatus = &fraudStatus
			}
		}
	}

	if err := s.paymentRepo.Update(payment); err != nil {
		log.Printf("❌ Failed to update payment: %v", err)
		return err
	}

	// Update order status if payment is successful
	if paymentStatus == model.PaymentStatusSuccess {
		order, err := s.orderRepo.FindByID(payment.OrderUUID)
		if err == nil && order.Status == "pending" {
			order.Status = "processing"
			if err := s.orderRepo.Update(order); err != nil {
				log.Printf("⚠️  Failed to update order status: %v", err)
			}
		}
	}

	return nil
}
