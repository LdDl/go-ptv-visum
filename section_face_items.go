package ptvvisum

import (
	"fmt"
	"math"
	"sort"
	"strconv"
)

// FaceItemSection represents $FACEITEM section
type FaceItemSection struct {
	BaseSection
	Items []FaceItem
}

// FaceItem represents a single edge in a face definition
type FaceItem struct {
	FaceID    int // ID of the face this edge belongs to
	Index     int // Sequence number in the face
	EdgeID    int // ID of the edge
	Direction int // Direction of the edge (0 or 1)
}

// GetItemsByFaceID retrieves all items for a specific face
func (s *FaceItemSection) GetItemsByFaceID(faceID int) []FaceItem {
	var result []FaceItem
	for _, item := range s.Items {
		if item.FaceID == faceID {
			result = append(result, item)
		}
	}

	// Sort by index to ensure correct order
	sort.Slice(result, func(i, j int) bool {
		return result[i].Index < result[j].Index
	})

	return result
}

// GetFaceGeometry builds the complete geometry of a face
func (s *FaceItemSection) GetFaceGeometry(faceID int, data *PTVData) [][2]float64 {
	items := s.GetItemsByFaceID(faceID)
	if len(items) == 0 {
		return nil
	}

	var coordinates [][2]float64

	// Process each edge in order
	for _, item := range items {
		if data.Edge == nil || data.Point == nil {
			continue // Can't get geometry without edge and point data
		}

		// Get the edge
		edge, found := data.Edge.GetEdgeByID(item.EdgeID)
		if !found {
			continue
		}

		// Get points for the edge
		var fromPointID, toPointID int
		if item.Direction == 0 {
			// Regular direction
			fromPointID = edge.FromPointID
			toPointID = edge.ToPointID
		} else {
			// Reverse direction
			fromPointID = edge.ToPointID
			toPointID = edge.FromPointID
		}

		fromPoint, foundFrom := data.Point.GetPointByID(fromPointID)
		if !foundFrom {
			continue
		}

		// Add the point to our geometry
		coordinates = append(coordinates, [2]float64{fromPoint.XCoord, fromPoint.YCoord})

		// If this is the last edge, add the last point too
		if item.Index == len(items) {
			toPoint, foundTo := data.Point.GetPointByID(toPointID)
			if foundTo {
				coordinates = append(coordinates, [2]float64{toPoint.XCoord, toPoint.YCoord})
			}
		}
	}

	return coordinates
}

// CalculateFaceArea calculates the approximate area of a face
func (s *FaceItemSection) CalculateFaceArea(faceID int, data *PTVData) float64 {
	coords := s.GetFaceGeometry(faceID, data)
	if len(coords) < 3 {
		return 0 // Need at least 3 points for a polygon
	}

	// Close the polygon if needed
	if coords[0][0] != coords[len(coords)-1][0] || coords[0][1] != coords[len(coords)-1][1] {
		coords = append(coords, coords[0])
	}

	// Calculate area using Shoelace formula
	var area float64
	for i := 0; i < len(coords)-1; i++ {
		area += (coords[i][0] * coords[i+1][1]) - (coords[i+1][0] * coords[i][1])
	}

	return math.Abs(area) / 2.0
}

// getFaceItem extracts data from FACEITEM section row
func getFaceItem(values []string) (FaceItem, error) {
	if len(values) < 4 {
		return FaceItem{}, fmt.Errorf("invalid FACEITEM data: %v", values)
	}

	var item FaceItem
	var err error

	// Parse FACEID (required field)
	if values[0] == "" {
		return FaceItem{}, fmt.Errorf("missing required field FACEID")
	}
	item.FaceID, err = strconv.Atoi(values[0])
	if err != nil {
		return FaceItem{}, fmt.Errorf("error parsing FACEID: %w", err)
	}

	// Parse INDEX (required field)
	if values[1] == "" {
		return FaceItem{}, fmt.Errorf("missing required field INDEX")
	}
	item.Index, err = strconv.Atoi(values[1])
	if err != nil {
		return FaceItem{}, fmt.Errorf("error parsing INDEX: %w", err)
	}

	// Parse EDGEID (required field)
	if values[2] == "" {
		return FaceItem{}, fmt.Errorf("missing required field EDGEID")
	}
	item.EdgeID, err = strconv.Atoi(values[2])
	if err != nil {
		return FaceItem{}, fmt.Errorf("error parsing EDGEID: %w", err)
	}

	// Parse DIRECTION (required field)
	if values[3] == "" {
		return FaceItem{}, fmt.Errorf("missing required field DIRECTION")
	}
	item.Direction, err = strconv.Atoi(values[3])
	if err != nil {
		return FaceItem{}, fmt.Errorf("error parsing DIRECTION: %w", err)
	}

	return item, nil
}
