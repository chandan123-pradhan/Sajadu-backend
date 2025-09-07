package userservices

import (
	usermodels "decoration_project/models/user_models"
	userrepo "decoration_project/repository/user_repo"
	"decoration_project/utils"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(user usermodels.User) (string, error) {
    // Hash the password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    user.Password = string(hashedPassword)

    // Add user in repository
    return userrepo.AddUser(user)
}


func LoginUser(email, password string) (map[string]interface{}, error) {
    user, err := userrepo.GetUserByEmail(email)
    if err != nil {
        return nil, errors.New("invalid email and password") // user not found
    }

    // Compare password
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
        return nil, errors.New("invalid email and password")
    }

    // Generate auth token
    token, err := utils.GenerateToken(user.UserID)
    if err != nil {
        return nil, err
    }

    // Prepare response (exclude password)
    resp := map[string]interface{}{
        "user_id":       user.UserID,
        "full_name":     user.FullName,
        "email":         user.Email,
        "mobile_number": user.MobileNumber,
        "auth_token":    token,
    }

    return resp, nil
}