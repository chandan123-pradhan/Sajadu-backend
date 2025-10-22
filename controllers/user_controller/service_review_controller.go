package usercontroller

import (
	"encoding/json"
	"net/http"
	usermodels "decoration_project/models/user_models"
	userservices "decoration_project/services/user_services"
	"decoration_project/utils"
)

// AddServiceReview handles adding a user review for a service
func AddServiceReview(w http.ResponseWriter, r *http.Request) {
	// Step 1: Validate JWT token to get user ID
	userID, err := utils.ValidateToken(r)
	if err != nil {
		utils.SendResponse(w, http.StatusUnauthorized, false, nil, "Unauthorized: "+err.Error())
		return
	}

	// Step 2: Parse incoming JSON body
	var req usermodels.ServiceReview
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid request body: "+err.Error())
		return
	}

	// Step 3: Validate required fields
	if req.ServiceID == "" {
		utils.SendResponse(w, http.StatusBadRequest, false, nil, "Service ID is required")
		return
	}
	if req.Rating < 1 || req.Rating > 5 {
		utils.SendResponse(w, http.StatusBadRequest, false, nil, "Rating must be between 1 and 5")
		return
	}
	if req.UserName == "" {
		utils.SendResponse(w, http.StatusBadRequest, false, nil, "User name is required")
		return
	}

	// Step 4: Add userId from JWT to the request (for security)
	req.UserID = userID

	// Step 5: Call service layer
	err = userservices.AddServiceReview(req)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to add review: "+err.Error())
		return
	}

	// Step 6: Send success response
	utils.SendResponse(w, http.StatusCreated, true, nil, "Review added successfully")
}
