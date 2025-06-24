package fynetest

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// Suite manages a collection of tests with shared configuration.
type Suite struct {
	tests  []Test
	runner *Runner
	config SuiteConfig
}

// SuiteConfig contains configuration options for a test suite.
type SuiteConfig struct {
	// Name of the test suite
	Name string
	
	// OutputDir for screenshots (default: "test-screenshots")
	OutputDir string
	
	// DefaultTheme for all tests (can be overridden per test)
	DefaultTheme fyne.Theme
	
	// DefaultSize for test windows (can be overridden per test)
	DefaultSize fyne.Size
	
	// Parallel enables concurrent test execution
	Parallel bool
	
	// MaxConcurrency limits parallel execution (default: 4)
	MaxConcurrency int
	
	// Verbose enables detailed output
	Verbose bool
	
	// GenerateReport enables HTML report generation
	GenerateReport bool
	
	// ReportTitle for the HTML report
	ReportTitle string
}

// NewSuite creates a new test suite with default configuration.
func NewSuite() *Suite {
	return &Suite{
		tests:  make([]Test, 0),
		runner: NewRunner(),
		config: SuiteConfig{
			Name:           "Fyne Visual Tests",
			OutputDir:      "test-screenshots",
			DefaultTheme:   theme.LightTheme(),
			DefaultSize:    fyne.NewSize(800, 600),
			Parallel:       false,
			MaxConcurrency: 4,
			Verbose:        false,
			GenerateReport: true,
			ReportTitle:    "Fyne Visual Test Results",
		},
	}
}

// NewSuiteWithConfig creates a new test suite with custom configuration.
func NewSuiteWithConfig(config SuiteConfig) *Suite {
	suite := &Suite{
		tests:  make([]Test, 0),
		runner: NewRunner(),
		config: config,
	}
	
	// Apply config to runner
	suite.runner.OutputDir = config.OutputDir
	suite.runner.DefaultTheme = config.DefaultTheme
	suite.runner.DefaultSize = config.DefaultSize
	suite.runner.Verbose = config.Verbose
	
	return suite
}

// Add adds a single test to the suite.
func (s *Suite) Add(test Test) *Suite {
	s.tests = append(s.tests, test)
	return s
}

// AddTests adds multiple tests to the suite.
func (s *Suite) AddTests(tests ...Test) *Suite {
	s.tests = append(s.tests, tests...)
	return s
}

// AddBuilder adds a test using a builder.
func (s *Suite) AddBuilder(builder *TestBuilder) *Suite {
	test, err := builder.Build()
	if err != nil {
		panic(fmt.Sprintf("failed to build test: %v", err))
	}
	return s.Add(test)
}

// WithConfig updates the suite configuration.
func (s *Suite) WithConfig(fn func(*SuiteConfig)) *Suite {
	fn(&s.config)
	
	// Update runner with new config
	s.runner.OutputDir = s.config.OutputDir
	s.runner.DefaultTheme = s.config.DefaultTheme
	s.runner.DefaultSize = s.config.DefaultSize
	s.runner.Verbose = s.config.Verbose
	
	return s
}

// FilterByTags returns tests that have any of the specified tags.
func (s *Suite) FilterByTags(tags ...string) []Test {
	if len(tags) == 0 {
		return s.tests
	}
	
	filtered := make([]Test, 0)
	for _, test := range s.tests {
		for _, tag := range tags {
			if contains(test.Tags, tag) {
				filtered = append(filtered, test)
				break
			}
		}
	}
	return filtered
}

// FilterByName returns tests whose names contain the given substring.
func (s *Suite) FilterByName(pattern string) []Test {
	filtered := make([]Test, 0)
	pattern = strings.ToLower(pattern)
	
	for _, test := range s.tests {
		if strings.Contains(strings.ToLower(test.Name), pattern) {
			filtered = append(filtered, test)
		}
	}
	return filtered
}

