package handler

import (
	"net/http"
	"time"
)

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
