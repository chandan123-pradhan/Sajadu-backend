package adminmodel

type Festival struct {
    FestivalID     string `json:"festival_id"`
    FestivalName   string `json:"festival_name"`
    StartDate      string `json:"start_date"`
    EndDate        string `json:"end_date"`
    CreatedAt      string `json:"created_at"`
}
