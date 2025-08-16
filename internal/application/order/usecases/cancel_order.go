package usecases

import (
	"github.com/juninhoitabh/clob-go/internal/domain/account"
	domainBook "github.com/juninhoitabh/clob-go/internal/domain/book"
	domainOrder "github.com/juninhoitabh/clob-go/internal/domain/order"
	"github.com/juninhoitabh/clob-go/internal/shared"
)

type CancelOrderInput struct {
	OrderID string
}

type CancelOrderOutput struct {
	Order *domainOrder.Order
}

type CancelOrderUseCase struct {
	BookRepo    domainBook.IBookRepository
	AccountRepo account.IAccountRepository
}

func (c *CancelOrderUseCase) Execute(input CancelOrderInput) (*CancelOrderOutput, error) {
	order, err := c.BookRepo.GetOrder(input.OrderID)
	if err != nil {
		return nil, err
	}

	if order == nil {
		return nil, shared.ErrNotFound
	}

	b, err := c.BookRepo.GetBook(order.Instrument)
	if err != nil {
		return nil, err
	}

	if b == nil {
		return nil, shared.ErrNotFound
	}

	b.RemoveOrder(order)
	c.BookRepo.SaveBook(b)

	base, quote, err := domainBook.SplitInstrument(order.Instrument)
	if err != nil {
		return nil, err
	}

	if order.Side == domainOrder.Buy {
		amount := shared.Mul(order.Price, order.Remaining)
		if err := c.AccountRepo.ReleaseReserved(order.AccountID, quote, amount); err != nil {
			return nil, err
		}
	} else {
		if err := c.AccountRepo.ReleaseReserved(order.AccountID, base, order.Remaining); err != nil {
			return nil, err
		}
	}

	order.Remaining = 0
	c.BookRepo.SaveOrder(order)

	return &CancelOrderOutput{Order: order}, nil
}
