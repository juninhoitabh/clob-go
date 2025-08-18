package routes

import (
	"net/http"

	controllerAccount "github.com/juninhoitabh/clob-go/internal/infra/controllers/account"
	daosAccount "github.com/juninhoitabh/clob-go/internal/infra/daos/account"
	repositoriesAccount "github.com/juninhoitabh/clob-go/internal/infra/repositories/account"
)

func AccountGenerate(router *http.ServeMux) {
	accountRepo := repositoriesAccount.NewInMemoryAccountRepository()
	accountDAO := daosAccount.NewInMemoryAccountDAO(accountRepo.Mutex(), accountRepo.AccountsMap())

	controller := controllerAccount.NewAccountController(
		accountDAO,
		accountRepo,
	)

	router.HandleFunc("POST /accounts", controller.Create)
	router.HandleFunc("POST /accounts/{id}/credit", controller.Credit)
	router.HandleFunc("GET /accounts/{id}", controller.Get)
}
