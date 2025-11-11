package data

import (
	"encoding/csv"
	"fmt"
	"os"
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
	UploadSumMs  int64
}

type InMemoryStorage struct {
	mu      sync.RWMutex
	devices map[string]*Device
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		devices: make(map[string]*Device),
	}
}

func NewInMemoryStorageFromCSV(path string) (*InMemoryStorage, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open devices csv: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("read devices csv: %w", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("devices csv: empty file")
	}

	storage := NewInMemoryStorage()

	for i, row := range records {
		if i == 0 {
			// header row, skip it
			continue
		}
		if len(row) == 0 {
			continue
		}
		deviceID := row[0]
		if deviceID == "" {
			continue
		}

		storage.AddDevice(deviceID)
	}

	return storage, nil
}

func (s *InMemoryStorage) AddDevice(deviceID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.devices[deviceID]; exists {
		return
	}

	s.devices[deviceID] = &Device{
		ID:           deviceID,
		HeartbeatMin: make(map[time.Time]struct{}),
	}
}

func (s *InMemoryStorage) GetDevice(deviceID string) (*Device, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	device, ok := s.devices[deviceID]
	return device, ok
}

func (s *InMemoryStorage) DevicesCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.devices)
}
