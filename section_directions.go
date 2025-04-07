package ptvvisum

import (
	"fmt"
	"strconv"
)

// DirectionSection represents $DIRECTION section
type DirectionSection struct {
	BaseSection
	Directions []Direction
}

// Direction represents a single direction entry
type Direction struct {
	No   int    // Direction number
	Code string // Direction code/symbol
	Name string // Direction name/description
}

// getDirection extracts data from DIRECTION section row
func getDirection(values []string) (Direction, error) {
	if len(values) < 3 {
		return Direction{}, fmt.Errorf("invalid DIRECTION data: %v", values)
	}

	// Parse NO (required field)
	no, err := strconv.Atoi(values[0])
	if err != nil {
		return Direction{}, fmt.Errorf("error parsing No: %w", err)
	}

	// Create direction with parsed values
	direction := Direction{
		No:   no,
		Code: values[1],
		Name: values[2],
	}

	return direction, nil
}
