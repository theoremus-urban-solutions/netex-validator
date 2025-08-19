package context

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/antchfx/xmlquery"
)

// ObjectValidationContext provides rich context for object model validation
type ObjectValidationContext struct {
	// Basic validation information
	FileName           string
	Codespace          string
	ValidationReportID string
	IsCommonFile       bool

	// Parsed NetEX data
	PublicationDelivery *PublicationDelivery
	Document            *xmlquery.Node

	// Element collections for fast lookup
	elementIndex   map[string]NetexObject
	referenceIndex map[string][]string // ref -> list of referencing element IDs

	// Frame-specific collections
	operators            map[string]*Operator
	authorities          map[string]*Authority
	networks             map[string]*Network
	lines                map[string]*Line
	flexibleLines        map[string]*FlexibleLine
	routes               map[string]*Route
	journeyPatterns      map[string]*JourneyPattern
	serviceJourneys      map[string]*ServiceJourney
	datedServiceJourneys map[string]*DatedServiceJourney
	scheduledStopPoints  map[string]*ScheduledStopPoint
	stopPlaces           map[string]*StopPlace
	quays                map[string]*Quay
	dayTypes             map[string]*DayType
	operatingDays        map[string]*OperatingDay
	blocks               map[string]*Block

	// Common data collections (shared across files)
	commonDataRepository *CommonDataRepository
}

// CommonDataRepository holds shared data across multiple files
type CommonDataRepository struct {
	sharedOperators   map[string]*Operator
	sharedAuthorities map[string]*Authority
	sharedStopPlaces  map[string]*StopPlace
	sharedQuays       map[string]*Quay
	sharedNetworks    map[string]*Network
}

// NewObjectValidationContext creates a new object validation context
func NewObjectValidationContext(fileName, codespace, reportID string, xmlData []byte, doc *xmlquery.Node) (*ObjectValidationContext, error) {
	ctx := &ObjectValidationContext{
		FileName:             fileName,
		Codespace:            codespace,
		ValidationReportID:   reportID,
		IsCommonFile:         strings.HasPrefix(fileName, "_"),
		Document:             doc,
		elementIndex:         make(map[string]NetexObject),
		referenceIndex:       make(map[string][]string),
		operators:            make(map[string]*Operator),
		authorities:          make(map[string]*Authority),
		networks:             make(map[string]*Network),
		lines:                make(map[string]*Line),
		flexibleLines:        make(map[string]*FlexibleLine),
		routes:               make(map[string]*Route),
		journeyPatterns:      make(map[string]*JourneyPattern),
		serviceJourneys:      make(map[string]*ServiceJourney),
		datedServiceJourneys: make(map[string]*DatedServiceJourney),
		scheduledStopPoints:  make(map[string]*ScheduledStopPoint),
		stopPlaces:           make(map[string]*StopPlace),
		quays:                make(map[string]*Quay),
		dayTypes:             make(map[string]*DayType),
		operatingDays:        make(map[string]*OperatingDay),
		blocks:               make(map[string]*Block),
	}

	// Parse XML into object model
	if err := ctx.parseNetexData(xmlData); err != nil {
		return nil, fmt.Errorf("failed to parse NetEX data: %w", err)
	}

	// Build indices for fast lookup
	ctx.buildIndices()

	return ctx, nil
}

// parseNetexData parses XML data into the object model
func (ctx *ObjectValidationContext) parseNetexData(xmlData []byte) error {
	var delivery PublicationDelivery
	if err := xml.Unmarshal(xmlData, &delivery); err != nil {
		return fmt.Errorf("failed to unmarshal XML: %w", err)
	}

	ctx.PublicationDelivery = &delivery
	return nil
}

