package rules

import "github.com/theoremus-urban-solutions/netex-validator/types"

// loadAdvancedTransportRules loads comprehensive transport mode and submode validation rules
func (r *RuleRegistry) loadAdvancedTransportRules() {
	// Transport mode validation rules based on Java ValidateAllowedTransportMode
	r.addTransportModeValidationRules()

	// Transport submode validation rules based on Java ValidateAllowedTransportSubMode
	r.addTransportSubModeValidationRules()

	// Contextual transport mode rules
	r.addContextualTransportModeRules()
}

// addTransportModeValidationRules adds basic transport mode validation
func (r *RuleRegistry) addTransportModeValidationRules() {
	// Valid transport modes according to EU NeTEx profile
	validModes := []string{
		"coach", "bus", "tram", "rail", "metro", "air", "taxi", "water", "cableway", "funicular", "unknown",
	}

	// Build XPath condition for invalid transport modes
	invalidModeCondition := "not(text() = '" + validModes[0] + "'"
	for _, mode := range validModes[1:] {
		invalidModeCondition += " or text() = '" + mode + "'"
	}
	invalidModeCondition += ")"

	r.addRule("TRANSPORT_MODE_INVALID_LINE", "Invalid transport mode on Line",
		"Line has invalid TransportMode value", types.ERROR,
		"//lines/*[self::Line or self::FlexibleLine]/TransportMode["+invalidModeCondition+"]")

	r.addRule("TRANSPORT_MODE_INVALID_SERVICE_JOURNEY", "Invalid transport mode on ServiceJourney",
		"ServiceJourney has invalid TransportMode value", types.ERROR,
		"//vehicleJourneys/*[self::ServiceJourney or self::DatedServiceJourney]/TransportMode["+invalidModeCondition+"]")
}

// addTransportSubModeValidationRules adds detailed transport submode validation
func (r *RuleRegistry) addTransportSubModeValidationRules() {
	// Bus submodes
	r.addBusSubModeRules()

	// Rail submodes
	r.addRailSubModeRules()

	// Tram submodes
	r.addTramSubModeRules()

	// Metro submodes
	r.addMetroSubModeRules()

	// Water submodes
	r.addWaterSubModeRules()

	// Air submodes
	r.addAirSubModeRules()

	// Coach submodes
	r.addCoachSubModeRules()

	// Other transport submodes
	r.addOtherTransportSubModeRules()
}

// addBusSubModeRules validates bus transport submodes
func (r *RuleRegistry) addBusSubModeRules() {
	validBusSubmodes := []string{
		"localBus", "regionalBus", "expressBus", "nightBus", "postBus", "specialNeedsBus",
		"mobilityBus", "mobilityBusForRegisteredDisabled", "sightseeingBus", "shuttleBus",
		"schoolBus", "schoolAndPublicServiceBus", "railReplacementBus", "demandAndResponseBus",
		"airportLinkBus", "highFrequencyBus", "dedicatedLaneBus",
	}

	invalidBusSubmodeCondition := r.buildInvalidSubmodeCondition(validBusSubmodes)

	r.addRule("TRANSPORT_SUBMODE_BUS_INVALID_LINE", "Invalid bus submode on Line",
		"Line with bus transport has invalid TransportSubmode", types.ERROR,
		"//lines/*[self::Line or self::FlexibleLine][TransportMode='bus']/TransportSubmode["+invalidBusSubmodeCondition+"]")

	r.addRule("TRANSPORT_SUBMODE_BUS_INVALID_SERVICE_JOURNEY", "Invalid bus submode on ServiceJourney",
		"ServiceJourney with bus transport has invalid TransportSubmode", types.ERROR,
		"//vehicleJourneys/*[self::ServiceJourney or self::DatedServiceJourney][TransportMode='bus']/TransportSubmode["+invalidBusSubmodeCondition+"]")
}

// addRailSubModeRules validates rail transport submodes
func (r *RuleRegistry) addRailSubModeRules() {
	validRailSubmodes := []string{
		"local", "highSpeedRail", "suburbanRailway", "regionalRail", "interregionalRail",
		"longDistance", "international", "sleeperRailService", "nightRail", "carTransportRailService",
		"touristRailway", "railShuttle", "rackAndPinionRailway",
	}

	invalidRailSubmodeCondition := r.buildInvalidSubmodeCondition(validRailSubmodes)

	r.addRule("TRANSPORT_SUBMODE_RAIL_INVALID_LINE", "Invalid rail submode on Line",
		"Line with rail transport has invalid TransportSubmode", types.ERROR,
		"//lines/*[self::Line or self::FlexibleLine][TransportMode='rail']/TransportSubmode["+invalidRailSubmodeCondition+"]")

	r.addRule("TRANSPORT_SUBMODE_RAIL_INVALID_SERVICE_JOURNEY", "Invalid rail submode on ServiceJourney",
		"ServiceJourney with rail transport has invalid TransportSubmode", types.ERROR,
		"//vehicleJourneys/*[self::ServiceJourney or self::DatedServiceJourney][TransportMode='rail']/TransportSubmode["+invalidRailSubmodeCondition+"]")
}

