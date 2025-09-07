package staffmodel

type VerifyOTPRequest struct {
	BookingID string `json:"booking_id"`
	OTP       string `json:"otp"`
}
