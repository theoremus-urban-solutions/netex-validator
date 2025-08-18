package netexvalidator

// ValidationRulesDocumentation provides comprehensive documentation for all NetEX validation rules.
//
// This package implements over 200 validation rules covering:
// - Schema validation against NetEX XSD
// - Business logic rules for Nordic NeTEx Profile
// - Data consistency and integrity checks
// - Cross-references validation
// - Transport mode and submode validation
// - Flexible service and booking validation
// - Timetable and journey pattern validation
// - Stop place and infrastructure validation
// - Calendar and validity period validation
// - Fare and pricing validation
// - Interchange and connection validation

// LINE Validation Rules
//
// These rules validate Line and FlexibleLine elements, which represent public transport lines
// or services. Lines are the core organizational unit for public transport services.

const LineValidationRulesDoc = `
LINE VALIDATION RULES

LINE_2: Line missing Name
  Severity: ERROR
  Description: Line elements must have a Name element with non-empty content
  XPath: //lines/*[self::Line or self::FlexibleLine][not(Name) or normalize-space(Name) = '']
  Example violation:
    <Line id="NO:Line:1" version="1">
      <PublicCode>1</PublicCode>
      <!-- Missing Name element -->
    </Line>
  Fix: Add Name element: <Name>Line 1</Name>

LINE_3: Line missing PublicCode
  Severity: WARNING
  Description: Lines should have a PublicCode for passenger information
  XPath: //lines/*[self::Line or self::FlexibleLine][not(PublicCode) or normalize-space(PublicCode) = '']
  Example violation:
    <Line id="NO:Line:1" version="1">
      <Name>Blue Line</Name>
      <!-- Missing PublicCode -->
    </Line>
  Fix: Add PublicCode: <PublicCode>1</PublicCode>

LINE_4: Line missing TransportMode
  Severity: ERROR
  Description: Lines must specify the mode of transport (bus, rail, tram, etc.)
  XPath: //lines/*[self::Line or self::FlexibleLine][not(TransportMode)]
  Example violation:
    <Line id="NO:Line:1" version="1">
      <Name>City Bus</Name>
      <!-- Missing TransportMode -->
    </Line>
  Fix: Add TransportMode: <TransportMode>bus</TransportMode>

LINE_5: Line missing TransportSubmode
  Severity: WARNING
  Description: Lines should specify transport submode for more detailed classification
  XPath: //lines/*[self::Line or self::FlexibleLine][not(TransportSubmode)]
  Example violation:
    <Line id="NO:Line:1" version="1">
      <Name>City Bus</Name>
      <TransportMode>bus</TransportMode>
      <!-- Missing TransportSubmode -->
    </Line>
  Fix: Add TransportSubmode: <TransportSubmode>localBus</TransportSubmode>

LINE_6: Line with incorrect use of Route
  Severity: ERROR
  Description: Line should not contain Route elements directly (use ServiceFrame/routes)
  XPath: //lines/*[self::Line or self::FlexibleLine]/routes/Route
  Example violation:
    <Line id="NO:Line:1" version="1">
      <routes>
        <Route><!-- Routes should be in ServiceFrame --></Route>
      </routes>
    </Line>
  Fix: Move routes to ServiceFrame/routes

LINE_7: Line missing Network or GroupOfLines
  Severity: ERROR  
  Description: Lines must be associated with a Network through RepresentedByGroupRef
  XPath: //lines/*[self::Line or self::FlexibleLine][not(RepresentedByGroupRef)]
  Example violation:
    <Line id="NO:Line:1" version="1">
      <Name>City Bus</Name>
      <!-- Missing RepresentedByGroupRef -->
    </Line>
  Fix: Add RepresentedByGroupRef: <RepresentedByGroupRef ref="NO:Network:1"/>

LINE_8: Line missing OperatorRef
  Severity: ERROR
  Description: Lines must reference the operating company/authority
  XPath: //lines/*[self::Line or self::FlexibleLine][not(OperatorRef)]
  Example violation:
    <Line id="NO:Line:1" version="1">
      <Name>City Bus</Name>
      <!-- Missing OperatorRef -->
    </Line>
  Fix: Add OperatorRef: <OperatorRef ref="NO:Operator:CityBus"/>

LINE_9: Line missing AuthorityRef
  Severity: WARNING
  Description: Lines should reference the transport authority
  XPath: //lines/*[self::Line or self::FlexibleLine][not(AuthorityRef)]
  Example violation:
    <Line id="NO:Line:1" version="1">
      <Name>City Bus</Name>
      <!-- Missing AuthorityRef -->
    </Line>
  Fix: Add AuthorityRef: <AuthorityRef ref="NO:Authority:Oslo"/>
`

// ROUTE Validation Rules
//
// Route elements define the geographical path and sequence of stops for a Line.
// Routes connect Lines to stop sequences and journey patterns.

