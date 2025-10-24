package paymentrepo

import (
	"bytes"
	"decoration_project/config"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// RazorpayRepository handles Razorpay operations and DB updates
type RazorpayRepository struct {
	APIKey    string
	APISecret string
}

func NewRazorpayRepository(apiKey, apiSecret string) *RazorpayRepository {
	return &RazorpayRepository{
		APIKey:    apiKey,
		APISecret: apiSecret,
	}
}

func (r *RazorpayRepository) Refund(paymentID string, transactionId string, amount int, bookingID string) (map[string]interface{}, error) {
	if paymentID == "" || bookingID == ""  || transactionId==""{
		return nil, errors.New("paymentID and bookingID and transaction id are required")
	}

	// Razorpay refund API endpoint
	url := fmt.Sprintf("https://api.razorpay.com/v1/payments/%s/refund", transactionId)

	body := map[string]interface{}{
		"amount": amount, // amount in paise
	}

	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	req.SetBasicAuth(r.APIKey, r.APISecret)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("refund failed: %+v", result)
	}

	// ✅ Update the specific payment by paymentID
	paymentUpdateQuery := `
		UPDATE payments
		SET payment_status = 'REFUNDED',
		    updated_at = NOW()
		WHERE payment_id = ?;
	`
	_, payErr := config.DB.Exec(paymentUpdateQuery, paymentID)
	if payErr != nil {
		fmt.Println("⚠️ Failed to update payment status:", payErr)
	}

	// ✅ Update booking's overall payment_status
	bookingUpdateQuery := `
		UPDATE bookings
		SET payment_status = 'REFUNDED',
		    status_id = ?, -- refunded status_id
		    updated_at = NOW()
		WHERE booking_id = ?;
	`
	_, dbErr := config.DB.Exec(bookingUpdateQuery, 4, bookingID) // 4 = refunded status
	if dbErr != nil {
		fmt.Println("⚠️ Failed to update booking refund status:", dbErr)
	}

	return result, nil
}
