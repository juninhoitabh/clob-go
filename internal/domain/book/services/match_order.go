package services

import (
	"github.com/juninhoitabh/clob-go/internal/domain/book"
	"github.com/juninhoitabh/clob-go/internal/domain/order"
)

type Trade struct {
	TakerOrderID string
	MakerOrderID string
	Price        int64
	Qty          int64
	Buyer        string
	Seller       string
}

type TradeReport struct {
	Trades []Trade
}

func MatchOrder(b *book.Book, o *order.Order) *TradeReport {
	report := &TradeReport{}

	if o.Side == order.Buy {
		for o.Remaining > 0 {
			ask := b.BestAsk()
			if ask == nil || ask.Price > o.Price {
				break
			}

			for len(ask.Orders) > 0 && o.Remaining > 0 && ask.Price <= o.Price {
				maker := ask.Orders[0]
				tradeQty := min(o.Remaining, maker.Remaining)
				execPrice := maker.Price

				report.Trades = append(report.Trades, Trade{
					TakerOrderID: o.ID,
					MakerOrderID: maker.ID,
					Price:        execPrice,
					Qty:          tradeQty,
					Buyer:        o.AccountID,
					Seller:       maker.AccountID,
				})

				o.Remaining -= tradeQty
				maker.Remaining -= tradeQty

				if maker.Remaining == 0 {
					b.RemoveOrder(maker)
				} else {
					ask.Orders[0] = maker
				}
			}
		}

		if o.Remaining > 0 {
			b.AddOrder(o)
		}
	} else {
		for o.Remaining > 0 {
			bid := b.BestBid()
			if bid == nil || bid.Price < o.Price {
				break
			}

			for len(bid.Orders) > 0 && o.Remaining > 0 && bid.Price >= o.Price {
				maker := bid.Orders[0]
				tradeQty := min(o.Remaining, maker.Remaining)
				execPrice := maker.Price

				report.Trades = append(report.Trades, Trade{
					TakerOrderID: o.ID,
					MakerOrderID: maker.ID,
					Price:        execPrice,
					Qty:          tradeQty,
					Buyer:        maker.AccountID,
					Seller:       o.AccountID,
				})

				o.Remaining -= tradeQty
				maker.Remaining -= tradeQty

				if maker.Remaining == 0 {
					b.RemoveOrder(maker)
				} else {
					bid.Orders[0] = maker
				}
			}
		}

		if o.Remaining > 0 {
			b.AddOrder(o)
		}
	}

	return report
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
