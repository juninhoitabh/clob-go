package httpServer

import (
	"net/http"
	"net/http/httptest"
	"sync"

	"github.com/juninhoitabh/clob-go/internal/infra/config"
)

type E2eTestHandle struct {
	HttpServerTest *httptest.Server
	HttpHeader     http.Header
}

func NewE2eTestHandle() *E2eTestHandle {
	e2eTestHandle := &E2eTestHandle{
		HttpServerTest: httpServerInitTest(),
		HttpHeader:     http.Header{},
	}

	return e2eTestHandle
}

var once sync.Once
var instance *httptest.Server

func httpServerInitTest() *httptest.Server {
	once.Do(func() {
		config.Init()

		httpServer := &HttpServer{}

		httpHandler := httpServer.generateRoutes(config.EnvConfigInstance.ApiPort)

		instance = httptest.NewServer(httpHandler)
	})

	return instance
}
