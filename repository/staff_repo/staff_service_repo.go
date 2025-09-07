package staffrepo

import (
	"database/sql"
	"decoration_project/config"
	restorantmodels "decoration_project/models/restorant_models"
	staffmodel "decoration_project/models/staff_model"
	usermodels "decoration_project/models/user_models"
	"decoration_project/utils"
	"errors"
	"fmt"
)

// GetAssignedBookings fetches all bookings assigned to a staff member
func GetAssignedBookings(staffID string) (staffmodel.AssignedBookingsWrapper, error) {
	query := `
		SELECT 
			b.booking_id,
			b.user_id,
			b.service_id,
			s.status_name,
			b.scheduled_date,
			b.address,
			b.service_name,
			b.created_at,

			-- Payment details (may be NULL if no payment yet)
			p.payment_id,
			COALESCE(p.amount, 0),
			COALESCE(p.currency, ''),
			COALESCE(p.payment_mode, ''),
			COALESCE(p.transaction_id, ''),
			p.payment_date,
			COALESCE(p.payment_status, '')
		FROM bookings b
		JOIN booking_status s ON b.status_id = s.status_id
		LEFT JOIN payments p ON b.booking_id = p.booking_id
		WHERE b.staff_id = ?
		ORDER BY b.created_at DESC
	`

	rows, err := config.DB.Query(query, staffID)
	if err != nil {
		return staffmodel.AssignedBookingsWrapper{Bookings: []staffmodel.AssignedBookingsResponse{}}, err
	}
	defer rows.Close()

	var bookings []staffmodel.AssignedBookingsResponse

	for rows.Next() {
		var booking staffmodel.AssignedBookingsResponse
		var payment usermodels.PaymentResponse

		var paymentID, transactionID, paymentMethod, currency, paymentStatus sql.NullString
		var paymentDate sql.NullTime
		var amount sql.NullFloat64

		err := rows.Scan(
			&booking.BookingID,
			&booking.UserID,
			&booking.ServiceID,
			&booking.Status,
			&booking.ScheduledDate,
			&booking.Address,
			&booking.ServiceName,
			&booking.CreatedAt,
			&paymentID,
			&amount,
			&currency,
			&paymentMethod,
			&transactionID,
			&paymentDate,
			&paymentStatus,
		)
		if err != nil {
			return staffmodel.AssignedBookingsWrapper{Bookings: []staffmodel.AssignedBookingsResponse{}}, err
		}

		// Map price from payment amount
		if amount.Valid {
			booking.Price = amount.Float64
		}

		// Populate payment details if exists
		if paymentID.Valid {
			payment.PaymentID = paymentID.String
			payment.Amount = amount.Float64
			payment.TransactionID = transactionID.String
			payment.PaymentMethod = paymentMethod.String
			payment.Currency = currency.String
			payment.PaymentStatus = paymentStatus.String
			if paymentDate.Valid {
				payment.PaymentDate = paymentDate.Time
			}
			booking.Payment = &payment
		} else {
			booking.Payment = nil
		}

		bookings = append(bookings, booking)
	}

	if err = rows.Err(); err != nil {
		return staffmodel.AssignedBookingsWrapper{Bookings: []staffmodel.AssignedBookingsResponse{}}, err
	}

	return staffmodel.AssignedBookingsWrapper{Bookings: bookings}, nil
}

