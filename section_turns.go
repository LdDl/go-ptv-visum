package ptvvisum

import (
	"fmt"
	"strconv"
	"strings"
)

// TurnSection represents $TURN section
type TurnSection struct {
	BaseSection
	Turns []Turn
}

// Turn represents a single turning movement in the transportation network
type Turn struct {
	FromNodeNo                             int     // Origin node ID
	ViaNodeNo                              int     // Intersection node ID
	ToNodeNo                               int     // Destination node ID
	TypeNo                                 int     // Turn type ID (1=left, 2=right, 3=through, 4=U-turn)
	TSysSet                                string  // Transport systems allowed
	CapPRT                                 int     // Capacity for private transport
	T0PRT                                  string  // Default travel time
	AddVal                                 [3]int  // Additional values 1-3
	SBAPresetCriticalGap                   string  // SBA critical gap time
	SBAUsePresetCriticalGap                int     // Flag to use preset critical gap
	SBAPresetFollowupGap                   string  // SBA followup gap time
	SBAUsePresetFollowupGap                int     // Flag to use preset followup gap
	SBAPresetCriticalGapTurnOnRed          string  // Critical gap for turn on red
	SBAUsePresetCriticalGapTurnOnRed       int     // Flag to use preset critical gap for turn on red
	SBAPresetFollowupGapTurnOnRed          string  // Followup gap for turn on red
	SBAUsePresetFollowupGapTurnOnRed       int     // Flag to use preset followup gap for turn on red
	ICAUsePresetSatFlowRate                int     // Flag to use preset saturation flow rate
	ICAPresetSatFlowRate                   float64 // Preset saturation flow rate
	ICAUsePresetCriticalGap                int     // Flag to use preset critical gap for ICA
	ICAPresetCriticalGap                   string  // Preset critical gap for ICA
	ICAPresetCriticalGapStageOne           string  // Preset critical gap stage one
	ICAPresetCriticalGapStageTwo           string  // Preset critical gap stage two
	ICAUsePresetFollowupTime               int     // Flag to use preset followup time for ICA
	ICAPresetFollowupTime                  string  // Preset followup time for ICA
	ICATurningRadius                       string  // Turning radius
	ICAUsePresentSatFlowAdjustment         int     // Flag to use preset saturation flow adjustment
	ICAPresetSatFlowAdjustment             float64 // Preset saturation flow adjustment
	ICAProtectedInnerSatFlowAdjustment     float64 // Protected inner saturation flow adjustment
	ICAUsePermissiveInnerSatFlowAdjustment int     // Flag to use permissive inner saturation flow adjustment
	ICAPermissiveInnerSatFlowAdjustment    float64 // Permissive inner saturation flow adjustment
	ICAUsePedestrianSatFlowAdjustment      int     // Flag to use pedestrian saturation flow adjustment
	ICAPedestrianSatFlowAdjustment         float64 // Pedestrian saturation flow adjustment
	ICAUsePresetLaneWidthAdjustment        int     // Flag to use preset lane width adjustment
	ICAPresetLaneWidthAdjustment           float64 // Preset lane width adjustment
	ICAUsePresetGradeAdjustment            int     // Flag to use preset grade adjustment
	ICAPresetGradeAdjustment               float64 // Preset grade adjustment
	ICAUsePresetTurningRadiusAdjustment    int     // Flag to use preset turning radius adjustment
	ICAPresetTurningRadiusAdjustment       float64 // Preset turning radius adjustment
	ICAUpstreamAdj                         float64 // Upstream adjustment factor
	ICAPHFVolAdj                           float64 // Peak hour factor volume adjustment
	ICAUnsignalizedDelay                   string  // Unsignalized delay
	AuxiliarySG                            string  // Auxiliary signal group
	IsChangeOfDirection                    int     // Flag indicating change of direction
	VISTROBaseVolInput                     int     // VISTRO base volume input
	VISTROBaseVolAdjustFactor              float64 // VISTRO base volume adjustment factor
	ShareHGV                               float64 // Share of heavy goods vehicles
	VISTROGrowthFactor                     float64 // VISTRO growth factor
	VISTROInProcessVol                     int     // VISTRO in-process volume
	VISTRODivTrips                         int     // VISTRO diverted trips
	VISTROPassByTrips                      int     // VISTRO pass-by trips
	VISTROSiteAdjustVol                    int     // VISTRO site adjustment volume
	VISTROOtherVol                         int     // VISTRO other volume
	VISTRORightTurnOnRedVol                int     // VISTRO right turn on red volume
	VISTROTurnOnRedPercentage              float64 // VISTRO turn on red percentage
	VISTROTurnOnRedVolumeCalculationMethod string  // VISTRO turn on red volume calculation method
	VISTROLRORderNo                        int     // VISTRO left/right order number
	VISTROOtherAdjustFactor                float64 // VISTRO other adjustment factor
	VISTROLaneWidth                        string  // VISTRO lane width
	UseVISTROLaneWidth                     int     // Flag to use VISTRO lane width
	VISTROOuterControl                     string  // VISTRO outer control type
	VISTROThruControl                      string  // VISTRO through control type
	VISTROInnerControl                     string  // VISTRO inner control type
	VISTROSGNo                             int     // VISTRO signal group number
	VISTROOVLNo                            int     // VISTRO overlap number
}

