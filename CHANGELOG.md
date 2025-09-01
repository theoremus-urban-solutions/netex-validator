# Changelog

All notable changes to the NetEX Validator project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.2] - 2025-09-01

### Added
- Grouped JSON output format for better user experience
- Consistent grouping between HTML and JSON reports
- JSON output now groups validation issues by:
  - File: Issues organized by source file
  - Severity: Issues organized by severity level (Critical/Error/Warning/Info)
  - Rule: Issues organized by validation rule name
- Enhanced JSON statistics with top issues and severity percentages
- Backwards compatibility with flat JSON format via ToFlatJSON() method

### Changed
- JSON output now uses grouped format by default
- Improved consistency between HTML and JSON report structures

### Added
- Initial release of NetEX Validator
- Support for NetEX XML validation against EU NeTEx Profile
- Over 200 built-in validation rules covering:
  - Schema validation
  - Business logic rules
  - Cross-file ID validation
  - Transport mode validation
  - Service journey validation
  - Stop point validation
  - Timing validation
  - Reference integrity checks
- Multiple validation modes:
  - Single XML file validation
  - ZIP dataset validation (multiple files)
  - In-memory content validation
- Configurable validation options:
  - Skip schema validation for faster processing
  - Custom codespace configuration
  - Rule severity overrides via YAML configuration
  - Selective rule enable/disable
- Multiple output formats:
  - JSON for programmatic processing
  - HTML with interactive interface
  - Plain text for CLI output
- Performance features:
  - Validation caching with TTL
  - Concurrent file processing
  - Memory-efficient processing for large datasets
- Comprehensive test coverage including:
  - Unit tests for all components
  - Integration tests
  - Performance benchmarks
  - Concurrent validation stress tests
- Go library API for integration into other applications
- Command-line interface with rich options

### Fixed
- Race conditions in concurrent validation scenarios
- XPath expression thread-safety issues
- Global logger concurrent access synchronization
- Validation cache result modification race conditions

### Security
- Safe XML parsing with protection against XXE attacks
- Secure file path handling
- Input validation and sanitization

## [0.1.0] - TBD

_Initial release - see Unreleased section above for features_

[Unreleased]: https://github.com/theoremus-urban-solutions/netex-validator/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/theoremus-urban-solutions/netex-validator/releases/tag/v0.1.0