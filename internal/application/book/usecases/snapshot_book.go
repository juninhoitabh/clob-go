package usecases

import (
	domainBook "github.com/juninhoitabh/clob-go/internal/domain/book"
	"github.com/juninhoitabh/clob-go/internal/shared"
)

type (
	Level struct {
		Price int64
		Qty   int64
	}
	SnapshotBookUseCase struct {
		BookRepo domainBook.IBookRepository
	}
)

func (s *SnapshotBookUseCase) Execute(input SnapshotBookInput) (*SnapshotBookOutput, error) {
	b, err := s.BookRepo.GetBook(input.Instrument)
	if err != nil {
		return nil, err
	}

	if b == nil {
		return nil, shared.ErrNotFound
	}

	out := &SnapshotBookOutput{
		Instrument: input.Instrument,
		Bids:       []Level{},
		Asks:       []Level{},
	}

	for _, p := range b.BidPrices() {
		pl := b.Bids()[p]
		q := pl.TotalQty()

		if q > 0 {
			out.Bids = append(out.Bids, Level{Price: p, Qty: q})
		}
	}

	for _, p := range b.AskPrices() {
		pl := b.Asks()[p]
		q := pl.TotalQty()

		if q > 0 {
			out.Asks = append(out.Asks, Level{Price: p, Qty: q})
		}
	}

	return out, nil
}

func NewSnapshotBookUseCase(
	bookRepo domainBook.IBookRepository,
) *SnapshotBookUseCase {
	return &SnapshotBookUseCase{
		BookRepo: bookRepo,
	}
}
