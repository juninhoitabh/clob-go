package router

import (
	"fmt"
	"net/http"

	"github.com/juninhoitabh/clob-go/docs"
	"github.com/juninhoitabh/clob-go/internal/infra/config"
	"github.com/juninhoitabh/clob-go/internal/infra/http-server/router/routes"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Clob API
// @description     Clob API without Authentication
// @termsOfService  http://swagger.io/terms/

// @contact.name   Junior Paz
func Generate(apiPort string) http.Handler {
	mux := http.NewServeMux()

	apiV1Prefix := "/api/v1"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", config.EnvConfigInstance.ApiHost, apiPort)
	docs.SwaggerInfo.BasePath = apiV1Prefix
	docs.SwaggerInfo.Schemes = []string{"http"}

	mux.Handle(apiV1Prefix+"/docs/", http.StripPrefix(apiV1Prefix, httpSwagger.WrapHandler))

	routes.AccountGenerate(mux, apiV1Prefix)
	routes.BookGenerate(mux, apiV1Prefix)
	routes.OrderGenerate(mux, apiV1Prefix)

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
