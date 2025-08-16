package repositories

import (
	"sync"

	"github.com/juninhoitabh/clob-go/internal/domain/book"
	"github.com/juninhoitabh/clob-go/internal/domain/order"
)

type InMemoryBookRepository struct {
	mu     sync.Mutex
	books  map[string]*book.Book
	orders map[string]*order.Order
}

func NewInMemoryBookRepository() *InMemoryBookRepository {
	return &InMemoryBookRepository{
		books:  make(map[string]*book.Book),
		orders: make(map[string]*order.Order),
	}
}

func (r *InMemoryBookRepository) GetBook(instrument string) (*book.Book, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.books[instrument], nil
}

func (r *InMemoryBookRepository) SaveBook(b *book.Book) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.books[b.Instrument] = b

	return nil
}

func (r *InMemoryBookRepository) GetOrder(orderID string) (*order.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.orders[orderID], nil
}

func (r *InMemoryBookRepository) SaveOrder(o *order.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.orders[o.ID] = o

	return nil
}

func (r *InMemoryBookRepository) RemoveOrder(orderID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.orders, orderID)

	return nil
}
