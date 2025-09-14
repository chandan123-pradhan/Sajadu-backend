package commonmodel
// package commonmodels

// import "time/"

type PaymentResponseBasic struct {
   PaymentID string  `json:"payment_id"`
   Amount    float64 `json:"amount"`
   Status    string  `json:"status"`
}

type UserBasic struct {
   UserID   string `json:"user_id"`
   Name     string `json:"name"`
   Phone    string `json:"phone"`
   Email    string `json:"email"`
}

type StaffBasic struct {
    StaffID *string `json:"staff_id,omitempty"`
    Name    *string `json:"name,omitempty"`
    Phone   *string `json:"phone,omitempty"`
    Images  *string `json:"images,omitempty"`
}

type RestorantBasic struct{
	RestorantId string `json:"restorant_id"`
	RestorantName string `json:"restorant_name"`

}
