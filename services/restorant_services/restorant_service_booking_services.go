package restorantservices

import (
	restorantmodels "decoration_project/models/restorant_models"
	restorantrepo "decoration_project/repository/restorant_repo"
	"fmt"
)

func FetchRestorantBookings(restorantId string, key string) (restorantmodels.RestaurantBookingsWrapper, error) {
	// Call repository to fetch bookings
	return restorantrepo.GetRestaurantBookings(restorantId,key)
}

func FetchRestorantBookingDetails(bookingId string, key string) (restorantmodels.BookingDetailsResponse, error) {
	// Call repository to fetch bookings
	return restorantrepo.GetRestaurantBookingDetails(bookingId,key)
}


// AcceptBooking accepts a booking and generates OTP for staff verification
func AcceptBooking(bookingID string, acceptedStatusID int, key string) (string, error) {
	// Call repository to accept booking and generate encrypted OTP
	otp, err := restorantrepo.AcceptBooking(bookingID, acceptedStatusID, key)
	if err != nil {
		return "", fmt.Errorf("failed to accept booking: %w", err)
	}

	// OTP returned so you can send it to staff via SMS/Notification
	return otp, nil
}


// RejectBooking rejects a booking and updates status with a reason
func RejectBooking(bookingID string, rejectedStatusID int, reason string) error {
	err := restorantrepo.RejectBooking(bookingID, rejectedStatusID, reason)
	if err != nil {
		return fmt.Errorf("failed to reject booking: %w", err)
	}

	return nil
}


// AssignStaff assigns a staff member to a booking
func AssignStaff(bookingID string, staffID string) error {
	err := restorantrepo.AssignStaff(bookingID, staffID)
	if err != nil {
		return fmt.Errorf("failed to assign staff: %w", err)
	}
	return nil
}

