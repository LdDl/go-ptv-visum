package ptvvisum

import (
	"fmt"
	"strconv"
	"strings"
)

// LinkTypeSection represents $LINKTYPE section
type LinkTypeSection struct {
	BaseSection
	LinkTypes []LinkType
}

// LinkType represents a single link type in the network
type LinkType struct {
	No                      int                // Link type ID (NO)
	GroupType               int                // Group type (GTYPE)
	Name                    string             // Link type name
	Strict                  int                // Strict flag
	Rank                    int                // Hierarchy rank
	TSysSet                 string             // Set of transport systems allowed
	NumLanes                int                // Number of lanes
	CapPRT                  int                // Capacity for private transport
	V0PRT                   string             // Default speed for private transport
	VMinPRT                 string             // Minimum speed for private transport
	CostRate1PUTSys         map[string]float64 // Cost rate 1 by public transport system
	CostRate2PUTSys         map[string]float64 // Cost rate 2 by public transport system
	CostRate3PUTSys         map[string]float64 // Cost rate 3 by public transport system
	HBEFARoadType           string             // Road type for emissions calculations
	VMaxPRTSys              map[string]string  // Maximum speed by private transport system
	VDefPUTSys              map[string]string  // Default speed by public transport system
	SBAUseOnlyOutermostLane map[string]int     // Lane use restrictions by system
	CapDay                  int                // Daily capacity
	CapHour                 int                // Hourly capacity
	RoadClass               int                // Road classification
}

// GetLinkTypeByID retrieves a link type by its ID
func (s *LinkTypeSection) GetLinkTypeByID(id int) (LinkType, bool) {
	for _, linkType := range s.LinkTypes {
		if linkType.No == id {
			return linkType, true
		}
	}
	return LinkType{}, false
}

// GetLinkTypesByGroupType retrieves all link types of a specified group
func (s *LinkTypeSection) GetLinkTypesByGroupType(groupType int) []LinkType {
	var result []LinkType
	for _, linkType := range s.LinkTypes {
		if linkType.GroupType == groupType {
			result = append(result, linkType)
		}
	}
	return result
}

// GetLinkTypesByRank retrieves all link types with a specified rank
func (s *LinkTypeSection) GetLinkTypesByRank(rank int) []LinkType {
	var result []LinkType
	for _, linkType := range s.LinkTypes {
		if linkType.Rank == rank {
			result = append(result, linkType)
		}
	}
	return result
}

// GetLinkTypesByLanes retrieves all link types with a specified number of lanes
func (s *LinkTypeSection) GetLinkTypesByLanes(numLanes int) []LinkType {
	var result []LinkType
	for _, linkType := range s.LinkTypes {
		if linkType.NumLanes == numLanes {
			result = append(result, linkType)
		}
	}
	return result
}

// GetLinkTypesAllowingTSys retrieves all link types allowing a specific transport system
func (s *LinkTypeSection) GetLinkTypesAllowingTSys(tsys string) []LinkType {
	var result []LinkType
	for _, linkType := range s.LinkTypes {
		systems := strings.Split(linkType.TSysSet, ",")
		for _, system := range systems {
			if system == tsys {
				result = append(result, linkType)
				break
			}
		}
	}
	return result
}

// Count returns the number of link types in the section
func (s *LinkTypeSection) Count() int {
	return len(s.LinkTypes)
}