const RouteValidationRulesDoc = `
ROUTE VALIDATION RULES

ROUTE_2: Route missing Name
  Severity: ERROR
  Description: Route elements must have a descriptive name
  XPath: //routes/Route[not(Name) or normalize-space(Name) = '']
  Example violation:
    <Route id="NO:Route:1" version="1">
      <!-- Missing Name element -->
      <LineRef ref="NO:Line:1"/>
    </Route>
  Fix: Add Name: <Name>City Center - Airport</Name>

ROUTE_3: Route missing LineRef
  Severity: ERROR
  Description: Routes must reference the Line they belong to
  XPath: //routes/Route[not(LineRef) and not(FlexibleLineRef)]
  Example violation:
    <Route id="NO:Route:1" version="1">
      <Name>Downtown Route</Name>
      <!-- Missing LineRef -->
    </Route>
  Fix: Add LineRef: <LineRef ref="NO:Line:1"/>

ROUTE_4: Route missing pointsInSequence
  Severity: ERROR
  Description: Routes must define the sequence of points (stops) along the route
  XPath: //routes/Route[not(pointsInSequence)]
  Example violation:
    <Route id="NO:Route:1" version="1">
      <Name>Downtown Route</Name>
      <LineRef ref="NO:Line:1"/>
      <!-- Missing pointsInSequence -->
    </Route>
  Fix: Add pointsInSequence with PointOnRoute elements

ROUTE_5: Route illegal DirectionRef
  Severity: WARNING
  Description: Routes should not use DirectionRef (deprecated)
  XPath: //routes/Route/DirectionRef
  Example violation:
    <Route id="NO:Route:1" version="1">
      <DirectionRef ref="NO:Direction:Outbound"/>
    </Route>
  Fix: Use DirectionType instead: <DirectionType>outbound</DirectionType>

ROUTE_6: Route duplicated order
  Severity: WARNING
  Description: PointOnRoute elements must have unique order values
  XPath: //routes/Route/pointsInSequence/PointOnRoute[@order = preceding-sibling::PointOnRoute/@order]
  Example violation:
    <pointsInSequence>
      <PointOnRoute id="NO:PointOnRoute:1" order="1"/>
      <PointOnRoute id="NO:PointOnRoute:2" order="1"/> <!-- Duplicate order -->
    </pointsInSequence>
  Fix: Use sequential order values: order="1", order="2", etc.

ROUTE_7: Route missing DirectionType
  Severity: WARNING
  Description: Routes should specify direction (inbound/outbound/clockwise/anticlockwise)
  XPath: //routes/Route[not(DirectionType)]
  Example violation:
    <Route id="NO:Route:1" version="1">
      <Name>Downtown Route</Name>
      <!-- Missing DirectionType -->
    </Route>
  Fix: Add DirectionType: <DirectionType>outbound</DirectionType>

ROUTE_8: Route invalid DirectionType
  Severity: ERROR
  Description: DirectionType must be one of: inbound, outbound, clockwise, anticlockwise
  XPath: //routes/Route[DirectionType and not(DirectionType = 'inbound' or DirectionType = 'outbound' or DirectionType = 'clockwise' or DirectionType = 'anticlockwise')]
  Example violation:
    <Route id="NO:Route:1" version="1">
      <DirectionType>forward</DirectionType> <!-- Invalid value -->
    </Route>
  Fix: Use valid value: <DirectionType>outbound</DirectionType>
`

// SERVICE_JOURNEY Validation Rules
//
// ServiceJourney elements represent individual trips/runs of a public transport service.
// They define the specific times, stops, and operational details for each journey.

const ServiceJourneyValidationRulesDoc = `
SERVICE_JOURNEY VALIDATION RULES

SERVICE_JOURNEY_2: ServiceJourney illegal element Call
  Severity: ERROR
  Description: ServiceJourney should use passingTimes, not calls
  XPath: //vehicleJourneys/ServiceJourney/calls
  Example violation:
    <ServiceJourney id="NO:ServiceJourney:1" version="1">
      <calls><!-- Use passingTimes instead --></calls>
    </ServiceJourney>
  Fix: Use passingTimes: <passingTimes><TimetabledPassingTime>...</TimetabledPassingTime></passingTimes>

SERVICE_JOURNEY_3: ServiceJourney missing element PassingTimes
  Severity: ERROR
  Description: ServiceJourneys must have passingTimes for timetable information
  XPath: //vehicleJourneys/*[self::ServiceJourney][not(passingTimes)]
  Example violation:
    <ServiceJourney id="NO:ServiceJourney:1" version="1">
      <JourneyPatternRef ref="NO:JourneyPattern:1"/>
      <!-- Missing passingTimes -->
    </ServiceJourney>
  Fix: Add passingTimes with TimetabledPassingTime elements

SERVICE_JOURNEY_4: ServiceJourney missing arrival and departure
  Severity: ERROR
  Description: Each stop must have either arrival or departure time (or both)
  XPath: //vehicleJourneys/ServiceJourney/passingTimes/TimetabledPassingTime[not(DepartureTime or EarliestDepartureTime) and not(ArrivalTime or LatestArrivalTime)]
  Example violation:
    <TimetabledPassingTime id="NO:TimetabledPassingTime:1" version="1">
      <StopPointInJourneyPatternRef ref="NO:StopPoint:1"/>
      <!-- Missing both arrival and departure times -->
    </TimetabledPassingTime>
  Fix: Add time: <DepartureTime>09:30:00</DepartureTime>

SERVICE_JOURNEY_5: ServiceJourney missing departure times
  Severity: ERROR
  Description: First stop must have departure time
  XPath: //vehicleJourneys/ServiceJourney[not(passingTimes/TimetabledPassingTime[1]/DepartureTime) and not(passingTimes/TimetabledPassingTime[1]/EarliestDepartureTime)]
  Example violation:
    <passingTimes>
      <TimetabledPassingTime id="NO:TimetabledPassingTime:1" version="1">
        <ArrivalTime>09:00:00</ArrivalTime>
        <!-- Missing DepartureTime for first stop -->
      </TimetabledPassingTime>
    </passingTimes>
  Fix: Add DepartureTime: <DepartureTime>09:00:00</DepartureTime>

SERVICE_JOURNEY_6: ServiceJourney missing arrival time for last stop
  Severity: ERROR
  Description: Last stop must have arrival time
  XPath: //vehicleJourneys/ServiceJourney[count(passingTimes/TimetabledPassingTime[last()]/ArrivalTime) = 0 and count(passingTimes/TimetabledPassingTime[last()]/LatestArrivalTime) = 0]
  Example violation:
    <TimetabledPassingTime id="NO:TimetabledPassingTime:Last" version="1">
      <DepartureTime>10:30:00</DepartureTime>
      <!-- Missing ArrivalTime for last stop -->
    </TimetabledPassingTime>
  Fix: Add ArrivalTime: <ArrivalTime>10:30:00</ArrivalTime>

SERVICE_JOURNEY_10: ServiceJourney missing reference to JourneyPattern
  Severity: ERROR
  Description: ServiceJourneys must reference a JourneyPattern
  XPath: //vehicleJourneys/*[self::ServiceJourney][not(JourneyPatternRef)]
  Example violation:
    <ServiceJourney id="NO:ServiceJourney:1" version="1">
      <!-- Missing JourneyPatternRef -->
      <passingTimes>...</passingTimes>
    </ServiceJourney>
  Fix: Add JourneyPatternRef: <JourneyPatternRef ref="NO:JourneyPattern:1"/>

SERVICE_JOURNEY_12: ServiceJourney missing OperatorRef
  Severity: ERROR
  Description: ServiceJourneys must reference the operating company
  XPath: //vehicleJourneys/ServiceJourney[not(OperatorRef) and not(//ServiceFrame/lines/*[self::Line or self::FlexibleLine]/OperatorRef)]
  Example violation:
    <ServiceJourney id="NO:ServiceJourney:1" version="1">
      <!-- Missing OperatorRef -->
    </ServiceJourney>
  Fix: Add OperatorRef: <OperatorRef ref="NO:Operator:CityBus"/>
`

