package rules

import "github.com/theoremus-urban-solutions/netex-validator/types"

// loadXPathBusinessRules loads additional XPath-based business validation rules
func (r *RuleRegistry) loadXPathBusinessRules() {
	// Advanced mandatory field validation
	r.addMandatoryFieldRules()

	// Advanced structural validation
	r.addStructuralValidationRules()

	// Advanced data consistency rules
	r.addDataConsistencyRules()

	// EU profile specific rules
	r.addEUProfileRules()
}

// addMandatoryFieldRules adds comprehensive mandatory field validation
func (r *RuleRegistry) addMandatoryFieldRules() {
	// PublicationDelivery must have correct structure
	r.addRule("PUBLICATION_DELIVERY_STRUCTURE", "Invalid PublicationDelivery structure",
		"PublicationDelivery must have proper structure", types.ERROR,
		"/PublicationDelivery[not(PublicationTimestamp) or not(ParticipantRef) or not(dataObjects)]")

	// Frame validation - frames must have proper structure
	r.addRule("COMPOSITE_FRAME_STRUCTURE", "Invalid CompositeFrame structure",
		"CompositeFrame must contain required child frames", types.ERROR,
		"//CompositeFrame[not(frames)]")

	r.addRule("SERVICE_FRAME_STRUCTURE", "Invalid ServiceFrame structure",
		"ServiceFrame must contain service-related elements", types.WARNING,
		"//ServiceFrame[not(lines) and not(routes) and not(journeyPatterns) and not(vehicleJourneys) and not(stopAssignments)]")

	r.addRule("RESOURCE_FRAME_STRUCTURE", "Invalid ResourceFrame structure",
		"ResourceFrame must contain resource elements", types.WARNING,
		"//ResourceFrame[not(organisations) and not(operationalContexts) and not(vehicleTypes) and not(equipments)]")

	// Mandatory names on key elements
	r.addRule("AUTHORITY_MISSING_NAME", "Authority missing Name",
		"Authority must have Name", types.ERROR,
		"//organisations/Authority[not(Name) or normalize-space(Name) = '']")

	r.addRule("ROUTE_MISSING_NAME", "Route missing Name",
		"Route should have descriptive Name", types.WARNING,
		"//routes/Route[not(Name) or normalize-space(Name) = '']")

	r.addRule("JOURNEY_PATTERN_MISSING_NAME", "JourneyPattern missing Name",
		"JourneyPattern should have descriptive Name", types.WARNING,
		"//journeyPatterns/*[self::JourneyPattern or self::ServiceJourneyPattern][not(Name) or normalize-space(Name) = '']")

	// Mandatory references
	r.addRule("ROUTE_MISSING_DIRECTION", "Route missing DirectionType",
		"Route should specify DirectionType", types.WARNING,
		"//routes/Route[not(DirectionType)]")

	r.addRule("STOP_ASSIGNMENT_MISSING_REFS", "StopAssignment missing references",
		"PassengerStopAssignment must reference both ScheduledStopPoint and StopPlace/Quay", types.ERROR,
		"//stopAssignments/PassengerStopAssignment[not(ScheduledStopPointRef) or (not(StopPlaceRef) and not(QuayRef))]")
}

// addStructuralValidationRules adds structural validation beyond basic XSD
func (r *RuleRegistry) addStructuralValidationRules() {
	// Frame hierarchy validation
	r.addRule("INVALID_FRAME_NESTING", "Invalid frame nesting",
		"Frames should not contain other frames of same type", types.ERROR,
		"//ServiceFrame[.//ServiceFrame] | //ResourceFrame[.//ResourceFrame] | //TimetableFrame[.//TimetableFrame]")

	// Element order validation
	r.addRule("ROUTE_POINT_ORDER_MISSING", "RoutePoint missing order",
		"PointOnRoute must have order attribute", types.ERROR,
		"//routes/Route/pointsInSequence/PointOnRoute[not(@order)]")

	r.addRule("STOP_POINT_ORDER_MISSING", "StopPoint missing order",
		"StopPointInJourneyPattern must have order attribute", types.ERROR,
		"//journeyPatterns/*/pointsInSequence/StopPointInJourneyPattern[not(@order)]")

	// Duplicate order values
	r.addRule("ROUTE_DUPLICATE_ORDER", "Duplicate order in Route",
		"PointOnRoute order values must be unique within Route", types.ERROR,
		"//routes/Route/pointsInSequence/PointOnRoute[@order = preceding-sibling::PointOnRoute/@order or @order = following-sibling::PointOnRoute/@order]")

	r.addRule("JOURNEY_PATTERN_DUPLICATE_ORDER", "Duplicate order in JourneyPattern",
		"StopPointInJourneyPattern order values must be unique", types.ERROR,
		"//journeyPatterns/*/pointsInSequence/StopPointInJourneyPattern[@order = preceding-sibling::StopPointInJourneyPattern/@order or @order = following-sibling::StopPointInJourneyPattern/@order]")

	// Minimum element counts
	r.addRule("ROUTE_INSUFFICIENT_POINTS", "Route has insufficient points",
		"Route must have at least 2 points", types.ERROR,
		"//routes/Route[count(pointsInSequence/PointOnRoute) < 2]")

	r.addRule("JOURNEY_PATTERN_INSUFFICIENT_STOPS", "JourneyPattern has insufficient stops",
		"JourneyPattern must have at least 2 stops", types.ERROR,
		"//journeyPatterns/*[self::JourneyPattern or self::ServiceJourneyPattern][count(pointsInSequence/StopPointInJourneyPattern) < 2]")

	// Block structure validation
	r.addRule("BLOCK_NO_JOURNEYS", "Block without journeys",
		"Block must contain at least one journey", types.ERROR,
		"//blocks/Block[not(journeys) or count(journeys/*) = 0]")
}

