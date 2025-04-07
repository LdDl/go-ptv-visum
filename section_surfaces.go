package ptvvisum

import (
	"fmt"
	"strconv"
)

// SurfaceSection represents $SURFACE section
type SurfaceSection struct {
	BaseSection
	Surfaces []Surface
}

// Surface represents a single surface in the network
type Surface struct {
	ID int // Surface identifier
}

// GetSurfaceByID retrieves a surface by its ID
func (s *SurfaceSection) GetSurfaceByID(id int) (Surface, bool) {
	for _, surface := range s.Surfaces {
		if surface.ID == id {
			return surface, true
		}
	}
	return Surface{}, false
}

// Contains checks if a surface ID exists in the section
func (s *SurfaceSection) Contains(id int) bool {
	_, found := s.GetSurfaceByID(id)
	return found
}

// Count returns the number of surfaces in the section
func (s *SurfaceSection) Count() int {
	return len(s.Surfaces)
}

// GetAllIDs returns a slice containing all surface IDs
func (s *SurfaceSection) GetAllIDs() []int {
	ids := make([]int, len(s.Surfaces))
	for i, surface := range s.Surfaces {
		ids[i] = surface.ID
	}
	return ids
}

// getSurface extracts data from SURFACE section row
func getSurface(values []string) (Surface, error) {
	if len(values) < 1 {
		return Surface{}, fmt.Errorf("invalid SURFACE data: %v", values)
	}

	// Parse ID (required field)
	if values[0] == "" {
		return Surface{}, fmt.Errorf("missing required field ID")
	}

	id, err := strconv.Atoi(values[0])
	if err != nil {
		return Surface{}, fmt.Errorf("error parsing ID: %w", err)
	}

	return Surface{ID: id}, nil
}
