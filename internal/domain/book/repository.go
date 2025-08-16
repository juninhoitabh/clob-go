package book

type IBookRepository interface {
	GetBook(instrument string) (*Book, error)
	SaveBook(book *Book) error
}