// addDataConsistencyRules adds advanced data consistency validation
func (r *RuleRegistry) addDataConsistencyRules() {
	// Calendar consistency
	r.addRule("OPERATING_PERIOD_INVALID_DATES", "OperatingPeriod invalid date range",
		"OperatingPeriod FromDate must be before ToDate", types.ERROR,
		"//operatingPeriods/OperatingPeriod[FromDate >= ToDate]")

	r.addRule("SERVICE_CALENDAR_MISSING_PERIODS", "ServiceCalendar missing periods",
		"ServiceCalendar must have operating periods or day type assignments", types.ERROR,
		"//serviceCalendar/ServiceCalendar[not(operatingPeriods) and not(dayTypeAssignments)]")

	// Timing consistency
	r.addRule("PASSING_TIME_SEQUENCE_ERROR", "Passing time sequence error",
		"Passing times should increase through the journey", types.WARNING,
		"//vehicleJourneys/*/passingTimes/TimetabledPassingTime[position() > 1 and DepartureTime <= preceding-sibling::TimetabledPassingTime[1]/DepartureTime]")

	// Geographic consistency
	r.addRule("STOP_PLACE_NO_LOCATION", "StopPlace missing location",
		"StopPlace should have geographic coordinates", types.WARNING,
		"//stopPlaces/StopPlace[not(Centroid) and not(gml:Point) and not(Location)]")

	r.addRule("QUAY_NO_LOCATION", "Quay missing location",
		"Quay should have geographic coordinates", types.WARNING,
		"//stopPlaces/StopPlace/quays/Quay[not(Centroid) and not(gml:Point) and not(Location)]")

	// Version consistency
	r.addRule("INCONSISTENT_ELEMENT_VERSIONS", "Inconsistent element versions",
		"Related elements should have consistent versions", types.WARNING,
		"//*//*[@version != 'any' and @version != ancestor::*[@version][1]/@version]")

	// Network topology consistency
	r.addRule("ROUTE_DIRECTION_INCONSISTENCY", "Route direction inconsistency",
		"Route DirectionType should be consistent with DestinationDisplay", types.WARNING,
		"//routes/Route[DirectionType = 'inbound' and DestinationDisplay[contains(translate(Name/text(), 'ABCDEFGHIJKLMNOPQRSTUVWXYZ', 'abcdefghijklmnopqrstuvwxyz'), 'outbound')]]")
}

// addEUProfileRules adds EU NeTEx profile specific rules
func (r *RuleRegistry) addEUProfileRules() {
	// EU accessibility requirements
	r.addRule("EU_ACCESSIBILITY_INFO_REQUIRED", "EU accessibility information required",
		"EU profile encourages accessibility information", types.INFO,
		"//stopPlaces/StopPlace[not(AccessibilityAssessment) and not(placeEquipments)]")

	// EU multilingual support
	r.addRule("EU_MULTILINGUAL_NAMES", "EU multilingual name support",
		"EU profile supports multilingual names", types.INFO,
		"//Name[not(@lang) and ancestor::PublicationDelivery[@locale != @lang]]")

	// EU interchange requirements
	r.addRule("EU_INTERCHANGE_TIME_REQUIRED", "EU interchange time requirements",
		"EU profile requires realistic interchange times", types.WARNING,
		"//interchanges/ServiceJourneyInterchange[not(StandardTransferTime) and not(MinimumTransferTime)]")

	// EU fare information requirements
	r.addRule("EU_FARE_INFO_ENCOURAGED", "EU fare information encouraged",
		"EU profile encourages fare information", types.INFO,
		"//lines/*[self::Line or self::FlexibleLine][not(//fareProducts) and not(//tariffZones)]")

	// EU environmental information
	r.addRule("EU_ENVIRONMENTAL_INFO", "EU environmental information",
		"EU profile supports environmental impact information", types.INFO,
		"//vehicleTypes/VehicleType[not(FuelType) and not(EmissionClassification)]")
}
