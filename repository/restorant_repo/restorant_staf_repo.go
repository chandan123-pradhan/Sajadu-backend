package restorantrepo

import (
	"decoration_project/config"
	restorantmodels "decoration_project/models/restorant_models"
	"github.com/google/uuid"
	"time"
)

// AddStaff inserts a new staff member into the Restaurant_Staff table
func AddStaff(staff restorantmodels.RestaurantStaff) (string, error) {
	staffID := uuid.New().String()
	now := time.Now()

	query := `
		INSERT INTO Restaurant_Staff 
		(staff_id, restaurant_id, name, email, whatsapp_no, designation, description, password, image_url, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := config.DB.Exec(query,
		staffID,
		staff.RestaurantID,
		staff.Name,
		staff.Email,
		staff.WhatsappNo,
		staff.Designation,
		staff.Description,
		staff.Password,   // already hashed before calling repo
		staff.ImageURL,   // single profile image
		now,
		now,
	)

	if err != nil {
		return "", err
	}

	return staffID, nil
}


func GetStaffByRestaurant(restaurantID string) ([]restorantmodels.RestaurantStaff, error) {
    query := "SELECT staff_id, restaurant_id, name, email, whatsapp_no, designation, description, password, image_url, created_at, updated_at FROM Restaurant_Staff WHERE restaurant_id = ?"
    rows, err := config.DB.Query(query, restaurantID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var staffList []restorantmodels.RestaurantStaff
    for rows.Next() {
        var staff restorantmodels.RestaurantStaff
        if err := rows.Scan(&staff.StaffID, &staff.RestaurantID, &staff.Name, &staff.Email, &staff.WhatsappNo, &staff.Designation, &staff.Description, &staff.Password, &staff.ImageURL, &staff.CreatedAt, &staff.UpdatedAt); err != nil {
            return nil, err
        }
        staffList = append(staffList, staff)
    }

    return staffList, nil
}