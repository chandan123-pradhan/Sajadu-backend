package userrepo

import (
	"decoration_project/config"
	restorantmodels "decoration_project/models/restorant_models"
)

func GetServicesByCategory(categoryID string) ([]restorantmodels.RestaurantService, error) {
	// Fetch services for the category
	query := `
		SELECT 
			service_id, category_id, service_name, service_description, 
			service_price, created_at, updated_at
		FROM Our_Services
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
		

		if err := rows.Scan(
			&service.ServiceID,
			&service.CategoryId,
			&service.ServiceName,
			&service.ServiceDescription,
			&service.ServicePrice,
			&service.CreatedAt,
			&service.UpdatedAt,
		); err != nil {
			return nil, err
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


func GetServiceDetails(serviceID string) (restorantmodels.ServiceWithRestaurant, error) {
	var result restorantmodels.ServiceWithRestaurant

	// Fetch only service details
	query := `
		SELECT 
			service_id, category_id, service_name, service_description, service_price, created_at, updated_at
		FROM Our_Services
		WHERE service_id = ?
	`

	row := config.DB.QueryRow(query, serviceID)
	err := row.Scan(
		&result.Service.ServiceID,
		&result.Service.CategoryId,
		&result.Service.ServiceName,
		&result.Service.ServiceDescription,
		&result.Service.ServicePrice,
		&result.Service.CreatedAt,
		&result.Service.UpdatedAt,
	)
	if err != nil {
		return result, err
	}

	// Fetch service images
	imageRows, err := config.DB.Query("SELECT image_url FROM Service_Images WHERE service_id = ?", result.Service.ServiceID)
	if err != nil {
		return result, err
	}
	defer imageRows.Close()

	images := []string{} // empty slice to avoid null
	for imageRows.Next() {
		var img string
		if err := imageRows.Scan(&img); err != nil {
			return result, err
		}
		images = append(images, img)
	}
	result.Service.Images = images

	return result, nil
}
