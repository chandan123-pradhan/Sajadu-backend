package restorantservices

import (
	restorantmodels "decoration_project/models/restorant_models"
	restorantrepo "decoration_project/repository/restorant_repo"
)

func AddStaffService(staff restorantmodels.RestaurantStaff) (string, error) {
    // Insert staff with image in one go
    staffID, err := restorantrepo.AddStaff(staff)
    if err != nil {
        return "", err
    }

    return staffID, nil
}

func GetStaffByRestaurant(restaurantID string) ([]restorantmodels.RestaurantStaff, error) {
    return restorantrepo.GetStaffByRestaurant(restaurantID)
}