package fakers

import (
	faker "github.com/brianvoe/gofakeit/v7"

	orderUsecases "github.com/juninhoitabh/clob-go/internal/application/order/usecases"
)

func CancelOrderInputFaker() orderUsecases.CancelOrderInput {
	faker := faker.New(0)

	return orderUsecases.CancelOrderInput{
		OrderID: faker.UUID(),
	}
}
