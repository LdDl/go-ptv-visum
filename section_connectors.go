package ptvvisum

import (
	"fmt"
	"strconv"
	"strings"
)

// ConnectorSection represents $CONNECTOR section
type ConnectorSection struct {
	BaseSection
	Connectors []Connector
}

// Connector represents a single connector in the transportation network
type Connector struct {
	ZoneNo       int               // Zone ID
	NodeNo       int               // Node ID
	Direction    string            // Direction (O=Origin, D=Destination)
	TypeNo       int               // Connector type ID
	TSysSet      string            // Transport systems allowed
	Length       string            // Length with unit (e.g., "0.903km")
	T0TSys       map[string]string // Travel time by transport system
	WeightPRT    float64           // Weight for private transport
	WeightPUT    float64           // Weight for public transport
	AddVal       [3]int            // Additional values 1-3
	LabelPosRelX float64           // X coordinate for label
	LabelPosRelY float64           // Y coordinate for label
}

// GetConnectorsByZone retrieves all connectors for a specific zone
func (s *ConnectorSection) GetConnectorsByZone(zoneNo int) []Connector {
	var result []Connector
	for _, connector := range s.Connectors {
		if connector.ZoneNo == zoneNo {
			result = append(result, connector)
		}
	}
	return result
}

// GetConnectorsByNode retrieves all connectors for a specific node
func (s *ConnectorSection) GetConnectorsByNode(nodeNo int) []Connector {
	var result []Connector
	for _, connector := range s.Connectors {
		if connector.NodeNo == nodeNo {
			result = append(result, connector)
		}
	}
	return result
}

// GetOriginConnectors retrieves all origin connectors
func (s *ConnectorSection) GetOriginConnectors() []Connector {
	var result []Connector
	for _, connector := range s.Connectors {
		if connector.IsOriginConnector() {
			result = append(result, connector)
		}
	}
	return result
}

// GetDestinationConnectors retrieves all destination connectors
func (s *ConnectorSection) GetDestinationConnectors() []Connector {
	var result []Connector
	for _, connector := range s.Connectors {
		if connector.IsDestinationConnector() {
			result = append(result, connector)
		}
	}
	return result
}

// GetConnectorsByType retrieves all connectors of a specified type
func (s *ConnectorSection) GetConnectorsByType(typeNo int) []Connector {
	var result []Connector
	for _, connector := range s.Connectors {
		if connector.TypeNo == typeNo {
			result = append(result, connector)
		}
	}
	return result
}

// GetConnectorsByTransportSystem retrieves all connectors allowing a specific transport system
func (s *ConnectorSection) GetConnectorsByTransportSystem(tsys string) []Connector {
	var result []Connector
	for _, connector := range s.Connectors {
		if connector.AllowsTransportSystem(tsys) {
			result = append(result, connector)
		}
	}
	return result
}

// GetConnector retrieves a specific connector by zone, node and direction
func (s *ConnectorSection) GetConnector(zoneNo, nodeNo int, direction string) (Connector, bool) {
	for _, connector := range s.Connectors {
		if connector.ZoneNo == zoneNo && connector.NodeNo == nodeNo && connector.Direction == direction {
			return connector, true
		}
	}
	return Connector{}, false
}

// GetTotalConnectorLength calculates the total length of all connectors in kilometers
func (s *ConnectorSection) GetTotalConnectorLength() float64 {
	var totalLength float64
	for _, connector := range s.Connectors {
		totalLength += connector.GetLengthInKm()
	}
	return totalLength
}

// GetAverageTravelTime calculates the average travel time over all connectors for a specific mode
func (s *ConnectorSection) GetAverageTravelTime(tsys string) float64 {
	var totalTime float64
	count := 0

	for _, connector := range s.Connectors {
		time := connector.GetTravelTimeSeconds(tsys)
		if time > 0 {
			totalTime += time
			count++
		}
	}

	if count > 0 {
		return totalTime / float64(count)
	}
	return 0
}

// Count returns the number of connectors in the section
func (s *ConnectorSection) Count() int {
	return len(s.Connectors)
}

// CountOriginConnectors returns the number of origin connectors
func (s *ConnectorSection) CountOriginConnectors() int {
	count := 0
	for _, connector := range s.Connectors {
		if connector.IsOriginConnector() {
			count++
		}
	}
	return count
}

// CountDestinationConnectors returns the number of destination connectors
func (s *ConnectorSection) CountDestinationConnectors() int {
	count := 0
	for _, connector := range s.Connectors {
		if connector.IsDestinationConnector() {
			count++
		}
	}
	return count
}

// GetZoneConnectivity returns a map showing how many connectors each zone has
func (s *ConnectorSection) GetZoneConnectivity() map[int]int {
	connectivity := make(map[int]int)
	for _, connector := range s.Connectors {
		connectivity[connector.ZoneNo]++
	}
	return connectivity
}

// GetNodeConnectivity returns a map showing how many connectors each node has
func (s *ConnectorSection) GetNodeConnectivity() map[int]int {
	connectivity := make(map[int]int)
	for _, connector := range s.Connectors {
		connectivity[connector.NodeNo]++
	}
	return connectivity
}

