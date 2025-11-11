package handler

import (
	"github.com/mesamhaider/edge-device-sample/internal/data"
	"go.uber.org/zap"
)

type CoreHandler struct {
	storage *data.InMemoryStorage
	logger  *zap.Logger
}
