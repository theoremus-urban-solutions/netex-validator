# NetEX Validator - Go Library & CLI

[![Go Version](https://img.shields.io/badge/Go-1.21%2B-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-EUPL--1.2-blue.svg)](LICENSE)
[![CI/CD](https://github.com/theoremus-urban-solutions/netex-validator/workflows/CI%2FCD/badge.svg)](https://github.com/theoremus-urban-solutions/netex-validator/actions)

A comprehensive Go implementation of the NetEX validator for validating NetEX datasets against the EU NeTEx Profile. Available as both a Go library and command-line tool with advanced features including concurrent processing, validation caching, and comprehensive reporting.

## ✨ Features

### 🏛️ Production-Ready
- **Thread-Safe Validation**: Concurrent validation with race condition protection
- **Validation Caching**: TTL-based caching for improved performance
- **Memory Efficient**: Optimized for large dataset processing
- **Comprehensive Logging**: Structured logging with configurable levels

### 🔍 Extensive Validation Coverage
- **240+ Validation Rules**: Comprehensive NetEX business rules
- **Schema Validation**: XML Schema (XSD) validation against NetEX specification
- **Cross-File Validation**: ID reference validation across multiple files
- **Transport Mode Validation**: Comprehensive transport mode and submode validation
- **Service Journey Validation**: Timing, passing times, and journey pattern validation
- **Reference Integrity**: ID uniqueness and reference consistency checks

### 📦 Multiple Input Formats
- **Single XML Files**: Individual NetEX XML file validation
- **ZIP Datasets**: Multi-file NetEX dataset validation
- **In-Memory Content**: Direct validation of XML content
- **Stream Processing**: Memory-efficient processing of large files

### 📊 Rich Output Options
- **JSON Reports**: Machine-readable validation results
- **HTML Reports**: Interactive reports with tabbed interface and statistics
- **Plain Text**: Human-readable console output
- **Detailed Diagnostics**: Location information with XPath and line numbers

### ⚙️ Flexible Configuration
- **YAML Configuration**: Rule severity overrides and custom settings
- **Selective Validation**: Skip schema or business rule validation
- **Performance Tuning**: Configurable concurrency and cache settings
- **Custom Codespaces**: Support for organization-specific codespaces

## 🚀 Installation

### Prerequisites

- Go 1.21 or later

### Install via Go

```bash
go install github.com/theoremus-urban-solutions/netex-validator/cmd/netex-validator@latest
```

### Build from Source

```bash
git clone https://github.com/theoremus-urban-solutions/netex-validator.git
cd netex-validator
go build -o netex-validator cmd/netex-validator/main.go
```

### Download Pre-built Binaries

Download the latest release from the [GitHub Releases page](https://github.com/theoremus-urban-solutions/netex-validator/releases).

## 📖 Usage

### Command Line Interface

#### Basic Validation
```bash
# Validate a single NetEX XML file
./netex-validator validate -i input.xml -c "MyCodespace"

# Validate a NetEX dataset (ZIP file)
./netex-validator validate -i dataset.zip -c "MyCodespace"

# Generate HTML report
./netex-validator validate -i input.xml -c "MyCodespace" --html-output report.html
```

#### Advanced Options
```bash
# Skip schema validation for faster processing
./netex-validator validate -i input.xml -c "MyCodespace" --skip-schema

# Enable validation caching with custom TTL
./netex-validator validate -i input.xml -c "MyCodespace" --cache --cache-ttl 24h

# Limit concurrent processing
./netex-validator validate -i dataset.zip -c "MyCodespace" --concurrent-files 2

# Verbose output with debug information
./netex-validator validate -i input.xml -c "MyCodespace" --verbose

# Custom configuration file
./netex-validator validate -i input.xml -c "MyCodespace" --config config.yaml
```

#### Configuration File Example

```yaml
# config.yaml
rules:
  severityOverrides:
    TRANSPORT_MODE_1: WARNING
    SERVICE_JOURNEY_3: INFO
  disabled:
    - OPTIONAL_RULE_1
    - DEPRECATED_RULE_2

performance:
  maxSchemaErrors: 50
  concurrentFiles: 4
  
cache:
  enabled: true
  maxEntries: 1000
  maxMemoryMB: 100
  ttlHours: 24
```

### Go Library API

#### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/theoremus-urban-solutions/netex-validator/validator"
)

func main() {
    // Create validator with default options
    v, err := validator.New()
    if err != nil {
        log.Fatal(err)
    }
    
    // Validate a file
    result, err := v.ValidateFile("input.xml")
    if err != nil {
        log.Fatal(err)
    }
    
    // Check results
    if result.IsValid() {
        fmt.Println("✅ Validation passed!")
    } else {
        fmt.Printf("❌ Found %d issues\\n", len(result.ValidationReportEntries))
        for _, entry := range result.ValidationReportEntries {
            fmt.Printf("- %s: %s\\n", entry.Severity, entry.Message)
        }
    }
}
```

#### Advanced Configuration

```go
package main

import (
    "github.com/theoremus-urban-solutions/netex-validator/validator"
)

func main() {
    // Create validator with custom options
    options := validator.DefaultValidationOptions().
        WithCodespace("MyOrg").
        WithSkipSchema(false).
        WithValidationCache(true, 1000, 100, 24).
        WithConcurrentFiles(4).
        WithMaxFindings(500).
        WithVerbose(true)
    
    v, err := validator.NewWithOptions(options)
    if err != nil {
        log.Fatal(err)
    }
    
    // Validate ZIP dataset
    result, err := v.ValidateZip("dataset.zip")
    if err != nil {
        log.Fatal(err)
    }
    
    // Generate HTML report
    htmlReport, err := result.ToHTML()
    if err != nil {
        log.Fatal(err)
    }
    
    // Save report
    if err := os.WriteFile("report.html", []byte(htmlReport), 0644); err != nil {
        log.Fatal(err)
    }
}
```

#### Concurrent Validation

```go
package main

import (
    "sync"
    "github.com/theoremus-urban-solutions/netex-validator/validator"
)

func main() {
    v, _ := validator.New()
    
    files := []string{"file1.xml", "file2.xml", "file3.xml"}
    results := make([]*validator.ValidationResult, len(files))
    
    var wg sync.WaitGroup
    for i, file := range files {
        wg.Add(1)
        go func(idx int, filename string) {
            defer wg.Done()
            result, err := v.ValidateFile(filename)
            if err == nil {
                results[idx] = result
            }
        }(i, file)
    }
    wg.Wait()
    
    // Process results...
}
```

## 🏗️ Architecture

The validator follows a modular architecture with clear separation of concerns:

```
netex-validator/
├── cmd/netex-validator/        # CLI application
├── validator/                  # Main validator library
├── validation/
│   ├── engine/                # Validation orchestration
│   ├── schema/                # XSD schema validation
│   ├── context/               # Validation context management
│   └── ids/                   # ID and reference validation
├── rules/                     # Business rule definitions
├── reporting/                 # Error reporting and formatting
├── config/                    # Configuration management
├── logging/                   # Structured logging
├── utils/                     # Utilities (caching, XPath, etc.)
├── types/                     # Common type definitions
└── interfaces/                # Interface definitions
```

### Core Components

- **ValidationEngine**: Orchestrates the validation process with concurrent execution
- **SchemaValidator**: XML schema validation against NetEX XSD files
- **XPathValidator**: Business rule validation using XPath expressions
- **IDValidator**: Cross-file ID reference and uniqueness validation
- **ReportGenerator**: Creates detailed validation reports in multiple formats
- **CacheManager**: TTL-based validation result caching
- **ConfigManager**: YAML-based configuration handling

## 📋 Validation Rules

The validator implements 240+ comprehensive validation rules covering:

### Schema Validation
- XML well-formedness and namespace validation
- NetEX XSD schema compliance
- Root element and structure validation

### Business Logic Rules
- **Transport Rules**: Mode and submode validation
- **Service Rules**: Journey patterns, passing times, timetables
- **Infrastructure Rules**: Stop places, quays, service links
- **Reference Rules**: ID formats, reference integrity
- **Frame Rules**: Proper frame organization and relationships
- **Timing Rules**: Validity periods, operating days
- **Group Rules**: Lines, services, and operator groups

### Cross-File Validation
- ID uniqueness across multiple files
- Reference consistency in ZIP datasets
- Common data file validation
- Circular reference detection

## 📊 Output Examples

### JSON Output
```json
{
  "validationReport": {
    "creationDate": "2024-01-15T14:30:00Z",
    "codespace": "MyCodespace",
    "processingTime": "1.23s",
    "cacheHit": false,
    "summary": {
      "totalIssues": 3,
      "errors": 2,
      "warnings": 1,
      "infos": 0
    },
    "validationReportEntries": [
      {
        "severity": "ERROR",
        "name": "TRANSPORT_MODE_1",
        "message": "Line is missing required TransportMode",
        "fileName": "input.xml",
        "location": {
          "lineNumber": 42,
          "xpath": "//Line[@id='MyCodespace:Line:123']",
          "elementId": "MyCodespace:Line:123"
        }
      }
    ],
    "statistics": {
      "filesProcessed": 1,
      "rulesExecuted": 240,
      "processingTimeMs": 1230
    }
  }
}
```

### HTML Report Features
- **Interactive Interface**: Tabbed navigation between issues, statistics, and files
- **Filtering**: Filter by severity, rule, or file
- **Statistics Dashboard**: Visual charts and metrics
- **Responsive Design**: Works on desktop and mobile devices
- **Export Options**: Print-friendly and shareable reports

## 🧪 Testing

### Run Tests
```bash
# Run all tests
go test ./...

# Run with race detection
go test -race ./...

# Run with coverage
go test -cover ./...

# Run specific test suites
go test ./validator
go test ./validation/engine
```

### Benchmarks
```bash
# Run performance benchmarks
go test -bench=. ./validator

# Memory profiling
go test -memprofile=mem.prof -bench=. ./validator
```

### Test Coverage
The project maintains high test coverage with:
- Unit tests for all components
- Integration tests for end-to-end validation
- Performance benchmarks and stress tests
- Race condition detection tests
- Cross-platform compatibility tests

## 🔧 Development

### Adding New Validation Rules

1. **Define Rule**: Create rule definition in `rules/`
```go
var TransportModeRequired = ValidationRule{
    Code:     "TRANSPORT_MODE_1",
    Name:     "Line TransportMode Required",
    Message:  "Line is missing required TransportMode",
    Severity: ERROR,
    XPath:    "//Line[not(TransportMode)]",
}
```

2. **Register Rule**: Add to rule registry
3. **Add Tests**: Create comprehensive test cases
4. **Update Documentation**: Add to rule documentation

### Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes with tests
4. Ensure all tests pass (`go test -race ./...`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

### Code Style
- Follow Go conventions and `gofmt` formatting
- Use meaningful variable and function names
- Add comprehensive tests for new functionality
- Include documentation for public APIs
- Ensure thread-safety for concurrent operations

## 📈 Performance

The Go implementation is optimized for performance:

- **Concurrent Processing**: Parallel validation of multiple files
- **Memory Efficiency**: Stream processing and memory pooling
- **Validation Caching**: TTL-based result caching with configurable limits
- **Fast Startup**: No JVM overhead, instant startup time
- **Optimized XPath**: Compiled XPath expressions with thread-safe access

### Benchmarks
- **Single File**: ~100ms for typical NetEX files
- **Large Datasets**: Linear scaling with concurrent processing
- **Memory Usage**: <50MB for most validation scenarios
- **Cache Hit Rate**: 80%+ cache hit rate in typical workflows

## 🔒 Security

- **Safe XML Processing**: Protection against XXE and billion laughs attacks
- **Path Validation**: Secure file path handling for ZIP datasets
- **Input Sanitization**: Validation of all user inputs
- **Memory Limits**: Configurable limits to prevent DoS attacks

## 📄 License

This project is licensed under the EUPL-1.2 License - see the [LICENSE](LICENSE) file for details.

## 🤝 Support & Community

- **Issues**: Report bugs and request features on [GitHub Issues](https://github.com/theoremus-urban-solutions/netex-validator/issues)
- **Discussions**: Join the conversation on [GitHub Discussions](https://github.com/theoremus-urban-solutions/netex-validator/discussions)
- **Documentation**: Additional documentation available in the [docs](docs/) directory

## 📝 Changelog

See [CHANGELOG.md](CHANGELOG.md) for a detailed list of changes and version history.

---

**Made with ❤️ by the Theoremus Urban Solutions team**