package ptvvisum

import (
	"fmt"
	"strconv"
)

// ValidDaysSection represents $VALIDDAYS section
type ValidDaysSection struct {
	BaseSection
	Days []ValidDay
}

// ValidDay represents a single valid day entry
type ValidDay struct {
	No            int
	Code          string
	Name          string
	DayVector     int
	PrfacHourCost float64
	PrfacSupply   float64
}

// getValidDay extracts data from VALIDDAYS section row
func getValidDay(values []string) (ValidDay, error) {
	if len(values) < 6 {
		return ValidDay{}, fmt.Errorf("invalid VALIDDAYS data: %v", values)
	}

	// Parse No (required field)
	no, err := strconv.Atoi(values[0])
	if err != nil {
		return ValidDay{}, fmt.Errorf("error parsing No: %w", err)
	}

	// Initialize with string values
	day := ValidDay{
		No:   no,
		Code: values[1],
		Name: values[2],
	}

	// Parse DayVector (required field)
	dayVector, err := strconv.Atoi(values[3])
	if err != nil {
		return ValidDay{}, fmt.Errorf("error parsing DayVector: %w", err)
	}
	day.DayVector = dayVector

	// Parse PrfacHourCost (required field)
	prfacHourCost, err := strconv.ParseFloat(values[4], 64)
	if err != nil {
		return ValidDay{}, fmt.Errorf("error parsing PrfacHourCost: %w", err)
	}
	day.PrfacHourCost = prfacHourCost

	// Parse PrfacSupply (required field)
	prfacSupply, err := strconv.ParseFloat(values[5], 64)
	if err != nil {
		return ValidDay{}, fmt.Errorf("error parsing PrfacSupply: %w", err)
	}
	day.PrfacSupply = prfacSupply

	return day, nil
}
