package idvalidation

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/theoremus-urban-solutions/netex-validator/types"
)

// NetexIdRepository manages NetEX ID tracking and validation across files
type NetexIdRepository struct {
	// Map of ID -> IdVersion for tracking all IDs across files
	ids map[string]types.IdVersion
	// Map of filename -> set of IDs for file-specific tracking
	fileIds map[string]map[string]bool
	// Map of ID -> references for tracking unresolved references
	references map[string][]types.IdVersion
	// Map of ID -> map[fileName]version for cross-file consistency checks
	idToFiles map[string]map[string]string
	// Map of filename -> bool for tracking common files
	commonFiles map[string]bool
	// Set of element names to ignore for ID uniqueness validation
	ignorableElements map[string]bool
	// Thread safety
	mu sync.RWMutex
}

// NewNetexIdRepository creates a new ID repository
func NewNetexIdRepository() *NetexIdRepository {
	return NewNetexIdRepositoryWithIgnorableElements(getDefaultIgnorableElements())
}

// NewNetexIdRepositoryWithIgnorableElements creates a new ID repository with custom ignorable elements
func NewNetexIdRepositoryWithIgnorableElements(ignorableElements []string) *NetexIdRepository {
	ignorableMap := make(map[string]bool)
	for _, elem := range ignorableElements {
		ignorableMap[elem] = true
	}

	return &NetexIdRepository{
		ids:               make(map[string]types.IdVersion),
		fileIds:           make(map[string]map[string]bool),
		references:        make(map[string][]types.IdVersion),
		idToFiles:         make(map[string]map[string]string),
		commonFiles:       make(map[string]bool),
		ignorableElements: ignorableMap,
	}
}

// getDefaultIgnorableElements returns the default set of elements to ignore for ID uniqueness
func getDefaultIgnorableElements() []string {
	return []string{
		"ResourceFrame",
		"SiteFrame",
		"CompositeFrame",
		"TimetableFrame",
		"ServiceFrame",
		"ServiceCalendarFrame",
		"VehicleScheduleFrame",
		"Block",
		"RoutePoint",
		"PointProjection",
		"ScheduledStopPoint",
		"PassengerStopAssignment",
		"NoticeAssignment",
		"ServiceLinkInJourneyPattern",
		"ServiceFacilitySet",
		"AvailabilityCondition",
	}
}

// AddId registers a NetEX ID in the repository
func (r *NetexIdRepository) AddId(id, version, fileName string) error {
	return r.AddIdWithElementType(id, version, fileName, "")
}

// AddIdWithElementType registers a NetEX ID in the repository with element type information
func (r *NetexIdRepository) AddIdWithElementType(id, version, fileName, elementType string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if this element type should be ignored
	if elementType != "" && r.ignorableElements[elementType] {
		return nil // Skip registration for ignorable elements
	}

	// Check for duplicates
	if existing, exists := r.ids[id]; exists {
		if existing.FileName != fileName {
			return fmt.Errorf("duplicate NetEX ID '%s' found in files '%s' and '%s'",
				id, existing.FileName, fileName)
		}
		// Same file, check version
		if existing.Version != version {
			return fmt.Errorf("NetEX ID '%s' has conflicting versions '%s' and '%s' in file '%s'",
				id, existing.Version, version, fileName)
		}
	}

	// Register the ID
	idVersion := types.NewIdVersion(id, version, fileName)
	r.ids[id] = idVersion

	// Track by file
	if r.fileIds[fileName] == nil {
		r.fileIds[fileName] = make(map[string]bool)
	}
	r.fileIds[fileName][id] = true

	// Track versions per file for cross-file consistency
	if r.idToFiles[id] == nil {
		r.idToFiles[id] = make(map[string]string)
	}
	r.idToFiles[id][fileName] = version

	return nil
}

// AddReference registers a reference to a NetEX ID
func (r *NetexIdRepository) AddReference(refId, version, fileName string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	refVersion := types.NewIdVersion(refId, version, fileName)
	r.references[refId] = append(r.references[refId], refVersion)
}