// GetTestNames returns a sorted list of all test names.
func (s *Suite) GetTestNames() []string {
	names := make([]string, len(s.tests))
	for i, test := range s.tests {
		names[i] = test.Name
	}
	sort.Strings(names)
	return names
}

// Run executes all tests in the suite and returns the results.
func (s *Suite) Run() (SuiteResult, error) {
	return s.RunTests(s.tests)
}

// RunTests executes specific tests and returns the results.
func (s *Suite) RunTests(tests []Test) (SuiteResult, error) {
	startTime := time.Now()
	
	// Create timestamped output directory
	var results []Result
	var outputDir string
	
	if s.config.Parallel && len(tests) > 1 {
		results, outputDir = s.runner.RunTestsWithTimestamp(tests)
	} else {
		results, outputDir = s.runner.RunTestsWithTimestamp(tests)
	}
	
	// Create suite result
	suiteResult := SuiteResult{
		Name:      s.config.Name,
		Results:   results,
		StartTime: startTime,
		EndTime:   time.Now(),
		OutputDir: outputDir,
	}
	
	// Generate report if enabled
	if s.config.GenerateReport {
		reportPath := filepath.Join(outputDir, "index.html")
		reporter := NewReportGenerator()
		reporter.Title = s.config.ReportTitle
		
		if err := reporter.GenerateHTMLReport(results, reportPath); err != nil {
			return suiteResult, fmt.Errorf("failed to generate report: %w", err)
		}
		
		suiteResult.ReportPath = reportPath
	}
	
	return suiteResult, nil
}

// RunCLI runs the test suite as a CLI application with flag parsing.
// This is the main entry point for command-line usage.
func (s *Suite) RunCLI() {
	// Parse command line flags
	outputDir := flag.String("output", s.config.OutputDir, "Output directory for screenshots")
	testName := flag.String("test", "", "Run specific test by name")
	testPattern := flag.String("pattern", "", "Run tests matching name pattern")
	listTests := flag.Bool("list", false, "List all available tests")
	listTags := flag.Bool("tags", false, "List all available tags")
	tagFilter := flag.String("tag", "", "Run tests with specific tag")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	parallel := flag.Bool("parallel", s.config.Parallel, "Run tests in parallel")
	reportTitle := flag.String("title", s.config.ReportTitle, "Title for HTML report")
	noReport := flag.Bool("no-report", false, "Disable HTML report generation")
	flag.Parse()
	
	// Apply CLI flags to config
	s.config.OutputDir = *outputDir
	s.config.Verbose = *verbose
	s.config.Parallel = *parallel
	s.config.ReportTitle = *reportTitle
	s.config.GenerateReport = !*noReport
	
	// Update runner
	s.runner.OutputDir = s.config.OutputDir
	s.runner.Verbose = s.config.Verbose
	
	// Handle list flags
	if *listTests {
		s.listTests()
		return
	}
	
	if *listTags {
		s.listTags()
		return
	}
	
	// Filter tests based on flags
	testsToRun := s.tests
	
	if *testName != "" {
		testsToRun = s.filterByExactName(*testName)
		if len(testsToRun) == 0 {
			fmt.Printf("❌ Test '%s' not found\n", *testName)
			s.listTests()
			os.Exit(1)
		}
	} else if *testPattern != "" {
		testsToRun = s.FilterByName(*testPattern)
		if len(testsToRun) == 0 {
			fmt.Printf("❌ No tests match pattern '%s'\n", *testPattern)
			s.listTests()
			os.Exit(1)
		}
	} else if *tagFilter != "" {
		testsToRun = s.FilterByTags(*tagFilter)
		if len(testsToRun) == 0 {
			fmt.Printf("❌ No tests with tag '%s'\n", *tagFilter)
			s.listTags()
			os.Exit(1)
		}
	}
	
	// Print header
	fmt.Println("🧪 Fyne Visual Test Runner")
	fmt.Println("==========================")
	fmt.Printf("Suite: %s\n", s.config.Name)
	fmt.Printf("Output directory: %s\n", s.config.OutputDir)
	if s.config.Parallel {
		fmt.Printf("Execution mode: Parallel (max %d)\n", s.config.MaxConcurrency)
	} else {
		fmt.Println("Execution mode: Sequential")
	}
	fmt.Printf("Tests to run: %d\n", len(testsToRun))
	fmt.Println()
	
	// Run tests
	result, err := s.RunTests(testsToRun)
	if err != nil {
		fmt.Printf("❌ Error running tests: %v\n", err)
		os.Exit(1)
	}
	
	// Print summary
	s.printSummary(result)
	
	// Exit with error code if tests failed
	if result.Failed() > 0 {
		os.Exit(1)
	}
}

