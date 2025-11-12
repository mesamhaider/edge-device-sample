package services_test

import (
	"math"
	"testing"
	"time"

	"github.com/mesamhaider/edge-device-sample/internal/services"
)

func TestCalculateUptime(t *testing.T) {
	base := time.Date(2025, time.January, 1, 12, 0, 0, 0, time.UTC)
	last := base.Add(4 * time.Minute) // 5 minutes inclusive

	tests := []struct {
		name         string
		sumHeartbeat int
		firstHB      time.Time
		lastHB       time.Time
		want         float64
	}{
		{
			name:         "no heartbeats",
			sumHeartbeat: 0,
			firstHB:      base,
			lastHB:       last,
			want:         0.0,
		},
		{
			name:         "full uptime",
			sumHeartbeat: 5,
			firstHB:      base,
			lastHB:       last,
			want:         100.0,
		},
		{
			name:         "partial uptime",
			sumHeartbeat: 3,
			firstHB:      base,
			lastHB:       last,
			want:         60.0,
		},
		{
			name:         "mismatched heartbeat order",
			sumHeartbeat: 2,
			firstHB:      last,
			lastHB:       base,
			want:         (float64(2) / float64(5)) * 100.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := services.CalculateUptime(tt.sumHeartbeat, tt.firstHB, tt.lastHB)
			if math.Abs(got-tt.want) > 1e-6 {
				t.Errorf("CalculateUptime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateAverageUploadTime(t *testing.T) {
	tests := []struct {
		name  string
		sumNs int64
		count int
		want  time.Duration
	}{
		{
			name:  "no uploads",
			sumNs: 0,
			count: 0,
			want:  0,
		},
		{
			name:  "valid average",
			sumNs: (5 * time.Second).Nanoseconds(),
			count: 5,
			want:  time.Second,
		},
		{
			name:  "fractional average truncates",
			sumNs: (3500 * time.Millisecond).Nanoseconds(),
			count: 3,
			want:  time.Duration((3500 * time.Millisecond).Nanoseconds() / 3),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := services.CalculateAverageUploadTime(tt.sumNs, tt.count)
			if got != tt.want {
				t.Errorf("CalculateAverageUploadTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
