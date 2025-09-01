package validator

import (
	"encoding/json"
	"sort"
	"time"

	"github.com/theoremus-urban-solutions/netex-validator/types"
)

// GroupedJSONResult represents the JSON output with grouping similar to HTML report
type GroupedJSONResult struct {
	// Metadata
	Codespace          string    `json:"codespace"`
	ValidationReportID string    `json:"validationReportId"`
	CreationDate       time.Time `json:"creationDate"`
	GeneratedAt        time.Time `json:"generatedAt"`

	// Summary and statistics (consistent with HTML)
	Summary    ValidationSummary             `json:"summary"`
	Statistics map[string]interface{}        `json:"statistics"`

	// Original flat list (for backwards compatibility)
	ValidationReportEntries []ValidationReportEntry `json:"validationReportEntries"`

	// Grouped views (matching HTML report tabs)
	GroupedByFile     map[string][]ValidationReportEntry `json:"groupedByFile"`
	GroupedBySeverity map[string][]ValidationReportEntry `json:"groupedBySeverity"`
	GroupedByRule     map[string][]ValidationReportEntry `json:"groupedByRule"`

	// Processing info
	FilesProcessed int           `json:"filesProcessed"`
	ProcessingTime time.Duration `json:"processingTimeMs"`
	CacheHit       bool          `json:"cacheHit,omitempty"`
	FileHash       string        `json:"fileHash,omitempty"`
}

// ToGroupedJSON converts the validation result to grouped JSON format (consistent with HTML)
func (r *ValidationResult) ToGroupedJSON() ([]byte, error) {
	grouped := r.createGroupedJSON()
	return json.MarshalIndent(grouped, "", "  ")
}

// createGroupedJSON creates a grouped JSON structure consistent with HTML report
func (r *ValidationResult) createGroupedJSON() *GroupedJSONResult {
	// Group by file (same as HTML)
	groupedByFile := make(map[string][]ValidationReportEntry)
	for _, entry := range r.ValidationReportEntries {
		fileName := entry.FileName
		if fileName == "" {
			fileName = "Unknown"
		}
		groupedByFile[fileName] = append(groupedByFile[fileName], entry)
	}

	// Group by severity (same as HTML)
	groupedBySeverity := make(map[string][]ValidationReportEntry)
	for _, entry := range r.ValidationReportEntries {
		severity := severityTextForJSON(entry.Severity)
		groupedBySeverity[severity] = append(groupedBySeverity[severity], entry)
	}

	// Group by rule (same as HTML)
	groupedByRule := make(map[string][]ValidationReportEntry)
	for _, entry := range r.ValidationReportEntries {
		groupedByRule[entry.Name] = append(groupedByRule[entry.Name], entry)
	}

	// Calculate statistics (same as HTML)
	statistics := r.calculateStatistics(groupedBySeverity)

	// Get summary
	summary := r.Summary()

	return &GroupedJSONResult{
		Codespace:               r.Codespace,
		ValidationReportID:      r.ValidationReportID,
		CreationDate:            r.CreationDate,
		GeneratedAt:             time.Now(),
		Summary:                 summary,
		Statistics:              statistics,
		ValidationReportEntries: r.ValidationReportEntries,
		GroupedByFile:           groupedByFile,
		GroupedBySeverity:       groupedBySeverity,
		GroupedByRule:           groupedByRule,
		FilesProcessed:          r.FilesProcessed,
		ProcessingTime:          r.ProcessingTime,
		CacheHit:                r.CacheHit,
		FileHash:                r.FileHash,
	}
}

// calculateStatistics calculates statistics consistent with HTML report
func (r *ValidationResult) calculateStatistics(issuesBySeverity map[string][]ValidationReportEntry) map[string]interface{} {
	stats := make(map[string]interface{})
	
	// Total issues
	totalIssues := len(r.ValidationReportEntries)
	stats["totalIssues"] = totalIssues
	stats["filesProcessed"] = r.FilesProcessed
	stats["processingTimeMs"] = r.ProcessingTime.Milliseconds()
	stats["hasErrors"] = !r.IsValid()
	
	// Severity counts
	severityCounts := make(map[string]int)
	severityPercents := make(map[string]float64)
	
	for severity, issues := range issuesBySeverity {
		count := len(issues)
		severityCounts[severity] = count
		if totalIssues > 0 {
			severityPercents[severity] = float64(count) / float64(totalIssues) * 100
		}
	}
	
	stats["severityCounts"] = severityCounts
	stats["severityPercents"] = severityPercents
	
	// Files with issues
	filesWithIssues := make(map[string]bool)
	for _, entry := range r.ValidationReportEntries {
		if entry.FileName != "" {
			filesWithIssues[entry.FileName] = true
		}
	}
	stats["filesWithIssues"] = len(filesWithIssues)
	
	// Most common issues (top 5 rules)
	ruleCounts := make(map[string]int)
	for _, entry := range r.ValidationReportEntries {
		ruleCounts[entry.Name]++
	}
	
	type ruleCount struct {
		Rule  string `json:"rule"`
		Count int    `json:"count"`
	}
	
	var topIssues []ruleCount
	for rule, count := range ruleCounts {
		topIssues = append(topIssues, ruleCount{Rule: rule, Count: count})
	}
	
	// Sort by count descending
	sort.Slice(topIssues, func(i, j int) bool {
		return topIssues[i].Count > topIssues[j].Count
	})
	
	// Keep top 5
	if len(topIssues) > 5 {
		topIssues = topIssues[:5]
	}
	
	stats["topIssues"] = topIssues
	
	return stats
}

// ToJSONWithGrouping allows choosing specific grouping for JSON output
func (r *ValidationResult) ToJSONWithGrouping(groupBy string) ([]byte, error) {
	switch groupBy {
	case "file":
		return json.MarshalIndent(r.GetIssuesByFile(), "", "  ")
	case "severity":
		return json.MarshalIndent(r.GetIssuesBySeverity(), "", "  ")
	case "rule":
		groupedByRule := make(map[string][]ValidationReportEntry)
		for _, entry := range r.ValidationReportEntries {
			groupedByRule[entry.Name] = append(groupedByRule[entry.Name], entry)
		}
		return json.MarshalIndent(groupedByRule, "", "  ")
	case "full", "all":
		return r.ToGroupedJSON()
	default:
		// Default to original flat format for backwards compatibility
		return r.ToFlatJSON()
	}
}

// severityTextForJSON converts severity to text (same as HTML reporter)
func severityTextForJSON(severity types.Severity) string {
	switch severity {
	case types.CRITICAL:
		return "Critical"
	case types.ERROR:
		return "Error"
	case types.WARNING:
		return "Warning"
	case types.INFO:
		return "Info"
	default:
		return "Unknown"
	}
}