// GetRestaurantBookingDetails fetches full details of a specific booking
func GetAssignedServiceDetails(bookingID string) (restorantmodels.BookingDetailsResponse, error) {
	query := `
	SELECT 
		b.booking_id,
		u.user_id,
		u.full_name,
		u.email,
		u.mobile_number,

		b.service_id,
		b.service_name,
		s.status_name,
		b.scheduled_date,
		b.address,
		b.pincode,
		b.state,
		b.city,
		b.latitude,
		b.longitude,
		b.created_at,


		p.payment_id,
		COALESCE(p.amount,0),
		COALESCE(p.currency,''),
		COALESCE(p.payment_mode,''),
		COALESCE(p.transaction_id,''),
		p.payment_date,
		COALESCE(p.payment_status,'')
	FROM bookings b
	JOIN booking_status s ON b.status_id = s.status_id
	JOIN users u ON b.user_id = u.user_id
	LEFT JOIN Restaurant_Staff rs ON rs.staff_id = b.staff_id
	LEFT JOIN payments p ON b.booking_id = p.booking_id
	WHERE b.booking_id = ?
	`

	row := config.DB.QueryRow(query, bookingID)

	var booking restorantmodels.BookingDetailsResponse
	var user usermodels.UserDetailsModel
	var staff restorantmodels.StaffDetails
	var payment usermodels.PaymentResponse

	var staffID, staffName, staffWhatsapp, staffDesignation, staffEmail, staffDescription, staffImage sql.NullString
	var paymentID, transactionID, paymentMethod, currency, paymentStatus sql.NullString
	var paymentDate sql.NullTime
	var amount sql.NullFloat64

	err := row.Scan(
		&booking.BookingID,
		&user.UserID,
		&user.FullName,
		&user.Email,
		&user.Email,

		&booking.ServiceID,
		&booking.ServiceName,
		&booking.Status,
		&booking.ScheduledDate,
		&booking.Address,
		&booking.Pincode,
		&booking.State,
		&booking.City,
		&booking.Latitude,
		&booking.Longitude,
		&booking.CreatedAt,

		&paymentID,
		&amount,
		&currency,
		&paymentMethod,
		&transactionID,
		&paymentDate,
		&paymentStatus,
	)
	if err != nil {
		return restorantmodels.BookingDetailsResponse{}, err
	}

	// Assign user
	booking.User = user

	// Assign staff if available
	if staffID.Valid {
		staff.StaffID = staffID.String
		staff.Name = staffName.String
		staff.WhatsappNo = staffWhatsapp.String
		staff.Designation = staffDesignation.String
		staff.Email = staffEmail.String
		staff.Description = staffDescription.String
		staff.ImageURL = staffImage.String
		booking.Staff = &staff
	} else {
		booking.Staff = nil
	}

	// Assign payment if available
	if paymentID.Valid {
		payment.PaymentID = paymentID.String
		payment.Amount = amount.Float64
		payment.Currency = currency.String
		payment.PaymentMethod = paymentMethod.String
		payment.TransactionID = transactionID.String
		payment.PaymentStatus = paymentStatus.String
		if paymentDate.Valid {
			payment.PaymentDate = paymentDate.Time
		}
		booking.Payment = &payment
	} else {
		booking.Payment = nil
	}

	return booking, nil
}
// VerifyToStartService verifies the staff OTP and updates booking status to "In Progress"
// and generates the completion OTP for the user
func VerifyToStartService(bookingID string, staffOTP string, key string) error {
	// 1. Fetch encrypted OTP, current status, and start_verified
	var encryptedOTP string
	var statusID int
	var startVerified bool
	query := `SELECT start_otp_hash, status_id, start_verified FROM bookings WHERE booking_id = ?`
	err := config.DB.QueryRow(query, bookingID).Scan(&encryptedOTP, &statusID, &startVerified)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("booking not found")
		}
		return err
	}

	// 2. Check if OTP is already verified or status is already "In Progress"
	if startVerified || statusID == 4 {
		return errors.New("service already started or OTP already verified")
	}

	// 3. Decrypt OTP
	decryptedOTP, err := utils.DecryptOTP(encryptedOTP, key)
	if err != nil {
		return errors.New("failed to decrypt OTP")
	}

	// 4. Compare OTP
	if decryptedOTP != staffOTP {
		return errors.New("invalid OTP")
	}

	// 5. Generate completion OTP for user
	completeOTP := utils.GenerateOTP() // implement a 4-6 digit OTP generator
	encryptedCompleteOTP, err := utils.EncryptOTP(completeOTP, key)
	if err != nil {
		return errors.New("failed to encrypt completion OTP")
	}

	// 6. Update status to "In Progress", mark start_verified, and save encrypted completion OTP
	updateQuery := `
		UPDATE bookings 
		SET status_id = 4, start_verified = TRUE, complete_otp_hash = ?, updated_at = CURRENT_TIMESTAMP
		WHERE booking_id = ?
	`
	_, err = config.DB.Exec(updateQuery, encryptedCompleteOTP, bookingID)
	if err != nil {
		return err
	}

	// 7. Optionally, return or log the plain completion OTP so it can be sent to the user
	fmt.Println("Completion OTP for booking", bookingID, "is", completeOTP)

	return nil
}
