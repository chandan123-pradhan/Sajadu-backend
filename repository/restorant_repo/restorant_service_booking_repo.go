package restorantrepo

import (
	"database/sql"
	"decoration_project/config"
	restorantmodels "decoration_project/models/restorant_models"
	usermodels "decoration_project/models/user_models"
	"decoration_project/utils"
	"fmt"
	"time"
)

// GetRestaurantBookings fetches all bookings of a given restaurant
func GetRestaurantBookings(restaurantID string, key string) (restorantmodels.RestaurantBookingsWrapper, error) {
	query := `
		SELECT 
			b.booking_id,
			b.user_id,
			b.service_id,
			s.status_name,
			b.scheduled_date,
			b.address,
			b.pincode,
			b.state,
			b.city,
			b.service_name,
			b.price,
			b.created_at,
			b.start_otp_hash,  -- Encrypted OTP for staff verification
			rs.service_description,  -- added service description

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
		LEFT JOIN Our_Services rs ON b.service_id = rs.service_id  -- join to get service description
		LEFT JOIN payments p ON b.booking_id = p.booking_id
		WHERE b.restaurant_id = ?
		ORDER BY b.created_at DESC
	`

	rows, err := config.DB.Query(query, restaurantID)
	if err != nil {
		return restorantmodels.RestaurantBookingsWrapper{Bookings: []restorantmodels.RestorantBookingsResponse{}}, err
	}
	defer rows.Close()

	var bookings []restorantmodels.RestorantBookingsResponse

	for rows.Next() {
		var booking restorantmodels.RestorantBookingsResponse
		var payment usermodels.PaymentResponse

		var paymentID, transactionID, paymentMethod, currency, paymentStatus, startOTP sql.NullString
		var paymentDate sql.NullTime
		var amount sql.NullFloat64
		var serviceDesc sql.NullString

		err := rows.Scan(
			&booking.BookingID,
			&booking.UserID,
			&booking.ServiceID,
			&booking.Status,
			&booking.ScheduledDate,
			&booking.Address,
			&booking.Pincode,
			&booking.State,
			&booking.City,
			&booking.ServiceName,
			&booking.Price,
			&booking.CreatedAt,
			&startOTP,  
			&serviceDesc,  // scan service description
			&paymentID,
			&amount,
			&currency,
			&paymentMethod,
			&transactionID,
			&paymentDate,
			&paymentStatus,
		)
		if err != nil {
			return restorantmodels.RestaurantBookingsWrapper{Bookings: []restorantmodels.RestorantBookingsResponse{}}, err
		}

		// Assign service description
		if serviceDesc.Valid {
			booking.ServiceDesc = serviceDesc.String
		} else {
			booking.ServiceDesc = ""
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

		// Assign decrypted OTP only if status is not Pending or Rejected
		if startOTP.Valid && booking.Status != "Pending" && booking.Status != "Rejected" {
			decryptedOtp, err := utils.DecryptOTP(startOTP.String, key)
			if err != nil {
				booking.StartOtp = ""
			} else {
				booking.StartOtp = decryptedOtp
			}
		} else {
			booking.StartOtp = ""
		}

		bookings = append(bookings, booking)
	}

	if err = rows.Err(); err != nil {
		return restorantmodels.RestaurantBookingsWrapper{Bookings: []restorantmodels.RestorantBookingsResponse{}}, err
	}

	return restorantmodels.RestaurantBookingsWrapper{Bookings: bookings}, nil
}


// GetRestaurantBookingDetails fetches full details of a specific booking
func GetRestaurantBookingDetails(bookingID, key string) (restorantmodels.BookingDetailsResponse, error) {
	query := `
	SELECT 
		b.booking_id,
		u.user_id,
		u.full_name,
		u.email,
		u.mobile_number,

		b.service_id,
		b.service_name,
		rs.service_description,   -- added service description
		b.price,

		s.status_name,
		b.scheduled_date,
		b.address,
		b.pincode,
		b.state,
		b.city,
		b.latitude,
		b.longitude,
		b.created_at,
		b.start_otp_hash,

		staff.staff_id,
		staff.name,
		staff.whatsapp_no,
		staff.designation,
		staff.email,
		staff.description,
		staff.image_url,

		p.payment_id,
		COALESCE(p.amount,0),
		COALESCE(p.currency,''),
		COALESCE(p.payment_mode,''),
		COALESCE(p.transaction_id,''),
		p.payment_date,
		COALESCE(p.payment_status,'')
	FROM bookings b
	JOIN booking_status s ON b.status_id = s.status_id
	JOIN Users u ON b.user_id = u.user_id
	LEFT JOIN Our_Services rs ON b.service_id = rs.service_id   -- join to get service_description
	LEFT JOIN Restaurant_Staff staff ON staff.staff_id = b.staff_id
	LEFT JOIN payments p ON b.booking_id = p.booking_id
	WHERE b.booking_id = ?
	`

	row := config.DB.QueryRow(query, bookingID)

	var booking restorantmodels.BookingDetailsResponse
	var user usermodels.UserDetailsModel
	var staff restorantmodels.StaffDetails
	var payment usermodels.PaymentResponse

	// Nullable fields
	var startOTP sql.NullString
	var staffID, staffName, staffWhatsapp, staffDesignation, staffEmail, staffDescription, staffImage sql.NullString
	var serviceDescription sql.NullString
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
		&serviceDescription,   // scanned service_description
		&booking.Price,

		&booking.Status,
		&booking.ScheduledDate,
		&booking.Address,
		&booking.Pincode,
		&booking.State,
		&booking.City,
		&booking.Latitude,
		&booking.Longitude,
		&booking.CreatedAt,
		&startOTP,

		&staffID,
		&staffName,
		&staffWhatsapp,
		&staffDesignation,
		&staffEmail,
		&staffDescription,
		&staffImage,

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

	// Assign service description
	if serviceDescription.Valid {
		booking.ServiceDesc = serviceDescription.String
	}

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

	// Decrypt Start OTP if available
	if startOTP.Valid && booking.Status != "Pending" && booking.Status != "Rejected" {
		decryptedOtp, err := utils.DecryptOTP(startOTP.String, key)
		if err == nil {
			booking.StartOtp = decryptedOtp
		} else {
			booking.StartOtp = ""
		}
	} else {
		booking.StartOtp = ""
	}

	return booking, nil
}


// AcceptBooking generates staff OTP and updates booking status to Accepted
func AcceptBooking(bookingID string, acceptedStatusID int, key string) (string, error) {
	otp := utils.GenerateOTP()

	// Encrypt OTP before storing
	otpEncrypted, err := utils.EncryptOTP(otp, key)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt OTP: %w", err)
	}

	now := time.Now()

	query := `
		UPDATE bookings
		SET status_id = ?, start_otp_hash = ?, start_verified = FALSE, updated_at = ?
		WHERE booking_id = ?;
	`

	_, err = config.DB.Exec(query, acceptedStatusID, otpEncrypted, now, bookingID)
	if err != nil {
		return "", err
	}

	// Return plain OTP for sending to staff
	return otp, nil
}

