package staffrepo

import (
	"database/sql"
	"decoration_project/config"
	staffmodel "decoration_project/models/staff_model"
)

// GetStaffByEmail fetches staff details by email
func GetStaffByEmail(email string) (staffmodel.Staff, error) {
	query := `
		SELECT 
			staff_id,
			restaurant_id,
			name,
			email,
			whatsapp_no,
			designation,
			description,
			password,
			image_url,
			created_at,
			updated_at
		FROM Restaurant_Staff
		WHERE email = ?;
	`

	var staff staffmodel.Staff
	row := config.DB.QueryRow(query, email)

	err := row.Scan(
		&staff.StaffID,
		&staff.RestaurantID,
		&staff.Name,
		&staff.Email,
		&staff.WhatsappNo,
		&staff.Designation,
		&staff.Description,
		&staff.Password,
		&staff.ImageURL,
		&staff.CreatedAt,
		&staff.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return staffmodel.Staff{}, nil // return empty if not found
		}
		return staffmodel.Staff{}, err
	}

	return staff, nil
}
