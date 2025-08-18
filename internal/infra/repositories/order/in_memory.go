package repositories

import (
	"sync"

	"github.com/juninhoitabh/clob-go/internal/domain/order"
)

var (
	instance *InMemoryOrderRepository
	once     sync.Once
)

type InMemoryOrderRepository struct {
	mu     sync.Mutex
	orders map[string]*order.Order
}

func NewInMemoryOrderRepository() *InMemoryOrderRepository {
	once.Do(func() {
		instance = &InMemoryOrderRepository{
			orders: make(map[string]*order.Order),
		}
	})

	return instance
}

func (r *InMemoryOrderRepository) GetOrder(orderID string) (*order.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.orders[orderID], nil
}

func (r *InMemoryOrderRepository) SaveOrder(o *order.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.orders[o.GetID()] = o

	return nil
}

func (r *InMemoryOrderRepository) RemoveOrder(orderID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.orders, orderID)

	return nil
}