// getLinkType extracts data from LINKTYPE section row
func getLinkType(values []string, headers []string) (LinkType, error) {
	if len(values) < 8 {
		return LinkType{}, fmt.Errorf("invalid LINKTYPE data (insufficient fields): %v", values)
	}

	var linkType LinkType
	var err error

	// Initialize maps
	linkType.CostRate1PUTSys = make(map[string]float64)
	linkType.CostRate2PUTSys = make(map[string]float64)
	linkType.CostRate3PUTSys = make(map[string]float64)
	linkType.VMaxPRTSys = make(map[string]string)
	linkType.VDefPUTSys = make(map[string]string)
	linkType.SBAUseOnlyOutermostLane = make(map[string]int)

	// Parse NO (required field)
	if values[0] == "" {
		return LinkType{}, fmt.Errorf("missing required field NO")
	}
	linkType.No, err = strconv.Atoi(values[0])
	if err != nil {
		return LinkType{}, fmt.Errorf("error parsing NO: %w", err)
	}

	// Parse GTYPE (required field)
	if values[1] != "" {
		linkType.GroupType, err = strconv.Atoi(values[1])
		if err != nil {
			return LinkType{}, fmt.Errorf("error parsing GTYPE: %w", err)
		}
	}

	// Parse NAME (optional)
	linkType.Name = values[2]

	// Parse STRICT (required field)
	if values[3] != "" {
		linkType.Strict, err = strconv.Atoi(values[3])
		if err != nil {
			return LinkType{}, fmt.Errorf("error parsing STRICT: %w", err)
		}
	}

	// Parse RANK (required field)
	if values[4] != "" {
		linkType.Rank, err = strconv.Atoi(values[4])
		if err != nil {
			return LinkType{}, fmt.Errorf("error parsing RANK: %w", err)
		}
	}

	// Parse TSYSSET (optional)
	linkType.TSysSet = values[5]

	// Parse NUMLANES (required field)
	if values[6] != "" {
		linkType.NumLanes, err = strconv.Atoi(values[6])
		if err != nil {
			return LinkType{}, fmt.Errorf("error parsing NUMLANES: %w", err)
		}
	}

	// Parse CAPPRT (required field)
	if values[7] != "" {
		linkType.CapPRT, err = strconv.Atoi(values[7])
		if err != nil {
			return LinkType{}, fmt.Errorf("error parsing CAPPRT: %w", err)
		}
	}

	// Parse V0PRT (required field)
	if len(values) > 8 {
		linkType.V0PRT = values[8]
	}

	// Parse VMINPRT (optional)
	if len(values) > 9 {
		linkType.VMinPRT = values[9]
	}

	// Parse COSTRATE fields for different PUTSYSs
	for i, header := range headers {
		if i >= len(values) {
			break
		}

		// Process COSTRATE1_PUTSYS fields
		if strings.HasPrefix(header, "COSTRATE1_PUTSYS(") {
			system := extractTransportSystem(header)
			if system != "" && values[i] != "" {
				rate, err := strconv.ParseFloat(values[i], 64)
				if err == nil {
					linkType.CostRate1PUTSys[system] = rate
				}
			}
		}

		// Process COSTRATE2_PUTSYS fields
		if strings.HasPrefix(header, "COSTRATE2_PUTSYS(") {
			system := extractTransportSystem(header)
			if system != "" && values[i] != "" {
				rate, err := strconv.ParseFloat(values[i], 64)
				if err == nil {
					linkType.CostRate2PUTSys[system] = rate
				}
			}
		}

		// Process COSTRATE3_PUTSYS fields
		if strings.HasPrefix(header, "COSTRATE3_PUTSYS(") {
			system := extractTransportSystem(header)
			if system != "" && values[i] != "" {
				rate, err := strconv.ParseFloat(values[i], 64)
				if err == nil {
					linkType.CostRate3PUTSys[system] = rate
				}
			}
		}

		// Process VMAX_PRTSYS fields
		if strings.HasPrefix(header, "VMAX_PRTSYS(") {
			system := extractTransportSystem(header)
			if system != "" && values[i] != "" {
				linkType.VMaxPRTSys[system] = values[i]
			}
		}

		// Process VDEF_PUTSYS fields
		if strings.HasPrefix(header, "VDEF_PUTSYS(") {
			system := extractTransportSystem(header)
			if system != "" && values[i] != "" {
				linkType.VDefPUTSys[system] = values[i]
			}
		}

		// Process SBAUSEONLYOUTERMOSTLANE fields
		if strings.HasPrefix(header, "SBAUSEONLYOUTERMOSTLANE(") {
			system := extractTransportSystem(header)
			if system != "" && values[i] != "" {
				val, err := strconv.Atoi(values[i])
				if err == nil {
					linkType.SBAUseOnlyOutermostLane[system] = val
				}
			}
		}

		// Process HBEFA_ROADTYPE
		if header == "HBEFA_ROADTYPE" && values[i] != "" {
			linkType.HBEFARoadType = values[i]
		}

		// Process CAPDAY
		if header == "CAPDAY" && values[i] != "" {
			linkType.CapDay, _ = strconv.Atoi(values[i])
		}

		// Process CAPHOUR
		if header == "CAPHOUR" && values[i] != "" {
			linkType.CapHour, _ = strconv.Atoi(values[i])
		}

		// Process ROADCLASS
		if header == "ROADCLASS" && values[i] != "" {
			linkType.RoadClass, _ = strconv.Atoi(values[i])
		}
	}

	return linkType, nil
}

// Helper function to extract transport system name from header
func extractTransportSystem(header string) string {
	start := strings.Index(header, "(")
	end := strings.Index(header, ")")
	if start != -1 && end != -1 && end > start {
		return header[start+1 : end]
	}
	return ""
}
