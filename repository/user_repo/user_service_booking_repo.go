package userrepo

import (
	"database/sql"
	"decoration_project/config"
	usermodels "decoration_project/models/user_models"
	"decoration_project/utils"
	"fmt"
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
		(booking_id, user_id, restaurant_id, service_id, status_id, scheduled_date, address, latitude, longitude, pincode, state, city, service_name) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err = tx.Exec(bookingQuery,
		bookingID,
		req.UserID,
		req.RestaurantID,
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
		RestaurantID:  req.RestaurantID,
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
		WHERE b.user_id = ?
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

		err := rows.Scan(
			&booking.BookingID,
			&booking.UserID,
			&booking.RestaurantID,
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
			&paymentID,
			&amount,
			&currency,
			&paymentMethod,
			&transactionID,
			&paymentDate,
			&paymentStatus,
		)
		if err != nil {
			return usermodels.UserBookingsWrapper{Bookings: []usermodels.BookingResponse{}}, err
		}

		// Map price from payment amount
		if amount.Valid {
			booking.Price = amount.Float64
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

		bookings = append(bookings, booking)
	}

	if err = rows.Err(); err != nil {
		return usermodels.UserBookingsWrapper{Bookings: []usermodels.BookingResponse{}}, err
	}

	return usermodels.UserBookingsWrapper{Bookings: bookings}, nil
}
