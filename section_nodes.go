package ptvvisum

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// NodeSection represents $NODE section
type NodeSection struct {
	BaseSection
	Nodes []Node
}

// Node represents a single node in the network (typically an intersection)
type Node struct {
	ID              int     // Node identifier (NO)
	Code            string  // Node code
	Name            string  // Node name
	TypeNo          int     // Node type number
	ControlType     int     // Control type (0=uncontrolled, 1=priority, 2=signalized, etc.)
	MainNodeNo      int     // Main node number (for complex intersections)
	XCoord          float64 // X-coordinate
	YCoord          float64 // Y-coordinate
	ZCoord          float64 // Z-coordinate (elevation)
	AddVal1         int     // Additional value 1
	AddVal2         int     // Additional value 2
	AddVal3         int     // Additional value 3
	T0PRT           string  // Base travel time for public transport
	CapPRT          int     // Capacity for private transport
	LaneDef         int     // Lane definition
	Notes           string  // Notes/comments
	RailwayCrossing int     // Railway crossing flag
}

// GetNodeByID retrieves a node by its ID
func (s *NodeSection) GetNodeByID(id int) (Node, bool) {
	for _, node := range s.Nodes {
		if node.ID == id {
			return node, true
		}
	}
	return Node{}, false
}

// GetNodesByType retrieves all nodes of a specified type
func (s *NodeSection) GetNodesByType(typeNo int) []Node {
	var result []Node
	for _, node := range s.Nodes {
		if node.TypeNo == typeNo {
			result = append(result, node)
		}
	}
	return result
}

// GetNodesByControlType retrieves all nodes with a specified control type
func (s *NodeSection) GetNodesByControlType(controlType int) []Node {
	var result []Node
	for _, node := range s.Nodes {
		if node.ControlType == controlType {
			result = append(result, node)
		}
	}
	return result
}

// GetNearbyNodes finds all nodes within a specified distance from given coordinates
func (s *NodeSection) GetNearbyNodes(x, y float64, maxDistance float64) []Node {
	var result []Node
	for _, node := range s.Nodes {
		dx := node.XCoord - x
		dy := node.YCoord - y
		distance := math.Sqrt(dx*dx + dy*dy)
		if distance <= maxDistance {
			result = append(result, node)
		}
	}
	return result
}

// CalculateBoundingBox returns the min/max coordinates of all nodes
func (s *NodeSection) CalculateBoundingBox() (minX, minY, maxX, maxY float64) {
	if len(s.Nodes) == 0 {
		return 0, 0, 0, 0
	}

	minX = s.Nodes[0].XCoord
	minY = s.Nodes[0].YCoord
	maxX = s.Nodes[0].XCoord
	maxY = s.Nodes[0].YCoord

	for _, node := range s.Nodes {
		if node.XCoord < minX {
			minX = node.XCoord
		}
		if node.YCoord < minY {
			minY = node.YCoord
		}
		if node.XCoord > maxX {
			maxX = node.XCoord
		}
		if node.YCoord > maxY {
			maxY = node.YCoord
		}
	}

	return minX, minY, maxX, maxY
}

// Count returns the number of nodes in the section
func (s *NodeSection) Count() int {
	return len(s.Nodes)
}

// getNode extracts data from NODE section row
func getNode(values []string) (Node, error) {
	if len(values) < 10 {
		return Node{}, fmt.Errorf("invalid NODE data (insufficient fields): %v", values)
	}

	var node Node
	var err error

	// Parse NO (required field)
	if values[0] == "" {
		return Node{}, fmt.Errorf("missing required field NO")
	}
	node.ID, err = strconv.Atoi(values[0])
	if err != nil {
		return Node{}, fmt.Errorf("error parsing NO: %w", err)
	}

	// Parse CODE (optional)
	node.Code = values[1]

	// Parse NAME (optional)
	node.Name = values[2]

	// Parse TYPENO (required field)
	if values[3] != "" {
		node.TypeNo, err = strconv.Atoi(values[3])
		if err != nil {
			return Node{}, fmt.Errorf("error parsing TYPENO: %w", err)
		}
	}

	// Parse CONTROLTYPE (required field)
	if values[4] != "" {
		node.ControlType, err = strconv.Atoi(values[4])
		if err != nil {
			return Node{}, fmt.Errorf("error parsing CONTROLTYPE: %w", err)
		}
	}

	// Parse MAINNODENO (required field)
	if values[5] != "" {
		node.MainNodeNo, err = strconv.Atoi(values[5])
		if err != nil {
			return Node{}, fmt.Errorf("error parsing MAINNODENO: %w", err)
		}
	}

	// Skip fields 6 and 7 (USEMETHODIMPATNODE, METHODIMPATNODE) for brevity

	// Parse XCOORD (required field)
	if values[9] == "" {
		return Node{}, fmt.Errorf("missing required field XCOORD")
	}
	node.XCoord, err = strconv.ParseFloat(strings.Replace(values[9], ",", ".", -1), 64)
	if err != nil {
		return Node{}, fmt.Errorf("error parsing XCOORD: %w", err)
	}

	// Parse YCOORD (required field)
	if values[10] == "" {
		return Node{}, fmt.Errorf("missing required field YCOORD")
	}
	node.YCoord, err = strconv.ParseFloat(strings.Replace(values[10], ",", ".", -1), 64)
	if err != nil {
		return Node{}, fmt.Errorf("error parsing YCOORD: %w", err)
	}

	// Parse ZCOORD (optional but usually present)
	if len(values) > 11 && values[11] != "" {
		node.ZCoord, err = strconv.ParseFloat(strings.Replace(values[11], ",", ".", -1), 64)
		if err != nil {
			return Node{}, fmt.Errorf("error parsing ZCOORD: %w", err)
		}
	}

	// Parse ADDVAL1, ADDVAL2, ADDVAL3 (optional)
	if len(values) > 12 && values[12] != "" {
		node.AddVal1, _ = strconv.Atoi(values[12])
	}
	if len(values) > 13 && values[13] != "" {
		node.AddVal2, _ = strconv.Atoi(values[13])
	}
	if len(values) > 14 && values[14] != "" {
		node.AddVal3, _ = strconv.Atoi(values[14])
	}

	// Parse T0PRT (optional) - contains time values like "13s"
	if len(values) > 15 && values[15] != "" {
		node.T0PRT = values[15]
	}

	// Parse CAPPRT (optional) - capacity value
	if len(values) > 16 && values[16] != "" {
		node.CapPRT, _ = strconv.Atoi(strings.Replace(values[16], ",", "", -1))
	}

	// Parse LANEDEF (optional)
	if len(values) > 17 && values[17] != "" {
		node.LaneDef, _ = strconv.Atoi(values[17])
	}

	// Parse NOTES (optional) - available at a specific index in your data
	if len(values) > 31 && values[31] != "" {
		node.Notes = values[31]
	}

	// Parse RAILWAY_CROSSING (optional) - available at the end of your data
	if len(values) > 45 && values[45] != "" {
		node.RailwayCrossing, _ = strconv.Atoi(values[45])
	}

	return node, nil
}
