package restorantmodels

type ServiceWithRestaurant struct {
	Service    RestaurantService `json:"service"`
	Restaurant RestaurantProfile `json:"restaurant"`
}
