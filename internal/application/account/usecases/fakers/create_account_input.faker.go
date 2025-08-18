package fakers

import (
	faker "github.com/brianvoe/gofakeit/v7"

	accountUsecases "github.com/juninhoitabh/clob-go/internal/application/account/usecases"
)

func CreateAccountInputFaker() accountUsecases.CreateAccountInput {
	faker := faker.New(0)

	return accountUsecases.CreateAccountInput{
		AccountName: faker.Username(),
	}
}
