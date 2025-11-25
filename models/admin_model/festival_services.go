package adminmodel

type FestivalService struct {
    ID         string `json:"id"`
    FestivalID string `json:"festival_id"`
    ServiceID  string `json:"service_id"`
    CreatedAt  string `json:"created_at"`
}

type AddFestivalServiceRequest struct {
    FestivalID      string  `json:"festival_id"`
    ServiceID       string  `json:"service_id"`
    DiscountPercent float64 `json:"discount_percent"`
}