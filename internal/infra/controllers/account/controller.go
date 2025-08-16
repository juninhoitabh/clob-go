package account

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/juninhoitabh/clob-go/internal/domain/account"
	"github.com/juninhoitabh/clob-go/internal/shared"
)

type AccountController struct {
	accountRepo account.AccountRepository
	accountDAO  account.AccountDAO
}

type (
	createReq struct {
		AccountID string `json:"account_id"`
	}
	creditReq struct {
		Asset  string `json:"asset"`
		Amount int64  `json:"amount"`
	}
)

// TODO: fazer um usecase
func (a *AccountController) Create(w http.ResponseWriter, req *http.Request) {
	var body createReq
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if body.AccountID == "" {
		http.Error(w, "account_id required", http.StatusBadRequest)
		return
	}

	created := a.accountRepo.Create(body.AccountID, "") // TODO: verrr
	if created {
		shared.WriteJSON(w, http.StatusCreated, map[string]any{"status": "created"})
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"status": "exists"})
}

func (a *AccountController) Get(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")

	acct, err := a.accountDAO.Snapshot(id)
	if err != nil {
		status := http.StatusBadRequest

		if errors.Is(err, shared.ErrNotFound) {
			status = http.StatusNotFound
		}

		http.Error(w, err.Error(), status)

		return
	}

	shared.WriteJSON(w, http.StatusOK, acct)
}

// TODO: fazer u usecase
func (a *AccountController) Credit(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	if id == "" {
		http.Error(w, "missing account id", http.StatusBadRequest)
		return
	}

	var body creditReq

	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if body.Asset == "" || body.Amount <= 0 {
		http.Error(w, "asset and positive amount required", http.StatusBadRequest)

		return
	}

	if err := a.accountRepo.Credit(id, strings.ToUpper(body.Asset), body.Amount); err != nil {
		status := http.StatusBadRequest

		if errors.Is(err, shared.ErrNotFound) {
			status = http.StatusNotFound
		}

		http.Error(w, err.Error(), status)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"status": "ok"})
}
