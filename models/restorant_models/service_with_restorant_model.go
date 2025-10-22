package restorantmodels

import "time"

type ServiceWithRestaurant struct {
	Service RestaurantService `json:"service"`
	Reviews []ServiceReview   `json:"reviews,omitempty"`
}

type ServiceReview struct {
	UserName   string    `json:"user_name"`
	Rating     float64   `json:"rating"`
	ReviewText string    `json:"review_text"`
	CreatedAt  time.Time `json:"created_at"`
}
