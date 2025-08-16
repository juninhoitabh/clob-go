package order

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/juninhoitabh/clob-go/internal/shared"
)

type OrderController struct {
}

type placeReq struct {
	AccountID  string `json:"account_id"`
	Instrument string `json:"instrument"`
	Side       string `json:"side"` // "buy" or "sell"
	Price      int64  `json:"price"`
	Qty        int64  `json:"qty"` // TODO: tem que ser float64??
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

	order, report, err := o.Eng.Place(body.AccountID, strings.ToUpper(body.Instrument), strings.ToLower(body.Side), body.Price, body.Qty)
	if err != nil {
		status := http.StatusBadRequest

		if errors.Is(err, engine.ErrNotFound) {
			status = http.StatusNotFound
		}

		http.Error(w, err.Error(), status)

		return
	}

	shared.WriteJSON(w, http.StatusCreated, map[string]any{
		"order":  order.Public(),
		"report": report,
	})
}

func (o *OrderController) Cancel(w http.ResponseWriter, req *http.Request) {
	oid := req.PathValue("id")
	if oid == "" {
		http.Error(w, "order id required", http.StatusBadRequest)
		return
	}

	order, err := o.Eng.Cancel(oid)
	if err != nil {
		status := http.StatusBadRequest

		if errors.Is(err, engine.ErrNotFound) {
			status = http.StatusNotFound
		}

		http.Error(w, err.Error(), status)

		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"order": order.Public(), "status": "canceled"})
}
