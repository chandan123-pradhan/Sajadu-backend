package restorantmodels

import "time"

type RestaurantService struct {
    ServiceID          string   `json:"service_id"`
    RestaurantID       string   `json:"restaurant_id"`
    CategoryId         string   `json:"category_id"`
    ServiceName        string   `json:"service_name"`
    ServiceDescription string   `json:"service_description"`
    ServicePrice       float64  `json:"service_price"`
    Items              string   `json:"items"` // JSON string
    Images             []string `json:"images,omitempty"`
    CreatedAt          time.Time `json:"created_at"`
    UpdatedAt          time.Time `json:"updated_at"`
}

type ServiceImage struct {
    ImageID   string    `json:"image_id"`
    ServiceID string    `json:"service_id"`
    ImageURL  string    `json:"image_url"`
    CreatedAt time.Time `json:"created_at"`
}
