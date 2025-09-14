package usercontroller

import (
	usermodels "decoration_project/models/user_models"
	// staffservices "decoration_project/services/staff_services"
	userservices "decoration_project/services/user_services"
	"decoration_project/utils"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gorilla/mux"
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





func GetBookingDetails(w http.ResponseWriter, r *http.Request) {
	// Validate JWT token and get userId
	_, err := utils.ValidateToken(r)
	if err != nil {
		utils.SendResponse(w, http.StatusUnauthorized, false, map[string]interface{}{}, "Unauthorized: "+err.Error())
		return
	}

	// Get booking_id from URL
	vars := mux.Vars(r)
	bookingID := vars["id"]
	if bookingID == "" {
		utils.SendResponse(w, http.StatusBadRequest, false, nil, "booking_id is required")
		return
	}

	// Call service layer
	booking, err := userservices.FetchBookingDetails(bookingID)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, false, map[string]interface{}{}, "Failed to fetch booking details: "+err.Error())
		return
	}

	// Success response
	utils.SendResponse(w, http.StatusOK, true, booking, "Booking details fetched successfully")
	

}


