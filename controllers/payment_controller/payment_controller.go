package paymentcontroller

import (
	"encoding/json"
	"fmt"
	"net/http"

	paymentrepo "decoration_project/repository/payment_repo"
	notificationservices "decoration_project/services/notification_services"
	paymentservices "decoration_project/services/payment_services"
	"decoration_project/utils"
)

type PaymentController struct {
	Service *paymentservices.PaymentService
}

// NewPaymentController initializes the controller with Razorpay credentials
func NewPaymentController() *PaymentController {
	apiKey := "rzp_test_RWtknBdIxDzT1t"       // Replace with your Razorpay key
	apiSecret := "FPQ60I8xVWINQcWZBdrvXcgN"   // Replace with your Razorpay secret

	repo := paymentrepo.NewRazorpayRepository(apiKey, apiSecret)
	service := paymentservices.NewPaymentService(repo)

	return &PaymentController{Service: service}
}

// RefundHandler handles user refund requests
func (c *PaymentController) RefundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// ✅ Step 1: Validate JWT token and extract user ID
	userID, err := utils.ValidateToken(r)
	if err != nil {
		utils.SendResponse(w, http.StatusUnauthorized, false, nil, "Unauthorized: "+err.Error())
		return
	}

	// ✅ Step 2: Parse request body
	var req struct {
		BookingID string `json:"booking_id"`
		PaymentID string `json:"payment_id"`
		TransactionId string `json:"transaction_id"`
		Amount    int    `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid request body: "+err.Error())
		return
	}

	// ✅ Step 3: Validate required fields
	if req.BookingID == "" {
		utils.SendResponse(w, http.StatusBadRequest, false, nil, "Booking ID is required")
		return
	}
	if req.TransactionId == "" {
		utils.SendResponse(w, http.StatusBadRequest, false, nil, "Transaction ID is required")
		return
	}
	if req.PaymentID == "" {
		utils.SendResponse(w, http.StatusBadRequest, false, nil, "Payment ID is required")
		return
	}


	// ✅ Step 4: Process refund via service layer
	result, err := c.Service.Refund(req.PaymentID, req.TransactionId, req.Amount, req.BookingID)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, false, nil, fmt.Sprintf("Refund failed: %v", err))
		return
	}

	// ✅ Step 5: Send success response
	utils.SendResponse(w, http.StatusOK, true, result, fmt.Sprintf("Refund for booking %s processed successfully", req.BookingID))

	// ✅ Step 6: Notify user about refund
	go func() {
		notifyErr := notificationservices.SendBookingNormalNotificationToUser(
			userID,
			"Refund Under Progress",
			fmt.Sprintf("Your cancelled order (Booking ID: %s) refund has been processed and will be credited back to your source account within 3–5 working days.", req.BookingID),
		)
		if notifyErr != nil {
			fmt.Println("⚠️ Failed to send refund notification:", notifyErr)
		}
	}()
}
