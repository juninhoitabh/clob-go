package account

import (
	"encoding/json"
	"errors"
	"net/http"

	accountUsecases "github.com/juninhoitabh/clob-go/internal/application/account/usecases"
	domainAccount "github.com/juninhoitabh/clob-go/internal/domain/account"
	"github.com/juninhoitabh/clob-go/internal/shared"
)

type (
	createReq struct {
		AccountName string `json:"account_name"`
	}
	creditReq struct {
		Asset  string `json:"asset"`
		Amount int64  `json:"amount"`
	}
	AccountController struct {
		accountDAO  domainAccount.IAccountDAO
		accountRepo domainAccount.IAccountRepository
	}
)

func (a *AccountController) Create(w http.ResponseWriter, req *http.Request) {
	var body createReq
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)

		return
	}

	if body.AccountName == "" {
		http.Error(w, "account_id required", http.StatusBadRequest)

		return
	}

	createAccountUseCase := accountUsecases.NewCreateAccountUseCase(a.accountRepo)

	createAccountOutput, err := createAccountUseCase.Execute(accountUsecases.CreateAccountInput{
		AccountName: body.AccountName,
	})
	if err != nil {
		if errors.Is(err, shared.ErrAlreadyExists) {
			shared.WriteJSON(w, http.StatusOK, map[string]any{"status": "exists"})

			return
		}

		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	shared.WriteJSON(w, http.StatusCreated, map[string]any{"account_id": createAccountOutput.ID})
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

	creditAccountUseCase := accountUsecases.NewCreditAccountUseCase(a.accountRepo)

	err := creditAccountUseCase.Execute(accountUsecases.CreditAccountInput{
		AccountID: id,
		Asset:     body.Asset,
		Amount:    body.Amount,
	})
	if err != nil {
		status := http.StatusBadRequest

		if errors.Is(err, shared.ErrNotFound) {
			status = http.StatusNotFound
		}

		http.Error(w, err.Error(), status)

		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"status": "ok"})
}

func NewAccountController(
	accountDAO domainAccount.IAccountDAO,
	accountRepo domainAccount.IAccountRepository,
) *AccountController {
	return &AccountController{
		accountDAO:  accountDAO,
		accountRepo: accountRepo,
	}
}
