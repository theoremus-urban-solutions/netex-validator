package rules

import "github.com/theoremus-urban-solutions/netex-validator/types"

// loadFlexibleServiceRules loads comprehensive flexible service and booking validation rules
func (r *RuleRegistry) loadFlexibleServiceRules() {
	// FlexibleLine validation rules
	r.addFlexibleLineRules()

	// Booking property validation rules
	r.addBookingPropertyRules()

	// FlexibleService validation rules
	r.addFlexibleServiceValidationRules()

	// Booking arrangement validation
	r.addBookingArrangementRules()

	// Flexible area validation
	r.addFlexibleAreaRules()
}

// addFlexibleLineRules adds FlexibleLine specific validation
func (r *RuleRegistry) addFlexibleLineRules() {
	// FlexibleLine must have FlexibleLineType
	r.addRule("FLEXIBLE_LINE_MISSING_TYPE", "FlexibleLine missing FlexibleLineType",
		"FlexibleLine must specify FlexibleLineType", types.ERROR,
		"//lines/FlexibleLine[not(FlexibleLineType)]")

	// FlexibleLineType must be valid
	validFlexibleLineTypes := []string{
		"fixedStop", "flexibleAreasOnly", "hailAndRideAreas", "flexibleAreasAndStops",
		"hailAndRideSections", "fixedStopAreaWide", "freeAreaAreaWide", "mixedFlexible",
		"mixedFlexibleAndFixed", "fixed", "mainRouteWithFlexibleEnds", "flexibleRoute",
		"demandAndResponseServices",
	}

	invalidFlexibleLineTypeCondition := r.buildInvalidValueCondition(validFlexibleLineTypes)

	r.addRule("FLEXIBLE_LINE_INVALID_TYPE", "FlexibleLine invalid FlexibleLineType",
		"FlexibleLine has invalid FlexibleLineType value", types.ERROR,
		"//lines/FlexibleLine/FlexibleLineType["+invalidFlexibleLineTypeCondition+"]")

	// FlexibleLine should have booking information for certain types
	r.addRule("FLEXIBLE_LINE_MISSING_BOOKING_INFO", "FlexibleLine missing booking information",
		"FlexibleLine with flexible areas should have booking information", types.ERROR,
		"//lines/FlexibleLine[FlexibleLineType = 'flexibleAreasOnly' or FlexibleLineType = 'hailAndRideAreas' or FlexibleLineType = 'demandAndResponseServices'][not(BookWhen) and not(MinimumBookingPeriod) and not(bookingArrangements) and not(BookingContact) and not(BookingUrl)]")

	// FlexibleLine transport mode restrictions
	r.addRule("FLEXIBLE_LINE_INVALID_TRANSPORT_MODE", "FlexibleLine invalid transport mode",
		"FlexibleLine should use appropriate transport modes", types.WARNING,
		"//lines/FlexibleLine[TransportMode and not(TransportMode = 'bus' or TransportMode = 'taxi' or TransportMode = 'water' or TransportMode = 'unknown')]")
}

