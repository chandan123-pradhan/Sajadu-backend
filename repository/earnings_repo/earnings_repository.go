package earningsrepo

import (
	"decoration_project/config"
	"decoration_project/models/earnings"
)



func GetEarnings() (
	earnings.EarningsResponse,
	error,
) {

	var response earnings.EarningsResponse
	query := `
	SELECT
		COALESCE(
			SUM(
				CASE
					WHEN bs.status_name='Completed'
					THEN b.price
					ELSE 0
				END
			),0
		) total_earning,

		COALESCE(
			SUM(
				CASE
					WHEN bs.status_name='Completed'
					AND MONTH(b.created_at)=MONTH(CURRENT_DATE())
					AND YEAR(b.created_at)=YEAR(CURRENT_DATE())
					THEN b.price
					ELSE 0
				END
			),0
		) month_earning,

		COALESCE(
			SUM(
				CASE
					WHEN bs.status_name='Completed'
					AND DATE(b.created_at)=CURRENT_DATE()
					THEN b.price
					ELSE 0
				END
			),0
		) today_earning

	FROM bookings b
	INNER JOIN booking_status bs
	ON b.status_id=bs.status_id
	`

	err := config.DB.QueryRow(query).Scan(
		&response.TotalEarning,
		&response.MonthEarning,
		&response.TodayEarning,
	)

	if err != nil {
		return response, err
	}

	return response, nil
}