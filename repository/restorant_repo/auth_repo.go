package restorantrepo

import (
	"decoration_project/config"
	restorantmodels "decoration_project/models/restorant_models"
	"fmt"

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