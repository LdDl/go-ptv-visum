package ptvvisum

import (
	"fmt"
	"strconv"
	"strings"
)

// DemandSegmentSection represents $DEMANDSEGMENT section
type DemandSegmentSection struct {
	BaseSection
	Segments []DemandSegment
}

// DemandSegment represents a single demand segment entry
type DemandSegment struct {
	Code          string  // Segment code
	Name          string  // Segment name
	Mode          string  // Transport mode
	OccupancyRate float64 // Occupancy rate
	PrFacAP       float64 // Price factor active days
	PrFacAH       float64 // Price factor active hours
}

// getDemandSegment extracts data from DEMANDSEGMENT section row
func getDemandSegment(values []string) (DemandSegment, error) {
	if len(values) < 6 {
		return DemandSegment{}, fmt.Errorf("invalid DEMANDSEGMENT data: %v", values)
	}

	// Initialize with string values
	segment := DemandSegment{
		Code: values[0],
		Name: values[1],
		Mode: values[2],
	}

	var err error

	// Parse OccupancyRate
	if values[3] != "" {
		segment.OccupancyRate, err = strconv.ParseFloat(strings.Replace(values[3], ",", ".", -1), 64)
		if err != nil {
			return DemandSegment{}, fmt.Errorf("error parsing OccupancyRate: %w", err)
		}
	}

	// Parse PrFacAP
	if values[4] != "" {
		segment.PrFacAP, err = strconv.ParseFloat(strings.Replace(values[4], ",", ".", -1), 64)
		if err != nil {
			return DemandSegment{}, fmt.Errorf("error parsing PrFacAP: %w", err)
		}
	}

	// Parse PrFacAH
	if values[5] != "" {
		segment.PrFacAH, err = strconv.ParseFloat(strings.Replace(values[5], ",", ".", -1), 64)
		if err != nil {
			return DemandSegment{}, fmt.Errorf("error parsing PrFacAH: %w", err)
		}
	}

	return segment, nil
}
