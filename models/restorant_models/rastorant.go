package restorantmodels

import "time"

type RestaurantProfile struct {
    RestaurantID string   `json:"restaurant_id"`
    Name         string   `json:"name"`
    Email        string   `json:"email"`
    PhoneNumber  string   `json:"phone_number"`
    Address      string   `json:"address"`
    City         string   `json:"city"`
    PostalCode string   `json:"postalCode"`
    State        string   `json:"state"`
    Country      string   `json:"country"`
    Images             []string `json:"images,omitempty"`
    Latitude     float64  `json:"latitude"`
    Longitude    float64  `json:"longitude"`
    Password     string   `json:"password,omitempty"` // omit from JSON response
}

type RestaurantImage struct {
    ImageID      string    `json:"image_id" db:"image_id"`
    RestaurantID string    `json:"restaurant_id" db:"restaurant_id"`
    ImageURL     string    `json:"image_url" db:"image_url"`
    CreatedAt    time.Time `json:"created_at" db:"created_at"`
}