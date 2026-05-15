package admincontroller

import (
	usermodels "decoration_project/models/user_models"
	adminserices "decoration_project/services/admin_serices"
	"decoration_project/utils"
	"fmt"
	"net/http"
)


func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	
	

	users, err := adminserices.GetAllUsers()
	if err != nil {
		fmt.Print(err)
		utils.SendResponse(w, http.StatusInternalServerError, false, map[string]interface{}{
			"users": []interface{}{},
		}, "Failed to fetch users")
		return
	}

	if users == nil {
		users = []usermodels.UserDetailsModel{}
	}

	data := map[string]interface{}{
		"users": users,
	}

	utils.SendResponse(w, http.StatusOK, true, data, "Users fetched successfully")
}
