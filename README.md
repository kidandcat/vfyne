# VFyne - Visual Testing Framework for Fyne Applications

VFyne is a powerful visual testing framework designed specifically for [Fyne](https://fyne.io) applications. It enables developers to create automated screenshots of UI components for visual regression testing, documentation, and AI-assisted analysis.

## 🌟 Key Features

- **📸 Automated Screenshot Capture** - Capture pixel-perfect screenshots of any Fyne UI component
- **🧪 Go Test Integration** - Works seamlessly with `go test` command
- **📷 Snapshot Testing** - Compare UI against baseline images with automatic diff generation
- **🔄 Consistent Naming** - Screenshots always use the same name and overwrite previous runs
- **🎨 Theme Support** - Test your UI across different themes (light, dark, custom)
- **📐 Flexible Sizing** - Test responsive layouts with custom window sizes
- **🏷️ Test Organization** - Tag and filter tests for better organization
- **📊 Beautiful Reports** - Generate interactive HTML reports with visual test results (CLI mode)
- **🚀 Simple API** - Intuitive builder pattern for creating tests
- **⚡ Fast Execution** - Efficient test runner with optional parallel execution
- **🤖 AI-Friendly** - Structured output perfect for AI analysis and automation

## 📦 Installation

```bash
go get github.com/jairo/vfyne
```

## 🚀 Quick Start

### Go Test Integration (Recommended)

VFyne now integrates seamlessly with Go's testing framework:

```go
package myapp_test

import (
    "testing"
    
    "fyne.io/fyne/v2/widget"
    vfyne "github.com/jairo/vfyne/testing"
)

func TestMyWidget(t *testing.T) {
    // Simple screenshot (always overwrites)
    vfyne.AssertScreenshot(t, "my_widget", widget.NewLabel("Hello, World!"))
    
    // Snapshot testing (compares with baseline)
    vfyne.AssertSnapshot(t, "my_snapshot", widget.NewButton("Click me", func() {}))
}
```

Run tests:
```bash
# Run tests and capture screenshots
go test -v

# Update snapshots
go test -v -update-snapshots
```

### Standalone CLI Mode

For standalone visual test suites:

```go
package main

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/widget"
    fynetest "github.com/jairo/vfyne"
)

func main() {
    // Create a test suite
    suite := fynetest.NewSuite()
    
    // Add a simple test
    suite.Add(fynetest.QuickTest("hello_world", func() fyne.CanvasObject {
        return widget.NewLabel("Hello, World!")
    }))
    
    // Run with CLI support
    suite.RunCLI()
}
```

### Using the Builder Pattern

```go
// Create a more detailed test
test := fynetest.NewTest("login_form").
    WithDescription("Login form with validation").
    WithSetup(func() fyne.CanvasObject {
        return createLoginForm()
    }).
    WithSize(400, 300).
    WithTags("forms", "authentication").
    Build()

suite.Add(test)
```

### Running Tests

```bash
# Run all tests
go run main.go

# Run specific test
go run main.go -test login_form

# Run tests with a specific tag
go run main.go -tag forms

# List available tests
go run main.go -list

# Generate report without running tests
go run main.go -no-report
```

## 📚 Advanced Usage

### Go Test Integration Features

#### Screenshot Testing
Screenshots are saved with consistent names and always overwrite previous runs:

```go
func TestScreenshots(t *testing.T) {
    vt := vfyne.New(t)
    
    // Basic screenshot
    vt.Screenshot("login_form", createLoginForm())
    
    // With custom size
    vt.Screenshot("dashboard", createDashboard(), vfyne.WithSize(1200, 800))
    
    // Mobile responsive
    vt.Screenshot("mobile_view", createMobileUI(), vfyne.WithMobileSize())
}
```

#### Snapshot Testing
Compare UI against baseline images:

```go
func TestSnapshots(t *testing.T) {
    vt := vfyne.New(t)
    
    // Will fail if image differs from snapshot
    vt.Snapshot("user_profile", createUserProfile())
    
    // Update snapshots with: go test -update-snapshots
}
```

When snapshots don't match:
- Test fails with detailed error
- Diff image is generated showing differences
- Actual output is saved for comparison

#### Theme Testing
```go
func TestThemes(t *testing.T) {
    content := createThemedContent()
    
    t.Run("Light", func(t *testing.T) {
        vt := vfyne.New(t)
        vt.SetTheme(theme.LightTheme())
        vt.Screenshot("theme_light", content)
    })
    
    t.Run("Dark", func(t *testing.T) {
        vt := vfyne.New(t)
        vt.SetTheme(theme.DarkTheme()) 
        vt.Screenshot("theme_dark", content)
    })
}
```

#### File Organization
```
myapp/
├── widget_test.go
└── testdata/
    ├── screenshots/      # Current test outputs (overwritten each run)
    │   ├── login_form.png
    │   ├── dashboard.png
    │   ├── actual_*.png  # Failed snapshot attempts
    │   └── diff_*.png    # Visual diffs
    └── snapshots/        # Baseline images for comparison
        ├── user_profile.png
        └── checkout_flow.png
```

### Custom Test Suite Configuration

```go
suite := fynetest.NewSuite().
    WithConfig(func(config *fynetest.SuiteConfig) {
        config.Name = "My App Visual Tests"
        config.OutputDir = "./screenshots"
        config.Verbose = true
        config.ReportTitle = "My App Test Results"
        config.Parallel = true
        config.MaxConcurrency = 4
    })
```

### Testing Different Themes

```go
// Test with dark theme
suite.AddBuilder(
    fynetest.NewTest("dark_theme_ui").
        WithDescription("UI in dark theme").
        WithTheme(theme.DarkTheme()).
        WithSetup(func() fyne.CanvasObject {
            return createMainUI()
        }),
)
```

### Testing Responsive Layouts

```go
// Test mobile layout
suite.AddBuilder(
    fynetest.NewTest("mobile_view").
        WithDescription("Mobile responsive layout").
        WithSize(375, 667). // iPhone SE size
        WithSetup(func() fyne.CanvasObject {
            return createMobileLayout()
        }).
        WithTags("mobile", "responsive"),
)
```

### Testing Form Validation

```go
suite.AddBuilder(
    fynetest.NewTest("form_errors").
        WithDescription("Form showing validation errors").
        WithSetup(func() fyne.CanvasObject {
            form := createRegistrationForm()
            // Trigger validation errors
            form.SubmitWithInvalidData()
            return form
        }).
        WithTags("forms", "validation", "error-states"),
)
```

## 📊 Output Structure

Tests generate organized output:
```
test-screenshots/
└── 20240119-143022/
    ├── hello_world_20240119-143022.png
    ├── login_form_20240119-143023.png
    ├── dark_theme_20240119-143024.png
    ├── index.html              # Interactive HTML report
    └── index.json             # Machine-readable JSON report
```

## 🤖 AI Integration

VFyne is designed to work seamlessly with AI tools:

```go
// The structured output makes it easy for AI to:
// 1. Understand UI structure
// 2. Analyze visual changes
// 3. Generate test descriptions
// 4. Detect anomalies

// JSON output includes:
// - Test metadata
// - Screenshot paths
// - Timing information
// - Success/failure status
// - Error messages
```

## 🛠️ API Reference

### Core Types

#### Test
```go
type Test struct {
    Name         string                    // Unique identifier
    Description  string                    // Human-readable description
    Setup        func() fyne.CanvasObject  // UI builder function
    Size         *fyne.Size                // Optional window size
    Theme        fyne.Theme                // Optional theme
    Tags         []string                  // Test categories
    WaitDuration time.Duration             // Render wait time
}
```

#### Runner
```go
type Runner struct {
    OutputDir           string      // Screenshot directory
    DefaultTheme        fyne.Theme  // Default theme
    DefaultSize         fyne.Size   // Default window size
    Verbose             bool        // Enable detailed logging
}
```

#### Suite
```go
type Suite struct {
    // Manages test collection and execution
}

// Key methods:
Add(test Test) *Suite
AddBuilder(builder *TestBuilder) *Suite
FilterByTags(tags ...string) []Test
Run() (SuiteResult, error)
RunCLI()
```

### Builder Pattern

```go
fynetest.NewTest(name string) *TestBuilder
    .WithDescription(string) *TestBuilder
    .WithSetup(func() fyne.CanvasObject) *TestBuilder
    .WithSize(width, height float32) *TestBuilder
    .WithTheme(fyne.Theme) *TestBuilder
    .WithTags(...string) *TestBuilder
    .WithWaitDuration(time.Duration) *TestBuilder
    .Build() (Test, error)
```

### Quick Helpers

```go
// Simple test with just name and setup
QuickTest(name string, setup func() fyne.CanvasObject) Test

// Test with description
QuickTestWithDescription(name, description string, setup func() fyne.CanvasObject) Test

// Test with specific theme
ThemeTest(name string, theme fyne.Theme, setup func() fyne.CanvasObject) Test

// Test with specific size
SizedTest(name string, width, height float32, setup func() fyne.CanvasObject) Test
```

## 🔧 Configuration Options

### Suite Configuration
```go
type SuiteConfig struct {
    Name            string      // Suite name
    OutputDir       string      // Output directory
    DefaultTheme    fyne.Theme  // Default theme
    DefaultSize     fyne.Size   // Default size
    Parallel        bool        // Enable parallel execution
    MaxConcurrency  int         // Max parallel tests
    Verbose         bool        // Verbose output
    GenerateReport  bool        // Generate HTML report
    ReportTitle     string      // Report title
}
```

### Command Line Flags
- `-output <dir>` - Output directory (default: "test-screenshots")
- `-test <name>` - Run specific test by name
- `-pattern <pattern>` - Run tests matching pattern
- `-tag <tag>` - Run tests with specific tag
- `-list` - List all available tests
- `-tags` - List all available tags
- `-verbose` - Enable verbose output
- `-parallel` - Run tests in parallel
- `-title <title>` - HTML report title
- `-no-report` - Skip HTML report generation

## 📝 Examples

See the [example directory](./example) for a comprehensive showcase of VFyne features, including:
- Basic component tests
- Form validation states  
- Theme variations
- Responsive layouts
- Complex dashboards
- Mobile interfaces

## 🤝 Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## 📄 License

This project is licensed under the same terms as the Fyne framework.

## 🙏 Acknowledgments

Built with ❤️ for the [Fyne](https://fyne.io) community.