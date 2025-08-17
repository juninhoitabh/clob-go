package usecases

import domainOrder "github.com/juninhoitabh/clob-go/internal/domain/order"

type (
	CancelOrderInput struct {
		OrderID string
	}
	CancelOrderOutput struct {
		Order *domainOrder.Order
	}
	ICancelOrderUseCase interface {
		Execute(input CancelOrderInput) (*CancelOrderOutput, error)
	}
	IPlaceOrderUseCase interface {
		Execute(input PlaceOrderInput) (*PlaceOrderOutput, error)
	}
)
