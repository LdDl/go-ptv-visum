package ptvvisum

import (
	"fmt"
	"strconv"
	"strings"
)

// TSysSection represents $TSYS section (Transport Systems)
type TSysSection struct {
	BaseSection
	Systems []TransportSystem
}

// TransportSystem represents a single transport system entry
type TransportSystem struct {
	Code                     string  // Transport system code
	Name                     string  // Transport system name
	Type                     string  // System type (PuT, PrT, etc.)
	PCU                      float64 // Passenger car unit
	SBAReactionTime          string  // SBA reaction time
	SBAEffVehLength          string  // SBA effective vehicle length
	SBAMaxWaitingTime        string  // SBA max waiting time
	IsRoundTripSystem        int     // Is round trip system flag
	IsStationBased           int     // Is station based flag
	AllowRelocations         int     // Allow relocations flag
	MaxNumRelocationsPerHour float64 // Max number of relocations per hour
	HasDepot                 int     // Has depot flag
	NumVehiclesInNetwork     float64 // Number of vehicles in network
	OccupancyRate            float64 // Occupancy rate
}

// getTransportSystem extracts data from TSYS section row
func getTransportSystem(values []string) (TransportSystem, error) {
	if len(values) < 3 {
		return TransportSystem{}, fmt.Errorf("insufficient TSYS data (need at least CODE, NAME, TYPE): %v", values)
	}

	// Always initialize these fields
	ts := TransportSystem{
		Code: values[0],
		Name: values[1],
		Type: values[2],
	}

	// Parse PCU (Passenger Car Unit) if available
	if len(values) > 3 && values[3] != "" {
		pcu, err := strconv.ParseFloat(strings.Replace(values[3], ",", ".", -1), 64)
		if err != nil {
			return TransportSystem{}, fmt.Errorf("error parsing PCU: %w", err)
		}
		ts.PCU = pcu
	}

	// SBA fields (as strings with units)
	if len(values) > 4 {
		ts.SBAReactionTime = values[4]
	}

	if len(values) > 5 {
		ts.SBAEffVehLength = values[5]
	}

	if len(values) > 6 {
		ts.SBAMaxWaitingTime = values[6]
	}

	// Boolean flags (stored as integers)
	if len(values) > 7 && values[7] != "" {
		isRoundTrip, err := strconv.Atoi(values[7])
		if err != nil {
			return TransportSystem{}, fmt.Errorf("error parsing IsRoundTripSystem: %w", err)
		}
		ts.IsRoundTripSystem = isRoundTrip
	}

	if len(values) > 8 && values[8] != "" {
		isStation, err := strconv.Atoi(values[8])
		if err != nil {
			return TransportSystem{}, fmt.Errorf("error parsing IsStationBased: %w", err)
		}
		ts.IsStationBased = isStation
	}

	if len(values) > 9 && values[9] != "" {
		allowRel, err := strconv.Atoi(values[9])
		if err != nil {
			return TransportSystem{}, fmt.Errorf("error parsing AllowRelocations: %w", err)
		}
		ts.AllowRelocations = allowRel
	}

	// Numeric fields
	if len(values) > 10 && values[10] != "" {
		maxRel, err := strconv.ParseFloat(strings.Replace(values[10], ",", ".", -1), 64)
		if err != nil {
			return TransportSystem{}, fmt.Errorf("error parsing MaxNumRelocationsPerHour: %w", err)
		}
		ts.MaxNumRelocationsPerHour = maxRel
	}

	if len(values) > 11 && values[11] != "" {
		hasDepot, err := strconv.Atoi(values[11])
		if err != nil {
			return TransportSystem{}, fmt.Errorf("error parsing HasDepot: %w", err)
		}
		ts.HasDepot = hasDepot
	}

	if len(values) > 12 && values[12] != "" {
		numVeh, err := strconv.ParseFloat(strings.Replace(values[12], ",", ".", -1), 64)
		if err != nil {
			return TransportSystem{}, fmt.Errorf("error parsing NumVehiclesInNetwork: %w", err)
		}
		ts.NumVehiclesInNetwork = numVeh
	}

	if len(values) > 13 && values[13] != "" {
		occRate, err := strconv.ParseFloat(strings.Replace(values[13], ",", ".", -1), 64)
		if err != nil {
			return TransportSystem{}, fmt.Errorf("error parsing OccupancyRate: %w", err)
		}
		ts.OccupancyRate = occRate
	}

	return ts, nil
}
