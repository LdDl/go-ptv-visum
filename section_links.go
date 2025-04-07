package ptvvisum

import (
	"fmt"
	"strconv"
	"strings"
)

// LinkSection represents $LINK section
type LinkSection struct {
	BaseSection
	Links []Link
}

// Link represents a single link in the transportation network
type Link struct {
	No                      int                        // Link ID
	FromNodeNo              int                        // Origin node ID
	ToNodeNo                int                        // Destination node ID
	Name                    string                     // Link name (optional)
	TypeNo                  int                        // Link type ID
	TSysSet                 string                     // Transport systems allowed on this link
	UserDirection           int                        // Direction restriction (0=both directions, 1=from→to, 2=to→from)
	Length                  string                     // Length with unit (e.g., "0.081km")
	NumLanes                int                        // Number of lanes
	PlanNo                  int                        // Plan number
	CapPRT                  int                        // Capacity for private transport
	V0PRT                   string                     // Default speed for private transport
	TPuTSys                 map[string]string          // Travel time by public transport system
	TModelSpecial           int                        // Special travel time model flag
	TModelMainNodeSpecial   int                        // Special travel time model for main node
	AddVal                  [3]int                     // Additional values 1-3
	AddValTSys              map[string]int             // Additional values by transport system
	RestrTrafAreaSet        string                     // Restricted traffic area set
	TollPRTSys              map[string]float64         // Toll costs by private transport system
	CostRatePUTSys          map[string]map[int]float64 // Cost rates by public transport system
	NumFarePointsTSys       map[string]int             // Number of fare points by transport system
	FromNodeOrientation     string                     // Orientation at the from-node
	ToNodeOrientation       string                     // Orientation at the to-node
	FromMainNodeOrientation string                     // Orientation at the from main node
	ToMainNodeOrientation   string                     // Orientation at the to main node
	EWSType                 int                        // Environmental sensitivity type
	EWSClass                int                        // Environmental sensitivity class
	SurfaceType             int                        // Road surface type
	NoiseImmisHeight        float64                    // Noise immission height
	ShareHGV                float64                    // Share of heavy goods vehicles
	Slope                   float64                    // Slope (percent)
	ShowBarText             int                        // Show bar text flag
	BarTextRelPos           float64                    // Bar text relative position
	LabelPosRelX            float64                    // Label position X
	LabelPosRelY            float64                    // Label position Y
	SpacePerkPCU            float64                    // Space per PCU (passenger car unit)
	DUEvWave                string                     // Dynamic user equilibrium wave speed
	Urban                   int                        // Urban flag (0=rural, 1=urban)
	SpeedLimit              int                        // Posted speed limit
	Bridge                  int                        // Bridge flag
	Overpass                int                        // Overpass flag
}

// GetLinkByID retrieves a link by its ID
func (s *LinkSection) GetLinkByID(id int) (Link, bool) {
	for _, link := range s.Links {
		if link.No == id {
			return link, true
		}
	}
	return Link{}, false
}

// GetLinksByType retrieves all links of a specified link type
func (s *LinkSection) GetLinksByType(typeNo int) []Link {
	var result []Link
	for _, link := range s.Links {
		if link.TypeNo == typeNo {
			result = append(result, link)
		}
	}
	return result
}

// GetLinksByFromNode retrieves all links originating from a specific node
func (s *LinkSection) GetLinksByFromNode(nodeNo int) []Link {
	var result []Link
	for _, link := range s.Links {
		if link.FromNodeNo == nodeNo {
			result = append(result, link)
		}
	}
	return result
}

// GetLinksByToNode retrieves all links ending at a specific node
func (s *LinkSection) GetLinksByToNode(nodeNo int) []Link {
	var result []Link
	for _, link := range s.Links {
		if link.ToNodeNo == nodeNo {
			result = append(result, link)
		}
	}
	return result
}

// GetLinksBetweenNodes retrieves all links connecting two specific nodes
func (s *LinkSection) GetLinksBetweenNodes(fromNodeNo, toNodeNo int) []Link {
	var result []Link
	for _, link := range s.Links {
		if (link.FromNodeNo == fromNodeNo && link.ToNodeNo == toNodeNo) ||
			(link.FromNodeNo == toNodeNo && link.ToNodeNo == fromNodeNo) {
			result = append(result, link)
		}
	}
	return result
}

// GetLinksByName retrieves all links with a specific name
func (s *LinkSection) GetLinksByName(name string) []Link {
	var result []Link
	for _, link := range s.Links {
		if link.Name == name {
			result = append(result, link)
		}
	}
	return result
}

