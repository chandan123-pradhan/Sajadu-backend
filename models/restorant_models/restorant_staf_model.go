package restorantmodels

import "time"

type RestaurantStaff struct {
	StaffID      string    `json:"staff_id"`
	RestaurantID string    `json:"restaurant_id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	WhatsappNo   string    `json:"whatsapp_no,omitempty"`
	Designation  string    `json:"designation"` // "Decorator" or "Staff"
	Description  string    `json:"description,omitempty"`
	Password     string    `json:"password,omitempty"`  // hashed password
	ImageURL     string    `json:"image_url,omitempty"` // profile image path
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
