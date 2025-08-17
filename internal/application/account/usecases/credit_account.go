package usecases

import (
	domainAccount "github.com/juninhoitabh/clob-go/internal/domain/account"
)

type (
	CreditAccountInput struct {
		AccountID string
		Asset     string
		Amount    int64
	}
	CreditAccountUseCase struct {
		accountRepo domainAccount.IAccountRepository
	}
)

func (c *CreditAccountUseCase) Execute(input CreditAccountInput) error {
	acct, err := c.accountRepo.Get(input.AccountID)
	if err != nil {
		return err
	}

	acct.Credit(input.Asset, input.Amount)

	if err := c.accountRepo.Save(acct); err != nil {
		return err
	}

	return nil
}
