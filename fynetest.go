// Package fynetest provides a visual testing framework for Fyne applications.
// It allows developers to create automated screenshots of UI components for
// visual regression testing, documentation, and AI-assisted understanding.
//
// Key features:
//   - Automated screenshot capture of Fyne UI components
//   - Support for different themes and window sizes
//   - HTML report generation with visual test results
//   - Simple builder API for creating tests
//   - AI-friendly metadata and state capture
//
// Basic usage:
//
//	test := fynetest.NewTest("my_test").
//		WithDescription("Test description").
//		WithSetup(func() fyne.CanvasObject {
//			return widget.NewLabel("Hello")
//		}).
//		Build()
//
//	runner := fynetest.NewRunner()
//	result := runner.RunTest(test)
package fynetest

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
)

// Test represents a visual test case for a Fyne UI component.
type Test struct {
	// Name is the unique identifier for this test (required)
	Name string
	
	// Description provides a human-readable explanation of what this test validates
	Description string
	
	// Tags allow categorization and filtering of tests
	Tags []string
	
	// Setup returns the Fyne canvas object to be tested (required)
	Setup func() fyne.CanvasObject
	
	// Size optionally specifies the window size for this test
	Size *fyne.Size
	
	// Theme optionally specifies a custom theme for this test
	Theme fyne.Theme
	
	// WaitDuration specifies how long to wait after showing the window (default: 100ms)
	WaitDuration time.Duration
	
	// Metadata allows storing additional information about the test
	Metadata map[string]interface{}
}

// Validate checks if the test configuration is valid
func (t *Test) Validate() error {
	if t.Name == "" {
		return fmt.Errorf("test name cannot be empty")
	}
	
	// Sanitize name for filesystem
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range invalidChars {
		if strings.Contains(t.Name, char) {
			return fmt.Errorf("test name contains invalid character: %s", char)
		}
	}
	
	if t.Setup == nil {
		return fmt.Errorf("test setup function cannot be nil")
	}
	
	if t.WaitDuration < 0 {
		return fmt.Errorf("wait duration cannot be negative")
	}
	
	return nil
}

// Result contains the outcome of running a visual test.
type Result struct {
	// Test is the test that was run
	Test Test
	
	// Success indicates whether the test passed
	Success bool
	
	// Error contains any error that occurred during the test
	Error error
	
	// ScreenshotPath is the file path where the screenshot was saved
	ScreenshotPath string
	
	// Screenshot contains the captured image data
	Screenshot image.Image
	
	// ImageSize is the size of the captured image
	ImageSize fyne.Size
	
	// Duration is how long the test took to run
	Duration time.Duration
	
	// Timestamp is when the test was run
	Timestamp time.Time
	
	// Metadata contains additional information about the test run
	Metadata map[string]interface{}
}

// Runner manages the execution of visual tests.
type Runner struct {
	// OutputDir is the directory where screenshots will be saved
	OutputDir string
	
	// DefaultTheme is the theme to use for tests that don't specify one
	DefaultTheme fyne.Theme
	
	// DefaultSize is the default window size for tests that don't specify one
	DefaultSize fyne.Size
	
	// DefaultWaitDuration is the default time to wait for window rendering
	DefaultWaitDuration time.Duration
	
	// Verbose enables detailed logging
	Verbose bool
	
	// app instance (reused across tests for efficiency)
	app fyne.App
	
	// mutex for thread safety
	mu sync.Mutex
}

// NewRunner creates a new test runner with sensible defaults.
func NewRunner() *Runner {
	return &Runner{
		OutputDir:           "test-screenshots",
		DefaultTheme:        theme.LightTheme(),
		DefaultSize:         fyne.NewSize(800, 600),
		DefaultWaitDuration: 100 * time.Millisecond,
		Verbose:             false,
	}
}

// ensureApp creates or returns the app instance
func (r *Runner) ensureApp() fyne.App {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if r.app == nil {
		r.app = app.New()
	}
	return r.app
}

