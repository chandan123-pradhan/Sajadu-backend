package staffcontrollers

import (
	staffmodel "decoration_project/models/staff_model"
	
	staffservices "decoration_project/services/staff_services"
	"decoration_project/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func GetAllAssignedBookings(w http.ResponseWriter, r *http.Request) {
	// Validate staff token and get staff ID
	staffID, err := utils.ValidateStaffToken(r)
	if err != nil {
		utils.SendResponse(w, http.StatusUnauthorized, false, map[string]interface{}{}, "Unauthorized: "+err.Error())
		return
	}
	fmt.Println(staffID)
	// Call service layer to fetch assigned bookings
	bookings, err := staffservices.GetAssignedServices(staffID)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, false, map[string]interface{}{}, "Failed to fetch bookings: "+err.Error())
		return
	}

	// Success response
	utils.SendResponse(w, http.StatusOK, true, bookings, "Assigned bookings fetched successfully")
}

func GetAssignedServicesDetails(w http.ResponseWriter, r *http.Request) {
	// Validate restaurant token
	_, err := utils.ValidateStaffToken(r)
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
	booking, err := staffservices.GetAssignedServiesDetails(bookingID)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, false, map[string]interface{}{}, "Failed to fetch booking details: "+err.Error())
		return
	}

	// Success response
	utils.SendResponse(w, http.StatusOK, true, booking, "Booking details fetched successfully")
}

func StartService(w http.ResponseWriter, r *http.Request) {
	// Validate staff token
	staffId, err := utils.ValidateStaffToken(r)
	if err != nil {
		utils.SendResponse(w, http.StatusUnauthorized, false, nil, "Unauthorized: "+err.Error())
		return
	}

	// Parse request body
	var req staffmodel.VerifyOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid request body: "+err.Error())
		return
	}

	if req.BookingID == "" {
		utils.SendResponse(w, http.StatusBadRequest, false, nil, "BookingID is required")
		return
	}
// Read OTP encryption key from environment
	key := os.Getenv("OTP_ENCRYPTION_KEY")
	fmt.Println("OTP_ENCRYPTION_KEY:", key, "Length:", len(key))
	if len(key) != 32 {
		utils.SendResponse(w, http.StatusInternalServerError, false, nil, "Invalid OTP encryption key configuration")
		return
	}
	
	// Call service layer to verify OTP
	err = staffservices.StartService(req.BookingID, staffId, key)
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "Error"+err.Error())
		return
	}
	// Success response
	utils.SendResponse(w, http.StatusOK, true, nil, "Service has been started successfully, booking status updated to In Progress")
}


func UpdateStaffLocation(w http.ResponseWriter, r *http.Request) {
    // Validate staff token
    staffID, err := utils.ValidateStaffToken(r)
    if err != nil {
        utils.SendResponse(w, http.StatusUnauthorized, false, map[string]interface{}{}, "Unauthorized: "+err.Error())
        return
    }

    // Parse request body
    var req staffmodel.UpdateLocationRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "Invalid request body: "+err.Error())
        return
    }

    if req.BookingID == "" || req.Latitude == 0 || req.Longitude == 0 {
        utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "BookingID, Latitude, and Longitude are required")
        return
    }

    // Call service layer
    err = staffservices.UpdateStaffLocation(staffID, req.BookingID, req.Latitude, req.Longitude)
    if err != nil {
        utils.SendResponse(w, http.StatusInternalServerError, false, map[string]interface{}{}, "Failed to update location: "+err.Error())
        return
    }

    // Success response
    utils.SendResponse(w, http.StatusOK, true, map[string]interface{}{}, "Location updated successfully")
}






// ✅ Verify OTP and mark booking as Completed
func VerifyCompletionOTP(w http.ResponseWriter, r *http.Request) {
	// Validate staff token
	_, err := utils.ValidateStaffToken(r)
	if err != nil {
		utils.SendResponse(w, http.StatusUnauthorized, false, nil, "Unauthorized: "+err.Error())
		return
	}

	// Parse request body
	var req staffmodel.VerifyOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid request body: "+err.Error())
		return
	}

	if req.BookingID == "" || req.OTP == "" {
		utils.SendResponse(w, http.StatusBadRequest, false, nil, "BookingID and OTP are required")
		return
	}

	// Read OTP encryption key from environment
	key := os.Getenv("OTP_ENCRYPTION_KEY")
	if len(key) != 32 {
		utils.SendResponse(w, http.StatusInternalServerError, false, nil, "Invalid OTP encryption key configuration")
		return
	}

	// Call service layer → verify OTP & complete booking
	err = staffservices.VerifyOtpToCompleteService(req.BookingID, req.OTP, key)
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "OTP verification failed: "+err.Error())
		return
	}

	// ✅ Success response
	utils.SendResponse(w, http.StatusOK, true, nil, "OTP verified successfully, booking status updated to Completed")
}