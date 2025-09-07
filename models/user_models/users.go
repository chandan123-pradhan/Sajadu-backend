package usermodels

type User struct {
    UserID       string `json:"user_id"`
    FullName     string `json:"full_name"`
    Email        string `json:"email"`
    Password     string `json:"password"`        // allow decoding from JSON
    MobileNumber string `json:"mobile_number,omitempty"`
}




type UserDetailsModel struct {
    UserID       string `json:"user_id"`
    FullName     string `json:"full_name"`
    Email        string `json:"email"`      // allow decoding from JSON
    MobileNumber string `json:"mobile_number,omitempty"`
}

