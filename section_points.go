package ptvvisum

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// PointSection represents $POINT section
type PointSection struct {
	BaseSection
	Points []Point
}

// Point represents a single point with coordinates
type Point struct {
	ID     int     // Point ID
	XCoord float64 // X-coordinate
	YCoord float64 // Y-coordinate
}

// GetPointByID retrieves a point by its ID
func (s *PointSection) GetPointByID(id int) (Point, bool) {
	for _, p := range s.Points {
		if p.ID == id {
			return p, true
		}
	}
	return Point{}, false
}

// GetPointsInArea returns all points within a bounding box
func (s *PointSection) GetPointsInArea(minX, minY, maxX, maxY float64) []Point {
	var result []Point
	for _, p := range s.Points {
		if p.XCoord >= minX && p.XCoord <= maxX &&
			p.YCoord >= minY && p.YCoord <= maxY {
			result = append(result, p)
		}
	}
	return result
}

// Distance returns the Euclidean distance between two points
func (p Point) Distance(other Point) float64 {
	dx := p.XCoord - other.XCoord
	dy := p.YCoord - other.YCoord
	return math.Sqrt(dx*dx + dy*dy)
}

// getPoint extracts data from POINT section row
func getPoint(values []string) (Point, error) {
	if len(values) < 3 {
		return Point{}, fmt.Errorf("invalid POINT data: %v", values)
	}

	var point Point
	var err error

	// Parse ID (required field)
	if values[0] == "" {
		return Point{}, fmt.Errorf("missing required field ID")
	}
	point.ID, err = strconv.Atoi(values[0])
	if err != nil {
		return Point{}, fmt.Errorf("error parsing ID: %w", err)
	}

	// Parse XCOORD (required field)
	if values[1] == "" {
		return Point{}, fmt.Errorf("missing required field XCOORD")
	}
	point.XCoord, err = strconv.ParseFloat(strings.Replace(values[1], ",", ".", -1), 64)
	if err != nil {
		return Point{}, fmt.Errorf("error parsing XCOORD: %w", err)
	}

	// Parse YCOORD (required field)
	if values[2] == "" {
		return Point{}, fmt.Errorf("missing required field YCOORD")
	}
	point.YCoord, err = strconv.ParseFloat(strings.Replace(values[2], ",", ".", -1), 64)
	if err != nil {
		return Point{}, fmt.Errorf("error parsing YCOORD: %w", err)
	}

	return point, nil
}
