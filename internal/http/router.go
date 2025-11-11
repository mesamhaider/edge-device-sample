package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/mesamhaider/edge-device-sample/internal/handler"
)

func NewRouter(logger *zap.Logger, coreHandler *handler.CoreHandler) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(RequestContextMiddleware(logger))

	router.Get("/health-check", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("I'm alive!"))
	})

	router.Route("/api/v1/devices/{device_id}", func(r chi.Router) {
		r.Post("/heartbeat", coreHandler.AddBeat)

		r.Route("/stats", func(ro chi.Router) {
			ro.Post("/", coreHandler.NewDeviceStats)
			ro.Get("/", coreHandler.GetDeviceStats)
		})
	})

	return router
}
