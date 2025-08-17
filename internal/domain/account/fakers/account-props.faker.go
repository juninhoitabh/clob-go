package fakers

import (
	faker "github.com/brianvoe/gofakeit/v7"

	"github.com/juninhoitabh/clob-go/internal/domain/account"
)

func AccountPropsFaker() account.AccountProps {
	faker := faker.New(0)

	return account.AccountProps{
		Name: faker.Name(),
	}
}
