package book

import (
	"sort"

	"github.com/juninhoitabh/clob-go/internal/domain/order"
)

// TODO: falta id
type Book struct {
	Instrument string
	bids       map[int64]*PriceLevel
	asks       map[int64]*PriceLevel
	bidPrices  []int64
	askPrices  []int64
}

func (b *Book) AddOrder(o *order.Order) {
	if o.Side == order.Buy {
		pl := b.bids[o.Price]
		if pl == nil {
			pl = NewPriceLevel(o.Price)
			b.bids[o.Price] = pl
			b.bidPrices = append(b.bidPrices, o.Price)
			sort.Slice(b.bidPrices, func(i, j int) bool { return b.bidPrices[i] > b.bidPrices[j] })
		}
		pl.Orders = append(pl.Orders, o)
	} else {
		pl := b.asks[o.Price]
		if pl == nil {
			pl = NewPriceLevel(o.Price)
			b.asks[o.Price] = pl
			b.askPrices = append(b.askPrices, o.Price)
			sort.Slice(b.askPrices, func(i, j int) bool { return b.askPrices[i] < b.askPrices[j] })
		}
		pl.Orders = append(pl.Orders, o)
	}
}

func (b *Book) RemoveOrder(o *order.Order) {
	m := b.asks
	prices := &b.askPrices

	if o.Side == order.Buy {
		m = b.bids
		prices = &b.bidPrices
	}

	pl := m[o.Price]
	if pl == nil {
		return
	}

	for i, oo := range pl.Orders {
		if oo == o {
			pl.Orders = append(pl.Orders[:i], pl.Orders[i+1:]...)

			break
		}
	}

	if len(pl.Orders) == 0 {
		delete(m, o.Price)

		for i, p := range *prices {
			if p == o.Price {
				*prices = append((*prices)[:i], (*prices)[i+1:]...)

				break
			}
		}
	}
}

func (b *Book) BestBid() *PriceLevel {
	if len(b.bidPrices) == 0 {
		return nil
	}

	return b.bids[b.bidPrices[0]]
}

func (b *Book) BestAsk() *PriceLevel {
	if len(b.askPrices) == 0 {
		return nil
	}

	return b.asks[b.askPrices[0]]
}

func (b *Book) BidPrices() []int64 {
	return b.bidPrices
}

func (b *Book) AskPrices() []int64 {
	return b.askPrices
}

func (b *Book) Bids() map[int64]*PriceLevel {
	return b.bids
}

func (b *Book) Asks() map[int64]*PriceLevel {
	return b.asks
}

func NewBook(instrument string) *Book {
	return &Book{
		Instrument: instrument,
		bids:       make(map[int64]*PriceLevel),
		asks:       make(map[int64]*PriceLevel),
	}
}
