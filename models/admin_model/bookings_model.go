package adminmodel

import "time"

// PaymentResponse represents payment details for a booking
type PaymentResponse struct {
	PaymentID     string    `json:"payment_id"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	PaymentMethod string    `json:"payment_method"`
	TransactionID string    `json:"transaction_id"`
	PaymentDate   time.Time `json:"payment_date"`
	PaymentStatus string    `json:"payment_status"`
}

// BookingResponse represents detailed booking info
type BookingResponse struct {
	BookingID       string           `json:"booking_id"`
	UserID          string           `json:"user_id"`
	RestaurantID    string           `json:"restaurant_id"`
	ServiceID       string           `json:"service_id"`
	ServiceName     string           `json:"service_name"`	
	ServiceDesc     string           `json:"service_description,omitempty"`
	Price           float64          `json:"price"`
	Status          string           `json:"status"`
	ScheduledDate   string           `json:"scheduled_date"`
	Address         string           `json:"address"`
	Latitude        float64          `json:"latitude,omitempty"`
	Longitude       float64          `json:"longitude,omitempty"`
	Pincode         string           `json:"pincode"`
	State           string           `json:"state"`
	City            string           `json:"city"`
	CreatedAt       time.Time        `json:"created_at"`
	Payment         *PaymentResponse `json:"payment,omitempty"` // nil if no payment
	Images          []string         `json:"images"`            // service images
	RestaurantName  string           `json:"restaurant_name,omitempty"`
	RestaurantImages []string        `json:"restaurant_images,omitempty"`
}
