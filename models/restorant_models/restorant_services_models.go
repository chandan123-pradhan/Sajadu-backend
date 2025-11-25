package restorantmodels

import "time"

// RestaurantService represents a service offered by a restaurant
type RestaurantService struct {
    ServiceID          string    `json:"service_id"`
    CategoryId         string    `json:"category_id"`
    ServiceName        string    `json:"service_name"`
    ServiceDescription string    `json:"service_description"`
    ServicePrice       float64   `json:"service_price"`
    AverageRating      float64   `json:"average_rating"`          // new field for average rating
    Images             []string  `json:"images,omitempty"`
    CreatedAt          time.Time `json:"created_at"`
    UpdatedAt          time.Time `json:"updated_at"`
    ProposedRestorantId string   `json:"proposed_restorant_id"`
    DiscountPercent     float64   `json:"discount_percent,omitempty"`
}

// ServiceImage represents an image associated with a service
type ServiceImage struct {
    ImageID   string    `json:"image_id"`
    ServiceID string    `json:"service_id"`
    ImageURL  string    `json:"image_url"`
    CreatedAt time.Time `json:"created_at"`
}
