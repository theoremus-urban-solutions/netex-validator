package main

import (
	"fmt"
	"log"

	netexvalidator "github.com/theoremus-urban-solutions/netex-validator/netexvalidator"
)

func main() {
	// Example 1: Simple file validation with default options
	fmt.Println("=== Example 1: Simple File Validation ===")

	options := netexvalidator.DefaultValidationOptions().
		WithCodespace("MyCodespace")

	result, err := netexvalidator.ValidateFile("germany.xml", options)
	if err != nil {
		log.Fatalf("Validation failed: %v", err)
	}

	summary := result.Summary()
	fmt.Printf("Validation completed: %d issues found\n", summary.TotalIssues)
	fmt.Printf("Processing time: %v\n", summary.ProcessingTime)
	fmt.Printf("Has errors: %v\n", summary.HasErrors)

	// Example 2: Custom validation with rule overrides
	fmt.Println("\n=== Example 2: Custom Validation Options ===")

	customOptions := netexvalidator.DefaultValidationOptions().
		WithCodespace("CustomCodespace").
		WithVerbose(true).
		WithSkipSchema(true) // Skip schema validation for faster processing

	customResult, err := netexvalidator.ValidateFile("germany.xml", customOptions)
	if err != nil {
		log.Fatalf("Custom validation failed: %v", err)
	}

	fmt.Printf("Custom validation: %d issues found\n", len(customResult.ValidationReportEntries))

	// Example 3: Using validator instance for multiple validations
	fmt.Println("\n=== Example 3: Validator Instance ===")

	validator, err := netexvalidator.NewWithOptions(options)
	if err != nil {
		log.Fatalf("Failed to create validator: %v", err)
	}

	// Validate the file using the validator instance
	instanceResult, err := validator.ValidateFile("germany.xml")
	if err != nil {
		log.Printf("Failed to validate with instance: %v", err)
	} else {
		fmt.Printf("Instance validation: %d issues found\n", len(instanceResult.ValidationReportEntries))
	}

	// Example 4: Detailed issue analysis
	fmt.Println("\n=== Example 4: Detailed Issue Analysis ===")

	detailedResult, err := netexvalidator.ValidateFile("germany.xml", options)
	if err != nil {
		log.Fatalf("Detailed validation failed: %v", err)
	}

	// Group by severity
	issuesBySeverity := detailedResult.GetIssuesBySeverity()
	for severity, issues := range issuesBySeverity {
		fmt.Printf("Severity %v: %d issues\n", severity, len(issues))
		for i, issue := range issues {
			if i < 2 { // Show first 2 issues per severity
				fmt.Printf("  - %s: %s\n", issue.Name, issue.Message)
			}
		}
	}

	// Group by file
	issuesByFile := detailedResult.GetIssuesByFile()
	fmt.Printf("\nIssues by file:\n")
	for fileName, issues := range issuesByFile {
		fmt.Printf("File %s: %d issues\n", fileName, len(issues))
	}

	// Example 5: JSON and HTML output
	fmt.Println("\n=== Example 5: JSON and HTML Output ===")

	jsonData, err := detailedResult.ToJSON()
	if err != nil {
		log.Fatalf("Failed to convert to JSON: %v", err)
	}

	fmt.Printf("JSON output length: %d bytes\n", len(jsonData))
	fmt.Printf("First 200 characters: %s...\n", string(jsonData[:200]))

	// HTML output
	htmlData, err := detailedResult.ToHTML()
	if err != nil {
		log.Fatalf("Failed to convert to HTML: %v", err)
	}

	fmt.Printf("HTML output length: %d bytes\n", len(htmlData))
	fmt.Printf("HTML starts with: %s...\n", string(htmlData[:50]))

	fmt.Println("\n=== Library Usage Examples Complete ===")
}
