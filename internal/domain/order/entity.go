package order

import (
	"errors"
	"time"

	baseEntity "github.com/juninhoitabh/clob-go/internal/shared/domain/entities"
	idObjValue "github.com/juninhoitabh/clob-go/internal/shared/domain/value-objects/id"
)

var (
	ErrInvalidOrder     = errors.New("invalid order")
	ErrInvalidSideOrder = errors.New("invalid side order")
)

type Side int

const (
	Buy Side = iota + 1
	Sell
)

type OrderProps struct {
	AccountID  string
	Instrument string
	Side       Side
	Price      int64
	Qty        int64
	Remaining  int64
}

type Order struct {
	CreatedAt time.Time
	baseEntity.BaseEntity
	AccountID  string
	Instrument string
	Side       Side
	Price      int64
	Qty        int64
	Remaining  int64
}

func (o *Order) Prepare(typeId idObjValue.TypeIdEnum) error {
	err := o.Validate()
	if err != nil {
		return err
	}

	o.BaseEntity.NewBaseEntity("", typeId)

	o.CreatedAt = time.Now()

	return nil
}

func (o *Order) Validate() error {
	if o.AccountID == "" || o.Instrument == "" {
		return ErrInvalidOrder
	}

	if o.Side != Buy && o.Side != Sell {
		return ErrInvalidSideOrder
	}

	if o.Price <= 0 || o.Qty <= 0 || o.Remaining < 0 || o.Remaining > o.Qty {
		return ErrInvalidOrder
	}

	return nil
}

func (o *Order) Public() map[string]any {
	side := "buy"
	if o.Side == Sell {
		side = "sell"
	}

	return map[string]any{
		"id":         o.BaseEntity.ID.ID,
		"account_id": o.AccountID,
		"instrument": o.Instrument,
		"side":       side,
		"price":      o.Price,
		"qty":        o.Qty,
		"remaining":  o.Remaining,
		"created_at": o.CreatedAt.UTC().Format(time.RFC3339Nano),
	}
}

func NewOrder(props OrderProps, typeId idObjValue.TypeIdEnum) (*Order, error) {
	order := Order{
		AccountID:  props.AccountID,
		Instrument: props.Instrument,
		Side:       props.Side,
		Price:      props.Price,
		Qty:        props.Qty,
		Remaining:  props.Remaining,
	}

	err := order.Prepare(typeId)
	if err != nil {
		return nil, err
	}

	return &order, nil
}
