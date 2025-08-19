package routes

import (
	"net/http"

	controllerBook "github.com/juninhoitabh/clob-go/internal/infra/controllers/book"
	repositoriesBook "github.com/juninhoitabh/clob-go/internal/infra/repositories/book"
)

func BookGenerate(router *http.ServeMux, apiV1Prefix string) {
	bookRepo := repositoriesBook.NewInMemoryBookRepository()

	controller := controllerBook.NewBookController(bookRepo)

	router.HandleFunc("GET "+apiV1Prefix+"/books", controller.Get)
}
