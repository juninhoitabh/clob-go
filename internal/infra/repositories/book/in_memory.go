package repositories

import (
	"sync"

	"github.com/juninhoitabh/clob-go/internal/domain/book"
)

var (
	instance *InMemoryBookRepository
	once     sync.Once
)

type InMemoryBookRepository struct {
	mu    sync.Mutex
	books map[string]*book.Book
}

func NewInMemoryBookRepository() *InMemoryBookRepository {
	once.Do(func() {
		instance = &InMemoryBookRepository{
			books: make(map[string]*book.Book),
		}
	})

	return instance
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
