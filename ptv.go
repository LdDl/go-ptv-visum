package ptvvisum

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// PTVData represents the complete PTV Visum network file data
type PTVData struct {
	Version          *VersionSection
	Info             *InfoSection
	POICategory      *POICategorySection
	UserAttDef       *UserAttDefSection
	CalendarPeriod   *CalendarPeriodSection
	ValidDays        *ValidDaysSection
	Network          *NetworkSection
	TSys             *TSysSection
	Mode             *ModeSection
	DemandSegment    *DemandSegmentSection
	BlockItemType    *BlockItemTypeSection
	FareModel        *FareModelSection
	VehUnit          *VehUnitSection
	VehComb          *VehCombSection
	VehUnitToVehComb *VehUnitToVehCombSection
	Direction        *DirectionSection
	Point            *PointSection
	Edge             *EdgeSection
	EdgeItem         *EdgeItemSection
	Face             *FaceSection
	FaceItem         *FaceItemSection
	Surface          *SurfaceSection
	SurfaceItem      *SurfaceItemSection
	Node             *NodeSection
	Zone             *ZoneSection
	LinkType         *LinkTypeSection
	Link             *LinkSection
	LinkPoly         *LinkPolySection
	Turn             *TurnSection
	Connector        *ConnectorSection

	Sections map[string]Section // Generic access to all sections
}

