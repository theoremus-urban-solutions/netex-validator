package model

import (
	"encoding/xml"
	"time"
)

// NetexObject represents the base interface for all NetEX elements
type NetexObject interface {
	GetID() string
	GetVersion() string
	GetDataLocation() *DataLocation
}

// DataLocation provides location information for validation errors
type DataLocation struct {
	FileName    string `json:"fileName"`
	LineNumber  int    `json:"lineNumber"`
	XPath       string `json:"xpath"`
	ElementID   string `json:"elementId"`
}

// BaseNetexObject provides common fields for all NetEX objects
type BaseNetexObject struct {
	ID           string        `xml:"id,attr"`
	Version      string        `xml:"version,attr"`
	DataLocation *DataLocation `xml:"-"`
}

func (b *BaseNetexObject) GetID() string {
	return b.ID
}

func (b *BaseNetexObject) GetVersion() string {
	return b.Version
}

func (b *BaseNetexObject) GetDataLocation() *DataLocation {
	return b.DataLocation
}

// PublicationDelivery represents the root NetEX element
type PublicationDelivery struct {
	XMLName              xml.Name               `xml:"PublicationDelivery"`
	PublicationTimestamp time.Time              `xml:"PublicationTimestamp"`
	ParticipantRef       string                 `xml:"ParticipantRef"`
	DataObjects          *DataObjects           `xml:"dataObjects"`
}

// DataObjects contains the frames
type DataObjects struct {
	CompositeFrame *CompositeFrame `xml:"CompositeFrame"`
}

// CompositeFrame represents a NetEX composite frame
type CompositeFrame struct {
	BaseNetexObject
	XMLName     xml.Name      `xml:"CompositeFrame"`
	ValidityConditions *ValidityConditions `xml:"validityConditions"`
	Frames      *Frames       `xml:"frames"`
}

// Frames contains different types of frames
type Frames struct {
	ResourceFrame           *ResourceFrame           `xml:"ResourceFrame"`
	ServiceFrame            *ServiceFrame            `xml:"ServiceFrame"`
	TimetableFrame          *TimetableFrame          `xml:"TimetableFrame"`
	SiteFrame              *SiteFrame               `xml:"SiteFrame"`
	ServiceCalendarFrame   *ServiceCalendarFrame    `xml:"ServiceCalendarFrame"`
	VehicleScheduleFrame   *VehicleScheduleFrame    `xml:"VehicleScheduleFrame"`
}

// ResourceFrame contains organizational data
type ResourceFrame struct {
	BaseNetexObject
	XMLName       xml.Name       `xml:"ResourceFrame"`
	Organisations *Organisations `xml:"organisations"`
	VehicleTypes  *VehicleTypes  `xml:"vehicleTypes"`
}

// ServiceFrame contains service-related data
type ServiceFrame struct {
	BaseNetexObject
	XMLName                xml.Name                `xml:"ServiceFrame"`
	Networks               *Networks               `xml:"networks"`
	Lines                  *Lines                  `xml:"lines"`
	Routes                 *Routes                 `xml:"routes"`
	JourneyPatterns        *JourneyPatterns        `xml:"journeyPatterns"`
	VehicleJourneys        *VehicleJourneys        `xml:"vehicleJourneys"`
	ScheduledStopPoints    *ScheduledStopPoints    `xml:"scheduledStopPoints"`
	StopAssignments        *StopAssignments        `xml:"stopAssignments"`
	Interchanges           *Interchanges           `xml:"interchanges"`
}

// TimetableFrame contains timetable-related data
type TimetableFrame struct {
	BaseNetexObject
	XMLName         xml.Name         `xml:"TimetableFrame"`
	VehicleJourneys *VehicleJourneys `xml:"vehicleJourneys"`
}

// SiteFrame contains stop place data
type SiteFrame struct {
	BaseNetexObject
	XMLName    xml.Name    `xml:"SiteFrame"`
	StopPlaces *StopPlaces `xml:"stopPlaces"`
}

// ServiceCalendarFrame contains calendar data
type ServiceCalendarFrame struct {
	BaseNetexObject
	XMLName        xml.Name        `xml:"ServiceCalendarFrame"`
	ServiceCalendar *ServiceCalendar `xml:"serviceCalendar"`
	DayTypes       *DayTypes       `xml:"dayTypes"`
	OperatingDays  *OperatingDays  `xml:"operatingDays"`
}

