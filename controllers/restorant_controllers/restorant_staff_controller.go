package restorantcontrollers

import (
	restorantmodels "decoration_project/models/restorant_models"
	restorantservices "decoration_project/services/restorant_services"
	"decoration_project/utils"
	"fmt"
	"net/http"
)


func AddStaff(w http.ResponseWriter, r *http.Request) {
	// Validate token
	restaurantID, err := utils.ValidateRestaurantToken(r)
	if err != nil {
		utils.SendResponse(w, http.StatusUnauthorized, false, nil, err.Error())
		return
	}

	// Parse multipart form data (max 5 MB for safety)
	err = r.ParseMultipartForm(1 << 20)
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, false, nil, "Failed to parse form data. Max 5 MB allowed")
		return
	}

	// Build staff model from form fields
	staff := restorantmodels.RestaurantStaff{
		RestaurantID: restaurantID,
		Name:         r.FormValue("name"),
		Email:        r.FormValue("email"),
		WhatsappNo:   r.FormValue("whatsapp_no"),
		Designation:  r.FormValue("designation"),
		Description:  r.FormValue("description"),
		Password:     r.FormValue("password"),
	}

	// Validate required fields
	if staff.Name == "" || staff.Email == "" || staff.Password == "" || staff.Designation == "" {
		utils.SendResponse(w, http.StatusBadRequest, false, nil, "Name, Email, Password, and Designation are required")
		return
	}

	// Validate designation
	if staff.Designation != "Decorator" && staff.Designation != "Staff" {
		utils.SendResponse(w, http.StatusBadRequest, false, nil, "Designation must be either 'Decorator' or 'Staff'")
		return
	}

	// Handle image upload (required)
	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		utils.SendResponse(w, http.StatusBadRequest, false, nil, "Staff image is required")
		return
	}
	defer file.Close()

	imagePath, err := utils.SaveFile(fileHeader, "uploads/staff")
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to save image")
		return
	}
	staff.ImageURL = imagePath

	// Save staff
	staffID, err := restorantservices.AddStaffService(staff)
	if err != nil {
		fmt.Println("Error adding staff:", err.Error())
		utils.SendResponse(w, http.StatusInternalServerError, false, nil, err.Error())
		return
	}

	// Prepare response
	responseData := map[string]interface{}{
		"staff_id":    staffID,
		
	}

	utils.SendResponse(w, http.StatusCreated, true, responseData, "Staff created successfully")
}


func GetAllStaff(w http.ResponseWriter, r *http.Request) {
    // Validate restaurant token
    restaurantID, err := utils.ValidateRestaurantToken(r)
    if err != nil {
        utils.SendResponse(w, http.StatusUnauthorized, false, nil, err.Error())
        return
    }

    // Fetch staff from service
    staffList, err := restorantservices.GetStaffByRestaurant(restaurantID)
    if err != nil {
        fmt.Println("Error fetching staff:", err.Error())
        utils.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to fetch staff")
        return
    }

    // If no staff found, return empty array instead of null
    if staffList == nil {
        staffList = []restorantmodels.RestaurantStaff{}
    }

    // Prepare response
    responseData := map[string]interface{}{
        "staff": staffList,
    }

    utils.SendResponse(w, http.StatusOK, true, responseData, "Staff fetched successfully")
}


