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





func SendOtp(mobileNo string) (string, error) {

    // 1. Check user exists
    restorant, err := restorantrepo.GetRestaurantByMobileNo(mobileNo)
    if err != nil {
        return "", errors.New("mobile number not registered")
    }

    // 2. Now user exists → pass userID to SendOTP
    restorantId := &restorant.RestaurantID

    otp, err := restorantrepo.SendOTP(mobileNo, restorantId)
    if err != nil {
        return "", err
    }

    return otp, nil
}


func VerifyOtp(mobile string, otp string) (map[string]interface{}, error) {

    // 1. Verify OTP
    _, err := restorantrepo.VerifyOTP(mobile, otp)
    if err != nil {
        return nil, err
    }

    // 2. Fetch restaurant data using mobile number
    restorant, err := restorantrepo.GetRestaurantByMobileNo(mobile)
    if err != nil {
        return nil, errors.New("restaurant not found")
    }

    // 3. Now fetch restaurant details + images + hashed password (same as login)
    restaurantData, images, _, err := restorantrepo.GetRestaurantWithImages(restorant.Email)
    if err != nil {
        return nil, errors.New("failed to fetch restaurant data")
    }

    // 4. Generate token (same as login)
    token, err := utils.GenerateRestaurantToken(restaurantData.RestaurantID)
    if err != nil {
        return nil, errors.New("failed to generate token")
    }

    // 5. Final response (same structure as LoginRestaurant)
    response := map[string]interface{}{
        "token":      token,
        "restaurant": restaurantData,
        "images":     images,
    }

    return response, nil
}