// VehicleScheduleFrame contains vehicle schedule data
type VehicleScheduleFrame struct {
	BaseNetexObject
	XMLName xml.Name `xml:"VehicleScheduleFrame"`
	Blocks  *Blocks  `xml:"blocks"`
}

// Organisations contains operators and authorities
type Organisations struct {
	Operators   []*Operator   `xml:"Operator"`
	Authorities []*Authority  `xml:"Authority"`
}

// Operator represents a transport operator
type Operator struct {
	BaseNetexObject
	XMLName        xml.Name        `xml:"Operator"`
	Name           string          `xml:"Name"`
	ShortName      string          `xml:"ShortName"`
	ContactDetails *ContactDetails `xml:"ContactDetails"`
}

// Authority represents a transport authority
type Authority struct {
	BaseNetexObject
	XMLName        xml.Name        `xml:"Authority"`
	Name           string          `xml:"Name"`
	ShortName      string          `xml:"ShortName"`
	ContactDetails *ContactDetails `xml:"ContactDetails"`
}

// ContactDetails represents contact information
type ContactDetails struct {
	Phone string `xml:"Phone"`
	Email string `xml:"Email"`
	Url   string `xml:"Url"`
}

// Networks contains network information
type Networks struct {
	Networks []*Network `xml:"Network"`
}

// Network represents a transport network
type Network struct {
	BaseNetexObject
	XMLName      xml.Name      `xml:"Network"`
	Name         string        `xml:"Name"`
	AuthorityRef *AuthorityRef `xml:"AuthorityRef"`
	GroupsOfLines *GroupsOfLines `xml:"groupsOfLines"`
}

// GroupsOfLines contains groups of lines
type GroupsOfLines struct {
	GroupsOfLines []*GroupOfLines `xml:"GroupOfLines"`
}

// GroupOfLines represents a group of lines
type GroupOfLines struct {
	BaseNetexObject
	XMLName xml.Name `xml:"GroupOfLines"`
	Name    string   `xml:"Name"`
}

// Lines contains line information
type Lines struct {
	Lines         []*Line         `xml:"Line"`
	FlexibleLines []*FlexibleLine `xml:"FlexibleLine"`
}

// Line represents a public transport line
type Line struct {
	BaseNetexObject
	XMLName              xml.Name              `xml:"Line"`
	Name                 string                `xml:"Name"`
	ShortName            string                `xml:"ShortName"`
	PublicCode           string                `xml:"PublicCode"`
	TransportMode        string                `xml:"TransportMode"`
	TransportSubmode     string                `xml:"TransportSubmode"`
	OperatorRef          *OperatorRef          `xml:"OperatorRef"`
	AuthorityRef         *AuthorityRef         `xml:"AuthorityRef"`
	RepresentedByGroupRef *RepresentedByGroupRef `xml:"RepresentedByGroupRef"`
	Presentation         *Presentation         `xml:"Presentation"`
}

// FlexibleLine represents a flexible transport line
type FlexibleLine struct {
	BaseNetexObject
	XMLName                xml.Name                `xml:"FlexibleLine"`
	Name                   string                  `xml:"Name"`
	ShortName              string                  `xml:"ShortName"`
	PublicCode             string                  `xml:"PublicCode"`
	TransportMode          string                  `xml:"TransportMode"`
	TransportSubmode       string                  `xml:"TransportSubmode"`
	FlexibleLineType       string                  `xml:"FlexibleLineType"`
	BookWhen               string                  `xml:"BookWhen"`
	MinimumBookingPeriod   string                  `xml:"MinimumBookingPeriod"`
	LatestBookingTime      string                  `xml:"LatestBookingTime"`
	BookingContact         *BookingContact         `xml:"BookingContact"`
	BookingUrl             string                  `xml:"BookingUrl"`
	BookingArrangements    *BookingArrangements    `xml:"bookingArrangements"`
	OperatorRef            *OperatorRef            `xml:"OperatorRef"`
	AuthorityRef           *AuthorityRef           `xml:"AuthorityRef"`
}

// BookingContact represents booking contact information
type BookingContact struct {
	Phone string `xml:"Phone"`
	Email string `xml:"Email"`
	Url   string `xml:"Url"`
}

