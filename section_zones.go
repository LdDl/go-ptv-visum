package ptvvisum

import (
	"fmt"
	"strconv"
	"strings"
)

// ZoneSection represents $ZONE section
type ZoneSection struct {
	BaseSection
	Zones []Zone
}

// Zone represents a single zone in the transportation model
type Zone struct {
	No               int     // Zone ID (NO)
	Code             string  // Zone code
	Name             string  // Zone name
	MainZoneNo       int     // Main zone number (for aggregation)
	TypeNo           int     // Zone type number
	XCoord           float64 // X-coordinate
	YCoord           float64 // Y-coordinate
	SurfaceID        int     // ID of the surface that defines zone boundary
	RelativeState    int     // Relative state
	SharePRTOrig     float64 // Park and ride share for origins
	SharePRTDest     float64 // Park and ride share for destinations
	SharePUT         float64 // Public transport share
	MethodConnShares int     // Method for connection shares
	Population       int     // Population
	Employment       int     // Employment (Workplaces)
	Workers          int     // Workers
	Students         int     // Students
	StudyPlaces      int     // Study places
	PopDens          float64 // Population density
	Comment          string  // Comment/description
}

// GetZoneByID retrieves a zone by its ID
func (s *ZoneSection) GetZoneByID(id int) (Zone, bool) {
	for _, zone := range s.Zones {
		if zone.No == id {
			return zone, true
		}
	}
	return Zone{}, false
}

// GetZoneByCode retrieves a zone by its code
func (s *ZoneSection) GetZoneByCode(code string) (Zone, bool) {
	for _, zone := range s.Zones {
		if zone.Code == code {
			return zone, true
		}
	}
	return Zone{}, false
}

// GetZonesByType retrieves all zones of a specified type
func (s *ZoneSection) GetZonesByType(typeNo int) []Zone {
	var result []Zone
	for _, zone := range s.Zones {
		if zone.TypeNo == typeNo {
			result = append(result, zone)
		}
	}
	return result
}

// GetZonesBySurface retrieves all zones associated with a particular surface
func (s *ZoneSection) GetZonesBySurface(surfaceID int) []Zone {
	var result []Zone
	for _, zone := range s.Zones {
		if zone.SurfaceID == surfaceID {
			result = append(result, zone)
		}
	}
	return result
}

// CalculateBoundingBox returns the min/max coordinates of all zones
func (s *ZoneSection) CalculateBoundingBox() (minX, minY, maxX, maxY float64) {
	if len(s.Zones) == 0 {
		return 0, 0, 0, 0
	}

	minX = s.Zones[0].XCoord
	minY = s.Zones[0].YCoord
	maxX = s.Zones[0].XCoord
	maxY = s.Zones[0].YCoord

	for _, zone := range s.Zones {
		if zone.XCoord < minX {
			minX = zone.XCoord
		}
		if zone.YCoord < minY {
			minY = zone.YCoord
		}
		if zone.XCoord > maxX {
			maxX = zone.XCoord
		}
		if zone.YCoord > maxY {
			maxY = zone.YCoord
		}
	}

	return minX, minY, maxX, maxY
}

// GetTotalPopulation calculates the total population across all zones
func (s *ZoneSection) GetTotalPopulation() int {
	total := 0
	for _, zone := range s.Zones {
		total += zone.Population
	}
	return total
}

// GetTotalEmployment calculates the total employment across all zones
func (s *ZoneSection) GetTotalEmployment() int {
	total := 0
	for _, zone := range s.Zones {
		total += zone.Employment
	}
	return total
}

// Count returns the number of zones in the section
func (s *ZoneSection) Count() int {
	return len(s.Zones)
}