// addTramSubModeRules validates tram transport submodes
func (r *RuleRegistry) addTramSubModeRules() {
	validTramSubmodes := []string{
		"cityTram", "localTram", "regionalTram", "sightseeingTram", "shuttleTram", "trainTram",
	}

	invalidTramSubmodeCondition := r.buildInvalidSubmodeCondition(validTramSubmodes)

	r.addRule("TRANSPORT_SUBMODE_TRAM_INVALID_LINE", "Invalid tram submode on Line",
		"Line with tram transport has invalid TransportSubmode", types.ERROR,
		"//lines/*[self::Line or self::FlexibleLine][TransportMode='tram']/TransportSubmode["+invalidTramSubmodeCondition+"]")

	r.addRule("TRANSPORT_SUBMODE_TRAM_INVALID_SERVICE_JOURNEY", "Invalid tram submode on ServiceJourney",
		"ServiceJourney with tram transport has invalid TransportSubmode", types.ERROR,
		"//vehicleJourneys/*[self::ServiceJourney or self::DatedServiceJourney][TransportMode='tram']/TransportSubmode["+invalidTramSubmodeCondition+"]")
}

// addMetroSubModeRules validates metro transport submodes
func (r *RuleRegistry) addMetroSubModeRules() {
	validMetroSubmodes := []string{
		"metro", "tube", "urbanRailway",
	}

	invalidMetroSubmodeCondition := r.buildInvalidSubmodeCondition(validMetroSubmodes)

	r.addRule("TRANSPORT_SUBMODE_METRO_INVALID_LINE", "Invalid metro submode on Line",
		"Line with metro transport has invalid TransportSubmode", types.ERROR,
		"//lines/*[self::Line or self::FlexibleLine][TransportMode='metro']/TransportSubmode["+invalidMetroSubmodeCondition+"]")
}

// addWaterSubModeRules validates water transport submodes
func (r *RuleRegistry) addWaterSubModeRules() {
	validWaterSubmodes := []string{
		"internationalCarFerry", "nationalCarFerry", "regionalCarFerry", "localCarFerry",
		"internationalPassengerFerry", "nationalPassengerFerry", "regionalPassengerFerry",
		"localPassengerFerry", "postBoat", "trainFerry", "roadFerryLink", "airportBoatLink",
		"highSpeedVehicleService", "highSpeedPassengerService", "sightseeingService",
		"schoolBoat", "cableFerry", "riverBus", "scheduledFerry", "shuttleFerryService",
	}

	invalidWaterSubmodeCondition := r.buildInvalidSubmodeCondition(validWaterSubmodes)

	r.addRule("TRANSPORT_SUBMODE_WATER_INVALID_LINE", "Invalid water submode on Line",
		"Line with water transport has invalid TransportSubmode", types.ERROR,
		"//lines/*[self::Line or self::FlexibleLine][TransportMode='water']/TransportSubmode["+invalidWaterSubmodeCondition+"]")
}

// addAirSubModeRules validates air transport submodes
func (r *RuleRegistry) addAirSubModeRules() {
	validAirSubmodes := []string{
		"internationalFlight", "domesticFlight", "intercontinentalFlight", "domesticScheduledFlight",
		"shuttleFlight", "intercontinentalCharterFlight", "internationalCharterFlight",
		"roundTripCharterFlight", "sightseeingFlight", "helicopterService", "domesticCharterFlight",
		"SchengenAreaFlight", "airshipService", "shortHaulInternationalFlight",
	}

	invalidAirSubmodeCondition := r.buildInvalidSubmodeCondition(validAirSubmodes)

	r.addRule("TRANSPORT_SUBMODE_AIR_INVALID_LINE", "Invalid air submode on Line",
		"Line with air transport has invalid TransportSubmode", types.ERROR,
		"//lines/*[self::Line or self::FlexibleLine][TransportMode='air']/TransportSubmode["+invalidAirSubmodeCondition+"]")
}

