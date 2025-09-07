package usercontroller

import (
	usermodels "decoration_project/models/user_models"
	userservices "decoration_project/services/user_services"
	"decoration_project/utils"
	"encoding/json"
	"net/http"

	"github.com/go-sql-driver/mysql"
)

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
    var user usermodels.User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "Invalid request body")
        return
    }

    // Validate required fields
    if user.FullName == "" || user.Email == "" || user.Password == "" {
        utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "Full name, email, and password are required")
        return
    }

    // Register user
    userID, err := userservices.RegisterUser(user)
    if err != nil {
        if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
            utils.SendResponse(w, http.StatusBadRequest, false, map[string]interface{}{}, "Email | Mobile No. already registered")
            return
        }
        utils.SendResponse(w, http.StatusInternalServerError, false, map[string]interface{}{}, err.Error())
        return
    }

    // Generate auth token
    token, err := utils.GenerateToken(userID)
    if err != nil {
        utils.SendResponse(w, http.StatusInternalServerError, false, map[string]interface{}{}, "Failed to generate auth token")
        return
    }

    // Prepare response data (exclude password)
    responseData := map[string]interface{}{
        "user_id":       userID,
        "full_name":     user.FullName,
        "email":         user.Email,
        "mobile_number": user.MobileNumber,
        "auth_token":    token,
    }

    utils.SendResponse(w, http.StatusCreated, true, responseData, "User registered successfully")
}


func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
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

    user, err := userservices.LoginUser(req.Email, req.Password)
    if err != nil {
        utils.SendResponse(w, http.StatusUnauthorized, false, map[string]interface{}{}, err.Error())
        return
    }

    utils.SendResponse(w, http.StatusOK, true, user, "Login successful")
}