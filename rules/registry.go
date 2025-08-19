package rules

import (
	"github.com/theoremus-urban-solutions/netex-validator/config"
	"github.com/theoremus-urban-solutions/netex-validator/types"
)

// Rule represents a single validation rule
type Rule struct {
	Code        string
	Name        string
	Message     string
	Severity    types.Severity
	XPath       string
	Category    string
	Description string
}

// RuleRegistry manages all validation rules
type RuleRegistry struct {
	rules   []Rule
	config  *config.ValidatorConfig
	profile string
}

// NewRuleRegistry creates a new rule registry with all built-in rules
func NewRuleRegistry(cfg *config.ValidatorConfig) *RuleRegistry {
	registry := &RuleRegistry{
		rules:  make([]Rule, 0),
		config: cfg,
	}

	// Load all built-in rules
	registry.loadBuiltinRules()

	return registry
}

// WithProfile allows selecting a ruleset profile (e.g., "eu", "custom").
func (r *RuleRegistry) WithProfile(profile string) *RuleRegistry {
	r.profile = profile
	return r
}

// GetEnabledRules returns all enabled rules based on configuration
func (r *RuleRegistry) GetEnabledRules() []Rule {
	var enabled []Rule

	for _, rule := range r.rules {
		// Filter by profile: for EU, include only generic EU-safe categories
		if r.profile == "eu" && !isEUCategory(rule.Category) {
			continue
		}
		if r.config.IsRuleEnabled(rule.Code) {
			// Apply severity overrides from config
			effectiveSeverity := r.config.GetRuleSeverity(rule.Code, rule.Severity)
			rule.Severity = effectiveSeverity
			enabled = append(enabled, rule)
		}
	}

	// Add custom rules from config
	for _, customRule := range r.config.GetCustomRules() {
		enabled = append(enabled, Rule{
			Code:     customRule.Code,
			Name:     customRule.Name,
			Message:  customRule.Message,
			Severity: customRule.Severity,
			XPath:    customRule.XPath,
			Category: "custom",
		})
	}

	return enabled
}

// isEUCategory returns true if a rule category is part of the generic EU profile
func isEUCategory(category string) bool {
	// Conservative allow-list; expand as EU set is curated
	switch category {
	case "line", "route", "transport_mode", "version", "journey_pattern", "stop_point", "calendar", "validity", "interchange", "group", "tariff_zone", "responsibility_set", "type_of_service":
		return true
	default:
		return false
	}
}

// GetRuleByCode returns a specific rule by its code
func (r *RuleRegistry) GetRuleByCode(code string) (Rule, bool) {
	for _, rule := range r.rules {
		if rule.Code == code {
			return rule, true
		}
	}
	return Rule{}, false
}

// GetRulesByCategory returns all rules in a specific category
func (r *RuleRegistry) GetRulesByCategory(category string) []Rule {
	var categoryRules []Rule

	for _, rule := range r.rules {
		if rule.Category == category {
			categoryRules = append(categoryRules, rule)
		}
	}

	return categoryRules
}

// GetAllCategories returns all available rule categories
func (r *RuleRegistry) GetAllCategories() []string {
	categories := make(map[string]bool)

	for _, rule := range r.rules {
		categories[rule.Category] = true
	}

	result := make([]string, 0, len(categories))
	for category := range categories {
		result = append(result, category)
	}

	return result
}

// addRule is a helper function to add a rule to the registry
func (r *RuleRegistry) addRule(code, name, message string, severity types.Severity, xpath string) {
	rule := Rule{
		Code:     code,
		Name:     name,
		Message:  message,
		Severity: severity,
		XPath:    xpath,
		Category: r.getCategoryFromCode(code),
	}

	r.rules = append(r.rules, rule)
}

// getCategoryFromCode determines the category from rule code
func (r *RuleRegistry) getCategoryFromCode(code string) string {
	categoryMap := map[string]string{
		"LINE_":                  "line",
		"ROUTE_":                 "route",
		"SERVICE_JOURNEY_":       "service_journey",
		"FLEXIBLE_LINE_":         "flexible_line",
		"NETWORK_":               "network",
		"AUTHORITY_":             "network",
		"OPERATOR_":              "network",
		"JOURNEY_PATTERN_":       "journey_pattern",
		"STOP_POINT_":            "stop_point",
		"SCHEDULED_STOP_":        "stop_point",
		"VERSION_":               "version",
		"TRANSPORT_MODE_":        "transport_mode",
		"TRANSPORT_SUB_MODE_":    "transport_mode",
		"BOOKING_":               "booking",
		"FLEXIBLE_LINE_TYPE_":    "flexible_line",
		"LINE_INVALID_COLOR":     "line",
		"SERVICE_CALENDAR_":      "calendar",
		"VALIDITY_CONDITIONS_":   "validity",
		"DATED_SERVICE_JOURNEY_": "dated_service_journey",
		"DEAD_RUN_":              "dead_run",
		"INTERCHANGE_":           "interchange",
		"NOTICE_":                "notice",
		"COMPOSITE_FRAME_":       "frame",
		"TIMETABLE_FRAME_":       "frame",
		"SERVICE_FRAME_":         "frame",
		"RESOURCE_FRAME_":        "frame",
		"SITE_FRAME_":            "frame",
		"INFRASTRUCTURE_FRAME_":  "frame",
		"FLEXIBLE_SERVICE_":      "flexible_service",
		"FLEXIBLE_STOP_":         "flexible_service",
		"FLEXIBLE_AREA_":         "flexible_service",
		"BLOCK_":                 "block",
		"COURSE_OF_JOURNEYS_":    "course_of_journeys",
		"TARIFF_ZONE_":           "tariff_zone",
		"RESPONSIBILITY_SET_":    "responsibility_set",
		"TYPE_OF_SERVICE_":       "type_of_service",
		"GROUP_OF_":              "group",
		"FARE_":                  "group",
	}

	for prefix, category := range categoryMap {
		if len(code) >= len(prefix) && code[:len(prefix)] == prefix {
			return category
		}
	}

	return "custom"
}
