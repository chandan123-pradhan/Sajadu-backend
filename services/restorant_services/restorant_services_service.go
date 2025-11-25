package restorantservices

import (
	restorantmodels "decoration_project/models/restorant_models"
	restorantrepo "decoration_project/repository/restorant_repo"
	"fmt"
)


func GetAllProposedServices(restorantId string) ([]restorantmodels.RestaurantService, error) {
    // Fetch all services with images from repository
    services, err := restorantrepo.GetAllProposedServices(restorantId)
    if err != nil {
        return nil, err
    }
    return services, nil
}



func DeleteProposedService(serviceID string, restaurantID string) error {

	// Check if service belongs to this restaurant
	isOwner, err := restorantrepo.CheckServiceOwner(serviceID, restaurantID)
	if err != nil {
		return err
	}

	if !isOwner {
		return fmt.Errorf("You are not authorized to delete this service")
	}

	// Delete service + related images
	err = restorantrepo.DeleteServiceWithImages(serviceID)
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("Failed to delete service")
	}

	return nil
}
