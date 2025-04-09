package graph

import (
	"fmt"
	"sort"

	ptvvisum "github.com/lddl/go-ptv-visum"
	"github.com/lddl/go-ptv-visum/utils"
)

// Edge is a struct representing a graph edge
type Edge struct {
	ID       int
	Source   int
	Target   int
	Geometry [][]float64
	LinkID   int

	LanesNum int
	// meters
	Length float64
	// km/h
	FreeFlowSpeed float64
	Capacity      int
}

// Vertex is a struct representing a graph vertex
type Vertex struct {
	ID int
	X  float64
	Y  float64
}

// Graph contains set of vertices and edges
type Graph struct {
	// Vertices represented as a set
	Vertices map[int]*Vertex
	// Edges represented as a set rather than an adjacency matrix
	Edges map[int]*Edge
}

// ExtractGraph prepares set of vertices and edges with geometry from the given PTV data
func ExtractGraph(ptv *ptvvisum.PTVData) (Graph, error) {
	vertices := make(map[int]*Vertex)
	edges := make(map[int]*Edge)
	edgeID := 0
	// Map edges to node pairs for using reversed edges
	mapEdges := make(map[int]map[int]*Edge)

	// Process nodes
	if ptv.Node == nil {
		return Graph{}, fmt.Errorf("no nodes found in the data")
	}
	for _, node := range ptv.Node.Nodes {
		vertex := &Vertex{
			ID: node.ID,
			X:  node.XCoord,
			Y:  node.YCoord,
		}
		vertices[node.ID] = vertex
	}

	// Process edges
	if ptv.Link == nil {
		return Graph{}, fmt.Errorf("no edges found in the data")
	}
	for _, link := range ptv.Link.Links {
		edgeID++
		fromNodeID := link.FromNodeNo
		toNodeID := link.ToNodeNo

		fromNode, ok := vertices[fromNodeID]
		if !ok {
			return Graph{}, fmt.Errorf("from node %d not found for link %d", fromNodeID, link.No)
		}
		toNode, ok := vertices[toNodeID]
		if !ok {
			return Graph{}, fmt.Errorf("to node %d not found for link %d", toNodeID, link.No)
		}
		geometry := buildLinkGeometry(ptv, link, fromNode, toNode, mapEdges)
		if geometry == nil {
			return Graph{}, fmt.Errorf("no geometry found for link %d", link.No)
		}
		// Create edge
		edge := &Edge{
			ID:       edgeID,
			Source:   fromNodeID,
			Target:   toNodeID,
			Geometry: geometry,
			LinkID:   link.No,

			LanesNum: link.NumLanes,
			Capacity: link.CapPRT,
		}
		length, err := utils.ParseLengthValue(link.Length)
		if err != nil {
			return Graph{}, fmt.Errorf("failed to parse length for link %d: %w", link.No, err)
		}
		edge.Length = length
		freeFlowSpeed, err := utils.ParseSpeedValue(link.V0PRT)
		if err != nil {
			return Graph{}, fmt.Errorf("failed to parse free flow speed for link %d: %w", link.No, err)
		}
		edge.FreeFlowSpeed = freeFlowSpeed
		edges[edgeID] = edge
		if _, ok := mapEdges[fromNodeID]; !ok {
			mapEdges[fromNodeID] = make(map[int]*Edge)
		}
		mapEdges[fromNodeID][toNodeID] = edge
	}
	return Graph{
		Vertices: vertices,
		Edges:    edges,
	}, nil
}

func buildLinkGeometry(ptv *ptvvisum.PTVData, link ptvvisum.Link, fromNode, toNode *Vertex, mapEdges map[int]map[int]*Edge) [][]float64 {
	// 1. Start with default straight-line geometry
	defaultGeometry := [][]float64{
		{fromNode.X, fromNode.Y},
		{toNode.X, toNode.Y},
	}

	// 2. Check if we have intermediate points in EdgeItem section
	if ptv.EdgeItem != nil {
		// Get intermediate points for this link
		edgeItems := ptv.EdgeItem.GetItemsByEdgeID(link.No)
		// If we have intermediate points, create a curved geometry
		if len(edgeItems) > 0 {
			// Sort points by index to ensure correct order
			sort.Slice(edgeItems, func(i, j int) bool {
				return edgeItems[i].Index < edgeItems[j].Index
			})
			// Create complete geometry: source node + intermediate points + target node
			geometry := [][]float64{{fromNode.X, fromNode.Y}}
			// Add intermediate points
			for _, item := range edgeItems {
				geometry = append(geometry, []float64{item.XCoord, item.YCoord})
			}
			geometry = append(geometry, []float64{toNode.X, toNode.Y})
			return geometry
		}
	}

	// 3. Check if we have geometry in LinkPoly section
	if ptv.LinkPoly != nil && ptv.LinkPoly.HasLinkGeometry(link.FromNodeNo, link.ToNodeNo) {
		linkGeom := ptv.LinkPoly.GetLinkGeometry(link.FromNodeNo, link.ToNodeNo)
		if len(linkGeom) >= 2 {
			geometry := [][]float64{{fromNode.X, fromNode.Y}} // Start with from-node
			for _, point := range linkGeom {
				geometry = append(geometry, []float64{point[0], point[1]})
			}
			geometry = append(geometry, []float64{toNode.X, toNode.Y})
			// Reverse geometry for the opposite direction:
			// If reversed edge has only two points, reverse the geometry
			if _, ok := mapEdges[toNode.ID]; ok {
				if existingEdge, ok := mapEdges[toNode.ID][fromNode.ID]; ok {
					if len(existingEdge.Geometry) == 2 {
						existingEdge.Geometry = utils.ReverseGeometry(geometry)
					}
				}
			}
			return geometry
		}
	}

	// 4. Check if reversed geometry already exists
	if _, ok := mapEdges[toNode.ID]; ok {
		if existingEdge, ok := mapEdges[toNode.ID][fromNode.ID]; ok {
			// Reverse the geometry
			if len(existingEdge.Geometry) != 2 {
				reversedGeometry := utils.ReverseGeometry(existingEdge.Geometry)
				return reversedGeometry
			}
		}
	}

	// 5. Return simple straight line for default case
	return defaultGeometry
}
