package usecases

import (
	domainAccount "github.com/juninhoitabh/clob-go/internal/domain/account"
)

type (
	CreditAccountUseCase struct {
		accountRepo domainAccount.IAccountRepository
	}
)

func (c *CreditAccountUseCase) Execute(input CreditAccountInput) error {
	acct, err := c.accountRepo.Get(input.AccountID)
	if err != nil {
		return err
	}

	err = acct.Credit(input.Asset, input.Amount)
	if err != nil {
		return err
	}

	if err := c.accountRepo.Save(acct); err != nil {
		return err
	}

	return nil
}

func NewCreditAccountUseCase(
	accountRepo domainAccount.IAccountRepository,
) *CreditAccountUseCase {
	return &CreditAccountUseCase{
		accountRepo: accountRepo,
	}
}
