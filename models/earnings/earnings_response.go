package earnings

type EarningsResponse struct {
	TotalEarning float64 `json:"total_earning"`
	MonthEarning float64 `json:"month_earning"`
	TodayEarning float64 `json:"today_earning"`
}