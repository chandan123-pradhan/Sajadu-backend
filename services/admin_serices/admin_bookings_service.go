package adminserices

import (
	adminmodel "decoration_project/models/admin_model"
	adminrepo "decoration_project/repository/admin_repo"
)

func GetAllRestorant() ([]adminmodel.RestorantModel, error) {
    // Fetch all services with images from repository
    restorants, err := adminrepo.GetAllRestorants()
    if err != nil {
        return nil, err
    }
    return restorants, nil
}


func GetAllActiveBookings(status []string ) ([]adminmodel.BookingResponse, error) {
    // Fetch all services with images from repository
    bookings, err := adminrepo.GetAllActiveBookings(status)
    if err != nil {
        return nil, err
    }
    return bookings, nil
}



func GetBookingsDetails(booking_id string) (adminmodel.BookingResponse, error) {
    // Fetch all services with images from repository
    bookings, err := adminrepo.GetBookingDetailsByID(booking_id)
    if err != nil {
        return adminmodel.BookingResponse{}, err
    }
    return bookings, nil
}


// ✅ Update booking status and assign restaurant if accepted
func UpdateBookingStatus(bookingID, restaurantID, newStatus string) error {
	err := adminrepo.UpdateBookingStatus(bookingID, restaurantID, newStatus)
	if err != nil {
		return err
	}
	return nil
}