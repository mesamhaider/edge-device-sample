package data

import (
	"sync"
	"time"
)

type Device struct {
	sync.RWMutex

	ID           string
	HeartbeatMin map[time.Time]struct{}
	FirstHB      time.Time
	LastHB       time.Time
	UploadCount  int
	UploadSumNs  int64
}

type InMemoryStorage struct {
	mu      sync.RWMutex
	devices map[string]*Device
}
