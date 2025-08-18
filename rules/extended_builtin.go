package rules

import "github.com/theoremus-urban-solutions/netex-validator/types"

// loadExtendedBuiltinRules loads additional validation rules to reach parity with Java version
func (r *RuleRegistry) loadExtendedBuiltinRules() {
	// Transport Mode Validation Rules
	r.addTransportModeRules()

	// Booking and Flexible Service Rules
	r.addBookingRules()

	// Service Journey and Timetable Rules
	r.addServiceJourneyRules()

	// Stop Place and Quay Rules
	r.addStopPlaceRules()

	// Journey Pattern Rules
	r.addJourneyPatternRules()

	// Network and Operator Rules
	r.addNetworkRules()

	// Fare and Pricing Rules
	r.addFareRules()

	// Calendar and Validity Rules
	r.addCalendarRules()

	// Vehicle and Equipment Rules
	r.addVehicleRules()

	// Interchange and Connection Rules
	r.addInterchangeRules()

	// Reference consistency rules across frames (EU)
	r.addReferenceConsistencyRules()

	// Load advanced rule sets for comprehensive validation
	r.loadAdvancedTransportRules()
	r.loadServiceJourneyRules()
	r.loadFlexibleServiceRules()
	r.loadXPathBusinessRules()

	// Note: Some advanced business logic rules with complex cross-references may be limited
	// due to XPath current() function limitations in the xmlquery library
}

// addTransportModeRules adds comprehensive transport mode validation rules
func (r *RuleRegistry) addTransportModeRules() {
	// Use or-based XPath expressions instead of parenthesized lists (not supported by xmlquery)

	r.addRule("TRANSPORT_MODE_1", "Invalid transport mode on Line", "Line has invalid transport mode", types.ERROR,
		"//lines/*[self::Line or self::FlexibleLine]/TransportMode[not(text() = 'coach' or text() = 'bus' or text() = 'tram' or text() = 'rail' or text() = 'metro' or text() = 'air' or text() = 'taxi' or text() = 'water' or text() = 'cableway' or text() = 'funicular' or text() = 'unknown')]")

	r.addRule("TRANSPORT_MODE_2", "Invalid transport mode on ServiceJourney", "ServiceJourney has invalid transport mode", types.ERROR,
		"//vehicleJourneys/*[self::ServiceJourney or self::DatedServiceJourney]/TransportMode[not(text() = 'coach' or text() = 'bus' or text() = 'tram' or text() = 'rail' or text() = 'metro' or text() = 'air' or text() = 'taxi' or text() = 'water' or text() = 'cableway' or text() = 'funicular' or text() = 'unknown')]")

	// Transport submode validation
	r.addRule("TRANSPORT_SUBMODE_1", "Missing transport submode", "Transport submode is required when transport mode is specified", types.WARNING,
		"//*/TransportMode[text() != 'unknown' and not(following-sibling::TransportSubmode or preceding-sibling::TransportSubmode)]")

	r.addRule("TRANSPORT_SUBMODE_2", "Inconsistent transport mode and submode", "Transport submode does not match transport mode", types.ERROR,
		"//*/TransportMode[text()='bus']/following-sibling::TransportSubmode[not(text() = 'localBus' or text() = 'expressBus' or text() = 'nightBus' or text() = 'postBus' or text() = 'specialNeedsBus' or text() = 'mobilityBus' or text() = 'mobilityBusForRegisteredDisabled' or text() = 'sightseeingBus' or text() = 'shuttleBus' or text() = 'highFrequencyBus' or text() = 'dedicatedLaneBus' or text() = 'schoolBus' or text() = 'schoolAndPublicBus' or text() = 'railReplacementBus' or text() = 'demandAndResponseBus' or text() = 'airportLinkBus')]")
}

