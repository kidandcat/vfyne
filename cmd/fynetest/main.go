package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"plugin"

	fynetest "github.com/jairo/vfyne"
)

func main() {
	// Parse command line flags
	outputDir := flag.String("output", "test-screenshots", "Output directory for screenshots")
	testName := flag.String("test", "", "Run specific test by name")
	listTests := flag.Bool("list", false, "List all available tests")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	reportTitle := flag.String("title", "Fyne Visual Test Results", "Title for HTML report")
	pluginPath := flag.String("plugin", "", "Path to test plugin (.so file)")
	flag.Parse()

	if *pluginPath == "" {
		fmt.Fprintln(os.Stderr, "Error: -plugin flag is required")
		fmt.Fprintln(os.Stderr, "Usage: fynetest -plugin <path-to-test-plugin>")
		flag.Usage()
		os.Exit(1)
	}

	// Load the plugin
	p, err := plugin.Open(*pluginPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading plugin: %v\n", err)
		os.Exit(1)
	}

	// Look for the GetTests function
	getTestsSymbol, err := p.Lookup("GetTests")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: plugin must export 'GetTests' function: %v\n", err)
		os.Exit(1)
	}

	getTests, ok := getTestsSymbol.(func() []fynetest.Test)
	if !ok {
		fmt.Fprintln(os.Stderr, "Error: GetTests must have signature 'func() []fynetest.Test'")
		os.Exit(1)
	}

	// Get all tests from the plugin
	allTests := getTests()

	// Handle list flag
	if *listTests {
		fmt.Println("Available visual tests:")
		fmt.Println("======================")
		for i, test := range allTests {
			fmt.Printf("%d. %s - %s\n", i+1, test.Name, test.Description)
		}
		return
	}

	// Filter tests if specific test requested
	testsToRun := allTests
	if *testName != "" {
		testsToRun = []fynetest.Test{}
		for _, test := range allTests {
			if test.Name == *testName {
				testsToRun = append(testsToRun, test)
				break
			}
		}
		if len(testsToRun) == 0 {
			fmt.Printf("❌ Test '%s' not found\n", *testName)
			os.Exit(1)
		}
	}

	// Create runner
	runner := fynetest.NewRunner()
	runner.OutputDir = *outputDir
	runner.Verbose = *verbose

	// Print header
	fmt.Println("🧪 Fyne Visual Test Runner")
	fmt.Println("==========================")
	fmt.Printf("Plugin: %s\n", *pluginPath)
	fmt.Printf("Output directory: %s\n", runner.OutputDir)
	fmt.Println()

	// Run tests with timestamp
	results, runDir := runner.RunTestsWithTimestamp(testsToRun)

	// Count successes and failures
	successCount := 0
	failureCount := 0
	for _, result := range results {
		if result.Success {
			successCount++
		} else {
			failureCount++
			if !*verbose {
				fmt.Printf("❌ Test '%s' failed: %v\n", result.Test.Name, result.Error)
			}
		}
	}

	// Summary
	fmt.Println("\n📊 Test Summary")
	fmt.Println("===============")
	fmt.Printf("Total tests: %d\n", len(testsToRun))
	fmt.Printf("✅ Passed: %d\n", successCount)
	fmt.Printf("❌ Failed: %d\n", failureCount)
	fmt.Printf("\nScreenshots saved to: %s\n", runDir)

	// Generate HTML report
	reportGen := fynetest.NewReportGenerator()
	reportGen.Title = *reportTitle
	reportPath := filepath.Join(runDir, "index.html")
	if err := reportGen.GenerateHTMLReport(results, reportPath); err != nil {
		fmt.Printf("Warning: Failed to create HTML report: %v\n", err)
	} else {
		fmt.Printf("View results: file://%s\n", reportPath)
	}

	// Exit with error code if tests failed
	if failureCount > 0 {
		os.Exit(1)
	}
}