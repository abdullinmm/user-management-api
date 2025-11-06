package usecases

import (
	"github.com/abdullinmm/user-management-api/internal/repository"
)

type UserUseCase struct {
	userRepo        repository.UserRepository
	balanceRepo     repository.BalanceRepository
	transactionRepo repository.TransactionRepository
	referralBonus   int64
	refereeBonus    int64
}

func NewUserUseCase(userRepo repository.UserRepository, balanceRepo repository.BalanceRepository, transactionRepo repository.TransactionRepository, referralBonus int64, refereeBonus int64) *UserUseCase {
	return &UserUseCase{
		userRepo:        userRepo,
		balanceRepo:     balanceRepo,
		transactionRepo: transactionRepo,
		referralBonus:   referralBonus,
		refereeBonus:    refereeBonus,
	}
}