// TRANSPORT_MODE Validation Rules
//
// Transport mode validation ensures proper classification of services by transport type.
// Valid modes include bus, rail, tram, metro, coach, water, air, taxi, cableway, funicular.

const TransportModeValidationRulesDoc = `
TRANSPORT_MODE VALIDATION RULES

TRANSPORT_MODE_ON_LINE: Line with illegal TransportMode
  Severity: ERROR
  Description: Line TransportMode must be one of the predefined values
  Valid values: coach, bus, tram, rail, metro, air, taxi, water, cableway, funicular, unknown
  XPath: //lines/*[self::Line or self::FlexibleLine][TransportMode and not(TransportMode = 'coach' or TransportMode = 'bus' or ...)]
  Example violation:
    <Line id="NO:Line:1" version="1">
      <TransportMode>automobile</TransportMode> <!-- Invalid mode -->
    </Line>
  Fix: Use valid mode: <TransportMode>bus</TransportMode>

TRANSPORT_MODE_ON_SERVICE_JOURNEY: ServiceJourney with illegal TransportMode
  Severity: ERROR
  Description: ServiceJourney TransportMode must be one of the predefined values
  Valid values: coach, bus, tram, rail, metro, air, taxi, water, cableway, funicular, unknown
  XPath: //vehicleJourneys/ServiceJourney[TransportMode and not(TransportMode = 'coach' or TransportMode = 'bus' or ...)]
  Example violation:
    <ServiceJourney id="NO:ServiceJourney:1" version="1">
      <TransportMode>ship</TransportMode> <!-- Use 'water' instead -->
    </ServiceJourney>
  Fix: Use valid mode: <TransportMode>water</TransportMode>

TRANSPORT_SUB_MODE_BUS_INVALID: Line with invalid bus TransportSubMode
  Severity: ERROR
  Description: Bus lines must use valid bus submodes
  Valid bus submodes: localBus, regionalBus, expressBus, nightBus, postBus, specialNeedsBus, 
                     mobilityBus, mobilityBusForRegisteredDisabled, sightseeingBus, shuttleBus,
                     schoolBus, schoolAndPublicServiceBus, railReplacementBus, demandAndResponseBus, airportLinkBus
  Example violation:
    <Line id="NO:Line:1" version="1">
      <TransportMode>bus</TransportMode>
      <TransportSubmode>cityBus</TransportSubmode> <!-- Invalid submode -->
    </Line>
  Fix: Use valid submode: <TransportSubmode>localBus</TransportSubmode>

TRANSPORT_SUB_MODE_RAIL_INVALID: Line with invalid rail TransportSubMode
  Severity: ERROR
  Description: Rail lines must use valid rail submodes
  Valid rail submodes: local, highSpeedRail, suburbanRailway, regionalRail, interregionalRail,
                      longDistance, international, sleeperRailService, nightRail, 
                      carTransportRailService, touristRailway, railShuttle
  Example violation:
    <Line id="NO:Line:1" version="1">
      <TransportMode>rail</TransportMode>
      <TransportSubmode>fastTrain</TransportSubmode> <!-- Invalid submode -->
    </Line>
  Fix: Use valid submode: <TransportSubmode>highSpeedRail</TransportSubmode>

TRANSPORT_SUB_MODE_TRAM_INVALID: Line with invalid tram TransportSubMode
  Severity: ERROR
  Description: Tram lines must use valid tram submodes
  Valid tram submodes: cityTram, localTram, regionalTram, sightseeingTram, shuttleTram, trainTram
  Example violation:
    <Line id="NO:Line:1" version="1">
      <TransportMode>tram</TransportMode>
      <TransportSubmode>metroTram</TransportSubmode> <!-- Invalid submode -->
    </Line>
  Fix: Use valid submode: <TransportSubmode>cityTram</TransportSubmode>
`

// FLEXIBLE_LINE Validation Rules
//
// Flexible lines represent demand-responsive transport services with booking requirements.
// These rules ensure proper configuration of flexible service properties.

