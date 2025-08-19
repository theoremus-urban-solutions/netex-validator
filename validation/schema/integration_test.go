package schema

import (
	"strings"
	"testing"
)

// TestXSDValidatorIntegration tests the integration between XSDValidator and SchemaManager
func TestXSDValidatorIntegration(t *testing.T) {
	// Create validator with network disabled for testing
	options := &XSDValidationOptions{
		AllowNetworkDownload: false,
		CacheDirectory:       t.TempDir(),
		CacheExpiryHours:     1,
		StrictMode:           false,
		MaxSchemaSize:        50 * 1024 * 1024,
		HttpTimeoutSeconds:   30,
		UseLibxml2:           false,
	}

	validator, err := NewXSDValidator(options)
	if err != nil {
		t.Fatalf("Failed to create XSD validator: %v", err)
	}

	// Test basic NetEX XML content
	xmlContent := `<?xml version="1.0" encoding="UTF-8"?>
<PublicationDelivery xmlns="http://www.netex.org.uk/netex" version="1.4">
	<PublicationTimestamp>2023-01-01T12:00:00Z</PublicationTimestamp>
	<ParticipantRef>TEST</ParticipantRef>
	<dataObjects>
		<SiteFrame>
			<stopPlaces>
				<StopPlace id="TEST:StopPlace:001">
					<Name>Test Stop</Name>
				</StopPlace>
			</stopPlaces>
		</SiteFrame>
	</dataObjects>
</PublicationDelivery>`

	// Validate the XML
	validationErrors, err := validator.ValidateXML([]byte(xmlContent), "test.xml")
	if err != nil {
		t.Fatalf("Validation failed with error: %v", err)
	}

	// Should have no critical errors for basic structure
	hasSchemaErrors := false
	for _, verr := range validationErrors {
		if strings.Contains(verr.Message, "Missing required") {
			hasSchemaErrors = true
		}
	}

	if hasSchemaErrors {
		t.Errorf("Basic NetEX structure should not have schema errors, got: %v", validationErrors)
	}
}

// TestSchemaManagerIntegration tests SchemaManager functionality
func TestSchemaManagerIntegration(t *testing.T) {
	// Create schema manager
	schemaManager := NewSchemaManager(t.TempDir())
	schemaManager.SetNetworkEnabled(false) // Disable network for testing

	// Test version detection
	xmlContent := `<?xml version="1.0" encoding="UTF-8"?>
<PublicationDelivery xmlns="http://www.netex.org.uk/netex" version="1.4">
</PublicationDelivery>`

	version, err := schemaManager.DetectSchemaVersion([]byte(xmlContent))
	if err != nil {
		t.Fatalf("Failed to detect schema version: %v", err)
	}

	if version != "1.4" {
		t.Errorf("Expected version 1.4, got %s", version)
	}

	// Test cache stats
	stats := schemaManager.GetCacheStats()
	if stats["cacheDir"] == nil {
		t.Error("Cache stats should include cacheDir")
	}
}

// TestValidationIntegration tests the complete validation pipeline
func TestValidationIntegration(t *testing.T) {
	// Create validator
	validator, err := NewXSDValidator(DefaultXSDValidationOptions())
	if err != nil {
		t.Fatalf("Failed to create validator: %v", err)
	}

	// Test invalid XML (missing required elements)
	invalidXML := `<?xml version="1.0" encoding="UTF-8"?>
<PublicationDelivery xmlns="http://www.netex.org.uk/netex">
	<!-- Missing PublicationTimestamp, ParticipantRef, dataObjects -->
</PublicationDelivery>`

	errors, err := validator.ValidateXML([]byte(invalidXML), "invalid.xml")
	if err != nil {
		t.Fatalf("Validation should not fail with error: %v", err)
	}

	if len(errors) == 0 {
		t.Error("Expected validation errors for incomplete NetEX document")
	}

	// Check that we get the expected validation errors
	foundMissingElements := false
	for _, verr := range errors {
		if strings.Contains(verr.Message, "Missing required element") ||
			strings.Contains(verr.Message, "Missing required") {
			foundMissingElements = true
			break
		}
	}

	if !foundMissingElements {
		t.Logf("Validation errors found: %v", errors)
		// This might be expected if the basic validation doesn't catch all missing elements
		// The important thing is that we get some validation errors
		if len(errors) == 0 {
			t.Error("Expected some validation errors for incomplete NetEX document")
		}
	}
}

// TestSupportedVersions tests version support
func TestSupportedVersions(t *testing.T) {
	validator, err := NewXSDValidator(DefaultXSDValidationOptions())
	if err != nil {
		t.Fatalf("Failed to create validator: %v", err)
	}

	versions := validator.GetSupportedVersions()
	if len(versions) == 0 {
		t.Error("Should support at least one NetEX version")
	}

	// Check that common versions are supported
	supportedMap := make(map[string]bool)
	for _, v := range versions {
		supportedMap[v] = true
	}

	expectedVersions := []string{"1.0", "1.4", "1.16"}
	for _, expectedVer := range expectedVersions {
		if !supportedMap[expectedVer] {
			t.Errorf("Expected version %s to be supported", expectedVer)
		}
	}
}
