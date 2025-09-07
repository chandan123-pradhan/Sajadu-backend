package staffcontrollers

import (
	staffservices "decoration_project/services/staff_services"
	"decoration_project/utils"
	"encoding/json"
	"net/http"
)


func LoginStaffHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "Invalid request body")
        return
    }

    if req.Email == "" || req.Password == "" {
        utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "Email and password are required")
        return
    }

    user, err := staffservices.LoginStaff(req.Email, req.Password)
    if err != nil {
        utils.SendResponse(w, http.StatusUnauthorized, false, map[string]interface{}{}, err.Error())
        return
    }

    utils.SendResponse(w, http.StatusOK, true, user, "Login successful")
}