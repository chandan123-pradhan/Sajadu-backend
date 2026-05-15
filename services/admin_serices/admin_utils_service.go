package adminserices

import (
	usermodels "decoration_project/models/user_models"
	adminrepo "decoration_project/repository/admin_repo"
)


func GetAllUsers() ([]usermodels.UserDetailsModel, error) {
    // Fetch all services with images from repository
    users, err := adminrepo.GetAllUsers()
    if err != nil {
        return nil, err
    }
    return users, nil
}