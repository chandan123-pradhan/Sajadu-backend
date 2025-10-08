package adminmodel

type BookingServiceSummary struct {
	ImageURL    string  `json:"image_url"`
	ServiceName string  `json:"service_name"`
	UserID      string  `json:"user_id"`
	Price       float64 `json:"price"`
}
