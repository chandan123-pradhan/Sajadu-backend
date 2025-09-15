package staffrepo

import (
	"database/sql"
	"decoration_project/config"
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


// GetAssignedServiceDetails fetches full details of a specific booking (without OTP & Staff)
func GetAssignedServiceDetails(bookingID string) (staffmodel.StaffAssignedServicesDetails, error) {
	query := `
	SELECT 
		b.booking_id,
		u.user_id,
		u.full_name,
		u.email,
		u.mobile_number,

		b.service_id,
		b.service_name,
		rs.service_description,
		COALESCE(rs.service_price, 0),
		COALESCE(rs.items, JSON_ARRAY()),

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
	JOIN Restaurant_Services rs ON b.service_id = rs.service_id
	LEFT JOIN payments p ON b.booking_id = p.booking_id
	WHERE b.booking_id = ?
	`

	row := config.DB.QueryRow(query, bookingID)

	var booking staffmodel.StaffAssignedServicesDetails
	var user usermodels.UserDetailsModel
	var payment usermodels.PaymentResponse

	// Nullable fields
	var serviceDesc sql.NullString
	var items sql.NullString
	var paymentID, transactionID, paymentMethod, currency, paymentStatus sql.NullString
	var paymentDate sql.NullTime
	var amount sql.NullFloat64

	err := row.Scan(
		&booking.BookingID,
		&user.UserID,
		&user.FullName,
		&user.Email,
		&user.MobileNumber,

		&booking.ServiceID,
		&booking.ServiceName,
		&serviceDesc,
		&booking.ServicePrice,
		&items,

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
		return staffmodel.StaffAssignedServicesDetails{}, err
	}

	// Assign user
	booking.User = user

	// Service description
	if serviceDesc.Valid {
		booking.ServiceDescription = serviceDesc.String
	}

	// Items (JSON stored as string)
	if items.Valid {
		booking.Items = items.String
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

	// Fetch service images
	imgQuery := `SELECT image_url FROM Service_Images WHERE service_id = ?`
	rows, err := config.DB.Query(imgQuery, booking.ServiceID)
	if err == nil {
		defer rows.Close()
		var images []string
		for rows.Next() {
			var img string
			if err := rows.Scan(&img); err == nil {
				images = append(images, img)
			}
		}
		booking.Images = images
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


func SaveStaffLocation(staffID, bookingID string, latitude, longitude float64) error {
    query := `
        INSERT INTO service_partner_location (location_id, staff_id, booking_id, latitude, longitude)
        VALUES (UUID(), ?, ?, ?, ?)
        ON DUPLICATE KEY UPDATE latitude=?, longitude=?, updated_at=NOW()
    `
    _, err := config.DB.Exec(query, staffID, bookingID, latitude, longitude, latitude, longitude)
    return err
}



func FetchPartnerLocation(bookingID string) (staffmodel.PartnerLocationResponse, error) {
    var location staffmodel.PartnerLocationResponse

    query := `
        SELECT spl.staff_id, spl.booking_id, spl.latitude, spl.longitude, spl.updated_at
        FROM service_partner_location spl
        WHERE spl.booking_id = ?
        ORDER BY spl.updated_at DESC
        LIMIT 1
    `

    err := config.DB.QueryRow(query, bookingID).Scan(
        &location.PartnerID,
        &location.BookingID,
        &location.Latitude,
        &location.Longitude,
        &location.UpdatedAt,
    )

    if err == sql.ErrNoRows {
        return location, errors.New("no location found for this booking")
    } else if err != nil {
        return location, err
    }

    return location, nil
}

// VerifyToCompleteService verifies the completion OTP and updates booking status to "Completed"
func VerifyToCompleteService(bookingID string, staffOTP string, key string) error {
	// 1. Fetch encrypted completion OTP, current status, and completion_verified flag
	var encryptedCompleteOTP string
	var statusID int
	var completeVerified bool

	query := `SELECT complete_otp_hash, status_id, complete_verified 
	          FROM bookings WHERE booking_id = ?`
	err := config.DB.QueryRow(query, bookingID).Scan(&encryptedCompleteOTP, &statusID, &completeVerified)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("booking not found")
		}
		return err
	}

	// 2. Check if already completed
	if completeVerified || statusID == 5 {
		return errors.New("service already completed or OTP already verified")
	}

	// 3. Decrypt OTP
	decryptedOTP, err := utils.DecryptOTP(encryptedCompleteOTP, key)
	if err != nil {
		return errors.New("failed to decrypt completion OTP")
	}

	// 4. Compare OTP
	if decryptedOTP != staffOTP {
		return errors.New("invalid OTP")
	}

	// 5. Update booking status to Completed (status_id = 5), mark completion_verified
	updateQuery := `
		UPDATE bookings 
		SET status_id = 5, complete_verified = TRUE, updated_at = CURRENT_TIMESTAMP
		WHERE booking_id = ?
	`
	_, err = config.DB.Exec(updateQuery, bookingID)
	if err != nil {
		return err
	}

	return nil
}
