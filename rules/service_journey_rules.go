package rules

import "github.com/theoremus-urban-solutions/netex-validator/types"

// loadServiceJourneyRules loads comprehensive service journey and timetable validation rules
func (r *RuleRegistry) loadServiceJourneyRules() {
	// Basic service journey validation
	r.addBasicServiceJourneyRules()

	// Timetabled passing time validation
	r.addTimetablePassingTimeRules()

	// Service journey timing validation
	r.addServiceJourneyTimingRules()

	// Service journey reference validation
	r.addServiceJourneyReferenceRules()

	// Dated service journey validation
	r.addDatedServiceJourneyRules()

	// Journey pattern consistency validation
	r.addJourneyPatternConsistencyRules()
}

// addBasicServiceJourneyRules adds fundamental service journey validation
func (r *RuleRegistry) addBasicServiceJourneyRules() {
	// ServiceJourney must have JourneyPatternRef
	r.addRule("SERVICE_JOURNEY_MISSING_PATTERN_REF", "ServiceJourney missing JourneyPatternRef",
		"ServiceJourney must reference a JourneyPattern", types.ERROR,
		"//vehicleJourneys/*[self::ServiceJourney][not(JourneyPatternRef)]")

	// ServiceJourney must have LineRef (directly or via JourneyPattern)
	r.addRule("SERVICE_JOURNEY_MISSING_LINE_REF", "ServiceJourney missing LineRef",
		"ServiceJourney must reference a Line", types.ERROR,
		"//vehicleJourneys/*[self::ServiceJourney][not(LineRef) and not(JourneyPatternRef)]")

	// ServiceJourney must have OperatorRef
	r.addRule("SERVICE_JOURNEY_MISSING_OPERATOR_REF", "ServiceJourney missing OperatorRef",
		"ServiceJourney must reference an Operator", types.ERROR,
		"//vehicleJourneys/*[self::ServiceJourney][not(OperatorRef) and not(//lines/*[self::Line or self::FlexibleLine]/OperatorRef)]")

	// ServiceJourney must have either DayTypes or be referenced by DatedServiceJourney
	r.addRule("SERVICE_JOURNEY_MISSING_CALENDAR_REF", "ServiceJourney missing calendar reference",
		"ServiceJourney must have calendar information", types.ERROR,
		"//vehicleJourneys/*[self::ServiceJourney][not(dayTypes) and not(@id = //vehicleJourneys/DatedServiceJourney/ServiceJourneyRef/@ref)]")

	// ServiceJourney should have Name for clarity
	r.addRule("SERVICE_JOURNEY_MISSING_NAME", "ServiceJourney missing Name",
		"ServiceJourney should have a descriptive Name", types.WARNING,
		"//vehicleJourneys/*[self::ServiceJourney][not(Name) or normalize-space(Name)='']")
}

// addTimetablePassingTimeRules adds timetabled passing time validation
func (r *RuleRegistry) addTimetablePassingTimeRules() {
	// TimetabledPassingTime must have either ArrivalTime or DepartureTime
	r.addRule("TIMETABLED_PASSING_TIME_NO_TIMES", "TimetabledPassingTime missing times",
		"TimetabledPassingTime must have ArrivalTime or DepartureTime", types.ERROR,
		"//vehicleJourneys/*[self::ServiceJourney]/passingTimes/TimetabledPassingTime[not(ArrivalTime) and not(DepartureTime) and not(EarliestDepartureTime) and not(LatestArrivalTime)]")

	// First stop must have DepartureTime
	r.addRule("SERVICE_JOURNEY_FIRST_STOP_NO_DEPARTURE", "First stop missing departure time",
		"First TimetabledPassingTime must have DepartureTime", types.ERROR,
		"//vehicleJourneys/*[self::ServiceJourney]/passingTimes/TimetabledPassingTime[1][not(DepartureTime) and not(EarliestDepartureTime)]")

	// Last stop must have ArrivalTime
	r.addRule("SERVICE_JOURNEY_LAST_STOP_NO_ARRIVAL", "Last stop missing arrival time",
		"Last TimetabledPassingTime must have ArrivalTime", types.ERROR,
		"//vehicleJourneys/*[self::ServiceJourney]/passingTimes/TimetabledPassingTime[position()=last()][not(ArrivalTime) and not(LatestArrivalTime)]")

	// TimetabledPassingTime must have StopPointInJourneyPatternRef
	r.addRule("TIMETABLED_PASSING_TIME_NO_STOP_REF", "TimetabledPassingTime missing stop reference",
		"TimetabledPassingTime must reference a StopPointInJourneyPattern", types.ERROR,
		"//vehicleJourneys/*[self::ServiceJourney]/passingTimes/TimetabledPassingTime[not(StopPointInJourneyPatternRef)]")

	// Duplicate TimetabledPassingTime IDs within ServiceJourney
	r.addRule("SERVICE_JOURNEY_DUPLICATE_TPT_ID", "Duplicate TimetabledPassingTime ID",
		"TimetabledPassingTime IDs must be unique within ServiceJourney", types.ERROR,
		"//vehicleJourneys/*[self::ServiceJourney]/passingTimes/TimetabledPassingTime[@id = preceding-sibling::TimetabledPassingTime/@id or @id = following-sibling::TimetabledPassingTime/@id]")

	// TimetabledPassingTime should have version
	r.addRule("TIMETABLED_PASSING_TIME_MISSING_VERSION", "TimetabledPassingTime missing version",
		"TimetabledPassingTime should have version attribute", types.WARNING,
		"//vehicleJourneys/*[self::ServiceJourney]/passingTimes/TimetabledPassingTime[not(@version)]")
}

