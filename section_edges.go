package ptvvisum

import (
	"fmt"
	"strconv"
)

// EdgeSection represents $EDGE section
type EdgeSection struct {
	BaseSection
	Edges []Edge
}

// Edge represents a single edge/connection between points
type Edge struct {
	ID          int // Edge ID
	FromPointID int // Starting point ID
	ToPointID   int // Ending point ID
}

// GetEdgeByID retrieves an edge by its ID
func (s *EdgeSection) GetEdgeByID(id int) (Edge, bool) {
	for _, e := range s.Edges {
		if e.ID == id {
			return e, true
		}
	}
	return Edge{}, false
}

// GetEdgesByPointID retrieves all edges connected to a specific point
func (s *EdgeSection) GetEdgesByPointID(pointID int) []Edge {
	var result []Edge
	for _, e := range s.Edges {
		if e.FromPointID == pointID || e.ToPointID == pointID {
			result = append(result, e)
		}
	}
	return result
}

// GetOutgoingEdges retrieves all edges starting from a specific point
func (s *EdgeSection) GetOutgoingEdges(pointID int) []Edge {
	var result []Edge
	for _, e := range s.Edges {
		if e.FromPointID == pointID {
			result = append(result, e)
		}
	}
	return result
}

// GetIncomingEdges retrieves all edges ending at a specific point
func (s *EdgeSection) GetIncomingEdges(pointID int) []Edge {
	var result []Edge
	for _, e := range s.Edges {
		if e.ToPointID == pointID {
			result = append(result, e)
		}
	}
	return result
}

// getEdge extracts data from EDGE section row
func getEdge(values []string) (Edge, error) {
	if len(values) < 3 {
		return Edge{}, fmt.Errorf("invalid EDGE data: %v", values)
	}

	var edge Edge
	var err error

	// Parse ID (required field)
	if values[0] == "" {
		return Edge{}, fmt.Errorf("missing required field ID")
	}
	edge.ID, err = strconv.Atoi(values[0])
	if err != nil {
		return Edge{}, fmt.Errorf("error parsing ID: %w", err)
	}

	// Parse FROMPOINTID (required field)
	if values[1] == "" {
		return Edge{}, fmt.Errorf("missing required field FROMPOINTID")
	}
	edge.FromPointID, err = strconv.Atoi(values[1])
	if err != nil {
		return Edge{}, fmt.Errorf("error parsing FROMPOINTID: %w", err)
	}

	// Parse TOPOINTID (required field)
	if values[2] == "" {
		return Edge{}, fmt.Errorf("missing required field TOPOINTID")
	}
	edge.ToPointID, err = strconv.Atoi(values[2])
	if err != nil {
		return Edge{}, fmt.Errorf("error parsing TOPOINTID: %w", err)
	}

	return edge, nil
}