// addCoachSubModeRules validates coach transport submodes
func (r *RuleRegistry) addCoachSubModeRules() {
	validCoachSubmodes := []string{
		"internationalCoach", "nationalCoach", "shuttleCoach", "regionalCoach", "specialCoach",
		"schoolCoach", "sightseeingCoach", "touristCoach", "commuterCoach",
	}

	invalidCoachSubmodeCondition := r.buildInvalidSubmodeCondition(validCoachSubmodes)

	r.addRule("TRANSPORT_SUBMODE_COACH_INVALID_LINE", "Invalid coach submode on Line",
		"Line with coach transport has invalid TransportSubmode", types.ERROR,
		"//lines/*[self::Line or self::FlexibleLine][TransportMode='coach']/TransportSubmode["+invalidCoachSubmodeCondition+"]")
}

// addOtherTransportSubModeRules validates other transport submodes
func (r *RuleRegistry) addOtherTransportSubModeRules() {
	// Taxi submodes
	validTaxiSubmodes := []string{
		"communalTaxi", "waterTaxi", "railTaxi", "bikeTaxi", "licensedTaxi", "privateTaxi",
		"allTaxiServices",
	}

	invalidTaxiSubmodeCondition := r.buildInvalidSubmodeCondition(validTaxiSubmodes)

	r.addRule("TRANSPORT_SUBMODE_TAXI_INVALID_LINE", "Invalid taxi submode on Line",
		"Line with taxi transport has invalid TransportSubmode", types.ERROR,
		"//lines/*[self::Line or self::FlexibleLine][TransportMode='taxi']/TransportSubmode["+invalidTaxiSubmodeCondition+"]")

	// Cableway submodes
	validCablewaySubmodes := []string{
		"telecabin", "cableCar", "lift", "chairLift", "dragLift", "smallTelecabin",
	}

	invalidCablewaySubmodeCondition := r.buildInvalidSubmodeCondition(validCablewaySubmodes)

	r.addRule("TRANSPORT_SUBMODE_CABLEWAY_INVALID_LINE", "Invalid cableway submode on Line",
		"Line with cableway transport has invalid TransportSubmode", types.ERROR,
		"//lines/*[self::Line or self::FlexibleLine][TransportMode='cableway']/TransportSubmode["+invalidCablewaySubmodeCondition+"]")

	// Funicular submodes
	r.addRule("TRANSPORT_SUBMODE_FUNICULAR_INVALID_LINE", "Invalid funicular submode on Line",
		"Line with funicular transport has invalid TransportSubmode", types.ERROR,
		"//lines/*[self::Line or self::FlexibleLine][TransportMode='funicular']/TransportSubmode[not(text()='funicular')]")
}

// addContextualTransportModeRules adds rules that validate transport modes in context
func (r *RuleRegistry) addContextualTransportModeRules() {
	// ServiceJourney transport mode should be compatible with Line transport mode
	r.addRule("TRANSPORT_MODE_INCOMPATIBLE_SERVICE_JOURNEY", "Incompatible transport modes",
		"ServiceJourney transport mode incompatible with Line transport mode", types.ERROR,
		"//vehicleJourneys/*[self::ServiceJourney or self::DatedServiceJourney][TransportMode and TransportMode != //lines/*[self::Line or self::FlexibleLine][@id=current()/LineRef/@ref]/TransportMode]")

	// Line missing transport mode when required
	r.addRule("TRANSPORT_MODE_MISSING_REQUIRED", "Missing required transport mode",
		"TransportMode is required for all Lines and FlexibleLines", types.ERROR,
		"//lines/*[self::Line or self::FlexibleLine][not(TransportMode)]")

	// Transport submode without transport mode
	r.addRule("TRANSPORT_SUBMODE_WITHOUT_MODE", "Transport submode without mode",
		"TransportSubmode specified without TransportMode", types.WARNING,
		"//*/TransportSubmode[not(../TransportMode)]")

	// Inconsistent transport mode hierarchy (Line vs ServiceJourney vs Route)
	r.addRule("TRANSPORT_MODE_HIERARCHY_INCONSISTENT", "Inconsistent transport mode hierarchy",
		"Transport mode inconsistent between Line, Route, and ServiceJourney", types.WARNING,
		"//routes/Route[TransportMode and TransportMode != //lines/*[self::Line or self::FlexibleLine][@id=current()/LineRef/@ref]/TransportMode]")
}

// buildInvalidSubmodeCondition builds XPath condition for invalid submodes
func (r *RuleRegistry) buildInvalidSubmodeCondition(validSubmodes []string) string {
	if len(validSubmodes) == 0 {
		return "true()"
	}

	condition := "not(text() = '" + validSubmodes[0] + "'"
	for _, submode := range validSubmodes[1:] {
		condition += " or text() = '" + submode + "'"
	}
	condition += ")"

	return condition
}
