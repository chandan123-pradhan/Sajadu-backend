package userservices

import (
	restorantmodels "decoration_project/models/restorant_models"
	
	userrepo "decoration_project/repository/user_repo"
)

func GetServicesByCategory(categoryID string) ([]restorantmodels.RestaurantService, error) {
	return userrepo.GetServicesByCategory(categoryID)
}

func GetServiceDetails(restaurantID, serviceID string) (restorantmodels.ServiceWithRestaurant, error) {
	// Call repository
	return userrepo.GetServiceDetails(restaurantID, serviceID)
}