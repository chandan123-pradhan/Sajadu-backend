package staffmodel

type UpdateLocationRequest struct {
    BookingID string  `json:"booking_id"`
    Latitude  float64 `json:"latitude"`
    Longitude float64 `json:"longitude"`
}
