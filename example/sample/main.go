package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	ptvvisum "github.com/lddl/go-ptv-visum"
)

func main() {
	file, err := os.Open("./example/sample/example.net")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	ptvData, err := ptvvisum.ReadPTVFromFile(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Version:")
	fmt.Printf("\tVersion: %s\n", ptvData.Version.Version)
	fmt.Printf("\tFileType: %s\n", ptvData.Version.FileType)
	fmt.Printf("\tLanguage: %s\n", ptvData.Version.Language)
	fmt.Printf("\tUnit: %s\n", ptvData.Version.Unit)

	fmt.Println("\nInfo:")
	for _, line := range ptvData.Info.Lines {
		fmt.Printf("\tIndex: %d, Text: %s\n", line.Index, line.Text)
	}

	fmt.Println("\nPOI Categories:")
	for _, cat := range ptvData.POICategory.Categories {
		fmt.Printf("\tNo: %d, Code: %s, Name: %s, ParentCatNo: %d\n",
			cat.No, cat.Code, cat.Name, cat.ParentCatNo)
	}

	fmt.Println("\nUser Attributes:")
	for _, attr := range ptvData.UserAttDef.Attributes {
		fmt.Printf("\tObjID: %s, AttID: %s, Code: %s, Name: %s, ValueType: %s\n",
			attr.ObjID, attr.AttID, attr.Code, attr.Name, attr.ValueType)
	}

	fmt.Println("\nCalendar Periods:")
	for _, period := range ptvData.CalendarPeriod.Periods {
		fmt.Printf("\tType: %s, ValidFrom: %s, ValidUntil: %s, AnalysisPeriodStartDayIndex: %d, AnalysisPeriodEndDayIndex: %d, AnalysisTimeIntervalSetNo: %d\n",
			period.Type, period.ValidFrom.Format("02.01.2006"), period.ValidUntil.Format("02.01.2006"),
			period.AnalysisPeriodStartDayIndex, period.AnalysisPeriodEndDayIndex, period.AnalysisTimeIntervalSetNo,
		)
	}

	fmt.Println("\nValid Days:")
	for _, day := range ptvData.ValidDays.Days {
		fmt.Printf("\tNo: %d, Code: %s, Name: %s, DayVector: %d, "+
			"PrfacHourCost: %.3f, PrfacSupply: %.3f\n",
			day.No, day.Code, day.Name, day.DayVector,
			day.PrfacHourCost, day.PrfacSupply)
	}

	fmt.Println("\nNetwork Configuration:")
	fmt.Printf("\tVersion: %s (%s)\n", ptvData.Network.Network.NetVersionID, ptvData.Network.Network.NetVersionName)
	fmt.Printf("\tScale: %.3f %s\n", ptvData.Network.Network.Scale, ptvData.Network.Network.Unit)
	fmt.Printf("\tLeft-hand traffic: %v\n", ptvData.Network.Network.LeftHandTraffic == 1)
	fmt.Printf("\tProjection: %s\n", ptvData.Network.Network.ProjectionDefinition)
	fmt.Printf("\tNetwork name: %s\n", ptvData.Network.Network.Name)

	fmt.Println("\nTransport Systems:")
	for _, system := range ptvData.TSys.Systems {
		fmt.Printf("\tCode: %s, Name: %s, Type: %s, PCU: %.3f\n",
			system.Code, system.Name, system.Type, system.PCU)
	}
	// Filter by system type
	fmt.Println("Public Transport Systems:")
	for _, system := range ptvData.TSys.Systems {
		if system.Type == "PuT" {
			fmt.Printf("\t- %s: %s\n", system.Code, system.Name)
		}
	}

	fmt.Println("\nTransport Modes:")
	for _, mode := range ptvData.Mode.Modes {
		fmt.Printf("\tCode: %s, Name: %s, Interchangeable: %t\n",
			mode.Code, mode.Name, mode.Interchangeable == 1)

		fmt.Printf("\t\tTransport Systems: %s\n", strings.Join(mode.TSysSet, ", "))
	}
	// Find a specific mode
	fmt.Println("Details for PuT mode:")
	for _, mode := range ptvData.Mode.Modes {
		if mode.Code == "PuT" {
			fmt.Printf("\tName: %s\n", mode.Name)
			fmt.Printf("\tTransport Systems: %s\n", strings.Join(mode.TSysSet, ", "))
			fmt.Printf("\tInterchangeable: %t\n", mode.Interchangeable == 1)
			break
		}
	}

	fmt.Println("\nDemand Segments:")
	for _, segment := range ptvData.DemandSegment.Segments {
		fmt.Printf("\tCode: %s, Name: %s, Mode: %s\n",
			segment.Code, segment.Name, segment.Mode)
		fmt.Printf("\t\tOccupancy Rate: %.3f, PrFacAP: %.3f, PrFacAH: %.3f\n",
			segment.OccupancyRate, segment.PrFacAP, segment.PrFacAH)
	}
	// Find segments for a specific mode
	fmt.Println("Freight Segments:")
	for _, segment := range ptvData.DemandSegment.Segments {
		if segment.Mode == "TB" || segment.Mode == "TM" || segment.Mode == "TS" {
			fmt.Printf("\t- %s: %s\n", segment.Code, segment.Name)
		}
	}

	fmt.Println("\nBlock Item Types:")
	for _, itemType := range ptvData.BlockItemType.Types {
		fmt.Printf("\tNo: %d, Name: %s, Default Length: %s\n",
			itemType.No, itemType.Name, itemType.DefLength)
		fmt.Printf("\t\tShort layover weight: %.3f, threshold: %s\n",
			itemType.WeightForLayoversShort, itemType.LayoverThresholdShort)
		fmt.Printf("\t\tLong layover weight: %.3f, threshold: %s\n",
			itemType.WeightForLayoversLong, itemType.LayoverThresholdLong)
	}
	// Find a specific block item type
	fmt.Println("Details for Vehicle Journey:")
	for _, itemType := range ptvData.BlockItemType.Types {
		if itemType.Name == "Vehicle journey" {
			fmt.Printf("\tNo: %d\n", itemType.No)
			fmt.Printf("\tDefault Length: %s\n", itemType.DefLength)
			fmt.Printf("\tShare Before: %.3f\n", itemType.ShareBefore)
			break
		}
	}

	fmt.Println("\nFare Model:")
	fmt.Printf("\tFallback Fare: %.2f\n", ptvData.FareModel.FallbackFare)

	fmt.Println("\nVehicle Units:")
	for _, unit := range ptvData.VehUnit.Units {
		fmt.Printf("\tNo: %d, Code: %s, Name: %s, System: %s\n",
			unit.No, unit.Code, unit.Name, unit.TSysSet)
		fmt.Printf("\t\tCapacities: %d seats, %d total\n",
			unit.SeatCap, unit.TotalCap)
	}
	// Filter by transport system
	fmt.Println("Bus Units:")
	for _, unit := range ptvData.VehUnit.Units {
		if unit.TSysSet == "BUS" {
			fmt.Printf("\t- %s: %s (Seats: %d, Total: %d)\n",
				unit.Code, unit.Name, unit.SeatCap, unit.TotalCap)
		}
	}
	// Find largest vehicle by capacity
	var largestUnit ptvvisum.VehicleUnit
	for _, unit := range ptvData.VehUnit.Units {
		if unit.TotalCap > largestUnit.TotalCap {
			largestUnit = unit
		}
	}
	fmt.Printf("Largest vehicle: %s with capacity for %d passengers\n",
		largestUnit.Name, largestUnit.TotalCap)

	fmt.Println("\nVehicle Combinations:")
	for _, comb := range ptvData.VehComb.Combinations {
		fmt.Printf("\tNo: %d, Code: %s, Set: %s, Name: %s\n",
			comb.No, comb.Code, comb.VehCombSet, comb.Name)
		fmt.Printf("\t\tHour costs: Service=%.2f, Empty=%.2f, Layover=%.2f, Depot=%.2f\n",
			comb.CostRateHourService, comb.CostRateHourEmpty,
			comb.CostRateHourLayover, comb.CostRateHourDepot)
		fmt.Printf("\t\tKm costs: Service=%.2f, Empty=%.2f\n",
			comb.CostRateKmService, comb.CostRateKmEmpty)
	}
	// Find combinations by type
	fmt.Println("Bus Combinations:")
	for _, comb := range ptvData.VehComb.Combinations {
		if strings.HasPrefix(comb.Code, "B") {
			fmt.Printf("\t- %s: %s\n", comb.Code, comb.Name)
		}
	}
	// Find combinations with non-empty set
	fmt.Println("Combinations with defined set:")
	for _, comb := range ptvData.VehComb.Combinations {
		if comb.VehCombSet != "" {
			fmt.Printf("\t- %s: Set %s\n", comb.Name, comb.VehCombSet)
		}
	}

	fmt.Println("\nVehicle Unit to Combination Mappings:")
	for _, mapping := range ptvData.VehUnitToVehComb.Mappings {
		fmt.Printf("\tCombination %d contains %d units of type %d\n",
			mapping.VehCombNo, mapping.NumVehUnits, mapping.VehUnitNo)
	}
	// Find all units for a specific combination
	combNo := 10
	fmt.Printf("Units in combination %d:\n", combNo)
	for _, mapping := range ptvData.VehUnitToVehComb.Mappings {
		if mapping.VehCombNo == combNo {
			fmt.Printf("\t- %d x Unit %d\n", mapping.NumVehUnits, mapping.VehUnitNo)
		}
	}
	// Get detailed information by joining with other sections
	fmt.Println("Detailed Vehicle Combinations:")
	for _, comb := range ptvData.VehComb.Combinations {
		fmt.Printf("\tCombination %d (%s):\n", comb.No, comb.Name)
		// Find all units in this combination
		for _, mapping := range ptvData.VehUnitToVehComb.Mappings {
			if mapping.VehCombNo == comb.No {
				// Find unit details
				for _, unit := range ptvData.VehUnit.Units {
					if unit.No == mapping.VehUnitNo {
						fmt.Printf("\t\t- %d x %s (%d seats each)\n",
							mapping.NumVehUnits, unit.Name, unit.SeatCap)
						break
					}
				}
			}
		}
	}

	fmt.Println("\nDirections:")
	for _, dir := range ptvData.Direction.Directions {
		fmt.Printf("\tNo: %d, Code: %s, Name: %s\n",
			dir.No, dir.Code, dir.Name)
	}
	// You can also look up a direction by number
	dirNo := 1
	fmt.Printf("Direction %d:\n", dirNo)
	for _, dir := range ptvData.Direction.Directions {
		if dir.No == dirNo {
			fmt.Printf("\tCode: %s, Name: %s\n", dir.Code, dir.Name)
			break
		}
	}
	// Or find direction by code
	fmt.Println("Direction by code '>'")
	dirCode := ">"
	for _, dir := range ptvData.Direction.Directions {
		if dir.Code == dirCode {
			fmt.Printf("\tDirection '%s' is number %d with name '%s'\n",
				dirCode, dir.No, dir.Name)
			break
		}
	}

	fmt.Println("\nPoints (first 5 values):")
	for i, point := range ptvData.Point.Points {
		if i < 5 { // Just show the first 5 points
			fmt.Printf("\tID: %d, X: %.4f, Y: %.4f\n",
				point.ID, point.XCoord, point.YCoord)
		}
	}
	fmt.Printf("Total points: %d\n", len(ptvData.Point.Points))
	// Find a specific point by ID
	pointID := 11349
	fmt.Printf("Point %d:\n", pointID)
	for _, point := range ptvData.Point.Points {
		if point.ID == pointID {
			fmt.Printf("\tCoordinates: (%.4f, %.4f)\n", point.XCoord, point.YCoord)
			break
		}
	}
	// Find points within a bounding box
	minX, maxX := 65.65, 65.70
	minY, maxY := 64.55, 64.60
	fmt.Printf("Points within (%.4f, %.4f) - (%.4f, %.4f):\n",
		minX, minY, maxX, maxY)
	count := 0
	for _, point := range ptvData.Point.Points {
		if point.XCoord >= minX && point.XCoord <= maxX &&
			point.YCoord >= minY && point.YCoord <= maxY {
			fmt.Printf("\tID: %d, X: %.4f, Y: %.4f\n",
				point.ID, point.XCoord, point.YCoord)
			count++
		}
	}
	fmt.Printf("Total points in area: %d\n", count)

	fmt.Println("\nEdges (first 5):")
	for i, edge := range ptvData.Edge.Edges {
		if i < 5 { // Just show the first 5 edges
			fmt.Printf("\tID: %d, From: %d, To: %d\n",
				edge.ID, edge.FromPointID, edge.ToPointID)
		}
	}
	fmt.Printf("Total edges: %d\n", len(ptvData.Edge.Edges))
	// Find edges connected to a specific point
	pointID = 11197
	connectedEdges := ptvData.Edge.GetEdgesByPointID(pointID)
	fmt.Printf("Edges connected to point %d: %d edges\n", pointID, len(connectedEdges))
	for _, edge := range connectedEdges {
		fmt.Printf("\tEdge %d: %d → %d\n", edge.ID, edge.FromPointID, edge.ToPointID)
	}
	// Find a specific edge
	edgeID := 188814
	if edge, found := ptvData.Edge.GetEdgeByID(edgeID); found {
		fmt.Printf("Edge %d connects points %d and %d\n",
			edge.ID, edge.FromPointID, edge.ToPointID)
		// If points are also loaded, we can get their coordinates
		if ptvData.Point != nil {
			if fromPoint, foundFrom := ptvData.Point.GetPointByID(edge.FromPointID); foundFrom {
				if toPoint, foundTo := ptvData.Point.GetPointByID(edge.ToPointID); foundTo {
					fmt.Printf("\tCoordinates: (%.4f, %.4f) → (%.4f, %.4f)\n",
						fromPoint.XCoord, fromPoint.YCoord,
						toPoint.XCoord, toPoint.YCoord)
					// Calculate edge length
					distance := fromPoint.Distance(toPoint)
					fmt.Printf("\tLength: %.4f units\n", distance)
				}
			}
		}
	}

	fmt.Println("\tEdge Items:")
	// Get unique edge IDs
	edgeIDs := make(map[int]bool)
	for _, item := range ptvData.EdgeItem.Items {
		edgeIDs[item.EdgeID] = true
	}
	fmt.Printf("Intermediate points for %d edges\n", len(edgeIDs))
	// Print a few examples
	for edgeID := range edgeIDs {
		items := ptvData.EdgeItem.GetItemsByEdgeID(edgeID)
		fmt.Printf("\tEdge %d has %d intermediate points (first 5):\n", edgeID, len(items))
		// Print first 5 points as example
		for i, item := range items {
			if i < 5 {
				fmt.Printf("\t\t%d: (%.4f, %.4f)\n", item.Index, item.XCoord, item.YCoord)
			}
		}
		// Calculate length
		length := ptvData.EdgeItem.CalculateEdgeLength(edgeID, ptvData)
		fmt.Printf("\tTotal length: %.4f units\n", length)
		// Only show a few examples
		if len(edgeIDs) > 3 {
			break
		}
	}

	fmt.Println("\nFaces:")
	fmt.Printf("Total faces: %d\n", ptvData.Face.Count())
	// Print the first 10 face IDs
	fmt.Println("Face IDs (first 10):")
	for i, face := range ptvData.Face.Faces {
		if i < 10 {
			fmt.Printf("  %d", face.ID)
			// Add newlines for better formatting
			if i%5 == 4 {
				fmt.Println()
			}
		}
	}
	// Check if a specific face exists
	faceID := 428
	if ptvData.Face.Contains(faceID) {
		fmt.Printf("Face %d exists in the data\n", faceID)
	} else {
		fmt.Printf("Face %d does not exist in the data\n", faceID)
	}

	fmt.Println("\nFace Items:")
	// Count items per face
	faceItems := make(map[int]int)
	for _, item := range ptvData.FaceItem.Items {
		faceItems[item.FaceID]++
	}
	fmt.Printf("Total faces with items: %d\n", len(faceItems))
	fmt.Println("Show details for a first 3 faces:")
	count = 0
	// Show details for a few faces
	for faceID, numItems := range faceItems {
		items := ptvData.FaceItem.GetItemsByFaceID(faceID)
		fmt.Printf("\tFace %d has %d edges:\n", faceID, numItems)
		for _, item := range items {
			fmt.Printf("\t\tEdge %d (index %d, direction %d)\n",
				item.EdgeID, item.Index, item.Direction)
		}
		// Calculate area if possible
		area := ptvData.FaceItem.CalculateFaceArea(faceID, ptvData)
		if area > 0 {
			fmt.Printf("  Approximate area: %.2f square units\n", area)
		}
		count++
		if count >= 3 {
			break // Only show first 3 faces
		}
	}

	fmt.Println("\nSurfaces:")
	fmt.Printf("Total surfaces: %d\n", ptvData.Surface.Count())
	// Print all surface IDs
	ids := ptvData.Surface.GetAllIDs()
	fmt.Println("Surface IDs:")
	for i, id := range ids {
		fmt.Printf("\t%d", id)
		if i < len(ids)-1 {
			fmt.Print(", ")
		}
		// Add line breaks for readability
		if (i+1)%10 == 0 {
			fmt.Print("\n")
		}
	}
	fmt.Println()
	// Check if a specific surface exists
	surfaceID := 196
	if ptvData.Surface.Contains(surfaceID) {
		fmt.Printf("Surface %d exists in the data\n", surfaceID)
	} else {
		fmt.Printf("Surface %d does not exist in the data\n", surfaceID)
	}

	fmt.Println("\nSurface Items:")
	// Get unique surface IDs
	surfaceIDs := ptvData.SurfaceItem.GetSurfaceIDs()
	fmt.Printf("Total surfaces with items: %d\n", len(surfaceIDs))
	// Print details for each surface
	for _, surfaceID := range surfaceIDs {
		outer, inner := ptvData.SurfaceItem.GetBoundariesBySurfaceID(surfaceID)
		fmt.Printf("\tSurface %d:\n", surfaceID)
		fmt.Printf("\t\tOuter boundaries: %d faces\n", len(outer))
		fmt.Printf("\t\tInner holes (enclaves): %d faces\n", len(inner))
		// Print the face IDs for smaller surfaces (limit output)
		if len(outer)+len(inner) < 10 {
			fmt.Printf("\t\tOuter faces: %v\n", outer)
			if len(inner) > 0 {
				fmt.Printf("\t\tInner faces: %v\n", inner)
			}
		}
	}
	// Check which surfaces a particular face belongs to
	faceID = 415
	items := ptvData.SurfaceItem.GetItemsByFaceID(faceID)
	fmt.Printf("Face %d belongs to %d surfaces:\n", faceID, len(items))
	for _, item := range items {
		enclave := "outer boundary"
		if item.Enclave == 1 {
			enclave = "inner hole"
		}
		fmt.Printf("\tSurface %d (%s)\n", item.SurfaceID, enclave)
	}

	fmt.Println("\nNodes:")
	fmt.Printf("Total nodes: %d\n", ptvData.Node.Count())
	// Print the bounding box of the network
	minX, minY, maxX, maxY = ptvData.Node.CalculateBoundingBox()
	fmt.Printf("Network extent: (%.4f, %.4f) to (%.4f, %.4f)\n",
		minX, minY, maxX, maxY)
	// Analyze node types
	controlTypes := make(map[int]int)
	for _, node := range ptvData.Node.Nodes {
		controlTypes[node.ControlType]++
	}
	fmt.Println("Node control types:")
	for controlType, count := range controlTypes {
		var typeName string
		switch controlType {
		case 0:
			typeName = "Uncontrolled"
		case 1:
			typeName = "Priority"
		case 2:
			typeName = "Signalized"
		case 3:
			typeName = "Stop sign"
		default:
			typeName = "Other"
		}
		fmt.Printf("\t%s: %d nodes\n", typeName, count)
	}
	// Find the highest capacity node
	var maxCapNode ptvvisum.Node
	maxCap := -1
	for _, node := range ptvData.Node.Nodes {
		if node.CapPRT > maxCap {
			maxCap = node.CapPRT
			maxCapNode = node
		}
	}
	if maxCap > 0 {
		fmt.Printf("Highest capacity node: ID=%d with capacity %d\n",
			maxCapNode.ID, maxCapNode.CapPRT)
	}
	// Find nearby nodes for spatial analysis
	centerX := (minX + maxX) / 2
	centerY := (minY + maxY) / 2
	radius := (maxX - minX) * 0.1 // 10% of the network width
	nearbyNodes := ptvData.Node.GetNearbyNodes(centerX, centerY, radius)
	fmt.Printf("Found %d nodes within %.2f units of the network center\n", len(nearbyNodes), radius)

	fmt.Println("\nZones:")
	fmt.Printf("Total zones: %d\n", ptvData.Zone.Count())
	// Print overall statistics
	totalPop := ptvData.Zone.GetTotalPopulation()
	totalEmp := ptvData.Zone.GetTotalEmployment()
	fmt.Printf("Total population: %d\n", totalPop)
	fmt.Printf("Total employment: %d\n", totalEmp)
	fmt.Printf("Average jobs per person: %.2f\n", float64(totalEmp)/float64(totalPop))
	// Get bounding box
	minX, minY, maxX, maxY = ptvData.Zone.CalculateBoundingBox()
	fmt.Printf("Study area extent: (%.4f, %.4f) to (%.4f, %.4f)\n",
		minX, minY, maxX, maxY)
	// Print details of a specific zone
	zoneID := 4 // Губкинский
	if zone, found := ptvData.Zone.GetZoneByID(zoneID); found {
		fmt.Printf("\tZone %d: %s\n", zone.No, zone.Name)
		fmt.Printf("\t\tLocation: (%.4f, %.4f)\n", zone.XCoord, zone.YCoord)
		fmt.Printf("\t\tPopulation: %d\n", zone.Population)
		fmt.Printf("\t\tEmployment: %d\n", zone.Employment)
		fmt.Printf("\t\tWorkers: %d\n", zone.Workers)
		fmt.Printf("\t\tStudents: %d\n", zone.Students)
		// Employment self-sufficiency ratio
		if zone.Workers > 0 {
			fmt.Printf("\t\tJobs/Workers ratio: %.2f\n", float64(zone.Employment)/float64(zone.Workers))
		}
		// Get geometry if available
		if zone.SurfaceID > 0 && ptvData.Surface != nil && ptvData.SurfaceItem != nil {
			if ptvData.Surface.Contains(zone.SurfaceID) {
				outer, inner := ptvData.SurfaceItem.GetBoundariesBySurfaceID(zone.SurfaceID)
				fmt.Printf("\t\tZone boundary: %d outer faces, %d inner holes\n", len(outer), len(inner))
			}
		}
	}
	// Find zones by type
	typeNo := 2
	typeZones := ptvData.Zone.GetZonesByType(typeNo)
	fmt.Printf("Zones of type %d: %d zones\n", typeNo, len(typeZones))
	for _, zone := range typeZones {
		fmt.Printf("\t%d: %s\n", zone.No, zone.Name)
	}
	// Sort zones by population (largest first)
	sort.Slice(ptvData.Zone.Zones, func(i, j int) bool {
		return ptvData.Zone.Zones[i].Population > ptvData.Zone.Zones[j].Population
	})
	fmt.Println("Top 5 zones by population:")
	for i, zone := range ptvData.Zone.Zones {
		if i < 5 {
			fmt.Printf("\t%d: %s - %d residents\n",
				zone.No, zone.Name, zone.Population)
		}
	}

	fmt.Println("\nLink Types:")
	fmt.Printf("Total link types: %d\n", ptvData.LinkType.Count())
	// Get distribution by group type
	groupTypes := make(map[int]int)
	for _, linkType := range ptvData.LinkType.LinkTypes {
		groupTypes[linkType.GroupType]++
	}
	fmt.Println("Link types by group type:")
	for groupType, count := range groupTypes {
		fmt.Printf("\tGroup %d: %d link types\n", groupType, count)
		// Show sample link types from this group
		types := ptvData.LinkType.GetLinkTypesByGroupType(groupType)
		fmt.Printf("\t\tSample link types number: %d\n", len(types))
	}

	fmt.Println("\nLinks:")
	fmt.Printf("Total links: %d\n", ptvData.Link.Count())
	// Calculate network statistics
	totalLength := ptvData.Link.GetTotalNetworkLength()
	avgSpeed := ptvData.Link.GetAverageSpeed()
	fmt.Printf("Total network length: %.2f km\n", totalLength)
	fmt.Printf("Average network speed: %.2f km/h\n", avgSpeed)
	// Get bidirectional links
	biLinks := ptvData.Link.GetBidirectionalLinks()
	fmt.Printf("Bidirectional links: %d (%.1f%%)\n",
		len(biLinks), float64(len(biLinks))*100/float64(ptvData.Link.Count()))
	// Find links by type
	typeNo = 16 // Highway links
	typeLinks := ptvData.Link.GetLinksByType(typeNo)
	fmt.Printf("Links of type %d: %d links\n", typeNo, len(typeLinks))
	// Get links by transport system
	busTSys := "BUS"
	busLinks := ptvData.Link.GetLinksByTransportSystem(busTSys)
	fmt.Printf("Links allowing %s: %d links (%.1f%%)\n",
		busTSys, len(busLinks), float64(len(busLinks))*100/float64(ptvData.Link.Count()))
	// Analyze connectivity at a specific node
	nodeNo := 7
	fromLinks := ptvData.Link.GetLinksByFromNode(nodeNo)
	toLinks := ptvData.Link.GetLinksByToNode(nodeNo)
	fmt.Printf("Node %d connectivity:\n", nodeNo)
	fmt.Printf("\tOutgoing links: %d\n", len(fromLinks))
	fmt.Printf("\tIncoming links: %d\n", len(toLinks))
	// Get connected nodes
	adjacentNodes := ptvData.Link.GetAdjacentNodes(nodeNo)
	fmt.Printf("\tConnected to %d nodes: %v\n", len(adjacentNodes), adjacentNodes)

	// Calculate network density
	connectedNodes := ptvData.Link.GetConnectedNodes()
	fmt.Printf("Network density:\n")
	fmt.Printf("\tLinks: %d\n", ptvData.Link.Count())
	fmt.Printf("\tConnected nodes: %d\n", len(connectedNodes))
	fmt.Printf("\tLinks/node ratio: %.2f\n", float64(ptvData.Link.Count())/float64(len(connectedNodes)))

	// Get details for a specific link
	linkID := 6
	if link, found := ptvData.Link.GetLinkByID(linkID); found {
		fmt.Printf("\nLink %d details:\n", linkID)
		fmt.Printf("\tFrom node %d to node %d\n", link.FromNodeNo, link.ToNodeNo)
		fmt.Printf("\tName: %s\n", link.Name)
		fmt.Printf("\tType: %d\n", link.TypeNo)
		fmt.Printf("\tLength: %s (%.3f km)\n", link.Length, link.GetLengthInKm())
		fmt.Printf("\tSpeed: %s (%.1f km/h)\n", link.V0PRT, link.GetSpeedInKmh())
		fmt.Printf("\tCapacity: %d vehicles/hour\n", link.CapPRT)
		fmt.Printf("\tLanes: %d\n", link.NumLanes)
		fmt.Printf("\tTransport systems: %s\n", link.TSysSet)
		if link.SpeedLimit > 0 {
			fmt.Printf("\tSpeed limit: %d km/h\n", link.SpeedLimit)
		}
		// If we have node data, show node names
		if ptvData.Node != nil {
			if fromNode, found := ptvData.Node.GetNodeByID(link.FromNodeNo); found {
				fmt.Printf("\tFrom node name: %s\n", fromNode.Name)
			}
			if toNode, found := ptvData.Node.GetNodeByID(link.ToNodeNo); found {
				fmt.Printf("\tTo node name: %s\n", toNode.Name)
			}
		}
	}
}
