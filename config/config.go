package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/theoremus-urban-solutions/netex-validator/types"
	"gopkg.in/yaml.v3"
)

// ValidatorConfig represents the complete validator configuration
type ValidatorConfig struct {
	Validator ValidatorSettings `yaml:"validator"`
	Rules     RulesConfig       `yaml:"rules"`
	Output    OutputConfig      `yaml:"output"`
}

// ValidatorSettings contains general validator settings
type ValidatorSettings struct {
	Profile         string `yaml:"profile"`         // e.g., "eu", "custom"
	MaxFileSize     int64  `yaml:"maxFileSize"`     // Maximum file size in bytes
	MaxSchemaErrors int    `yaml:"maxSchemaErrors"` // Maximum schema errors to report
	ConcurrentFiles int    `yaml:"concurrentFiles"` // Number of files to process concurrently
	EnableCache     bool   `yaml:"enableCache"`     // Enable validation result caching
	CacheTimeout    int    `yaml:"cacheTimeout"`    // Cache timeout in minutes
}

// RulesConfig contains rule-specific configuration
type RulesConfig struct {
	Categories map[string]RuleCategoryConfig `yaml:"categories"`
	Custom     []CustomRuleConfig            `yaml:"custom"`
}

// RuleCategoryConfig configures an entire category of rules
type RuleCategoryConfig struct {
	Enabled         bool                  `yaml:"enabled"`
	DefaultSeverity *types.Severity       `yaml:"defaultSeverity,omitempty"`
	Rules           map[string]RuleConfig `yaml:"rules,omitempty"`
}

// RuleConfig configures a specific validation rule
type RuleConfig struct {
	Enabled  bool            `yaml:"enabled"`
	Severity *types.Severity `yaml:"severity,omitempty"`
	Message  string          `yaml:"message,omitempty"`
	XPath    string          `yaml:"xpath,omitempty"`
}

// CustomRuleConfig defines a custom validation rule
type CustomRuleConfig struct {
	Code     string         `yaml:"code"`
	Name     string         `yaml:"name"`
	Message  string         `yaml:"message"`
	Severity types.Severity `yaml:"severity"`
	XPath    string         `yaml:"xpath"`
	Enabled  bool           `yaml:"enabled"`
}

// OutputConfig configures output settings
type OutputConfig struct {
	Format          string `yaml:"format"`          // json, text, html
	IncludeDetails  bool   `yaml:"includeDetails"`  // Include detailed location info
	GroupBySeverity bool   `yaml:"groupBySeverity"` // Group output by severity
	MaxEntries      int    `yaml:"maxEntries"`      // Maximum entries to output (0 = unlimited)
}

// DefaultConfig returns a default configuration
func DefaultConfig() *ValidatorConfig {
	return &ValidatorConfig{
		Validator: ValidatorSettings{
			Profile:         "eu",
			MaxFileSize:     100 * 1024 * 1024, // 100MB
			MaxSchemaErrors: 100,
			ConcurrentFiles: 4,
			EnableCache:     false,
			CacheTimeout:    30,
		},
		Rules: RulesConfig{
			Categories: map[string]RuleCategoryConfig{
				"line": {
					Enabled:         true,
					DefaultSeverity: nil, // Use rule default
				},
				"route": {
					Enabled:         true,
					DefaultSeverity: nil,
				},
				"service_journey": {
					Enabled:         true,
					DefaultSeverity: nil,
				},
				"flexible_line": {
					Enabled:         true,
					DefaultSeverity: nil,
				},
				"network": {
					Enabled:         true,
					DefaultSeverity: nil,
				},
				"journey_pattern": {
					Enabled:         true,
					DefaultSeverity: nil,
				},
				"stop_point": {
					Enabled:         true,
					DefaultSeverity: nil,
				},
				"version": {
					Enabled:         true,
					DefaultSeverity: nil,
				},
				"transport_mode": {
					Enabled:         true,
					DefaultSeverity: nil,
				},
				"booking": {
					Enabled:         true,
					DefaultSeverity: nil,
				},
				"calendar": {
					Enabled:         true,
					DefaultSeverity: nil,
				},
				"validity": {
					Enabled:         true,
					DefaultSeverity: nil,
				},
				"dated_service_journey": {
					Enabled:         true,
					DefaultSeverity: nil,
				},
				"dead_run": {
					Enabled:         true,
					DefaultSeverity: nil,
				},
				"interchange": {
					Enabled:         true,
					DefaultSeverity: nil,
				},
				"notice": {
					Enabled:         true,
					DefaultSeverity: nil,
				},
				"frame": {
					Enabled:         true,
					DefaultSeverity: nil,
				},
				"flexible_service": {
					Enabled:         true,
					DefaultSeverity: nil,
				},
				"block": {
					Enabled:         true,
					DefaultSeverity: nil,
				},
				"course_of_journeys": {
					Enabled:         true,
					DefaultSeverity: nil,
				},
				"tariff_zone": {
					Enabled:         true,
					DefaultSeverity: nil,
				},
				"responsibility_set": {
					Enabled:         true,
					DefaultSeverity: nil,
				},
				"type_of_service": {
					Enabled:         true,
					DefaultSeverity: nil,
				},
				"group": {
					Enabled:         true,
					DefaultSeverity: nil,
				},
			},
			Custom: []CustomRuleConfig{},
		},
		Output: OutputConfig{
			Format:          "json",
			IncludeDetails:  true,
			GroupBySeverity: true,
			MaxEntries:      0, // Unlimited
		},
	}
}

