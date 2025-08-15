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

func (r *InMemoryBookRepository) GetBook(instrument string) *book.Book {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.books[instrument]
}

func (r *InMemoryBookRepository) SaveBook(b *book.Book) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.books[b.Instrument] = b
}

func (r *InMemoryBookRepository) GetOrder(orderID string) *order.Order {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.orders[orderID]
}

func (r *InMemoryBookRepository) SaveOrder(o *order.Order) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.orders[o.ID] = o
}

func (r *InMemoryBookRepository) RemoveOrder(orderID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.orders, orderID)
}
