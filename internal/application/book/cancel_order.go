package book

import (
	"time"

	"github.com/google/uuid"
	"github.com/juninhoitabh/clob-go/internal/domain/account"
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
	BookRepo    domainBook.BookRepository
	AccountRepo account.AccountsRepository
}

func (uc *PlaceOrderUseCase) Execute(input PlaceOrderInput) (*PlaceOrderOutput, error) {
	if input.Price <= 0 || input.Qty <= 0 {
		return nil, shared.ErrInvalidParam
	}

	side, err := order.ParseSide(input.Side)
	if err != nil {
		return nil, err
	}

	_, err = uc.AccountRepo.Snapshot(input.AccountID)
	if err != nil {
		return nil, shared.ErrNotFound
	}

	base, quote, err := domainBook.SplitInstrument(input.Instrument)
	if err != nil {
		return nil, err
	}

	if side == order.Buy {
		cost := domainBook.Mul(input.Price, input.Qty)
		if err := uc.AccountRepo.Reserve(input.AccountID, quote, cost); err != nil {
			return nil, err
		}
	} else {
		if err := uc.AccountRepo.Reserve(input.AccountID, base, input.Qty); err != nil {
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

	uc.BookRepo.SaveOrder(o)

	b := uc.BookRepo.GetBook(input.Instrument)
	if b == nil {
		b = domainBook.NewBook(input.Instrument)
		uc.BookRepo.SaveBook(b)
	}

	report := services.MatchOrder(b, o)
	uc.BookRepo.SaveBook(b)

	// Aqui você pode adicionar a lógica de settleTrade como serviço de domínio também

	return &PlaceOrderOutput{
		Order:       o,
		TradeReport: report,
	}, nil
}