// GetLinksByTransportSystem retrieves all links allowing a specific transport system
func (s *LinkSection) GetLinksByTransportSystem(tsys string) []Link {
	var result []Link
	for _, link := range s.Links {
		if link.AllowsTransportSystem(tsys) {
			result = append(result, link)
		}
	}
	return result
}

// GetBidirectionalLinks retrieves all links that allow travel in both directions
func (s *LinkSection) GetBidirectionalLinks() []Link {
	var result []Link
	for _, link := range s.Links {
		if link.IsBidirectional() {
			result = append(result, link)
		}
	}
	return result
}

// GetTotalNetworkLength calculates the total length of all links in kilometers
func (s *LinkSection) GetTotalNetworkLength() float64 {
	var totalLength float64
	for _, link := range s.Links {
		totalLength += link.GetLengthInKm()
	}
	return totalLength
}

// GetAverageSpeed calculates the average speed over all links, weighted by length
func (s *LinkSection) GetAverageSpeed() float64 {
	var totalSpeedDistance float64
	var totalDistance float64

	for _, link := range s.Links {
		length := link.GetLengthInKm()
		speed := link.GetSpeedInKmh()

		if length > 0 && speed > 0 {
			totalSpeedDistance += length * speed
			totalDistance += length
		}
	}

	if totalDistance > 0 {
		return totalSpeedDistance / totalDistance
	}
	return 0
}

// Count returns the number of links in the section
func (s *LinkSection) Count() int {
	return len(s.Links)
}

// GetConnectedNodes returns the set of nodes that are connected by this network
func (s *LinkSection) GetConnectedNodes() map[int]bool {
	nodeSet := make(map[int]bool)
	for _, link := range s.Links {
		nodeSet[link.FromNodeNo] = true
		nodeSet[link.ToNodeNo] = true
	}
	return nodeSet
}

// GetAdjacentNodes returns all nodes directly connected to the given node
func (s *LinkSection) GetAdjacentNodes(nodeNo int) []int {
	nodeSet := make(map[int]bool)
	for _, link := range s.Links {
		if link.FromNodeNo == nodeNo {
			nodeSet[link.ToNodeNo] = true
		}
		if link.ToNodeNo == nodeNo {
			nodeSet[link.FromNodeNo] = true
		}
	}

	nodes := make([]int, 0, len(nodeSet))
	for node := range nodeSet {
		nodes = append(nodes, node)
	}
	return nodes
}