const FlexibleLineValidationRulesDoc = `
FLEXIBLE_LINE VALIDATION RULES

FLEXIBLE_LINE_1: FlexibleLine missing FlexibleLineType
  Severity: ERROR
  Description: FlexibleLines must specify the type of flexible service
  XPath: //lines/FlexibleLine[not(FlexibleLineType)]
  Valid types: fixedStop, flexibleAreasOnly, hailAndRideAreas, flexibleAreasAndStops,
              hailAndRideSections, fixedStopAreaWide, freeAreaAreaWide, mixedFlexible,
              mixedFlexibleAndFixed, fixed, mainRouteWithFlexibleEnds, flexibleRoute
  Example violation:
    <FlexibleLine id="NO:FlexibleLine:1" version="1">
      <Name>Demand Service</Name>
      <!-- Missing FlexibleLineType -->
    </FlexibleLine>
  Fix: Add FlexibleLineType: <FlexibleLineType>demandAndResponseServices</FlexibleLineType>

FLEXIBLE_LINE_10: FlexibleLine illegal use of both BookWhen and MinimumBookingPeriod
  Severity: WARNING
  Description: Use either BookWhen or MinimumBookingPeriod, not both
  XPath: //lines/FlexibleLine[BookWhen and MinimumBookingPeriod]
  Example violation:
    <FlexibleLine id="NO:FlexibleLine:1" version="1">
      <BookWhen>dayOfTravelOnly</BookWhen>
      <MinimumBookingPeriod>PT1H</MinimumBookingPeriod> <!-- Conflicting -->
    </FlexibleLine>
  Fix: Choose appropriate booking method based on service type

BOOKING_INVALID_BOOK_WHEN: Invalid BookWhen value
  Severity: ERROR
  Description: BookWhen must be one of the predefined values
  Valid values: dayOfTravelOnly, untilPreviousDay, advanceAndDayOfTravel
  XPath: //lines/FlexibleLine[BookWhen and not(BookWhen = 'dayOfTravelOnly' or BookWhen = 'untilPreviousDay' or BookWhen = 'advanceAndDayOfTravel')]
  Example violation:
    <FlexibleLine id="NO:FlexibleLine:1" version="1">
      <BookWhen>immediately</BookWhen> <!-- Invalid value -->
    </FlexibleLine>
  Fix: Use valid value: <BookWhen>dayOfTravelOnly</BookWhen>

BOOKING_MISSING_PROPERTIES: Mandatory booking property missing
  Severity: ERROR
  Description: Flexible lines must specify booking arrangements
  XPath: //lines/FlexibleLine[FlexibleLineType and (FlexibleLineType = 'flexibleAreasOnly' or FlexibleLineType = 'hailAndRideAreas' or FlexibleLineType = 'demandAndResponseServices') and not(BookWhen or MinimumBookingPeriod)]
  Example violation:
    <FlexibleLine id="NO:FlexibleLine:1" version="1">
      <FlexibleLineType>flexibleAreasOnly</FlexibleLineType>
      <!-- Missing booking properties -->
    </FlexibleLine>
  Fix: Add booking property: <BookWhen>dayOfTravelOnly</BookWhen>
`

// STOP_PLACE Validation Rules
//
// Stop places represent physical locations where passengers can board/alight vehicles.
// These rules ensure stop places have required identification and location information.

const StopPlaceValidationRulesDoc = `
STOP_PLACE VALIDATION RULES

STOP_PLACE_1: StopPlace missing Name
  Severity: ERROR
  Description: All stop places must have a descriptive name
  XPath: //stopPlaces/StopPlace[not(Name) or normalize-space(Name) = '']
  Example violation:
    <StopPlace id="NO:StopPlace:1" version="1">
      <!-- Missing Name element -->
      <Centroid>...</Centroid>
    </StopPlace>
  Fix: Add Name: <Name>Central Station</Name>

STOP_PLACE_2: StopPlace missing Centroid
  Severity: WARNING
  Description: Stop places should have location coordinates for mapping/navigation
  XPath: //stopPlaces/StopPlace[not(Centroid)]
  Example violation:
    <StopPlace id="NO:StopPlace:1" version="1">
      <Name>Central Station</Name>
      <!-- Missing Centroid with coordinates -->
    </StopPlace>
  Fix: Add Centroid: <Centroid><Location><Longitude>10.123</Longitude><Latitude>59.456</Latitude></Location></Centroid>

STOP_PLACE_3: Quay missing Name
  Severity: ERROR
  Description: All quays (platforms/boarding areas) must have names
  XPath: //stopPlaces/StopPlace/quays/Quay[not(Name) or normalize-space(Name) = '']
  Example violation:
    <StopPlace id="NO:StopPlace:1" version="1">
      <quays>
        <Quay id="NO:Quay:1" version="1">
          <!-- Missing Name element -->
        </Quay>
      </quays>
    </StopPlace>
  Fix: Add Name: <Name>Platform A</Name>

STOP_PLACE_4: Quay missing Centroid
  Severity: WARNING
  Description: Quays should have precise location coordinates
  XPath: //stopPlaces/StopPlace/quays/Quay[not(Centroid)]
  Example violation:
    <Quay id="NO:Quay:1" version="1">
      <Name>Platform A</Name>
      <!-- Missing Centroid -->
    </Quay>
  Fix: Add Centroid with coordinates

STOP_PLACE_5: Invalid StopPlaceType
  Severity: ERROR
  Description: StopPlaceType must be one of the predefined values
  Valid types: onstreetBus, onstreetTram, airport, railStation, metroStation, busStation,
              coachStation, tramStation, harbourPort, ferryPort, ferryStop, liftStation,
              vehicleRailInterchange, other
  Example violation:
    <StopPlace id="NO:StopPlace:1" version="1">
      <StopPlaceType>trainStation</StopPlaceType> <!-- Use 'railStation' -->
    </StopPlace>
  Fix: Use valid type: <StopPlaceType>railStation</StopPlaceType>
`

// JOURNEY_PATTERN Validation Rules
//
// Journey patterns define the sequence of stops that vehicles follow on their routes.
// They connect routes to specific stop sequences and timing points.

