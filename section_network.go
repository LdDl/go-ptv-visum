package ptvvisum

import (
	"fmt"
	"strconv"
	"strings"
)

// NetworkSection represents the $NETWORK section
type NetworkSection struct {
	BaseSection
	Network NetworkData
}

// NetworkData represents network configuration parameters
type NetworkData struct {
	NetVersionID                    string  // Version ID
	NetVersionName                  string  // Version name
	Scale                           float64 // Scale factor
	Unit                            string  // Length unit (KM, M, etc.)
	LeftHandTraffic                 int     // 1 if left-hand traffic, 0 otherwise
	CoordDecPlaces                  int     // Decimal places for coordinates
	DecPlacesOther                  int     // Decimal places for other values
	CurrencyDecPlaces               int     // Decimal places for currency values
	LongLengthDecPlaces             int     // Decimal places for long lengths
	ShortLengthDecPlaces            int     // Decimal places for short lengths
	TurnT0DecPlaces                 int     // Decimal places for turn times
	SpeedDecPlaces                  int     // Decimal places for speeds
	MaxFloatPrecisionFileExport     int     // Max precision for float export
	ConcatMaxLen                    int     // Max length for concatenation
	ConcatSeparator                 string  // Separator for concatenation
	CreateModedSeg                  int     // Create moded segments flag
	ProjectionDefinition            string  // GIS projection definition
	TurnTypeDefault                 string  // Default turn type
	LinkOrientationCalculationType  string  // Link orientation calculation type
	TransferWaitTimeLimitForReached string  // Transfer wait time limit (reached)
	TransferWaitTimeLimitForMissed  string  // Transfer wait time limit (missed)
	TransfersOnlyDifferentLines     int     // Transfers only between different lines flag
	StrongLineRouteLengthsAdaption  int     // Strong line route lengths adaption flag
	Name                            string  // Network name
}

// getNetwork extracts data from the NETWORK section row
func getNetwork(values []string) (NetworkData, error) {
	if len(values) < 14 {
		return NetworkData{}, fmt.Errorf("insufficient NETWORK data (need at least 14 fields): %v", values)
	}

	network := NetworkData{
		NetVersionID:   values[0],
		NetVersionName: values[1],
		Unit:           values[3],
	}

	// Parse numeric values
	var err error

	// Scale
	if values[2] != "" {
		network.Scale, err = strconv.ParseFloat(strings.Replace(values[2], ",", ".", -1), 64)
		if err != nil {
			return NetworkData{}, fmt.Errorf("error parsing Scale: %w", err)
		}
	}

	// Integer values
	intFields := []struct {
		index int
		dest  *int
		name  string
	}{
		{4, &network.LeftHandTraffic, "LeftHandTraffic"},
		{5, &network.CoordDecPlaces, "CoordDecPlaces"},
		{6, &network.DecPlacesOther, "DecPlacesOther"},
		{7, &network.CurrencyDecPlaces, "CurrencyDecPlaces"},
		{8, &network.LongLengthDecPlaces, "LongLengthDecPlaces"},
		{9, &network.ShortLengthDecPlaces, "ShortLengthDecPlaces"},
		{10, &network.TurnT0DecPlaces, "TurnT0DecPlaces"},
		{11, &network.SpeedDecPlaces, "SpeedDecPlaces"},
		{12, &network.MaxFloatPrecisionFileExport, "MaxFloatPrecisionFileExport"},
		{13, &network.ConcatMaxLen, "ConcatMaxLen"},
	}

	for _, field := range intFields {
		if values[field.index] != "" {
			*field.dest, err = strconv.Atoi(values[field.index])
			if err != nil {
				return NetworkData{}, fmt.Errorf("error parsing %s: %w", field.name, err)
			}
		}
	}

	// Optional fields
	if len(values) > 14 {
		network.ConcatSeparator = values[14]
	}

	if len(values) > 15 && values[15] != "" {
		network.CreateModedSeg, err = strconv.Atoi(values[15])
		if err != nil {
			return NetworkData{}, fmt.Errorf("error parsing CreateModedSeg: %w", err)
		}
	}

	if len(values) > 16 {
		network.ProjectionDefinition = values[16]
	}

	if len(values) > 17 {
		network.TurnTypeDefault = values[17]
	}

	if len(values) > 18 {
		network.LinkOrientationCalculationType = values[18]
	}

	if len(values) > 19 {
		network.TransferWaitTimeLimitForReached = values[19]
	}

	if len(values) > 20 {
		network.TransferWaitTimeLimitForMissed = values[20]
	}

	if len(values) > 21 && values[21] != "" {
		network.TransfersOnlyDifferentLines, err = strconv.Atoi(values[21])
		if err != nil {
			return NetworkData{}, fmt.Errorf("error parsing TransfersOnlyDifferentLines: %w", err)
		}
	}

	if len(values) > 22 && values[22] != "" {
		network.StrongLineRouteLengthsAdaption, err = strconv.Atoi(values[22])
		if err != nil {
			return NetworkData{}, fmt.Errorf("error parsing StrongLineRouteLengthsAdaption: %w", err)
		}
	}

	if len(values) > 23 {
		network.Name = values[23]
	}

	return network, nil
}
