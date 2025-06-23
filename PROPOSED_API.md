# FyneTest: AI-Friendly Visual Testing API

## Overview

FyneTest is a visual testing framework for Fyne applications designed to provide AI models with comprehensive UI feedback and interaction capabilities.

## Core API Design

### 1. Test Definition

```go
// Test represents a single UI test scenario
type Test struct {
    Name        string
    Description string
    Tags        []string // For categorization and filtering
    Setup       SetupFunc
    Scenario    *Scenario // Optional: predefined interactions
}

// SetupFunc creates the UI component to test
type SetupFunc func() (fyne.CanvasObject, error)

// Scenario defines a sequence of interactions and validations
type Scenario struct {
    Steps []Step
}

// Step represents a single test step
type Step struct {
    Action     Action
    Validation *Validation
    Wait       time.Duration
}
```

### 2. UI State Capture

```go
// UIState represents a complete snapshot of the UI
type UIState struct {
    Timestamp   time.Time
    Screenshot  Screenshot
    Elements    []UIElement
    Theme       fyne.Theme
    Size        fyne.Size
    Metadata    map[string]interface{}
}

// Screenshot contains image data and metadata
type Screenshot struct {
    Data      []byte
    Format    string // "png", "jpeg"
    Width     int
    Height    int
    Path      string // Optional: file path if saved
}

// UIElement represents a UI component with semantic information
type UIElement struct {
    ID          string
    Type        string // "button", "entry", "label", etc.
    Text        string
    Position    fyne.Position
    Size        fyne.Size
    Visible     bool
    Enabled     bool
    Parent      *UIElement
    Children    []*UIElement
    Properties  map[string]interface{}
}
```

### 3. Interaction API

```go
// Runner provides test execution and interaction capabilities
type Runner interface {
    // Execute a test and return results
    Run(test Test) (*Result, error)
    
    // Capture current UI state
    CaptureState() (*UIState, error)
    
    // Perform an action
    DoAction(action Action) error
    
    // Find element by various criteria
    FindElement(selector Selector) (*UIElement, error)
    
    // Wait for condition
    WaitFor(condition Condition, timeout time.Duration) error
}

// Action represents a user interaction
type Action interface {
    Execute(r Runner) error
}

// Common actions
type ClickAction struct {
    Target Selector
}

type TypeAction struct {
    Target Selector
    Text   string
}

type DragAction struct {
    From Selector
    To   Selector
}

// Selector for finding elements
type Selector interface {
    Find(elements []UIElement) (*UIElement, error)
}

// Common selectors
type TextSelector struct {
    Text    string
    Partial bool
}

type TypeSelector struct {
    Type string // "button", "entry", etc.
}

type IDSelector struct {
    ID string
}

type PositionSelector struct {
    X, Y float32
}
```

### 4. AI-Friendly Features

```go
// AIAssistant provides high-level AI-friendly operations
type AIAssistant interface {
    // Describe the current UI in natural language
    DescribeUI() (string, error)
    
    // Find element by natural language description
    FindByDescription(description string) (*UIElement, error)
    
    // Compare two UI states
    CompareStates(before, after *UIState) (*StateDiff, error)
    
    // Suggest next actions based on current state
    SuggestActions() ([]Action, error)
    
    // Validate UI against natural language expectation
    ValidateExpectation(expectation string) (*ValidationResult, error)
}

// StateDiff represents differences between UI states
type StateDiff struct {
    Added      []UIElement
    Removed    []UIElement
    Modified   []ElementChange
    Visual     *VisualDiff
}

// VisualDiff provides pixel-level comparison
type VisualDiff struct {
    DiffImage     []byte
    DiffPercent   float64
    ChangedRegions []Region
}
```

### 5. Test Builder API

```go
// TestBuilder provides fluent API for test creation
type TestBuilder struct {
    test *Test
}

func NewTest(name string) *TestBuilder {
    return &TestBuilder{
        test: &Test{Name: name},
    }
}

func (b *TestBuilder) WithDescription(desc string) *TestBuilder
func (b *TestBuilder) WithSetup(setup SetupFunc) *TestBuilder
func (b *TestBuilder) WithScenario() *ScenarioBuilder
func (b *TestBuilder) Build() Test

// ScenarioBuilder for building test scenarios
type ScenarioBuilder struct {
    scenario *Scenario
}

func (b *ScenarioBuilder) Click(selector Selector) *ScenarioBuilder
func (b *ScenarioBuilder) Type(selector Selector, text string) *ScenarioBuilder
func (b *ScenarioBuilder) Wait(duration time.Duration) *ScenarioBuilder
func (b *ScenarioBuilder) Validate(validation Validation) *ScenarioBuilder
func (b *ScenarioBuilder) End() *TestBuilder
```

### 6. Configuration

```go
// Config for test execution
type Config struct {
    OutputDir      string
    ScreenshotFormat string
    Theme          fyne.Theme
    DefaultSize    fyne.Size
    Parallel       bool
    Verbose        bool
    AIMode         bool // Enable AI-specific features
}

// Suite manages multiple tests
type Suite struct {
    Tests  []Test
    Config Config
}

func (s *Suite) Run() (*SuiteResult, error)
func (s *Suite) RunWithAI(assistant AIAssistant) (*SuiteResult, error)
```

## Usage Examples

### Basic Test

```go
test := fynetest.NewTest("login_validation").
    WithDescription("Validate login form error handling").
    WithSetup(func() (fyne.CanvasObject, error) {
        return createLoginForm(), nil
    }).
    WithScenario().
        Type(fynetest.ID("username"), "user@example.com").
        Type(fynetest.ID("password"), "short").
        Click(fynetest.Text("Login")).
        Wait(100 * time.Millisecond).
        Validate(fynetest.TextVisible("Password must be at least 8 characters")).
    End().
    Build()

runner := fynetest.NewRunner()
result, err := runner.Run(test)
```

### AI-Driven Testing

```go
ai := fynetest.NewAIAssistant()
runner := fynetest.NewRunner()

// Natural language interaction
element, _ := ai.FindByDescription("the blue submit button")
runner.DoAction(fynetest.Click(element))

// Describe current state
description := ai.DescribeUI()
fmt.Println("Current UI:", description)

// Validate expectation
result, _ := ai.ValidateExpectation("error message should be visible")
```

### YAML Test Definition

```yaml
name: user_registration
description: Test user registration flow
setup:
  widget: RegistrationForm
  size: [600, 400]
scenario:
  - action: type
    target: {id: "email"}
    value: "test@example.com"
  - action: type
    target: {id: "password"}
    value: "securepass123"
  - action: click
    target: {text: "Register"}
  - wait: 500ms
  - validate:
      text_visible: "Registration successful"
```

## Migration from Current API

### Before
```go
suite := fynetest.NewTestSuite()
suite.Add(fynetest.Test{
    Name: "test",
    Setup: func() fyne.CanvasObject { return widget },
})
suite.RunCLI()
```

### After
```go
test := fynetest.NewTest("test").
    WithSetup(func() (fyne.CanvasObject, error) {
        return widget, nil
    }).
    Build()

suite := fynetest.NewSuite()
suite.AddTest(test)
suite.Run()
```

## Benefits for AI Integration

1. **Rich Metadata**: Every UI element includes semantic information
2. **Natural Language**: Support for describing and finding elements naturally
3. **State Management**: Complete UI state capture for analysis
4. **Interactive Testing**: Programmatic interaction capabilities
5. **Visual Analysis**: Built-in comparison and diffing tools
6. **Flexible Input**: Support for code, YAML, and JSON test definitions