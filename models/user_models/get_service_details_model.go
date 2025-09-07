package usermodels

type GetServiceDetails struct {
	RestaurantId string `json:"restaurant_id"`
	ServiceId    string `json:"service_id"`
}

var GetServiceDetailsRequest GetServiceDetails
