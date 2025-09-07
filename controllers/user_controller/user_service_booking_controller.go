package usercontroller

import (
	usermodels "decoration_project/models/user_models"
	userservices "decoration_project/services/user_services"
	"decoration_project/utils"
	"encoding/json"
	"net/http"
	"os"
)



func CreateBooking(w http.ResponseWriter, r *http.Request) {
	// Validate JWT token and get userId
	userId, err := utils.ValidateToken(r)
	if err != nil {
		utils.SendResponse(w, http.StatusUnauthorized, false, nil, "Unauthorized: "+err.Error())
		return
	}

	// Decode request body
	var bookingReq usermodels.BookingRequest
	if err := json.NewDecoder(r.Body).Decode(&bookingReq); err != nil {
		utils.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid request body: "+err.Error())
		return
	}

	// Override UserID from token (don’t trust client input)
	bookingReq.UserID = userId

	// ✅ Centralized validation
	if err := utils.ValidateBookingRequest(bookingReq); err != nil {
		utils.SendResponse(w, http.StatusBadRequest, false, nil, err.Error())
		return
	}

	// Call service layer
	bookingRes, err := userservices.CreateBooking(bookingReq)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to create booking: "+err.Error())
		return
	}

	// Success
	utils.SendResponse(w, http.StatusOK, true, bookingRes, "Booking created successfully")
}


func GetUsersBookings(w http.ResponseWriter, r *http.Request) {
	// Validate JWT token and get userId
	userId, err := utils.ValidateToken(r)
	if err != nil {
		utils.SendResponse(w, http.StatusUnauthorized, false, map[string]interface{}{}, "Unauthorized: "+err.Error())
		return
	}

	// Read OTP encryption key from environment
	key := os.Getenv("OTP_ENCRYPTION_KEY")
	if len(key) != 32 {
		utils.SendResponse(w, http.StatusInternalServerError, false, nil, "Invalid OTP encryption key configuration")
		return
	}
	// Call service layer to get bookings
	bookings, err := userservices.FetchBookings(userId,key)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, false, map[string]interface{}{}, "Failed to fetch bookings: "+err.Error())
		return
	}

	// Success response
	utils.SendResponse(w, http.StatusOK, true, bookings, "Bookings fetched successfully")
}




func GetServiceDetails(w http.ResponseWriter, r *http.Request) {
	// Validate JWT token
	_, err := utils.ValidateToken(r)
	if err != nil {
		utils.SendResponse(w, http.StatusUnauthorized, false, map[string]interface{}{}, "Unauthorized: "+err.Error())
		return
	}

	// Decode request body
	err = json.NewDecoder(r.Body).Decode(&usermodels.GetServiceDetailsRequest)
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "Invalid request body")
		return
	}

	// Validate required fields
	if usermodels.GetServiceDetailsRequest.RestaurantId == "" || usermodels.GetServiceDetailsRequest.ServiceId == "" {
		utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "Restaurant ID and Service ID are required")
		return
	}

	// Call service layer
	serviceWithRest, err := userservices.GetServiceDetails(
		usermodels.GetServiceDetailsRequest.RestaurantId,
		usermodels.GetServiceDetailsRequest.ServiceId,
	)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, false, map[string]interface{}{}, "Failed to fetch details: "+err.Error())
		return
	}

	// If not found
	if serviceWithRest.Service.ServiceID == "" {
		utils.SendResponse(w, http.StatusNotFound, false, map[string]interface{}{}, "No service found for given Restaurant ID and Service ID")
		return
	}

	// Success response
	utils.SendResponse(w, http.StatusOK, true, map[string]interface{}{
		"service":    serviceWithRest.Service,
		"restaurant": serviceWithRest.Restaurant,
	}, "Service details fetched successfully")
}