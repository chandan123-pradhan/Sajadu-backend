package notificationservices

import (
	adminrepo "decoration_project/repository/admin_repo"
	userrepo "decoration_project/repository/user_repo"
	"decoration_project/services"
	"log"
)

// SendBookingNotificationToAdmin sends a notification to the admin when a new booking is created
func SendBookingNotificationToAdmin(userId string) error {
	// Fetch the admin FCM token
	adminToken, err := adminrepo.GetAdminFCMToken()
	if err != nil {
		log.Println("Failed to fetch admin FCM token:", err)
		return err
	}

	if adminToken == "" {
		log.Println("No admin FCM token found")
		return nil
	}
	

	// Prepare notification
	title := "New Booking Received"
	body := "A new booking has been created by user ID: " + userId

	// Send notification
	if err := services.SendNotification(adminToken, title, body); err != nil {
		log.Println("Failed to send notification:", err)
		return err
	}

	log.Println("✅ Notification sent to admin successfully")
	return nil
}



// SendBookingNotificationToAdmin sends a notification to the admin when a new booking is created
func SendBookingNotificationToUser(userId string, item string, reson string) error {
	// Fetch the admin FCM token
	adminToken, err := userrepo.GetUserFcmToken(userId)
	if err != nil {
		log.Println("Failed to fetch admin FCM token:", err)
		return err
	}

	if adminToken == "" {
		log.Println("No admin FCM token found")
		return nil
	}
	
 
	// Prepare notification
	title := "Booking Update"
	body := "Your new order item "+item+"has been "+reson+" by Us"

	// Send notification
	if err := services.SendNotification(adminToken, title, body); err != nil {
		log.Println("Failed to send notification:", err)
		return err
	}

	log.Println("✅ Notification sent to user successfully")
	return nil
}
