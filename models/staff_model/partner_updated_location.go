package staffmodel

import "time"

type PartnerLocationResponse struct {
	PartnerID  string    `json:"partner_id"`
	BookingID  string    `json:"booking_id"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	UpdatedAt  time.Time `json:"updated_at"`
}
