package ptvvisum

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
)

// EdgeItemSection represents $EDGEITEM section (intermediate points of edges)
type EdgeItemSection struct {
	BaseSection
	Items []EdgeItem
}

// EdgeItem represents a single intermediate point on an edge
type EdgeItem struct {
	EdgeID int     // Edge ID this point belongs to
	Index  int     // Sequence number within the edge
	XCoord float64 // X-coordinate
	YCoord float64 // Y-coordinate
}

// GetItemsByEdgeID retrieves all intermediate points for a specific edge
func (s *EdgeItemSection) GetItemsByEdgeID(edgeID int) []EdgeItem {
	var result []EdgeItem
	for _, item := range s.Items {
		if item.EdgeID == edgeID {
			result = append(result, item)
		}
	}

	// Sort by index to ensure correct order
	sort.Slice(result, func(i, j int) bool {
		return result[i].Index < result[j].Index
	})

	return result
}

// GetEdgeGeometry returns the complete geometry of an edge as a series of coordinates
func (s *EdgeItemSection) GetEdgeGeometry(edgeID int, data *PTVData) [][2]float64 {
	// Get all intermediate points for this edge
	items := s.GetItemsByEdgeID(edgeID)

	// Find the edge to get start and end points
	var edge Edge
	var found bool
	if data.Edge != nil {
		edge, found = data.Edge.GetEdgeByID(edgeID)
	}

	var geometry [][2]float64

	// If we have the edge and points data, add the actual start point
	if found && data.Point != nil {
		if fromPoint, foundFrom := data.Point.GetPointByID(edge.FromPointID); foundFrom {
			geometry = append(geometry, [2]float64{fromPoint.XCoord, fromPoint.YCoord})
		}
	}

	// Add all intermediate points in order
	for _, item := range items {
		geometry = append(geometry, [2]float64{item.XCoord, item.YCoord})
	}

	// If we have the edge and points data, add the actual end point
	if found && data.Point != nil {
		if toPoint, foundTo := data.Point.GetPointByID(edge.ToPointID); foundTo {
			geometry = append(geometry, [2]float64{toPoint.XCoord, toPoint.YCoord})
		}
	}

	return geometry
}

// CalculateEdgeLength calculates the total length of an edge based on its geometry
func (s *EdgeItemSection) CalculateEdgeLength(edgeID int, data *PTVData) float64 {
	geometry := s.GetEdgeGeometry(edgeID, data)
	if len(geometry) < 2 {
		return 0.0
	}

	var length float64
	for i := 1; i < len(geometry); i++ {
		dx := geometry[i][0] - geometry[i-1][0]
		dy := geometry[i][1] - geometry[i-1][1]
		segmentLength := math.Sqrt(dx*dx + dy*dy)
		length += segmentLength
	}

	return length
}

// getEdgeItem extracts data from EDGEITEM section row
func getEdgeItem(values []string) (EdgeItem, error) {
	if len(values) < 4 {
		return EdgeItem{}, fmt.Errorf("invalid EDGEITEM data (need EDGEID;INDEX;XCOORD;YCOORD): %v", values)
	}

	var item EdgeItem
	var err error

	// Parse EDGEID (required field)
	if values[0] == "" {
		return EdgeItem{}, fmt.Errorf("missing required field EDGEID")
	}
	item.EdgeID, err = strconv.Atoi(values[0])
	if err != nil {
		return EdgeItem{}, fmt.Errorf("error parsing EDGEID: %w", err)
	}

	// Parse INDEX (required field)
	if values[1] == "" {
		return EdgeItem{}, fmt.Errorf("missing required field INDEX")
	}
	item.Index, err = strconv.Atoi(values[1])
	if err != nil {
		return EdgeItem{}, fmt.Errorf("error parsing INDEX: %w", err)
	}

	// Parse XCOORD (required field)
	if values[2] == "" {
		return EdgeItem{}, fmt.Errorf("missing required field XCOORD")
	}
	item.XCoord, err = strconv.ParseFloat(strings.Replace(values[2], ",", ".", -1), 64)
	if err != nil {
		return EdgeItem{}, fmt.Errorf("error parsing XCOORD: %w", err)
	}

	// Parse YCOORD (required field)
	if values[3] == "" {
		return EdgeItem{}, fmt.Errorf("missing required field YCOORD")
	}
	item.YCoord, err = strconv.ParseFloat(strings.Replace(values[3], ",", ".", -1), 64)
	if err != nil {
		return EdgeItem{}, fmt.Errorf("error parsing YCOORD: %w", err)
	}

	return item, nil
}