// getLink extracts data from LINK section row
func getLink(values []string, headers []string) (Link, error) {
	if len(values) < 12 {
		return Link{}, fmt.Errorf("invalid LINK data (insufficient fields): %v", values)
	}

	var link Link
	var err error

	// Initialize maps
	link.TPuTSys = make(map[string]string)
	link.AddValTSys = make(map[string]int)
	link.TollPRTSys = make(map[string]float64)
	link.CostRatePUTSys = make(map[string]map[int]float64)
	link.NumFarePointsTSys = make(map[string]int)

	// Parse NO (required field)
	if values[0] == "" {
		return Link{}, fmt.Errorf("missing required field NO")
	}
	link.No, err = strconv.Atoi(values[0])
	if err != nil {
		return Link{}, fmt.Errorf("error parsing NO: %w", err)
	}

	// Parse FROMNODENO (required field)
	if values[1] == "" {
		return Link{}, fmt.Errorf("missing required field FROMNODENO")
	}
	link.FromNodeNo, err = strconv.Atoi(values[1])
	if err != nil {
		return Link{}, fmt.Errorf("error parsing FROMNODENO: %w", err)
	}

	// Parse TONODENO (required field)
	if values[2] == "" {
		return Link{}, fmt.Errorf("missing required field TONODENO")
	}
	link.ToNodeNo, err = strconv.Atoi(values[2])
	if err != nil {
		return Link{}, fmt.Errorf("error parsing TONODENO: %w", err)
	}

	// Parse NAME (optional)
	link.Name = values[3]

	// Parse TYPENO (required field)
	if values[4] == "" {
		return Link{}, fmt.Errorf("missing required field TYPENO")
	}
	link.TypeNo, err = strconv.Atoi(values[4])
	if err != nil {
		return Link{}, fmt.Errorf("error parsing TYPENO: %w", err)
	}

	// Parse TSYSSET (optional but usually present)
	link.TSysSet = values[5]

	// Parse USERDIRECTION (optional)
	if values[6] != "" {
		link.UserDirection, err = strconv.Atoi(values[6])
		if err != nil {
			return Link{}, fmt.Errorf("error parsing USERDIRECTION: %w", err)
		}
	}

	// Parse LENGTH (required field)
	link.Length = values[7]

	// Parse NUMLANES (required field)
	if values[8] != "" {
		link.NumLanes, err = strconv.Atoi(values[8])
		if err != nil {
			return Link{}, fmt.Errorf("error parsing NUMLANES: %w", err)
		}
	}

	// Parse PLANNO (optional)
	if values[9] != "" {
		link.PlanNo, err = strconv.Atoi(values[9])
		if err != nil {
			return Link{}, fmt.Errorf("error parsing PLANNO: %w", err)
		}
	}

	// Parse CAPPRT (required field)
	if values[10] != "" {
		link.CapPRT, err = strconv.Atoi(values[10])
		if err != nil {
			return Link{}, fmt.Errorf("error parsing CAPPRT: %w", err)
		}
	}

	// Parse V0PRT (required field)
	link.V0PRT = values[11]

	// Process remaining fields based on headers
	for i := 0; i < len(headers) && i < len(values); i++ {
		headerName := headers[i]
		value := values[i]

		if value == "" {
			continue // Skip empty values
		}

		// Process T_PUTSYS fields
		if strings.HasPrefix(headerName, "T_PUTSYS(") {
			tsys := extractSystemName(headerName)
			if tsys != "" {
				link.TPuTSys[tsys] = value
			}
		}

		// Process TMODELSPECIAL
		if headerName == "TMODELSPECIAL" {
			link.TModelSpecial, _ = strconv.Atoi(value)
		}

		// Process TMODELMAINNODESPECIAL
		if headerName == "TMODELMAINNODESPECIAL" {
			link.TModelMainNodeSpecial, _ = strconv.Atoi(value)
		}

		// Process ADDVAL1, ADDVAL2, ADDVAL3
		if headerName == "ADDVAL1" {
			link.AddVal[0], _ = strconv.Atoi(value)
		}
		if headerName == "ADDVAL2" {
			link.AddVal[1], _ = strconv.Atoi(value)
		}
		if headerName == "ADDVAL3" {
			link.AddVal[2], _ = strconv.Atoi(value)
		}

		// Process ADDVAL_TSYS fields
		if strings.HasPrefix(headerName, "ADDVAL_TSYS(") {
			tsys := extractSystemName(headerName)
			if tsys != "" {
				val, err := strconv.Atoi(value)
				if err == nil {
					link.AddValTSys[tsys] = val
				}
			}
		}

		// Process RESTRTRAFAREASET
		if headerName == "RESTRTRAFAREASET" {
			link.RestrTrafAreaSet = value
		}

		// Process TOLL_PRTSYS fields
		if strings.HasPrefix(headerName, "TOLL_PRTSYS(") {
			tsys := extractSystemName(headerName)
			if tsys != "" {
				val, err := strconv.ParseFloat(value, 64)
				if err == nil {
					link.TollPRTSys[tsys] = val
				}
			}
		}

		// Process COSTRATE fields for different PUTSYSs
		if strings.HasPrefix(headerName, "COSTRATE") && strings.Contains(headerName, "PUTSYS") {
			// Extract system and rate number (1, 2, or 3)
			parts := strings.Split(headerName, "_")
			if len(parts) >= 2 {
				rateStr := strings.TrimPrefix(parts[0], "COSTRATE")
				rateNum, err := strconv.Atoi(rateStr)
				if err == nil && parts[1] != "" {
					tsys := extractSystemName(parts[1])
					if tsys != "" {
						rate, err := strconv.ParseFloat(value, 64)
						if err == nil {
							if link.CostRatePUTSys[tsys] == nil {
								link.CostRatePUTSys[tsys] = make(map[int]float64)
							}
							link.CostRatePUTSys[tsys][rateNum] = rate
						}
					}
				}
			}
		}

		// Process NUMFAREPOINTS_TSYS fields
		if strings.HasPrefix(headerName, "NUMFAREPOINTS_TSYS(") {
			tsys := extractSystemName(headerName)
			if tsys != "" {
				val, err := strconv.Atoi(value)
				if err == nil {
					link.NumFarePointsTSys[tsys] = val
				}
			}
		}

		// Process orientation fields
		if headerName == "FROMNODEORIENTATION" {
			link.FromNodeOrientation = value
		}
		if headerName == "TONODEORIENTATION" {
			link.ToNodeOrientation = value
		}
		if headerName == "FROMMAINNODEORIENTATION" {
			link.FromMainNodeOrientation = value
		}
		if headerName == "TOMAINNODEORIENTATION" {
			link.ToMainNodeOrientation = value
		}

		// Process EWSTYPE
		if headerName == "EWSTYPE" {
			link.EWSType, _ = strconv.Atoi(value)
		}

		// Process EWSCLASS
		if headerName == "EWSCLASS" {
			link.EWSClass, _ = strconv.Atoi(value)
		}

		// Process SURFACETYPE
		if headerName == "SURFACETYPE" {
			link.SurfaceType, _ = strconv.Atoi(value)
		}

		// Process NOISEIMMISHEIGHT
		if headerName == "NOISEIMMISHEIGHT" {
			link.NoiseImmisHeight, _ = strconv.ParseFloat(strings.Replace(value, "m", "", 1), 64)
		}

		// Process SHAREHGV
		if headerName == "SHAREHGV" {
			link.ShareHGV, _ = strconv.ParseFloat(value, 64)
		}

		// Process SLOPE
		if headerName == "SLOPE" {
			link.Slope, _ = strconv.ParseFloat(value, 64)
		}

		// Process SHOWBARTEXT
		if headerName == "SHOWBARTEXT" {
			link.ShowBarText, _ = strconv.Atoi(value)
		}

		// Process BARTEXTRELPOS
		if headerName == "BARTEXTRELPOS" {
			link.BarTextRelPos, _ = strconv.ParseFloat(value, 64)
		}

		// Process LABELPOSRELX
		if headerName == "LABELPOSRELX" {
			link.LabelPosRelX, _ = strconv.ParseFloat(value, 64)
		}

		// Process LABELPOSRELY
		if headerName == "LABELPOSRELY" {
			link.LabelPosRelY, _ = strconv.ParseFloat(value, 64)
		}

		// Process SPACEPERPCU
		if headerName == "SPACEPERPCU" {
			link.SpacePerkPCU, _ = strconv.ParseFloat(value, 64)
		}

		// Process DUEVWAVE
		if headerName == "DUEVWAVE" {
			link.DUEvWave = value
		}

		// Process URBAN
		if headerName == "URBAN" {
			link.Urban, _ = strconv.Atoi(value)
		}

		// Process SPEEDLIMIT
		if headerName == "SPEEDLIMIT" {
			link.SpeedLimit, _ = strconv.Atoi(value)
		}

		// Process BRIDGE
		if headerName == "BRIDGE" {
			link.Bridge, _ = strconv.Atoi(value)
		}

		// Process OVERPASS
		if headerName == "OVERPASS" {
			link.Overpass, _ = strconv.Atoi(value)
		}
	}

	return link, nil
}