// RunTest executes a single visual test and captures a screenshot.
func (r *Runner) RunTest(test Test) Result {
	startTime := time.Now()
	result := Result{
		Test:      test,
		Success:   false,
		Timestamp: startTime,
		Metadata:  make(map[string]interface{}),
	}
	
	// Validate test
	if err := test.Validate(); err != nil {
		result.Error = fmt.Errorf("invalid test configuration: %w", err)
		result.Duration = time.Since(startTime)
		return result
	}
	
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(r.OutputDir, 0755); err != nil {
		result.Error = fmt.Errorf("failed to create output directory: %w", err)
		result.Duration = time.Since(startTime)
		return result
	}
	
	// Get or create app instance
	testApp := r.ensureApp()
	
	// Set theme
	theme := test.Theme
	if theme == nil {
		theme = r.DefaultTheme
	}
	if theme != nil {
		testApp.Settings().SetTheme(theme)
	}
	
	// Create window
	window := testApp.NewWindow(test.Name)
	defer window.Close()
	
	// Get the content to test
	content := test.Setup()
	if content == nil {
		result.Error = fmt.Errorf("test setup returned nil content")
		result.Duration = time.Since(startTime)
		return result
	}
	
	// Set window content
	window.SetContent(content)
	
	// Calculate appropriate size
	size := r.calculateWindowSize(test, content)
	window.Resize(size)
	
	// Center window on screen (helps with consistency)
	window.CenterOnScreen()
	
	// Show the window to ensure it's rendered
	window.Show()
	
	// Wait for rendering
	waitDuration := test.WaitDuration
	if waitDuration == 0 {
		waitDuration = r.DefaultWaitDuration
	}
	time.Sleep(waitDuration)
	
	// Capture the image
	canvas := window.Canvas()
	if canvas == nil {
		result.Error = fmt.Errorf("failed to get canvas from window")
		result.Duration = time.Since(startTime)
		return result
	}
	
	img := canvas.Capture()
	if img == nil {
		result.Error = fmt.Errorf("failed to capture canvas image")
		result.Duration = time.Since(startTime)
		return result
	}
	
	result.Screenshot = img
	
	// Save the image
	timestamp := time.Now().Format("20060102-150405")
	filename := fmt.Sprintf("%s_%s.png", sanitizeFilename(test.Name), timestamp)
	filepath := filepath.Join(r.OutputDir, filename)
	
	if err := r.saveImage(img, filepath); err != nil {
		result.Error = fmt.Errorf("failed to save screenshot: %w", err)
		result.Duration = time.Since(startTime)
		return result
	}
	
	// Set result data
	result.Success = true
	result.ScreenshotPath = filepath
	result.ImageSize = fyne.NewSize(float32(img.Bounds().Dx()), float32(img.Bounds().Dy()))
	result.Duration = time.Since(startTime)
	
	// Add metadata
	result.Metadata["theme"] = getThemeName(theme)
	result.Metadata["window_size"] = size
	
	if r.Verbose {
		r.logTestResult(result)
	}
	
	return result
}

// RunTests executes multiple visual tests sequentially.
func (r *Runner) RunTests(tests []Test) []Result {
	results := make([]Result, 0, len(tests))
	
	for i, test := range tests {
		if r.Verbose {
			fmt.Printf("[%d/%d] Running test: %s\n", i+1, len(tests), test.Name)
		}
		result := r.RunTest(test)
		results = append(results, result)
		
		// Small delay between tests to ensure clean state
		if i < len(tests)-1 {
			time.Sleep(50 * time.Millisecond)
		}
	}
	
	return results
}

// RunTestsWithTimestamp executes tests in a timestamped subdirectory.
func (r *Runner) RunTestsWithTimestamp(tests []Test) ([]Result, string) {
	// Create timestamp for this test run
	timestamp := time.Now().Format("20060102-150405")
	originalOutputDir := r.OutputDir
	r.OutputDir = filepath.Join(originalOutputDir, timestamp)
	defer func() { r.OutputDir = originalOutputDir }()
	
	results := r.RunTests(tests)
	return results, r.OutputDir
}

// RunTestsConcurrent executes tests in parallel with a specified concurrency level.
func (r *Runner) RunTestsConcurrent(tests []Test, maxConcurrency int) []Result {
	if maxConcurrency <= 0 {
		maxConcurrency = 1
	}
	
	results := make([]Result, len(tests))
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, maxConcurrency)
	
	for i, test := range tests {
		wg.Add(1)
		go func(index int, t Test) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			
			if r.Verbose {
				fmt.Printf("Running test (concurrent): %s\n", t.Name)
			}
			results[index] = r.RunTest(t)
		}(i, test)
	}
	
	wg.Wait()
	return results
}

// Cleanup should be called when done with the runner to release resources
func (r *Runner) Cleanup() {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if r.app != nil {
		r.app.Quit()
		r.app = nil
	}
}

// Helper functions

func (r *Runner) calculateWindowSize(test Test, content fyne.CanvasObject) fyne.Size {
	if test.Size != nil {
		return *test.Size
	}
	
	minSize := content.MinSize()
	width := max(minSize.Width, r.DefaultSize.Width)
	height := max(minSize.Height, r.DefaultSize.Height)
	
	// Add some padding
	width += 20
	height += 20
	
	return fyne.NewSize(width, height)
}

func (r *Runner) saveImage(img image.Image, filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	
	return png.Encode(file, img)
}

func (r *Runner) logTestResult(result Result) {
	status := "✅ PASS"
	if !result.Success {
		status = "❌ FAIL"
	}
	
	fmt.Printf("%s Test '%s' completed in %v\n", status, result.Test.Name, result.Duration)
	
	if result.Test.Description != "" {
		fmt.Printf("   Description: %s\n", result.Test.Description)
	}
	
	if result.Success {
		fmt.Printf("   Screenshot: %s\n", result.ScreenshotPath)
		fmt.Printf("   Size: %dx%d pixels\n", int(result.ImageSize.Width), int(result.ImageSize.Height))
	} else {
		fmt.Printf("   Error: %v\n", result.Error)
	}
	
	fmt.Println()
}

func sanitizeFilename(name string) string {
	// Replace invalid characters with underscores
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|", " "}
	result := name
	for _, char := range invalidChars {
		result = strings.ReplaceAll(result, char, "_")
	}
	return result
}

func getThemeName(t fyne.Theme) string {
	if t == nil {
		return "default"
	}
	
	switch t {
	case theme.LightTheme():
		return "light"
	case theme.DarkTheme():
		return "dark"
	default:
		return "custom"
	}
}

func max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}