// addBookingRules adds flexible service and booking validation rules
func (r *RuleRegistry) addBookingRules() {
	r.addRule("BOOKING_1", "Missing booking arrangements", "Flexible service missing booking arrangements", types.ERROR,
		"//lines/FlexibleLine[not(BookingContact) and not(BookingUrl) and not(bookingArrangements)]")

	r.addRule("BOOKING_2", "Invalid booking method", "Invalid booking method specified", types.ERROR,
		"//bookingArrangements/BookingMethod[not(text() = 'callDriver' or text() = 'callOffice' or text() = 'online' or text() = 'phoneAtStop' or text() = 'text')]")

	r.addRule("BOOKING_3", "Invalid booking access", "Invalid booking access method", types.ERROR,
		"//bookingArrangements/BookingAccess[not(text() = 'public' or text() = 'authorisedPublic' or text() = 'staff' or text() = 'other')]")

	r.addRule("BOOKING_4", "Invalid booking when", "Invalid booking time specification", types.ERROR,
		"//bookingArrangements/BookWhen[not(text() = 'timeOfTravelOnly' or text() = 'dayOfTravelOnly' or text() = 'untilPreviousDay' or text() = 'advanceOnly' or text() = 'other')]")

	r.addRule("BOOKING_5", "Missing minimum booking period", "Minimum booking period required for advance booking", types.WARNING,
		"//bookingArrangements[BookWhen='advanceOnly' and not(MinimumBookingPeriod)]")
}

// addServiceJourneyRules adds service journey and timetable validation rules
func (r *RuleRegistry) addServiceJourneyRules() {
	r.addRule("SERVICE_JOURNEY_1", "ServiceJourney missing JourneyPatternRef", "ServiceJourney must reference a JourneyPattern", types.ERROR,
		"//vehicleJourneys/*[self::ServiceJourney or self::DatedServiceJourney][not(JourneyPatternRef)]")

	r.addRule("SERVICE_JOURNEY_2", "ServiceJourney missing OperatorRef", "ServiceJourney must reference an Operator", types.ERROR,
		"//vehicleJourneys/*[self::ServiceJourney or self::DatedServiceJourney][not(OperatorRef)]")

	r.addRule("SERVICE_JOURNEY_3", "ServiceJourney missing LineRef", "ServiceJourney must reference a Line", types.ERROR,
		"//vehicleJourneys/*[self::ServiceJourney or self::DatedServiceJourney][not(LineRef)]")

	r.addRule("SERVICE_JOURNEY_5", "Missing departure time", "First stop must have departure time", types.ERROR,
		"//vehicleJourneys/*[self::ServiceJourney or self::DatedServiceJourney]/passingTimes/TimetabledPassingTime[1][not(DepartureTime)]")

	r.addRule("SERVICE_JOURNEY_6", "Missing arrival time", "Last stop must have arrival time", types.ERROR,
		"//vehicleJourneys/*[self::ServiceJourney or self::DatedServiceJourney]/passingTimes/TimetabledPassingTime[last()][not(ArrivalTime)]")

	r.addRule("SERVICE_JOURNEY_7", "Duplicate TimetabledPassingTime ID", "TimetabledPassingTime IDs must be unique within ServiceJourney", types.ERROR,
		"//vehicleJourneys/*[self::ServiceJourney or self::DatedServiceJourney]/passingTimes/TimetabledPassingTime[@id = preceding-sibling::TimetabledPassingTime/@id or @id = following-sibling::TimetabledPassingTime/@id]")
}