// ValidateReferences validates all references against registered IDs
func (r *NetexIdRepository) ValidateReferences() []types.ValidationIssue {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var issues []types.ValidationIssue

	for refId, references := range r.references {
		// Check if the referenced ID exists
		if _, exists := r.ids[refId]; !exists {
			for _, ref := range references {
				issues = append(issues, types.ValidationIssue{
					Rule: types.ValidationRule{
						Code:     "NETEX_ID_5",
						Name:     "NeTEx ID unresolved reference",
						Message:  fmt.Sprintf("Unresolved reference to NetEX ID '%s'", refId),
						Severity: types.ERROR,
					},
					Location: types.DataLocation{
						FileName:  ref.FileName,
						ElementID: refId,
					},
					Message: fmt.Sprintf("Unresolved reference to NetEX ID '%s' from file '%s'", refId, ref.FileName),
				})
			}
		} else {
			// ID exists, validate version consistency if specified
			targetId := r.ids[refId]
			for _, ref := range references {
				if ref.Version != "" && ref.Version != "any" && targetId.Version != "" && ref.Version != targetId.Version {
					issues = append(issues, types.ValidationIssue{
						Rule: types.ValidationRule{
							Code:     "NETEX_ID_9",
							Name:     "NeTEx ID version mismatch on reference",
							Message:  fmt.Sprintf("Version mismatch for reference to NetEX ID '%s'", refId),
							Severity: types.WARNING,
						},
						Location: types.DataLocation{
							FileName:  ref.FileName,
							ElementID: refId,
						},
						Message: fmt.Sprintf("Reference to NetEX ID '%s' has version '%s' but target has version '%s'",
							refId, ref.Version, targetId.Version),
					})
				}
				// If target has a version but reference does not, warn about missing version
				if (ref.Version == "" || ref.Version == "any") && targetId.Version != "" {
					issues = append(issues, types.ValidationIssue{
						Rule: types.ValidationRule{
							Code:     "NETEX_ID_11",
							Name:     "NeTEx ID missing version on reference",
							Message:  fmt.Sprintf("Reference to NetEX ID '%s' is missing version while target has version '%s'", refId, targetId.Version),
							Severity: types.WARNING,
						},
						Location: types.DataLocation{
							FileName:  ref.FileName,
							ElementID: refId,
						},
						Message: fmt.Sprintf("Reference to NetEX ID '%s' is missing version while target has version '%s'",
							refId, targetId.Version),
					})
				}
			}
		}
	}

	return issues
}

// ValidateIdFormat validates NetEX ID format compliance
func (r *NetexIdRepository) ValidateIdFormat() []types.ValidationIssue {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var issues []types.ValidationIssue

	for id, idVersion := range r.ids {
		if !r.isValidNetexIdFormat(id) {
			issues = append(issues, types.ValidationIssue{
				Rule: types.ValidationRule{
					Code:     "NETEX_ID_7",
					Name:     "NeTEx ID invalid value",
					Message:  fmt.Sprintf("NetEX ID '%s' has invalid format", id),
					Severity: types.ERROR,
				},
				Location: types.DataLocation{
					FileName:  idVersion.FileName,
					ElementID: id,
				},
				Message: fmt.Sprintf("NetEX ID '%s' does not follow the required format", id),
			})
		}
	}

	return issues
}

// ValidateVersions validates version information on IDs
func (r *NetexIdRepository) ValidateVersions() []types.ValidationIssue {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var issues []types.ValidationIssue

	for id, idVersion := range r.ids {
		// Check for missing version
		if idVersion.Version == "" {
			issues = append(issues, types.ValidationIssue{
				Rule: types.ValidationRule{
					Code:     "NETEX_ID_8",
					Name:     "NeTEx ID missing version on elements",
					Message:  fmt.Sprintf("NetEX ID '%s' is missing version information", id),
					Severity: types.WARNING,
				},
				Location: types.DataLocation{
					FileName:  idVersion.FileName,
					ElementID: id,
				},
				Message: fmt.Sprintf("NetEX ID '%s' is missing version attribute", id),
			})
		}

		// Check for non-numeric version
		if idVersion.Version != "" && idVersion.Version != "any" && !r.isNumericVersion(idVersion.Version) {
			issues = append(issues, types.ValidationIssue{
				Rule: types.ValidationRule{
					Code:     "VERSION_NON_NUMERIC",
					Name:     "Non-numeric NeTEx version",
					Message:  fmt.Sprintf("NetEX ID '%s' has non-numeric version", id),
					Severity: types.WARNING,
				},
				Location: types.DataLocation{
					FileName:  idVersion.FileName,
					ElementID: id,
				},
				Message: fmt.Sprintf("NetEX ID '%s' has non-numeric version '%s'", id, idVersion.Version),
			})
		}
	}

	return issues
}

