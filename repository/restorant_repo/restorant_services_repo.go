package restorantrepo

import (
	"database/sql"
	"decoration_project/config"
	restorantmodels "decoration_project/models/restorant_models"
)

func GetAllProposedServices(restaurantID string) ([]restorantmodels.RestaurantService, error) {
	var services []restorantmodels.RestaurantService

	query := `SELECT 
				service_id, category_id, service_name, service_description, 
				service_price, proposed_restaurant_id, created_at, updated_at
			  FROM Our_Services
			  WHERE proposed_restaurant_id = ?
			  AND is_deleted = 0`   // <-- return only non-deleted services

	rows, err := config.DB.Query(query, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var service restorantmodels.RestaurantService
		var proposedRestaurantID sql.NullString

		err := rows.Scan(
			&service.ServiceID,
			&service.CategoryId,
			&service.ServiceName,
			&service.ServiceDescription,
			&service.ServicePrice,
			&proposedRestaurantID,
			&service.CreatedAt,
			&service.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Set proposed restaurant ID
		if proposedRestaurantID.Valid {
			service.ProposedRestorantId = proposedRestaurantID.String
		} else {
			service.ProposedRestorantId = ""
		}

		// Fetch images
		imgQuery := `SELECT image_url FROM Service_Images WHERE service_id = ?`
		imgRows, err := config.DB.Query(imgQuery, service.ServiceID)
		if err != nil {
			return nil, err
		}

		var images []string
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



func CheckServiceOwner(serviceID string, restaurantID string) (bool, error) {
	var count int

	query := `SELECT COUNT(*) 
              FROM Our_Services 
              WHERE service_id = ? 
              AND proposed_restaurant_id = ?
              AND is_deleted = 0`

	err := config.DB.QueryRow(query, serviceID, restaurantID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}


func DeleteServiceWithImages(serviceID string) error {
	tx, err := config.DB.Begin()
	if err != nil {
		return err
	}

	// Delete images (optional)
	_, err = tx.Exec(`DELETE FROM Service_Images WHERE service_id = ?`, serviceID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Soft delete service
	_, err = tx.Exec(`UPDATE Our_Services SET is_deleted = 1 WHERE service_id = ?`, serviceID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
