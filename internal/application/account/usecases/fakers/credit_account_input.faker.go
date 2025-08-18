package fakers

import (
	faker "github.com/brianvoe/gofakeit/v7"

	accountUsecases "github.com/juninhoitabh/clob-go/internal/application/account/usecases"
)

func CreditAccountInputFaker() accountUsecases.CreditAccountInput {
	faker := faker.New(0)

	return accountUsecases.CreditAccountInput{
		AccountID: faker.UUID(),
		Asset:     faker.Word(),
		Amount:    int64(faker.Price(1, 1000)),
	}
}