const JourneyPatternValidationRulesDoc = `
JOURNEY_PATTERN VALIDATION RULES

JOURNEY_PATTERN_1: JourneyPattern missing Name
  Severity: WARNING
  Description: Journey patterns should have descriptive names
  XPath: //journeyPatterns/*[self::JourneyPattern or self::ServiceJourneyPattern][not(Name) or normalize-space(Name) = '']
  Example violation:
    <JourneyPattern id="NO:JourneyPattern:1" version="1">
      <!-- Missing Name -->
      <RouteRef ref="NO:Route:1"/>
    </JourneyPattern>
  Fix: Add Name: <Name>City Center via Main Street</Name>

JOURNEY_PATTERN_2: JourneyPattern missing RouteRef
  Severity: ERROR
  Description: Journey patterns must reference the route they follow
  XPath: //journeyPatterns/*[self::JourneyPattern or self::ServiceJourneyPattern][not(RouteRef)]
  Example violation:
    <JourneyPattern id="NO:JourneyPattern:1" version="1">
      <Name>Downtown Pattern</Name>
      <!-- Missing RouteRef -->
    </JourneyPattern>
  Fix: Add RouteRef: <RouteRef ref="NO:Route:Downtown"/>

JOURNEY_PATTERN_3: JourneyPattern missing StopPoints
  Severity: ERROR
  Description: Journey patterns must define at least 2 stop points
  XPath: //journeyPatterns/*[self::JourneyPattern or self::ServiceJourneyPattern][not(pointsInSequence) or count(pointsInSequence/StopPointInJourneyPattern) < 2]
  Example violation:
    <JourneyPattern id="NO:JourneyPattern:1" version="1">
      <RouteRef ref="NO:Route:1"/>
      <pointsInSequence>
        <StopPointInJourneyPattern order="1" id="NO:StopPoint:1" version="1"/>
        <!-- Need at least 2 stops -->
      </pointsInSequence>
    </JourneyPattern>
  Fix: Add more stops or use different pattern type

JOURNEY_PATTERN_4: StopPointInJourneyPattern missing order
  Severity: ERROR
  Description: All stop points must have order attribute for sequencing
  XPath: //journeyPatterns/*[self::JourneyPattern or self::ServiceJourneyPattern]/pointsInSequence/StopPointInJourneyPattern[not(@order)]
  Example violation:
    <StopPointInJourneyPattern id="NO:StopPoint:1" version="1">
      <!-- Missing order attribute -->
      <ScheduledStopPointRef ref="NO:ScheduledStopPoint:1"/>
    </StopPointInJourneyPattern>
  Fix: Add order: <StopPointInJourneyPattern order="1" ...>

JOURNEY_PATTERN_5: Duplicate order in JourneyPattern
  Severity: ERROR
  Description: Order values must be unique within each journey pattern
  XPath: //journeyPatterns/*[self::JourneyPattern or self::ServiceJourneyPattern]/pointsInSequence/StopPointInJourneyPattern[@order = preceding-sibling::StopPointInJourneyPattern/@order or @order = following-sibling::StopPointInJourneyPattern/@order]
  Example violation:
    <pointsInSequence>
      <StopPointInJourneyPattern order="1" id="NO:StopPoint:1" version="1"/>
      <StopPointInJourneyPattern order="1" id="NO:StopPoint:2" version="1"/> <!-- Duplicate -->
    </pointsInSequence>
  Fix: Use sequential order: order="1", order="2", order="3", etc.
`

// NETWORK and OPERATOR Validation Rules
//
// Networks and operators represent the organizational structure of public transport services.
// Networks group lines together, while operators run the actual services.

const NetworkOperatorValidationRulesDoc = `
NETWORK AND OPERATOR VALIDATION RULES

NETWORK_1: Network missing Name
  Severity: ERROR
  Description: Transport networks must have descriptive names
  XPath: //networks/Network[not(Name) or normalize-space(Name) = '']
  Example violation:
    <Network id="NO:Network:1" version="1">
      <!-- Missing Name -->
      <AuthorityRef ref="NO:Authority:Oslo"/>
    </Network>
  Fix: Add Name: <Name>Oslo Public Transport Network</Name>

NETWORK_2: Network missing AuthorityRef
  Severity: WARNING
  Description: Networks should reference the responsible transport authority
  XPath: //networks/Network[not(AuthorityRef)]
  Example violation:
    <Network id="NO:Network:1" version="1">
      <Name>City Network</Name>
      <!-- Missing AuthorityRef -->
    </Network>
  Fix: Add AuthorityRef: <AuthorityRef ref="NO:Authority:Oslo"/>

OPERATOR_1: Operator missing Name
  Severity: ERROR
  Description: Transport operators must have company names
  XPath: //operators/Operator[not(Name) or normalize-space(Name) = '']
  Example violation:
    <Operator id="NO:Operator:1" version="1">
      <!-- Missing Name -->
      <CompanyNumber>123456789</CompanyNumber>
    </Operator>
  Fix: Add Name: <Name>Oslo Bus Company AS</Name>

OPERATOR_2: Operator missing ContactDetails
  Severity: WARNING
  Description: Operators should provide contact information for passengers
  XPath: //operators/Operator[not(ContactDetails)]
  Example violation:
    <Operator id="NO:Operator:1" version="1">
      <Name>City Bus</Name>
      <!-- Missing ContactDetails -->
    </Operator>
  Fix: Add ContactDetails: <ContactDetails><Phone>+47 123 45 678</Phone></ContactDetails>

AUTHORITY_1: Authority missing Name
  Severity: ERROR
  Description: Transport authorities must have official names
  XPath: //organisations/Authority[not(Name) or normalize-space(Name) = '']
  Example violation:
    <Authority id="NO:Authority:1" version="1">
      <!-- Missing Name -->
    </Authority>
  Fix: Add Name: <Name>Oslo County Municipality</Name>
`

// CALENDAR and VALIDITY Validation Rules
//
// Calendar rules ensure proper definition of operating periods and service dates.
// These rules validate day types, operating periods, and service calendars.

