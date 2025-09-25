package adminmodel

type UpdateBookingStatusRequest struct {
	BookingID    string `json:"booking_id"`
	NewStatus    string `json:"new_status"`    // e.g. "Accepted" or "Rejected"
	RestaurantID string `json:"restaurant_id"` // optional, only if Accepted
}