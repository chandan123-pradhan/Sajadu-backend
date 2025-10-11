package restorantrepo

import (
	"decoration_project/config"
	restorantmodels "decoration_project/models/restorant_models"
	"fmt"
)

// SaveNotification inserts a new notification into the database
func SaveNotification(notification restorantmodels.RestorantNotifications) error {
	query := `
		INSERT INTO RestorantNotification (restaurant_id, title, description, image)
		VALUES (?, ?, ?, ?)
	`

	_, err := config.DB.Exec(query, notification.RestorantId, notification.Title, notification.Description, notification.Image)
	if err != nil {
		return fmt.Errorf("failed to insert notification: %v", err)
	}

	return nil
}

// GetUserNotifications retrieves notifications for a given user

func GetRestorantNotifications(restorantId string) ([]restorantmodels.RestorantNotifications, error) {
	query := `
		SELECT notification_id, restaurant_id, title, description, image, created_at
		FROM RestorantNotification
		WHERE restaurant_id = ?
		ORDER BY created_at DESC
	`

	rows, err := config.DB.Query(query, restorantId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch notifications: %v", err)
	}
	defer rows.Close()

	notifications := []restorantmodels.RestorantNotifications{} 
	for rows.Next() {
		var n restorantmodels.RestorantNotifications
		if err := rows.Scan(
			&n.ID,
			&n.RestorantId,
			&n.Title,
			&n.Description,
			&n.Image,
			&n.CreatedAt,
		); err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}

	return notifications, nil
}
