package ptvvisum

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
)

// LinkPolySection represents $LINKPOLY section
type LinkPolySection struct {
	BaseSection
	Points []LinkPolyPoint
}

// LinkPolyPoint represents a single point in a link polygon
type LinkPolyPoint struct {
	FromNodeNo int     // Origin node ID
	ToNodeNo   int     // Destination node ID
	Index      int     // Sequential index of the point
	XCoord     float64 // X coordinate
	YCoord     float64 // Y coordinate
	ZCoord     float64 // Z coordinate (elevation)
}

// GetPointsByLink retrieves all points for a specific link (from node to node)
func (s *LinkPolySection) GetPointsByLink(fromNodeNo, toNodeNo int) []LinkPolyPoint {
	var result []LinkPolyPoint
	for _, point := range s.Points {
		if point.FromNodeNo == fromNodeNo && point.ToNodeNo == toNodeNo {
			result = append(result, point)
		}
	}

	// Sort by index to ensure proper order
	sort.Slice(result, func(i, j int) bool {
		return result[i].Index < result[j].Index
	})

	return result
}

// GetLinkGeometry returns the complete geometry of a link as a series of coordinates
func (s *LinkPolySection) GetLinkGeometry(fromNodeNo, toNodeNo int) [][3]float64 {
	points := s.GetPointsByLink(fromNodeNo, toNodeNo)
	if len(points) == 0 {
		return nil
	}

	geometry := make([][3]float64, len(points))
	for i, point := range points {
		geometry[i] = [3]float64{point.XCoord, point.YCoord, point.ZCoord}
	}

	return geometry
}

// CalculateLinkLength calculates the length of a link based on its geometry
func (s *LinkPolySection) CalculateLinkLength(fromNodeNo, toNodeNo int) float64 {
	geometry := s.GetLinkGeometry(fromNodeNo, toNodeNo)
	if len(geometry) < 2 {
		return 0
	}

	var length float64
	for i := 1; i < len(geometry); i++ {
		dx := geometry[i][0] - geometry[i-1][0]
		dy := geometry[i][1] - geometry[i-1][1]
		dz := geometry[i][2] - geometry[i-1][2]
		segmentLength := math.Sqrt(dx*dx + dy*dy + dz*dz)
		length += segmentLength
	}

	return length
}

// GetAllLinks returns all unique link identifiers (fromNodeNo, toNodeNo pairs)
func (s *LinkPolySection) GetAllLinks() [][2]int {
	linkMap := make(map[[2]int]bool)

	for _, point := range s.Points {
		key := [2]int{point.FromNodeNo, point.ToNodeNo}
		linkMap[key] = true
	}

	links := make([][2]int, 0, len(linkMap))
	for link := range linkMap {
		links = append(links, link)
	}

	return links
}

// HasLinkGeometry checks if detailed geometry exists for a specific link
func (s *LinkPolySection) HasLinkGeometry(fromNodeNo, toNodeNo int) bool {
	for _, point := range s.Points {
		if point.FromNodeNo == fromNodeNo && point.ToNodeNo == toNodeNo {
			return true
		}
	}
	return false
}

// GetBoundingBox returns the min/max coordinates of all link polygon points
func (s *LinkPolySection) GetBoundingBox() (minX, minY, minZ, maxX, maxY, maxZ float64) {
	if len(s.Points) == 0 {
		return 0, 0, 0, 0, 0, 0
	}

	minX = s.Points[0].XCoord
	minY = s.Points[0].YCoord
	minZ = s.Points[0].ZCoord
	maxX = s.Points[0].XCoord
	maxY = s.Points[0].YCoord
	maxZ = s.Points[0].ZCoord

	for _, point := range s.Points {
		if point.XCoord < minX {
			minX = point.XCoord
		}
		if point.YCoord < minY {
			minY = point.YCoord
		}
		if point.ZCoord < minZ {
			minZ = point.ZCoord
		}
		if point.XCoord > maxX {
			maxX = point.XCoord
		}
		if point.YCoord > maxY {
			maxY = point.YCoord
		}
		if point.ZCoord > maxZ {
			maxZ = point.ZCoord
		}
	}

	return minX, minY, minZ, maxX, maxY, maxZ
}

// Count returns the number of link polygon points in the section
func (s *LinkPolySection) Count() int {
	return len(s.Points)
}

// CountLinks returns the number of unique links with geometry
func (s *LinkPolySection) CountLinks() int {
	return len(s.GetAllLinks())
}

// getLinkPolyPoint extracts data from LINKPOLY section row
func getLinkPolyPoint(values []string) (LinkPolyPoint, error) {
	if len(values) < 6 {
		return LinkPolyPoint{}, fmt.Errorf("invalid LINKPOLY data (insufficient fields): %v", values)
	}

	var point LinkPolyPoint
	var err error

	// Parse FROMNODENO (required field)
	if values[0] == "" {
		return LinkPolyPoint{}, fmt.Errorf("missing required field FROMNODENO")
	}
	point.FromNodeNo, err = strconv.Atoi(values[0])
	if err != nil {
		return LinkPolyPoint{}, fmt.Errorf("error parsing FROMNODENO: %w", err)
	}

	// Parse TONODENO (required field)
	if values[1] == "" {
		return LinkPolyPoint{}, fmt.Errorf("missing required field TONODENO")
	}
	point.ToNodeNo, err = strconv.Atoi(values[1])
	if err != nil {
		return LinkPolyPoint{}, fmt.Errorf("error parsing TONODENO: %w", err)
	}

	// Parse INDEX (required field)
	if values[2] == "" {
		return LinkPolyPoint{}, fmt.Errorf("missing required field INDEX")
	}
	point.Index, err = strconv.Atoi(values[2])
	if err != nil {
		return LinkPolyPoint{}, fmt.Errorf("error parsing INDEX: %w", err)
	}

	// Parse XCOORD (required field)
	if values[3] == "" {
		return LinkPolyPoint{}, fmt.Errorf("missing required field XCOORD")
	}
	point.XCoord, err = strconv.ParseFloat(strings.Replace(values[3], ",", ".", -1), 64)
	if err != nil {
		return LinkPolyPoint{}, fmt.Errorf("error parsing XCOORD: %w", err)
	}

	// Parse YCOORD (required field)
	if values[4] == "" {
		return LinkPolyPoint{}, fmt.Errorf("missing required field YCOORD")
	}
	point.YCoord, err = strconv.ParseFloat(strings.Replace(values[4], ",", ".", -1), 64)
	if err != nil {
		return LinkPolyPoint{}, fmt.Errorf("error parsing YCOORD: %w", err)
	}

	// Parse ZCOORD (required field)
	if values[5] == "" {
		return LinkPolyPoint{}, fmt.Errorf("missing required field ZCOORD")
	}
	point.ZCoord, err = strconv.ParseFloat(strings.Replace(values[5], ",", ".", -1), 64)
	if err != nil {
		return LinkPolyPoint{}, fmt.Errorf("error parsing ZCOORD: %w", err)
	}

	return point, nil
}
