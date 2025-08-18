package main

import (
	"github.com/juninhoitabh/clob-go/internal/infra/config"
	httpServer "github.com/juninhoitabh/clob-go/internal/infra/http-server"
)

func main() {
	config.Init()
	httpServer.Start()
}