// RejectBooking updates booking status to Rejected with reason
func RejectBooking(bookingID string, rejectedStatusID int, reason string) error {
	query := `
		UPDATE bookings
		SET status_id = ?, cancel_reason = ?, updated_at = NOW()
		WHERE booking_id = ?;
	`
	_, err := config.DB.Exec(query, rejectedStatusID, reason, bookingID)
	if err != nil {
		return fmt.Errorf("failed to reject booking: %w", err)
	}

	return nil
}

// VerifyStartOTP verifies staff OTP and updates booking status to In Progress
func VerifyStartOTP(bookingID string, inputOTP string, inProgressStatusID int, key string) error {
	var storedEncrypted string
	var verified bool

	err := config.DB.QueryRow(`
		SELECT start_otp_hash, start_verified 
		FROM bookings 
		WHERE booking_id = ?`,
		bookingID,
	).Scan(&storedEncrypted, &verified)
	if err != nil {
		return err
	}

	if verified {
		return fmt.Errorf("OTP already verified")
	}

	// Decrypt stored OTP
	storedOTP, err := utils.DecryptOTP(storedEncrypted, key)
	if err != nil {
		return fmt.Errorf("failed to decrypt OTP: %w", err)
	}

	if inputOTP != storedOTP {
		return fmt.Errorf("invalid OTP")
	}

	// Update status to In Progress
	_, err = config.DB.Exec(`
		UPDATE bookings
		SET status_id = ?, start_verified = TRUE, updated_at = ?
		WHERE booking_id = ?`,
		inProgressStatusID, time.Now(), bookingID,
	)
	if err != nil {
		return err
	}

	return nil
}





//assign booking to staff

// AssignStaff assigns a staff member to a booking
func AssignStaff(bookingID string, staffID string) error {
	query := `
		UPDATE bookings
		SET staff_id = ?, updated_at = ?
		WHERE booking_id = ?;
	`
	_, err := config.DB.Exec(query, staffID, time.Now(), bookingID)
	if err != nil {
		return fmt.Errorf("failed to assign staff: %w", err)
	}
	return nil
}