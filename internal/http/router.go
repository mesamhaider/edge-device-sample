package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mesamhaider/edge-device-sample/internal/handler"
)

func NewRouter(coreHandler *handler.CoreHandler) *chi.Mux {
	router := chi.NewRouter()

	router.Get("/health-check", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("I'm alive!"))
	})

	router.Route("/devices/{device_id}", func(r chi.Router) {
		r.Post("/heartbeat", coreHandler.AddBeat)

		r.Route("/stats", func(ro chi.Router) {
			ro.Post("/", coreHandler.NewDeviceStats)
			ro.Get("/", coreHandler.GetDeviceStats)
		})
	})

	return router
}
