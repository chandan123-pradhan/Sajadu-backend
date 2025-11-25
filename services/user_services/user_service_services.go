package userservices

import (
	adminmodel "decoration_project/models/admin_model"
	restorantmodels "decoration_project/models/restorant_models"
	"time"

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

// ======================= GetAllFestivals =======================
// Fetches all festivals (active or all) for users
func GetAllFestivals() ([]adminmodel.Festival, error) {
	return userrepo.GetAllFestivals()
}

// ======================= GetFestivalServices =======================
// Fetches all services linked to a given festival including discount info
func GetFestivalServices(festivalID string) ([]restorantmodels.RestaurantService, error) {
	services, err := userrepo.GetFestivalServices(festivalID)
	if err != nil {
		return nil, err
	}

	// Ensure DiscountPercent is set (0 if not present)
	for i := range services {
		if services[i].DiscountPercent == 0 {
			services[i].DiscountPercent = 0
		}
		// Ensure CreatedAt/UpdatedAt are valid (optional)
		if services[i].CreatedAt.IsZero() {
			services[i].CreatedAt = time.Now()
		}
		if services[i].UpdatedAt.IsZero() {
			services[i].UpdatedAt = time.Now()
		}
	}

	return services, nil
}