// GetTurnsByIntersection retrieves all turns at a specified intersection node
func (s *TurnSection) GetTurnsByIntersection(nodeNo int) []Turn {
	var result []Turn
	for _, turn := range s.Turns {
		if turn.ViaNodeNo == nodeNo {
			result = append(result, turn)
		}
	}
	return result
}

// GetTurnsByOrigin retrieves all turns from a specified origin node
func (s *TurnSection) GetTurnsByOrigin(nodeNo int) []Turn {
	var result []Turn
	for _, turn := range s.Turns {
		if turn.FromNodeNo == nodeNo {
			result = append(result, turn)
		}
	}
	return result
}

// GetTurnsByDestination retrieves all turns to a specified destination node
func (s *TurnSection) GetTurnsByDestination(nodeNo int) []Turn {
	var result []Turn
	for _, turn := range s.Turns {
		if turn.ToNodeNo == nodeNo {
			result = append(result, turn)
		}
	}
	return result
}

// GetTurnsByType retrieves all turns of a specified type
func (s *TurnSection) GetTurnsByType(typeNo int) []Turn {
	var result []Turn
	for _, turn := range s.Turns {
		if turn.TypeNo == typeNo {
			result = append(result, turn)
		}
	}
	return result
}

// GetTurn retrieves a specific turn by its nodes
func (s *TurnSection) GetTurn(fromNodeNo, viaNodeNo, toNodeNo int) (Turn, bool) {
	for _, turn := range s.Turns {
		if turn.FromNodeNo == fromNodeNo && turn.ViaNodeNo == viaNodeNo && turn.ToNodeNo == toNodeNo {
			return turn, true
		}
	}
	return Turn{}, false
}

// GetTurnsByTransportSystem retrieves all turns allowing a specific transport system
func (s *TurnSection) GetTurnsByTransportSystem(tsys string) []Turn {
	var result []Turn
	for _, turn := range s.Turns {
		systems := strings.Split(turn.TSysSet, ",")
		for _, system := range systems {
			if system == tsys {
				result = append(result, turn)
				break
			}
		}
	}
	return result
}

// GetLeftTurns retrieves all left turns (assuming TypeNo=1 is left turn)
func (s *TurnSection) GetLeftTurns() []Turn {
	return s.GetTurnsByType(1)
}

// GetRightTurns retrieves all right turns (assuming TypeNo=2 is right turn)
func (s *TurnSection) GetRightTurns() []Turn {
	return s.GetTurnsByType(2)
}

// GetThroughTurns retrieves all through movements (assuming TypeNo=3 is through)
func (s *TurnSection) GetThroughTurns() []Turn {
	return s.GetTurnsByType(3)
}

// GetUTurns retrieves all U-turns (assuming TypeNo=4 is U-turn)
func (s *TurnSection) GetUTurns() []Turn {
	return s.GetTurnsByType(4)
}

