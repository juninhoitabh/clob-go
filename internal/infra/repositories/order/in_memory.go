package repositories

import (
	"sync"

	"github.com/juninhoitabh/clob-go/internal/domain/order"
)

type InMemoryBookRepository struct {
	mu     sync.Mutex
	orders map[string]*order.Order
}

func NewInMemoryBookRepository() *InMemoryBookRepository {
	return &InMemoryBookRepository{
		orders: make(map[string]*order.Order),
	}
}

func (r *InMemoryBookRepository) GetOrder(orderID string) (*order.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.orders[orderID], nil
}

func (r *InMemoryBookRepository) SaveOrder(o *order.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.orders[o.GetID()] = o

	return nil
}

func (r *InMemoryBookRepository) RemoveOrder(orderID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.orders, orderID)

	return nil
}
