package staffservices

import (
	staffmodel "decoration_project/models/staff_model"
	staffrepo "decoration_project/repository/staff_repo"
)


func GetAssignedServices(restorantId string) (staffmodel.AssignedBookingsWrapper, error) {
	// Call repository to fetch bookings
	return staffrepo.GetAssignedBookings(restorantId)
}



func GetAssignedServiesDetails(bookingId string) (staffmodel.StaffAssignedServicesDetails, error) {
	// Call repository to fetch bookings
	return staffrepo.GetBookingFullDetails(bookingId)
}


func StartService(bookingId string, staffId string, key string) error {
	// Call repository to fetch bookings
	return staffrepo.StartService(
		bookingId, staffId,key)
}

func UpdateStaffLocation(staffID, bookingID string, latitude, longitude float64) error {
    return staffrepo.SaveStaffLocation(staffID, bookingID, latitude, longitude)
}

func GetPartnerLocationByBookingID(bookingID string) (staffmodel.PartnerLocationResponse, error) {
	return staffrepo.FetchPartnerLocation(bookingID)
}
func VerifyOtpToCompleteService(bookingId string, stafOtp string, key string) error {
	// Call repository to verify OTP and update status → Completed
	return staffrepo.VerifyToCompleteService(
		bookingId, stafOtp, key,
	)
}