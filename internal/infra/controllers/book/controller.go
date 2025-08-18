package book

import (
	"fmt"
	"net/http"
	"strings"

	bookUsecases "github.com/juninhoitabh/clob-go/internal/application/book/usecases"
	domainBook "github.com/juninhoitabh/clob-go/internal/domain/book"
	"github.com/juninhoitabh/clob-go/internal/shared"
)

type BookController struct {
	bookRepo domainBook.IBookRepository
}

func (b *BookController) Get(w http.ResponseWriter, req *http.Request) {
	inst := req.PathValue("instrument")

	inst = strings.ToUpper(inst)
	if inst == "" {
		http.Error(w, "instrument required", http.StatusBadRequest)

		return
	}

	snapshotBookUseCase := bookUsecases.NewSnapshotBookUseCase(b.bookRepo)

	book, err := snapshotBookUseCase.Execute(bookUsecases.SnapshotBookInput{
		Instrument: inst,
	})
	if err != nil {
		status := http.StatusBadRequest

		if strings.Contains(err.Error(), "not found") {
			status = http.StatusNotFound
		}

		http.Error(w, err.Error(), status)

		return
	}

	if book == nil {
		http.Error(w, fmt.Sprintf("instrument %s not found (empty book yet?)", inst), http.StatusNotFound)

		return
	}

	shared.WriteJSON(w, http.StatusOK, book)
}

func NewBookController(
	bookRepo domainBook.IBookRepository,
) *BookController {
	return &BookController{
		bookRepo: bookRepo,
	}
}
