package userservices

import (
	// staffmodel "decoration_project/models/staff_model"
	"decoration_project/models/user_models"
	userrepo "decoration_project/repository/user_repo"
)

// CreateBooking handles booking creation logic
func CreateBooking(req usermodels.BookingRequest) (usermodels.BookingResponse, error) {
	// Call repository to insert booking
	return userrepo.CreateBooking(req)
}

func FetchBookings(userId string,key string) (usermodels.UserBookingsWrapper, error) {
	// Call repository to fetch bookings
	return userrepo.GetUserBookings(userId,key)
}


func FetchBookingDetails(bookingId string) (usermodels.GetUsersBookingDetailsResponse, error) {
	// Call repository to fetch bookings
	return userrepo.GetBookingDetails(bookingId)
}