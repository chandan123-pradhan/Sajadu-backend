package staffmodel

import (
	"time"
	usermodels "decoration_project/models/user_models"
)

type StaffAssignedServicesDetails struct {
	BookingID          string                      `json:"booking_id"`
	User               usermodels.UserDetailsModel `json:"user"`
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
	Payment            *usermodels.PaymentResponse `json:"payment,omitempty"`
}