// ReadPTVFromFile parses a PTV Visum network file
func ReadPTVFromFile(reader io.Reader) (*PTVData, error) {
	data := &PTVData{
		Sections: make(map[string]Section),
	}

	scanner := bufio.NewScanner(reader)
	var currentSection *BaseSection

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines
		if line == "" {
			continue
		}

		// Handle comments
		if strings.HasPrefix(line, "*") {
			// Process comments if needed
			continue
		}

		// Handle section headers
		if strings.HasPrefix(line, "$") {
			sectionParts := strings.SplitN(line, ":", 2)
			sectionName := strings.TrimPrefix(sectionParts[0], "$")

			// Create new section
			currentSection = &BaseSection{
				name: sectionName,
				rows: [][]string{},
			}

			// Parse headers if present
			if len(sectionParts) > 1 && sectionParts[1] != "" {
				currentSection.headers = strings.Split(sectionParts[1], ";")
			}

			// Store section in the data structure
			data.Sections[sectionName] = currentSection

			// Create specialized section if supported
			switch sectionName {
			case "VERSION":
				data.Version = &VersionSection{BaseSection: *currentSection}
			case "INFO":
				data.Info = &InfoSection{BaseSection: *currentSection}
			case "POICATEGORY":
				data.POICategory = &POICategorySection{BaseSection: *currentSection}
			case "USERATTDEF":
				data.UserAttDef = &UserAttDefSection{BaseSection: *currentSection}
			case "CALENDARPERIOD":
				data.CalendarPeriod = &CalendarPeriodSection{BaseSection: *currentSection}
			case "VALIDDAYS":
				data.ValidDays = &ValidDaysSection{BaseSection: *currentSection}
			case "NETWORK":
				data.Network = &NetworkSection{BaseSection: *currentSection}
			case "TSYS":
				data.TSys = &TSysSection{BaseSection: *currentSection}
			case "MODE":
				data.Mode = &ModeSection{BaseSection: *currentSection}
			case "DEMANDSEGMENT":
				data.DemandSegment = &DemandSegmentSection{BaseSection: *currentSection}
			case "BLOCKITEMTYPE":
				data.BlockItemType = &BlockItemTypeSection{BaseSection: *currentSection}
			case "FAREMODEL":
				data.FareModel = &FareModelSection{BaseSection: *currentSection}
			case "VEHUNIT":
				data.VehUnit = &VehUnitSection{BaseSection: *currentSection}
			case "VEHCOMB":
				data.VehComb = &VehCombSection{BaseSection: *currentSection}
			case "VEHUNITTOVEHCOMB":
				data.VehUnitToVehComb = &VehUnitToVehCombSection{BaseSection: *currentSection}
			case "DIRECTION":
				data.Direction = &DirectionSection{BaseSection: *currentSection}
			case "POINT":
				data.Point = &PointSection{BaseSection: *currentSection}
			case "EDGE":
				data.Edge = &EdgeSection{BaseSection: *currentSection}
			case "EDGEITEM":
				data.EdgeItem = &EdgeItemSection{BaseSection: *currentSection}
			case "FACE":
				data.Face = &FaceSection{BaseSection: *currentSection}
			case "FACEITEM":
				data.FaceItem = &FaceItemSection{BaseSection: *currentSection}
			case "SURFACE", "SURFACE:ID": // Handle both formats
				data.Surface = &SurfaceSection{BaseSection: *currentSection}
			case "SURFACEITEM":
				data.SurfaceItem = &SurfaceItemSection{BaseSection: *currentSection}
			case "NODE":
				data.Node = &NodeSection{BaseSection: *currentSection}
			case "ZONE":
				data.Zone = &ZoneSection{BaseSection: *currentSection}
			case "LINKTYPE":
				data.LinkType = &LinkTypeSection{BaseSection: *currentSection}
			case "LINK":
				data.Link = &LinkSection{BaseSection: *currentSection}
			case "LINKPOLY":
				data.LinkPoly = &LinkPolySection{BaseSection: *currentSection}
			case "TURN":
				data.Turn = &TurnSection{BaseSection: *currentSection}
			case "CONNECTOR":
				data.Connector = &ConnectorSection{BaseSection: *currentSection}

			default:
				// return nil, fmt.Errorf("unsupported section: %s", sectionName)
			}

			continue
		}

		// Process data rows
		if currentSection != nil {
			values := strings.Split(line, ";")
			currentSection.AddRow(values)

			// Process specific section data
			switch currentSection.name {
			case "VERSION":
				if data.Version != nil && len(values) >= 4 {
					version, fileType, language, unit, err := getVersion(values)
					if err != nil {
						return nil, fmt.Errorf("error parsing VERSION data: %w", err)
					}
					data.Version.Version = version
					data.Version.FileType = fileType
					data.Version.Language = language
					data.Version.Unit = unit
				}
			case "INFO":
				if data.Info != nil && len(values) >= 2 {
					infoLine, err := getInfoLine(values)
					if err != nil {
						return nil, fmt.Errorf("error parsing INFO data: %w", err)
					}
					data.Info.Lines = append(data.Info.Lines, infoLine)
				}
			case "POICATEGORY":
				if data.POICategory != nil && len(values) >= 5 {
					poiCategory, err := getPoiCategory(values)
					if err != nil {
						return nil, fmt.Errorf("error parsing POICATEGORY data: %w", err)
					}
					data.POICategory.Categories = append(data.POICategory.Categories, poiCategory)
				}
			case "USERATTDEF":
				if data.UserAttDef != nil {
					attr, err := getUserAttDef(values)
					if err != nil {
						return nil, fmt.Errorf("error parsing USERATTDEF data: %w", err)
					}
					data.UserAttDef.Attributes = append(data.UserAttDef.Attributes, attr)
				}
			case "CALENDARPERIOD":
				if data.CalendarPeriod != nil {
					period, err := getCalendarPeriod(values)
					if err != nil {
						return nil, fmt.Errorf("error parsing CALENDARPERIOD data: %w", err)
					}
					data.CalendarPeriod.Periods = append(data.CalendarPeriod.Periods, period)
				}
			case "VALIDDAYS":
				if data.ValidDays != nil {
					day, err := getValidDay(values)
					if err != nil {
						return nil, fmt.Errorf("error parsing VALIDDAYS data: %w", err)
					}
					data.ValidDays.Days = append(data.ValidDays.Days, day)
				}
			case "NETWORK":
				if data.Network != nil {
					network, err := getNetwork(values)
					if err != nil {
						return nil, fmt.Errorf("error parsing NETWORK data: %w", err)
					}
					data.Network.Network = network
				}
			case "TSYS":
				if data.TSys != nil {
					system, err := getTransportSystem(values)
					if err != nil {
						return nil, fmt.Errorf("error parsing TSYS data: %w", err)
					}
					data.TSys.Systems = append(data.TSys.Systems, system)
				}
			case "MODE":
				if data.Mode != nil {
					mode, err := getMode(values)
					if err != nil {
						return nil, fmt.Errorf("error parsing MODE data: %w", err)
					}
					data.Mode.Modes = append(data.Mode.Modes, mode)
				}
			case "DEMANDSEGMENT":
				if data.DemandSegment != nil {
					segment, err := getDemandSegment(values)
					if err != nil {
						return nil, fmt.Errorf("error parsing DEMANDSEGMENT data: %w", err)
					}
					data.DemandSegment.Segments = append(data.DemandSegment.Segments, segment)
				}
			case "BLOCKITEMTYPE":
				if data.BlockItemType != nil {
					itemType, err := getBlockItemType(values)
					if err != nil {
						return nil, fmt.Errorf("error parsing BLOCKITEMTYPE data: %w", err)
					}
					data.BlockItemType.Types = append(data.BlockItemType.Types, itemType)
				}
			case "FAREMODEL":
				if data.FareModel != nil {
					fallbackFare, err := getFallbackFare(values)
					if err != nil {
						return nil, fmt.Errorf("error parsing FAREMODEL data: %w", err)
					}
					data.FareModel.FallbackFare = fallbackFare
				}
			case "VEHUNIT":
				if data.VehUnit != nil {
					unit, err := getVehicleUnit(values)
					if err != nil {
						return nil, fmt.Errorf("error parsing VEHUNIT data: %w", err)
					}
					data.VehUnit.Units = append(data.VehUnit.Units, unit)
				}
			case "VEHCOMB":
				if data.VehComb != nil {
					comb, err := getVehicleCombination(values)
					if err != nil {
						return nil, fmt.Errorf("error parsing VEHCOMB data: %w", err)
					}
					data.VehComb.Combinations = append(data.VehComb.Combinations, comb)
				}
			case "VEHUNITTOVEHCOMB":
				if data.VehUnitToVehComb != nil {
					mapping, err := getVehUnitToVehCombMapping(values)
					if err != nil {
						return nil, fmt.Errorf("error parsing VEHUNITTOVEHCOMB data: %w", err)
					}
					data.VehUnitToVehComb.Mappings = append(data.VehUnitToVehComb.Mappings, mapping)
				}
			case "DIRECTION":
				if data.Direction != nil {
					direction, err := getDirection(values)
					if err != nil {
						return nil, fmt.Errorf("error parsing DIRECTION data: %w", err)
					}
					data.Direction.Directions = append(data.Direction.Directions, direction)
				}
			case "POINT":
				if data.Point != nil {
					point, err := getPoint(values)
					if err != nil {
						return nil, fmt.Errorf("error parsing POINT data: %w", err)
					}
					data.Point.Points = append(data.Point.Points, point)
				}
			case "EDGE":
				if data.Edge != nil {
					edge, err := getEdge(values)
					if err != nil {
						return nil, fmt.Errorf("error parsing EDGE data: %w", err)
					}
					data.Edge.Edges = append(data.Edge.Edges, edge)
				}
			case "EDGEITEM":
				if data.EdgeItem != nil {
					item, err := getEdgeItem(values)
					if err != nil {
						return nil, fmt.Errorf("error parsing EDGEITEM data: %w", err)
					}
					data.EdgeItem.Items = append(data.EdgeItem.Items, item)
				}
			case "FACE":
				if data.Face != nil {
					face, err := getFace(values)
					if err != nil {
						return nil, fmt.Errorf("error parsing FACE data: %w", err)
					}
					data.Face.Faces = append(data.Face.Faces, face)
				}
			case "FACEITEM":
				if data.FaceItem != nil {
					item, err := getFaceItem(values)
					if err != nil {
						return nil, fmt.Errorf("error parsing FACEITEM data: %w", err)
					}
					data.FaceItem.Items = append(data.FaceItem.Items, item)
				}
			case "SURFACE", "SURFACE:ID": // Handle both formats
				if data.Surface != nil {
					surface, err := getSurface(values)
					if err != nil {
						return nil, fmt.Errorf("error parsing SURFACE data: %w", err)
					}
					data.Surface.Surfaces = append(data.Surface.Surfaces, surface)
				}
			case "SURFACEITEM":
				if data.SurfaceItem != nil {
					item, err := getSurfaceItem(values)
					if err != nil {
						return nil, fmt.Errorf("error parsing SURFACEITEM data: %w", err)
					}
					data.SurfaceItem.Items = append(data.SurfaceItem.Items, item)
				}
			case "NODE":
				if data.Node != nil {
					node, err := getNode(values)
					if err != nil {
						return nil, fmt.Errorf("error parsing NODE data: %w", err)
					}
					data.Node.Nodes = append(data.Node.Nodes, node)
				}
			case "ZONE":
				if data.Zone != nil {
					zone, err := getZone(values)
					if err != nil {
						return nil, fmt.Errorf("error parsing ZONE data: %w", err)
					}
					data.Zone.Zones = append(data.Zone.Zones, zone)
				}
			case "LINKTYPE":
				if data.LinkType != nil {
					// FIXED: Call Headers() method to get the string slice
					linkType, err := getLinkType(values, data.LinkType.Headers())
					if err != nil {
						return nil, fmt.Errorf("error parsing LINKTYPE data: %w", err)
					}
					data.LinkType.LinkTypes = append(data.LinkType.LinkTypes, linkType)
				}
			case "LINK":
				if data.Link != nil {
					link, err := getLink(values, data.Link.Headers())
					if err != nil {
						return nil, fmt.Errorf("error parsing LINK data: %w", err)
					}
					data.Link.Links = append(data.Link.Links, link)
				}
			case "LINKPOLY":
				if data.LinkPoly != nil {
					point, err := getLinkPolyPoint(values)
					if err != nil {
						return nil, fmt.Errorf("error parsing LINKPOLY data: %w", err)
					}
					data.LinkPoly.Points = append(data.LinkPoly.Points, point)
				}
			case "TURN":
				if data.Turn != nil {
					turn, err := getTurn(values)
					if err != nil {
						return nil, fmt.Errorf("error parsing TURN data: %w", err)
					}
					data.Turn.Turns = append(data.Turn.Turns, turn)
				}
			case "CONNECTOR":
				if data.Connector != nil {
					connector, err := getConnector(values, data.Connector.Headers())
					if err != nil {
						return nil, fmt.Errorf("error parsing CONNECTOR data: %w", err)
					}
					data.Connector.Connectors = append(data.Connector.Connectors, connector)
				}
			default:
				// return nil, fmt.Errorf("unsupported section: %s", currentSection.name)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading PTV file: %w", err)
	}

	return data, nil
}
