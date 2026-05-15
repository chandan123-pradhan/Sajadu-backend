package adminrepo

import (
	"decoration_project/config"
	usermodels "decoration_project/models/user_models"
)



func GetAllUsers() ([]usermodels.UserDetailsModel, error) {
	var users []usermodels.UserDetailsModel

	query := `
		SELECT 
			user_id,
			full_name,
			email,
			mobile_number
		FROM Users
	`

	rows, err := config.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user usermodels.UserDetailsModel

		err := rows.Scan(
			&user.UserID,
			&user.FullName,
			&user.Email,
			&user.MobileNumber,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}