package fakers

import (
	faker "github.com/brianvoe/gofakeit/v7"

	bookUsecases "github.com/juninhoitabh/clob-go/internal/application/book/usecases"
)

func SnapshotBookInputFaker() bookUsecases.SnapshotBookInput {
	faker := faker.New(0)

	return bookUsecases.SnapshotBookInput{
		Instrument: faker.Word(),
	}
}
