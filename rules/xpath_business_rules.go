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

	// Nordic profile specific rules
	r.addNordicProfileRules()

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

// addNordicProfileRules adds Nordic NeTEx profile specific rules
func (r *RuleRegistry) addNordicProfileRules() {
	// Nordic ID format requirements
	r.addRule("NORDIC_ID_FORMAT_CODESPACE", "Invalid Nordic codespace format",
		"Nordic profile requires specific codespace format", types.WARNING,
		"//@id[not(matches(., '^[A-Z]{2,6}:.*'))]")

	// Nordic authority requirements
	r.addRule("NORDIC_AUTHORITY_REQUIRED", "Nordic profile missing Authority",
		"Nordic profile requires Authority definition", types.ERROR,
		"//PublicationDelivery[not(.//Authority)]")

	// Nordic transport mode restrictions
	r.addRule("NORDIC_TRANSPORT_MODE_RESTRICTED", "Nordic transport mode restrictions",
		"Nordic profile has restrictions on certain transport modes", types.WARNING,
		"//lines/*/TransportMode[text() = 'air' or text() = 'water'][not(ancestor::*/Name[contains(text(), 'Airport') or contains(text(), 'Ferry') or contains(text(), 'Boat')])]")

	// Nordic operator requirements
	r.addRule("NORDIC_OPERATOR_CONTACT_REQUIRED", "Nordic operator contact required",
		"Nordic profile requires operator contact information", types.WARNING,
		"//organisations/Operator[not(ContactDetails)]")

	// Nordic stop place requirements
	r.addRule("NORDIC_STOP_PLACE_TYPE_REQUIRED", "Nordic StopPlace type required",
		"Nordic profile requires StopPlaceType", types.WARNING,
		"//stopPlaces/StopPlace[not(StopPlaceType)]")
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

// addAdvancedReferenceRules adds complex reference validation
func (r *RuleRegistry) addAdvancedReferenceRules() {
	// Cross-frame reference validation
	r.addRule("CROSS_FRAME_LINE_REF", "Cross-frame LineRef validation",
		"LineRef should reference Line in ServiceFrame", types.ERROR,
		"//TimetableFrame//LineRef[@ref and not(@ref = //ServiceFrame//lines/*/@id)]")

	r.addRule("CROSS_FRAME_OPERATOR_REF", "Cross-frame OperatorRef validation",
		"OperatorRef should reference Operator in ResourceFrame", types.ERROR,
		"//ServiceFrame//OperatorRef[@ref and not(@ref = //ResourceFrame//organisations/Operator/@id)]")

	r.addRule("CROSS_FRAME_STOP_PLACE_REF", "Cross-frame StopPlaceRef validation",
		"StopPlaceRef should reference StopPlace in SiteFrame", types.ERROR,
		"//ServiceFrame//StopPlaceRef[@ref and not(@ref = //SiteFrame//stopPlaces/StopPlace/@id)]")

	// Circular reference detection
	r.addRule("CIRCULAR_REFERENCE_LINE_GROUP", "Circular reference in Line-Group",
		"Lines and GroupOfLines should not have circular references", types.ERROR,
		"//lines/*[RepresentedByGroupRef/@ref = @id]")

	// Orphaned reference detection
	r.addRule("ORPHANED_STOP_ASSIGNMENT", "Orphaned stop assignment",
		"PassengerStopAssignment references non-existent elements", types.ERROR,
		"//stopAssignments/PassengerStopAssignment[ScheduledStopPointRef/@ref and not(ScheduledStopPointRef/@ref = //scheduledStopPoints/ScheduledStopPoint/@id)]")
}

// addValidationRuleDocumentationExamples adds rules that demonstrate validation capabilities
func (r *RuleRegistry) addValidationRuleDocumentationExamples() {
	// Complex business rules requiring multiple element validation
	r.addRule("BUS_ROUTE_STOP_VALIDATION", "Bus route stop validation",
		"Bus routes should connect logical stops", types.WARNING,
		"//routes/Route[//lines/*[@id=current()/LineRef/@ref]/TransportMode='bus']/pointsInSequence[count(PointOnRoute) < 2]")

	// Timetable feasibility validation
	r.addRule("TIMETABLE_FEASIBILITY", "Timetable feasibility check",
		"Service journey times should be realistic", types.WARNING,
		"//vehicleJourneys/ServiceJourney/passingTimes[TimetabledPassingTime[1]/DepartureTime > TimetabledPassingTime[last()]/ArrivalTime]")

	// Network connectivity validation
	r.addRule("NETWORK_CONNECTIVITY", "Network connectivity validation",
		"Routes should form connected network", types.INFO,
		"//routes/Route[not(pointsInSequence/PointOnRoute[1]/ScheduledStopPointRef = //routes/Route[position() != current()/position()]/pointsInSequence/PointOnRoute[last()]/ScheduledStopPointRef)]")
}
