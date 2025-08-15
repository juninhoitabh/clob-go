package trade

import "time"

type Trade struct {
	ID           string
	TakerOrderID string
	MakerOrderID string
	Price        int64
	Qty          int64
	Buyer        string
	Seller       string
	CreatedAt    time.Time
}