// Helper methods

func (s *Suite) filterByExactName(name string) []Test {
	for _, test := range s.tests {
		if test.Name == name {
			return []Test{test}
		}
	}
	return []Test{}
}

func (s *Suite) listTests() {
	fmt.Println("Available visual tests:")
	fmt.Println("======================")
	
	for i, test := range s.tests {
		fmt.Printf("%d. %s", i+1, test.Name)
		if test.Description != "" {
			fmt.Printf(" - %s", test.Description)
		}
		if len(test.Tags) > 0 {
			fmt.Printf(" [%s]", strings.Join(test.Tags, ", "))
		}
		fmt.Println()
	}
}

func (s *Suite) listTags() {
	tagMap := make(map[string]int)
	for _, test := range s.tests {
		for _, tag := range test.Tags {
			tagMap[tag]++
		}
	}
	
	if len(tagMap) == 0 {
		fmt.Println("No tags defined in test suite")
		return
	}
	
	fmt.Println("Available tags:")
	fmt.Println("===============")
	
	// Sort tags
	tags := make([]string, 0, len(tagMap))
	for tag := range tagMap {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	
	for _, tag := range tags {
		fmt.Printf("- %s (%d tests)\n", tag, tagMap[tag])
	}
}

func (s *Suite) printSummary(result SuiteResult) {
	fmt.Println("\n📊 Test Summary")
	fmt.Println("===============")
	fmt.Printf("Total tests: %d\n", result.Total())
	fmt.Printf("✅ Passed: %d\n", result.Passed())
	fmt.Printf("❌ Failed: %d\n", result.Failed())
	fmt.Printf("⏱️  Duration: %v\n", result.Duration())
	fmt.Printf("\nScreenshots saved to: %s\n", result.OutputDir)
	
	if result.ReportPath != "" {
		fmt.Printf("View results: file://%s\n", result.ReportPath)
	}
	
	// List failed tests
	if result.Failed() > 0 {
		fmt.Println("\nFailed tests:")
		for _, r := range result.Results {
			if !r.Success {
				fmt.Printf("- %s: %v\n", r.Test.Name, r.Error)
			}
		}
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// SuiteResult contains the results of running a test suite.
type SuiteResult struct {
	Name       string
	Results    []Result
	StartTime  time.Time
	EndTime    time.Time
	OutputDir  string
	ReportPath string
}

// Total returns the total number of tests run.
func (sr SuiteResult) Total() int {
	return len(sr.Results)
}

// Passed returns the number of tests that passed.
func (sr SuiteResult) Passed() int {
	count := 0
	for _, r := range sr.Results {
		if r.Success {
			count++
		}
	}
	return count
}

// Failed returns the number of tests that failed.
func (sr SuiteResult) Failed() int {
	return sr.Total() - sr.Passed()
}

// Duration returns how long the suite took to run.
func (sr SuiteResult) Duration() time.Duration {
	return sr.EndTime.Sub(sr.StartTime)
}

// PassRate returns the percentage of tests that passed.
func (sr SuiteResult) PassRate() float64 {
	if sr.Total() == 0 {
		return 0
	}
	return float64(sr.Passed()) / float64(sr.Total()) * 100
}