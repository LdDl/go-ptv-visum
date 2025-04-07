package ptvvisum

import (
	"fmt"
	"strconv"
	"strings"
)

// BlockItemTypeSection represents $BLOCKITEMTYPE section
type BlockItemTypeSection struct {
	BaseSection
	Types []BlockItemType
}

// BlockItemType represents a single block item type
type BlockItemType struct {
	No                           int     // Item number
	Name                         string  // Item name
	DefLength                    string  // Default length (including unit)
	ShareBefore                  float64 // Share before value
	WeightForLayoversShort       float64 // Weight for short layovers
	WeightForLayoversLong        float64 // Weight for long layovers
	LayoverThresholdShort        string  // Threshold for short layovers (including unit)
	LayoverThresholdLong         string  // Threshold for long layovers (including unit)
	ChargingFunctionInitGradient string  // Initial gradient for charging function
	DischargingFunction          string  // Discharging function specification
}

// getBlockItemType extracts data from BLOCKITEMTYPE section row
func getBlockItemType(values []string) (BlockItemType, error) {
	if len(values) < 2 {
		return BlockItemType{}, fmt.Errorf("invalid BLOCKITEMTYPE data: %v", values)
	}

	// Parse the No field
	no, err := strconv.Atoi(values[0])
	if err != nil {
		return BlockItemType{}, fmt.Errorf("error parsing No: %w", err)
	}

	// Initialize with required fields
	itemType := BlockItemType{
		No:   no,
		Name: values[1],
	}

	// Handle optional fields
	if len(values) > 2 {
		itemType.DefLength = values[2]
	}

	if len(values) > 3 && values[3] != "" {
		shareBefore, err := strconv.ParseFloat(strings.Replace(values[3], ",", ".", -1), 64)
		if err != nil {
			return BlockItemType{}, fmt.Errorf("error parsing ShareBefore: %w", err)
		}
		itemType.ShareBefore = shareBefore
	}

	if len(values) > 4 && values[4] != "" {
		weightShort, err := strconv.ParseFloat(strings.Replace(values[4], ",", ".", -1), 64)
		if err != nil {
			return BlockItemType{}, fmt.Errorf("error parsing WeightForLayoversShort: %w", err)
		}
		itemType.WeightForLayoversShort = weightShort
	}

	if len(values) > 5 && values[5] != "" {
		weightLong, err := strconv.ParseFloat(strings.Replace(values[5], ",", ".", -1), 64)
		if err != nil {
			return BlockItemType{}, fmt.Errorf("error parsing WeightForLayoversLong: %w", err)
		}
		itemType.WeightForLayoversLong = weightLong
	}

	if len(values) > 6 {
		itemType.LayoverThresholdShort = values[6]
	}

	if len(values) > 7 {
		itemType.LayoverThresholdLong = values[7]
	}

	if len(values) > 8 {
		itemType.ChargingFunctionInitGradient = values[8]
	}

	if len(values) > 9 {
		itemType.DischargingFunction = values[9]
	}

	return itemType, nil
}
