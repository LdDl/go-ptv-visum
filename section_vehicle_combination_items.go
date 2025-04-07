package ptvvisum

import (
	"fmt"
	"strconv"
)

// VehUnitToVehCombSection represents $VEHUNITTOVEHCOMB section
type VehUnitToVehCombSection struct {
	BaseSection
	Mappings []VehUnitToVehCombMapping
}

// VehUnitToVehCombMapping represents a mapping between vehicle combinations and units
type VehUnitToVehCombMapping struct {
	VehCombNo   int // Vehicle combination number
	VehUnitNo   int // Vehicle unit number
	NumVehUnits int // Number of vehicle units of this type in the combination
}

// getVehUnitToVehCombMapping extracts data from VEHUNITTOVEHCOMB section row
func getVehUnitToVehCombMapping(values []string) (VehUnitToVehCombMapping, error) {
	if len(values) < 3 {
		return VehUnitToVehCombMapping{}, fmt.Errorf("invalid VEHUNITTOVEHCOMB data: %v", values)
	}

	var mapping VehUnitToVehCombMapping
	var err error

	// Parse VehCombNo (required field)
	if values[0] == "" {
		return VehUnitToVehCombMapping{}, fmt.Errorf("missing required field VEHCOMBNO")
	}
	mapping.VehCombNo, err = strconv.Atoi(values[0])
	if err != nil {
		return VehUnitToVehCombMapping{}, fmt.Errorf("error parsing VehCombNo: %w", err)
	}

	// Parse VehUnitNo (required field)
	if values[1] == "" {
		return VehUnitToVehCombMapping{}, fmt.Errorf("missing required field VEHUNITNO")
	}
	mapping.VehUnitNo, err = strconv.Atoi(values[1])
	if err != nil {
		return VehUnitToVehCombMapping{}, fmt.Errorf("error parsing VehUnitNo: %w", err)
	}

	// Parse NumVehUnits (required field)
	if values[2] == "" {
		return VehUnitToVehCombMapping{}, fmt.Errorf("missing required field NUMVEHUNITS")
	}
	mapping.NumVehUnits, err = strconv.Atoi(values[2])
	if err != nil {
		return VehUnitToVehCombMapping{}, fmt.Errorf("error parsing NumVehUnits: %w", err)
	}

	return mapping, nil
}