// LoadConfig loads configuration from a YAML file
func LoadConfig(configPath string) (*ValidatorConfig, error) {
	// Start with default config
	config := DefaultConfig()

	// If no config file specified, return default
	if configPath == "" {
		return config, nil
	}

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("configuration file not found: %s", configPath)
	}

	// Validate file path to prevent path traversal
	if !filepath.IsAbs(configPath) && strings.Contains(configPath, "..") {
		return nil, fmt.Errorf("invalid config file path: %s", configPath)
	}

	// Read file
	data, err := os.ReadFile(configPath) //nolint:gosec // Path is validated above
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file: %w", err)
	}

	// Parse YAML
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse configuration file: %w", err)
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// SaveConfig saves configuration to a YAML file
func (c *ValidatorConfig) SaveConfig(configPath string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal to YAML
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}

	// Write file
	if err := os.WriteFile(configPath, data, 0o600); err != nil {
		return fmt.Errorf("failed to write configuration file: %w", err)
	}

	return nil
}

// Validate validates the configuration
func (c *ValidatorConfig) Validate() error {
	// Validate validator settings
	if c.Validator.MaxFileSize <= 0 {
		return fmt.Errorf("maxFileSize must be positive")
	}

	if c.Validator.MaxSchemaErrors < 0 {
		return fmt.Errorf("maxSchemaErrors cannot be negative")
	}

	if c.Validator.ConcurrentFiles <= 0 {
		return fmt.Errorf("concurrentFiles must be positive")
	}

	// Validate output format
	validFormats := map[string]bool{"json": true, "text": true, "html": true}
	if !validFormats[c.Output.Format] {
		return fmt.Errorf("invalid output format: %s (valid: json, text, html)", c.Output.Format)
	}

	// Validate custom rules
	for i, rule := range c.Rules.Custom {
		if rule.Code == "" {
			return fmt.Errorf("custom rule %d: code cannot be empty", i)
		}
		if rule.Name == "" {
			return fmt.Errorf("custom rule %d: name cannot be empty", i)
		}
		if rule.XPath == "" {
			return fmt.Errorf("custom rule %d: xpath cannot be empty", i)
		}
	}

	return nil
}

// IsRuleEnabled checks if a specific rule is enabled
func (c *ValidatorConfig) IsRuleEnabled(ruleCode string) bool {
	// Determine category from rule code
	category := getRuleCategoryFromCode(ruleCode)

	// Check if category is enabled
	if catConfig, exists := c.Rules.Categories[category]; exists {
		if !catConfig.Enabled {
			return false
		}

		// Check specific rule override
		if ruleConfig, exists := catConfig.Rules[ruleCode]; exists {
			return ruleConfig.Enabled
		}

		// Default to enabled if category is enabled
		return true
	}

	// Check custom rules
	for _, customRule := range c.Rules.Custom {
		if customRule.Code == ruleCode {
			return customRule.Enabled
		}
	}

	// Default to enabled
	return true
}

// GetRuleSeverity gets the effective severity for a rule
func (c *ValidatorConfig) GetRuleSeverity(ruleCode string, defaultSeverity types.Severity) types.Severity {
	// Determine category from rule code
	category := getRuleCategoryFromCode(ruleCode)

	// Check category configuration
	if catConfig, exists := c.Rules.Categories[category]; exists {
		// Check specific rule override
		if ruleConfig, exists := catConfig.Rules[ruleCode]; exists {
			if ruleConfig.Severity != nil {
				return *ruleConfig.Severity
			}
		}

		// Check category default
		if catConfig.DefaultSeverity != nil {
			return *catConfig.DefaultSeverity
		}
	}

	// Check custom rules
	for _, customRule := range c.Rules.Custom {
		if customRule.Code == ruleCode {
			return customRule.Severity
		}
	}

	// Return original default
	return defaultSeverity
}

// GetCustomRules returns all enabled custom rules
func (c *ValidatorConfig) GetCustomRules() []CustomRuleConfig {
	var enabled []CustomRuleConfig
	for _, rule := range c.Rules.Custom {
		if rule.Enabled {
			enabled = append(enabled, rule)
		}
	}
	return enabled
}

// getRuleCategoryFromCode determines the rule category from rule code
func getRuleCategoryFromCode(ruleCode string) string {
	if len(ruleCode) == 0 {
		return "unknown"
	}

	// Map rule codes to categories
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
		if len(ruleCode) >= len(prefix) && ruleCode[:len(prefix)] == prefix {
			return category
		}
	}

	return "custom"
}

// GenerateDefaultConfigFile creates a default configuration file
func GenerateDefaultConfigFile(configPath string) error {
	config := DefaultConfig()
	return config.SaveConfig(configPath)
}
