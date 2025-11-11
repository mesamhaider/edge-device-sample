package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/mesamhaider/edge-device-sample/internal/handler"
)

func NewRouter(coreHandler *handler.CoreHandler) *chi.Mux {
	router := chi.NewRouter()

	router.Route("/devices/{device_id}", func(r chi.Router) {
		r.Post("/heartbeat", coreHandler.AddBeat)

		r.Route("/stats", func(ro chi.Router) {
			ro.Post("/", coreHandler.NewDeviceStats)
			ro.Get("/", coreHandler.GetDeviceStats)
		})
	})

	return router
}
