package restorantcontrollers

import (
	restorantmodels "decoration_project/models/restorant_models"
	restorantservices "decoration_project/services/restorant_services"
	"decoration_project/utils"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-sql-driver/mysql"
)

// ======================= REGISTER RESTAURANT =======================
func RegisterRestaurant(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form data (for images + JSON fields)
	err := r.ParseMultipartForm(10 << 20) // 10MB max upload
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "Failed to parse form data")
		return
	}

	// Extract restaurant details from "restaurant" field (JSON string)
	var restaurant restorantmodels.RestaurantProfile
	err = json.Unmarshal([]byte(r.FormValue("restaurant")), &restaurant)
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "Invalid restaurant data")
		return
	}

	// Validate required fields
	if restaurant.Name == "" || restaurant.Email == "" || restaurant.Password == "" {
		utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "Name, email and password are required")
		return
	}

	// Handle multiple image files
	var imagePaths []string
	files := r.MultipartForm.File["images"] // key = "images"
	for _, fileHeader := range files {
		path, err := utils.SaveFile(fileHeader, "uploads/restaurant_images")
		if err != nil {
			utils.SendResponse(w, http.StatusInternalServerError, false, map[string]interface{}{}, "Failed to save image")
			return
		}
		imagePaths = append(imagePaths, path)
	}

	// Register restaurant with images
	restaurantID, err := restorantservices.RegisterRestaurant(restaurant, restaurant.Password, imagePaths)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "Email already registered")
			return
		}
		fmt.Println(err.Error())
		utils.SendResponse(w, http.StatusInternalServerError, false, map[string]interface{}{}, err.Error())
		return
	}

	// Generate auth token
	token, err := utils.GenerateRestaurantToken(restaurantID)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, false, map[string]interface{}{}, "Failed to generate auth token")
		return
	}

	// Prepare response with images
	responseData := map[string]interface{}{
		"restaurant_id": restaurantID,
		"name":          restaurant.Name,
		"email":         restaurant.Email,
		"phone_number":  restaurant.PhoneNumber,
		"address":       restaurant.Address,
		"city":          restaurant.City,
		"state":         restaurant.State,
		"postal_code":   restaurant.PostalCode,
		"country":       restaurant.Country,
		"latitude":      restaurant.Latitude,
		"longitude":     restaurant.Longitude,
		"auth_token":    token,
		"images":        imagePaths,
	}

	utils.SendResponse(w, http.StatusCreated, true, responseData, "Restaurant registered successfully")
}


func LoginRestaurant(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "Invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "Email and password are required")
		return
	}

	// Call service
	loginData, err := restorantservices.LoginRestaurant(req.Email, req.Password)
	if err != nil {
		utils.SendResponse(w, http.StatusUnauthorized, false, map[string]interface{}{}, err.Error())
		return
	}

	// Extract service result
	restaurant := loginData["restaurant"].(restorantmodels.RestaurantProfile)
	images := loginData["images"].([]string)
	token := loginData["token"].(string)

	// Prepare response in same format as RegisterRestaurant
	responseData := map[string]interface{}{
		"restaurant_id": restaurant.RestaurantID,
		"name":          restaurant.Name,
		"email":         restaurant.Email,
		"phone_number":  restaurant.PhoneNumber,
		"address":       restaurant.Address,
		"city":          restaurant.City,
		"state":         restaurant.State,
		"postal_code":   restaurant.PostalCode,
		"country":       restaurant.Country,
		"latitude":      restaurant.Latitude,
		"longitude":     restaurant.Longitude,
		"auth_token":    token,
		"images":        images,
	}

	utils.SendResponse(w, http.StatusOK, true, responseData, "Login successful")
}



