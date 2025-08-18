package fakers

import (
	faker "github.com/brianvoe/gofakeit/v7"

	orderUsecases "github.com/juninhoitabh/clob-go/internal/application/order/usecases"
)

func PlaceOrderInputFaker() orderUsecases.PlaceOrderInput {
	faker := faker.New(0)

	return orderUsecases.PlaceOrderInput{
		AccountID:  faker.UUID(),
		Instrument: "BTC/USDT",
		Side:       "buy",
		Price:      int64(faker.Price(100, 10000)),
		Qty:        int64(faker.Number(1, 100)),
	}
}
