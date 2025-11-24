package adminrepo

import (
	"database/sql"
	"decoration_project/config"
	restorantmodels "decoration_project/models/restorant_models"

	"github.com/google/uuid"
)

// AddService inserts a new service including CategoryId
func AddService(service restorantmodels.RestaurantService) (string, error) {
	serviceID := uuid.New().String()
	query := `
        INSERT INTO Our_Services 
        (service_id, category_id, service_name, service_description, service_price, proposed_restaurant_id)
        VALUES (?, ?, ?, ?, ?, ?)
    `
	_, err := config.DB.Exec(query,
		serviceID,
		service.CategoryId,
		service.ServiceName,
		service.ServiceDescription,
		service.ServicePrice,
		service.ProposedRestorantId,
	)
	if err != nil {
		return "", err
	}
	return serviceID, nil
}

// AddServiceImages inserts multiple images for a service
func AddServiceImages(serviceID string, images []string) error {
	query := `INSERT INTO Service_Images (image_id, service_id, image_url) VALUES (?, ?, ?)`
	for _, img := range images {
		_, err := config.DB.Exec(query, uuid.New().String(), serviceID, img)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetServiceWithImages fetches a single service along with its images
func GetServiceWithImages(serviceID string) (restorantmodels.RestaurantService, []string, error) {
	var service restorantmodels.RestaurantService
	var proposedRestaurantID sql.NullString
	query := `
        SELECT service_id, category_id, service_name, service_description, service_price, proposed_restaurant_id
        FROM Our_Services
        WHERE service_id = ?
    `
	row := config.DB.QueryRow(query, serviceID)
	err := row.Scan(&service.ServiceID, &service.CategoryId,
		&service.ServiceName, &service.ServiceDescription, &service.ServicePrice, &proposedRestaurantID)
	if err != nil {
		return service, nil, err
	}
	if proposedRestaurantID.Valid {
		service.ProposedRestorantId = proposedRestaurantID.String
	} else {
		service.ProposedRestorantId = "" // default empty string
	}

	// Fetch images
	rows, err := config.DB.Query(`SELECT image_url FROM Service_Images WHERE service_id = ?`, serviceID)
	if err != nil {
		return service, nil, err
	}
	defer rows.Close()

	var images []string
	for rows.Next() {
		var img string
		if err := rows.Scan(&img); err != nil {
			return service, nil, err
		}
		images = append(images, img)
	}

	return service, images, nil
}

// GetAllServicesWithImages fetches all services for a restaurant with images
func GetAllServiceCategoryWise(categoryId string) ([]restorantmodels.RestaurantService, error) {
	var services []restorantmodels.RestaurantService
	var proposedRestaurantID sql.NullString

	query := `SELECT service_id, category_id, service_name, service_description, service_price, proposed_restaurant_id, created_at, updated_at
              FROM Our_Services
              WHERE category_id = ?`
	rows, err := config.DB.Query(query, categoryId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var service restorantmodels.RestaurantService
		if err := rows.Scan(
			&service.ServiceID,
			&service.CategoryId,
			&service.ServiceName,
			&service.ServiceDescription,
			&service.ServicePrice,
			&proposedRestaurantID,
			&service.CreatedAt,
			&service.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if proposedRestaurantID.Valid {
			service.ProposedRestorantId = proposedRestaurantID.String
		} else {
			service.ProposedRestorantId = "" // default empty string
		}

		// Fetch images for this service
		imgRows, err := config.DB.Query(`SELECT image_url FROM Service_Images WHERE service_id = ?`, service.ServiceID)
		if err != nil {
			return nil, err
		}

		images := []string{}
		for imgRows.Next() {
			var url string
			if err := imgRows.Scan(&url); err != nil {
				imgRows.Close()
				return nil, err
			}
			images = append(images, url)
		}
		imgRows.Close()

		service.Images = images
		services = append(services, service)
	}

	return services, nil
}
