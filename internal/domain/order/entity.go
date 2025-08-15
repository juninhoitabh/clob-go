package order

import "time"

type Side int

const (
	Buy Side = iota + 1
	Sell
)

type Order struct {
	ID         string
	AccountID  string
	Instrument string
	Side       Side
	Price      int64
	Qty        int64
	Remaining  int64
	CreatedAt  time.Time
}

func (o *Order) Public() map[string]any {
	side := "buy"
	if o.Side == Sell {
		side = "sell"
	}

	return map[string]any{
		"id":         o.ID,
		"account_id": o.AccountID,
		"instrument": o.Instrument,
		"side":       side,
		"price":      o.Price,
		"qty":        o.Qty,
		"remaining":  o.Remaining,
		"created_at": o.CreatedAt.UTC().Format(time.RFC3339Nano),
	}
}
