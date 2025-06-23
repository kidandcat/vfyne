package fynetest

import (
	"time"

	"fyne.io/fyne/v2"
)

// TestBuilder provides a fluent interface for creating tests.
type TestBuilder struct {
	test *Test
}

// NewTest creates a new test builder with the given name.
// The name must be unique and will be used as the filename for screenshots.
func NewTest(name string) *TestBuilder {
	return &TestBuilder{
		test: &Test{
			Name:     name,
			Tags:     make([]string, 0),
			Metadata: make(map[string]interface{}),
		},
	}
}

// WithDescription sets the test description.
func (b *TestBuilder) WithDescription(description string) *TestBuilder {
	b.test.Description = description
	return b
}

// WithSetup sets the function that creates the UI to test.
// This is required and must return a non-nil fyne.CanvasObject.
func (b *TestBuilder) WithSetup(setup func() fyne.CanvasObject) *TestBuilder {
	b.test.Setup = setup
	return b
}

// WithSize sets a custom window size for this test.
// If not set, the window will use the content's minimum size or the runner's default.
func (b *TestBuilder) WithSize(width, height float32) *TestBuilder {
	size := fyne.NewSize(width, height)
	b.test.Size = &size
	return b
}

// WithTheme sets a custom theme for this test.
// If not set, the runner's default theme will be used.
func (b *TestBuilder) WithTheme(theme fyne.Theme) *TestBuilder {
	b.test.Theme = theme
	return b
}

// WithWaitDuration sets how long to wait after showing the window before capturing.
// This can be useful for animations or async rendering. Default is 100ms.
func (b *TestBuilder) WithWaitDuration(duration time.Duration) *TestBuilder {
	b.test.WaitDuration = duration
	return b
}

// WithTags adds tags for categorizing and filtering tests.
func (b *TestBuilder) WithTags(tags ...string) *TestBuilder {
	b.test.Tags = append(b.test.Tags, tags...)
	return b
}

// WithMetadata adds custom metadata to the test.
func (b *TestBuilder) WithMetadata(key string, value interface{}) *TestBuilder {
	b.test.Metadata[key] = value
	return b
}

// Build creates the final Test instance.
// This will validate the test configuration and return an error if invalid.
func (b *TestBuilder) Build() (Test, error) {
	if err := b.test.Validate(); err != nil {
		return Test{}, err
	}
	return *b.test, nil
}

// MustBuild creates the final Test instance, panicking if validation fails.
// Use this when you're certain the test configuration is valid.
func (b *TestBuilder) MustBuild() Test {
	test, err := b.Build()
	if err != nil {
		panic(err)
	}
	return test
}

// Quick helper functions for common test patterns

// QuickTest creates a simple test with just a name and setup function.
func QuickTest(name string, setup func() fyne.CanvasObject) Test {
	return NewTest(name).WithSetup(setup).MustBuild()
}

// QuickTestWithDescription creates a test with a name, description, and setup.
func QuickTestWithDescription(name, description string, setup func() fyne.CanvasObject) Test {
	return NewTest(name).
		WithDescription(description).
		WithSetup(setup).
		MustBuild()
}

// ThemeTest creates a test that runs with a specific theme.
func ThemeTest(name string, theme fyne.Theme, setup func() fyne.CanvasObject) Test {
	return NewTest(name).
		WithTheme(theme).
		WithSetup(setup).
		MustBuild()
}

// SizedTest creates a test with a specific window size.
func SizedTest(name string, width, height float32, setup func() fyne.CanvasObject) Test {
	return NewTest(name).
		WithSize(width, height).
		WithSetup(setup).
		MustBuild()
}