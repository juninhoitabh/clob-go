package usecases

import (
	"github.com/juninhoitabh/clob-go/internal/domain/account"
	accountServices "github.com/juninhoitabh/clob-go/internal/domain/account/services"
	domainBook "github.com/juninhoitabh/clob-go/internal/domain/book"
	"github.com/juninhoitabh/clob-go/internal/domain/book/services"
	domainOrder "github.com/juninhoitabh/clob-go/internal/domain/order"
	"github.com/juninhoitabh/clob-go/internal/shared"
	idObjValue "github.com/juninhoitabh/clob-go/internal/shared/domain/value-objects/id"
)

type PlaceOrderInput struct {
	AccountID  string
	Instrument string
	Side       string
	Price      int64
	Qty        int64
}

type PlaceOrderOutput struct {
	Order       *domainOrder.Order
	TradeReport *services.TradeReport
}

type PlaceOrderUseCase struct {
	BookRepo    domainBook.IBookRepository
	OrderRepo   domainOrder.IOrderRepository
	AccountRepo account.IAccountRepository
}

func (p *PlaceOrderUseCase) Execute(input PlaceOrderInput) (*PlaceOrderOutput, error) {
	if input.Price <= 0 || input.Qty <= 0 {
		return nil, shared.ErrInvalidParam
	}

	side, err := domainOrder.ParseSide(input.Side)
	if err != nil {
		return nil, err
	}

	acct, err := p.AccountRepo.Get(input.AccountID)
	if err != nil {
		return nil, shared.ErrNotFound
	}

	base, quote, err := domainBook.SplitInstrument(input.Instrument)
	if err != nil {
		return nil, err
	}

	if side == domainOrder.Buy {
		cost := shared.Mul(input.Price, input.Qty)
		if err := acct.Reserve(quote, cost); err != nil {
			return nil, err
		}
	} else {
		if err := acct.Reserve(base, input.Qty); err != nil {
			return nil, err
		}
	}

	if err := p.AccountRepo.Save(acct); err != nil {
		return nil, err
	}

	order, err := domainOrder.NewOrder(domainOrder.OrderProps{
		AccountID:  input.AccountID,
		Instrument: input.Instrument,
		Side:       side,
		Price:      input.Price,
		Qty:        input.Qty,
		Remaining:  input.Qty,
	}, idObjValue.Uuid)
	if err != nil {
		return nil, err
	}

	p.OrderRepo.SaveOrder(order)

	b, err := p.BookRepo.GetBook(input.Instrument)
	if err != nil {
		return nil, err
	}

	if b == nil {
		b, err = domainBook.NewBook(domainBook.BookProps{
			Instrument: input.Instrument,
		}, idObjValue.Uuid)
		if err != nil {
			return nil, err
		}

		err = p.BookRepo.SaveBook(b)
		if err != nil {
			return nil, err
		}
	}

	report := services.MatchOrder(b, order)
	err = p.BookRepo.SaveBook(b)
	if err != nil {
		return nil, err
	}

	for _, trade := range report.Trades {
		err := accountServices.SettleTrade(
			p.AccountRepo,
			trade.Buyer,
			trade.Seller,
			base,
			quote,
			trade.Price,
			trade.Qty,
		)
		if err != nil {
			return nil, err
		}
	}

	return &PlaceOrderOutput{
		Order:       order,
		TradeReport: report,
	}, nil
}