const CalendarValidityValidationRulesDoc = `
CALENDAR AND VALIDITY VALIDATION RULES

CALENDAR_1: DayType missing Name
  Severity: ERROR
  Description: Day types must have descriptive names (e.g., "Weekdays", "Saturdays")
  XPath: //dayTypes/DayType[not(Name) or normalize-space(Name) = '']
  Example violation:
    <DayType id="NO:DayType:1" version="1">
      <!-- Missing Name -->
      <properties>
        <PropertyOfDay>
          <DaysOfWeek>Monday Tuesday Wednesday Thursday Friday</DaysOfWeek>
        </PropertyOfDay>
      </properties>
    </DayType>
  Fix: Add Name: <Name>Weekdays</Name>

CALENDAR_2: OperatingDay missing CalendarDate
  Severity: ERROR
  Description: Operating days must specify the calendar date
  XPath: //operatingDays/OperatingDay[not(CalendarDate)]
  Example violation:
    <OperatingDay id="NO:OperatingDay:1" version="1">
      <Name>Christmas Day</Name>
      <!-- Missing CalendarDate -->
    </OperatingDay>
  Fix: Add CalendarDate: <CalendarDate>2023-12-25</CalendarDate>

CALENDAR_3: ServiceCalendar missing FromDate
  Severity: ERROR
  Description: Service calendars must specify start date of validity
  XPath: //serviceCalendar/ServiceCalendar[not(FromDate)]
  Example violation:
    <ServiceCalendar id="NO:ServiceCalendar:1" version="1">
      <!-- Missing FromDate -->
      <ToDate>2023-12-31</ToDate>
    </ServiceCalendar>
  Fix: Add FromDate: <FromDate>2023-01-01</FromDate>

CALENDAR_4: ServiceCalendar missing ToDate
  Severity: ERROR
  Description: Service calendars must specify end date of validity
  XPath: //serviceCalendar/ServiceCalendar[not(ToDate)]
  Example violation:
    <ServiceCalendar id="NO:ServiceCalendar:1" version="1">
      <FromDate>2023-01-01</FromDate>
      <!-- Missing ToDate -->
    </ServiceCalendar>
  Fix: Add ToDate: <ToDate>2023-12-31</ToDate>

CALENDAR_5: Invalid date range
  Severity: ERROR
  Description: FromDate must be before or equal to ToDate
  XPath: //serviceCalendar/ServiceCalendar[FromDate >= ToDate]
  Example violation:
    <ServiceCalendar id="NO:ServiceCalendar:1" version="1">
      <FromDate>2023-12-31</FromDate>
      <ToDate>2023-01-01</ToDate> <!-- ToDate before FromDate -->
    </ServiceCalendar>
  Fix: Correct date order: FromDate before ToDate

VALIDITY_CONDITIONS_1: ValidityConditions missing AvailabilityCondition
  Severity: WARNING
  Description: Validity conditions should specify availability periods
  XPath: //validityConditions[not(AvailabilityCondition)]
  Example violation:
    <validityConditions>
      <!-- Missing AvailabilityCondition -->
    </validityConditions>
  Fix: Add AvailabilityCondition with FromDate/ToDate
`

// FARE and PRICING Validation Rules
//
// Fare rules validate pricing information and fare zones for public transport services.
// These rules ensure proper fare structure definition and pricing information.

const FareValidationRulesDoc = `
FARE AND PRICING VALIDATION RULES

FARE_1: FareZone missing Name
  Severity: ERROR
  Description: Fare zones must have descriptive names for passenger information
  XPath: //fareZones/FareZone[not(Name) or normalize-space(Name) = '']
  Example violation:
    <FareZone id="NO:FareZone:1" version="1">
      <!-- Missing Name -->
      <members>...</members>
    </FareZone>
  Fix: Add Name: <Name>Zone 1 - City Center</Name>

FARE_2: TariffZone missing Name
  Severity: ERROR
  Description: Tariff zones must have descriptive names
  XPath: //tariffZones/TariffZone[not(Name) or normalize-space(Name) = '']
  Example violation:
    <TariffZone id="NO:TariffZone:1" version="1">
      <!-- Missing Name -->
    </TariffZone>
  Fix: Add Name: <Name>Inner City Tariff Zone</Name>

FARE_3: FareProduct missing Name
  Severity: ERROR
  Description: Fare products must have customer-facing names
  XPath: //fareProducts/FareProduct[not(Name) or normalize-space(Name) = '']
  Example violation:
    <FareProduct id="NO:FareProduct:1" version="1">
      <!-- Missing Name -->
      <prices>...</prices>
    </FareProduct>
  Fix: Add Name: <Name>Adult Single Journey</Name>

FARE_4: Missing price for FareProduct
  Severity: WARNING
  Description: Fare products should include pricing information
  XPath: //fareProducts/FareProduct[not(prices)]
  Example violation:
    <FareProduct id="NO:FareProduct:1" version="1">
      <Name>Adult Single</Name>
      <!-- Missing prices -->
    </FareProduct>
  Fix: Add prices: <prices><FarePrice><Amount>35.00</Amount><Currency>NOK</Currency></FarePrice></prices>
`

// VEHICLE and EQUIPMENT Validation Rules
//
// Vehicle rules validate vehicle types, equipment, and block assignments.
// These rules ensure proper vehicle specification and operational planning.

const VehicleValidationRulesDoc = `
VEHICLE AND EQUIPMENT VALIDATION RULES

VEHICLE_1: VehicleType missing Name
  Severity: ERROR
  Description: Vehicle types must have descriptive names
  XPath: //vehicleTypes/VehicleType[not(Name) or normalize-space(Name) = '']
  Example violation:
    <VehicleType id="NO:VehicleType:1" version="1">
      <!-- Missing Name -->
      <Length>12.0</Length>
    </VehicleType>
  Fix: Add Name: <Name>Standard City Bus</Name>

VEHICLE_2: Vehicle missing VehicleTypeRef
  Severity: ERROR
  Description: Vehicles must reference their vehicle type
  XPath: //vehicles/Vehicle[not(VehicleTypeRef)]
  Example violation:
    <Vehicle id="NO:Vehicle:1001" version="1">
      <Name>Bus 1001</Name>
      <!-- Missing VehicleTypeRef -->
    </Vehicle>
  Fix: Add VehicleTypeRef: <VehicleTypeRef ref="NO:VehicleType:CityBus"/>

VEHICLE_3: Block missing Name
  Severity: WARNING
  Description: Vehicle blocks should have identifying names
  XPath: //blocks/Block[not(Name) or normalize-space(Name) = '']
  Example violation:
    <Block id="NO:Block:1" version="1">
      <!-- Missing Name -->
      <journeys>...</journeys>
    </Block>
  Fix: Add Name: <Name>Morning Block A</Name>
`

// INTERCHANGE Validation Rules
//
// Interchange rules validate connections between different services and stop points.
// These rules ensure proper definition of transfer opportunities and timing.

