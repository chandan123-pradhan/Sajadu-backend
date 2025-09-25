package userrepo

import (
	"database/sql"
	"decoration_project/config"
	commonmodel "decoration_project/models/common_model"
	// staffmodel "decoration_project/models/staff_model"
	usermodels "decoration_project/models/user_models"
	"decoration_project/utils"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// CreateBooking inserts a new booking into the DB
func CreateBooking(req usermodels.BookingRequest) (usermodels.BookingResponse, error) {
	bookingID := uuid.New().String()

	formattedDate := req.ScheduledDate.Format("2006-01-02 15:04:05")

	// ✅ Start transaction (to ensure atomic insert for booking + payment)
	tx, err := config.DB.Begin()
	if err != nil {
		return usermodels.BookingResponse{}, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()
	var statusID int
	statusName := "Pending"

	err = config.DB.QueryRow(
		"SELECT status_id FROM booking_status WHERE status_name = ?",
		statusName,
	).Scan(&statusID)
	if err != nil {
		return usermodels.BookingResponse{}, fmt.Errorf("failed to get status id: %v", err)
	}
	// Insert booking
	bookingQuery := `
		INSERT INTO bookings 
		(booking_id, user_id, service_id, status_id, scheduled_date, address, latitude, longitude, pincode, state, city, service_name, price) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err = tx.Exec(bookingQuery,
		bookingID,
		req.UserID,
		req.ServiceID,
		statusID,
		formattedDate,
		req.Address,
		sql.NullFloat64{Float64: req.Latitude, Valid: req.Latitude != 0},
		sql.NullFloat64{Float64: req.Longitude, Valid: req.Longitude != 0},
		req.Pincode,
		req.State,
		req.City,
		req.ServiceName,
		req.Price,
	)
	if err != nil {
		return usermodels.BookingResponse{}, err
	}

	// If payment details are provided → insert payment
	var paymentID string
	if req.PaymentMethod != "" {
		paymentID = uuid.New().String()
		paymentQuery := `
			INSERT INTO payments 
			(payment_id, booking_id, amount, currency, payment_mode, transaction_id, payment_date, payment_status) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`
		_, err = tx.Exec(paymentQuery,
			paymentID,
			bookingID,
			req.AmountPaid,
			req.Currency,
			req.PaymentMethod,
			req.TransactionID,
			time.Now().Format("2006-01-02 15:04:05"),
			"Success", // assuming success, update later if integrating gateway
		)
		if err != nil {
			return usermodels.BookingResponse{}, err
		}
	}

	// Build response object
	res := usermodels.BookingResponse{
		BookingID:     bookingID,
		UserID:        req.UserID,
		ServiceID:     req.ServiceID,
		Status:        "Pending", // mapped from status_id
		ScheduledDate: formattedDate,
		Address:       req.Address,
		Pincode:       req.Pincode,
		State:         req.State,
		City:          req.City,
		Price:         req.Price,
		ServiceName: req.ServiceName,
		CreatedAt:     time.Now(),
	}

	// If payment added, include in response
	if paymentID != "" {
		res.Payment = &usermodels.PaymentResponse{
			PaymentID:     paymentID,
			Amount:        req.AmountPaid,
			Currency:      req.Currency,
			PaymentMethod: req.PaymentMethod,
			TransactionID: req.TransactionID,
			PaymentDate:   time.Now(),
			PaymentStatus: "Success",
		}
	}

	return res, nil
}


// GetUserBookings fetches all bookings of a given user
func GetUserBookings(userID, otpKey string) (usermodels.UserBookingsWrapper, error) {
	query := `
	SELECT 
		b.booking_id,
		b.user_id,
		b.restaurant_id,
		b.service_id,
		s.status_name,
		b.scheduled_date,
		b.address,
		b.pincode,
		b.state,
		b.city,
		b.service_name,
		b.created_at,
		b.start_verified,
		b.complete_otp_hash,
		b.price,

		-- Payment details (may be NULL if no payment yet)
		MAX(p.payment_id),
		COALESCE(MAX(p.amount), 0),
		COALESCE(MAX(p.currency), ''),
		COALESCE(MAX(p.payment_mode), ''),
		COALESCE(MAX(p.transaction_id), ''),
		MAX(p.payment_date),
		COALESCE(MAX(p.payment_status), ''),

		-- ✅ Fetch service images (aggregated)
		COALESCE(GROUP_CONCAT(si.image_url), '')
	FROM bookings b
	JOIN booking_status s ON b.status_id = s.status_id
	LEFT JOIN payments p ON b.booking_id = p.booking_id
	LEFT JOIN Service_Images si ON b.service_id = si.service_id
	WHERE b.user_id = ?
	GROUP BY b.booking_id
	ORDER BY b.created_at DESC
`

	rows, err := config.DB.Query(query, userID)
	if err != nil {
		return usermodels.UserBookingsWrapper{Bookings: []usermodels.BookingResponse{}}, err
	}
	defer rows.Close()

	var bookings []usermodels.BookingResponse

	for rows.Next() {
		var booking usermodels.BookingResponse
		var payment usermodels.PaymentResponse
		var paymentID, transactionID, paymentMethod, currency, paymentStatus, completeOtpHash sql.NullString
		var paymentDate sql.NullTime
		var amount sql.NullFloat64
		var startVerified bool
		var serviceImages sql.NullString
        var restaurantID sql.NullString
		err := rows.Scan(
			&booking.BookingID,
			&booking.UserID,
			&restaurantID,
			&booking.ServiceID,
			&booking.Status,
			&booking.ScheduledDate,
			&booking.Address,
			&booking.Pincode,
			&booking.State,
			&booking.City,
			&booking.ServiceName,
			&booking.CreatedAt,
			&startVerified,
			&completeOtpHash,
			&booking.Price,
			&paymentID,
			&amount,
			&currency,
			&paymentMethod,
			&transactionID,
			&paymentDate,
			&paymentStatus,
			&serviceImages, // NEW
		)
		if err != nil {
			return usermodels.UserBookingsWrapper{Bookings: []usermodels.BookingResponse{}}, err
		}

		
		// If payment exists, populate PaymentResponse
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

		// Assign decrypted complete OTP only if start_verified is true
		if startVerified && completeOtpHash.Valid {
			decryptedOtp, err := utils.DecryptOTP(completeOtpHash.String, otpKey)
			if err == nil {
				booking.CompleteOtp = decryptedOtp
			} else {
				booking.CompleteOtp = ""
			}
		}

		// Parse service images
		if serviceImages.Valid && serviceImages.String != "" {
			booking.Images = strings.Split(serviceImages.String, ",")
		} else {
			booking.Images = []string{}
		}

		bookings = append(bookings, booking)
	}

	if err = rows.Err(); err != nil {
		return usermodels.UserBookingsWrapper{Bookings: []usermodels.BookingResponse{}}, err
	}

	return usermodels.UserBookingsWrapper{Bookings: bookings}, nil
}




// GetBookingDetails fetches full details of a specific booking (without OTP & Staff mandatory)
func GetBookingDetails(bookingID string) (usermodels.GetUsersBookingDetailsResponse, error) {
	query := `
	SELECT 
		b.booking_id,
		u.staff_id,
		u.name,
		u.whatsapp_no,
		u.image_url,

		b.service_id,
		b.service_name,
		rs.service_description,
		COALESCE(rs.service_price, 0),

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
	LEFT JOIN Restaurant_Staff u ON b.staff_id = u.staff_id
	LEFT JOIN Our_Services rs ON b.service_id = rs.service_id
	LEFT JOIN payments p ON b.booking_id = p.booking_id
	WHERE b.booking_id = ?
	`

	row := config.DB.QueryRow(query, bookingID)

	var booking usermodels.GetUsersBookingDetailsResponse
	var payment usermodels.PaymentResponse

	// Nullable fields
	var serviceDesc sql.NullString
	var paymentID, transactionID, paymentMethod, currency, paymentStatus sql.NullString
	var paymentDate sql.NullTime
	var amount sql.NullFloat64

	// Staff nullable
	// var staff commonmodel.StaffBasic
	var staffID, staffName, staffPhone, staffImg sql.NullString

	err := row.Scan(
		&booking.BookingID,
		&staffID,
		&staffName,
		&staffPhone,
		&staffImg,

		&booking.ServiceID,
		&booking.ServiceName,
		&serviceDesc,
		&booking.ServicePrice,

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
		return usermodels.GetUsersBookingDetailsResponse{}, err
	}

if staffID.Valid {
		booking.User = commonmodel.StaffBasic{
			StaffID: &staffID.String,
			Name:    &staffName.String,
			Phone:   &staffPhone.String,
			Images:  &staffImg.String,
		}
	} else {
		fmt.Println("emtpy")
		// Staff not assigned → empty object
		booking.User = commonmodel.StaffBasic{}
	}

	// Service description
	if serviceDesc.Valid {
		booking.ServiceDescription = serviceDesc.String
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
