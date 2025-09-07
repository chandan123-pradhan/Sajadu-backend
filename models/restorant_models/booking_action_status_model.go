package restorantmodels

type BookingActionRequest struct {
	BookingID string `json:"booking_id"`
	Status    string `json:"status"`        // "accept" or "reject"
	Reason    string `json:"reason"`        // optional, for rejection
}