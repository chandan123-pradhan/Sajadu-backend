package admincontroller

import (
	adminserices "decoration_project/services/admin_serices"
	"decoration_project/utils"
	"encoding/json"
	"net/http"
)

// ======================= SET ADMIN FCM TOKEN =======================
func SetAdminFCMToken(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Token string `json:"token"`
	}

	// Parse JSON body
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil || requestBody.Token == "" {
		utils.SendResponse(w, http.StatusBadRequest, false, nil, "Invalid or missing FCM token")
		return
	}

	// Save or update token in DB
	err = adminserices.SaveOrUpdateAdminFCMToken(requestBody.Token)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to update FCM token: "+err.Error())
		return
	}

	utils.SendResponse(w, http.StatusOK, true, nil, "Admin FCM token updated successfully")
}

// ======================= GET ADMIN FCM TOKEN =======================
func GetAdminFCMToken(w http.ResponseWriter, r *http.Request) {
	token, err := adminserices.GetAdminFCMToken()
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to fetch FCM token: "+err.Error())
		return
	}

	data := map[string]interface{}{
		"fcm_token": token,
	}

	utils.SendResponse(w, http.StatusOK, true, data, "Admin FCM token fetched successfully")
}
