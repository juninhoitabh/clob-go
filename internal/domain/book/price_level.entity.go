package book

import "github.com/juninhoitabh/clob-go/internal/domain/order"

type PriceLevel struct {
	Price  int64
	Orders []*order.Order // FIFO
}

func (pl *PriceLevel) TotalQty() int64 {
	var t int64

	for _, o := range pl.Orders {
		t += o.Remaining
	}

	return t
}

func NewPriceLevel(price int64) *PriceLevel {
	return &PriceLevel{
		Price:  price,
		Orders: []*order.Order{},
	}
}
