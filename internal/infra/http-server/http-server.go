package httpServer

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/juninhoitabh/clob-go/internal/infra/config"
	"github.com/juninhoitabh/clob-go/internal/infra/http-server/router"
)

type HttpServer struct{}

func (s *HttpServer) generateRoutes(apiPort string) http.Handler {
	return router.Generate(apiPort)
}

func Start() {
	httpServer := &HttpServer{}

	apiPort := config.EnvConfigInstance.ApiPort

	server := &http.Server{Addr: fmt.Sprintf(":%s", apiPort), Handler: httpServer.generateRoutes(apiPort)}

	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go gracefulShutdown(sig, serverCtx, server, serverStopCtx)

	fmt.Printf("Http Server is starting on port %s\n", apiPort)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	<-serverCtx.Done()
}

func gracefulShutdown(sig chan os.Signal, serverCtx context.Context, server *http.Server, serverStopCtx context.CancelFunc) {
	<-sig

	//nolint:govet
	shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

	go func() {
		<-shutdownCtx.Done()

		if shutdownCtx.Err() == context.DeadlineExceeded {
			log.Fatal("graceful shutdown timed out.. forcing exit.")
		}
	}()

	err := server.Shutdown(shutdownCtx)
	if err != nil {
		log.Fatal(err)
	}

	serverStopCtx()
}
