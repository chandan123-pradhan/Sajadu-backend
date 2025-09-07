package restorantservices

import (
	restorantmodels "decoration_project/models/restorant_models"
	restorantrepo "decoration_project/repository/restorant_repo"
)

func CreateService(service restorantmodels.RestaurantService, images []string) (string, error) {
    serviceID, err := restorantrepo.AddService(service)
    if err != nil {
        return "", err
    }

    if len(images) > 0 {
        err = restorantrepo.AddServiceImages(serviceID, images)
        if err != nil {
            return "", err
        }
    }

    return serviceID, nil
}

func GetService(serviceID string) (restorantmodels.RestaurantService, error) {
    service, images, err := restorantrepo.GetServiceWithImages(serviceID)
    if err != nil {
        return service, err
    }
    service.Images = images
    return service, nil
}

func GetAllServicesForRestaurant(restaurantID string) ([]restorantmodels.RestaurantService, error) {
    // Fetch all services with images from repository
    services, err := restorantrepo.GetAllServicesWithImages(restaurantID)
    if err != nil {
        return nil, err
    }
    return services, nil
}