// BookingArrangements represents booking arrangements
type BookingArrangements struct {
	BookingMethods []string `xml:"BookingMethod"`
	BookingAccess  string   `xml:"BookingAccess"`
	BookWhen       string   `xml:"BookWhen"`
	BookingNote    string   `xml:"BookingNote"`
}

// Presentation represents line presentation information
type Presentation struct {
	Colour     string `xml:"Colour"`
	TextColour string `xml:"TextColour"`
}

// Routes contains route information
type Routes struct {
	Routes []*Route `xml:"Route"`
}

// Route represents a route
type Route struct {
	BaseNetexObject
	XMLName          xml.Name          `xml:"Route"`
	Name             string            `xml:"Name"`
	ShortName        string            `xml:"ShortName"`
	LineRef          *LineRef          `xml:"LineRef"`
	DirectionType    string            `xml:"DirectionType"`
	PointsInSequence *PointsInSequence `xml:"pointsInSequence"`
}

// PointsInSequence contains points on a route
type PointsInSequence struct {
	PointOnRoutes []*PointOnRoute `xml:"PointOnRoute"`
}

// PointOnRoute represents a point on a route
type PointOnRoute struct {
	BaseNetexObject
	XMLName                   xml.Name                   `xml:"PointOnRoute"`
	Order                     int                        `xml:"order,attr"`
	ScheduledStopPointRef     *ScheduledStopPointRef     `xml:"ScheduledStopPointRef"`
}

// JourneyPatterns contains journey pattern information
type JourneyPatterns struct {
	JourneyPatterns        []*JourneyPattern        `xml:"JourneyPattern"`
	ServiceJourneyPatterns []*ServiceJourneyPattern `xml:"ServiceJourneyPattern"`
}

// JourneyPattern represents a journey pattern
type JourneyPattern struct {
	BaseNetexObject
	XMLName          xml.Name                     `xml:"JourneyPattern"`
	Name             string                       `xml:"Name"`
	RouteRef         *RouteRef                    `xml:"RouteRef"`
	PointsInSequence *StopPointsInSequence        `xml:"pointsInSequence"`
}

// ServiceJourneyPattern represents a service journey pattern
type ServiceJourneyPattern struct {
	BaseNetexObject
	XMLName          xml.Name              `xml:"ServiceJourneyPattern"`
	Name             string                `xml:"Name"`
	RouteRef         *RouteRef             `xml:"RouteRef"`
	PointsInSequence *StopPointsInSequence `xml:"pointsInSequence"`
}

// StopPointsInSequence contains stop points in a journey pattern
type StopPointsInSequence struct {
	StopPointInJourneyPatterns []*StopPointInJourneyPattern `xml:"StopPointInJourneyPattern"`
}

// StopPointInJourneyPattern represents a stop point in a journey pattern
type StopPointInJourneyPattern struct {
	BaseNetexObject
	XMLName               xml.Name               `xml:"StopPointInJourneyPattern"`
	Order                 int                    `xml:"order,attr"`
	ScheduledStopPointRef *ScheduledStopPointRef `xml:"ScheduledStopPointRef"`
}

// VehicleJourneys contains service journeys
type VehicleJourneys struct {
	ServiceJourneys      []*ServiceJourney      `xml:"ServiceJourney"`
	DatedServiceJourneys []*DatedServiceJourney `xml:"DatedServiceJourney"`
	DeadRuns             []*DeadRun             `xml:"DeadRun"`
}

// ServiceJourney represents a service journey
type ServiceJourney struct {
	BaseNetexObject
	XMLName            xml.Name            `xml:"ServiceJourney"`
	Name               string              `xml:"Name"`
	PrivateCode        string              `xml:"PrivateCode"`
	TransportMode      string              `xml:"TransportMode"`
	TransportSubmode   string              `xml:"TransportSubmode"`
	JourneyPatternRef  *JourneyPatternRef  `xml:"JourneyPatternRef"`
	LineRef            *LineRef            `xml:"LineRef"`
	OperatorRef        *OperatorRef        `xml:"OperatorRef"`
	FlexibleServiceRef *FlexibleServiceRef `xml:"FlexibleServiceRef"`
	BlockRef           *BlockRef           `xml:"BlockRef"`
	DayTypes           *DayTypes           `xml:"dayTypes"`
	PassingTimes       *PassingTimes       `xml:"passingTimes"`
}

