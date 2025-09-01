package validator

import (
	"encoding/xml"
	"strings"
)

// NetEXStatistics provides counts of key NetEX elements similar to GTFS statistics
type NetEXStatistics struct {
	// Core transport elements
	OperatorCount       int `json:"operatorCount"`
	LineCount           int `json:"lineCount"`
	RouteCount          int `json:"routeCount"`
	ServiceJourneyCount int `json:"serviceJourneyCount"`
	StopPlaceCount      int `json:"stopPlaceCount"`
	StopPointCount      int `json:"stopPointCount"`

	// Timetabling
	TimetabledPassingTimeCount int `json:"timetabledPassingTimeCount"`
	DayTypeCount               int `json:"dayTypeCount"`
	ServiceCalendarCount       int `json:"serviceCalendarCount"`

	// Network structure
	RoutePointCount     int `json:"routePointCount"`
	RouteLinkCount      int `json:"routeLinkCount"`
	JourneyPatternCount int `json:"journeyPatternCount"`

	// Fare and accessibility
	FareZoneCount      int `json:"fareZoneCount"`
	TariffZoneCount    int `json:"tariffZoneCount"`
	AccessibilityCount int `json:"accessibilityCount"`

	// Infrastructure
	QuayCount     int `json:"quayCount"`
	EntranceCount int `json:"entranceCount"`
	PathLinkCount int `json:"pathLinkCount"`

	// Total counts
	TotalElements int `json:"totalElements"`
	TotalFiles    int `json:"totalFiles"`
}

// ExtractStatistics analyzes the XML content and extracts NetEX element statistics
func ExtractStatistics(xmlContent []byte, fileName string) NetEXStatistics {
	stats := NetEXStatistics{
		TotalFiles: 1,
	}

	if len(xmlContent) == 0 {
		return stats
	}

	// Parse XML to count elements
	decoder := xml.NewDecoder(strings.NewReader(string(xmlContent)))

	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}

		if startElement, ok := token.(xml.StartElement); ok {
			elementName := startElement.Name.Local
			stats.TotalElements++

			// Count NetEX-specific elements
			switch elementName {
			case "Operator", "Authority":
				stats.OperatorCount++
			case "Line":
				stats.LineCount++
			case "Route":
				stats.RouteCount++
			case "ServiceJourney", "VehicleJourney":
				stats.ServiceJourneyCount++
			case "StopPlace":
				stats.StopPlaceCount++
			case "StopPoint", "ScheduledStopPoint":
				stats.StopPointCount++
			case "TimetabledPassingTime":
				stats.TimetabledPassingTimeCount++
			case "DayType", "OperatingDay":
				stats.DayTypeCount++
			case "ServiceCalendar", "Calendar":
				stats.ServiceCalendarCount++
			case "RoutePoint":
				stats.RoutePointCount++
			case "RouteLink":
				stats.RouteLinkCount++
			case "JourneyPattern":
				stats.JourneyPatternCount++
			case "FareZone":
				stats.FareZoneCount++
			case "TariffZone":
				stats.TariffZoneCount++
			case "AccessibilityAssessment", "AccessibilityFeature":
				stats.AccessibilityCount++
			case "Quay":
				stats.QuayCount++
			case "Entrance", "StopPlaceEntrance":
				stats.EntranceCount++
			case "PathLink", "SitePathLink":
				stats.PathLinkCount++
			}
		}
	}

	return stats
}

// MergeStatistics combines statistics from multiple files
func MergeStatistics(stats []NetEXStatistics) NetEXStatistics {
	merged := NetEXStatistics{}

	for _, s := range stats {
		merged.OperatorCount += s.OperatorCount
		merged.LineCount += s.LineCount
		merged.RouteCount += s.RouteCount
		merged.ServiceJourneyCount += s.ServiceJourneyCount
		merged.StopPlaceCount += s.StopPlaceCount
		merged.StopPointCount += s.StopPointCount
		merged.TimetabledPassingTimeCount += s.TimetabledPassingTimeCount
		merged.DayTypeCount += s.DayTypeCount
		merged.ServiceCalendarCount += s.ServiceCalendarCount
		merged.RoutePointCount += s.RoutePointCount
		merged.RouteLinkCount += s.RouteLinkCount
		merged.JourneyPatternCount += s.JourneyPatternCount
		merged.FareZoneCount += s.FareZoneCount
		merged.TariffZoneCount += s.TariffZoneCount
		merged.AccessibilityCount += s.AccessibilityCount
		merged.QuayCount += s.QuayCount
		merged.EntranceCount += s.EntranceCount
		merged.PathLinkCount += s.PathLinkCount
		merged.TotalElements += s.TotalElements
		merged.TotalFiles += s.TotalFiles
	}

	return merged
}
