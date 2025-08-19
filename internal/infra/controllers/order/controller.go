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

type (
	placeInputDto struct {
		AccountID  string `json:"account_id" example:"123e4567-e89b-12d3-a456-426614174000" validate:"required"`
		Instrument string `json:"instrument" example:"BTC-USD" validate:"required"`
		Side       string `json:"side" example:"buy" validate:"required,oneof=buy sell"`
		Price      int64  `json:"price" example:"50000" validate:"required,gte=1"`
		Qty        int64  `json:"qty" example:"1" validate:"required,gte=1"`
	}
	placeTradeOutputDto struct {
		TakerOrderID string `json:"taker_order_id" example:"123e4567-e89b-12d3-a456-426614174000"`
		MakerOrderID string `json:"maker_order_id" example:"123e4567-e89b-12d3-a456-426614174000"`
		Price        int64  `json:"price" example:"50000" validate:"required,gte=1"`
		Qty          int64  `json:"qty" example:"1" validate:"required,gte=1"`
		BuyerID      string `json:"buyer_id" example:"123e4567-e89b-12d3-a456-426614174000"`
		SellerID     string `json:"seller_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	}
	placeTradeReportOutputDto struct {
		Trades []placeTradeOutputDto `json:"trades"`
	}

	placeOutputDto struct {
		Order  map[string]any            `json:"order"`
		Report placeTradeReportOutputDto `json:"report"`
	}
	OrderController struct {
		bookRepo    domainBook.IBookRepository
		orderRepo   domainOrder.IOrderRepository
		accountRepo account.IAccountRepository
	}
)

// Orders godoc
// @Summary      Orders
// @Description  Orders
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        request   body      placeInputDto  true  "placeInputDto request"
// @Success      201       {object}  placeOutputDto
// @Failure      500       {object}  shared.Errors
// @Router       /orders   [post]
func (o *OrderController) Place(w http.ResponseWriter, req *http.Request) {
	var body placeInputDto
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		shared.BadRequestError(w, "Invalid JSON", err.Error())

		return
	}

	if body.AccountID == "" || body.Instrument == "" || (body.Side != "buy" && body.Side != "sell") || body.Price <= 0 || body.Qty <= 0 {
		shared.BadRequestError(w, "invalid fields")

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
		shared.HandleError(w, err)

		return
	}

	placeOutputDtoResult := placeOutputDto{
		Order:  placeOrderOutput.Order.Public(),
		Report: placeTradeReportOutputDto{},
	}

	for _, trade := range placeOrderOutput.TradeReport.Trades {
		placeOutputDtoResult.Report.Trades = append(placeOutputDtoResult.Report.Trades, placeTradeOutputDto{
			TakerOrderID: trade.TakerOrderID,
			MakerOrderID: trade.MakerOrderID,
			Price:        trade.Price,
			Qty:          trade.Qty,
			BuyerID:      trade.BuyerID,
			SellerID:     trade.SellerID,
		})
	}

	shared.WriteJSON(w, http.StatusCreated, placeOutputDtoResult)
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
