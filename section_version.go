package ptvvisum

import "fmt"

// VersionSection represents $VERSION section
type VersionSection struct {
	BaseSection
	Version  string
	FileType string
	Language string
	Unit     string
}

func getVersion(values []string) (version string, sileType string, language string, unit string, err error) {
	if len(values) < 4 {
		return "", "", "", "", fmt.Errorf("invalid VERSION data: %v", values)
	}
	version = values[0]
	sileType = values[1]
	language = values[2]
	unit = values[3]
	return
}