// DatedServiceJourney represents a dated service journey
type DatedServiceJourney struct {
	BaseNetexObject
	XMLName            xml.Name            `xml:"DatedServiceJourney"`
	ServiceJourneyRef  *ServiceJourneyRef  `xml:"ServiceJourneyRef"`
	OperatingDayRef    *OperatingDayRef    `xml:"OperatingDayRef"`
}

// DeadRun represents a dead run
type DeadRun struct {
	BaseNetexObject
	XMLName      xml.Name      `xml:"DeadRun"`
	Name         string        `xml:"Name"`
	RouteRef     *RouteRef     `xml:"RouteRef"`
	PassingTimes *PassingTimes `xml:"passingTimes"`
}

// PassingTimes contains timetabled passing times
type PassingTimes struct {
	TimetabledPassingTimes []*TimetabledPassingTime `xml:"TimetabledPassingTime"`
}

// TimetabledPassingTime represents a timetabled passing time
type TimetabledPassingTime struct {
	BaseNetexObject
	XMLName                       xml.Name                       `xml:"TimetabledPassingTime"`
	StopPointInJourneyPatternRef  *StopPointInJourneyPatternRef  `xml:"StopPointInJourneyPatternRef"`
	ArrivalTime                   string                         `xml:"ArrivalTime"`
	DepartureTime                 string                         `xml:"DepartureTime"`
	EarliestDepartureTime         string                         `xml:"EarliestDepartureTime"`
	LatestArrivalTime             string                         `xml:"LatestArrivalTime"`
}

// ScheduledStopPoints contains scheduled stop points
type ScheduledStopPoints struct {
	ScheduledStopPoints []*ScheduledStopPoint `xml:"ScheduledStopPoint"`
}

// ScheduledStopPoint represents a scheduled stop point
type ScheduledStopPoint struct {
	BaseNetexObject
	XMLName xml.Name `xml:"ScheduledStopPoint"`
	Name    string   `xml:"Name"`
}

// StopAssignments contains stop assignments
type StopAssignments struct {
	PassengerStopAssignments []*PassengerStopAssignment `xml:"PassengerStopAssignment"`
}

// PassengerStopAssignment represents a passenger stop assignment
type PassengerStopAssignment struct {
	BaseNetexObject
	XMLName               xml.Name               `xml:"PassengerStopAssignment"`
	ScheduledStopPointRef *ScheduledStopPointRef `xml:"ScheduledStopPointRef"`
	StopPlaceRef          *StopPlaceRef          `xml:"StopPlaceRef"`
	QuayRef               *QuayRef               `xml:"QuayRef"`
}

// Interchanges contains service journey interchanges
type Interchanges struct {
	ServiceJourneyInterchanges []*ServiceJourneyInterchange `xml:"ServiceJourneyInterchange"`
}

// ServiceJourneyInterchange represents an interchange between service journeys
type ServiceJourneyInterchange struct {
	BaseNetexObject
	XMLName               xml.Name               `xml:"ServiceJourneyInterchange"`
	FromStopPointRef      *ScheduledStopPointRef `xml:"FromStopPointRef"`
	ToStopPointRef        *ScheduledStopPointRef `xml:"ToStopPointRef"`
	FromServiceJourneyRef *ServiceJourneyRef     `xml:"FromServiceJourneyRef"`
	ToServiceJourneyRef   *ServiceJourneyRef     `xml:"ToServiceJourneyRef"`
	StandardTransferTime  string                 `xml:"StandardTransferTime"`
}

// StopPlaces contains stop places
type StopPlaces struct {
	StopPlaces []*StopPlace `xml:"StopPlace"`
}

// StopPlace represents a stop place
type StopPlace struct {
	BaseNetexObject
	XMLName       xml.Name       `xml:"StopPlace"`
	Name          string         `xml:"Name"`
	ShortName     string         `xml:"ShortName"`
	StopPlaceType string         `xml:"StopPlaceType"`
	Centroid      *Centroid      `xml:"Centroid"`
	Quays         *Quays         `xml:"quays"`
}

// Quays contains quays
type Quays struct {
	Quays []*Quay `xml:"Quay"`
}

// Quay represents a quay
type Quay struct {
	BaseNetexObject
	XMLName     xml.Name     `xml:"Quay"`
	Name        string       `xml:"Name"`
	ShortName   string       `xml:"ShortName"`
	Centroid    *Centroid    `xml:"Centroid"`
}

