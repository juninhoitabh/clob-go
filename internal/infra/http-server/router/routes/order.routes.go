package routes

import (
	"net/http"

	controllerOrder "github.com/juninhoitabh/clob-go/internal/infra/controllers/order"
	repositoriesAccount "github.com/juninhoitabh/clob-go/internal/infra/repositories/account"
	repositoriesBook "github.com/juninhoitabh/clob-go/internal/infra/repositories/book"
	repositoriesOrder "github.com/juninhoitabh/clob-go/internal/infra/repositories/order"
)

func OrderGenerate(router *http.ServeMux, apiV1Prefix string) {
	accountRepo := repositoriesAccount.NewInMemoryAccountRepository()
	bookRepo := repositoriesBook.NewInMemoryBookRepository()
	orderRepo := repositoriesOrder.NewInMemoryOrderRepository()

	controller := controllerOrder.NewOrderController(
		bookRepo,
		orderRepo,
		accountRepo,
	)

	router.HandleFunc("POST "+apiV1Prefix+"/orders", controller.Place)
	router.HandleFunc("POST "+apiV1Prefix+"/orders/{id}/cancel", controller.Cancel)
}
