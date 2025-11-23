package userrepo

import (
	"decoration_project/config"
	restorantmodels "decoration_project/models/restorant_models"
	"fmt"
)

func GetServicesByCategory(categoryID string) ([]restorantmodels.RestaurantService, error) {
	// Fetch services for the category
query := `
    SELECT 
        service_id, category_id, service_name, service_description,
        service_price, average_rating,
        COALESCE(proposed_restaurant_id, '') AS proposed_restaurant_id,
        created_at, updated_at
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
			&service.AverageRating,
			&service.ProposedRestorantId,
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

	// Fetch service details including average_rating
	query := `
		SELECT 
			service_id, category_id, service_name, service_description, service_price, average_rating,proposed_restaurant_id, created_at, updated_at
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
		&result.Service.AverageRating,
		&result.Service.ProposedRestorantId,
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

	var images []string
	for imageRows.Next() {
		var img string
		if err := imageRows.Scan(&img); err != nil {
			return result, err
		}
		images = append(images, img)
	}
	result.Service.Images = images

	// ✅ Fetch reviews for this service
	reviewRows, err := config.DB.Query(`
		SELECT user_name, rating, review_text, created_at
		FROM Service_Reviews
		WHERE service_id = ?
		ORDER BY created_at DESC
	`, result.Service.ServiceID)
	if err != nil {
		fmt.Println(err.Error())
		return result, err
	}
	defer reviewRows.Close()

	var reviews []restorantmodels.ServiceReview
	for reviewRows.Next() {
		var r restorantmodels.ServiceReview
		if err := reviewRows.Scan(&r.UserName, &r.Rating, &r.ReviewText, &r.CreatedAt); err != nil {
			return result, err
		}
		reviews = append(reviews, r)
	}
	result.Reviews = reviews
	fmt.Println("here we have")
	return result, nil
}

func SearchServicesByName(search string) ([]restorantmodels.RestaurantService, error) {
	query := `
		SELECT 
			service_id, category_id, service_name, service_description, 
			service_price,average_rating, proposed_restaurant_id, created_at, updated_at
		FROM Our_Services
		WHERE service_name LIKE ?
	`

	// %search% pattern for LIKE query
	rows, err := config.DB.Query(query, "%"+search+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to search services: %w", err)
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
			&service.AverageRating,
			&service.ProposedRestorantId,
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
