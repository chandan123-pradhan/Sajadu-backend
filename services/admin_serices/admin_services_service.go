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

