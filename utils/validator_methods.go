package utils

import (
	usermodels "decoration_project/models/user_models"
	"fmt"
	"io"
	"mime/multipart"
)

// ValidateImageSize ensures uploaded file size is within limit
func ValidateImageSize(file multipart.File, maxSize int64) error {
	// Read into buffer to check size
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file")
	}

	if int64(len(fileBytes)) > maxSize {
		return fmt.Errorf("image cannot be more than %d MB", maxSize/(1<<20))
	}

	// Reset reader for re-use
	file.Seek(0, 0)
	return nil
}

func ValidateBookingRequest(req usermodels.BookingRequest) error {
	if req.RestaurantID == "" {
		return fmt.Errorf("Restaurant ID is required")
	}
	if req.ServiceID == "" {
		return fmt.Errorf("Service ID is required")
	}
	if req.Address == "" {
		return fmt.Errorf("Address is required")
	}
	if req.Latitude == 0 || req.Longitude == 0 {
		return fmt.Errorf("Valid latitude and longitude are required")
	}
	if req.Pincode == "" {
		return fmt.Errorf("Pincode is required")
	}
	if req.ServiceName == "" {
		return fmt.Errorf("Service name is required")
	}
	if req.State == "" {
		return fmt.Errorf("State is required")
	}
	if req.City == "" {
		return fmt.Errorf("City is required")
	}
	if req.Price <= 0 {
		return fmt.Errorf("Price must be greater than 0")
	}
	if req.ScheduledDate.IsZero() {
		return fmt.Errorf("Scheduled date is required")
	}

	// ✅ Extra: Payment validations
	if req.PaymentMethod != "COD" {
		// If user provided payment details, validate them
		if req.TransactionID == "" {
			return fmt.Errorf("Transaction ID is required when payment method is provided")
		}
		if req.AmountPaid <= 0 {
			return fmt.Errorf("AmountPaid must be greater than 0 when payment method is provided")
		}
		if req.Currency == "" {
			return fmt.Errorf("Currency is required when payment method is provided")
		}
	}

	return nil
}
