package restorantservices

import (
	restorantmodels "decoration_project/models/restorant_models"
	restorantrepo "decoration_project/repository/restorant_repo"
	"decoration_project/utils"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// RegisterRestaurant handles restaurant registration with password hashing
func RegisterRestaurant(restaurant restorantmodels.RestaurantProfile, password string, images []string) (string, error) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// Assign hashed password
	restaurantPassword := string(hashedPassword)

	// Add restaurant in repository
	restaurantID, err := restorantrepo.AddRestaurant(restaurant, restaurantPassword)
	if err != nil {
		return "", err
	}

	// Save images if provided
	if len(images) > 0 {
		err = restorantrepo.AddRestaurantImages(restaurantID, images)
		if err != nil {
			return "", err
		}
	}

	return restaurantID, nil
}

func LoginRestaurant(email, password string) (map[string]interface{}, error) {
	// Fetch restaurant with images + hashed password
	restaurant, images, storedPassword, err := restorantrepo.GetRestaurantWithImages(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate auth token
	token, err := utils.GenerateRestaurantToken(restaurant.RestaurantID)
	if err != nil {
		return nil, err
	}

	// Prepare response
	response := map[string]interface{}{
		"token":      token,
		"restaurant": restaurant,
		"images":     images,
	}

	return response, nil
}



func UpdateFcmToken(restorantId, fcmToken string) error {
    // Ensure userID and fcmToken are valid
    if restorantId == "" || fcmToken == "" {
        return errors.New("invalid Restorant ID or FCM token")
    }

    // Call repository layer
    err := restorantrepo.UpdateFcmToken(restorantId, fcmToken)
    if err != nil {
        return err
    }

    return nil
}

func GetRestorantFcmTokenService(restorantId string) (string, error) {
	return restorantrepo.GetRestorantFcmToken(restorantId)
}



func UpdateRestaurantProfileService(restaurant restorantmodels.RestaurantProfile, imagePaths []string) error {
	if restaurant.RestaurantID == "" {
		return errors.New("restaurant ID is required")
	}

	// Update restaurant profile
	err := restorantrepo.UpdateRestaurantProfile(restaurant)
	if err != nil {
		return err
	}

	// If new images provided → update them
	if len(imagePaths) > 0 {
		err = restorantrepo.UpdateRestaurantImages(restaurant.RestaurantID, imagePaths)
		if err != nil {
			return err
		}
	}

	return nil
}
