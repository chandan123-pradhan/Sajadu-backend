package staffservices

import (
	restorantmodels "decoration_project/models/restorant_models"
	staffmodel "decoration_project/models/staff_model"
	staffrepo "decoration_project/repository/staff_repo"
)


func GetAssignedServices(restorantId string) (staffmodel.AssignedBookingsWrapper, error) {
	// Call repository to fetch bookings
	return staffrepo.GetAssignedBookings(restorantId)
}



func GetAssignedServiesDetails(bookingId string) (restorantmodels.BookingDetailsResponse, error) {
	// Call repository to fetch bookings
	return staffrepo.GetAssignedServiceDetails(bookingId)
}


func VerifyOtpToStartService(bookingId string, stafOtp string, key string) error {
	// Call repository to fetch bookings
	return staffrepo.VerifyToStartService(
		bookingId,stafOtp,key,)
}

func UpdateStaffLocation(staffID, bookingID string, latitude, longitude float64) error {
    return staffrepo.SaveStaffLocation(staffID, bookingID, latitude, longitude)
}

func GetPartnerLocationByBookingID(bookingID string) (staffmodel.PartnerLocationResponse, error) {
	return staffrepo.FetchPartnerLocation(bookingID)
}