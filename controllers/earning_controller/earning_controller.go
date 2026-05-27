package earningcontroller

import (
	"net/http"

	earningsservice "decoration_project/services/earnings_service"
	"decoration_project/utils"
)

type EarningsController struct {
	Service *earningsservice.EarningsService
}

func NewEarningsController() *EarningsController {
	service := earningsservice.NewEarningsService()

	return &EarningsController{
		Service: service,
	}
}

func (c *EarningsController) GetEarningsHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	earnings, err := c.Service.GetEarnings()
	if err != nil {
		utils.SendResponse(
			w,
			http.StatusInternalServerError,
			false,
			nil,
			err.Error(),
		)
		return
	}

	utils.SendResponse(
		w,
		http.StatusOK,
		true,
		earnings,
		"Earnings fetched successfully",
	)
}