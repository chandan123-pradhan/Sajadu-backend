package usermodels

import "time"

type ServiceReview struct {
	ReviewID   string    `json:"review_id"`
	ServiceID  string    `json:"service_id"`
	UserID     string    `json:"user_id"`
	UserName   string    `json:"user_name"`
	Rating     int       `json:"rating"`
	ReviewText string    `json:"review_text"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
