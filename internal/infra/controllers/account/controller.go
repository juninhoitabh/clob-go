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
	createInputDto struct {
		AccountName string `json:"account_name" example:"test" validate:"required"`
	}
	createOutputDto struct {
		AccountId string `json:"account_id" example:"123e4567-e89b-12d3-a456-426614174000"`
		Status    string `json:"status" example:"created"`
	}
	getAllByIdBalanceOutputDto struct {
		Available int64 `json:"available" example:"1000"`
		Reserved  int64 `json:"reserved" example:"0"`
	}
	getAllByIdOutputDto struct {
		AccountID string                                `json:"account_id" example:"123e4567-e89b-12d3-a456-426614174000"`
		Balances  map[string]getAllByIdBalanceOutputDto `json:"balances"`
	}
	creditInputDto struct {
		Asset  string `json:"asset" example:"USD" validate:"required"`
		Amount int64  `json:"amount" example:"1000" validate:"required,gt=0"`
	}
	creditBalanceOutputDto struct {
		Available int64 `json:"available" example:"1000"`
		Reserved  int64 `json:"reserved" example:"0"`
	}
	creditOutputDto struct {
		AccountID string                                `json:"account_id" example:"123e4567-e89b-12d3-a456-426614174000"`
		Balances  map[string]getAllByIdBalanceOutputDto `json:"balances"`
	}
	AccountController struct {
		accountDAO  domainAccount.IAccountDAO
		accountRepo domainAccount.IAccountRepository
	}
)

// Accounts godoc
// @Summary      Accounts
// @Description  Accounts
// @Tags         Accounts
// @Accept       json
// @Produce      json
// @Param        request   body      createInputDto  true  "createInputDto request"
// @Success      200       {object}  createOutputDto
// @Success      201       {object}  createOutputDto
// @Failure      500       {object}  shared.Errors
// @Router       /accounts [post]
func (a *AccountController) Create(w http.ResponseWriter, req *http.Request) {
	var body createInputDto
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		shared.BadRequestError(w, "Invalid JSON", err.Error())

		return
	}

	if body.AccountName == "" {
		shared.BadRequestError(w, "account_name is required")

		return
	}

	createAccountUseCase := accountUsecases.NewCreateAccountUseCase(a.accountRepo)

	createAccountOutput, err := createAccountUseCase.Execute(accountUsecases.CreateAccountInput{
		AccountName: body.AccountName,
	})
	if err != nil {
		if errors.Is(err, shared.ErrAlreadyExists) {
			shared.WriteJSON(w, http.StatusOK, createOutputDto{Status: "exists"})

			return
		}

		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	shared.WriteJSON(w, http.StatusCreated, createOutputDto{AccountId: createAccountOutput.ID, Status: "created"})
}

// GetAllById godoc
// @Summary      Get All By Id
// @Description  Get All By Id
// @Tags         Accounts
// @Accept       json
// @Produce      json
// @Param        id        path      string               true  "account_id" Format(uuid)
// @Success      200       {object}  getAllByIdOutputDto
// @Failure      500       {object}  shared.Errors
// @Router       /accounts/{id} [get]
func (a *AccountController) GetAllById(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")

	acct, err := a.accountDAO.Snapshot(id)
	if err != nil {
		shared.HandleError(w, err)

		return
	}

	getAllByIdOutputDtoResponse := getAllByIdOutputDto{
		AccountID: acct.AccountID,
		Balances:  make(map[string]getAllByIdBalanceOutputDto),
	}

	for asset, balance := range acct.Balances {
		getAllByIdOutputDtoResponse.Balances[asset] = getAllByIdBalanceOutputDto{
			Available: balance.Available,
			Reserved:  balance.Reserved,
		}
	}

	shared.WriteJSON(w, http.StatusOK, getAllByIdOutputDtoResponse)
}

// Credit godoc
// @Summary      Credit
// @Description  Credit
// @Tags         Accounts
// @Accept       json
// @Produce      json
// @Param        id        path      string          true  "account_id" Format(uuid)
// @Param        request   body      creditInputDto  true  "creditInputDto request"
// @Success      200       {object}  creditOutputDto
// @Failure      500       {object}  shared.Errors
// @Router       /accounts/{id}/credit [post]
func (a *AccountController) Credit(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	if id == "" {
		shared.BadRequestError(w, "missing account_id")

		return
	}

	var body creditInputDto

	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		shared.BadRequestError(w, "Invalid JSON", err.Error())

		return
	}

	if body.Asset == "" || body.Amount <= 0 {
		shared.BadRequestError(w, "asset and positive amount required")

		return
	}

	creditAccountUseCase := accountUsecases.NewCreditAccountUseCase(a.accountRepo)

	err := creditAccountUseCase.Execute(accountUsecases.CreditAccountInput{
		AccountID: id,
		Asset:     body.Asset,
		Amount:    body.Amount,
	})
	if err != nil {
		shared.HandleError(w, err)

		return
	}

	updatedAccount, _ := a.accountDAO.Snapshot(id)

	shared.WriteJSON(w, http.StatusOK, updatedAccount)
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