// getZone extracts data from ZONE section row
func getZone(values []string) (Zone, error) {
	if len(values) < 7 {
		return Zone{}, fmt.Errorf("invalid ZONE data (insufficient fields): %v", values)
	}

	var zone Zone
	var err error

	// Parse NO (required field)
	if values[0] == "" {
		return Zone{}, fmt.Errorf("missing required field NO")
	}
	zone.No, err = strconv.Atoi(values[0])
	if err != nil {
		return Zone{}, fmt.Errorf("error parsing NO: %w", err)
	}

	// Parse CODE (optional)
	zone.Code = values[1]

	// Parse NAME (optional)
	zone.Name = values[2]

	// Parse MAINZONENO (optional)
	if values[3] != "" {
		zone.MainZoneNo, err = strconv.Atoi(values[3])
		if err != nil {
			return Zone{}, fmt.Errorf("error parsing MAINZONENO: %w", err)
		}
	}

	// Parse TYPENO (optional)
	if values[4] != "" {
		zone.TypeNo, err = strconv.Atoi(values[4])
		if err != nil {
			return Zone{}, fmt.Errorf("error parsing TYPENO: %w", err)
		}
	}

	// Parse XCOORD (required field)
	if values[5] == "" {
		return Zone{}, fmt.Errorf("missing required field XCOORD")
	}
	zone.XCoord, err = strconv.ParseFloat(strings.Replace(values[5], ",", ".", -1), 64)
	if err != nil {
		return Zone{}, fmt.Errorf("error parsing XCOORD: %w", err)
	}

	// Parse YCOORD (required field)
	if values[6] == "" {
		return Zone{}, fmt.Errorf("missing required field YCOORD")
	}
	zone.YCoord, err = strconv.ParseFloat(strings.Replace(values[6], ",", ".", -1), 64)
	if err != nil {
		return Zone{}, fmt.Errorf("error parsing YCOORD: %w", err)
	}

	// Parse SURFACEID (optional)
	if len(values) > 7 && values[7] != "" {
		zone.SurfaceID, err = strconv.Atoi(values[7])
		if err != nil {
			return Zone{}, fmt.Errorf("error parsing SURFACEID: %w", err)
		}
	}

	// Parse RELATIVESTATE (optional)
	if len(values) > 8 && values[8] != "" {
		zone.RelativeState, err = strconv.Atoi(values[8])
		if err != nil {
			return Zone{}, fmt.Errorf("error parsing RELATIVESTATE: %w", err)
		}
	}

	// Parse SHAREPRTORIG (optional)
	if len(values) > 9 && values[9] != "" {
		zone.SharePRTOrig, err = strconv.ParseFloat(strings.Replace(values[9], ",", ".", -1), 64)
		if err != nil {
			return Zone{}, fmt.Errorf("error parsing SHAREPRTORIG: %w", err)
		}
	}

	// Parse SHAREPRTDEST (optional)
	if len(values) > 10 && values[10] != "" {
		zone.SharePRTDest, err = strconv.ParseFloat(strings.Replace(values[10], ",", ".", -1), 64)
		if err != nil {
			return Zone{}, fmt.Errorf("error parsing SHAREPRTDEST: %w", err)
		}
	}

	// Parse SHAREPUT (optional)
	if len(values) > 11 && values[11] != "" {
		zone.SharePUT, err = strconv.ParseFloat(strings.Replace(values[11], ",", ".", -1), 64)
		if err != nil {
			return Zone{}, fmt.Errorf("error parsing SHAREPUT: %w", err)
		}
	}

	// Parse METHODCONNSHARES (optional)
	if len(values) > 12 && values[12] != "" {
		zone.MethodConnShares, err = strconv.Atoi(values[12])
		if err != nil {
			return Zone{}, fmt.Errorf("error parsing METHODCONNSHARES: %w", err)
		}
	}

	// Parse socioeconomic data - using your file indices
	// POPULATION
	if len(values) > 46 && values[46] != "" {
		zone.Population, err = strconv.Atoi(strings.Replace(values[46], ",", "", -1))
		if err != nil {
			// Just log warning and continue, as this is optional
			fmt.Printf("Warning: could not parse POPULATION for zone %d: %v\n", zone.No, err)
		}
	}

	// WORKPLACES
	if len(values) > 57 && values[57] != "" {
		zone.Employment, err = strconv.Atoi(strings.Replace(values[57], ",", "", -1))
		if err != nil {
			// Just log warning and continue, as this is optional
			fmt.Printf("Warning: could not parse WORKPLACES for zone %d: %v\n", zone.No, err)
		}
	}

	// WORKERS
	if len(values) > 56 && values[56] != "" {
		zone.Workers, err = strconv.Atoi(strings.Replace(values[56], ",", "", -1))
		if err != nil {
			// Just log warning and continue, as this is optional
			fmt.Printf("Warning: could not parse WORKERS for zone %d: %v\n", zone.No, err)
		}
	}

	// STUDENTS
	if len(values) > 49 && values[49] != "" {
		zone.Students, err = strconv.Atoi(strings.Replace(values[49], ",", "", -1))
		if err != nil {
			// Just log warning and continue, as this is optional
			fmt.Printf("Warning: could not parse STUDENTS for zone %d: %v\n", zone.No, err)
		}
	}

	// STUDYPLACES
	if len(values) > 50 && values[50] != "" {
		zone.StudyPlaces, err = strconv.Atoi(strings.Replace(values[50], ",", "", -1))
		if err != nil {
			// Just log warning and continue, as this is optional
			fmt.Printf("Warning: could not parse STUDYPLACES for zone %d: %v\n", zone.No, err)
		}
	}

	// POPDENS - Population density
	if len(values) > 45 && values[45] != "" {
		zone.PopDens, err = strconv.ParseFloat(strings.Replace(values[45], ",", ".", -1), 64)
		if err != nil {
			// Just log warning and continue, as this is optional
			fmt.Printf("Warning: could not parse POPDENS for zone %d: %v\n", zone.No, err)
		}
	}

	// COMMENT
	if len(values) > 29 && values[29] != "" {
		zone.Comment = values[29]
	}

	return zone, nil
}