// addStopPlaceRules adds stop place and quay validation rules
func (r *RuleRegistry) addStopPlaceRules() {
	r.addRule("STOP_PLACE_1", "StopPlace missing Name", "StopPlace must have a Name", types.ERROR,
		"//stopPlaces/StopPlace[not(Name) or normalize-space(Name) = '']")

	r.addRule("STOP_PLACE_2", "StopPlace missing Centroid", "StopPlace should have location coordinates", types.WARNING,
		"//stopPlaces/StopPlace[not(Centroid)]")

	r.addRule("STOP_PLACE_3", "Quay missing Name", "Quay must have a Name", types.ERROR,
		"//stopPlaces/StopPlace/quays/Quay[not(Name) or normalize-space(Name) = '']")

	r.addRule("STOP_PLACE_4", "Quay missing Centroid", "Quay should have location coordinates", types.WARNING,
		"//stopPlaces/StopPlace/quays/Quay[not(Centroid)]")

	r.addRule("STOP_PLACE_5", "Invalid StopPlaceType", "StopPlace has invalid type", types.ERROR,
		"//stopPlaces/StopPlace/StopPlaceType[not(text() = 'onstreetBus' or text() = 'onstreetTram' or text() = 'airport' or text() = 'railStation' or text() = 'metroStation' or text() = 'busStation' or text() = 'coachStation' or text() = 'tramStation' or text() = 'harbourPort' or text() = 'ferryPort' or text() = 'ferryStop' or text() = 'liftStation' or text() = 'vehicleRailInterchange' or text() = 'other')]")
}

// addJourneyPatternRules adds journey pattern validation rules
func (r *RuleRegistry) addJourneyPatternRules() {
	r.addRule("JOURNEY_PATTERN_1", "JourneyPattern missing Name", "JourneyPattern should have a Name", types.WARNING,
		"//journeyPatterns/*[self::JourneyPattern or self::ServiceJourneyPattern][not(Name) or normalize-space(Name) = '']")

	r.addRule("JOURNEY_PATTERN_2", "JourneyPattern missing RouteRef", "JourneyPattern must reference a Route", types.ERROR,
		"//journeyPatterns/*[self::JourneyPattern or self::ServiceJourneyPattern][not(RouteRef)]")

	r.addRule("JOURNEY_PATTERN_3", "JourneyPattern missing StopPoints", "JourneyPattern must have stop points", types.ERROR,
		"//journeyPatterns/*[self::JourneyPattern or self::ServiceJourneyPattern][not(pointsInSequence) or count(pointsInSequence/StopPointInJourneyPattern) < 2]")

	r.addRule("JOURNEY_PATTERN_4", "StopPointInJourneyPattern missing order", "StopPointInJourneyPattern must have order", types.ERROR,
		"//journeyPatterns/*[self::JourneyPattern or self::ServiceJourneyPattern]/pointsInSequence/StopPointInJourneyPattern[not(@order)]")

	r.addRule("JOURNEY_PATTERN_5", "Duplicate order in JourneyPattern", "Order values must be unique within JourneyPattern", types.ERROR,
		"//journeyPatterns/*[self::JourneyPattern or self::ServiceJourneyPattern]/pointsInSequence/StopPointInJourneyPattern[@order = preceding-sibling::StopPointInJourneyPattern/@order or @order = following-sibling::StopPointInJourneyPattern/@order]")
}

// addNetworkRules adds network and operator validation rules
func (r *RuleRegistry) addNetworkRules() {
	r.addRule("NETWORK_1", "Network missing Name", "Network must have a Name", types.ERROR,
		"//networks/Network[not(Name) or normalize-space(Name) = '']")

	r.addRule("NETWORK_2", "Network missing AuthorityRef", "Network should reference an Authority", types.WARNING,
		"//networks/Network[not(AuthorityRef)]")

	r.addRule("OPERATOR_1", "Operator missing Name", "Operator must have a Name", types.ERROR,
		"//operators/Operator[not(Name) or normalize-space(Name) = '']")

	r.addRule("OPERATOR_2", "Operator missing ContactDetails", "Operator should have contact details", types.WARNING,
		"//operators/Operator[not(ContactDetails)]")

	r.addRule("AUTHORITY_1", "Authority missing Name", "Authority must have a Name", types.ERROR,
		"//organisations/Authority[not(Name) or normalize-space(Name) = '']")
}

