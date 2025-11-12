package services

import (
	"time"
)

func CalculateUptime(sumHeartbeat int, firstHB, lastHB time.Time) float64 {
	if sumHeartbeat == 0 {
		return 0.0
	}

	denominator := diffMinutesInclusive(firstHB, lastHB)
	if denominator < 1 {
		denominator = 1
	}

	return (float64(sumHeartbeat) / float64(denominator)) * 100.0
}

func CalculateAverageUploadTime(uploadSumNs int64, uploadCount int) time.Duration {
	if uploadCount <= 0 || uploadSumNs <= 0 {
		return 0
	}

	avg := uploadSumNs / int64(uploadCount)
	return time.Duration(avg)
}

func diffMinutesInclusive(start, end time.Time) int {
	if start.IsZero() || end.IsZero() {
		return 0
	}

	if end.Before(start) {
		start, end = end, start
	}

	diff := int(end.Sub(start) / time.Minute)
	return diff + 1
}
