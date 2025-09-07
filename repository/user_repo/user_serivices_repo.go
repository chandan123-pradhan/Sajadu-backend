package userrepo

import (
	"database/sql"
	"decoration_project/config"
	restorantmodels "decoration_project/models/restorant_models"
)

func GetServicesByCategory(categoryID string) ([]restorantmodels.RestaurantService, error) {
	// Fetch services for the category
	query := `
		SELECT 
			service_id, restaurant_id, category_id, service_name, service_description, 
			service_price, items, created_at, updated_at
		FROM Restaurant_Services
		WHERE category_id = ?
	`

	rows, err := config.DB.Query(query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []restorantmodels.RestaurantService

	for rows.Next() {
		var service restorantmodels.RestaurantService
		var items sql.NullString

		if err := rows.Scan(
			&service.ServiceID,
			&service.RestaurantID,
			&service.CategoryId,
			&service.ServiceName,
			&service.ServiceDescription,
			&service.ServicePrice,
			&items,
			&service.CreatedAt,
			&service.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if items.Valid {
			service.Items = items.String
		} else {
			service.Items = ""
		}

		// Fetch images for this service
		imageRows, err := config.DB.Query("SELECT image_url FROM Service_Images WHERE service_id = ?", service.ServiceID)
		if err != nil {
			return nil, err
		}

		var images []string
		for imageRows.Next() {
			var img string
			if err := imageRows.Scan(&img); err != nil {
				imageRows.Close()
				return nil, err
			}
			images = append(images, img)
		}
		imageRows.Close()
		service.Images = images

		services = append(services, service)
	}

	return services, nil
}


func GetServiceDetails(restaurantID, serviceID string) (restorantmodels.ServiceWithRestaurant, error) {
	var result restorantmodels.ServiceWithRestaurant

	query := `
		SELECT 
			rs.service_id, rs.restaurant_id, rs.category_id, rs.service_name, rs.service_description, rs.service_price, rs.items, rs.created_at, rs.updated_at,
			r.restaurant_id, r.name, r.email, r.phone_number, r.address, r.city, r.state, r.country, r.latitude, r.longitude, r.postalCode
		FROM Restaurant_Services rs
		INNER JOIN Restaurant_Profile r ON rs.restaurant_id = r.restaurant_id
		WHERE rs.restaurant_id = ? AND rs.service_id = ?
	`

	row := config.DB.QueryRow(query, restaurantID, serviceID)

	var items sql.NullString
	err := row.Scan(
		&result.Service.ServiceID,
		&result.Service.RestaurantID,
		&result.Service.CategoryId,
		&result.Service.ServiceName,
		&result.Service.ServiceDescription,
		&result.Service.ServicePrice,
		&items,
		&result.Service.CreatedAt,
		&result.Service.UpdatedAt,
		&result.Restaurant.RestaurantID,
		&result.Restaurant.Name,
		&result.Restaurant.Email,
		&result.Restaurant.PhoneNumber,
		&result.Restaurant.Address,
		&result.Restaurant.City,
		&result.Restaurant.State,
		&result.Restaurant.Country,
		&result.Restaurant.Latitude,
		&result.Restaurant.Longitude,
		&result.Restaurant.PostalCode,
	)

	if err != nil {
		return result, err
	}

	// Handle JSON items column
	if items.Valid {
		result.Service.Items = items.String
	} else {
		result.Service.Items = ""
	}

	// Fetch service images
	imageRows, err := config.DB.Query("SELECT image_url FROM Service_Images WHERE service_id = ?", result.Service.ServiceID)
	if err != nil {
		return result, err
	}
	defer imageRows.Close()

	var images []string
	for imageRows.Next() {
		var img string
		if err := imageRows.Scan(&img); err != nil {
			return result, err
		}
		images = append(images, img)
	}
	result.Service.Images = images


	// ✅ Fetch restaurant images (not service images)
restorantImageRow, err := config.DB.Query("SELECT image_url FROM Restaurant_Images WHERE restaurant_id = ?", result.Restaurant.RestaurantID)
if err != nil {
    return result, err
}
defer restorantImageRow.Close()

var restorantImage []string
for restorantImageRow.Next() { // 🔥 fixed here
    var img string
    if err := restorantImageRow.Scan(&img); err != nil { // 🔥 fixed here
        return result, err
    }
    restorantImage = append(restorantImage, img)
}

result.Restaurant.Images = restorantImage


	return result, nil
}