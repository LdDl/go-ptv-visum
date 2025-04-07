package ptvvisum

import (
	"fmt"
	"strconv"
	"strings"
)

// ModeSection represents $MODE section
type ModeSection struct {
	BaseSection
	Modes []Mode
}

// Mode represents a single transport mode
type Mode struct {
	Code            string   // Mode code
	Name            string   // Mode name
	TSysSet         []string // Set of transport systems used by this mode
	Interchangeable int      // Whether mode is interchangeable with other modes
}

// getMode extracts data from MODE section row
func getMode(values []string) (Mode, error) {
	if len(values) < 4 {
		return Mode{}, fmt.Errorf("invalid MODE data (need CODE;NAME;TSYSSET;INTERCHANGEABLE): %v", values)
	}

	// Process mode fields
	mode := Mode{
		Code:    values[0],
		Name:    values[1],
		TSysSet: []string{},
	}

	// Parse TSysSet (comma-separated list of transport systems)
	if values[2] != "" {
		mode.TSysSet = strings.Split(values[2], ",")
	}

	// Parse Interchangeable flag
	if values[3] != "" {
		interchangeable, err := strconv.Atoi(values[3])
		if err != nil {
			return Mode{}, fmt.Errorf("error parsing Interchangeable flag: %w", err)
		}
		mode.Interchangeable = interchangeable
	}

	return mode, nil
}