// ValidateVersionConsistencyAcrossFiles checks that the same ID has consistent version across files
func (r *NetexIdRepository) ValidateVersionConsistencyAcrossFiles() []types.ValidationIssue {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var issues []types.ValidationIssue
	for id, fileToVersion := range r.idToFiles {
		var seenVersion string
		for _, v := range fileToVersion {
			if v == "" || v == "any" {
				continue
			}
			if seenVersion == "" {
				seenVersion = v
			} else if v != seenVersion {
				// version mismatch across files
				// collect files for context
				var files []string
				for fn := range fileToVersion {
					files = append(files, fn)
				}
				issues = append(issues, types.ValidationIssue{
					Rule: types.ValidationRule{
						Code:     "NETEX_ID_10",
						Name:     "NeTEx ID version mismatch across files",
						Message:  fmt.Sprintf("ID '%s' has conflicting versions across files", id),
						Severity: types.ERROR,
					},
					Location: types.DataLocation{ElementID: id},
					Message:  fmt.Sprintf("ID '%s' appears with different versions across files: %v", id, files),
				})
				break
			}
		}
	}
	return issues
}

// GetIdsByFile returns all IDs registered for a specific file
func (r *NetexIdRepository) GetIdsByFile(fileName string) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var ids []string
	if fileIds, exists := r.fileIds[fileName]; exists {
		for id := range fileIds {
			ids = append(ids, id)
		}
	}
	return ids
}

// GetAllIds returns all registered IDs
func (r *NetexIdRepository) GetAllIds() map[string]types.IdVersion {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[string]types.IdVersion)
	for id, idVersion := range r.ids {
		result[id] = idVersion
	}
	return result
}

// MarkAsCommonFile marks a file as a common file for special duplicate ID handling
func (r *NetexIdRepository) MarkAsCommonFile(fileName string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.commonFiles[fileName] = true
}

// IsCommonFile returns true if the file is marked as a common file
func (r *NetexIdRepository) IsCommonFile(fileName string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.commonFiles[fileName]
}

// Clear resets the repository
func (r *NetexIdRepository) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.ids = make(map[string]types.IdVersion)
	r.fileIds = make(map[string]map[string]bool)
	r.references = make(map[string][]types.IdVersion)
	r.idToFiles = make(map[string]map[string]string)
	r.commonFiles = make(map[string]bool)
}

// isValidNetexIdFormat validates NetEX ID format (flexible validation)
func (r *NetexIdRepository) isValidNetexIdFormat(id string) bool {
	// NetEX ID format supports multiple patterns:
	// 1. EU format: Codespace:EntityType:Identifier[:...]
	// 2. French format: FR:NumericCode:EntityType:ComplexIdentifier:RIV
	// 3. Simple names: monomodalStopPlace, multimodalStopPlace, etc.

	if len(id) < 3 {
		return false
	}

	// Allow simple descriptive names (used in French data)
	simpleNames := regexp.MustCompile(`^(monomodalStopPlace|multimodalStopPlace|monomodalHub|multimodalHub)$`)
	if simpleNames.MatchString(id) {
		return true
	}

	// Must contain at least one colon for structured IDs
	if !strings.Contains(id, ":") {
		return false
	}

	// Split and normalize tokens by removing empty segments (handles '::')
	raw := strings.Split(id, ":")
	tokens := make([]string, 0, len(raw))
	for _, t := range raw {
		if t != "" {
			tokens = append(tokens, t)
		}
	}

	if len(tokens) < 3 {
		return false
	}

	// Find entity type token (can be at position 1 or 2 depending on format)
	var entity string
	if len(tokens) >= 3 {
		// Check if token[1] looks like an entity type
		entity = tokens[1]
		// If token[1] is numeric, entity type is likely at token[2] (French format)
		if regexp.MustCompile(`^\d+$`).MatchString(tokens[1]) && len(tokens) >= 4 {
			entity = tokens[2]
		}
	}

	// Allowed entity types (EU/French NetEX profile) - comprehensive list
	allowed := regexp.MustCompile(`^(Line|FlexibleLine|Route|JourneyPattern|ServiceJourney|DatedServiceJourney|Operator|Authority|Network|ScheduledStopPoint|StopPlace|Quay|TariffZone|FareZone|GroupOfLines|GroupOfServices|Block|CourseOfJourneys|DeadRun|Interchange|Notice|FlexibleService|Location|Centroid|PostalAddress|AccessibilityAssessment|PlaceEquipment|SiteEquipment|AccessSpace|BoardingPosition|PathLink|PathJunction|Connection|SiteConnection|TopographicPlace|AddressablePlace|PointOfInterest|Parking|Zone|TransportZone|AccessZone|StopArea|Area|MultimodalStopPlace|StopPlaceSpace|StopPlaceComponent|StopPlaceEntrance|PathLinkEnd|EquipmentPlace|LocalService|StopPlaceRef|QuayRef|ScheduledStopPointRef|LocationRef|GroupOfStopPlaces|Timetable|StopPointInJourneyPattern|PassengerStopAssignment|ServiceJourneyPattern|DestinationDisplay|Direction|GroupOfLine|Company|CompositeFrame|GeneralFrame|ResourceFrame|ValueSet|TypeOfFrame)$`)

	if !allowed.MatchString(entity) {
		return false
	}

	return true
}

