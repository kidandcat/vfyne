# AI Integration Example

This example demonstrates how an AI model could use the improved FyneTest API to understand and interact with a Fyne application.

## Scenario: AI Testing a Calculator App

```go
package main

import (
    "fmt"
    "github.com/fynetest/fynetest"
    "github.com/fynetest/fynetest/ai"
)

func main() {
    // Initialize AI-enabled runner
    runner := fynetest.NewAIRunner(fynetest.Config{
        AIMode: true,
        OutputDir: "./test-results",
    })
    
    // Create test for calculator
    test := fynetest.NewTest("calculator_operations").
        WithDescription("Test basic calculator operations").
        WithSetup(func() (fyne.CanvasObject, error) {
            return createCalculatorUI(), nil
        }).
        Build()
    
    // Start test
    runner.StartTest(test)
    
    // AI explores the UI
    exploreCalculator(runner)
}

func exploreCalculator(runner *fynetest.AIRunner) {
    // 1. Get initial UI description
    state, _ := runner.CaptureState()
    description := runner.DescribeUI()
    fmt.Println("Initial UI:", description)
    // Output: "Calculator interface with number buttons 0-9, operation buttons (+, -, *, /), 
    //          equals button, and a display showing '0'"
    
    // 2. AI performs calculation: 15 + 7
    performCalculation(runner, "15", "+", "7")
    
    // 3. Verify result
    verifyResult(runner, "22")
    
    // 4. Test edge cases
    testEdgeCases(runner)
}

func performCalculation(runner *fynetest.AIRunner, num1, op, num2 string) {
    // AI finds and clicks number buttons for first number
    for _, digit := range num1 {
        element, _ := runner.FindByDescription(fmt.Sprintf("button with text '%c'", digit))
        runner.Click(element)
    }
    
    // Click operation
    opElement, _ := runner.FindByDescription(fmt.Sprintf("operation button '%s'", op))
    runner.Click(opElement)
    
    // Enter second number
    for _, digit := range num2 {
        element, _ := runner.FindByDescription(fmt.Sprintf("button with text '%c'", digit))
        runner.Click(element)
    }
    
    // Click equals
    equals, _ := runner.FindByDescription("equals button")
    runner.Click(equals)
}

func verifyResult(runner *fynetest.AIRunner, expected string) {
    // Capture current state
    state, _ := runner.CaptureState()
    
    // Find display element
    display, _ := runner.FindByDescription("calculator display")
    
    // Verify the result
    validation := runner.ValidateExpectation(
        fmt.Sprintf("display should show '%s'", expected),
    )
    
    if validation.Success {
        fmt.Printf("✓ Calculation correct: display shows %s\n", display.Text)
    } else {
        fmt.Printf("✗ Calculation failed: expected %s, got %s\n", expected, display.Text)
    }
}

func testEdgeCases(runner *fynetest.AIRunner) {
    testCases := []struct {
        name        string
        actions     []string
        expectation string
    }{
        {
            name:        "Division by zero",
            actions:     []string{"5", "/", "0", "="},
            expectation: "display should show error or infinity",
        },
        {
            name:        "Large numbers",
            actions:     []string{"9999999", "*", "9999999", "="},
            expectation: "display should handle large number or show overflow",
        },
        {
            name:        "Decimal operations",
            actions:     []string{"3", ".", "1", "4", "+", "2", ".", "7", "="},
            expectation: "display should show approximately 5.84",
        },
    }
    
    for _, tc := range testCases {
        fmt.Printf("\nTesting: %s\n", tc.name)
        
        // Clear calculator
        runner.Click(runner.FindByDescription("clear button"))
        
        // Perform actions
        for _, action := range tc.actions {
            if element, err := runner.FindByDescription(
                fmt.Sprintf("button or key '%s'", action),
            ); err == nil {
                runner.Click(element)
                runner.Wait(100 * time.Millisecond)
            }
        }
        
        // Validate expectation
        result := runner.ValidateExpectation(tc.expectation)
        fmt.Printf("Result: %v\n", result.Message)
    }
}
```

## AI Assistant Interaction Example

```go
// Example of how an AI model could interact with the test framework
func aiTestSession(assistant fynetest.AIAssistant) {
    // AI can ask about the current UI
    question := "What buttons are visible on the screen?"
    response := assistant.AnswerQuestion(question)
    // Response: "Number buttons 0-9, operation buttons +, -, *, /, 
    //           equals button (=), clear button (C), and decimal point (.)"
    
    // AI can request specific actions
    assistant.ExecuteNaturalCommand("Click the number 7")
    assistant.ExecuteNaturalCommand("Click the plus button")
    assistant.ExecuteNaturalCommand("Click the number 3")
    assistant.ExecuteNaturalCommand("Press equals")
    
    // AI can verify results
    verification := assistant.VerifyNaturalExpectation(
        "The display should show 10",
    )
    
    // AI can compare states
    beforeState := assistant.GetCurrentState()
    assistant.ExecuteNaturalCommand("Clear the calculator")
    afterState := assistant.GetCurrentState()
    
    diff := assistant.CompareStates(beforeState, afterState)
    // diff.Description: "Display changed from '10' to '0', all buttons remain in default state"
}
```

## Integration with AI Testing Frameworks

```go
// Example integration with a hypothetical AI testing framework
type AITestDriver struct {
    fynetest *fynetest.AIRunner
}

func (d *AITestDriver) ProcessTestScript(script string) {
    // Parse natural language test script
    // Example script:
    // "Open calculator and verify it shows 0
    //  Calculate 15 + 7 and verify result is 22
    //  Try dividing by zero and verify error handling"
    
    steps := d.parseScript(script)
    for _, step := range steps {
        d.executeStep(step)
    }
}

func (d *AITestDriver) GenerateTestReport() TestReport {
    return TestReport{
        Screenshots: d.fynetest.GetAllScreenshots(),
        Actions:     d.fynetest.GetActionHistory(),
        Validations: d.fynetest.GetValidationResults(),
        Summary:     d.generateSummary(),
    }
}
```

## Benefits for AI Models

1. **Natural Language Understanding**: AI can describe what it sees and request actions in natural language
2. **Visual Feedback**: Every action produces a screenshot that AI can analyze
3. **Semantic Information**: UI elements have meaning beyond just pixels
4. **State Management**: AI can track changes and understand UI flow
5. **Validation Framework**: AI can verify expectations in human terms
6. **Exploration Support**: AI can discover UI capabilities through interaction