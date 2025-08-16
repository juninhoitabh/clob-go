package book

import (
	"github.com/juninhoitabh/clob-go/internal/domain/account"
	domainBook "github.com/juninhoitabh/clob-go/internal/domain/book"
	"github.com/juninhoitabh/clob-go/internal/domain/order"
	"github.com/juninhoitabh/clob-go/internal/shared"
)

type CancelOrderInput struct {
	OrderID string
}

type CancelOrderOutput struct {
	Order *order.Order
}

type CancelOrderUseCase struct {
	BookRepo    domainBook.BookRepository
	AccountRepo account.AccountRepository
}

func (uc *CancelOrderUseCase) Execute(input CancelOrderInput) (*CancelOrderOutput, error) {
	o := uc.BookRepo.GetOrder(input.OrderID)
	if o == nil {
		return nil, shared.ErrNotFound
	}

	b := uc.BookRepo.GetBook(o.Instrument)
	if b == nil {
		return nil, shared.ErrNotFound
	}

	b.RemoveOrder(o)
	uc.BookRepo.SaveBook(b)

	base, quote, err := domainBook.SplitInstrument(o.Instrument)
	if err != nil {
		return nil, err
	}

	if o.Side == order.Buy {
		amount := shared.Mul(o.Price, o.Remaining)
		if err := uc.AccountRepo.ReleaseReserved(o.AccountID, quote, amount); err != nil {
			return nil, err
		}
	} else {
		if err := uc.AccountRepo.ReleaseReserved(o.AccountID, base, o.Remaining); err != nil {
			return nil, err
		}
	}

	o.Remaining = 0
	uc.BookRepo.SaveOrder(o)

	return &CancelOrderOutput{Order: o}, nil
}