// ======================= UPDATE FCM TOKEN =======================
func UpdateFcmTokenHandler(w http.ResponseWriter, r *http.Request) {
    // Validate token and extract user ID
    userID, err := utils.ValidateRestaurantToken(r)
    if err != nil {
        utils.SendResponse(w, http.StatusUnauthorized, false, map[string]interface{}{}, "Unauthorized or invalid token")
        return
    }

    // Parse request body
    var req struct {
        FcmToken string `json:"fcm_token"`
    }
    err = json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "Invalid request body")
        return
    }

    // Validate input
    if req.FcmToken == "" {
        utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "FCM token is required")
        return
    }

    // Update FCM token in DB
    err = restorantservices.UpdateFcmToken(userID, req.FcmToken)
    if err != nil {
        utils.SendResponse(w, http.StatusInternalServerError, false, map[string]interface{}{}, "Failed to update FCM token")
        return
    }

    utils.SendResponse(w, http.StatusOK, true, map[string]interface{}{
        "restorant_id":   userID,
        "fcm_token": req.FcmToken,
    }, "FCM token updated successfully")
}




func UpdateRestaurantProfile(w http.ResponseWriter, r *http.Request) {
	// Validate token
	restaurantID, err := utils.ValidateRestaurantToken(r)
	if err != nil {
		utils.SendResponse(w, http.StatusUnauthorized, false, nil, "Unauthorized or invalid token")
		return
	}

	// Parse multipart form data
	err = r.ParseMultipartForm(20 << 20) // 20MB max
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, false, nil, "Failed to parse form data")
		return
	}

	// Parse restaurant JSON field
	var updatedData restorantmodels.RestaurantProfile
	err = json.Unmarshal([]byte(r.FormValue("restaurant")), &updatedData)
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid restaurant data")
		return
	}
	updatedData.RestaurantID = restaurantID

	// Handle uploaded images (optional)
	var imagePaths []string
	files := r.MultipartForm.File["images"]
	for _, fileHeader := range files {
		path, err := utils.SaveFile(fileHeader, "uploads/restaurant_images")
		if err != nil {
			utils.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to save image")
			return
		}
		imagePaths = append(imagePaths, path)
	}

	// Call service layer
	err = restorantservices.UpdateRestaurantProfileService(updatedData, imagePaths)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to update restaurant profile: "+err.Error())
		return
	}

	// Success
	utils.SendResponse(w, http.StatusOK, true, map[string]interface{}{
		"restaurant_id": updatedData.RestaurantID,
		"name":          updatedData.Name,
		"email":         updatedData.Email,
		"phone_number":  updatedData.PhoneNumber,
		"address":       updatedData.Address,
		"city":          updatedData.City,
		"state":         updatedData.State,
		"country":       updatedData.Country,
		"postal_code":   updatedData.PostalCode,
		"latitude":      updatedData.Latitude,
		"longitude":     updatedData.Longitude,
		"images":        imagePaths,
	}, "Restaurant profile and images updated successfully")
}




func SentOtp(w http.ResponseWriter, r *http.Request) {
	var req struct {
		MobileNumber string `json:"mobile_no"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "Invalid request body")
		return
	}

	if req.MobileNumber == "" {
		utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "Mobile number is required")
		return
	}

	if len(req.MobileNumber) != 10 {
		utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "Mobile number must be 10 digits")
		return
	}

	for _, c := range req.MobileNumber {
		if c < '0' || c > '9' {
			utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "Mobile number must contain only digits")
			return
		}
	}

	otp, err := restorantservices.SendOtp(req.MobileNumber)
	if err != nil {
		// NEVER return nil → return empty map
		utils.SendResponse(w, http.StatusUnauthorized, false, map[string]interface{}{}, err.Error())
		return
	}

	resp := map[string]interface{}{
		"otp": otp,
	}

	utils.SendResponse(w, http.StatusOK, true, resp, "OTP sent successfully.")
}


func VerifyOtp(w http.ResponseWriter, r *http.Request) {

	var req struct {
		Mobile string `json:"mobile_no"`
		Otp    string `json:"otp"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "Invalid request body")
		return
	}

	if req.Mobile == "" || req.Otp == "" {
		utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "Mobile number and OTP are required")
		return
	}

	data, err := restorantservices.VerifyOtp(req.Mobile, req.Otp)
	if err != nil {
		utils.SendResponse(w, http.StatusUnauthorized, false, map[string]interface{}{}, err.Error())
		return
	}

	utils.SendResponse(w, http.StatusOK, true, data, "Login successful")
}