// Centroid represents geographic coordinates
type Centroid struct {
	Location *Location `xml:"Location"`
}

// Location represents a geographic location
type Location struct {
	Longitude float64 `xml:"Longitude"`
	Latitude  float64 `xml:"Latitude"`
}

// DayTypes contains day types
type DayTypes struct {
	DayTypes    []*DayType    `xml:"DayType"`
	DayTypeRefs []*DayTypeRef `xml:"DayTypeRef"`
}

// DayType represents a day type
type DayType struct {
	BaseNetexObject
	XMLName xml.Name `xml:"DayType"`
	Name    string   `xml:"Name"`
}

// OperatingDays contains operating days
type OperatingDays struct {
	OperatingDays []*OperatingDay `xml:"OperatingDay"`
}

// OperatingDay represents an operating day
type OperatingDay struct {
	BaseNetexObject
	XMLName      xml.Name `xml:"OperatingDay"`
	CalendarDate string   `xml:"CalendarDate"`
}

// ServiceCalendar represents a service calendar
type ServiceCalendar struct {
	BaseNetexObject
	XMLName          xml.Name          `xml:"ServiceCalendar"`
	FromDate         string            `xml:"FromDate"`
	ToDate           string            `xml:"ToDate"`
	OperatingPeriods *OperatingPeriods `xml:"operatingPeriods"`
}

// OperatingPeriods contains operating periods
type OperatingPeriods struct {
	OperatingPeriods []*OperatingPeriod `xml:"OperatingPeriod"`
}

// OperatingPeriod represents an operating period
type OperatingPeriod struct {
	BaseNetexObject
	XMLName  xml.Name `xml:"OperatingPeriod"`
	FromDate string   `xml:"FromDate"`
	ToDate   string   `xml:"ToDate"`
}

// Blocks contains blocks
type Blocks struct {
	Blocks []*Block `xml:"Block"`
}

// Block represents a vehicle block
type Block struct {
	BaseNetexObject
	XMLName  xml.Name  `xml:"Block"`
	Name     string    `xml:"Name"`
	Journeys *Journeys `xml:"journeys"`
}

// Journeys contains vehicle journey references
type Journeys struct {
	VehicleJourneyRefs []*VehicleJourneyRef `xml:"VehicleJourneyRef"`
}

// VehicleTypes contains vehicle types
type VehicleTypes struct {
	VehicleTypes []*VehicleType `xml:"VehicleType"`
}

// VehicleType represents a vehicle type
type VehicleType struct {
	BaseNetexObject
	XMLName xml.Name `xml:"VehicleType"`
	Name    string   `xml:"Name"`
}

// ValidityConditions contains validity conditions
type ValidityConditions struct {
	AvailabilityConditions []*AvailabilityCondition `xml:"AvailabilityCondition"`
}

// AvailabilityCondition represents an availability condition
type AvailabilityCondition struct {
	BaseNetexObject
	XMLName  xml.Name `xml:"AvailabilityCondition"`
	FromDate string   `xml:"FromDate"`
	ToDate   string   `xml:"ToDate"`
}

// Reference types
type OperatorRef struct {
	Ref string `xml:"ref,attr"`
}

type AuthorityRef struct {
	Ref string `xml:"ref,attr"`
}

type RepresentedByGroupRef struct {
	Ref string `xml:"ref,attr"`
}

type LineRef struct {
	Ref string `xml:"ref,attr"`
}

type RouteRef struct {
	Ref string `xml:"ref,attr"`
}

type JourneyPatternRef struct {
	Ref string `xml:"ref,attr"`
}

type ServiceJourneyRef struct {
	Ref string `xml:"ref,attr"`
}

type OperatingDayRef struct {
	Ref string `xml:"ref,attr"`
}

type FlexibleServiceRef struct {
	Ref string `xml:"ref,attr"`
}

type BlockRef struct {
	Ref string `xml:"ref,attr"`
}

type ScheduledStopPointRef struct {
	Ref string `xml:"ref,attr"`
}

type StopPlaceRef struct {
	Ref string `xml:"ref,attr"`
}

type QuayRef struct {
	Ref string `xml:"ref,attr"`
}

type StopPointInJourneyPatternRef struct {
	Ref string `xml:"ref,attr"`
}

type VehicleJourneyRef struct {
	Ref string `xml:"ref,attr"`
}

type DayTypeRef struct {
	Ref string `xml:"ref,attr"`
}