package restorantcontrollers

import (
	restorantservices "decoration_project/services/restorant_services"
	"decoration_project/utils"
	"net/http"
)

func GetRestorantNotifications(w http.ResponseWriter, r *http.Request) {
    // Step 1: Validate JWT token to get userId
    restorantId, err := utils.ValidateRestaurantToken(r)
    if err != nil {
        utils.SendResponse(w, http.StatusUnauthorized, false, nil, "Unauthorized: "+err.Error())
        return
    }

    // Step 2: Call service layer to fetch notifications
    notifications, err := restorantservices.GetRestorantNotificationsService(restorantId)
    if err != nil {
        utils.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to fetch notifications: "+err.Error())
        return
    }

    // Step 3: Success response
    utils.SendResponse(w, http.StatusOK, true, notifications, "Notifications fetched successfully")
}
