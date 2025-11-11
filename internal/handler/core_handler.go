package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/mesamhaider/edge-device-sample/internal/data"
	"github.com/mesamhaider/edge-device-sample/internal/services"
)

func NewCoreHandler(storage *data.InMemoryStorage) *CoreHandler {
	return &CoreHandler{storage: storage}
}

func (h *CoreHandler) AddBeat(w http.ResponseWriter, r *http.Request) {
	device, err := h.getDeviceFromRequest(r)
	if err != nil {
		h.writeError(w, err)
		return
	}

	var payload struct {
		SentAt string `json:"sent_at"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.writeError(w, newBadRequestError("invalid JSON body"))
		return
	}
	defer r.Body.Close()

	if payload.SentAt == "" {
		h.writeError(w, newBadRequestError("sent_at is required"))
		return
	}

	sentAt, err := time.Parse(time.RFC3339, payload.SentAt)
	if err != nil {
		h.writeError(w, newBadRequestError("sent_at must be RFC3339 formatted"))
		return
	}

	sentAt = sentAt.UTC().Truncate(time.Minute)

	device.Lock()
	defer device.Unlock()

	if device.HeartbeatMin == nil {
		device.HeartbeatMin = make(map[time.Time]struct{})
	}

	if _, exists := device.HeartbeatMin[sentAt]; !exists {
		device.HeartbeatMin[sentAt] = struct{}{}
	}

	if device.FirstHB.IsZero() || sentAt.Before(device.FirstHB) {
		device.FirstHB = sentAt
	}
	if device.LastHB.IsZero() || sentAt.After(device.LastHB) {
		device.LastHB = sentAt
	}

	w.WriteHeader(http.StatusOK)
}

// NewDeviceStats handles POST /devices/{device_id}/stats.
func (h *CoreHandler) NewDeviceStats(w http.ResponseWriter, r *http.Request) {
	device, err := h.getDeviceFromRequest(r)
	if err != nil {
		h.writeError(w, err)
		return
	}

	var payload struct {
		SentAt     string `json:"sent_at"`
		UploadTime int64  `json:"upload_time"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.writeError(w, newBadRequestError("invalid JSON body"))
		return
	}
	defer r.Body.Close()

	if payload.SentAt != "" {
		if _, err := time.Parse(time.RFC3339, payload.SentAt); err != nil {
			h.writeError(w, newBadRequestError("sent_at must be RFC3339 formatted"))
			return
		}
	}

	if payload.UploadTime < 0 {
		h.writeError(w, newBadRequestError("upload_time must be >= 0"))
		return
	}

	device.Lock()
	defer device.Unlock()

	device.UploadCount++
	device.UploadSumMs += payload.UploadTime

	w.WriteHeader(http.StatusOK)
}

// GetDeviceStats handles GET /devices/{device_id}/stats.
func (h *CoreHandler) GetDeviceStats(w http.ResponseWriter, r *http.Request) {
	device, err := h.getDeviceFromRequest(r)
	if err != nil {
		h.writeError(w, err)
		return
	}

	device.RLock()

	sumHeartbeat := len(device.HeartbeatMin)
	uptime := services.CalculateUptime(sumHeartbeat, device.FirstHB, device.LastHB)
	avgUploadMs := services.CalculateAverageUploadTime(device.UploadSumMs, device.UploadCount)

	device.RUnlock()

	resp := struct {
		Uptime        float64 `json:"uptime"`
		AvgUploadTime string  `json:"avg_upload_time"`
	}{
		Uptime:        uptime,
		AvgUploadTime: formatUploadTime(avgUploadMs),
	}

	h.writeJSON(w, http.StatusOK, resp)
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

func (h *CoreHandler) writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *CoreHandler) writeError(w http.ResponseWriter, err error) {
	var apiErr *apiError
	if errors.As(err, &apiErr) {
		h.writeJSON(w, apiErr.Status, map[string]string{
			"error": apiErr.Message,
		})
		return
	}

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

func formatUploadTime(ms int64) string {
	return time.Duration(ms * int64(time.Millisecond)).String()
}
