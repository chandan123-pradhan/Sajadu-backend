package adminmodel

import restorantmodels "decoration_project/models/restorant_models"

type ServiceDetailsResponse struct {
    Service           restorantmodels.RestaurantService         `json:"service"`
    RestaurantProfile RestaurantProfileResponse `json:"restaurant_profile,omitempty"`
}


type RestaurantProfileResponse struct {
    RestaurantID string   `json:"restaurant_id"`
    Name         string   `json:"name"`
    Images       []string `json:"images"`
}
