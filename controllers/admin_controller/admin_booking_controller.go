package admincontroller

import (
	adminmodel "decoration_project/models/admin_model"
	adminserices "decoration_project/services/admin_serices"
	"decoration_project/utils"
	"encoding/json"
	"fmt"
	"net/http"
)


func GetRestorants(w http.ResponseWriter, r *http.Request) {
	restorants, err := adminserices.GetAllRestorant()
	if err != nil {
		fmt.Print(err)
		utils.SendResponse(w, http.StatusInternalServerError, false, map[string]interface{}{
			"category": []interface{}{}, // always array, even on error
		}, "Failed to fetch categories")
		return
	}

	if restorants == nil {
		restorants = []adminmodel.RestorantModel{}
	}

	data := map[string]interface{}{
		"restorants": restorants,
	}

	utils.SendResponse(w, http.StatusOK, true, data, "restorants fetched successfully")
}

func GetAllBookings(w http.ResponseWriter, r *http.Request) {
    type RequestBody struct {
        Status []string `json:"status"` // e.g., ["Accepted", "Pending"]
    }

    var reqBody RequestBody
    if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
        utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{
            "bookings": []interface{}{},
        }, "Invalid request body")
        return
    }

    bookings, err := adminserices.GetAllActiveBookings(reqBody.Status)
    if err != nil {
        fmt.Print(err)
        utils.SendResponse(w, http.StatusInternalServerError, false, map[string]interface{}{
            "bookings": []interface{}{},
        }, "Failed to fetch bookings")
        return
    }

    if bookings == nil {
        bookings = []adminmodel.BookingResponse{}
    }

    data := map[string]interface{}{
        "bookings": bookings,
    }

    utils.SendResponse(w, http.StatusOK, true, data, "Bookings fetched successfully")
}






// ======================= GET SERVICE DETAILS =======================
func GetBookingsDetails(w http.ResponseWriter, r *http.Request) {
	

	// Get serviceID from query parameter
	bookingid := r.URL.Query().Get("booking_id")
	if bookingid == "" {
		utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "service_id query parameter is required")
		return
	}

	// Fetch service details from service layer
	bookingDetails, err := adminserices.GetBookingsDetails(bookingid)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, false, map[string]interface{}{}, err.Error())
		return
	}


	utils.SendResponse(w, http.StatusOK, true, bookingDetails, "booking details fetched successfully")
}


func UpdateBookingStatus(w http.ResponseWriter, r *http.Request) {
	var req adminmodel.UpdateBookingStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "Invalid request body")
		return
	}

	if req.BookingID == "" || req.NewStatus == "" {
		utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "booking_id and new_status are required")
		return
	}

	if req.NewStatus == "Accepted" && req.RestaurantID == "" {
		utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "restaurant_id is required when status is Accepted")
		return
	}

	err := adminserices.UpdateBookingStatus(req.BookingID, req.RestaurantID, req.NewStatus)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, false, map[string]interface{}{}, err.Error())
		return
	}

	utils.SendResponse(w, http.StatusOK, true, map[string]interface{}{}, "Booking status updated successfully")
}


