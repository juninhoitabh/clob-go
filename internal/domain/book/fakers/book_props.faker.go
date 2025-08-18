package fakers

import (
	faker "github.com/brianvoe/gofakeit/v7"

	"github.com/juninhoitabh/clob-go/internal/domain/book"
)

func BookPropsFaker() book.BookProps {
	faker := faker.New(0)

	return book.BookProps{
		Instrument: faker.Name(),
	}
}
