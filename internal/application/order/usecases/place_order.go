package usecases

import (
	"time"

	"github.com/google/uuid"
	"github.com/juninhoitabh/clob-go/internal/domain/account"
	accountServices "github.com/juninhoitabh/clob-go/internal/domain/account/services"
	domainBook "github.com/juninhoitabh/clob-go/internal/domain/book"
	"github.com/juninhoitabh/clob-go/internal/domain/book/services"
	"github.com/juninhoitabh/clob-go/internal/domain/order"
	"github.com/juninhoitabh/clob-go/internal/shared"
)

type PlaceOrderInput struct {
	AccountID  string
	Instrument string
	Side       string
	Price      int64
	Qty        int64
}

type PlaceOrderOutput struct {
	Order       *order.Order
	TradeReport *services.TradeReport
}

type PlaceOrderUseCase struct {
	BookRepo    domainBook.IBookRepository
	AccountRepo account.IAccountRepository
	AccountDAO  account.IAccountDAO
}

func (p *PlaceOrderUseCase) Execute(input PlaceOrderInput) (*PlaceOrderOutput, error) {
	if input.Price <= 0 || input.Qty <= 0 {
		return nil, shared.ErrInvalidParam
	}

	side, err := order.ParseSide(input.Side)
	if err != nil {
		return nil, err
	}

	_, err = p.AccountDAO.Snapshot(input.AccountID)
	if err != nil {
		return nil, shared.ErrNotFound
	}

	base, quote, err := domainBook.SplitInstrument(input.Instrument)
	if err != nil {
		return nil, err
	}

	if side == order.Buy {
		cost := shared.Mul(input.Price, input.Qty)
		if err := p.AccountRepo.Reserve(input.AccountID, quote, cost); err != nil {
			return nil, err
		}
	} else {
		if err := p.AccountRepo.Reserve(input.AccountID, base, input.Qty); err != nil {
			return nil, err
		}
	}

	o := &order.Order{
		ID:         uuid.NewString(),
		AccountID:  input.AccountID,
		Instrument: input.Instrument,
		Side:       side,
		Price:      input.Price,
		Qty:        input.Qty,
		Remaining:  input.Qty,
		CreatedAt:  time.Now(),
	}

	p.BookRepo.SaveOrder(o)

	b, err := p.BookRepo.GetBook(input.Instrument)
	if err != nil {
		return nil, err
	}

	if b == nil {
		b = domainBook.NewBook(input.Instrument)
		err = p.BookRepo.SaveBook(b)
		if err != nil {
			return nil, err
		}
	}

	report := services.MatchOrder(b, o)
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
		Order:       o,
		TradeReport: report,
	}, nil
}