// addServiceJourneyTimingRules adds timing consistency validation
func (r *RuleRegistry) addServiceJourneyTimingRules() {
	// ArrivalTime should be before or equal to DepartureTime at same stop
	r.addRule("SERVICE_JOURNEY_ARRIVAL_AFTER_DEPARTURE", "Arrival time after departure time",
		"ArrivalTime should not be after DepartureTime at same stop", types.ERROR,
		"//vehicleJourneys/*[self::ServiceJourney]/passingTimes/TimetabledPassingTime[ArrivalTime and DepartureTime and ArrivalTime > DepartureTime]")

	// Times should progress monotonically through journey
	r.addRule("SERVICE_JOURNEY_TIME_REGRESSION", "Service journey time regression",
		"Times should progress monotonically through the journey", types.ERROR,
		"//vehicleJourneys/*[self::ServiceJourney]/passingTimes/TimetabledPassingTime[position() > 1 and ((DepartureTime and preceding-sibling::TimetabledPassingTime[1]/DepartureTime and DepartureTime <= preceding-sibling::TimetabledPassingTime[1]/DepartureTime) or (ArrivalTime and preceding-sibling::TimetabledPassingTime[1]/ArrivalTime and ArrivalTime <= preceding-sibling::TimetabledPassingTime[1]/ArrivalTime))]")

	// Warn about same arrival and departure times (except for terminals)
	r.addRule("SERVICE_JOURNEY_IDENTICAL_TIMES", "Identical arrival and departure times",
		"ArrivalTime equals DepartureTime - verify if stop is terminal", types.WARNING,
		"//vehicleJourneys/*[self::ServiceJourney]/passingTimes/TimetabledPassingTime[ArrivalTime = DepartureTime]")

	// Check for reasonable stop times (not too long or too short)
	r.addRule("SERVICE_JOURNEY_UNREASONABLE_STOP_TIME", "Unreasonable stop time",
		"Stop time appears unreasonable (too long or too short)", types.WARNING,
		"//vehicleJourneys/*[self::ServiceJourney]/passingTimes/TimetabledPassingTime[ArrivalTime and DepartureTime and (DepartureTime - ArrivalTime > 'PT30M' or DepartureTime - ArrivalTime < 'PT0S')]")

	// Validate time format (should be valid ISO 8601 time)
	r.addRule("SERVICE_JOURNEY_INVALID_TIME_FORMAT", "Invalid time format",
		"Time values should be valid ISO 8601 format", types.ERROR,
		"//vehicleJourneys/*[self::ServiceJourney]/passingTimes/TimetabledPassingTime/*[self::ArrivalTime or self::DepartureTime or self::EarliestDepartureTime or self::LatestArrivalTime][not(matches(text(), '^[0-9]{2}:[0-9]{2}:[0-9]{2}$'))]")
}