// buildIndices builds lookup indices for fast element access
func (ctx *ObjectValidationContext) buildIndices() {
	if ctx.PublicationDelivery == nil || ctx.PublicationDelivery.DataObjects == nil {
		return
	}

	dataObjects := ctx.PublicationDelivery.DataObjects

	// Check for frames inside CompositeFrame first
	if dataObjects.CompositeFrame != nil && dataObjects.CompositeFrame.Frames != nil {
		frames := dataObjects.CompositeFrame.Frames

		// Index frames from CompositeFrame
		if frames.ResourceFrame != nil {
			ctx.indexResourceFrame(frames.ResourceFrame)
		}
		if frames.ServiceFrame != nil {
			ctx.indexServiceFrame(frames.ServiceFrame)
		}
		if frames.TimetableFrame != nil {
			ctx.indexTimetableFrame(frames.TimetableFrame)
		}
		if frames.SiteFrame != nil {
			ctx.indexSiteFrame(frames.SiteFrame)
		}
		if frames.ServiceCalendarFrame != nil {
			ctx.indexServiceCalendarFrame(frames.ServiceCalendarFrame)
		}
		if frames.VehicleScheduleFrame != nil {
			ctx.indexVehicleScheduleFrame(frames.VehicleScheduleFrame)
		}
	}

	// Check for direct frames in DataObjects (common in simple cases)
	if dataObjects.ResourceFrame != nil {
		ctx.indexResourceFrame(dataObjects.ResourceFrame)
	}
	if dataObjects.ServiceFrame != nil {
		ctx.indexServiceFrame(dataObjects.ServiceFrame)
	}
	if dataObjects.TimetableFrame != nil {
		ctx.indexTimetableFrame(dataObjects.TimetableFrame)
	}
	if dataObjects.SiteFrame != nil {
		ctx.indexSiteFrame(dataObjects.SiteFrame)
	}
	if dataObjects.ServiceCalendarFrame != nil {
		ctx.indexServiceCalendarFrame(dataObjects.ServiceCalendarFrame)
	}
	if dataObjects.VehicleScheduleFrame != nil {
		ctx.indexVehicleScheduleFrame(dataObjects.VehicleScheduleFrame)
	}
}

// indexResourceFrame indexes elements from ResourceFrame
func (ctx *ObjectValidationContext) indexResourceFrame(frame *ResourceFrame) {
	if frame.Organisations != nil {
		for _, operator := range frame.Organisations.Operators {
			if operator.ID != "" {
				ctx.operators[operator.ID] = operator
				ctx.elementIndex[operator.ID] = operator
			}
		}
		for _, authority := range frame.Organisations.Authorities {
			if authority.ID != "" {
				ctx.authorities[authority.ID] = authority
				ctx.elementIndex[authority.ID] = authority
			}
		}
	}
}

// indexServiceFrame indexes elements from ServiceFrame
func (ctx *ObjectValidationContext) indexServiceFrame(frame *ServiceFrame) {
	// Index networks
	if frame.Networks != nil {
		for _, network := range frame.Networks.Networks {
			if network.ID != "" {
				ctx.networks[network.ID] = network
				ctx.elementIndex[network.ID] = network
			}
		}
	}

	// Index lines
	if frame.Lines != nil {
		for _, line := range frame.Lines.Lines {
			if line.ID != "" {
				ctx.lines[line.ID] = line
				ctx.elementIndex[line.ID] = line
			}
		}
		for _, flexLine := range frame.Lines.FlexibleLines {
			if flexLine.ID != "" {
				ctx.flexibleLines[flexLine.ID] = flexLine
				ctx.elementIndex[flexLine.ID] = flexLine
			}
		}
	}

	// Index routes
	if frame.Routes != nil {
		for _, route := range frame.Routes.Routes {
			if route.ID != "" {
				ctx.routes[route.ID] = route
				ctx.elementIndex[route.ID] = route
			}
		}
	}

	// Index journey patterns
	if frame.JourneyPatterns != nil {
		for _, jp := range frame.JourneyPatterns.JourneyPatterns {
			if jp.ID != "" {
				ctx.journeyPatterns[jp.ID] = jp
				ctx.elementIndex[jp.ID] = jp
			}
		}
	}

	// Index vehicle journeys
	if frame.VehicleJourneys != nil {
		for _, sj := range frame.VehicleJourneys.ServiceJourneys {
			if sj.ID != "" {
				ctx.serviceJourneys[sj.ID] = sj
				ctx.elementIndex[sj.ID] = sj
			}
		}
	}

	// Index scheduled stop points
	if frame.ScheduledStopPoints != nil {
		for _, ssp := range frame.ScheduledStopPoints.ScheduledStopPoints {
			if ssp.ID != "" {
				ctx.scheduledStopPoints[ssp.ID] = ssp
				ctx.elementIndex[ssp.ID] = ssp
			}
		}
	}
}

