package userservices

import (
	usermodels "decoration_project/models/user_models"
	userrepo "decoration_project/repository/user_repo"
	"fmt"
)

// AddServiceReview handles the business logic for adding a review
func AddServiceReview(req usermodels.ServiceReview) error {
	// Here you can add validation or pre-checks if needed
	if req.Rating < 1 || req.Rating > 5 {
		return fmt.Errorf("rating must be between 1 and 5")
	}

	// Call repository layer to insert review
	return userrepo.AddServiceReview(req)
}
