package book

import (
	"net/http"
	"strings"

	bookUsecases "github.com/juninhoitabh/clob-go/internal/application/book/usecases"
	domainBook "github.com/juninhoitabh/clob-go/internal/domain/book"
	"github.com/juninhoitabh/clob-go/internal/shared"
)

type (
	getByInstrumentLevelOutputDto struct {
		Price int64 `json:"price" example:"50000"`
		Qty   int64 `json:"qty" example:"1"`
	}
	getByInstrumentOutputDto struct {
		Instrument string                          `json:"instrument" example:"BTC/USDT"`
		Bids       []getByInstrumentLevelOutputDto `json:"bids"`
		Asks       []getByInstrumentLevelOutputDto `json:"asks"`
	}
	BookController struct {
		bookRepo domainBook.IBookRepository
	}
)

// GetByInstrument godoc
// @Summary      Get by Instrument
// @Description  Get by Instrument
// @Tags         Books
// @Accept       json
// @Produce      json
// @Param        instrument path      string true "instrument" example:"BTC/USDT"
// @Success      200       {object}  getByInstrumentOutputDto
// @Failure      500       {object}  shared.Errors
// @Router       /books/{instrument} [get]
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
		shared.HandleError(w, err)

		return
	}

	if book == nil {
		shared.HandleError(w, shared.ErrNotFound)

		return
	}

	getByInstrumentOutputDtoResponse := getByInstrumentOutputDto{
		Instrument: book.Instrument,
		Bids:       []getByInstrumentLevelOutputDto{},
		Asks:       []getByInstrumentLevelOutputDto{},
	}

	for _, bid := range book.Bids {
		getByInstrumentOutputDtoResponse.Bids = append(getByInstrumentOutputDtoResponse.Bids, getByInstrumentLevelOutputDto{
			Price: bid.Price,
			Qty:   bid.Qty,
		})
	}

	for _, ask := range book.Asks {
		getByInstrumentOutputDtoResponse.Asks = append(getByInstrumentOutputDtoResponse.Asks, getByInstrumentLevelOutputDto{
			Price: ask.Price,
			Qty:   ask.Qty,
		})
	}

	shared.WriteJSON(w, http.StatusOK, getByInstrumentOutputDtoResponse)
}

func NewBookController(
	bookRepo domainBook.IBookRepository,
) *BookController {
	return &BookController{
		bookRepo: bookRepo,
	}
}
