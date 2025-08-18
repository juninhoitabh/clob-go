package fakers

import (
	faker "github.com/brianvoe/gofakeit/v7"

	"github.com/juninhoitabh/clob-go/internal/domain/order"
)

func OrderPropsFaker() order.OrderProps {
	faker := faker.New(0)

	return order.OrderProps{
		AccountID:  faker.UUID(),
		Instrument: faker.Name(),
		Side:       order.Side(faker.RandomInt([]int{1, 2})),
		Price:      int64(faker.Price(0, 10000)),
		Qty:        int64(faker.Price(0, 1000)),
		Remaining:  int64(faker.Price(0, 1000)),
	}
}
