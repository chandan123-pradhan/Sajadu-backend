package restorantrepo

import (
	"crypto/rand"
	"decoration_project/config"
	restorantmodels "decoration_project/models/restorant_models"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Add new restaurant
func AddRestaurant(restaurant restorantmodels.RestaurantProfile, hashedPassword string) (string, error) {
	newUUID := uuid.New().String()
	query := `
    INSERT INTO Restaurant_Profile 
    (restaurant_id, name, email, phone_number, address, city, state, country, latitude, longitude, postalCode, password) 
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`
_, err := config.DB.Exec(query,
    newUUID,
    restaurant.Name,
    restaurant.Email,
    restaurant.PhoneNumber,
    restaurant.Address,
    restaurant.City,
    restaurant.State,
    restaurant.Country,
    restaurant.Latitude,
    restaurant.Longitude,
    restaurant.PostalCode,
    hashedPassword,
)
	if err != nil {
        fmt.Println(err.Error())
		return "", err
	}

	return newUUID, nil
}

// Get restaurant by email (for login)
func GetRestaurantByEmail(email string) (restorantmodels.RestaurantProfile, string, error) {
	query := `
        SELECT restaurant_id, name, email, phone_number, address, city,postalCode, state, country, latitude, longitude, password
        FROM Restaurant_Profile 
        WHERE email = ?
    `
	row := config.DB.QueryRow(query, email)

	var restaurant restorantmodels.RestaurantProfile
	var hashedPassword string

	err := row.Scan(
		&restaurant.RestaurantID,
		&restaurant.Name,
		&restaurant.Email,
		&restaurant.PhoneNumber,
		&restaurant.Address,
		&restaurant.City,
        &restaurant.PostalCode,
		&restaurant.State,
		&restaurant.Country,
		&restaurant.Latitude,
		&restaurant.Longitude,
		&hashedPassword,
	)
	if err != nil {
		return restorantmodels.RestaurantProfile{}, "", err
	}

	return restaurant, hashedPassword, nil
}

func AddRestaurantImages(restaurantID string, imageUrls []string) error {
    query := `
        INSERT INTO Restaurant_Images (image_id, restaurant_id, image_url)
        VALUES (?, ?, ?)
    `
    for _, url := range imageUrls {
        newUUID := uuid.New().String()
        _, err := config.DB.Exec(query, newUUID, restaurantID, url)
        if err != nil {
            return err
        }
    }
    return nil
}


// Get restaurant details with images
func GetRestaurantWithImages(email string) (restorantmodels.RestaurantProfile, []string, string, error) {
    // First fetch restaurant
    restaurant, hashedPassword, err := GetRestaurantByEmail(email)
    if err != nil {
        return restorantmodels.RestaurantProfile{}, nil, "", err
    }

    // Fetch images
    query := `SELECT image_url FROM Restaurant_Images WHERE restaurant_id = ?`
    rows, err := config.DB.Query(query, restaurant.RestaurantID)
    if err != nil {
        return restaurant, nil, hashedPassword, err
    }
    defer rows.Close()

    var images []string
    for rows.Next() {
        var url string
        if err := rows.Scan(&url); err != nil {
            return restaurant, nil, hashedPassword, err
        }
        images = append(images, url)
    }

    return restaurant, images, hashedPassword, nil
}


func UpdateFcmToken(restorantId, fcmToken string) error {
	query := "UPDATE Restaurant_Profile SET fcm_token = ? WHERE restaurant_id = ?"
	_, err := config.DB.Exec(query, fcmToken, restorantId)
	return err
}

func GetRestorantFcmToken(restorantId string) (string, error) {
	query := "SELECT fcm_token FROM Restaurant_Profile WHERE restaurant_id = ?"
	row := config.DB.QueryRow(query, restorantId)

	var fcmToken string
	err := row.Scan(&fcmToken)
	if err != nil {
		return "", err
	}
	return fcmToken, nil
}

func UpdateRestaurantProfile(restaurant restorantmodels.RestaurantProfile) error {
	query := `
		UPDATE Restaurant_Profile 
		SET 
			name = ?, 
			email = ?, 
			phone_number = ?, 
			address = ?, 
			city = ?, 
			state = ?, 
			country = ?, 
			latitude = ?, 
			longitude = ?, 
			postalCode = ?
		WHERE restaurant_id = ?
	`

	_, err := config.DB.Exec(query,
		restaurant.Name,
		restaurant.Email,
		restaurant.PhoneNumber,
		restaurant.Address,
		restaurant.City,
		restaurant.State,
		restaurant.Country,
		restaurant.Latitude,
		restaurant.Longitude,
		restaurant.PostalCode,
		restaurant.RestaurantID,
	)

	if err != nil {
		fmt.Println("Error updating restaurant profile:", err)
		return err
	}

	return nil
}


func UpdateRestaurantImages(restaurantID string, imageUrls []string) error {
    // Delete old images first (optional)
    deleteQuery := `DELETE FROM Restaurant_Images WHERE restaurant_id = ?`
    _, err := config.DB.Exec(deleteQuery, restaurantID)
    if err != nil {
        return err
    }

    // Insert new images
    insertQuery := `
        INSERT INTO Restaurant_Images (image_id, restaurant_id, image_url)
        VALUES (?, ?, ?)
    `
    for _, url := range imageUrls {
        newUUID := uuid.New().String()
        _, err := config.DB.Exec(insertQuery, newUUID, restaurantID, url)
        if err != nil {
            return err
        }
    }

    return nil
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
		INSERT INTO restorant_otp_logs (restorant_id, phone_number, otp_code, expires_at)
		VALUES (?, ?, ?, ?)
	`
	_, err = config.DB.Exec(query, userID, mobileNumber, otp, expiry)
	if err != nil {
		return "", err
	}

	return otp, nil
}

// -----------------------------
//  GET RESTAURANT BY MOBILE FUNCTION
// -----------------------------
func GetRestaurantByMobileNo(mobile string) (restorantmodels.RestaurantProfile, error) {

	query := `
	SELECT 
		restaurant_id,
		name,
		email,
		phone_number,
		address,
		city,
		state,
		country,
		latitude,
		longitude,
		postalCode,
		password
	FROM Restaurant_Profile 
	WHERE phone_number = ?
	LIMIT 1
	`

	row := config.DB.QueryRow(query, mobile)

	var restaurant restorantmodels.RestaurantProfile

	err := row.Scan(
		&restaurant.RestaurantID,
		&restaurant.Name,
		&restaurant.Email,
		&restaurant.PhoneNumber,
		&restaurant.Address,
		&restaurant.City,
		&restaurant.State,
		&restaurant.Country,
		&restaurant.Latitude,
		&restaurant.Longitude,
		&restaurant.PostalCode,
		&restaurant.Password,
	)

	if err != nil {
		fmt.Println(err.Error())
		return restorantmodels.RestaurantProfile{}, err
	}

	return restaurant, nil
}


func VerifyOTP(mobile string, otp string) (string, error) {

	fmt.Println("=== VERIFY OTP DEBUG ===")
	fmt.Println("INPUT Mobile:", mobile)
	fmt.Println("INPUT OTP:", otp)

	// Step 1 — DEBUG: Check latest OTP (optional)
	debugQuery := `
		SELECT id, restorant_id, phone_number, otp_code, expires_at
		FROM restorant_otp_logs
		WHERE phone_number = ?
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
		SELECT id, restorant_id, phone_number, otp_code, expires_at
		FROM restorant_otp_logs
		WHERE phone_number = ?
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