// indexTimetableFrame indexes elements from TimetableFrame
func (ctx *ObjectValidationContext) indexTimetableFrame(frame *TimetableFrame) {
	if frame.VehicleJourneys != nil {
		for _, dsj := range frame.VehicleJourneys.DatedServiceJourneys {
			if dsj.ID != "" {
				ctx.datedServiceJourneys[dsj.ID] = dsj
				ctx.elementIndex[dsj.ID] = dsj
			}
		}
	}
}

// indexSiteFrame indexes elements from SiteFrame
func (ctx *ObjectValidationContext) indexSiteFrame(frame *SiteFrame) {
	if frame.StopPlaces != nil {
		for _, sp := range frame.StopPlaces.StopPlaces {
			if sp.ID != "" {
				ctx.stopPlaces[sp.ID] = sp
				ctx.elementIndex[sp.ID] = sp
			}

			// Index quays
			if sp.Quays != nil {
				for _, quay := range sp.Quays.Quays {
					if quay.ID != "" {
						ctx.quays[quay.ID] = quay
						ctx.elementIndex[quay.ID] = quay
					}
				}
			}
		}
	}
}

// indexServiceCalendarFrame indexes elements from ServiceCalendarFrame
func (ctx *ObjectValidationContext) indexServiceCalendarFrame(frame *ServiceCalendarFrame) {
	if frame.DayTypes != nil {
		for _, dt := range frame.DayTypes.DayTypes {
			if dt.ID != "" {
				ctx.dayTypes[dt.ID] = dt
				ctx.elementIndex[dt.ID] = dt
			}
		}
	}

	if frame.OperatingDays != nil {
		for _, od := range frame.OperatingDays.OperatingDays {
			if od.ID != "" {
				ctx.operatingDays[od.ID] = od
				ctx.elementIndex[od.ID] = od
			}
		}
	}
}

// indexVehicleScheduleFrame indexes elements from VehicleScheduleFrame
func (ctx *ObjectValidationContext) indexVehicleScheduleFrame(frame *VehicleScheduleFrame) {
	if frame.Blocks != nil {
		for _, block := range frame.Blocks.Blocks {
			if block.ID != "" {
				ctx.blocks[block.ID] = block
				ctx.elementIndex[block.ID] = block
			}
		}
	}
}

// Element access methods

// GetElementByID returns any NetEX element by ID
func (ctx *ObjectValidationContext) GetElementByID(id string) NetexObject {
	return ctx.elementIndex[id]
}

// GetOperator returns an operator by ID
func (ctx *ObjectValidationContext) GetOperator(id string) *Operator {
	return ctx.operators[id]
}

// GetAuthority returns an authority by ID
func (ctx *ObjectValidationContext) GetAuthority(id string) *Authority {
	return ctx.authorities[id]
}

// GetLine returns a line by ID
func (ctx *ObjectValidationContext) GetLine(id string) *Line {
	return ctx.lines[id]
}

// GetFlexibleLine returns a flexible line by ID
func (ctx *ObjectValidationContext) GetFlexibleLine(id string) *FlexibleLine {
	return ctx.flexibleLines[id]
}

// GetRoute returns a route by ID
func (ctx *ObjectValidationContext) GetRoute(id string) *Route {
	return ctx.routes[id]
}

// GetJourneyPattern returns a journey pattern by ID
func (ctx *ObjectValidationContext) GetJourneyPattern(id string) *JourneyPattern {
	return ctx.journeyPatterns[id]
}

// GetServiceJourney returns a service journey by ID
func (ctx *ObjectValidationContext) GetServiceJourney(id string) *ServiceJourney {
	return ctx.serviceJourneys[id]
}

