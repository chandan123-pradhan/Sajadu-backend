package adminserices

import (
	adminmodel "decoration_project/models/admin_model"
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


func GetServicesDetails(serviceID string) (adminmodel.ServiceDetailsResponse, error) {
    // Step 1: Fetch service + service images + restaurant profile data
    service, serviceImages, restaurantProfile, restaurantImages, err :=
        adminrepo.GetServiceWithImagesAndRestaurant(serviceID)

    if err != nil {
        return adminmodel.ServiceDetailsResponse{}, err
    }

    // Step 2: Attach service images
    service.Images = serviceImages

    // Step 3: Prepare final response object
    response := adminmodel.ServiceDetailsResponse{
        Service: service,
    }

    // Step 4: Add restaurant only if exists
    if restaurantProfile.RestaurantID != "" {
        response.RestaurantProfile = adminmodel.RestaurantProfileResponse{
            RestaurantID: restaurantProfile.RestaurantID,
            Name:         restaurantProfile.Name,
            Images:       restaurantImages,
        }
    }

    return response, nil
}


func GetAllServiceCategoryWise(categoryId string) ([]restorantmodels.RestaurantService, error) {
    // Fetch all services with images from repository
    services, err := adminrepo.GetAllServiceCategoryWise(categoryId)
    if err != nil {
        return nil, err
    }
    return services, nil
}



func CreateFestival(f adminmodel.Festival) error {
    return adminrepo.CreateFestival(f)
}


func GetAllFestivals() ([]adminmodel.Festival, error) {
    return adminrepo.GetAllFestivals()
}


func AddServiceToFestival(festivalID string, serviceID string, discountPercent float64) error {
    return adminrepo.AddServiceToFestival(festivalID, serviceID, discountPercent)
}



func GetFestivalServices(festivalId string) ([]restorantmodels.RestaurantService, error) {
	return adminrepo.GetFestivalServices(festivalId)
}

