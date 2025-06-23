# VFyne Visual Testing Examples

This directory contains examples of how to use VFyne with Go's testing framework.

## Running the Tests

### Basic Screenshots
To capture screenshots (overwrites previous ones):
```bash
go test -v
```

### Snapshot Testing
To update snapshots:
```bash
go test -v -update-snapshots
```

To compare against existing snapshots:
```bash
go test -v
```

## Test Files

- **basic_test.go** - Simple widget screenshots and basic snapshot testing
- **theme_test.go** - Testing with different themes
- **advanced_test.go** - Complex layouts, forms with validation, and dashboard examples

## Directory Structure

After running tests, you'll see:
```
testdata/
├── screenshots/     # Current test outputs (always overwritten)
│   ├── basic_label.png
│   ├── basic_button.png
│   └── ...
└── snapshots/       # Baseline images for comparison
    ├── button_snapshot.png
    ├── card_snapshot.png
    └── ...
```

## Key Features

1. **Screenshot Mode**: Simple screenshot capture, always overwrites
2. **Snapshot Mode**: Compares against baseline, fails if different
3. **Diff Generation**: Creates diff images when snapshots don't match
4. **Responsive Testing**: Test with different window sizes
5. **Theme Support**: Test with light/dark themes

## Example Usage

```go
func TestMyWidget(t *testing.T) {
    vt := vfyne.New(t)
    
    // Simple screenshot (always overwrites)
    widget := widget.NewLabel("Hello")
    vt.Screenshot("my_label", widget)
    
    // Snapshot testing (compares with baseline)
    vt.Snapshot("my_snapshot", widget)
    
    // With custom size
    vt.Screenshot("my_large_widget", widget, vfyne.WithSize(1024, 768))
    
    // Quick assertions
    vfyne.AssertSnapshot(t, "quick_test", widget)
}
```