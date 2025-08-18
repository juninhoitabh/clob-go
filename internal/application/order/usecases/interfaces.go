package usecases

import (
	"github.com/juninhoitabh/clob-go/internal/domain/book/services"
	domainOrder "github.com/juninhoitabh/clob-go/internal/domain/order"
)

type (
	CancelOrderInput struct {
		OrderID string
	}
	CancelOrderOutput struct {
		Order *domainOrder.Order
	}
	PlaceOrderInput struct {
		AccountID  string
		Instrument string
		Side       string
		Price      int64
		Qty        int64
	}
	PlaceOrderOutput struct {
		Order       *domainOrder.Order
		TradeReport *services.TradeReport
	}
	ICancelOrderUseCase interface {
		Execute(input CancelOrderInput) (*CancelOrderOutput, error)
	}
	IPlaceOrderUseCase interface {
		Execute(input PlaceOrderInput) (*PlaceOrderOutput, error)
	}
)
