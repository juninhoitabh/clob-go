package usecases

type (
	ICancelOrderUseCase interface {
		Execute(input CancelOrderInput) (*CancelOrderOutput, error)
	}
	IPlaceOrderUseCase interface {
		Execute(input PlaceOrderInput) (*PlaceOrderOutput, error)
	}
)