const InterchangeValidationRulesDoc = `
INTERCHANGE AND CONNECTION VALIDATION RULES

INTERCHANGE_1: ServiceJourneyInterchange missing FromPointRef
  Severity: ERROR
  Description: Interchanges must specify the departure stop point
  XPath: //serviceJourneyInterchanges/ServiceJourneyInterchange[not(FromPointRef)]
  Example violation:
    <ServiceJourneyInterchange id="NO:Interchange:1" version="1">
      <!-- Missing FromPointRef -->
      <ToPointRef ref="NO:StopPoint:2"/>
    </ServiceJourneyInterchange>
  Fix: Add FromPointRef: <FromPointRef ref="NO:StopPoint:1"/>

INTERCHANGE_2: ServiceJourneyInterchange missing ToPointRef
  Severity: ERROR
  Description: Interchanges must specify the arrival stop point
  XPath: //serviceJourneyInterchanges/ServiceJourneyInterchange[not(ToPointRef)]
  Example violation:
    <ServiceJourneyInterchange id="NO:Interchange:1" version="1">
      <FromPointRef ref="NO:StopPoint:1"/>
      <!-- Missing ToPointRef -->
    </ServiceJourneyInterchange>
  Fix: Add ToPointRef: <ToPointRef ref="NO:StopPoint:2"/>

INTERCHANGE_3: ServiceJourneyInterchange missing FromJourneyRef
  Severity: ERROR
  Description: Interchanges must specify the departure service journey
  XPath: //serviceJourneyInterchanges/ServiceJourneyInterchange[not(FromJourneyRef)]
  Example violation:
    <ServiceJourneyInterchange id="NO:Interchange:1" version="1">
      <!-- Missing FromJourneyRef -->
      <ToJourneyRef ref="NO:ServiceJourney:2"/>
    </ServiceJourneyInterchange>
  Fix: Add FromJourneyRef: <FromJourneyRef ref="NO:ServiceJourney:1"/>

INTERCHANGE_4: ServiceJourneyInterchange missing ToJourneyRef
  Severity: ERROR
  Description: Interchanges must specify the arrival service journey
  XPath: //serviceJourneyInterchanges/ServiceJourneyInterchange[not(ToJourneyRef)]
  Example violation:
    <ServiceJourneyInterchange id="NO:Interchange:1" version="1">
      <FromJourneyRef ref="NO:ServiceJourney:1"/>
      <!-- Missing ToJourneyRef -->
    </ServiceJourneyInterchange>
  Fix: Add ToJourneyRef: <ToJourneyRef ref="NO:ServiceJourney:2"/>

INTERCHANGE_5: Invalid interchange duration
  Severity: ERROR
  Description: Transfer time must be positive duration
  XPath: //serviceJourneyInterchanges/ServiceJourneyInterchange/StandardTransferTime[number(.) <= 0]
  Example violation:
    <ServiceJourneyInterchange id="NO:Interchange:1" version="1">
      <StandardTransferTime>-PT5M</StandardTransferTime> <!-- Negative time -->
    </ServiceJourneyInterchange>
  Fix: Use positive duration: <StandardTransferTime>PT5M</StandardTransferTime>
`

// VERSION Validation Rules
//
// Version validation ensures proper version numbering across NetEX elements.
// These rules help maintain data consistency and change tracking.

const VersionValidationRulesDoc = `
VERSION VALIDATION RULES

VERSION_NON_NUMERIC: Non-numeric NeTEx version
  Severity: WARNING
  Description: Element versions should be numeric for proper versioning
  XPath: .//*[@version != 'any' and number(@version) != number(@version)]
  Example violation:
    <Line id="NO:Line:1" version="v1.0"> <!-- Non-numeric version -->
      <Name>City Line</Name>
    </Line>
  Fix: Use numeric version: version="1"

Special version values:
  'any' - Matches any version (used in references)
  Numeric values - Specific version numbers (1, 2, 3, etc.)
  Non-numeric values trigger warnings but don't invalidate the file
`

// FRAME Validation Rules
//
// Frame rules validate the overall structure of NetEX files and required frame types.
// These rules ensure proper file organization and completeness.

const FrameValidationRulesDoc = `
FRAME VALIDATION RULES

COMPOSITE_FRAME_1: CompositeFrame must be exactly one
  Severity: ERROR
  Description: NetEX files must contain exactly one CompositeFrame as root container
  XPath: //CompositeFrame[count(//CompositeFrame) != 1]
  Example violation:
    <dataObjects>
      <CompositeFrame id="NO:CompositeFrame:1" version="1">...</CompositeFrame>
      <CompositeFrame id="NO:CompositeFrame:2" version="1">...</CompositeFrame> <!-- Extra frame -->
    </dataObjects>
  Fix: Combine content into single CompositeFrame

RESOURCE_FRAME_IN_LINE_FILE: ResourceFrame must be exactly one in line file
  Severity: ERROR
  Description: Line files must contain exactly one ResourceFrame for operator/authority data
  XPath: //ResourceFrame[count(//ResourceFrame) != 1]
  Example violation:
    <CompositeFrame>
      <ResourceFrame>...</ResourceFrame>
      <ResourceFrame>...</ResourceFrame> <!-- Duplicate -->
    </CompositeFrame>
  Fix: Merge ResourceFrame content

SERVICE_FRAME_1: ServiceFrame missing in CompositeFrame
  Severity: ERROR
  Description: CompositeFrames must contain a ServiceFrame with line/route data
  XPath: //CompositeFrame[not(ServiceFrame)]
  Example violation:
    <CompositeFrame id="NO:CompositeFrame:1" version="1">
      <ResourceFrame>...</ResourceFrame>
      <!-- Missing ServiceFrame -->
    </CompositeFrame>
  Fix: Add ServiceFrame: <ServiceFrame>...</ServiceFrame>

TIMETABLE_FRAME_1: TimetableFrame missing in CompositeFrame
  Severity: WARNING
  Description: CompositeFrames should contain TimetableFrame for journey data
  XPath: //CompositeFrame[not(TimetableFrame)]
  Example violation:
    <CompositeFrame id="NO:CompositeFrame:1" version="1">
      <ServiceFrame>...</ServiceFrame>  
      <!-- Missing TimetableFrame -->
    </CompositeFrame>
  Fix: Add TimetableFrame: <TimetableFrame>...</TimetableFrame>
`

