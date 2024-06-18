package usecases

import (
	userUC "github.com/qaultsabit/wallet/src/app/usecases/user"
)

type AllUseCases struct {
	UserUC userUC.UserUCInterface
	// BalanceUC     balanceUC.BalanceUCInterface
	// TransactionUC transactionUC.TransactionUCInterface
}
