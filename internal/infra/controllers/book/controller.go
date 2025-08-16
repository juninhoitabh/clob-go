package book

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/juninhoitabh/clob-go/internal/shared"
)

type BookController struct {
}

func (b *BookController) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	inst := req.PathValue("instrument")

	inst = strings.ToUpper(inst)
	if inst == "" {
		http.Error(w, "instrument required", http.StatusBadRequest)

		return
	}

	book := b.Eng.SnapshotBook(inst)
	if book == nil {
		http.Error(w, fmt.Sprintf("instrument %s not found (empty book yet?)", inst), http.StatusNotFound)

		return
	}

	shared.WriteJSON(w, http.StatusOK, book)
}
