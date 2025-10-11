package userrepo


import (
	"decoration_project/config"
	usermodels "decoration_project/models/user_models"
	"fmt"
)

// SaveNotification inserts a new notification into the database
func SaveNotification(notification usermodels.UserNotification) error {
	query := `
		INSERT INTO UserNotification (user_id, title, description, image)
		VALUES (?, ?, ?, ?)
	`

	_, err := config.DB.Exec(query, notification.UserID, notification.Title, notification.Description, notification.Image)
	if err != nil {
		return fmt.Errorf("failed to insert notification: %v", err)
	}

	return nil
}

// GetUserNotifications retrieves notifications for a given user

func GetUserNotifications(userID string) ([]usermodels.UserNotification, error) {
	query := `
		SELECT notification_id, user_id, title, description, image, created_at
		FROM UserNotification
		WHERE user_id = ?
		ORDER BY created_at DESC
	`

	rows, err := config.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch notifications: %v", err)
	}
	defer rows.Close()

	notifications := []usermodels.UserNotification{} 
	for rows.Next() {
		var n usermodels.UserNotification
		if err := rows.Scan(
			&n.ID,
			&n.UserID,
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