// addBookingPropertyRules adds booking property validation based on Java validators
func (r *RuleRegistry) addBookingPropertyRules() {
	// BookWhen must be valid value
	validBookWhenValues := []string{
		"timeOfTravelOnly", "dayOfTravelOnly", "untilPreviousDay", "advanceAndDayOfTravel", "other",
	}

	invalidBookWhenCondition := r.buildInvalidValueCondition(validBookWhenValues)

	r.addRule("BOOKING_INVALID_BOOK_WHEN", "Invalid BookWhen property",
		"BookWhen has invalid value", types.ERROR,
		"//*/BookWhen["+invalidBookWhenCondition+"]")

	// BookingAccess must be valid value
	validBookingAccessValues := []string{
		"public", "authorisedPublic", "staff", "other",
	}

	invalidBookingAccessCondition := r.buildInvalidValueCondition(validBookingAccessValues)

	r.addRule("BOOKING_INVALID_ACCESS", "Invalid BookingAccess property",
		"BookingAccess has invalid value", types.ERROR,
		"//*/BookingAccess["+invalidBookingAccessCondition+"]")

	// BookingMethod must be valid value
	validBookingMethodValues := []string{
		"callDriver", "callOffice", "online", "phoneAtStop", "text", "none", "other",
	}

	invalidBookingMethodCondition := r.buildInvalidValueCondition(validBookingMethodValues)

	r.addRule("BOOKING_INVALID_METHOD", "Invalid BookingMethod property",
		"BookingMethod has invalid value", types.ERROR,
		"//*/BookingMethod["+invalidBookingMethodCondition+"]")

	// BuyWhen must be valid value
	validBuyWhenValues := []string{
		"timeOfTravelOnly", "dayOfTravelOnly", "untilPreviousDay", "advanceAndDayOfTravel", "other",
	}

	invalidBuyWhenCondition := r.buildInvalidValueCondition(validBuyWhenValues)

	r.addRule("BOOKING_INVALID_BUY_WHEN", "Invalid BuyWhen property",
		"BuyWhen has invalid value", types.ERROR,
		"//*/BuyWhen["+invalidBuyWhenCondition+"]")

	// Mandatory booking properties for flexible services
	r.addRule("BOOKING_MANDATORY_PROPERTIES_MISSING", "Missing mandatory booking properties",
		"Flexible service missing required booking properties", types.ERROR,
		"//lines/FlexibleLine[FlexibleLineType = 'flexibleAreasOnly' or FlexibleLineType = 'hailAndRideAreas' or FlexibleLineType = 'demandAndResponseServices'][not(BookWhen) and not(MinimumBookingPeriod)]")

	// BookWhen and MinimumBookingPeriod should not both be present
	r.addRule("BOOKING_CONFLICTING_PROPERTIES", "Conflicting booking properties",
		"BookWhen and MinimumBookingPeriod should not both be specified", types.WARNING,
		"//lines/FlexibleLine[BookWhen and MinimumBookingPeriod]")

	// LatestBookingTime requires BookWhen
	r.addRule("BOOKING_LATEST_TIME_WITHOUT_WHEN", "LatestBookingTime without BookWhen",
		"LatestBookingTime specified without BookWhen", types.WARNING,
		"//lines/FlexibleLine[LatestBookingTime and not(BookWhen)]")

	// MinimumBookingPeriod should have reasonable value
	r.addRule("BOOKING_UNREASONABLE_MINIMUM_PERIOD", "Unreasonable MinimumBookingPeriod",
		"MinimumBookingPeriod should be reasonable", types.WARNING,
		"//lines/FlexibleLine/MinimumBookingPeriod[. > 'P7D' or . < 'PT15M']")
}

// addFlexibleServiceValidationRules adds FlexibleService validation
func (r *RuleRegistry) addFlexibleServiceValidationRules() {
	// FlexibleService must have FlexibleServiceType
	r.addRule("FLEXIBLE_SERVICE_MISSING_TYPE", "FlexibleService missing FlexibleServiceType",
		"FlexibleService must specify FlexibleServiceType", types.ERROR,
		"//FlexibleService[not(FlexibleServiceType)]")

	// FlexibleServiceType must be valid
	validFlexibleServiceTypes := []string{
		"dynamicPassengerInformation", "fixedHeadwayService", "fixedPassingTimes", "notFixedPassingTimes",
	}

	invalidFlexibleServiceTypeCondition := r.buildInvalidValueCondition(validFlexibleServiceTypes)

	r.addRule("FLEXIBLE_SERVICE_INVALID_TYPE", "FlexibleService invalid FlexibleServiceType",
		"FlexibleService has invalid FlexibleServiceType value", types.ERROR,
		"//FlexibleService/FlexibleServiceType["+invalidFlexibleServiceTypeCondition+"]")

	// FlexibleService with fixedPassingTimes should have timetabled passing times
	r.addRule("FLEXIBLE_SERVICE_MISSING_TIMES", "FlexibleService missing passing times",
		"FlexibleService with fixedPassingTimes should have timetabled passing times", types.WARNING,
		"//FlexibleService[FlexibleServiceType = 'fixedPassingTimes' and not(//vehicleJourneys/ServiceJourney[FlexibleServiceRef/@ref = current()/@id]/passingTimes)]")

	// FlexibleService should be referenced by ServiceJourney
	r.addRule("FLEXIBLE_SERVICE_UNREFERENCED", "Unreferenced FlexibleService",
		"FlexibleService should be referenced by at least one ServiceJourney", types.WARNING,
		"//FlexibleService[not(@id = //vehicleJourneys/ServiceJourney/FlexibleServiceRef/@ref)]")
}

