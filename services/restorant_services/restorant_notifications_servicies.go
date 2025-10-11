package restorantservices

import (
	restorantmodels "decoration_project/models/restorant_models"
	restorantrepo "decoration_project/repository/restorant_repo"
	
)

// SaveNotificationService — handles logic before saving notification
func SaveNotificationService(notification restorantmodels.RestorantNotifications) error {
	// Any business logic can be added here later
	return restorantrepo.SaveNotification(notification)
}

// GetUserNotificationsService — fetch notifications for a specific user
func GetRestorantNotificationsService(restorantId string) ([]restorantmodels.RestorantNotifications, error) {
	return restorantrepo.GetRestorantNotifications(restorantId)
}
