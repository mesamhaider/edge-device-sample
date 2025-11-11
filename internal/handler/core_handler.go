package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/mesamhaider/edge-device-sample/internal/data"
	"github.com/mesamhaider/edge-device-sample/internal/logging"
	"github.com/mesamhaider/edge-device-sample/internal/services"
)

func NewCoreHandler(storage *data.InMemoryStorage, logger *zap.Logger) *CoreHandler {
	return &CoreHandler{
		storage: storage,
		logger:  logger,
	}
}

func (h *CoreHandler) AddBeat(w http.ResponseWriter, r *http.Request) {
	device, err := h.getDeviceFromRequest(r)
	if err != nil {
		h.writeError(w, r, err)
		return
	}

	var payload struct {
		SentAt string `json:"sent_at"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.writeError(w, r, newBadRequestError("invalid JSON body"))
		return
	}
	defer r.Body.Close()

	if payload.SentAt == "" {
		h.writeError(w, r, newBadRequestError("sent_at is required"))
		return
	}

	sentAt, err := time.Parse(time.RFC3339, payload.SentAt)
	if err != nil {
		h.writeError(w, r, newBadRequestError("sent_at must be RFC3339 formatted"))
		return
	}

	sentAt = sentAt.UTC().Truncate(time.Minute)

	device.Lock()
	defer device.Unlock()

	if device.HeartbeatMin == nil {
		device.HeartbeatMin = make(map[time.Time]struct{})
	}

	created := false
	if _, exists := device.HeartbeatMin[sentAt]; !exists {
		device.HeartbeatMin[sentAt] = struct{}{}
		created = true
	}

	if device.FirstHB.IsZero() || sentAt.Before(device.FirstHB) {
		device.FirstHB = sentAt
	}
	if device.LastHB.IsZero() || sentAt.After(device.LastHB) {
		device.LastHB = sentAt
	}

	logger := logging.LoggerFromContext(r.Context(), h.logger)
	logger.Info("heartbeat recorded",
		zap.String("device_id", device.ID),
		zap.Time("sent_at", sentAt),
		zap.Bool("new_minute", created),
	)

	w.WriteHeader(http.StatusOK)
}

// NewDeviceStats handles POST /devices/{device_id}/stats.
func (h *CoreHandler) NewDeviceStats(w http.ResponseWriter, r *http.Request) {
	device, err := h.getDeviceFromRequest(r)
	if err != nil {
		h.writeError(w, r, err)
		return
	}

	var payload struct {
		SentAt     string `json:"sent_at"`
		UploadTime int64  `json:"upload_time"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.writeError(w, r, newBadRequestError("invalid JSON body"))
		return
	}
	defer r.Body.Close()

	if payload.SentAt != "" {
		if _, err := time.Parse(time.RFC3339, payload.SentAt); err != nil {
			h.writeError(w, r, newBadRequestError("sent_at must be RFC3339 formatted"))
			return
		}
	}

	if payload.UploadTime < 0 {
		h.writeError(w, r, newBadRequestError("upload_time must be >= 0"))
		return
	}

	device.Lock()
	defer device.Unlock()

	device.UploadCount++
	device.UploadSumNs += payload.UploadTime

	logger := logging.LoggerFromContext(r.Context(), h.logger)
	logger.Info("stats recorded",
		zap.String("device_id", device.ID),
		zap.Int64("upload_time_ns", payload.UploadTime),
		zap.Int("upload_count", device.UploadCount),
		zap.Int64("upload_sum_ns", device.UploadSumNs),
	)

	w.WriteHeader(http.StatusOK)
}

// GetDeviceStats handles GET /devices/{device_id}/stats.
func (h *CoreHandler) GetDeviceStats(w http.ResponseWriter, r *http.Request) {
	device, err := h.getDeviceFromRequest(r)
	if err != nil {
		h.writeError(w, r, err)
		return
	}

	device.RLock()

	sumHeartbeat := len(device.HeartbeatMin)
	uptime := services.CalculateUptime(sumHeartbeat, device.FirstHB, device.LastHB)
	avgUploadDur := services.CalculateAverageUploadTime(device.UploadSumNs, device.UploadCount)

	device.RUnlock()

	resp := struct {
		Uptime        float64 `json:"uptime"`
		AvgUploadTime string  `json:"avg_upload_time"`
	}{
		Uptime:        uptime,
		AvgUploadTime: formatUploadTime(avgUploadDur),
	}

	logger := logging.LoggerFromContext(r.Context(), h.logger)
	logger.Info("stats retrieved",
		zap.String("device_id", device.ID),
		zap.Float64("uptime", uptime),
		zap.Duration("avg_upload_time", avgUploadDur),
	)

	h.writeJSON(r.Context(), w, http.StatusOK, resp)
}

func (h *CoreHandler) getDeviceFromRequest(r *http.Request) (*data.Device, error) {
	deviceID := chi.URLParam(r, "device_id")
	if deviceID == "" {
		return nil, newBadRequestError("device_id is required")
	}

	device, ok := h.storage.GetDevice(deviceID)
	if !ok {
		return nil, newNotFoundError("device not found")
	}

	return device, nil
}

func (h *CoreHandler) writeJSON(ctx context.Context, w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		logger := logging.LoggerFromContext(ctx, h.logger)
		logger.Error("failed to encode response payload",
			zap.Error(err),
		)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *CoreHandler) writeError(w http.ResponseWriter, r *http.Request, err error) {
	var apiErr *apiError
	if errors.As(err, &apiErr) {
		h.writeJSON(r.Context(), w, apiErr.Status, map[string]string{
			"error": apiErr.Message,
		})
		logger := logging.LoggerFromContext(r.Context(), h.logger)
		logger.Warn("request error",
			zap.String("error", apiErr.Message),
			zap.Int("status", apiErr.Status),
		)
		return
	}

	logger := logging.LoggerFromContext(r.Context(), h.logger)
	logger.Error("unexpected error handling request", zap.Error(err))
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

type apiError struct {
	Status  int
	Message string
}

func (e *apiError) Error() string {
	return e.Message
}

func newBadRequestError(message string) *apiError {
	return &apiError{Status: http.StatusBadRequest, Message: message}
}

func newNotFoundError(message string) *apiError {
	return &apiError{Status: http.StatusNotFound, Message: message}
}

func formatUploadTime(d time.Duration) string {
	return d.String()
}