// addServiceJourneyReferenceRules adds reference consistency validation
func (r *RuleRegistry) addServiceJourneyReferenceRules() {
	// JourneyPatternRef should exist
	r.addRule("SERVICE_JOURNEY_INVALID_PATTERN_REF", "Invalid JourneyPatternRef",
		"JourneyPatternRef does not reference existing JourneyPattern", types.ERROR,
		"//vehicleJourneys/*[self::ServiceJourney]/JourneyPatternRef[@ref and not(@ref = //journeyPatterns/*[self::JourneyPattern or self::ServiceJourneyPattern]/@id)]")

	// LineRef should exist
	r.addRule("SERVICE_JOURNEY_INVALID_LINE_REF", "Invalid LineRef",
		"LineRef does not reference existing Line", types.ERROR,
		"//vehicleJourneys/*[self::ServiceJourney]/LineRef[@ref and not(@ref = //lines/*[self::Line or self::FlexibleLine]/@id)]")

	// OperatorRef should exist
	r.addRule("SERVICE_JOURNEY_INVALID_OPERATOR_REF", "Invalid OperatorRef",
		"OperatorRef does not reference existing Operator", types.ERROR,
		"//vehicleJourneys/*[self::ServiceJourney]/OperatorRef[@ref and not(@ref = //organisations/Operator/@id)]")

	// DayTypeRef should exist
	r.addRule("SERVICE_JOURNEY_INVALID_DAYTYPE_REF", "Invalid DayTypeRef",
		"DayTypeRef does not reference existing DayType", types.ERROR,
		"//vehicleJourneys/*[self::ServiceJourney]/dayTypes/DayTypeRef[@ref and not(@ref = //dayTypes/DayType/@id)]")

	// StopPointInJourneyPatternRef should exist
	r.addRule("TIMETABLED_PASSING_TIME_INVALID_STOP_REF", "Invalid StopPointInJourneyPatternRef",
		"StopPointInJourneyPatternRef does not reference existing StopPointInJourneyPattern", types.ERROR,
		"//vehicleJourneys/*[self::ServiceJourney]/passingTimes/TimetabledPassingTime/StopPointInJourneyPatternRef[@ref and not(@ref = //journeyPatterns/*/pointsInSequence/StopPointInJourneyPattern/@id)]")
}

// addDatedServiceJourneyRules adds dated service journey validation
func (r *RuleRegistry) addDatedServiceJourneyRules() {
	// DatedServiceJourney must have ServiceJourneyRef
	r.addRule("DATED_SERVICE_JOURNEY_NO_SERVICE_REF", "DatedServiceJourney missing ServiceJourneyRef",
		"DatedServiceJourney must reference a ServiceJourney", types.ERROR,
		"//vehicleJourneys/*[self::DatedServiceJourney][not(ServiceJourneyRef)]")

	// DatedServiceJourney must have OperatingDayRef
	r.addRule("DATED_SERVICE_JOURNEY_NO_OPERATING_DAY_REF", "DatedServiceJourney missing OperatingDayRef",
		"DatedServiceJourney must reference an OperatingDay", types.ERROR,
		"//vehicleJourneys/*[self::DatedServiceJourney][not(OperatingDayRef)]")

	// ServiceJourneyRef should be valid
	r.addRule("DATED_SERVICE_JOURNEY_INVALID_SERVICE_REF", "Invalid ServiceJourneyRef",
		"ServiceJourneyRef does not reference existing ServiceJourney", types.ERROR,
		"//vehicleJourneys/*[self::DatedServiceJourney]/ServiceJourneyRef[@ref and not(@ref = //vehicleJourneys/ServiceJourney/@id)]")

	// OperatingDayRef should be valid
	r.addRule("DATED_SERVICE_JOURNEY_INVALID_OPERATING_DAY_REF", "Invalid OperatingDayRef",
		"OperatingDayRef does not reference existing OperatingDay", types.ERROR,
		"//vehicleJourneys/*[self::DatedServiceJourney]/OperatingDayRef[@ref and not(@ref = //operatingDays/OperatingDay/@id)]")

	// DatedServiceJourney should not have both calendar references and OperatingDay
	r.addRule("DATED_SERVICE_JOURNEY_MIXED_CALENDAR", "Mixed calendar references",
		"DatedServiceJourney should not mix DayTypes with OperatingDay", types.WARNING,
		"//vehicleJourneys/*[self::DatedServiceJourney][dayTypes and OperatingDayRef]")

	// DatedServiceJourney should not override timetabled passing times
	r.addRule("DATED_SERVICE_JOURNEY_OVERRIDE_TIMES", "DatedServiceJourney overrides times",
		"DatedServiceJourney should not override ServiceJourney timetabled passing times", types.WARNING,
		"//vehicleJourneys/*[self::DatedServiceJourney][passingTimes]")
}

