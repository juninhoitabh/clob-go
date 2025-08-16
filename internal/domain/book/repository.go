package book

import "github.com/juninhoitabh/clob-go/internal/domain/order"

type BookRepository interface {
	GetBook(instrument string) (*Book, error)
	SaveBook(book *Book) error
	GetOrder(orderID string) (*order.Order, error)
	SaveOrder(o *order.Order) error
	RemoveOrder(orderID string) error
}
