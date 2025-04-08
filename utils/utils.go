package utils

import (
	"fmt"
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
