package restorantcontrollers

import (
	restorantmodels "decoration_project/models/restorant_models"
	restorantservices "decoration_project/services/restorant_services"
	"decoration_project/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func GetAllBookedServices(w http.ResponseWriter, r *http.Request) {
	// Validate restaurant token and get restaurant ID
	restaurantID, err := utils.ValidateRestaurantToken(r)
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
	bookings, err := restorantservices.FetchRestorantBookings(restaurantID,key)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, false, map[string]interface{}{}, "Failed to fetch bookings: "+err.Error())
		return
	}

	// Success response
	utils.SendResponse(w, http.StatusOK, true, bookings, "Bookings fetched successfully")

}






func GetBookingDetails(w http.ResponseWriter, r *http.Request) {
	// Validate restaurant token
	_, err := utils.ValidateRestaurantToken(r)
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

	// Read OTP encryption key
	key := os.Getenv("OTP_ENCRYPTION_KEY")
	if len(key) != 32 {
		utils.SendResponse(w, http.StatusInternalServerError, false, nil, "Invalid OTP encryption key configuration")
		return
	}

	// Call service layer
	booking, err := restorantservices.FetchRestorantBookingDetails(bookingID, key)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, false, map[string]interface{}{}, "Failed to fetch booking details: "+err.Error())
		return
	}

	// Success response
	utils.SendResponse(w, http.StatusOK, true, booking, "Booking details fetched successfully")
}



func HandleBookingAction(w http.ResponseWriter, r *http.Request) {
	// Validate restaurant token
	_, err := utils.ValidateRestaurantToken(r)
	if err != nil {
		utils.SendResponse(w, http.StatusUnauthorized, false, nil, "Unauthorized: "+err.Error())
		return
	}

	// Parse request body
	var req restorantmodels.BookingActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid request body: "+err.Error())
		return
	}

	// Read encryption key from environment
	key := os.Getenv("OTP_ENCRYPTION_KEY")
	if len(key) != 32 {
		utils.SendResponse(w, http.StatusInternalServerError, false, nil, "Invalid OTP encryption key configuration")
		return
	}

	// Map status string to status_id and handler function
	var statusID int
	var actionFunc func(string, int, string) (interface{}, error)

	switch req.Status {
	case "accept":
		statusID = utils.GetStatusID("Accepted")
		actionFunc = func(bookingID string, statusID int, key string) (interface{}, error) {
			return restorantservices.AcceptBooking(bookingID, statusID, key)
		}
	case "reject":
		statusID = utils.GetStatusID("Rejected")
		actionFunc = func(bookingID string, statusID int, key string) (interface{}, error) {
			return nil, restorantservices.RejectBooking(bookingID, statusID, req.Reason)
		}
	default:
		utils.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid status value")
		return
	}

	// Execute action
	result, err := actionFunc(req.BookingID, statusID, key)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, false, nil, fmt.Sprintf("Failed to %s booking: %s", req.Status, err.Error()))
		return
	}

	// Send success response
	respData := map[string]interface{}{}
	if req.Status == "accept" {
		respData["otp"] = result // only accept returns OTP
	}
	utils.SendResponse(w, http.StatusOK, true, respData, fmt.Sprintf("Booking %sed successfully", req.Status))
}







func AssignStaffToBooking(w http.ResponseWriter, r *http.Request) {
	// Validate restaurant token
	_, err := utils.ValidateRestaurantToken(r)
	if err != nil {
		utils.SendResponse(w, http.StatusUnauthorized, false, nil, "Unauthorized: "+err.Error())
		return
	}

	// Parse request body
	var req struct {
		BookingID string `json:"booking_id"`
		StaffID   string `json:"staff_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid request body: "+err.Error())
		return
	}

	// Call service to assign staff
	err = restorantservices.AssignStaff(req.BookingID, req.StaffID)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to assign staff: "+err.Error())
		return
	}

	utils.SendResponse(w, http.StatusOK, true, map[string]interface{}{
		"booking_id": req.BookingID,
		"staff_id":   req.StaffID,
	}, "Staff assigned successfully")
}
