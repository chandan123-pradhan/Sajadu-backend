package usermodels

import (
	// staffmodel "decoration_project/models/staff_model"

	commonmodel "decoration_project/models/common_model"
	"time"
	// usermodels "decoration_project/models/user_models"
)

type GetUsersBookingDetailsResponse struct {
	BookingID          string                      `json:"booking_id"`
	User               commonmodel.StaffBasic `json:"staff"`
	ServiceID          string                      `json:"service_id"`
	ServiceName        string                      `json:"service_name"`
	ServiceDescription string                      `json:"service_description"`
	ServicePrice       float64                     `json:"service_price"`
	Items              string                      `json:"items"`
	Images             []string                    `json:"images"`
	Status             string                      `json:"status"`
	ScheduledDate      string                      `json:"scheduled_date"`
	Address            string                      `json:"address"`
	Pincode            string                      `json:"pincode"`
	State              string                      `json:"state"`
	City               string                      `json:"city"`
	Latitude           float64                     `json:"latitude"`
	Longitude          float64                     `json:"longitude"`
	CreatedAt          time.Time                   `json:"created_at"`
	Payment            *PaymentResponse `json:"payment,omitempty"`
}
