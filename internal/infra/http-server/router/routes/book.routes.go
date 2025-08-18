package routes

import (
	"net/http"

	controllerBook "github.com/juninhoitabh/clob-go/internal/infra/controllers/book"
	repositoriesBook "github.com/juninhoitabh/clob-go/internal/infra/repositories/book"
)

func BookGenerate(router *http.ServeMux) {
	bookRepo := repositoriesBook.NewInMemoryBookRepository()

	controller := controllerBook.NewBookController(bookRepo)

	router.HandleFunc("GET /book/{instrument...}", controller.Get)
}
