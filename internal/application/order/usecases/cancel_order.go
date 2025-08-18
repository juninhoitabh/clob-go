package usecases

import (
	"github.com/juninhoitabh/clob-go/internal/domain/account"
	domainBook "github.com/juninhoitabh/clob-go/internal/domain/book"
	domainOrder "github.com/juninhoitabh/clob-go/internal/domain/order"
	"github.com/juninhoitabh/clob-go/internal/shared"
)

type CancelOrderUseCase struct {
	BookRepo    domainBook.IBookRepository
	OrderRepo   domainOrder.IOrderRepository
	AccountRepo account.IAccountRepository
}

func (c *CancelOrderUseCase) Execute(input CancelOrderInput) (*CancelOrderOutput, error) {
	order, err := c.OrderRepo.GetOrder(input.OrderID)
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

	err = c.BookRepo.SaveBook(b)
	if err != nil {
		return nil, err
	}

	base, quote, err := domainBook.SplitInstrument(order.Instrument)
	if err != nil {
		return nil, err
	}

	acct, err := c.AccountRepo.Get(order.AccountID)
	if err != nil {
		return nil, shared.ErrNotFound
	}

	if order.Side == domainOrder.Buy {
		amount := shared.Mul(order.Price, order.Remaining)

		if err := acct.ReleaseReserved(quote, amount); err != nil {
			return nil, err
		}
	} else {
		if err := acct.ReleaseReserved(base, order.Remaining); err != nil {
			return nil, err
		}
	}

	if err := c.AccountRepo.Save(acct); err != nil {
		return nil, err
	}

	order.Remaining = 0
	err = c.OrderRepo.SaveOrder(order)
	if err != nil {
		return nil, err
	}

	return &CancelOrderOutput{Order: order}, nil
}

func NewCancelOrderUseCase(
	bookRepo domainBook.IBookRepository,
	orderRepo domainOrder.IOrderRepository,
	accountRepo account.IAccountRepository,
) *CancelOrderUseCase {
	return &CancelOrderUseCase{
		BookRepo:    bookRepo,
		OrderRepo:   orderRepo,
		AccountRepo: accountRepo,
	}
}
