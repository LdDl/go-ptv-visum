package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Helper function to check if a rune is a digit
func IsDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

// ReverseGeometry reverses the order of points in a geometry without modifying source slice
func ReverseGeometry(geom [][]float64) [][]float64 {
	result := make([][]float64, len(geom))
	for i, point := range geom {
		result[len(geom)-1-i] = make([]float64, 2)
		copy(result[len(geom)-1-i], point)
	}
	return result
}

// PointToWKT converts given point to Well-Known Text format
func PointToWKT(coord []float64) string {
	if len(coord) < 2 {
		return "POINT EMPTY"
	}
	return fmt.Sprintf("POINT(%g %g)", coord[0], coord[1])
}

// LineStringToWKT converts array of points to Well-Known Text format
func LineStringToWKT(coords [][]float64) string {
	if len(coords) == 0 {
		return "LINESTRING EMPTY"
	}
	var parts []string
	for _, coord := range coords {
		if len(coord) >= 2 {
			parts = append(parts, fmt.Sprintf("%g %g", coord[0], coord[1]))
		}
	}
	if len(parts) == 0 {
		return "LINESTRING EMPTY"
	}
	return fmt.Sprintf("LINESTRING(%s)", strings.Join(parts, ", "))
}

// ParseSpeedValue extracts numeric value from string with units and converts to km/h
func ParseSpeedValue(speedStr string) (float64, error) {
	if speedStr == "" {
		return 0, nil
	}

	// Use regex to extract number and unit
	re := regexp.MustCompile(`^([\d.,]+)\s*([a-zA-Z/]+)`)
	matches := re.FindStringSubmatch(speedStr)

	if len(matches) < 2 {
		return 0, fmt.Errorf("failed to parse speed value '%s'", speedStr)
	}

	// Parse numeric part
	numStr := strings.ReplaceAll(matches[1], ",", ".")
	value, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse numeric part of '%s': %w", speedStr, err)
	}

	// Get unit (if present)
	unit := ""
	if len(matches) > 2 {
		unit = strings.ToLower(matches[2])
	}

	// Common speed units
	switch {
	case strings.Contains(unit, "km/h") || unit == "":
		return value, nil // Already km/h
	case strings.Contains(unit, "m/s"):
		return value * 3.6, nil // m/s to km/h (1 m/s = 3.6 km/h)
	case strings.Contains(unit, "mph") || strings.Contains(unit, "mi/h"):
		return value * 1.60934, nil // mph to km/h
	case strings.Contains(unit, "km/min"):
		return value * 60, nil // km/min to km/h
	case strings.Contains(unit, "m/min"):
		return value * 0.06, nil // m/min to km/h
	case strings.Contains(unit, "ft/s"):
		return value * 1.09728, nil // ft/s to km/h
	default:
		// If unrecognized, default to km/h
		return value, nil
	}
}

// ParseLengthValue extracts numeric value from string with units and converts to meters
func ParseLengthValue(lengthStr string) (float64, error) {
	if lengthStr == "" {
		return 0, nil
	}

	// Use regex to extract number and unit
	re := regexp.MustCompile(`^([\d.,]+)\s*([a-zA-Z]*)`)
	matches := re.FindStringSubmatch(lengthStr)

	if len(matches) < 2 {
		return 0, fmt.Errorf("failed to parse length value '%s'", lengthStr)
	}

	// Parse numeric part
	numStr := strings.ReplaceAll(matches[1], ",", ".")
	value, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse numeric part of '%s': %w", lengthStr, err)
	}

	// Get unit (if present)
	unit := ""
	if len(matches) > 2 {
		unit = strings.ToLower(matches[2])
	}

	// Convert to meters based on unit
	switch {
	case strings.Contains(unit, "km"):
		return value * 1000, nil
	case strings.Contains(unit, "m") || unit == "":
		return value, nil
	case strings.Contains(unit, "cm"):
		return value / 100, nil
	case strings.Contains(unit, "mm"):
		return value / 1000, nil
	case strings.Contains(unit, "mi"):
		return value * 1609.34, nil // miles to meters
	case strings.Contains(unit, "ft"):
		return value * 0.3048, nil // feet to meters
	default:
		// For unrecognized unit, assume meters
		return value, nil
	}
}
