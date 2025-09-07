package restorantmodels

import (
	usermodels "decoration_project/models/user_models"
	"time"
)

type RestorantBookingsResponse struct {
	BookingID     string           `json:"booking_id"`
	UserID        string           `json:"user_id"`
	ServiceID     string           `json:"service_id"`
	StartOtp 	  string 		   `json:"start_otp"`
	ServiceName   string           `json:"service_name"`
	Price         float64          `json:"price"`
	Status        string           `json:"status"`
	ScheduledDate string           `json:"scheduled_date"`
	Address       string           `json:"address"`
	Pincode       string           `json:"pincode"`
	State         string           `json:"state"`
	City          string           `json:"city"`
	CreatedAt     time.Time        `json:"created_at"`
	Payment       *usermodels.PaymentResponse `json:"payment,omitempty"` // nil if no payment
}



// Wrapper response for restaurant bookings
type RestaurantBookingsWrapper struct {
	Bookings []RestorantBookingsResponse `json:"bookings"`
}
