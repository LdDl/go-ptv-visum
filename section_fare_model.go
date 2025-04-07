package ptvvisum

import (
	"fmt"
	"strconv"
	"strings"
)

// FareModelSection represents $FAREMODEL section
type FareModelSection struct {
	BaseSection
	FallbackFare float64 // Default fare value when no specific fare is defined
}

// getFallbackFare extracts fallback fare value from FAREMODEL section
func getFallbackFare(values []string) (float64, error) {
	if len(values) < 1 || values[0] == "" {
		return 0.0, fmt.Errorf("invalid FAREMODEL data: %v", values)
	}

	// Parse the fallback fare value
	fallbackFare, err := strconv.ParseFloat(strings.Replace(values[0], ",", ".", -1), 64)
	if err != nil {
		return 0.0, fmt.Errorf("error parsing fallback fare: %w", err)
	}

	return fallbackFare, nil
}
