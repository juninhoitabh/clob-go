package book

import "github.com/juninhoitabh/clob-go/internal/domain/order"

type BookRepository interface {
	GetBook(instrument string) *Book
	SaveBook(book *Book)
	GetOrder(orderID string) *order.Order
	SaveOrder(o *order.Order)
	RemoveOrder(orderID string)
}
