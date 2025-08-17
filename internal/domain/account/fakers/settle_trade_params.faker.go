package fakers

import (
	faker "github.com/brianvoe/gofakeit/v7"
)

type SettleTradeParams struct {
	BuyerID  string
	SellerID string
	Base     string
	Quote    string
	Price    int64
	Qty      int64
}

func SettleTradeParamsFaker() SettleTradeParams {
	f := faker.New(0)
	return SettleTradeParams{
		BuyerID:  f.UUID(),
		SellerID: f.UUID(),
		Base:     f.CurrencyShort(),
		Quote:    f.CurrencyShort(),
		Price:    int64(f.Number(1, 100000)),
		Qty:      int64(f.Number(1, 100000)),
	}
}
