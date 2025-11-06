package usecases

import (
	"github.com/abdullinmm/user-management-api/internal/repositories"
)

type UserUseCase struct {
	userRepo        repositories.UserRepository
	balanceRepo     repositories.BalanceRepository
	transactionRepo repositories.TransactionRepository
	referralBonus   int64
	refereeBonus    int64
}

func NewUserUseCase(userRepo repositories.UserRepository, balanceRepo repositories.BalanceRepository, transactionRepo repositories.TransactionRepository, referralBonus int64, refereeBonus int64) *UserUseCase {
	return &UserUseCase{
		userRepo:        userRepo,
		balanceRepo:     balanceRepo,
		transactionRepo: transactionRepo,
		referralBonus:   referralBonus,
		refereeBonus:    refereeBonus,
	}
}
