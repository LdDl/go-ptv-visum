package ptvvisum

import "fmt"

// UserAttDefSection represents $USERATTDEF section
type UserAttDefSection struct {
	BaseSection
	Attributes []UserAttDef
}

// UserAttDef represents a user-defined attribute
type UserAttDef struct {
	ObjID              string
	AttID              string
	Code               string
	Name               string
	ValueType          string
	MinValue           string
	MaxValue           string
	DefaultValue       string
	DefaultStringValue string
	Comment            string
	MaxStringLength    string
	NumDecPlaces       string
	DataSourceType     string
	Formula            string
	ScaledByLength     string
	CrossSectionLogic  string
	CSLIgnoreClosed    string
	SubAttrs           string
	CanBeEmpty         string
	OperationReference string
}

// getUserAttDef extracts data from USERATTDEF section row
func getUserAttDef(values []string) (UserAttDef, error) {
	// Check if we have at least the required fields
	if len(values) < 5 {
		return UserAttDef{}, fmt.Errorf("invalid USERATTDEF data: %v", values)
	}

	// Create attribute with initial required fields
	attr := UserAttDef{
		ObjID:     values[0],
		AttID:     values[1],
		Code:      values[2],
		Name:      values[3],
		ValueType: values[4],
	}

	// Add optional fields if available
	if len(values) > 5 {
		attr.MinValue = values[5]
	}
	if len(values) > 6 {
		attr.MaxValue = values[6]
	}
	if len(values) > 7 {
		attr.DefaultValue = values[7]
	}
	if len(values) > 8 {
		attr.DefaultStringValue = values[8]
	}
	if len(values) > 9 {
		attr.Comment = values[9]
	}
	if len(values) > 10 {
		attr.MaxStringLength = values[10]
	}
	if len(values) > 11 {
		attr.NumDecPlaces = values[11]
	}
	if len(values) > 12 {
		attr.DataSourceType = values[12]
	}
	if len(values) > 13 {
		attr.Formula = values[13]
	}
	if len(values) > 14 {
		attr.ScaledByLength = values[14]
	}
	if len(values) > 15 {
		attr.CrossSectionLogic = values[15]
	}
	if len(values) > 16 {
		attr.CSLIgnoreClosed = values[16]
	}
	if len(values) > 17 {
		attr.SubAttrs = values[17]
	}
	if len(values) > 18 {
		attr.CanBeEmpty = values[18]
	}
	if len(values) > 19 {
		attr.OperationReference = values[19]
	}

	return attr, nil
}
