package ptvvisum

import (
	"fmt"
	"strconv"
	"strings"
)

// VehCombSection represents $VEHCOMB section
type VehCombSection struct {
	BaseSection
	Combinations []VehicleCombination
}

// VehicleCombination represents a single vehicle combination entry
type VehicleCombination struct {
	No                  int     // Combination number
	Code                string  // Combination code
	VehCombSet          string  // Vehicle combination set
	Name                string  // Combination name/description
	CostRateHourService float64 // Cost rate per hour in service
	CostRateHourEmpty   float64 // Cost rate per hour when empty
	CostRateKmService   float64 // Cost rate per kilometer in service
	CostRateKmEmpty     float64 // Cost rate per kilometer when empty
	CostRateHourLayover float64 // Cost rate per hour during layover
	CostRateHourDepot   float64 // Cost rate per hour at depot
}

// getVehicleCombination extracts data from VEHCOMB section row
func getVehicleCombination(values []string) (VehicleCombination, error) {
	if len(values) < 10 {
		return VehicleCombination{}, fmt.Errorf("invalid VEHCOMB data: %v", values)
	}

	var comb VehicleCombination
	var err error

	// Parse NO (required field)
	if values[0] == "" {
		return VehicleCombination{}, fmt.Errorf("missing required field NO")
	}
	comb.No, err = strconv.Atoi(values[0])
	if err != nil {
		return VehicleCombination{}, fmt.Errorf("error parsing No: %w", err)
	}

	// Set string values
	comb.Code = values[1]
	comb.VehCombSet = values[2]
	comb.Name = values[3]

	// Parse float fields
	floatFields := []struct {
		index int
		dest  *float64
		name  string
	}{
		{4, &comb.CostRateHourService, "CostRateHourService"},
		{5, &comb.CostRateHourEmpty, "CostRateHourEmpty"},
		{6, &comb.CostRateKmService, "CostRateKmService"},
		{7, &comb.CostRateKmEmpty, "CostRateKmEmpty"},
		{8, &comb.CostRateHourLayover, "CostRateHourLayover"},
		{9, &comb.CostRateHourDepot, "CostRateHourDepot"},
	}

	for _, field := range floatFields {
		if values[field.index] != "" {
			*field.dest, err = strconv.ParseFloat(strings.Replace(values[field.index], ",", ".", -1), 64)
			if err != nil {
				return VehicleCombination{}, fmt.Errorf("error parsing %s: %w", field.name, err)
			}
		}
	}

	return comb, nil
}