// GetValidationRulesDocumentation returns comprehensive documentation for all validation rules
func GetValidationRulesDocumentation() string {
	return LineValidationRulesDoc + "\n\n" +
		RouteValidationRulesDoc + "\n\n" +
		ServiceJourneyValidationRulesDoc + "\n\n" +
		TransportModeValidationRulesDoc + "\n\n" +
		FlexibleLineValidationRulesDoc + "\n\n" +
		StopPlaceValidationRulesDoc + "\n\n" +
		JourneyPatternValidationRulesDoc + "\n\n" +
		NetworkOperatorValidationRulesDoc + "\n\n" +
		CalendarValidityValidationRulesDoc + "\n\n" +
		FareValidationRulesDoc + "\n\n" +
		VehicleValidationRulesDoc + "\n\n" +
		InterchangeValidationRulesDoc + "\n\n" +
		VersionValidationRulesDoc + "\n\n" +
		FrameValidationRulesDoc
}

// GetRulesBySeverity returns rules grouped by severity level
func GetRulesBySeverity() map[string][]string {
	return map[string][]string{
		"ERROR": {
			"LINE_2", "LINE_4", "LINE_6", "LINE_7", "LINE_8",
			"ROUTE_2", "ROUTE_3", "ROUTE_4", "ROUTE_8", 
			"SERVICE_JOURNEY_2", "SERVICE_JOURNEY_3", "SERVICE_JOURNEY_4", "SERVICE_JOURNEY_5", "SERVICE_JOURNEY_6", "SERVICE_JOURNEY_10", "SERVICE_JOURNEY_12",
			"TRANSPORT_MODE_ON_LINE", "TRANSPORT_MODE_ON_SERVICE_JOURNEY", "TRANSPORT_SUB_MODE_BUS_INVALID", "TRANSPORT_SUB_MODE_RAIL_INVALID", "TRANSPORT_SUB_MODE_TRAM_INVALID",
			"FLEXIBLE_LINE_1", "BOOKING_INVALID_BOOK_WHEN", "BOOKING_MISSING_PROPERTIES",
			"STOP_PLACE_1", "STOP_PLACE_3", "STOP_PLACE_5",
			"JOURNEY_PATTERN_2", "JOURNEY_PATTERN_3", "JOURNEY_PATTERN_4", "JOURNEY_PATTERN_5",
			"NETWORK_1", "OPERATOR_1", "AUTHORITY_1",
			"CALENDAR_2", "CALENDAR_3", "CALENDAR_4", "CALENDAR_5",
			"FARE_1", "FARE_2", "FARE_3",
			"VEHICLE_1", "VEHICLE_2",
			"INTERCHANGE_1", "INTERCHANGE_2", "INTERCHANGE_3", "INTERCHANGE_4", "INTERCHANGE_5",
			"COMPOSITE_FRAME_1", "RESOURCE_FRAME_IN_LINE_FILE", "SERVICE_FRAME_1",
		},
		"WARNING": {
			"LINE_3", "LINE_5", "LINE_9",
			"ROUTE_5", "ROUTE_6", "ROUTE_7",
			"FLEXIBLE_LINE_10",
			"STOP_PLACE_2", "STOP_PLACE_4",
			"JOURNEY_PATTERN_1",
			"NETWORK_2", "OPERATOR_2",
			"VALIDITY_CONDITIONS_1",
			"FARE_4",
			"VEHICLE_3",
			"VERSION_NON_NUMERIC",
			"TIMETABLE_FRAME_1",
		},
	}
}

// GetRulesByCategory returns rules grouped by functional category
func GetRulesByCategory() map[string][]string {
	return map[string][]string{
		"Line and Route": {
			"LINE_2", "LINE_3", "LINE_4", "LINE_5", "LINE_6", "LINE_7", "LINE_8", "LINE_9",
			"ROUTE_2", "ROUTE_3", "ROUTE_4", "ROUTE_5", "ROUTE_6", "ROUTE_7", "ROUTE_8",
		},
		"Service Journey and Timetables": {
			"SERVICE_JOURNEY_2", "SERVICE_JOURNEY_3", "SERVICE_JOURNEY_4", "SERVICE_JOURNEY_5", 
			"SERVICE_JOURNEY_6", "SERVICE_JOURNEY_10", "SERVICE_JOURNEY_12",
		},
		"Transport Modes": {
			"TRANSPORT_MODE_ON_LINE", "TRANSPORT_MODE_ON_SERVICE_JOURNEY", 
			"TRANSPORT_SUB_MODE_BUS_INVALID", "TRANSPORT_SUB_MODE_RAIL_INVALID", "TRANSPORT_SUB_MODE_TRAM_INVALID",
		},
		"Flexible Services": {
			"FLEXIBLE_LINE_1", "FLEXIBLE_LINE_10", "BOOKING_INVALID_BOOK_WHEN", "BOOKING_MISSING_PROPERTIES",
		},
		"Stop Places and Infrastructure": {
			"STOP_PLACE_1", "STOP_PLACE_2", "STOP_PLACE_3", "STOP_PLACE_4", "STOP_PLACE_5",
		},
		"Journey Patterns": {
			"JOURNEY_PATTERN_1", "JOURNEY_PATTERN_2", "JOURNEY_PATTERN_3", "JOURNEY_PATTERN_4", "JOURNEY_PATTERN_5",
		},
		"Network and Organization": {
			"NETWORK_1", "NETWORK_2", "OPERATOR_1", "OPERATOR_2", "AUTHORITY_1",
		},
		"Calendar and Validity": {
			"CALENDAR_2", "CALENDAR_3", "CALENDAR_4", "CALENDAR_5", "VALIDITY_CONDITIONS_1",
		},
		"Fares and Pricing": {
			"FARE_1", "FARE_2", "FARE_3", "FARE_4",
		},
		"Vehicles and Equipment": {
			"VEHICLE_1", "VEHICLE_2", "VEHICLE_3",
		},
		"Interchanges": {
			"INTERCHANGE_1", "INTERCHANGE_2", "INTERCHANGE_3", "INTERCHANGE_4", "INTERCHANGE_5",
		},
		"File Structure": {
			"COMPOSITE_FRAME_1", "RESOURCE_FRAME_IN_LINE_FILE", "SERVICE_FRAME_1", "TIMETABLE_FRAME_1",
		},
		"Versioning": {
			"VERSION_NON_NUMERIC",
		},
	}
}