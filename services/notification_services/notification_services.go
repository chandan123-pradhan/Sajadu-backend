package notificationservices

import (
	restorantmodels "decoration_project/models/restorant_models"
	usermodels "decoration_project/models/user_models"
	adminrepo "decoration_project/repository/admin_repo"
	restorantrepo "decoration_project/repository/restorant_repo"
	userrepo "decoration_project/repository/user_repo"
	"decoration_project/services"
	restorantservices "decoration_project/services/restorant_services"
	userservices "decoration_project/services/user_services"
	"fmt"
	"log"
)

/*
SendBookingNotificationToAdmin
---------------------------------
Sends a push notification to the admin when a new booking is created.

Parameters:
- userId (string): The ID of the user who created the booking.

Returns:
- error: Returns an error if fetching the admin FCM token or sending the notification fails.

Notes:
- Uses the admin repository to fetch the admin FCM token.
- Notification contains the user ID who created the booking.
*/
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

	// Prepare notification content
	title := "New Booking Received"
	body := "A new booking has been created by user ID: " + userId

	// Send notification via FCM service
	if err := services.SendNotification(adminToken, title, body); err != nil {
		log.Println("Failed to send notification:", err)
		return err
	}

	log.Println("✅ Notification sent to admin successfully")
	return nil
}

/*
SendBookingNotificationToUser
---------------------------------
Sends a booking-related notification to a specific user and logs it in the database.

Parameters:
- userId (string): The ID of the user to notify.
- item (string): Name of the booked item or service.
- reason (string): Status update or reason to include in the notification.
- imageUrl (string): Optional image URL to save along with the notification log.

Returns:
  - error: Returns an error if fetching the user's FCM token, sending the notification,
    or saving the notification log fails.

Notes:
- Saves the notification to the UserNotification table after sending.
*/
func SendBookingNotificationToUser(userId string, item string, reason string, imageUrl string) error {
	// Fetch the user's FCM token
	userToken, err := userrepo.GetUserFcmToken(userId)
	if err != nil {
		log.Println("Failed to fetch user FCM token:", err)
		return err
	}

	if userToken == "" {
		log.Println("No user FCM token found")
		return nil
	}

	// Prepare notification content
	title := "Booking Update"
	body := fmt.Sprintf("Your new order item %s has been %s by us.", item, reason)

	// Save notification log in DB
	notification := usermodels.UserNotification{
		UserID:      userId,
		Title:       title,
		Description: body,
		Image:       imageUrl,
	}

	if err := userservices.SaveNotificationService(notification); err != nil {
		log.Println("⚠️ Failed to save notification log:", err)
		return err
	}
	log.Println("💾 Notification saved to UserNotification table successfully")

	// Send notification via FCM service
	if err := services.SendNotification(userToken, title, body); err != nil {
		log.Println("Failed to send notification:", err)
		return err
	}

	log.Println("✅ Notification sent to user successfully")
	return nil
}



func SendBookingNormalNotificationToUser(userId string, title string, description string) error {
	// Fetch the user's FCM token
	userToken, err := userrepo.GetUserFcmToken(userId)
	if err != nil {
		log.Println("Failed to fetch user FCM token:", err)
		return err
	}

	if userToken == "" {
		log.Println("No user FCM token found")
		return nil
	}

	// Prepare notification content
	
	// Save notification log in DB
	notification := usermodels.UserNotification{
		UserID:      userId,
		Title:       title,
		Description: description,
		Image:       "",
	}

	if err := userservices.SaveNotificationService(notification); err != nil {
		log.Println("⚠️ Failed to save notification log:", err)
		return err
	}
	log.Println("💾 Notification saved to UserNotification table successfully")

	// Send notification via FCM service
	if err := services.SendNotification(userToken, title, description); err != nil {
		log.Println("Failed to send notification:", err)
		return err
	}

	log.Println("✅ Notification sent to user successfully")
	return nil
}


/*
SendBookingNotificationToRestorant
---------------------------------
Sends a booking notification to a specific restaurant.

Parameters:
- restorantId (string): The ID of the restaurant to notify.
- item (string): Name of the booked item/service.
- price (string): Price of the booking to include in the notification.
- reason (string): Optional reason/status for the notification.

Returns:
- error: Returns an error if fetching the restaurant's FCM token or sending the notification fails.

Notes:
- Notification informs the restaurant about a new booking with item name and price.
*/
func SendBookingNotificationToRestorant(restorantId string, item string, price string, reason string, imageUrl string) error {
	// Fetch the restaurant's FCM token
	fcmToken, err := restorantrepo.GetRestorantFcmToken(restorantId)
	if err != nil {
		log.Println("Failed to fetch restaurant FCM token:", err)
		return err
	}

	if fcmToken == "" {
		log.Println("No restaurant FCM token found")
		return nil
	}

	// Prepare notification content
	title := "New Booking"
	body := fmt.Sprintf("Hey, you got another booking %s at ₹%s by Sajadu. Please check your bookings page.", item, price)

	notification := restorantmodels.RestorantNotifications{
		RestorantId: restorantId,
		Title:       title,
		Description: body,
		Image:       imageUrl,
	}

	if err := restorantservices.SaveNotificationService(notification); err != nil {
		log.Println("⚠️ Failed to save notification log:", err)
		return err
	}
	log.Println("💾 Notification saved to UserNotification table successfully")

	// Send notification via FCM service
	if err := services.SendNotification(fcmToken, title, body); err != nil {
		log.Println("Failed to send notification:", err)
		return err
	}

	log.Println("✅ Notification sent to restaurant successfully")
	return nil
}