// Helper function to extract transport system name from header
func extractSystemName(header string) string {
	start := strings.Index(header, "(")
	end := strings.Index(header, ")")
	if start != -1 && end != -1 && end > start {
		return header[start+1 : end]
	}
	return ""
}

// GetLengthInKm parses the Length string and returns the value in kilometers
func (l *Link) GetLengthInKm() float64 {
	if l.Length == "" {
		return 0
	}

	// Extract numeric part
	numStr := l.Length
	unit := ""

	// Find where the numeric part ends
	for i, c := range l.Length {
		if !isDigit(c) && c != '.' && c != ',' {
			numStr = l.Length[:i]
			unit = strings.ToLower(l.Length[i:])
			break
		}
	}

	// Convert to float
	numStr = strings.Replace(numStr, ",", ".", 1)
	value, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0
	}

	// Convert to km based on unit
	switch {
	case strings.Contains(unit, "km"):
		return value
	case strings.Contains(unit, "m"):
		return value / 1000
	default:
		return value // Assume km if no unit specified
	}
}

// Helper function to check if a rune is a digit
func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

// GetSpeedInKmh parses the V0PRT string and returns the value in km/h
func (l *Link) GetSpeedInKmh() float64 {
	if l.V0PRT == "" {
		return 0
	}

	// Extract numeric part
	numStr := l.V0PRT

	// Find where the numeric part ends
	for i, c := range l.V0PRT {
		if !isDigit(c) && c != '.' && c != ',' {
			numStr = l.V0PRT[:i]
			break
		}
	}

	// Convert to float
	numStr = strings.Replace(numStr, ",", ".", 1)
	value, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0
	}

	return value
}

// IsBidirectional checks if the link allows travel in both directions
func (l *Link) IsBidirectional() bool {
	return l.UserDirection == 0
}

// AllowsTransportSystem checks if the link allows the specified transport system
func (l *Link) AllowsTransportSystem(tsys string) bool {
	systems := strings.Split(l.TSysSet, ",")
	for _, system := range systems {
		if system == tsys {
			return true
		}
	}
	return false
}
