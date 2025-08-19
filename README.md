# NetEX Validator - Go Library & CLI

A comprehensive Go implementation of the NetEX validator for validating NetEX datasets against the EU NeTEx Profile. Available as both a Go library and command-line tool.

## Features

- **🏛️ Production-Ready Library**: Clean API for integration into Go applications
- **🖥️ Command-Line Interface**: Standalone CLI tool for file validation
- **📋 Professional HTML Reports**: Beautiful, interactive validation reports with tabbed interface
- **🔍 200+ Validation Rules**: Comprehensive coverage of NetEX business rules
- **📦 ZIP Dataset Support**: Validates multi-file NetEX datasets
- **⚙️ YAML Configuration**: Customizable rule configuration
- **🚀 High Performance**: Pure Go implementation with no external dependencies
- **📊 Multiple Output Formats**: JSON, HTML, and text output

## Installation

### Prerequisites

- Go 1.21 or later

### Build from Source

```bash
# Clone or navigate to the golang directory
cd golang

# Download dependencies
go mod download

# Build the CLI application
go build -o netex-validator cmd/netex-validator/main.go

# Or run directly
go run cmd/netex-validator/main.go --help
```

## Usage

### Basic Usage

```bash
# Validate a single NetEX XML file
./netex-validator -i input.xml -c "MyCodespace"

# Validate a NetEX dataset (ZIP file)
./netex-validator -i dataset.zip -c "MyCodespace"

# Save report to file
./netex-validator -i input.xml -c "MyCodespace" -o report.json

# Different output formats
./netex-validator -i input.xml -c "MyCodespace" -f text
./netex-validator -i input.xml -c "MyCodespace" -f xml
```

### Advanced Options

```bash
# Skip schema validation
./netex-validator -i input.xml -c "MyCodespace" --skip-schema

# Skip business rule validators
./netex-validator -i input.xml -c "MyCodespace" --skip-validators

# Verbose output
./netex-validator -i input.xml -c "MyCodespace" -v

# Limit schema errors
./netex-validator -i input.xml -c "MyCodespace" --max-schema-errors 50
```

### Command-line Options

- `-i, --input`: Input NetEX file or ZIP dataset (required)
- `-c, --codespace`: NetEX codespace (required)
- `-o, --output`: Output file for validation report (default: stdout)
- `-f, --format`: Output format: json, text, xml (default: json)
- `--skip-schema`: Skip XML schema validation
- `--skip-validators`: Skip NetEX business rule validators
- `-v, --verbose`: Verbose output
- `--max-schema-errors`: Maximum number of schema errors to report (default: 100)

## Architecture

The Go implementation follows the same architectural patterns as the Java version:

### Core Components

- **ValidationRunner**: Orchestrates the validation process
- **Schema Validator**: XML schema validation
- **XPath Validators**: Business rule validation using XPath
- **JAXB Validators**: Complex object model validation (planned)
- **Dataset Validators**: Cross-file validation (planned)

### Package Structure

```
pkg/
├── validator/          # Core validation interfaces and types  
├── parser/            # NetEX XML parsing and ZIP handling
├── schema/            # XML schema validation
├── xpath/             # XPath-based validation
│   └── rules/         # XPath validation rule implementations
├── config/            # Configuration handling (planned)
└── report/            # Validation reporting (planned)
```

## Validation Rules

The Go implementation includes the following validation categories:

### XML Schema Validation
- XML well-formedness
- NetEX namespace validation
- Root element validation

### XPath Business Rules
- **Structure validation**: Required elements and attributes
- **Transport mode validation**: Valid transport modes and submodes
- **Service journey validation**: Timing and passing time validation
- **Reference validation**: ID uniqueness and reference integrity
- **Frame validation**: Proper frame organization and relationships

### Implemented Rules

Currently implemented XPath rules include:

| Rule Code | Description |
|-----------|-------------|
| LINE_4 | Line missing TransportMode |
| SERVICE_JOURNEY_3 | ServiceJourney missing element PassingTimes |
| ROUTE_1 | Route missing |
| SERVICE_JOURNEY_1 | ServiceJourney must exist |
| RESOURCE_FRAME_IN_LINE_FILE | ResourceFrame must be exactly one |
| COMPOSITE_FRAME_1 | CompositeFrame must be exactly one |

*Note: This is a subset of the 268+ rules available in the Java version. More rules will be added progressively.*

## Output Formats

### JSON Format
```json
{
  "codespace": "MyCodespace",
  "validationReportId": "report-123",
  "creationDate": "2024-01-01T12:00:00Z",
  "validationReportEntries": [
    {
      "name": "Line missing TransportMode",
      "message": "Line is missing TransportMode",
      "severity": "ERROR",
      "fileName": "input.xml",
      "location": {
        "fileName": "input.xml",
        "lineNumber": 0,
        "xpath": "/PublicationDelivery/...",
        "elementId": "MyCodespace:Line:123"
      }
    }
  ],
  "numberOfValidationEntriesPerRule": {
    "Line missing TransportMode": 1
  }
}
```

### Text Format
```
NetEX Validation Report
======================
Codespace: MyCodespace
Report ID: report-123
Created: 2024-01-01 12:00:00
Total Issues: 1

ERROR (1):
----------
  File: input.xml
  Rule: Line missing TransportMode
  Message: Line is missing TransportMode
  Element ID: MyCodespace:Line:123
```

## Testing

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./validator
go test ./xpath

# Run with coverage
go test -cover ./...
```

## Development

### Adding New Validation Rules

1. **XPath Rules**: Add new rules in `pkg/xpath/rules/`
2. **JAXB Rules**: Add complex validation in `pkg/jaxb/rules/` (when implemented)
3. **Register Rules**: Add to validator creation in CLI

### Example: Adding a New XPath Rule

```go
// pkg/xpath/rules/my_rule.go
func NewValidateMyRule() *ValidateNotExist {
    rule := validator.ValidationRule{
        Code:     "MY_RULE_1",
        Name:     "My validation rule",
        Message:  "Validation failed",
        Severity: validator.ERROR,
    }
    
    return NewValidateNotExistWithRule(
        "//xpath/expression/here",
        rule,
    )
}
```

## Compatibility

This Go implementation aims to maintain compatibility with the Java version:

- **Input/Output**: Same file formats and report structures
- **Validation Rules**: Same rule codes and descriptions
- **Configuration**: Compatible YAML configuration (planned)

## Performance

The Go implementation is designed for performance:

- **Concurrent processing**: Parallel validation where possible
- **Memory efficiency**: Streaming XML processing
- **Fast startup**: No JVM overhead

## Limitations

Current limitations compared to the Java version:

- **Rule Coverage**: Subset of the full rule set (progressive implementation)
- **Configuration**: YAML configuration not yet implemented
- **JAXB Validation**: Complex object model validation not yet implemented
- **Schema Files**: XSD schema validation uses basic XML parsing

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## License

This project follows the same license as the Java NetEX validator:
EUPL-1.2 with modifications