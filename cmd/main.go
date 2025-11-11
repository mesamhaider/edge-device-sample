package main

import (
	"log"
	"net/http"

	"github.com/mesamhaider/edge-device-sample/internal/data"
	"github.com/mesamhaider/edge-device-sample/internal/handler"
	coreHttp "github.com/mesamhaider/edge-device-sample/internal/http"
)

func main() {
	storage, err := data.NewInMemoryStorageFromCSV("etc/devices.csv")
	if err != nil {
		log.Fatalf("failed to load devices: %v", err)
	}

	coreHandler := handler.NewCoreHandler(storage)

	router := coreHttp.NewRouter(coreHandler)

	const addr = ":8080"
	log.Printf("Starting HTTP server on %s with %d devices", addr, storage.DevicesCount())

	if err := http.ListenAndServe(addr, router); err != nil && err != http.ErrServerClosed {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