// GetDatedServiceJourney returns a dated service journey by ID
func (ctx *ObjectValidationContext) GetDatedServiceJourney(id string) *DatedServiceJourney {
	return ctx.datedServiceJourneys[id]
}

// GetScheduledStopPoint returns a scheduled stop point by ID
func (ctx *ObjectValidationContext) GetScheduledStopPoint(id string) *ScheduledStopPoint {
	return ctx.scheduledStopPoints[id]
}

// GetStopPlace returns a stop place by ID
func (ctx *ObjectValidationContext) GetStopPlace(id string) *StopPlace {
	return ctx.stopPlaces[id]
}

// GetQuay returns a quay by ID
func (ctx *ObjectValidationContext) GetQuay(id string) *Quay {
	return ctx.quays[id]
}

// GetDayType returns a day type by ID
func (ctx *ObjectValidationContext) GetDayType(id string) *DayType {
	return ctx.dayTypes[id]
}

// GetOperatingDay returns an operating day by ID
func (ctx *ObjectValidationContext) GetOperatingDay(id string) *OperatingDay {
	return ctx.operatingDays[id]
}

// GetBlock returns a block by ID
func (ctx *ObjectValidationContext) GetBlock(id string) *Block {
	return ctx.blocks[id]
}

// Collection access methods

// ServiceJourneys returns all service journeys
func (ctx *ObjectValidationContext) ServiceJourneys() []*ServiceJourney {
	var journeys []*ServiceJourney
	for _, sj := range ctx.serviceJourneys {
		journeys = append(journeys, sj)
	}
	return journeys
}

// Lines returns all lines
func (ctx *ObjectValidationContext) Lines() []*Line {
	var lines []*Line
	for _, line := range ctx.lines {
		lines = append(lines, line)
	}
	return lines
}

// Operators returns all operators
func (ctx *ObjectValidationContext) Operators() []*Operator {
	var operators []*Operator
	for _, op := range ctx.operators {
		operators = append(operators, op)
	}
	return operators
}

// FlexibleLines returns all flexible lines
func (ctx *ObjectValidationContext) FlexibleLines() []*FlexibleLine {
	var lines []*FlexibleLine
	for _, line := range ctx.flexibleLines {
		lines = append(lines, line)
	}
	return lines
}

// Routes returns all routes
func (ctx *ObjectValidationContext) Routes() []*Route {
	var routes []*Route
	for _, route := range ctx.routes {
		routes = append(routes, route)
	}
	return routes
}

// JourneyPatterns returns all journey patterns
func (ctx *ObjectValidationContext) JourneyPatterns() []*JourneyPattern {
	var patterns []*JourneyPattern
	for _, jp := range ctx.journeyPatterns {
		patterns = append(patterns, jp)
	}
	return patterns
}

// StopPlaces returns all stop places
func (ctx *ObjectValidationContext) StopPlaces() []*StopPlace {
	var places []*StopPlace
	for _, sp := range ctx.stopPlaces {
		places = append(places, sp)
	}
	return places
}

// DatedServiceJourneys returns all dated service journeys
func (ctx *ObjectValidationContext) DatedServiceJourneys() []*DatedServiceJourney {
	var journeys []*DatedServiceJourney
	for _, dsj := range ctx.datedServiceJourneys {
		journeys = append(journeys, dsj)
	}
	return journeys
}

// Utility methods

// DataLocation creates a data location for an element
func (ctx *ObjectValidationContext) DataLocation(elementID string) *DataLocation {
	return &DataLocation{
		FileName:  ctx.FileName,
		ElementID: elementID,
		XPath:     fmt.Sprintf("//*[@id='%s']", elementID),
	}
}

// IsReferenceResolved checks if a reference can be resolved locally
func (ctx *ObjectValidationContext) IsReferenceResolved(ref string) bool {
	return ctx.elementIndex[ref] != nil
}

// GetReferencedElement resolves a reference to its target element
func (ctx *ObjectValidationContext) GetReferencedElement(ref string) NetexObject {
	return ctx.elementIndex[ref]
}

