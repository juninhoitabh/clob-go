package usecases

import (
	domainBook "github.com/juninhoitabh/clob-go/internal/domain/book"
	"github.com/juninhoitabh/clob-go/internal/shared"
)

type (
	SnapshotBookOutput struct {
		Instrument string
		Bids       []Level
		Asks       []Level
	}
	Level struct {
		Price int64
		Qty   int64
	}
	SnapshotBookUseCase struct {
		BookRepo domainBook.IBookRepository
	}
)

func (s *SnapshotBookUseCase) Execute(instrument string) (*SnapshotBookOutput, error) {
	b, err := s.BookRepo.GetBook(instrument)
	if err != nil {
		return nil, err
	}

	if b == nil {
		return nil, shared.ErrNotFound
	}

	out := &SnapshotBookOutput{
		Instrument: instrument,
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
