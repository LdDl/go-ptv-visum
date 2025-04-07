package ptvvisum

import (
	"fmt"
	"strconv"
	"strings"
)

// VehUnitSection represents $VEHUNIT section
type VehUnitSection struct {
	BaseSection
	Units []VehicleUnit
}

// VehicleUnit represents a single vehicle unit entry
type VehicleUnit struct {
	No                  int     // Unit number
	Code                string  // Unit code
	Name                string  // Unit name/description
	TSysSet             string  // Transport system set
	Powered             int     // Is powered flag (0/1)
	SeatCap             int     // Seating capacity
	TotalCap            int     // Total capacity (seated + standing)
	CostRateHourService float64 // Cost rate per hour in service
	CostRateHourEmpty   float64 // Cost rate per hour when empty
	CostRateHourLayover float64 // Cost rate per hour during layover
	CostRateHourDepot   float64 // Cost rate per hour at depot
	CostRateKmService   float64 // Cost rate per kilometer in service
	CostRateKmEmpty     float64 // Cost rate per kilometer when empty
	CostRateVehUnit     float64 // Fixed cost rate per vehicle unit
}

// getVehicleUnit extracts data from VEHUNIT section row
func getVehicleUnit(values []string) (VehicleUnit, error) {
	if len(values) < 14 {
		return VehicleUnit{}, fmt.Errorf("invalid VEHUNIT data: %v", values)
	}

	var unit VehicleUnit
	var err error

	// Parse NO (required field)
	if values[0] == "" {
		return VehicleUnit{}, fmt.Errorf("missing required field NO")
	}
	unit.No, err = strconv.Atoi(values[0])
	if err != nil {
		return VehicleUnit{}, fmt.Errorf("error parsing No: %w", err)
	}

	// Set string values
	unit.Code = values[1]
	unit.Name = values[2]
	unit.TSysSet = values[3]

	// Parse integer fields
	intFields := []struct {
		index int
		dest  *int
		name  string
	}{
		{4, &unit.Powered, "Powered"},
		{5, &unit.SeatCap, "SeatCap"},
		{6, &unit.TotalCap, "TotalCap"},
	}

	for _, field := range intFields {
		if values[field.index] != "" {
			*field.dest, err = strconv.Atoi(values[field.index])
			if err != nil {
				return VehicleUnit{}, fmt.Errorf("error parsing %s: %w", field.name, err)
			}
		}
	}

	// Parse float fields
	floatFields := []struct {
		index int
		dest  *float64
		name  string
	}{
		{7, &unit.CostRateHourService, "CostRateHourService"},
		{8, &unit.CostRateHourEmpty, "CostRateHourEmpty"},
		{9, &unit.CostRateHourLayover, "CostRateHourLayover"},
		{10, &unit.CostRateHourDepot, "CostRateHourDepot"},
		{11, &unit.CostRateKmService, "CostRateKmService"},
		{12, &unit.CostRateKmEmpty, "CostRateKmEmpty"},
		{13, &unit.CostRateVehUnit, "CostRateVehUnit"},
	}

	for _, field := range floatFields {
		if values[field.index] != "" {
			*field.dest, err = strconv.ParseFloat(strings.Replace(values[field.index], ",", ".", -1), 64)
			if err != nil {
				return VehicleUnit{}, fmt.Errorf("error parsing %s: %w", field.name, err)
			}
		}
	}

	return unit, nil
}
