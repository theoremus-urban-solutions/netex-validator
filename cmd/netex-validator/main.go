package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/pprof"
	"strings"

	"github.com/spf13/cobra"
	netexvalidator "github.com/theoremus-urban-solutions/netex-validator/netexvalidator"
	"github.com/theoremus-urban-solutions/netex-validator/types"
)

var (
	inputFile       string
	outputFile      string
	outputFormat    string
	codespace       string
	skipSchema      bool
	skipValidators  bool
	verbose         bool
	maxSchemaErrors int
	configFile      string
	generateConfig  bool
	profile         string
	maxFindings     int
	allowSchemaNet  bool
	schemaCacheDir  string
	schemaTimeout   int
	useLibxml2XSD   bool
	concurrentFiles int
	cpuProfile      string
	memProfile      string
	// Performance optimization flags
	enableCache      bool
	cacheMaxEntries  int
	cacheMaxMemoryMB int
	cacheTTLHours    int
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "netex-validator",
		Short: "NetEX validator for Nordic NeTEx Profile",
		Long: `A comprehensive NetEX validator with extensive rule coverage that supports:
- XML Schema validation
- 88+ XPath-based business rules covering all major NetEX categories
- ZIP dataset validation with cross-file ID validation
- YAML configuration for rule customization
- JSON and HTML output formats

Examples:
  netex-validator -i data.xml -c "MyCodespace"
  netex-validator -i dataset.zip -c "MyCodespace" --format json
  netex-validator -i data.xml -c "MyCodespace" --config custom-rules.yaml`,
		RunE: validateCommand,
	}

	// Add flags
	rootCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input NetEX file or ZIP dataset (required)")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file (default: stdout)")
	rootCmd.Flags().StringVar(&outputFormat, "format", "", "Output format: json or html (default: json)")
	rootCmd.Flags().StringVarP(&codespace, "codespace", "c", "", "Validation codespace (required)")
	rootCmd.Flags().BoolVar(&skipSchema, "skip-schema", false, "Skip XML Schema validation")
	rootCmd.Flags().BoolVar(&skipValidators, "skip-validators", false, "Skip XPath business rule validation")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.Flags().IntVar(&maxSchemaErrors, "max-schema-errors", 0, "Maximum schema errors to report (0 = use config default)")
	rootCmd.Flags().StringVar(&configFile, "config", "", "Configuration file path")
	rootCmd.Flags().BoolVar(&generateConfig, "generate-config", false, "Generate default configuration file")
	// Profile flag retained for compatibility but ignored (EU is default)
	rootCmd.Flags().StringVar(&profile, "profile", "", "(Deprecated) Validation profile – ignored; EU is default")
	rootCmd.Flags().IntVar(&maxFindings, "max-findings", 0, "Maximum number of findings to report (0 = unlimited)")
	rootCmd.Flags().BoolVar(&allowSchemaNet, "allow-schema-network", true, "Allow downloading NetEX schemas from the network")
	rootCmd.Flags().StringVar(&schemaCacheDir, "schema-cache-dir", "", "Directory to cache downloaded schemas")
	rootCmd.Flags().IntVar(&schemaTimeout, "schema-timeout", 30, "Schema download timeout in seconds")
	rootCmd.Flags().BoolVar(&useLibxml2XSD, "use-libxml2-xsd", false, "Use libxml2-backed XSD validation (experimental)")
	rootCmd.Flags().IntVar(&concurrentFiles, "concurrent", 0, "Number of files to validate in parallel for ZIP datasets (0 = default)")
	rootCmd.Flags().StringVar(&cpuProfile, "cpuprofile", "", "Write CPU profile to file")
	rootCmd.Flags().StringVar(&memProfile, "memprofile", "", "Write memory profile to file")

	// Performance optimization flags
	rootCmd.Flags().BoolVar(&enableCache, "enable-cache", false, "Enable validation result caching by file hash")
	rootCmd.Flags().IntVar(&cacheMaxEntries, "cache-max-entries", 1000, "Maximum number of cached validation results")
	rootCmd.Flags().IntVar(&cacheMaxMemoryMB, "cache-max-memory-mb", 50, "Maximum memory usage for cache in MB")
	rootCmd.Flags().IntVar(&cacheTTLHours, "cache-ttl", 24, "Cache time-to-live in hours")

	// Mark required flags
	rootCmd.MarkFlagRequired("input")
	rootCmd.MarkFlagRequired("codespace")

	// Add generate-config command
	var generateConfigCmd = &cobra.Command{
		Use:   "generate-config [file]",
		Short: "Generate default configuration file",
		Long:  "Generate a default YAML configuration file for customizing validation rules",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath := "netex-validator.yaml"
			if len(args) > 0 {
				configPath = args[0]
			}
			return generateDefaultConfig(configPath)
		},
	}
	rootCmd.AddCommand(generateConfigCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func validateCommand(cmd *cobra.Command, args []string) error {
	// Handle generate-config flag
	if generateConfig {
		return generateDefaultConfig("netex-validator.yaml")
	}

	// Validate input file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return fmt.Errorf("input file does not exist: %s", inputFile)
	}

	// Start CPU profiling if requested
	if cpuProfile != "" {
		f, err := os.Create(cpuProfile)
		if err != nil {
			return fmt.Errorf("could not create CPU profile: %w", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			_ = f.Close()
			return fmt.Errorf("could not start CPU profile: %w", err)
		}
		defer pprof.StopCPUProfile()
		defer f.Close()
	}

	if verbose {
		fmt.Printf("NetEX Validator - Starting validation\n")
		fmt.Printf("Input: %s\n", inputFile)
		fmt.Printf("Codespace: %s\n", codespace)
		if configFile != "" {
			fmt.Printf("Config: %s\n", configFile)
		}
	}

	// Create validation options
	options := netexvalidator.DefaultValidationOptions().
		WithCodespace(codespace).
		WithSkipSchema(skipSchema).
		WithVerbose(verbose).
		WithConfigFile(configFile)
	if profile != "" {
		options = options.WithProfile(profile)
	}
	if maxFindings > 0 {
		options = options.WithMaxFindings(maxFindings)
	}
	options = options.WithAllowSchemaNetwork(allowSchemaNet)
	if schemaCacheDir != "" {
		options = options.WithSchemaCacheDir(schemaCacheDir)
	}
	if schemaTimeout > 0 {
		options = options.WithSchemaTimeoutSeconds(schemaTimeout)
	}
	if useLibxml2XSD {
		options = options.WithUseLibxml2XSD(true)
	}
	if concurrentFiles > 0 {
		options = options.WithConcurrentFiles(concurrentFiles)
	}

	// Performance optimization options
	if enableCache {
		options = options.WithValidationCache(enableCache, cacheMaxEntries, cacheMaxMemoryMB, cacheTTLHours)
	}

	if maxSchemaErrors > 0 {
		options.MaxSchemaErrors = maxSchemaErrors
	}

	// Determine output format
	format := "json"
	if outputFormat != "" {
		format = outputFormat
	}
	options.OutputFormat = format

	// Perform validation
	var result *netexvalidator.ValidationResult
	var err error

	isZip := strings.ToLower(filepath.Ext(inputFile)) == ".zip"
	if isZip {
		if verbose {
			fmt.Printf("Processing ZIP dataset...\n")
		}
		result, err = netexvalidator.ValidateZip(inputFile, options)
	} else {
		if verbose {
			fmt.Printf("Processing single XML file...\n")
		}
		result, err = netexvalidator.ValidateFile(inputFile, options)
	}

	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Write memory profile if requested
	if memProfile != "" {
		f, err := os.Create(memProfile)
		if err == nil {
			pprof.WriteHeapProfile(f)
			f.Close()
		}
	}

	if verbose {
		summary := result.Summary()
		fmt.Printf("Validation completed: %d issues found (%d files processed)\n",
			summary.TotalIssues, summary.FilesProcessed)

		if len(summary.IssuesBySeverity) > 0 {
			fmt.Printf("Issues by severity: ")
			for severity, count := range summary.IssuesBySeverity {
				fmt.Printf("%s:%d ", severityToString(severity), count)
			}
			fmt.Printf("\n")
		}
	}

	// Output results
	if err := outputResult(result, format); err != nil {
		return fmt.Errorf("failed to output results: %w", err)
	}

	// Exit with error code if validation found errors
	if !result.IsValid() {
		if verbose {
			fmt.Printf("Validation completed with errors\n")
		}
		os.Exit(1)
	}

	if verbose {
		fmt.Printf("Validation completed successfully\n")
	}

	return nil
}

func outputResult(result *netexvalidator.ValidationResult, format string) error {
	var output []byte
	var err error

	switch format {
	case "json":
		output, err = result.ToJSON()
	case "html":
		output, err = result.ToHTML()
	default:
		return fmt.Errorf("unsupported output format: %s (supported: json, html)", format)
	}

	if err != nil {
		return err
	}

	// Write to file or stdout
	if outputFile != "" {
		return os.WriteFile(outputFile, output, 0644)
	} else {
		fmt.Print(string(output))
		return nil
	}
}

func generateDefaultConfig(configPath string) error {
	// For now, just create a simple default config
	// This could be enhanced to use the actual config generation from the library
	defaultConfig := `# NetEX Validator Configuration
validator:
  profile: "eu"
  maxFileSize: 104857600  # 100MB
  maxSchemaErrors: 100
  concurrentFiles: 4
  enableCache: false
  cacheTimeout: 30

rules:
  categories:
    line:
      enabled: true
    route:
      enabled: true
    service_journey:
      enabled: true
    # Add other categories as needed

output:
  format: "json"
  includeDetails: true
  groupBySeverity: true
  maxEntries: 0
`

	err := os.WriteFile(configPath, []byte(defaultConfig), 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	fmt.Printf("Generated default configuration file: %s\n", configPath)
	return nil
}

func severityToString(severity types.Severity) string {
	switch severity {
	case types.INFO:
		return "INFO"
	case types.WARNING:
		return "WARNING"
	case types.ERROR:
		return "ERROR"
	case types.CRITICAL:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}
