package ptvvisum

import "fmt"

// InfoSection represents $INFO section
type InfoSection struct {
	BaseSection
	Lines []InfoLine
}

type InfoLine struct {
	Index int
	Text  string
}

func getInfoLine(values []string) (InfoLine, error) {
	if len(values) < 2 {
		return InfoLine{}, fmt.Errorf("invalid INFO data: %v", values)
	}
	var index int
	if _, err := fmt.Sscanf(values[0], "%d", &index); err != nil {
		return InfoLine{}, fmt.Errorf("error parsing index: %w", err)
	}
	return InfoLine{
		Index: index,
		Text:  values[1],
	}, nil
}
