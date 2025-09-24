package adminrepo

import (
	"database/sql"
	"decoration_project/config"
	adminmodel "decoration_project/models/admin_model"
	"fmt"
	"strings"
	"time"
)

func GetAllRestorants() ([]adminmodel.RestorantModel, error) {
	var restorants []adminmodel.RestorantModel

	// Fetch all restaurants
	query := `SELECT restaurant_id, name FROM Restaurant_Profile`
	rows, err := config.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r adminmodel.RestorantModel
		if err := rows.Scan(&r.RestorantId, &r.RestorantName); err != nil {
			return nil, err
		}

		// Fetch restaurant images
		imgRows, err := config.DB.Query(`SELECT image_url FROM Restaurant_Images WHERE restaurant_id = ?`, r.RestorantId)
		if err != nil {
			return nil, err
		}

		var images []string
		for imgRows.Next() {
			var url string
			if err := imgRows.Scan(&url); err != nil {
				imgRows.Close()
				return nil, err
			}
			images = append(images, url)
		}
		imgRows.Close()
		if images == nil {
			images = []string{}
		}
		r.Images = images
		restorants = append(restorants, r)
	}

	return restorants, nil
}


func GetAllActiveBookings(statuses []string) ([]adminmodel.BookingResponse, error) {
	query := `
        SELECT 
            b.booking_id,
            b.user_id,
            b.restaurant_id,
            b.service_id,
            b.service_name,
            rs.service_description,
            b.price,
            s.status_name,
            b.scheduled_date,
            b.address,
            b.latitude,
            b.longitude,
            b.pincode,
            b.state,
            b.city,
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
        LEFT JOIN payments p ON b.booking_id = p.booking_id
        LEFT JOIN Our_Services rs ON b.service_id = rs.service_id
    `

	var args []interface{}
	if len(statuses) > 0 {
		placeholders := make([]string, len(statuses))
		for i, status := range statuses {
			placeholders[i] = "?"
			args = append(args, status)
		}
		query += " WHERE s.status_name IN (" + strings.Join(placeholders, ",") + ")"
	}

	query += " ORDER BY b.created_at DESC"

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []adminmodel.BookingResponse

	for rows.Next() {
		var b adminmodel.BookingResponse
		var paymentID, transactionID, paymentMethod, currency, paymentStatus sql.NullString
		var paymentDate sql.NullTime
		var amount sql.NullFloat64
		var serviceDesc sql.NullString
		var restaurantID sql.NullString // <-- nullable

		err := rows.Scan(
			&b.BookingID,
			&b.UserID,
			&restaurantID, // scan into NullString
			&b.ServiceID,
			&b.ServiceName,
			&serviceDesc,
			&b.Price,
			&b.Status,
			&b.ScheduledDate,
			&b.Address,
			&b.Latitude,
			&b.Longitude,
			&b.Pincode,
			&b.State,
			&b.City,
			&b.CreatedAt,
			&paymentID,
			&amount,
			&currency,
			&paymentMethod,
			&transactionID,
			&paymentDate,
			&paymentStatus,
		)
		if err != nil {
			return nil, err
		}

		// Assign restaurant ID only if valid
		if restaurantID.Valid {
			b.RestaurantID = restaurantID.String
		} else {
			b.RestaurantID = ""
		}

		// Service description
		if serviceDesc.Valid {
			b.ServiceDesc = serviceDesc.String
		}

		// Payment
		if paymentID.Valid {
			b.Payment = &adminmodel.PaymentResponse{
				PaymentID:     paymentID.String,
				Amount:        amount.Float64,
				Currency:      currency.String,
				PaymentMethod: paymentMethod.String,
				TransactionID: transactionID.String,
				PaymentStatus: paymentStatus.String,
			}
			if paymentDate.Valid {
				b.Payment.PaymentDate = paymentDate.Time
			}
		} else {
			b.Payment = nil
		}

		// Fetch service images
		imgQuery := `SELECT image_url FROM Service_Images WHERE service_id = ?`
		imgRows, err := config.DB.Query(imgQuery, b.ServiceID)
		if err == nil {
			var images []string
			defer imgRows.Close()
			for imgRows.Next() {
				var img string
				if err := imgRows.Scan(&img); err == nil {
					images = append(images, img)
				}
			}
			b.Images = images
		}

		bookings = append(bookings, b)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return bookings, nil
}




func GetBookingDetailsByID(bookingID string) (adminmodel.BookingResponse, error) {
	var booking adminmodel.BookingResponse

	query := `
	SELECT 
		b.booking_id,
		b.user_id,
		b.restaurant_id,
		b.service_id,
		srv.service_name,
		srv.service_description,
		srv.service_price,
		b.status_id,
		b.scheduled_date,
		b.address,
		b.latitude,
		b.longitude,
		b.pincode,
		b.state,
		b.city,
		b.created_at,
		bs.status_name
	FROM bookings b
	JOIN booking_status bs ON b.status_id = bs.status_id
	JOIN Our_Services srv ON b.service_id = srv.service_id
	WHERE b.booking_id = ?
	LIMIT 1
	`

	var restaurantID sql.NullString
	var latitude, longitude sql.NullFloat64
	var servicePrice sql.NullFloat64
	var scheduledDate sql.NullTime
	var statusName sql.NullString
	var serviceDescription sql.NullString

	err := config.DB.QueryRow(query, bookingID).Scan(
		&booking.BookingID,
		&booking.UserID,
		&restaurantID,
		&booking.ServiceID,
		&booking.ServiceName,
		&booking.ServiceDesc,
		&servicePrice,
		new(interface{}), // status_id, not needed
		&scheduledDate,
		&booking.Address,
		&latitude,
		&longitude,
		&booking.Pincode,
		&booking.State,
		&booking.City,
		&booking.CreatedAt,
		&statusName,
	)
	if err != nil {
		return booking, err
	}

	// Map nullable fields
	if restaurantID.Valid {
		booking.RestaurantID = restaurantID.String
	} else {
		booking.RestaurantID = ""
	}

	if scheduledDate.Valid {
		booking.ScheduledDate = scheduledDate.Time.Format("2006-01-02 15:04:05")
	} else {
		booking.ScheduledDate = ""
	}

	if latitude.Valid {
		booking.Latitude = latitude.Float64
	}
	if longitude.Valid {
		booking.Longitude = longitude.Float64
	}

	if servicePrice.Valid {
		booking.Price = servicePrice.Float64
	} else {
		booking.Price = 0
	}

	booking.Status = ""
	if statusName.Valid {
		booking.Status = statusName.String
	}

	// Service description
	if serviceDescription.Valid {
		booking.ServiceDesc = serviceDescription.String
	}

	// Fetch service images
	booking.Images = []string{}
	imgRows, err := config.DB.Query(`SELECT image_url FROM Service_Images WHERE service_id = ?`, booking.ServiceID)
	if err == nil {
		defer imgRows.Close()
		for imgRows.Next() {
			var url string
			if err := imgRows.Scan(&url); err == nil {
				booking.Images = append(booking.Images, url)
			}
		}
	}

	// Fetch payment info
	var payment adminmodel.PaymentResponse
	var paymentID, currency, paymentMode, transactionID, paymentStatus sql.NullString
	var amount sql.NullFloat64
	var paymentDate sql.NullTime

	paymentQuery := `
	SELECT payment_id, amount, currency, payment_mode, transaction_id, payment_status, payment_date
	FROM payments
	WHERE booking_id = ?
	ORDER BY payment_date DESC
	LIMIT 1
	`
	err = config.DB.QueryRow(paymentQuery, bookingID).Scan(
		&paymentID, &amount, &currency, &paymentMode, &transactionID, &paymentStatus, &paymentDate,
	)
	if err == nil && paymentID.Valid {
		payment.PaymentID = paymentID.String
		if amount.Valid {
			payment.Amount = amount.Float64
		}
		payment.Currency = currency.String
		payment.PaymentMethod = paymentMode.String
		payment.TransactionID = transactionID.String
		payment.PaymentStatus = paymentStatus.String
		if paymentDate.Valid {
			payment.PaymentDate = paymentDate.Time
		}
		booking.Payment = &payment
	}

	// Fetch restaurant profile if restaurant_id exists
	if booking.RestaurantID != "" {
		var restaurantName sql.NullString
		restaurantQuery := `SELECT name FROM Restaurant_Profile WHERE restaurant_id = ? LIMIT 1`
		err := config.DB.QueryRow(restaurantQuery, booking.RestaurantID).Scan(&restaurantName)
		if err == nil && restaurantName.Valid {
			booking.RestaurantName = restaurantName.String
		}

		// Fetch restaurant images
		booking.RestaurantImages = []string{}
		imgRows, err := config.DB.Query(`SELECT image_url FROM Restaurant_Images WHERE restaurant_id = ?`, booking.RestaurantID)
		if err == nil {
			defer imgRows.Close()
			for imgRows.Next() {
				var url string
				if err := imgRows.Scan(&url); err == nil {
					booking.RestaurantImages = append(booking.RestaurantImages, url)
				}
			}
		}
	}

	return booking, nil
}



func UpdateBookingStatus(bookingID, restaurantID, newStatus string) error {
	// Start transaction
	tx, err := config.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	// Get status_id from booking_status table
	var statusID int
	err = tx.QueryRow(`SELECT status_id FROM booking_status WHERE status_name = ?`, newStatus).Scan(&statusID)
	if err != nil {
		return fmt.Errorf("failed to get status id: %v", err)
	}

	// Update booking
	if newStatus == "Accepted" {
		_, err = tx.Exec(`
			UPDATE bookings 
			SET status_id = ?, restaurant_id = ?
			WHERE booking_id = ?
		`, statusID, restaurantID, bookingID)
	} else if newStatus == "Rejected" {
		_, err = tx.Exec(`
			UPDATE bookings 
			SET status_id = ?
			WHERE booking_id = ?
		`, statusID, bookingID)
	} else {
		return fmt.Errorf("invalid status: %s", newStatus)
	}
	if err != nil {
		return fmt.Errorf("failed to update booking: %v", err)
	}

	// Optional: Insert log for audit trail
	_, err = tx.Exec(`
		INSERT INTO booking_logs (log_id, booking_id, old_status, new_status, changed_by, created_at)
		SELECT UUID(), ?, b.status_id, ?, 'Restaurant', ?
		FROM bookings b WHERE b.booking_id = ?
	`, bookingID, statusID, time.Now(), bookingID)
	if err != nil {
		return fmt.Errorf("failed to insert booking log: %v", err)
	}

	return nil
}