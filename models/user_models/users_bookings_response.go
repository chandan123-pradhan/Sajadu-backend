package usermodels

import "time"

// BookingRequest
type BookingRequest struct {
	UserID        string    `json:"user_id"`
	RestaurantID  string    `json:"restaurant_id"`
	ServiceID     string    `json:"service_id"`
	ScheduledDate time.Time `json:"scheduled_date"`
	ServiceName   string    `json:"service_name"`

	// Location
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Pincode   string  `json:"pincode"`
	State     string  `json:"state"`
	City      string  `json:"city"`

	// Pricing
	Price float64 `json:"price"`

	// Payment Info (optional at booking time)
	PaymentMethod string  `json:"payment_method,omitempty"` // e.g. UPI, Card, Wallet
	TransactionID string  `json:"transaction_id,omitempty"` // from payment gateway
	AmountPaid    float64 `json:"amount_paid,omitempty"`
	Currency      string  `json:"currency,omitempty"` // e.g. INR, USD
}

// PaymentResponse (for nested payment details in BookingResponse)
type PaymentResponse struct {
	PaymentID     string    `json:"payment_id"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	PaymentMethod string    `json:"payment_method"`
	TransactionID string    `json:"transaction_id"`
	PaymentDate   time.Time `json:"payment_date"`
	PaymentStatus string    `json:"payment_status"`
}

// BookingResponse
type BookingResponse struct {
	BookingID     string           `json:"booking_id"`
	UserID        string           `json:"user_id"`
	RestaurantID  string           `json:"restaurant_id"`
	ServiceID     string           `json:"service_id"`
	ServiceName   string           `json:"service_name"`
	CompleteOtp	  string 		   `json:"complete_otp"`
	Price         float64          `json:"price"`
	Status        string           `json:"status"`
	ScheduledDate string           `json:"scheduled_date"`
	Address       string           `json:"address"`
	Pincode       string           `json:"pincode"`
	State         string           `json:"state"`
	City          string           `json:"city"`
	CreatedAt     time.Time        `json:"created_at"`
	Payment       *PaymentResponse `json:"payment,omitempty"` // nil if no payment
}


type UserBookingsWrapper struct {
	Bookings []BookingResponse `json:"bookings"`
}