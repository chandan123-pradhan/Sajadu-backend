package userrepo

import (
	"crypto/rand"
	"decoration_project/config"
	usermodels "decoration_project/models/user_models"
	"errors"
	"fmt"
	"time"

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

func GetUserByID(userID string) (usermodels.User, error) {
    query := "SELECT user_id, full_name, email, password, mobile_number FROM Users WHERE user_id = ?"
    row := config.DB.QueryRow(query, userID)

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










// -------------------
// Generate 6-digit OTP
// -------------------
func generateOTP() (string, error) {
	b := make([]byte, 3)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	otp := int(b[0])<<16 + int(b[1])<<8 + int(b[2])
	return fmt.Sprintf("%06d", otp%1000000), nil
}


func SendOTP(mobileNumber string, userID *string) (string, error) {

	otp, err := generateOTP()
	if err != nil {
		return "", err
	}

	// FIX: always save in UTC
	expiry := time.Now().UTC().Add(5 * time.Minute)

	query := `
		INSERT INTO otp_logs (user_id, mobile_number, otp_code, expires_at)
		VALUES (?, ?, ?, ?)
	`
	_, err = config.DB.Exec(query, userID, mobileNumber, otp, expiry)
	if err != nil {
		return "", err
	}

	return otp, nil
}


// -----------------------------
//  GET USER BY MOBILE FUNCTION
// -----------------------------
func GetUserByMobile(mobile string) (usermodels.User, error) {

	query := `SELECT user_id, full_name, email, password, mobile_number 
			  FROM Users WHERE mobile_number = ?`

	row := config.DB.QueryRow(query, mobile)

	var user usermodels.User
	err := row.Scan(&user.UserID, &user.FullName, &user.Email, &user.Password, &user.MobileNumber)
	if err != nil {
		return usermodels.User{}, err
	}

	return user, nil
}

func VerifyOTP(mobile string, otp string) (string, error) {

	fmt.Println("=== VERIFY OTP DEBUG ===")
	fmt.Println("INPUT Mobile:", mobile)
	fmt.Println("INPUT OTP:", otp)

	// Step 1 — DEBUG: Check latest OTP (optional)
	debugQuery := `
		SELECT id, user_id, mobile_number, otp_code, expires_at
		FROM otp_logs
		WHERE mobile_number = ?
		ORDER BY id DESC
		LIMIT 1
	`

	var dID int
	var dUserID *string
	var dMobile, dOTP string
	var dExpires time.Time

	err := config.DB.QueryRow(debugQuery, mobile).Scan(
		&dID, &dUserID, &dMobile, &dOTP, &dExpires,
	)

	if err == nil {
		fmt.Printf("DB → id: %d user_id: %v mobile: %s otp: %s expires: %v\n",
			dID, dUserID, dMobile, dOTP, dExpires)
		fmt.Println("Now (UTC):", time.Now().UTC())
	}

	// Step 2 — Main query (Use MySQL UTC_TIMESTAMP)
	query := `
		SELECT id, user_id, mobile_number, otp_code, expires_at
		FROM otp_logs
		WHERE mobile_number = ?
		  AND otp_code = ?
		  AND expires_at > UTC_TIMESTAMP()
		ORDER BY id DESC
		LIMIT 1
	`

	fmt.Println("MAIN QUERY running with:", mobile, otp)

	row := config.DB.QueryRow(query, mobile, otp)

	var id int
	var userID *string
	var mobileDB, otpDB string
	var expires time.Time

	err = row.Scan(
		&id,
		&userID,
		&mobileDB,
		&otpDB,
		&expires,
	)

	if err != nil {
		fmt.Println("MAIN QUERY ERROR:", err)
		return "", errors.New("invalid or expired OTP")
	}

	// OTP matched
	fmt.Println("OTP Verified Successfully!")

	// return userID to login the user
	if userID != nil {
		return *userID, nil
	}

	return "", nil
}
