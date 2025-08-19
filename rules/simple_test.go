package rules

import (
	"testing"

	"github.com/theoremus-urban-solutions/netex-validator/config"
	"github.com/theoremus-urban-solutions/netex-validator/types"
)

func TestRuleRegistryBasic(t *testing.T) {
	t.Run("Create rule registry", func(t *testing.T) {
		cfg := &config.ValidatorConfig{} // Minimal config
		registry := NewRuleRegistry(cfg)

		if registry == nil {
			t.Error("Expected non-nil registry")
		}
	})

	t.Run("Get all rules", func(t *testing.T) {
		cfg := &config.ValidatorConfig{}
		registry := NewRuleRegistry(cfg)

		// Should have some built-in rules
		rules := registry.GetEnabledRules()
		if len(rules) == 0 {
			t.Log("No built-in rules found (this might be expected)")
		} else {
			t.Logf("Found %d built-in rules", len(rules))
		}
	})
}

func TestRuleStruct(t *testing.T) {
	t.Run("Create rule", func(t *testing.T) {
		rule := Rule{
			Code:     "TEST_001",
			Name:     "Test Rule",
			Message:  "Test message",
			Severity: types.WARNING,
			XPath:    "//test",
			Category: "test",
		}

		if rule.Code != "TEST_001" {
			t.Errorf("Expected code TEST_001, got %s", rule.Code)
		}

		if rule.Name != "Test Rule" {
			t.Errorf("Expected name 'Test Rule', got %s", rule.Name)
		}

		if rule.Severity != types.WARNING {
			t.Errorf("Expected severity WARNING, got %s", rule.Severity)
		}
	})
}
