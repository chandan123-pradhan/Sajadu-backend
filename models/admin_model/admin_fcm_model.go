package adminmodel

type AdminFCM struct {
	ID        int    `json:"id"`
	FCMToken  string `json:"fcm_token"`
	UpdatedAt string `json:"updated_at"`
}
