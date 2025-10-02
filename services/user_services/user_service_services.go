package userservices

import (
	restorantmodels "decoration_project/models/restorant_models"
	
	userrepo "decoration_project/repository/user_repo"
)

func GetServicesByCategory(categoryID string) ([]restorantmodels.RestaurantService, error) {
	return userrepo.GetServicesByCategory(categoryID)
}

func GetServiceDetails(serviceID string) (restorantmodels.ServiceWithRestaurant, error) {
	// Call repository
	return userrepo.GetServiceDetails(serviceID)
}

func SearchServices(query string) ([]restorantmodels.RestaurantService, error) {
    return userrepo.SearchServicesByName(query)
}