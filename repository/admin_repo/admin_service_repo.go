package adminrepo

import (
	"database/sql"
	"decoration_project/config"
	adminmodel "decoration_project/models/admin_model"
	restorantmodels "decoration_project/models/restorant_models"
	"errors"
	"fmt"
	"strings"

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

// GetAllServicesWithImages fetches all non-deleted services for a category with images
func GetAllServiceCategoryWise(categoryId string) ([]restorantmodels.RestaurantService, error) {
	var services []restorantmodels.RestaurantService
	var proposedRestaurantID sql.NullString

	query := `SELECT 
				service_id, category_id, service_name, service_description, 
				service_price, proposed_restaurant_id, created_at, updated_at
			  FROM Our_Services
			  WHERE category_id = ? AND is_deleted = 0`

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

		// Handle NULL proposed_restaurant_id
		if proposedRestaurantID.Valid {
			service.ProposedRestorantId = proposedRestaurantID.String
		} else {
			service.ProposedRestorantId = ""
		}

		// Fetch images
		imgRows, err := config.DB.Query(`SELECT image_url FROM Service_Images WHERE service_id = ?`, service.ServiceID)
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



func CreateFestival(f adminmodel.Festival) error {

    query := `
        INSERT INTO Festivals (festival_name, start_date, end_date)
        VALUES (?, ?, ?)
    `

    _, err := config.DB.Exec(query, f.FestivalName, f.StartDate, f.EndDate)
    if err != nil {

        // Check for duplicate entry error
        if strings.Contains(err.Error(), "Error 1062") {
            return fmt.Errorf("festival name already exists")
        }

        // Any other DB error
        return fmt.Errorf("failed to create festival: %v", err)
    }

    return nil
}



func GetAllFestivals() ([]adminmodel.Festival, error) {
    query := `
        SELECT festival_id, festival_name, start_date, end_date
        FROM Festivals
        ORDER BY start_date DESC
    `

    rows, err := config.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var list []adminmodel.Festival

    for rows.Next() {
        var f adminmodel.Festival
        err := rows.Scan(&f.FestivalID, &f.FestivalName, &f.StartDate, &f.EndDate)
        if err != nil {
            return nil, err
        }
        list = append(list, f)
    }

    return list, nil
}

func AddServiceToFestival(festivalID string, serviceID string, discountPercent float64) error {
    query := `
        INSERT INTO festival_services (festival_id, service_id, discount_percent)
        VALUES (?, ?, ?)
    `

    _, err := config.DB.Exec(query, festivalID, serviceID, discountPercent)
    if err != nil {

        // Duplicate entry error
        if strings.Contains(err.Error(), "Duplicate entry") {
            return errors.New("this service is already added for this festival")
        }

        return err
    }

    return nil
}


// GetFestivalServices fetches all services linked to a given festival with images and discount
func GetFestivalServices(festivalId string) ([]restorantmodels.RestaurantService, error) {
	var services []restorantmodels.RestaurantService
	var proposedRestaurantID sql.NullString
	var discountPercent sql.NullFloat64

	query := `
		SELECT 
			s.service_id,
			s.category_id,
			s.service_name,
			s.service_description,
			s.service_price,
			s.proposed_restaurant_id,
			fs.discount_percent,
			s.created_at,
			s.updated_at
		FROM Our_Services s
		INNER JOIN festival_services fs ON s.service_id = fs.service_id
		WHERE fs.festival_id = ? AND s.is_deleted = 0
	`

	rows, err := config.DB.Query(query, festivalId)
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
			&discountPercent,
			&service.CreatedAt,
			&service.UpdatedAt,
		); err != nil {
			return nil, err
		}

		// Handle NULL proposed_restaurant_id
		if proposedRestaurantID.Valid {
			service.ProposedRestorantId = proposedRestaurantID.String
		} else {
			service.ProposedRestorantId = ""
		}

		// Handle discount percent
		if discountPercent.Valid {
			service.DiscountPercent = discountPercent.Float64
		} else {
			service.DiscountPercent = 0.0
		}

		// Fetch service images
		imgRows, err := config.DB.Query(
			`SELECT image_url FROM Service_Images WHERE service_id = ?`,
			service.ServiceID,
		)
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
