package userrepo

import (
	"decoration_project/config"
	usermodels "decoration_project/models/user_models"

	"github.com/google/uuid"
)

func AddUser(user usermodels.User) (string, error) {
    newUUID := uuid.New().String()
    query := "INSERT INTO Users (user_id, full_name, email, password, mobile_number) VALUES (?, ?, ?, ?, ?)"
    _, err := config.DB.Exec(query, newUUID, user.FullName, user.Email, user.Password, user.MobileNumber)
    if err != nil {
        return "", err
    }
    return newUUID, nil
}

func GetUserByEmail(email string) (usermodels.User, error) {
    query := "SELECT user_id, full_name, email, password, mobile_number FROM Users WHERE email = ?"
    row := config.DB.QueryRow(query, email)

    var user usermodels.User
    err := row.Scan(&user.UserID, &user.FullName, &user.Email, &user.Password, &user.MobileNumber)
    if err != nil {
        return usermodels.User{}, err
    }
    return user, nil
}


func UpdateFcmToken(userID, fcmToken string) error {
	query := "UPDATE Users SET fcm_token = ? WHERE user_id = ?"
	_, err := config.DB.Exec(query, fcmToken, userID)
	return err
}

func GetUserFcmToken(userID string) (string, error) {
	query := "SELECT fcm_token FROM Users WHERE user_id = ?"
	row := config.DB.QueryRow(query, userID)

	var fcmToken string
	err := row.Scan(&fcmToken)
	if err != nil {
		return "", err
	}
	return fcmToken, nil
}