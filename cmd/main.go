package main

import (
	"log"
	"net/http"

	"go.uber.org/zap"

	"github.com/mesamhaider/edge-device-sample/internal/data"
	"github.com/mesamhaider/edge-device-sample/internal/handler"
	coreHttp "github.com/mesamhaider/edge-device-sample/internal/http"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}

	defer func() {
		_ = logger.Sync()
	}()

	storage, err := data.NewInMemoryStorageFromCSV("etc/devices.csv")
	if err != nil {
		logger.Fatal("failed to load devices", zap.Error(err))
	}

	coreHandler := handler.NewCoreHandler(storage, logger)

	router := coreHttp.NewRouter(logger, coreHandler)

	const addr = ":8080"
	logger.Info("starting HTTP server",
		zap.String("addr", addr),
		zap.Int("device_count", storage.DevicesCount()),
	)

	if err := http.ListenAndServe(addr, router); err != nil && err != http.ErrServerClosed {
		logger.Fatal("HTTP server failed", zap.Error(err))
	}
}
