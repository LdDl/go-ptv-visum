package ptvvisum

import "fmt"

// POICategorySection represents $POICATEGORY section
type POICategorySection struct {
	BaseSection
	Categories []POICategory
}

type POICategory struct {
	No          int
	Code        string
	Name        string
	Comment     string
	ParentCatNo int
}

func getPoiCategory(values []string) (POICategory, error) {
	if len(values) < 5 {
		return POICategory{}, fmt.Errorf("invalid POICATEGORY data: %v", values)
	}
	var no int
	if _, err := fmt.Sscanf(values[0], "%d", &no); err != nil {
		return POICategory{}, fmt.Errorf("error parsing No: %w", err)
	}
	parentCatNo := 0
	if values[4] != "" {
		if _, err := fmt.Sscanf(values[4], "%d", &parentCatNo); err != nil {
			return POICategory{}, fmt.Errorf("error parsing ParentCatNo: %w", err)
		}
	}
	return POICategory{
		No:          no,
		Code:        values[1],
		Name:        values[2],
		Comment:     values[3],
		ParentCatNo: parentCatNo,
	}, nil
}
