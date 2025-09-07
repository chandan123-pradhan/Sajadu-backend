package staffmodel

import "time"

type Staff struct {
	StaffID      string    `json:"staff_id"`
	RestaurantID string    `json:"restaurant_id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	WhatsappNo   string    `json:"whatsapp_no"`
	Designation  string    `json:"designation"`
	Description  string    `json:"description"`
	Password     string    `json:"password"`
	ImageURL     string    `json:"image_url"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
