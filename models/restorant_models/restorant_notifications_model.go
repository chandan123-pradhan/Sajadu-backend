package restorantmodels

type RestorantNotifications struct {
	ID          string `json:"id"`
	RestorantId      string `json:"restorant_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
	CreatedAt   string `json:"created_at"`
}
