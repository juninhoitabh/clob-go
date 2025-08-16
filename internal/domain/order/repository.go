package order

type IOrderRepository interface {
	GetOrder(orderID string) (*Order, error)
	SaveOrder(o *Order) error
	RemoveOrder(orderID string) error
}