// addJourneyPatternConsistencyRules validates consistency between ServiceJourney and JourneyPattern
func (r *RuleRegistry) addJourneyPatternConsistencyRules() {
	// Number of TimetabledPassingTimes should match StopPointsInJourneyPattern
	r.addRule("SERVICE_JOURNEY_TPT_COUNT_MISMATCH", "Inconsistent passing time count",
		"Number of TimetabledPassingTimes should match JourneyPattern StopPoints", types.ERROR,
		"//vehicleJourneys/*[self::ServiceJourney][JourneyPatternRef and passingTimes and count(passingTimes/TimetabledPassingTime) != count(//journeyPatterns/*[@id=current()/JourneyPatternRef/@ref]/pointsInSequence/StopPointInJourneyPattern)]")

	// Order of TimetabledPassingTimes should match JourneyPattern order
	r.addRule("SERVICE_JOURNEY_TPT_ORDER_MISMATCH", "Inconsistent passing time order",
		"TimetabledPassingTime order should match JourneyPattern order", types.ERROR,
		"//vehicleJourneys/*[self::ServiceJourney]/passingTimes/TimetabledPassingTime[position() != //journeyPatterns/*[@id=current()/../../JourneyPatternRef/@ref]/pointsInSequence/StopPointInJourneyPattern[@id=current()/StopPointInJourneyPatternRef/@ref]/@order]")

	// ServiceJourney transport mode should be compatible with Line
	r.addRule("SERVICE_JOURNEY_INCOMPATIBLE_TRANSPORT_MODE", "Incompatible transport mode",
		"ServiceJourney transport mode incompatible with Line", types.ERROR,
		"//vehicleJourneys/*[self::ServiceJourney][TransportMode and LineRef and TransportMode != //lines/*[@id=current()/LineRef/@ref]/TransportMode]")

	// Block assignments should be consistent
	r.addRule("SERVICE_JOURNEY_INCONSISTENT_BLOCK", "Inconsistent block assignment",
		"ServiceJourney block assignment should be consistent", types.WARNING,
		"//vehicleJourneys/*[self::ServiceJourney][BlockRef and @id = //blocks/Block[@id != current()/BlockRef/@ref]/journeys/*/VehicleJourneyRef/@ref]")
}

// addServiceJourneyInterchangeRules adds service journey interchange validation
func (r *RuleRegistry) addServiceJourneyInterchangeRules() {
	// ServiceJourneyInterchange must have valid references
	r.addRule("INTERCHANGE_MISSING_FROM_JOURNEY_REF", "Interchange missing FromJourneyRef",
		"ServiceJourneyInterchange must have FromJourneyRef", types.ERROR,
		"//interchanges/ServiceJourneyInterchange[not(FromJourneyRef)]")

	r.addRule("INTERCHANGE_MISSING_TO_JOURNEY_REF", "Interchange missing ToJourneyRef",
		"ServiceJourneyInterchange must have ToJourneyRef", types.ERROR,
		"//interchanges/ServiceJourneyInterchange[not(ToJourneyRef)]")

	r.addRule("INTERCHANGE_MISSING_FROM_STOP_REF", "Interchange missing FromStopPointRef",
		"ServiceJourneyInterchange must have FromStopPointRef", types.ERROR,
		"//interchanges/ServiceJourneyInterchange[not(FromStopPointRef)]")

	r.addRule("INTERCHANGE_MISSING_TO_STOP_REF", "Interchange missing ToStopPointRef",
		"ServiceJourneyInterchange must have ToStopPointRef", types.ERROR,
		"//interchanges/ServiceJourneyInterchange[not(ToStopPointRef)]")

	// Interchange transfer time should be reasonable
	r.addRule("INTERCHANGE_UNREASONABLE_TRANSFER_TIME", "Unreasonable interchange transfer time",
		"Interchange transfer time should be reasonable", types.WARNING,
		"//interchanges/ServiceJourneyInterchange/StandardTransferTime[. > 'PT30M' or . < 'PT1M']")

	// Self-referencing interchanges should be flagged
	r.addRule("INTERCHANGE_SELF_REFERENCE", "Self-referencing interchange",
		"ServiceJourneyInterchange should not reference same journey", types.WARNING,
		"//interchanges/ServiceJourneyInterchange[FromJourneyRef/@ref = ToJourneyRef/@ref]")
}
