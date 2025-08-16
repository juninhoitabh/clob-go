package account

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	domainAccount "github.com/juninhoitabh/clob-go/internal/domain/account"
	"github.com/juninhoitabh/clob-go/internal/shared"
	idObjValue "github.com/juninhoitabh/clob-go/internal/shared/domain/value-objects/id"
)

type AccountController struct {
	accountRepo domainAccount.IAccountRepository
	accountDAO  domainAccount.IAccountDAO
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

	account, err := domainAccount.NewAccount(domainAccount.AccountProps{
		Name: body.AccountID,
	}, idObjValue.Uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	created := a.accountRepo.Create(account)
	if created {
		shared.WriteJSON(w, http.StatusCreated, map[string]any{"status": "created"})
		return
	}

	// TODO: retornar id
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

	acct, err := a.accountRepo.Get(id)
	if err != nil {
		status := http.StatusBadRequest

		if errors.Is(err, shared.ErrNotFound) {
			status = http.StatusNotFound
		}

		http.Error(w, err.Error(), status)

		return
	}

	acct.Credit(strings.ToUpper(body.Asset), body.Amount)

	if err := a.accountRepo.Save(acct); err != nil {
		status := http.StatusBadRequest

		if errors.Is(err, shared.ErrNotFound) {
			status = http.StatusNotFound
		}

		http.Error(w, err.Error(), status)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"status": "ok"})
}
