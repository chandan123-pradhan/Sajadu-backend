package userservices

import (
	usermodels "decoration_project/models/user_models"
	userrepo "decoration_project/repository/user_repo"
)

// SaveNotificationService — handles logic before saving notification
func SaveNotificationService(notification usermodels.UserNotification) error {
	// Any business logic can be added here later
	return userrepo.SaveNotification(notification)
}

// GetUserNotificationsService — fetch notifications for a specific user
func GetUserNotificationsService(userID string) ([]usermodels.UserNotification, error) {
	return userrepo.GetUserNotifications(userID)
}
