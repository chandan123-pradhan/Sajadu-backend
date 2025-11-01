package paymentservices

import paymentrepo "decoration_project/repository/payment_repo"

type PaymentService struct {
	repo *paymentrepo.RazorpayRepository
}

func NewPaymentService(repo *paymentrepo.RazorpayRepository) *PaymentService {
	return &PaymentService{repo: repo}
}

func (s *PaymentService) Refund(paymentID string,transactionId string, amount int, bookingID string) (map[string]interface{}, error) {
	return s.repo.Refund(paymentID,transactionId, amount, bookingID)
}
