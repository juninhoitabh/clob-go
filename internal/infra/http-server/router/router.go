package router

import (
	"net/http"

	"github.com/juninhoitabh/clob-go/internal/infra/http-server/router/routes"
)

func Generate(apiPort string) http.Handler {
	mux := http.NewServeMux()

	routes.AccountGenerate(mux)
	routes.BookGenerate(mux)
	routes.OrderGenerate(mux)

	// Handle 404 for any unmatched routes
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	return withJSON(withRecover(mux))
}

func withJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func withRecover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				http.Error(w, "internal error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
