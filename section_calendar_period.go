package ptvvisum

import (
	"fmt"
	"time"
)

// CalendarPeriodSection represents $CALENDARPERIOD section
type CalendarPeriodSection struct {
	BaseSection
	Periods []CalendarPeriod
}

// CalendarPeriod represents a single calendar period
type CalendarPeriod struct {
	Type                        string
	ValidFrom                   time.Time
	ValidUntil                  time.Time
	AnalysisPeriodStartDayIndex int
	AnalysisPeriodEndDayIndex   int
	AnalysisTimeIntervalSetNo   int
}

// getCalendarPeriod extracts data from CALENDARPERIOD section row
func getCalendarPeriod(values []string) (CalendarPeriod, error) {
	if len(values) < 5 {
		return CalendarPeriod{}, fmt.Errorf("invalid CALENDARPERIOD data: %v", values)
	}

	// Initialize with default values
	period := CalendarPeriod{
		Type: values[0],
	}

	// Parse dates (format: DD.MM.YYYY)
	if values[1] != "" {
		validFrom, err := time.Parse("02.01.2006", values[1])
		if err != nil {
			return CalendarPeriod{}, fmt.Errorf("error parsing ValidFrom date: %w", err)
		}
		period.ValidFrom = validFrom
	}

	if values[2] != "" {
		validUntil, err := time.Parse("02.01.2006", values[2])
		if err != nil {
			return CalendarPeriod{}, fmt.Errorf("error parsing ValidUntil date: %w", err)
		}
		period.ValidUntil = validUntil
	}

	// Parse integer values
	if values[3] != "" {
		if _, err := fmt.Sscanf(values[3], "%d", &period.AnalysisPeriodStartDayIndex); err != nil {
			return CalendarPeriod{}, fmt.Errorf("error parsing AnalysisPeriodStartDayIndex: %w", err)
		}
	}

	if values[4] != "" {
		if _, err := fmt.Sscanf(values[4], "%d", &period.AnalysisPeriodEndDayIndex); err != nil {
			return CalendarPeriod{}, fmt.Errorf("error parsing AnalysisPeriodEndDayIndex: %w", err)
		}
	}

	// Optional field
	if len(values) > 5 && values[5] != "" {
		if _, err := fmt.Sscanf(values[5], "%d", &period.AnalysisTimeIntervalSetNo); err != nil {
			return CalendarPeriod{}, fmt.Errorf("error parsing AnalysisTimeIntervalSetNo: %w", err)
		}
	}

	return period, nil
}
