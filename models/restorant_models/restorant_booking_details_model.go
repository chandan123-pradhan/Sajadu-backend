package restorantmodels

import (
	usermodels "decoration_project/models/user_models"
	"time"
)

// Detailed response for a single booking
type BookingDetailsResponse struct {
	BookingID     string                      `json:"booking_id"`
	User          usermodels.UserDetailsModel                 `json:"user"`
	Staff         *StaffDetails               `json:"staff,omitempty"`   // nil if not assigned
	ServiceID     string                      `json:"service_id"`
	ServiceName   string                      `json:"service_name"`
	ServiceDesc   string                      `json:"service_description"`
	Price         float64                     `json:"price"`
	Status        string                      `json:"status"`
	ScheduledDate string                      `json:"scheduled_date"`
	Address       string                      `json:"address"`
	Pincode       string                      `json:"pincode"`
	State         string                      `json:"state"`
	City          string                      `json:"city"`
	Latitude      float64                     `json:"latitude"`
	Longitude     float64                     `json:"longitude"`
	StartOtp      string                      `json:"start_otp"`          // decrypted OTP (empty if not available)
	CreatedAt     time.Time                   `json:"created_at"`
	Payment       *usermodels.PaymentResponse `json:"payment,omitempty"` // nil if no payment
}

// Staff details (subset)
type StaffDetails struct {
	StaffID      string `json:"staff_id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	WhatsappNo   string `json:"whatsapp_no"`
	Designation  string `json:"designation"`
	Description  string `json:"description"`
	ImageURL     string `json:"image_url"`
}