// addFareRules adds fare and pricing validation rules
func (r *RuleRegistry) addFareRules() {
	r.addRule("FARE_1", "FareZone missing Name", "FareZone must have a Name", types.ERROR,
		"//fareZones/FareZone[not(Name) or normalize-space(Name) = '']")

	r.addRule("FARE_2", "TariffZone missing Name", "TariffZone must have a Name", types.ERROR,
		"//tariffZones/TariffZone[not(Name) or normalize-space(Name) = '']")

	r.addRule("FARE_3", "FareProduct missing Name", "FareProduct must have a Name", types.ERROR,
		"//fareProducts/FareProduct[not(Name) or normalize-space(Name) = '']")

	r.addRule("FARE_4", "Missing price for FareProduct", "FareProduct should have a price", types.WARNING,
		"//fareProducts/FareProduct[not(prices)]")
}

// addCalendarRules adds calendar and validity period validation rules
func (r *RuleRegistry) addCalendarRules() {
	r.addRule("CALENDAR_1", "DayType missing Name", "DayType must have a Name", types.ERROR,
		"//dayTypes/DayType[not(Name) or normalize-space(Name) = '']")

	r.addRule("CALENDAR_2", "OperatingDay missing CalendarDate", "OperatingDay must have CalendarDate", types.ERROR,
		"//operatingDays/OperatingDay[not(CalendarDate)]")

	r.addRule("CALENDAR_3", "ServiceCalendar missing FromDate", "ServiceCalendar must have FromDate", types.ERROR,
		"//serviceCalendar/ServiceCalendar[not(FromDate)]")

	r.addRule("CALENDAR_4", "ServiceCalendar missing ToDate", "ServiceCalendar must have ToDate", types.ERROR,
		"//serviceCalendar/ServiceCalendar[not(ToDate)]")

	r.addRule("CALENDAR_5", "Invalid date range", "FromDate must be before ToDate", types.ERROR,
		"//serviceCalendar/ServiceCalendar[FromDate >= ToDate]")
}

// addVehicleRules adds vehicle and equipment validation rules
func (r *RuleRegistry) addVehicleRules() {
	r.addRule("VEHICLE_1", "VehicleType missing Name", "VehicleType must have a Name", types.ERROR,
		"//vehicleTypes/VehicleType[not(Name) or normalize-space(Name) = '']")

	r.addRule("VEHICLE_2", "Vehicle missing VehicleTypeRef", "Vehicle must reference a VehicleType", types.ERROR,
		"//vehicles/Vehicle[not(VehicleTypeRef)]")

	r.addRule("VEHICLE_3", "Block missing Name", "Block should have a Name", types.WARNING,
		"//blocks/Block[not(Name) or normalize-space(Name) = '']")
}

// addInterchangeRules adds interchange and connection validation rules
func (r *RuleRegistry) addInterchangeRules() {
	r.addRule("INTERCHANGE_1", "ServiceJourneyInterchange missing FromPointRef", "Interchange must have FromPointRef", types.ERROR,
		"//serviceJourneyInterchanges/ServiceJourneyInterchange[not(FromPointRef)]")

	r.addRule("INTERCHANGE_2", "ServiceJourneyInterchange missing ToPointRef", "Interchange must have ToPointRef", types.ERROR,
		"//serviceJourneyInterchanges/ServiceJourneyInterchange[not(ToPointRef)]")

	r.addRule("INTERCHANGE_3", "ServiceJourneyInterchange missing FromJourneyRef", "Interchange must have FromJourneyRef", types.ERROR,
		"//serviceJourneyInterchanges/ServiceJourneyInterchange[not(FromJourneyRef)]")

	r.addRule("INTERCHANGE_4", "ServiceJourneyInterchange missing ToJourneyRef", "Interchange must have ToJourneyRef", types.ERROR,
		"//serviceJourneyInterchanges/ServiceJourneyInterchange[not(ToJourneyRef)]")

	r.addRule("INTERCHANGE_5", "Invalid interchange duration", "Interchange duration must be positive", types.ERROR,
		"//serviceJourneyInterchanges/ServiceJourneyInterchange/StandardTransferTime[number(.) <= 0]")
}

