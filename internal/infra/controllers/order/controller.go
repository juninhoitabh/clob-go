package order

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	orderUsecases "github.com/juninhoitabh/clob-go/internal/application/order/usecases"
	"github.com/juninhoitabh/clob-go/internal/domain/account"
	domainBook "github.com/juninhoitabh/clob-go/internal/domain/book"
	domainOrder "github.com/juninhoitabh/clob-go/internal/domain/order"
	"github.com/juninhoitabh/clob-go/internal/shared"
)

type OrderController struct {
	bookRepo    domainBook.IBookRepository
	orderRepo   domainOrder.IOrderRepository
	accountRepo account.IAccountRepository
}

type placeReq struct {
	AccountID  string `json:"account_id"`
	Instrument string `json:"instrument"`
	Side       string `json:"side"` // "buy" or "sell"
	Price      int64  `json:"price"`
	Qty        int64  `json:"qty"`
}

func (o *OrderController) Place(w http.ResponseWriter, req *http.Request) {
	var body placeReq

	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if body.AccountID == "" || body.Instrument == "" || (body.Side != "buy" && body.Side != "sell") || body.Price <= 0 || body.Qty <= 0 {
		http.Error(w, "invalid fields", http.StatusBadRequest)

		return
	}

	placeOrderInput := orderUsecases.PlaceOrderInput{
		AccountID:  body.AccountID,
		Instrument: strings.ToUpper(body.Instrument),
		Side:       strings.ToLower(body.Side),
		Price:      body.Price,
		Qty:        body.Qty,
	}

	placeOrderUseCase := orderUsecases.NewPlaceOrderUseCase(o.bookRepo, o.orderRepo, o.accountRepo)

	placeOrderOutput, err := placeOrderUseCase.Execute(placeOrderInput)
	if err != nil {
		status := http.StatusBadRequest

		if errors.Is(err, shared.ErrNotFound) {
			status = http.StatusNotFound
		}

		http.Error(w, err.Error(), status)

		return
	}

	shared.WriteJSON(w, http.StatusCreated, map[string]any{
		"order":  placeOrderOutput.Order.Public(),
		"report": placeOrderOutput.TradeReport,
	})
}

func (o *OrderController) Cancel(w http.ResponseWriter, req *http.Request) {
	oid := req.PathValue("id")
	if oid == "" {
		http.Error(w, "order id required", http.StatusBadRequest)
		return
	}

	cancelOrderUseCase := orderUsecases.NewCancelOrderUseCase(o.bookRepo, o.orderRepo, o.accountRepo)

	cancelOrderInput := orderUsecases.CancelOrderInput{
		OrderID: oid,
	}

	cancelOrderOutput, err := cancelOrderUseCase.Execute(cancelOrderInput)
	if err != nil {
		status := http.StatusBadRequest

		if errors.Is(err, shared.ErrNotFound) {
			status = http.StatusNotFound
		}

		http.Error(w, err.Error(), status)

		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"order": cancelOrderOutput.Order.Public(), "status": "canceled"})
}

func NewOrderController(
	bookRepo domainBook.IBookRepository,
	orderRepo domainOrder.IOrderRepository,
	accountRepo account.IAccountRepository,
) *OrderController {
	return &OrderController{
		bookRepo:    bookRepo,
		orderRepo:   orderRepo,
		accountRepo: accountRepo,
	}
}
