package usecases

import (
	domainAccount "github.com/juninhoitabh/clob-go/internal/domain/account"
	idObjValue "github.com/juninhoitabh/clob-go/internal/shared/domain/value-objects/id"
)

type (
	CreateAccountUseCase struct {
		accountRepo domainAccount.IAccountRepository
	}
)

func (c *CreateAccountUseCase) Execute(input CreateAccountInput) (*CreateAccountOutput, error) {
	account, err := domainAccount.NewAccount(domainAccount.AccountProps{
		Name: input.AccountName,
	}, idObjValue.Uuid)
	if err != nil {
		return nil, err
	}

	err = c.accountRepo.Create(account)
	if err != nil {
		return nil, err
	}

	return &CreateAccountOutput{
		ID:   account.GetID(),
		Name: account.Name,
	}, nil
}

func NewCreateAccountUseCase(
	accountRepo domainAccount.IAccountRepository,
) *CreateAccountUseCase {
	return &CreateAccountUseCase{
		accountRepo: accountRepo,
	}
}