// HasFrame checks if a specific frame type exists
func (ctx *ObjectValidationContext) HasFrame(frameType string) bool {
	if ctx.PublicationDelivery == nil || ctx.PublicationDelivery.DataObjects == nil {
		return false
	}

	dataObjects := ctx.PublicationDelivery.DataObjects

	switch frameType {
	case "CompositeFrame":
		return dataObjects.CompositeFrame != nil
	case "ResourceFrame":
		// Check direct frame first, then in CompositeFrame
		if dataObjects.ResourceFrame != nil {
			return true
		}
		return dataObjects.CompositeFrame != nil &&
			dataObjects.CompositeFrame.Frames != nil &&
			dataObjects.CompositeFrame.Frames.ResourceFrame != nil
	case "ServiceFrame":
		if dataObjects.ServiceFrame != nil {
			return true
		}
		return dataObjects.CompositeFrame != nil &&
			dataObjects.CompositeFrame.Frames != nil &&
			dataObjects.CompositeFrame.Frames.ServiceFrame != nil
	case "TimetableFrame":
		if dataObjects.TimetableFrame != nil {
			return true
		}
		return dataObjects.CompositeFrame != nil &&
			dataObjects.CompositeFrame.Frames != nil &&
			dataObjects.CompositeFrame.Frames.TimetableFrame != nil
	case "SiteFrame":
		if dataObjects.SiteFrame != nil {
			return true
		}
		return dataObjects.CompositeFrame != nil &&
			dataObjects.CompositeFrame.Frames != nil &&
			dataObjects.CompositeFrame.Frames.SiteFrame != nil
	case "ServiceCalendarFrame":
		if dataObjects.ServiceCalendarFrame != nil {
			return true
		}
		return dataObjects.CompositeFrame != nil &&
			dataObjects.CompositeFrame.Frames != nil &&
			dataObjects.CompositeFrame.Frames.ServiceCalendarFrame != nil
	case "VehicleScheduleFrame":
		if dataObjects.VehicleScheduleFrame != nil {
			return true
		}
		return dataObjects.CompositeFrame != nil &&
			dataObjects.CompositeFrame.Frames != nil &&
			dataObjects.CompositeFrame.Frames.VehicleScheduleFrame != nil
	default:
		return false
	}
}

// SetCommonDataRepository sets the shared data repository
func (ctx *ObjectValidationContext) SetCommonDataRepository(repo *CommonDataRepository) {
	ctx.commonDataRepository = repo
}

// GetCommonDataRepository returns the shared data repository
func (ctx *ObjectValidationContext) GetCommonDataRepository() *CommonDataRepository {
	return ctx.commonDataRepository
}

// NewCommonDataRepository creates a new common data repository
func NewCommonDataRepository() *CommonDataRepository {
	return &CommonDataRepository{
		sharedOperators:   make(map[string]*Operator),
		sharedAuthorities: make(map[string]*Authority),
		sharedStopPlaces:  make(map[string]*StopPlace),
		sharedQuays:       make(map[string]*Quay),
		sharedNetworks:    make(map[string]*Network),
	}
}

// AddSharedOperator adds an operator to the common data repository
func (repo *CommonDataRepository) AddSharedOperator(operator *Operator) {
	if operator.ID != "" {
		repo.sharedOperators[operator.ID] = operator
	}
}

// GetSharedOperator gets an operator from the common data repository
func (repo *CommonDataRepository) GetSharedOperator(id string) *Operator {
	return repo.sharedOperators[id]
}

// AddSharedStopPlace adds a stop place to the common data repository
func (repo *CommonDataRepository) AddSharedStopPlace(stopPlace *StopPlace) {
	if stopPlace.ID != "" {
		repo.sharedStopPlaces[stopPlace.ID] = stopPlace
	}
}

// GetSharedStopPlace gets a stop place from the common data repository
func (repo *CommonDataRepository) GetSharedStopPlace(id string) *StopPlace {
	return repo.sharedStopPlaces[id]
}