// getConnector extracts data from CONNECTOR section row
func getConnector(values []string, headers []string) (Connector, error) {
	if len(values) < 13 {
		return Connector{}, fmt.Errorf("invalid CONNECTOR data (insufficient fields): %v", values)
	}

	var connector Connector
	var err error

	// Initialize maps
	connector.T0TSys = make(map[string]string)

	// Parse ZONENO (required field)
	if values[0] == "" {
		return Connector{}, fmt.Errorf("missing required field ZONENO")
	}
	connector.ZoneNo, err = strconv.Atoi(values[0])
	if err != nil {
		return Connector{}, fmt.Errorf("error parsing ZONENO: %w", err)
	}

	// Parse NODENO (required field)
	if values[1] == "" {
		return Connector{}, fmt.Errorf("missing required field NODENO")
	}
	connector.NodeNo, err = strconv.Atoi(values[1])
	if err != nil {
		return Connector{}, fmt.Errorf("error parsing NODENO: %w", err)
	}

	// Parse DIRECTION (required field)
	connector.Direction = values[2]
	if connector.Direction != "O" && connector.Direction != "D" {
		return Connector{}, fmt.Errorf("invalid DIRECTION value: %s (should be O or D)", connector.Direction)
	}

	// Parse TYPENO (required field)
	if values[3] != "" {
		connector.TypeNo, err = strconv.Atoi(values[3])
		if err != nil {
			return Connector{}, fmt.Errorf("error parsing TYPENO: %w", err)
		}
	}

	// Parse TSYSSET (optional)
	connector.TSysSet = values[4]

	// Parse LENGTH (required field)
	connector.Length = values[5]

	// Process T0_TSYS fields based on headers
	for i := 0; i < len(headers) && i < len(values); i++ {
		headerName := headers[i]
		value := values[i]

		if value == "" {
			continue // Skip empty values
		}

		// Process T0_TSYS fields
		if strings.HasPrefix(headerName, "T0_TSYS(") {
			tsys := extractSystemName(headerName)
			if tsys != "" {
				connector.T0TSys[tsys] = value
			}
		}
	}

	// Parse WEIGHT(PRT) (required field)
	if values[11] != "" {
		connector.WeightPRT, err = strconv.ParseFloat(strings.Replace(values[11], ",", ".", -1), 64)
		if err != nil {
			return Connector{}, fmt.Errorf("error parsing WEIGHT(PRT): %w", err)
		}
	}

	// Parse WEIGHT(PUT) (required field)
	if values[12] != "" {
		connector.WeightPUT, err = strconv.ParseFloat(strings.Replace(values[12], ",", ".", -1), 64)
		if err != nil {
			return Connector{}, fmt.Errorf("error parsing WEIGHT(PUT): %w", err)
		}
	}

	// Parse ADDVAL1, ADDVAL2, ADDVAL3 (if available)
	if len(values) > 13 && values[13] != "" {
		connector.AddVal[0], err = strconv.Atoi(values[13])
		if err != nil {
			return Connector{}, fmt.Errorf("error parsing ADDVAL1: %w", err)
		}
	}

	if len(values) > 14 && values[14] != "" {
		connector.AddVal[1], err = strconv.Atoi(values[14])
		if err != nil {
			return Connector{}, fmt.Errorf("error parsing ADDVAL2: %w", err)
		}
	}

	if len(values) > 15 && values[15] != "" {
		connector.AddVal[2], err = strconv.Atoi(values[15])
		if err != nil {
			return Connector{}, fmt.Errorf("error parsing ADDVAL3: %w", err)
		}
	}

	// Parse LABELPOSRELX (optional)
	if len(values) > 16 && values[16] != "" {
		connector.LabelPosRelX, err = strconv.ParseFloat(strings.Replace(values[16], ",", ".", -1), 64)
		if err != nil {
			return Connector{}, fmt.Errorf("error parsing LABELPOSRELX: %w", err)
		}
	}

	// Parse LABELPOSRELY (optional)
	if len(values) > 17 && values[17] != "" {
		connector.LabelPosRelY, err = strconv.ParseFloat(strings.Replace(values[17], ",", ".", -1), 64)
		if err != nil {
			return Connector{}, fmt.Errorf("error parsing LABELPOSRELY: %w", err)
		}
	}

	return connector, nil
}

// GetLengthInKm parses the Length string and returns the value in kilometers
func (c *Connector) GetLengthInKm() float64 {
	if c.Length == "" {
		return 0
	}

	// Extract numeric part
	numStr := c.Length
	unit := ""

	// Find where the numeric part ends
	for i, ch := range c.Length {
		if !isDigit(ch) && ch != '.' && ch != ',' {
			numStr = c.Length[:i]
			unit = strings.ToLower(c.Length[i:])
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

// GetTravelTimeSeconds returns the travel time in seconds for a specific transport system
func (c *Connector) GetTravelTimeSeconds(tsys string) float64 {
	timeStr, exists := c.T0TSys[tsys]
	if !exists || timeStr == "" {
		return 0
	}

	// Extract numeric part
	numStr := timeStr
	unit := ""

	// Find where the numeric part ends
	for i, ch := range timeStr {
		if !isDigit(ch) && ch != '.' && ch != ',' {
			numStr = timeStr[:i]
			unit = strings.ToLower(timeStr[i:])
			break
		}
	}

	// Convert to float
	numStr = strings.Replace(numStr, ",", ".", 1)
	value, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0
	}

	// Convert to seconds based on unit
	switch {
	case strings.Contains(unit, "s"):
		return value
	case strings.Contains(unit, "min"):
		return value * 60
	case strings.Contains(unit, "h"):
		return value * 3600
	default:
		return value // Assume seconds if no unit specified
	}
}

// IsOriginConnector returns true if this is an origin connector
func (c *Connector) IsOriginConnector() bool {
	return c.Direction == "O"
}

// IsDestinationConnector returns true if this is a destination connector
func (c *Connector) IsDestinationConnector() bool {
	return c.Direction == "D"
}

// AllowsTransportSystem checks if the connector allows the specified transport system
func (c *Connector) AllowsTransportSystem(tsys string) bool {
	systems := strings.Split(c.TSysSet, ",")
	for _, system := range systems {
		if system == tsys {
			return true
		}
	}
	return false
}
