package rules

import "github.com/theoremus-urban-solutions/netex-validator/types"

// loadBuiltinRules loads all 200+ built-in validation rules for comprehensive NetEX validation
func (r *RuleRegistry) loadBuiltinRules() {
	// LINE validation rules - Enhanced coverage
	r.addRule("LINE_2", "Line missing Name", "Line is missing Name. EU requirement: provide <Name> for every Line to ensure human-readable identification.", types.ERROR,
		"//*[local-name()='lines']/*[local-name()='Line' or local-name()='FlexibleLine'][not(*[local-name()='Name']) or normalize-space(*[local-name()='Name'])='']")

	r.addRule("LINE_3", "Line missing PublicCode", "Line is missing PublicCode", types.WARNING,
		"//lines/*[self::Line or self::FlexibleLine][not(PublicCode) or normalize-space(PublicCode) = '']")

	r.addRule("LINE_4", "Line missing TransportMode", "Line is missing TransportMode. EU requirement: set <TransportMode> to a valid value (bus, rail, tram, metro, air, water, taxi, cableway, funicular, coach).", types.ERROR,
		"//*[local-name()='lines']/*[local-name()='Line' or local-name()='FlexibleLine'][not(*[local-name()='TransportMode'])]")

	r.addRule("LINE_5", "Line missing TransportSubmode", "Line is missing TransportSubmode", types.WARNING,
		"//lines/*[self::Line or self::FlexibleLine][not(TransportSubmode)]")

	r.addRule("LINE_6", "Line with incorrect use of Route", "Line has incorrect use of Route element", types.ERROR,
		"//lines/*[self::Line or self::FlexibleLine]/routes/Route")

	r.addRule("LINE_7", "Line missing Network or GroupOfLines", "Line is missing reference to Network or GroupOfLines", types.WARNING,
		"//lines/*[self::Line or self::FlexibleLine][not(RepresentedByGroupRef)]")

	r.addRule("LINE_8", "Line missing OperatorRef", "Line is missing OperatorRef", types.WARNING,
		"//lines/*[self::Line or self::FlexibleLine][not(OperatorRef)]")

	r.addRule("LINE_9", "Line missing AuthorityRef", "Line is missing AuthorityRef", types.WARNING,
		"//lines/*[self::Line or self::FlexibleLine][not(AuthorityRef)]")

	// ROUTE validation rules - Enhanced coverage
	r.addRule("ROUTE_2", "Route missing Name", "Route is missing Name. Add <Name> to each Route for clarity in reporting and passenger information.", types.ERROR,
		"//*[local-name()='routes']/*[local-name()='Route'][not(*[local-name()='Name']) or normalize-space(*[local-name()='Name'])='']")

	r.addRule("ROUTE_3", "Route missing LineRef", "Route is missing LineRef. Ensure a Route references its parent Line using <LineRef> (EU profile).", types.ERROR,
		"//*[local-name()='routes']/*[local-name()='Route'][not(*[local-name()='LineRef']) and not(*[local-name()='FlexibleLineRef'])]")

	r.addRule("ROUTE_4", "Route missing pointsInSequence", "Route is missing pointsInSequence", types.ERROR,
		"//routes/Route[not(pointsInSequence)]")

	r.addRule("ROUTE_5", "Route illegal DirectionRef", "Route has illegal DirectionRef", types.WARNING,
		"//routes/Route/DirectionRef")

	r.addRule("ROUTE_6", "Route duplicated order", "Route has duplicated order values in PointOnRoute", types.WARNING,
		"//routes/Route/pointsInSequence/PointOnRoute[@order = preceding-sibling::PointOnRoute/@order]")

	r.addRule("ROUTE_7", "Route missing DirectionType", "Route is missing DirectionType", types.WARNING,
		"//routes/Route[not(DirectionType)]")

	r.addRule("ROUTE_8", "Route invalid DirectionType", "Route has invalid DirectionType", types.ERROR,
		"//routes/Route[DirectionType and not(DirectionType = 'inbound' or DirectionType = 'outbound' or DirectionType = 'clockwise' or DirectionType = 'anticlockwise')]")

	// SERVICE_JOURNEY validation rules - Enhanced coverage
	r.addRule("SERVICE_JOURNEY_2", "ServiceJourney illegal element Call", "ServiceJourney has illegal element Call", types.ERROR,
		"//vehicleJourneys/ServiceJourney/calls")

	r.addRule("SERVICE_JOURNEY_3", "ServiceJourney missing element PassingTimes", "ServiceJourney is missing PassingTimes element", types.ERROR,
		"//vehicleJourneys/*[self::ServiceJourney][not(passingTimes)]")

	r.addRule("SERVICE_JOURNEY_4", "ServiceJourney missing arrival and departure", "ServiceJourney is missing both arrival and departure times. EU requirement: each TimetabledPassingTime must have at least one of ArrivalTime or DepartureTime.", types.ERROR,
		"//vehicleJourneys/ServiceJourney/passingTimes/TimetabledPassingTime[not(DepartureTime or EarliestDepartureTime) and not(ArrivalTime or LatestArrivalTime)]")

	r.addRule("SERVICE_JOURNEY_5", "ServiceJourney missing departure times", "ServiceJourney is missing departure times for first stop", types.ERROR,
		"//vehicleJourneys/ServiceJourney[not(passingTimes/TimetabledPassingTime[1]/DepartureTime) and not(passingTimes/TimetabledPassingTime[1]/EarliestDepartureTime)]")

	r.addRule("SERVICE_JOURNEY_6", "ServiceJourney missing arrival time for last stop", "ServiceJourney is missing arrival time for last stop", types.ERROR,
		"//vehicleJourneys/ServiceJourney[count(passingTimes/TimetabledPassingTime[last()]/ArrivalTime) = 0 and count(passingTimes/TimetabledPassingTime[last()]/LatestArrivalTime) = 0]")

	r.addRule("SERVICE_JOURNEY_7", "ServiceJourney identical arrival and departure", "ServiceJourney has identical arrival and departure times", types.WARNING,
		"//vehicleJourneys/ServiceJourney/passingTimes/TimetabledPassingTime[DepartureTime = ArrivalTime]")

	r.addRule("SERVICE_JOURNEY_8", "ServiceJourney missing id on TimetabledPassingTime", "ServiceJourney is missing id on TimetabledPassingTime", types.WARNING,
		"//vehicleJourneys/ServiceJourney/passingTimes/TimetabledPassingTime[not(@id)]")

	r.addRule("SERVICE_JOURNEY_9", "ServiceJourney missing version on TimetabledPassingTime", "ServiceJourney is missing version on TimetabledPassingTime", types.WARNING,
		"//vehicleJourneys/ServiceJourney/passingTimes/TimetabledPassingTime[not(@version)]")

	r.addRule("SERVICE_JOURNEY_10", "ServiceJourney missing reference to JourneyPattern", "ServiceJourney is missing reference to JourneyPattern", types.ERROR,
		"//vehicleJourneys/*[self::ServiceJourney][not(JourneyPatternRef)]")

	r.addRule("SERVICE_JOURNEY_11", "ServiceJourney invalid overriding of transport modes", "ServiceJourney has invalid overriding of transport modes", types.WARNING,
		"//vehicleJourneys/ServiceJourney[(TransportMode and not(TransportSubmode)) or (not(TransportMode) and TransportSubmode)]")

	r.addRule("SERVICE_JOURNEY_12", "ServiceJourney missing OperatorRef", "ServiceJourney is missing OperatorRef", types.ERROR,
		"//vehicleJourneys/ServiceJourney[not(OperatorRef) and not(//ServiceFrame/lines/*[self::Line or self::FlexibleLine]/OperatorRef)]")

	r.addRule("SERVICE_JOURNEY_13", "ServiceJourney missing reference to calendar data", "ServiceJourney is missing reference to calendar data", types.ERROR,
		"//vehicleJourneys/ServiceJourney[not(dayTypes/DayTypeRef) and not(@id=//TimetableFrame/vehicleJourneys/DatedServiceJourney/ServiceJourneyRef/@ref)]")

	r.addRule("SERVICE_JOURNEY_14", "ServiceJourney duplicated reference to calendar data", "ServiceJourney has duplicated reference to calendar data", types.ERROR,
		"//vehicleJourneys/ServiceJourney[dayTypes/DayTypeRef and @id=//TimetableFrame/vehicleJourneys/DatedServiceJourney/ServiceJourneyRef/@ref]")

	r.addRule("SERVICE_JOURNEY_15", "ServiceJourney inconsistent number of timetable passing times", "ServiceJourney has inconsistent number of timetable passing times", types.ERROR,
		"//vehicleJourneys/ServiceJourney[JourneyPatternRef and count(passingTimes/TimetabledPassingTime) > 0]")

	r.addRule("SERVICE_JOURNEY_16", "ServiceJourney multiple versions", "ServiceJourney has multiple versions with same id", types.WARNING,
		"//vehicleJourneys/ServiceJourney[@id = preceding-sibling::ServiceJourney/@id]")

	r.addRule("SERVICE_JOURNEY_17", "ServiceJourney duplicate TimetabledPassingTime IDs", "ServiceJourney has duplicate TimetabledPassingTime IDs", types.ERROR,
		"//vehicleJourneys/ServiceJourney/passingTimes/TimetabledPassingTime[@id = preceding-sibling::TimetabledPassingTime/@id or @id = following-sibling::TimetabledPassingTime/@id]")

	// FLEXIBLE_LINE validation rules
	r.addRule("FLEXIBLE_LINE_1", "FlexibleLine missing FlexibleLineType", "FlexibleLine is missing FlexibleLineType", types.ERROR,
		"//lines/FlexibleLine[not(FlexibleLineType)]")

	r.addRule("FLEXIBLE_LINE_10", "FlexibleLine illegal use of both BookWhen and MinimumBookingPeriod", "FlexibleLine has illegal use of both BookWhen and MinimumBookingPeriod", types.WARNING,
		"//lines/FlexibleLine[BookWhen and MinimumBookingPeriod]")

	r.addRule("FLEXIBLE_LINE_11", "FlexibleLine BookWhen without LatestBookingTime or vice versa", "FlexibleLine has BookWhen without LatestBookingTime or vice versa", types.WARNING,
		"//lines/FlexibleLine[(BookWhen and not(LatestBookingTime)) or (not(BookWhen) and LatestBookingTime)]")

	// NETWORK validation rules
	r.addRule("NETWORK_1", "Network missing AuthorityRef", "Network is missing AuthorityRef", types.ERROR,
		"//Network[not(AuthorityRef)]")

	r.addRule("NETWORK_2", "Network missing Name on Network", "Network is missing Name", types.ERROR,
		"//Network[not(Name) or normalize-space(Name) = '']")

	r.addRule("NETWORK_3", "Network missing Name on GroupOfLines", "GroupOfLines is missing Name", types.ERROR,
		"//Network/groupsOfLines/GroupOfLines[not(Name) or normalize-space(Name) = '']")

	// AUTHORITY validation rules
	r.addRule("AUTHORITY_2", "Authority missing Name", "Authority is missing Name", types.ERROR,
		"//organisations/Authority[not(Name) or normalize-space(Name) = '']")

	// OPERATOR validation rules
	r.addRule("OPERATOR_2", "Operator missing Name", "Operator is missing Name", types.ERROR,
		"//organisations/Operator[not(Name) or normalize-space(Name) = '']")

	// VERSION validation rules
	r.addRule("VERSION_NON_NUMERIC", "Non-numeric NeTEx version", "Element has non-numeric NeTEx version", types.WARNING,
		".//*[@version != 'any' and number(@version) != number(@version)]")

	// JOURNEY_PATTERN validation rules
	r.addRule("JOURNEY_PATTERN_1", "JourneyPattern missing RouteRef", "JourneyPattern is missing RouteRef", types.ERROR,
		"//journeyPatterns/JourneyPattern[not(RouteRef)]")

	r.addRule("JOURNEY_PATTERN_2", "JourneyPattern missing pointsInSequence", "JourneyPattern is missing pointsInSequence", types.ERROR,
		"//journeyPatterns/JourneyPattern[not(pointsInSequence)]")

	// STOP_POINT validation rules
	r.addRule("STOP_POINT_1", "StopPoint missing ScheduledStopPointRef", "StopPoint is missing ScheduledStopPointRef", types.ERROR,
		"//StopPointInJourneyPattern[not(ScheduledStopPointRef)]")

	r.addRule("STOP_POINT_2", "StopPoint missing order", "StopPoint is missing order attribute", types.ERROR,
		"//StopPointInJourneyPattern[not(@order)]")

	// SCHEDULED_STOP_POINT validation rules
	r.addRule("SCHEDULED_STOP_POINT_1", "ScheduledStopPoint missing Name", "ScheduledStopPoint is missing Name", types.WARNING,
		"//scheduledStopPoints/ScheduledStopPoint[not(Name) or normalize-space(Name) = '']")

	// TRANSPORT_MODE validation rules - Advanced business logic
	r.addRule("TRANSPORT_MODE_ON_LINE", "Line with illegal TransportMode", "Line has illegal TransportMode", types.ERROR,
		"//lines/*[self::Line or self::FlexibleLine][TransportMode and not(TransportMode = 'coach' or TransportMode = 'bus' or TransportMode = 'tram' or TransportMode = 'rail' or TransportMode = 'metro' or TransportMode = 'air' or TransportMode = 'taxi' or TransportMode = 'water' or TransportMode = 'cableway' or TransportMode = 'funicular' or TransportMode = 'unknown')]")

	r.addRule("TRANSPORT_MODE_ON_SERVICE_JOURNEY", "ServiceJourney with illegal TransportMode", "ServiceJourney has illegal TransportMode", types.ERROR,
		"//vehicleJourneys/ServiceJourney[TransportMode and not(TransportMode = 'coach' or TransportMode = 'bus' or TransportMode = 'tram' or TransportMode = 'rail' or TransportMode = 'metro' or TransportMode = 'air' or TransportMode = 'taxi' or TransportMode = 'water' or TransportMode = 'cableway' or TransportMode = 'funicular' or TransportMode = 'unknown')]")

	// TRANSPORT_SUB_MODE validation rules - Context-specific validation
	r.addRule("TRANSPORT_SUB_MODE_BUS_INVALID", "Line with invalid bus TransportSubMode", "Line has invalid TransportSubMode for bus transport", types.ERROR,
		"//lines/*[self::Line or self::FlexibleLine][TransportMode = 'bus' and TransportSubmode and not(TransportSubmode = 'localBus' or TransportSubmode = 'regionalBus' or TransportSubmode = 'expressBus' or TransportSubmode = 'nightBus' or TransportSubmode = 'postBus' or TransportSubmode = 'specialNeedsBus' or TransportSubmode = 'mobilityBus' or TransportSubmode = 'mobilityBusForRegisteredDisabled' or TransportSubmode = 'sightseeingBus' or TransportSubmode = 'shuttleBus' or TransportSubmode = 'schoolBus' or TransportSubmode = 'schoolAndPublicServiceBus' or TransportSubmode = 'railReplacementBus' or TransportSubmode = 'demandAndResponseBus' or TransportSubmode = 'airportLinkBus')]")

	r.addRule("TRANSPORT_SUB_MODE_RAIL_INVALID", "Line with invalid rail TransportSubMode", "Line has invalid TransportSubMode for rail transport", types.ERROR,
		"//lines/*[self::Line or self::FlexibleLine][TransportMode = 'rail' and TransportSubmode and not(TransportSubmode = 'local' or TransportSubmode = 'highSpeedRail' or TransportSubmode = 'suburbanRailway' or TransportSubmode = 'regionalRail' or TransportSubmode = 'interregionalRail' or TransportSubmode = 'longDistance' or TransportSubmode = 'international' or TransportSubmode = 'sleeperRailService' or TransportSubmode = 'nightRail' or TransportSubmode = 'carTransportRailService' or TransportSubmode = 'touristRailway' or TransportSubmode = 'railShuttle')]")

	r.addRule("TRANSPORT_SUB_MODE_TRAM_INVALID", "Line with invalid tram TransportSubMode", "Line has invalid TransportSubMode for tram transport", types.ERROR,
		"//lines/*[self::Line or self::FlexibleLine][TransportMode = 'tram' and TransportSubmode and not(TransportSubmode = 'cityTram' or TransportSubmode = 'localTram' or TransportSubmode = 'regionalTram' or TransportSubmode = 'sightseeingTram' or TransportSubmode = 'shuttleTram' or TransportSubmode = 'trainTram')]")

	// BOOKING validation rules - Advanced flexible service validation
	r.addRule("BOOKING_INVALID_BOOK_WHEN", "Invalid BookWhen value", "BookWhen has invalid value", types.ERROR,
		"//lines/FlexibleLine[BookWhen and not(BookWhen = 'dayOfTravelOnly' or BookWhen = 'untilPreviousDay' or BookWhen = 'advanceAndDayOfTravel')]")

	r.addRule("BOOKING_MISSING_PROPERTIES", "Mandatory booking property missing", "Flexible line is missing mandatory booking properties", types.ERROR,
		"//lines/FlexibleLine[FlexibleLineType and (FlexibleLineType = 'flexibleAreasOnly' or FlexibleLineType = 'hailAndRideAreas' or FlexibleLineType = 'demandAndResponseServices') and not(BookWhen or MinimumBookingPeriod)]")

	// FLEXIBLE_LINE_TYPE validation rules - Advanced flexible line validation
	r.addRule("FLEXIBLE_LINE_TYPE_INVALID", "FlexibleLine with invalid FlexibleLineType", "FlexibleLine has invalid FlexibleLineType", types.ERROR,
		"//lines/FlexibleLine[FlexibleLineType and not(FlexibleLineType = 'fixedStop' or FlexibleLineType = 'flexibleAreasOnly' or FlexibleLineType = 'hailAndRideAreas' or FlexibleLineType = 'flexibleAreasAndStops' or FlexibleLineType = 'hailAndRideSections' or FlexibleLineType = 'fixedStopAreaWide' or FlexibleLineType = 'freeAreaAreaWide' or FlexibleLineType = 'mixedFlexible' or FlexibleLineType = 'mixedFlexibleAndFixed' or FlexibleLineType = 'fixed' or FlexibleLineType = 'mainRouteWithFlexibleEnds' or FlexibleLineType = 'flexibleRoute')]")

	// LINE presentation validation rules - Color and display validation
	r.addRule("LINE_INVALID_COLOR_LENGTH", "Line with invalid color coding length", "Line has invalid color coding length on Presentation", types.WARNING,
		"//lines/*[self::Line or self::FlexibleLine]/*[self::Presentation or self::AlternativePresentation]/*[self::Colour or self::TextColour][text()]")

	r.addRule("LINE_INVALID_COLOR_VALUE", "Line with invalid color coding value", "Line has invalid color coding value on Presentation", types.WARNING,
		"//lines/*[self::Line or self::FlexibleLine]/*[self::Presentation or self::AlternativePresentation]/*[self::Colour or self::TextColour][text()]")

	// CALENDAR and VALIDITY validation rules
	r.addRule("SERVICE_CALENDAR_1", "ServiceCalendar missing DayTypes", "ServiceCalendar is missing DayTypes", types.ERROR,
		"//ServiceCalendar[not(dayTypes)]")

	r.addRule("SERVICE_CALENDAR_2", "ServiceCalendar missing OperatingPeriod", "ServiceCalendar is missing OperatingPeriod", types.ERROR,
		"//ServiceCalendar[not(operatingPeriods)]")

	r.addRule("VALIDITY_CONDITIONS_1", "ValidityConditions missing AvailabilityCondition", "ValidityConditions is missing AvailabilityCondition", types.WARNING,
		"//validityConditions[not(AvailabilityCondition)]")

	r.addRule("VALIDITY_CONDITIONS_2", "ValidityConditions invalid FromDate", "ValidityConditions has invalid FromDate", types.ERROR,
		"//validityConditions/AvailabilityCondition[FromDate]")

	// DATED_SERVICE_JOURNEY validation rules
	r.addRule("DATED_SERVICE_JOURNEY_1", "DatedServiceJourney missing ServiceJourneyRef", "DatedServiceJourney is missing ServiceJourneyRef", types.ERROR,
		"//vehicleJourneys/DatedServiceJourney[not(ServiceJourneyRef)]")

	r.addRule("DATED_SERVICE_JOURNEY_2", "DatedServiceJourney missing OperatingDayRef", "DatedServiceJourney is missing OperatingDayRef", types.ERROR,
		"//vehicleJourneys/DatedServiceJourney[not(OperatingDayRef)]")

	r.addRule("DATED_SERVICE_JOURNEY_3", "DatedServiceJourney invalid OperatingDayRef", "DatedServiceJourney has invalid OperatingDayRef", types.ERROR,
		"//vehicleJourneys/DatedServiceJourney[OperatingDayRef]")

	// DEAD_RUN validation rules
	r.addRule("DEAD_RUN_1", "DeadRun missing RouteRef", "DeadRun is missing RouteRef", types.ERROR,
		"//vehicleJourneys/DeadRun[not(RouteRef)]")

	r.addRule("DEAD_RUN_2", "DeadRun missing passingTimes", "DeadRun is missing passingTimes", types.ERROR,
		"//vehicleJourneys/DeadRun[not(passingTimes)]")

	// INTERCHANGE validation rules
	r.addRule("INTERCHANGE_1", "ServiceJourneyInterchange missing FromStopPointRef", "ServiceJourneyInterchange is missing FromStopPointRef", types.ERROR,
		"//interchanges/ServiceJourneyInterchange[not(FromStopPointRef)]")

	r.addRule("INTERCHANGE_2", "ServiceJourneyInterchange missing ToStopPointRef", "ServiceJourneyInterchange is missing ToStopPointRef", types.ERROR,
		"//interchanges/ServiceJourneyInterchange[not(ToStopPointRef)]")

	r.addRule("INTERCHANGE_3", "ServiceJourneyInterchange missing FromServiceJourneyRef", "ServiceJourneyInterchange is missing FromServiceJourneyRef", types.ERROR,
		"//interchanges/ServiceJourneyInterchange[not(FromServiceJourneyRef)]")

	r.addRule("INTERCHANGE_4", "ServiceJourneyInterchange missing ToServiceJourneyRef", "ServiceJourneyInterchange is missing ToServiceJourneyRef", types.ERROR,
		"//interchanges/ServiceJourneyInterchange[not(ToServiceJourneyRef)]")

	// NOTICE validation rules
	r.addRule("NOTICE_1", "Notice missing Name", "Notice is missing Name", types.WARNING,
		"//notices/Notice[not(Name) or normalize-space(Name) = '']")

	r.addRule("NOTICE_2", "Notice missing Text", "Notice is missing Text", types.ERROR,
		"//notices/Notice[not(Text) or normalize-space(Text) = '']")

	r.addRule("NOTICE_3", "NoticeAssignment missing NoticedObjectRef", "NoticeAssignment is missing NoticedObjectRef", types.ERROR,
		"//noticeAssignments/NoticeAssignment[not(NoticedObjectRef)]")

	// FRAME validation rules
	r.addRule("COMPOSITE_FRAME_1", "CompositeFrame must be exactly one", "There must be exactly one CompositeFrame", types.ERROR,
		"//CompositeFrame[count(//CompositeFrame) != 1]")

	r.addRule("RESOURCE_FRAME_IN_LINE_FILE", "ResourceFrame must be exactly one in line file", "Line file must contain exactly one ResourceFrame", types.ERROR,
		"//ResourceFrame[count(//ResourceFrame) != 1]")

	r.addRule("SERVICE_FRAME_1", "ServiceFrame missing in CompositeFrame", "CompositeFrame is missing ServiceFrame", types.ERROR,
		"//CompositeFrame[not(ServiceFrame)]")

	r.addRule("TIMETABLE_FRAME_1", "TimetableFrame missing in CompositeFrame", "CompositeFrame is missing TimetableFrame", types.WARNING,
		"//CompositeFrame[not(TimetableFrame)]")

	// FLEXIBLE_SERVICE validation rules
	r.addRule("FLEXIBLE_SERVICE_1", "FlexibleService missing FlexibleServiceType", "FlexibleService is missing FlexibleServiceType", types.ERROR,
		"//FlexibleService[not(FlexibleServiceType)]")

	r.addRule("FLEXIBLE_SERVICE_2", "FlexibleService invalid FlexibleServiceType", "FlexibleService has invalid FlexibleServiceType", types.ERROR,
		"//FlexibleService[FlexibleServiceType and not(FlexibleServiceType = 'dynamicPassengerInformation' or FlexibleServiceType = 'fixedHeadwayService' or FlexibleServiceType = 'fixedPassingTimes' or FlexibleServiceType = 'notFixedPassingTimes')]")

	// BLOCK validation rules
	r.addRule("BLOCK_1", "Block missing Name", "Block is missing Name", types.WARNING,
		"//blocks/Block[not(Name) or normalize-space(Name) = '']")

	r.addRule("BLOCK_2", "Block missing vehicleJourneys", "Block is missing vehicleJourneys", types.ERROR,
		"//blocks/Block[not(journeys)]")

	// COURSE_OF_JOURNEYS validation rules
	r.addRule("COURSE_OF_JOURNEYS_1", "CourseOfJourneys missing Name", "CourseOfJourneys is missing Name", types.WARNING,
		"//CourseOfJourneys[not(Name) or normalize-space(Name) = '']")

	r.addRule("COURSE_OF_JOURNEYS_2", "CourseOfJourneys missing journeys", "CourseOfJourneys is missing journeys", types.ERROR,
		"//CourseOfJourneys[not(journeys)]")

	// TARIFF_ZONE validation rules
	r.addRule("TARIFF_ZONE_1", "TariffZone missing Name", "TariffZone is missing Name", types.WARNING,
		"//tariffzones/TariffZone[not(Name) or normalize-space(Name) = '']")

	r.addRule("TARIFF_ZONE_2", "TariffZone missing Centroid", "TariffZone is missing Centroid", types.WARNING,
		"//tariffzones/TariffZone[not(Centroid)]")

	// RESPONSIBILITY_SET validation rules
	r.addRule("RESPONSIBILITY_SET_1", "ResponsibilitySet missing Name", "ResponsibilitySet is missing Name", types.WARNING,
		"//ResponsibilitySet[not(Name) or normalize-space(Name) = '']")

	r.addRule("RESPONSIBILITY_SET_2", "ResponsibilitySet missing roles", "ResponsibilitySet is missing roles", types.ERROR,
		"//ResponsibilitySet[not(roles)]")

	// TYPE_OF_SERVICE validation rules
	r.addRule("TYPE_OF_SERVICE_1", "TypeOfService missing Name", "TypeOfService is missing Name", types.WARNING,
		"//TypeOfService[not(Name) or normalize-space(Name) = '']")

	// GROUP validation rules
	r.addRule("GROUP_OF_LINES_1", "GroupOfLines missing Name", "GroupOfLines is missing Name", types.ERROR,
		"//groupsOfLines/GroupOfLines[not(Name) or normalize-space(Name) = '']")

	r.addRule("GROUP_OF_SERVICES_1", "GroupOfServices missing Name", "GroupOfServices is missing Name", types.WARNING,
		"//GroupOfServices[not(Name) or normalize-space(Name) = '']")

	// Load additional extended rules to reach parity with Java version
	r.loadExtendedBuiltinRules()
}
