package ptvvisum

import (
	"fmt"
	"strconv"
)

// SurfaceItemSection represents $SURFACEITEM section
type SurfaceItemSection struct {
	BaseSection
	Items []SurfaceItem
}

// SurfaceItem represents a face that is part of a surface
type SurfaceItem struct {
	SurfaceID int // ID of the surface
	FaceID    int // ID of the face
	Enclave   int // Whether this face is an enclave (inner hole) (0=outer boundary, 1=inner hole)
}

// GetItemsBySurfaceID retrieves all items for a specific surface
func (s *SurfaceItemSection) GetItemsBySurfaceID(surfaceID int) []SurfaceItem {
	var result []SurfaceItem
	for _, item := range s.Items {
		if item.SurfaceID == surfaceID {
			result = append(result, item)
		}
	}
	return result
}

// GetItemsByFaceID retrieves all items for a specific face
func (s *SurfaceItemSection) GetItemsByFaceID(faceID int) []SurfaceItem {
	var result []SurfaceItem
	for _, item := range s.Items {
		if item.FaceID == faceID {
			result = append(result, item)
		}
	}
	return result
}

// GetBoundariesBySurfaceID gets the outer boundary and inner holes (enclaves) for a surface
func (s *SurfaceItemSection) GetBoundariesBySurfaceID(surfaceID int) (outerFaces []int, innerFaces []int) {
	items := s.GetItemsBySurfaceID(surfaceID)

	for _, item := range items {
		if item.Enclave == 0 {
			outerFaces = append(outerFaces, item.FaceID)
		} else {
			innerFaces = append(innerFaces, item.FaceID)
		}
	}

	return outerFaces, innerFaces
}

// GetSurfaceIDs returns all unique surface IDs in the section
func (s *SurfaceItemSection) GetSurfaceIDs() []int {
	surfaceMap := make(map[int]bool)
	for _, item := range s.Items {
		surfaceMap[item.SurfaceID] = true
	}

	result := make([]int, 0, len(surfaceMap))
	for id := range surfaceMap {
		result = append(result, id)
	}

	return result
}

// getSurfaceItem extracts data from SURFACEITEM section row
func getSurfaceItem(values []string) (SurfaceItem, error) {
	if len(values) < 3 {
		return SurfaceItem{}, fmt.Errorf("invalid SURFACEITEM data: %v", values)
	}

	var item SurfaceItem
	var err error

	// Parse SurfaceID (required field)
	if values[0] == "" {
		return SurfaceItem{}, fmt.Errorf("missing required field SURFACEID")
	}
	item.SurfaceID, err = strconv.Atoi(values[0])
	if err != nil {
		return SurfaceItem{}, fmt.Errorf("error parsing SURFACEID: %w", err)
	}

	// Parse FaceID (required field)
	if values[1] == "" {
		return SurfaceItem{}, fmt.Errorf("missing required field FACEID")
	}
	item.FaceID, err = strconv.Atoi(values[1])
	if err != nil {
		return SurfaceItem{}, fmt.Errorf("error parsing FACEID: %w", err)
	}

	// Parse Enclave (required field)
	if values[2] == "" {
		return SurfaceItem{}, fmt.Errorf("missing required field ENCLAVE")
	}
	item.Enclave, err = strconv.Atoi(values[2])
	if err != nil {
		return SurfaceItem{}, fmt.Errorf("error parsing ENCLAVE: %w", err)
	}

	return item, nil
}
