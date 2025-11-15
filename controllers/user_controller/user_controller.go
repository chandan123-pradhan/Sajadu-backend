package usercontroller

import (
	usermodels "decoration_project/models/user_models"
	userservices "decoration_project/services/user_services"
	"decoration_project/utils"
	"encoding/json"
	"net/http"

	"github.com/go-sql-driver/mysql"
)

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var user usermodels.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "Invalid request body")
		return
	}

	// Validate required fields
	if user.FullName == "" || user.Email == "" || user.Password == "" {
		utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "Full name, email, and password are required")
		return
	}

	// Register user
	userID, err := userservices.RegisterUser(user)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "Email | Mobile No. already registered")
			return
		}
		utils.SendResponse(w, http.StatusInternalServerError, false, map[string]interface{}{}, err.Error())
		return
	}

	// Generate auth token
	token, err := utils.GenerateToken(userID)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, false, map[string]interface{}{}, "Failed to generate auth token")
		return
	}

	// Prepare response data (exclude password)
	responseData := map[string]interface{}{
		"user_id":       userID,
		"full_name":     user.FullName,
		"email":         user.Email,
		"mobile_number": user.MobileNumber,
		"auth_token":    token,
	}

	utils.SendResponse(w, http.StatusCreated, true, responseData, "User registered successfully")
}

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
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

	user, err := userservices.LoginUser(req.Email, req.Password)
	if err != nil {
		utils.SendResponse(w, http.StatusUnauthorized, false, map[string]interface{}{}, err.Error())
		return
	}

	utils.SendResponse(w, http.StatusOK, true, user, "Login successful")
}

// ======================= UPDATE FCM TOKEN =======================
func UpdateFcmTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Validate token and extract user ID
	userID, err := utils.ValidateToken(r)
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
	err = userservices.UpdateFcmToken(userID, req.FcmToken)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, false, map[string]interface{}{}, "Failed to update FCM token")
		return
	}

	utils.SendResponse(w, http.StatusOK, true, map[string]interface{}{
		"user_id":   userID,
		"fcm_token": req.FcmToken,
	}, "FCM token updated successfully")
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

	otp, err := userservices.SendOtp(req.MobileNumber)
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

	data, err := userservices.VerifyOtp(req.Mobile, req.Otp)
	if err != nil {
		utils.SendResponse(w, http.StatusUnauthorized, false, map[string]interface{}{}, err.Error())
		return
	}

	utils.SendResponse(w, http.StatusOK, true, data, "Login successful")
}