// addReferenceConsistencyRules ensures that common *Ref elements point to existing targets
func (r *RuleRegistry) addReferenceConsistencyRules() {
	// OperatorRef must exist in ResourceFrame organisations/Operator
	r.addRule("REF_OPERATOR_1", "OperatorRef undefined", "OperatorRef does not reference an existing Operator in ResourceFrame", types.ERROR,
		"//*[local-name()='ServiceFrame']//*[local-name()='OperatorRef' and (@ref or normalize-space(text())!='')][not(@ref=//ResourceFrame//*[local-name()='Operator']/@id) and not(normalize-space(text())=//ResourceFrame//*[local-name()='Operator']/@id)]")

	// LineRef must exist in ServiceFrame lines/Line or FlexibleLine
	r.addRule("REF_LINE_1", "LineRef undefined", "LineRef does not reference an existing Line in ServiceFrame", types.ERROR,
		"//*[local-name()='ServiceFrame']//*[local-name()='LineRef' and (@ref or normalize-space(text())!='')][not(@ref=//ServiceFrame//*[local-name()='lines']/*[local-name()='Line' or local-name()='FlexibleLine']/@id) and not(normalize-space(text())=//ServiceFrame//*[local-name()='lines']/*[local-name()='Line' or local-name()='FlexibleLine']/@id)]")

	// JourneyPatternRef must exist in ServiceFrame journeyPatterns
	r.addRule("REF_JOURNEY_PATTERN_1", "JourneyPatternRef undefined", "JourneyPatternRef does not reference an existing JourneyPattern in ServiceFrame", types.ERROR,
		"//*[local-name()='ServiceFrame']//*[local-name()='JourneyPatternRef' and (@ref or normalize-space(text())!='')][not(@ref=//ServiceFrame//*[local-name()='journeyPatterns']/*[local-name()='JourneyPattern' or local-name()='ServiceJourneyPattern']/@id) and not(normalize-space(text())=//ServiceFrame//*[local-name()='journeyPatterns']/*[local-name()='JourneyPattern' or local-name()='ServiceJourneyPattern']/@id)]")

	// DatedServiceJourney/ServiceJourneyRef must exist in ServiceFrame vehicleJourneys/ServiceJourney
	r.addRule("REF_SERVICE_JOURNEY_1", "ServiceJourneyRef undefined", "ServiceJourneyRef in TimetableFrame does not reference an existing ServiceJourney in ServiceFrame", types.ERROR,
		"//*[local-name()='TimetableFrame']//*[local-name()='DatedServiceJourney']/*[local-name()='ServiceJourneyRef' and (@ref or normalize-space(text())!='')][not(@ref=//ServiceFrame//*[local-name()='vehicleJourneys']/*[local-name()='ServiceJourney']/@id) and not(normalize-space(text())=//ServiceFrame//*[local-name()='vehicleJourneys']/*[local-name()='ServiceJourney']/@id)]")

	// ScheduledStopPointRef must exist in ServiceFrame scheduledStopPoints
	r.addRule("REF_STOP_POINT_1", "ScheduledStopPointRef undefined", "ScheduledStopPointRef does not reference an existing ScheduledStopPoint in ServiceFrame", types.ERROR,
		"//*[local-name()='ServiceFrame']//*[local-name()='ScheduledStopPointRef' and (@ref or normalize-space(text())!='')][not(@ref=//ServiceFrame//*[local-name()='scheduledStopPoints']/*[local-name()='ScheduledStopPoint']/@id) and not(normalize-space(text())=//ServiceFrame//*[local-name()='scheduledStopPoints']/*[local-name()='ScheduledStopPoint']/@id)]")
}