// isNumericVersion checks if version is numeric
func (r *NetexIdRepository) isNumericVersion(version string) bool {
	if version == "" || version == "any" {
		return true
	}

	// Simple numeric check - can be enhanced
	for _, r := range version {
		if r < '0' || r > '9' {
			if r != '.' && r != '-' {
				return false
			}
		}
	}
	return true
}

// GetDuplicateIds returns IDs that appear in multiple files
func (r *NetexIdRepository) GetDuplicateIds() []types.ValidationIssue {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var issues []types.ValidationIssue
	fileCount := make(map[string]map[string]bool)

	// Count files per ID
	for fileName, ids := range r.fileIds {
		for id := range ids {
			if fileCount[id] == nil {
				fileCount[id] = make(map[string]bool)
			}
			fileCount[id][fileName] = true
		}
	}

	// Find duplicates
	for id, files := range fileCount {
		if len(files) > 1 {
			var fileNames []string
			var commonFileNames []string
			var regularFileNames []string

			for fileName := range files {
				fileNames = append(fileNames, fileName)
				if r.commonFiles[fileName] {
					commonFileNames = append(commonFileNames, fileName)
				} else {
					regularFileNames = append(regularFileNames, fileName)
				}
			}

			idVersion := r.ids[id]

			// Determine rule based on file types
			var rule types.ValidationRule
			if len(commonFileNames) > 0 && len(regularFileNames) > 0 {
				// Mixed common and regular files - use regular duplicate rule
				rule = types.ValidationRule{
					Code:     "NETEX_ID_1",
					Name:     "NeTEx ID duplicated across files",
					Message:  fmt.Sprintf("NetEX ID '%s' appears in multiple files", id),
					Severity: types.ERROR,
				}
			} else if len(commonFileNames) > 1 {
				// Only common files - use common duplicate rule (warning)
				rule = types.ValidationRule{
					Code:     "NETEX_ID_10",
					Name:     "Duplicate NeTEx ID across common files",
					Message:  fmt.Sprintf("NetEX ID '%s' appears in multiple common files", id),
					Severity: types.WARNING,
				}
			} else {
				// Regular duplicate across non-common files
				rule = types.ValidationRule{
					Code:     "NETEX_ID_1",
					Name:     "NeTEx ID duplicated across files",
					Message:  fmt.Sprintf("NetEX ID '%s' appears in multiple files", id),
					Severity: types.ERROR,
				}
			}

			issues = append(issues, types.ValidationIssue{
				Rule: rule,
				Location: types.DataLocation{
					FileName:  idVersion.FileName,
					ElementID: id,
				},
				Message: fmt.Sprintf("NetEX ID '%s' is duplicated in files: %s", id, strings.Join(fileNames, ", ")),
			})
		}
	}

	return issues
}
