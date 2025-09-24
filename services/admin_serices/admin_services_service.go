package adminserices

import (
	restorantmodels "decoration_project/models/restorant_models"
	adminrepo "decoration_project/repository/admin_repo"
)

func CreateService(service restorantmodels.RestaurantService, images []string) (string, error) {
    serviceID, err := adminrepo.AddService(service)
    if err != nil {
        return "", err
    }

    if len(images) > 0 {
        err = adminrepo.AddServiceImages(serviceID, images)
        if err != nil {
            return "", err
        }
    }

    return serviceID, nil
}

func GetServicesDetails(serviceID string) (restorantmodels.RestaurantService, error) {
    service, images, err := adminrepo.GetServiceWithImages(serviceID)
    if err != nil {
        return service, err
    }
    service.Images = images
    return service, nil
}

func GetAllServiceCategoryWise(categoryId string) ([]restorantmodels.RestaurantService, error) {
    // Fetch all services with images from repository
    services, err := adminrepo.GetAllServiceCategoryWise(categoryId)
    if err != nil {
        return nil, err
    }
    return services, nil
}