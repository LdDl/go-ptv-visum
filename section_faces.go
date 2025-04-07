package ptvvisum

import (
	"fmt"
	"strconv"
)

// FaceSection represents $FACE section
type FaceSection struct {
	BaseSection
	Faces []Face
}

// Face represents a single face/polygon identifier
type Face struct {
	ID int // Face ID
}

// GetFaceByID retrieves a face by its ID
func (s *FaceSection) GetFaceByID(id int) (Face, bool) {
	for _, f := range s.Faces {
		if f.ID == id {
			return f, true
		}
	}
	return Face{}, false
}

// Contains checks if a face ID exists in the section
func (s *FaceSection) Contains(id int) bool {
	_, found := s.GetFaceByID(id)
	return found
}

// Count returns the number of faces in the section
func (s *FaceSection) Count() int {
	return len(s.Faces)
}

// getFace extracts data from FACE section row
func getFace(values []string) (Face, error) {
	if len(values) < 1 {
		return Face{}, fmt.Errorf("invalid FACE data: %v", values)
	}

	// Parse ID (required field)
	if values[0] == "" {
		return Face{}, fmt.Errorf("missing required field ID")
	}

	id, err := strconv.Atoi(values[0])
	if err != nil {
		return Face{}, fmt.Errorf("error parsing ID: %w", err)
	}

	return Face{ID: id}, nil
}
