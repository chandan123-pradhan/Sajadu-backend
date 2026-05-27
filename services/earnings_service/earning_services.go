package earningsservice

import (
	"decoration_project/models/earnings"
	earningsrepo "decoration_project/repository/earnings_repo"
)


type EarningsService struct{}

func NewEarningsService() *EarningsService {
	return &EarningsService{}
}

func (s *EarningsService) GetEarnings() (
	earnings.EarningsResponse,
	error,
) {
	return earningsrepo.GetEarnings()
}