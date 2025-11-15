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


func UpdateFcmToken(userID, fcmToken string) error {
    // Ensure userID and fcmToken are valid
    if userID == "" || fcmToken == "" {
        return errors.New("invalid user ID or FCM token")
    }

    // Call repository layer
    err := userrepo.UpdateFcmToken(userID, fcmToken)
    if err != nil {
        return err
    }

    return nil
}

func GetUserFcmTokenService(userID string) (string, error) {
	return userrepo.GetUserFcmToken(userID)
}




func CancelBookingByUserService(userID, bookingID, reason string) error {
	// Validate input
	if userID == "" || bookingID == "" {
		return errors.New("userID and bookingID are required")
	}

	// Call repository function
	err := userrepo.CancelBookingByUser(bookingID, userID, reason)
	if err != nil {
		return err
	}

	return nil
}



func SendOtp(mobileNo string) (string, error) {

    // 1. Check user exists
    user, err := userrepo.GetUserByMobile(mobileNo)
    if err != nil {
        return "", errors.New("mobile number not registered")
    }

    // 2. Now user exists → pass userID to SendOTP
    userID := &user.UserID

    otp, err := userrepo.SendOTP(mobileNo, userID)
    if err != nil {
        return "", err
    }

    return otp, nil
}



func VerifyOtp(mobile string, otp string) (map[string]interface{}, error) {

    // Step 1: Verify OTP
    userID, err := userrepo.VerifyOTP(mobile, otp)
    if err != nil {
        return nil, err
    }

    // Step 2: Get User Details
    user, err := userrepo.GetUserByID(userID)
    if err != nil {
        return nil, errors.New("User not found")
    }

    // Step 3: Generate JWT Token
    token, err := utils.GenerateToken(user.UserID)
    if err != nil {
        return nil, errors.New("Failed to generate token")
    }

    // Step 4: Final Response
    data := map[string]interface{}{
        "auth_token":    token,
        "user_id":       user.UserID,
        "full_name":     user.FullName,
        "email":         user.Email,
        "mobile_number": user.MobileNumber,
    }

    return data, nil
}
