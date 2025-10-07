package adminrepo


import (
	"decoration_project/config"
	"fmt"
)

// SaveOrUpdateAdminFCM saves or updates the single admin FCM token
func SaveOrUpdateAdminFCM(token string) error {
	query := `
		INSERT INTO Admin_FCM_Token (id, fcm_token)
		VALUES (1, ?)
		ON DUPLICATE KEY UPDATE 
			fcm_token = VALUES(fcm_token),
			updated_at = CURRENT_TIMESTAMP
	`

	_, err := config.DB.Exec(query, token)
	if err != nil {
		return fmt.Errorf("failed to save/update FCM token: %v", err)
	}
	return nil
}



func GetAdminFCMToken() (string, error) {
	var token string
	query := `SELECT fcm_token FROM Admin_FCM_Token WHERE id = 1`
	err := config.DB.QueryRow(query).Scan(&token)
	if err != nil {
		return "", fmt.Errorf("failed to get admin FCM token: %v", err)
	}
	return token, nil
}