// addBookingArrangementRules adds booking arrangement validation
func (r *RuleRegistry) addBookingArrangementRules() {
	// BookingArrangements should have at least one booking method
	r.addRule("BOOKING_ARRANGEMENTS_NO_METHOD", "BookingArrangements without method",
		"BookingArrangements should specify at least one booking method", types.WARNING,
		"//bookingArrangements[not(BookingMethod) and not(BookingContact) and not(BookingUrl)]")

	// BookingContact should have valid contact information
	r.addRule("BOOKING_CONTACT_INCOMPLETE", "Incomplete booking contact",
		"BookingContact should have phone, email, or URL", types.WARNING,
		"//bookingArrangements/BookingContact[not(Phone) and not(Email) and not(Url)]")

	// BookingUrl should be valid URL format
	r.addRule("BOOKING_INVALID_URL", "Invalid booking URL",
		"BookingUrl should be valid URL format", types.WARNING,
		"//bookingArrangements/BookingUrl[not(starts-with(text(), 'http://') or starts-with(text(), 'https://'))]")

	// BookingNote should be meaningful
	r.addRule("BOOKING_NOTE_EMPTY", "Empty booking note",
		"BookingNote should contain meaningful information", types.WARNING,
		"//bookingArrangements/BookingNote[not(text()) or normalize-space(text()) = '']")

	// Multiple booking methods should be consistent
	r.addRule("BOOKING_INCONSISTENT_METHODS", "Inconsistent booking methods",
		"Multiple booking methods should be consistent", types.WARNING,
		"//bookingArrangements[count(BookingMethod) > 1 and BookingMethod = 'none']")
}

// addFlexibleAreaRules adds flexible area and zone validation
func (r *RuleRegistry) addFlexibleAreaRules() {
	// FlexibleArea should have geometry
	r.addRule("FLEXIBLE_AREA_NO_GEOMETRY", "FlexibleArea missing geometry",
		"FlexibleArea should have geographic definition", types.ERROR,
		"//FlexibleArea[not(Polygon) and not(gml:Polygon) and not(areas)]")

	// FlexibleStopPlace should be properly defined
	r.addRule("FLEXIBLE_STOP_PLACE_INCOMPLETE", "Incomplete FlexibleStopPlace",
		"FlexibleStopPlace should have proper geographic bounds", types.WARNING,
		"//FlexibleStopPlace[not(areas) and not(Polygon) and not(gml:Polygon)]")

	// FlexibleQuay should reference FlexibleStopPlace
	r.addRule("FLEXIBLE_QUAY_MISSING_STOP_PLACE", "FlexibleQuay missing StopPlace reference",
		"FlexibleQuay should reference a FlexibleStopPlace", types.WARNING,
		"//FlexibleQuay[not(StopPlaceRef)]")

	// HailAndRideArea should have reasonable size
	r.addRule("HAIL_AND_RIDE_AREA_SIZE", "HailAndRideArea size validation",
		"HailAndRideArea should have reasonable geographic bounds", types.WARNING,
		"//HailAndRideArea[not(Polygon) and not(gml:Polygon)]")
}

// addFlexibleServiceJourneyRules adds flexible service journey validation
func (r *RuleRegistry) addFlexibleServiceJourneyRules() {
	// FlexibleServiceJourney should have appropriate timing information
	r.addRule("FLEXIBLE_SJ_TIMING_INCONSISTENT", "Inconsistent flexible service journey timing",
		"FlexibleServiceJourney timing should match FlexibleServiceType", types.WARNING,
		"//vehicleJourneys/ServiceJourney[FlexibleServiceRef][//FlexibleService[@id = current()/FlexibleServiceRef/@ref]/FlexibleServiceType = 'notFixedPassingTimes' and passingTimes/TimetabledPassingTime[DepartureTime and ArrivalTime]]")

	// Dynamic flexible services should not have fixed times
	r.addRule("DYNAMIC_FLEXIBLE_SERVICE_FIXED_TIMES", "Dynamic flexible service with fixed times",
		"Dynamic flexible service should not have fixed passing times", types.ERROR,
		"//vehicleJourneys/ServiceJourney[FlexibleServiceRef][//FlexibleService[@id = current()/FlexibleServiceRef/@ref]/FlexibleServiceType = 'dynamicPassengerInformation' and passingTimes/TimetabledPassingTime[DepartureTime or ArrivalTime]]")

	// Flexible service journey should have appropriate stop assignments
	r.addRule("FLEXIBLE_SJ_STOP_ASSIGNMENT_MISSING", "Flexible service journey missing stop assignments",
		"Flexible service journey should have proper stop assignments", types.WARNING,
		"//vehicleJourneys/ServiceJourney[FlexibleServiceRef and not(passingTimes/TimetabledPassingTime/StopPointInJourneyPatternRef)]")
}

// buildInvalidValueCondition builds XPath condition for invalid enum values
func (r *RuleRegistry) buildInvalidValueCondition(validValues []string) string {
	if len(validValues) == 0 {
		return "true()"
	}

	condition := "not(text() = '" + validValues[0] + "'"
	for _, value := range validValues[1:] {
		condition += " or text() = '" + value + "'"
	}
	condition += ")"

	return condition
}
