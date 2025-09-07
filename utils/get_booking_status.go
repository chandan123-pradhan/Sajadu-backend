package utils

import (
	"database/sql"
	"decoration_project/config"
	"fmt"
)

// GetStatusID fetches the status_id for a given status_name from booking_status table
func GetStatusID(statusName string) int {
	var statusID int

	err := config.DB.QueryRow(`SELECT status_id FROM booking_status WHERE status_name = ?`, statusName).Scan(&statusID)
	if err != nil {
		if err == sql.ErrNoRows {
			// If status not found, panic or handle according to your app logic
			panic(fmt.Sprintf("Status '%s' not found in booking_status table", statusName))
		}
		panic(fmt.Sprintf("Failed to fetch status_id: %v", err))
	}

	return statusID
}