// GetTravelTime extracts the numeric travel time from the T0PRT string (e.g., "10s" -> 10)
func (t *Turn) GetTravelTime() float64 {
	if t.T0PRT == "" {
		return 0
	}

	// Extract numeric part
	numStr := t.T0PRT

	// Find where the numeric part ends
	for i, c := range t.T0PRT {
		if !isDigit(c) && c != '.' && c != ',' {
			numStr = t.T0PRT[:i]
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

// Count returns the number of turns in the section
func (s *TurnSection) Count() int {
	return len(s.Turns)
}

// CountByIntersection returns a map of turn counts by intersection node
func (s *TurnSection) CountByIntersection() map[int]int {
	counts := make(map[int]int)
	for _, turn := range s.Turns {
		counts[turn.ViaNodeNo]++
	}
	return counts
}

// getTurn extracts data from TURN section row
func getTurn(values []string) (Turn, error) {
	if len(values) < 7 {
		return Turn{}, fmt.Errorf("invalid TURN data (insufficient fields): %v", values)
	}

	var turn Turn
	var err error

	// Parse FROMNODENO (required field)
	if values[0] == "" {
		return Turn{}, fmt.Errorf("missing required field FROMNODENO")
	}
	turn.FromNodeNo, err = strconv.Atoi(values[0])
	if err != nil {
		return Turn{}, fmt.Errorf("error parsing FROMNODENO: %w", err)
	}

	// Parse VIANODENO (required field)
	if values[1] == "" {
		return Turn{}, fmt.Errorf("missing required field VIANODENO")
	}
	turn.ViaNodeNo, err = strconv.Atoi(values[1])
	if err != nil {
		return Turn{}, fmt.Errorf("error parsing VIANODENO: %w", err)
	}

	// Parse TONODENO (required field)
	if values[2] == "" {
		return Turn{}, fmt.Errorf("missing required field TONODENO")
	}
	turn.ToNodeNo, err = strconv.Atoi(values[2])
	if err != nil {
		return Turn{}, fmt.Errorf("error parsing TONODENO: %w", err)
	}

	// Parse TYPENO (required field)
	if values[3] == "" {
		return Turn{}, fmt.Errorf("missing required field TYPENO")
	}
	turn.TypeNo, err = strconv.Atoi(values[3])
	if err != nil {
		return Turn{}, fmt.Errorf("error parsing TYPENO: %w", err)
	}

	// Parse TSYSSET (optional)
	turn.TSysSet = values[4]

	// Parse CAPPRT (required field)
	if values[5] != "" {
		turn.CapPRT, err = strconv.Atoi(values[5])
		if err != nil {
			return Turn{}, fmt.Errorf("error parsing CAPPRT: %w", err)
		}
	}

	// Parse T0PRT (required field)
	turn.T0PRT = values[6]

	// Parse ADDVAL1, ADDVAL2, ADDVAL3 (if available)
	if len(values) > 7 && values[7] != "" {
		turn.AddVal[0], err = strconv.Atoi(values[7])
		if err != nil {
			return Turn{}, fmt.Errorf("error parsing ADDVAL1: %w", err)
		}
	}

	if len(values) > 8 && values[8] != "" {
		turn.AddVal[1], err = strconv.Atoi(values[8])
		if err != nil {
			return Turn{}, fmt.Errorf("error parsing ADDVAL2: %w", err)
		}
	}

	if len(values) > 9 && values[9] != "" {
		turn.AddVal[2], err = strconv.Atoi(values[9])
		if err != nil {
			return Turn{}, fmt.Errorf("error parsing ADDVAL3: %w", err)
		}
	}

	// Parse SBA fields (if available)
	if len(values) > 10 {
		turn.SBAPresetCriticalGap = values[10]
	}

	if len(values) > 11 && values[11] != "" {
		turn.SBAUsePresetCriticalGap, err = strconv.Atoi(values[11])
		if err != nil {
			return Turn{}, fmt.Errorf("error parsing SBAUSEPRESETCRITICALGAP: %w", err)
		}
	}

	if len(values) > 12 {
		turn.SBAPresetFollowupGap = values[12]
	}

	if len(values) > 13 && values[13] != "" {
		turn.SBAUsePresetFollowupGap, err = strconv.Atoi(values[13])
		if err != nil {
			return Turn{}, fmt.Errorf("error parsing SBAUSEPRESETFOLLOWUPGAP: %w", err)
		}
	}

	if len(values) > 14 {
		turn.SBAPresetCriticalGapTurnOnRed = values[14]
	}

	if len(values) > 15 && values[15] != "" {
		turn.SBAUsePresetCriticalGapTurnOnRed, err = strconv.Atoi(values[15])
		if err != nil {
			return Turn{}, fmt.Errorf("error parsing SBAUSEPRESETCRITICALGAPTURNONRED: %w", err)
		}
	}

	if len(values) > 16 {
		turn.SBAPresetFollowupGapTurnOnRed = values[16]
	}

	if len(values) > 17 && values[17] != "" {
		turn.SBAUsePresetFollowupGapTurnOnRed, err = strconv.Atoi(values[17])
		if err != nil {
			return Turn{}, fmt.Errorf("error parsing SBAUSEPRESETFOLLOWUPGAPTURNONRED: %w", err)
		}
	}

	// Parse ICA fields (if available)
	if len(values) > 18 && values[18] != "" {
		turn.ICAUsePresetSatFlowRate, err = strconv.Atoi(values[18])
		if err != nil {
			return Turn{}, fmt.Errorf("error parsing ICAUSEPRESETSATFLOWRATE: %w", err)
		}
	}

	if len(values) > 19 && values[19] != "" {
		turn.ICAPresetSatFlowRate, err = strconv.ParseFloat(strings.Replace(values[19], ",", ".", -1), 64)
		if err != nil {
			return Turn{}, fmt.Errorf("error parsing ICAPRESETSATFLOWRATE: %w", err)
		}
	}

	if len(values) > 20 && values[20] != "" {
		turn.ICAUsePresetCriticalGap, err = strconv.Atoi(values[20])
		if err != nil {
			return Turn{}, fmt.Errorf("error parsing ICAUSEPRESETCRITICALGAP: %w", err)
		}
	}

	if len(values) > 21 {
		turn.ICAPresetCriticalGap = values[21]
	}

	if len(values) > 22 {
		turn.ICAPresetCriticalGapStageOne = values[22]
	}

	if len(values) > 23 {
		turn.ICAPresetCriticalGapStageTwo = values[23]
	}

	if len(values) > 24 && values[24] != "" {
		turn.ICAUsePresetFollowupTime, err = strconv.Atoi(values[24])
		if err != nil {
			return Turn{}, fmt.Errorf("error parsing ICAUSEPRESETFOLLOWUPTIME: %w", err)
		}
	}

	if len(values) > 25 {
		turn.ICAPresetFollowupTime = values[25]
	}

	if len(values) > 26 {
		turn.ICATurningRadius = values[26]
	}

	// Parse additional ICA fields (if available)
	if len(values) > 27 && values[27] != "" {
		turn.ICAUsePresentSatFlowAdjustment, err = strconv.Atoi(values[27])
		if err != nil {
			return Turn{}, fmt.Errorf("error parsing ICAUSEPRESETSATFLOWADJUSTMENT: %w", err)
		}
	}

	if len(values) > 28 && values[28] != "" {
		turn.ICAPresetSatFlowAdjustment, err = strconv.ParseFloat(strings.Replace(values[28], ",", ".", -1), 64)
		if err != nil {
			return Turn{}, fmt.Errorf("error parsing ICAPRESETSATFLOWADJUSTMENT: %w", err)
		}
	}

	if len(values) > 29 && values[29] != "" {
		turn.ICAProtectedInnerSatFlowAdjustment, err = strconv.ParseFloat(strings.Replace(values[29], ",", ".", -1), 64)
		if err != nil {
			return Turn{}, fmt.Errorf("error parsing ICAPROTECTEDINNERSATFLOWADJUSTMENT: %w", err)
		}
	}

	if len(values) > 30 && values[30] != "" {
		turn.ICAUsePermissiveInnerSatFlowAdjustment, err = strconv.Atoi(values[30])
		if err != nil {
			return Turn{}, fmt.Errorf("error parsing ICAUSEPERMISSIVEINNERSATFLOWADJUSTMENT: %w", err)
		}
	}

	if len(values) > 31 && values[31] != "" {
		turn.ICAPermissiveInnerSatFlowAdjustment, err = strconv.ParseFloat(strings.Replace(values[31], ",", ".", -1), 64)
		if err != nil {
			return Turn{}, fmt.Errorf("error parsing ICAPERMISSIVEINNERSATFLOWADJUSTMENT: %w", err)
		}
	}

	if len(values) > 32 && values[32] != "" {
		turn.ICAUsePedestrianSatFlowAdjustment, err = strconv.Atoi(values[32])
		if err != nil {
			return Turn{}, fmt.Errorf("error parsing ICAUSEPEDESTRIANSATFLOWADJUSTMENT: %w", err)
		}
	}

	if len(values) > 33 && values[33] != "" {
		turn.ICAPedestrianSatFlowAdjustment, err = strconv.ParseFloat(strings.Replace(values[33], ",", ".", -1), 64)
		if err != nil {
			return Turn{}, fmt.Errorf("error parsing ICAPEDESTRIANSATFLOWADJUSTMENT: %w", err)
		}
	}

	if len(values) > 34 && values[34] != "" {
		turn.ICAUsePresetLaneWidthAdjustment, err = strconv.Atoi(values[34])
		if err != nil {
			return Turn{}, fmt.Errorf("error parsing ICAUSEPRESETLANEWIDTHADJUSTMENT: %w", err)
		}
	}

	if len(values) > 35 && values[35] != "" {
		turn.ICAPresetLaneWidthAdjustment, err = strconv.ParseFloat(strings.Replace(values[35], ",", ".", -1), 64)
		if err != nil {
			return Turn{}, fmt.Errorf("error parsing ICAPRESETLANEWIDTHADJUSTMENT: %w", err)
		}
	}

	if len(values) > 36 && values[36] != "" {
		turn.ICAUsePresetGradeAdjustment, err = strconv.Atoi(values[36])
		if err != nil {
			return Turn{}, fmt.Errorf("error parsing ICAUSEPRESETGRADEADJUSTMENT: %w", err)
		}
	}

	if len(values) > 37 && values[37] != "" {
		turn.ICAPresetGradeAdjustment, err = strconv.ParseFloat(strings.Replace(values[37], ",", ".", -1), 64)
		if err != nil {
			return Turn{}, fmt.Errorf("error parsing ICAPRESETGRADEADJUSTMENT: %w", err)
		}
	}

	if len(values) > 38 && values[38] != "" {
		turn.ICAUsePresetTurningRadiusAdjustment, err = strconv.Atoi(values[38])
		if err != nil {
			return Turn{}, fmt.Errorf("error parsing ICAUSEPRESETTURNINGRADIUSADJUSTMENT: %w", err)
		}
	}

	if len(values) > 39 && values[39] != "" {
		turn.ICAPresetTurningRadiusAdjustment, err = strconv.ParseFloat(strings.Replace(values[39], ",", ".", -1), 64)
		if err != nil {
			return Turn{}, fmt.Errorf("error parsing ICAPRESETTURNINGRADIUSADJUSTMENT: %w", err)
		}
	}

	if len(values) > 40 && values[40] != "" {
		turn.ICAUpstreamAdj, err = strconv.ParseFloat(strings.Replace(values[40], ",", ".", -1), 64)
		if err != nil {
			return Turn{}, fmt.Errorf("error parsing ICAUPSTREAMADJ: %w", err)
		}
	}

	if len(values) > 41 && values[41] != "" {
		turn.ICAPHFVolAdj, err = strconv.ParseFloat(strings.Replace(values[41], ",", ".", -1), 64)
		if err != nil {
			return Turn{}, fmt.Errorf("error parsing ICAPHFVOLADJ: %w", err)
		}
	}

	if len(values) > 42 {
		turn.ICAUnsignalizedDelay = values[42]
	}

	// Parse remaining fields
	if len(values) > 43 {
		turn.AuxiliarySG = values[43]
	}

	if len(values) > 44 && values[44] != "" {
		turn.IsChangeOfDirection, err = strconv.Atoi(values[44])
		if err != nil {
			return Turn{}, fmt.Errorf("error parsing ISCHANGEOFDIRECTION: %w", err)
		}
	}

	// Parse VISTRO fields
	if len(values) > 45 && values[45] != "" {
		turn.VISTROBaseVolInput, err = strconv.Atoi(values[45])
		if err != nil {
			return Turn{}, fmt.Errorf("error parsing VISTROBASEVOLINPUT: %w", err)
		}
	}

	if len(values) > 46 && values[46] != "" {
		turn.VISTROBaseVolAdjustFactor, err = strconv.ParseFloat(strings.Replace(values[46], ",", ".", -1), 64)
		if err != nil {
			return Turn{}, fmt.Errorf("error parsing VISTROBASEVOLADJUSTFACTOR: %w", err)
		}
	}

	if len(values) > 47 && values[47] != "" {
		turn.ShareHGV, err = strconv.ParseFloat(strings.Replace(values[47], ",", ".", -1), 64)
		if err != nil {
			return Turn{}, fmt.Errorf("error parsing SHAREHGV: %w", err)
		}
	}

	if len(values) > 48 && values[48] != "" {
		turn.VISTROGrowthFactor, err = strconv.ParseFloat(strings.Replace(values[48], ",", ".", -1), 64)
		if err != nil {
			return Turn{}, fmt.Errorf("error parsing VISTROGROWTHFACTOR: %w", err)
		}
	}

	if len(values) > 49 && values[49] != "" {
		turn.VISTROInProcessVol, err = strconv.Atoi(values[49])
		if err != nil {
			return Turn{}, fmt.Errorf("error parsing VISTROINPROCESSVOL: %w", err)
		}
	}

	if len(values) > 50 && values[50] != "" {
		turn.VISTRODivTrips, err = strconv.Atoi(values[50])
		if err != nil {
			return Turn{}, fmt.Errorf("error parsing VISTRODIVTRIPS: %w", err)
		}
	}

	if len(values) > 51 && values[51] != "" {
		turn.VISTROPassByTrips, err = strconv.Atoi(values[51])
		if err != nil {
			return Turn{}, fmt.Errorf("error parsing VISTROPASSBYTRIPS: %w", err)
		}
	}

	if len(values) > 52 && values[52] != "" {
		turn.VISTROSiteAdjustVol, err = strconv.Atoi(values[52])
		if err != nil {
			return Turn{}, fmt.Errorf("error parsing VISTROSITEADJUSTVOL: %w", err)
		}
	}

	if len(values) > 53 && values[53] != "" {
		turn.VISTROOtherVol, err = strconv.Atoi(values[53])
		if err != nil {
			return Turn{}, fmt.Errorf("error parsing VISTROOTHERVOL: %w", err)
		}
	}

	if len(values) > 54 && values[54] != "" {
		turn.VISTRORightTurnOnRedVol, err = strconv.Atoi(values[54])
		if err != nil {
			return Turn{}, fmt.Errorf("error parsing VISTRORIGHTTURNONREDVOL: %w", err)
		}
	}

	if len(values) > 55 && values[55] != "" {
		turn.VISTROTurnOnRedPercentage, err = strconv.ParseFloat(strings.Replace(values[55], ",", ".", -1), 64)
		if err != nil {
			return Turn{}, fmt.Errorf("error parsing VISTROTURNONREDPERCENTAGE: %w", err)
		}
	}

	if len(values) > 56 {
		turn.VISTROTurnOnRedVolumeCalculationMethod = values[56]
	}

	if len(values) > 57 && values[57] != "" {
		turn.VISTROLRORderNo, err = strconv.Atoi(values[57])
		if err != nil {
			return Turn{}, fmt.Errorf("error parsing VISTROLRORDERNO: %w", err)
		}
	}

	if len(values) > 58 && values[58] != "" {
		turn.VISTROOtherAdjustFactor, err = strconv.ParseFloat(strings.Replace(values[58], ",", ".", -1), 64)
		if err != nil {
			return Turn{}, fmt.Errorf("error parsing VISTROOTHERADJUSTFACTOR: %w", err)
		}
	}

	if len(values) > 59 {
		turn.VISTROLaneWidth = values[59]
	}

	if len(values) > 60 && values[60] != "" {
		turn.UseVISTROLaneWidth, err = strconv.Atoi(values[60])
		if err != nil {
			return Turn{}, fmt.Errorf("error parsing USEVISTROLANEWIDTH: %w", err)
		}
	}

	if len(values) > 61 {
		turn.VISTROOuterControl = values[61]
	}

	if len(values) > 62 {
		turn.VISTROThruControl = values[62]
	}

	if len(values) > 63 {
		turn.VISTROInnerControl = values[63]
	}

	if len(values) > 64 && values[64] != "" {
		turn.VISTROSGNo, err = strconv.Atoi(values[64])
		if err != nil {
			return Turn{}, fmt.Errorf("error parsing VISTROSGNO: %w", err)
		}
	}

	if len(values) > 65 && values[65] != "" {
		turn.VISTROOVLNo, err = strconv.Atoi(values[65])
		if err != nil {
			return Turn{}, fmt.Errorf("error parsing VISTROOVLNO: %w", err)
		}
	}

	return turn, nil
}
