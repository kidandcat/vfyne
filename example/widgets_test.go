package main

import (
	"testing"
	
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	vfyne "github.com/jairo/vfyne/testing"
)

func TestLoginForm(t *testing.T) {
	vt := vfyne.New(t)
	
	// Basic screenshot
	vt.Screenshot("login_form", CreateLoginForm())
	
	// With custom size
	vt.Screenshot("login_form_mobile", CreateLoginForm(), vfyne.WithSize(400, 600))
	
	// Snapshot testing (for regression testing)
	vt.Snapshot("login_form_snapshot", CreateLoginForm())
}

func TestDashboard(t *testing.T) {
	vt := vfyne.New(t)
	
	// Test dashboard in different sizes
	vt.Screenshot("dashboard_desktop", CreateDashboard(), vfyne.WithSize(1200, 800))
	vt.Screenshot("dashboard_tablet", CreateDashboard(), vfyne.WithSize(768, 1024))
	
	// Test with dark theme
	vt.SetTheme(theme.DarkTheme())
	vt.Screenshot("dashboard_dark", CreateDashboard())
	vt.SetTheme(theme.LightTheme())
}

func TestSettingsForm(t *testing.T) {
	vt := vfyne.New(t)
	
	// Basic settings form
	vt.Screenshot("settings_form", CreateSettingsForm())
	
	// Snapshot test
	vt.Snapshot("settings_snapshot", CreateSettingsForm(), vfyne.WithSize(600, 500))
}

func TestComplexLayout(t *testing.T) {
	vt := vfyne.New(t)
	
	// Test complex layout with different viewport sizes
	vt.Screenshot("complex_layout_full", CreateComplexLayout(), vfyne.WithSize(1400, 900))
	vt.Screenshot("complex_layout_compact", CreateComplexLayout(), vfyne.WithSize(800, 600))
}

func TestDataTable(t *testing.T) {
	vt := vfyne.New(t)
	
	// Test data table
	vt.Screenshot("data_table", CreateDataTable(), vfyne.WithSize(600, 400))
	
	// Snapshot for regression testing
	vt.Snapshot("data_table_snapshot", CreateDataTable())
}

func TestErrorDialog(t *testing.T) {
	vt := vfyne.New(t)
	
	// Test error dialog in both themes
	vt.Screenshot("error_dialog_light", CreateErrorDialog())
	vt.SetTheme(theme.DarkTheme())
	vt.Screenshot("error_dialog_dark", CreateErrorDialog())
	vt.SetTheme(theme.LightTheme())
}

func TestAllComponents(t *testing.T) {
	vt := vfyne.New(t)
	
	// Create a comprehensive view with all components
	allComponents := container.NewVScroll(container.NewVBox(
		widget.NewCard("Login Form", "", CreateLoginForm()),
		widget.NewSeparator(),
		widget.NewCard("Dashboard", "", CreateDashboard()),
		widget.NewSeparator(),
		widget.NewCard("Settings", "", CreateSettingsForm()),
		widget.NewSeparator(),
		CreateDataTable(),
		widget.NewSeparator(),
		CreateErrorDialog(),
	))
	
	// Capture full component showcase
	vt.Screenshot("all_components", allComponents, vfyne.WithSize(800, 1200))
}

// TestThemeComparison demonstrates theme testing
func TestThemeComparison(t *testing.T) {
	vt := vfyne.New(t)
	
	// Simple button in both themes
	button := widget.NewButton("Click Me!", func() {})
	button.Importance = widget.HighImportance
	
	vt.Screenshot("button_light", button)
	vt.SetTheme(theme.DarkTheme())
	vt.Screenshot("button_dark", button)
	vt.SetTheme(theme.LightTheme())
	
	// Form in both themes
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Name", Widget: widget.NewEntry()},
			{Text: "Email", Widget: widget.NewEntry()},
			{Text: "Subscribe", Widget: widget.NewCheck("Newsletter", nil)},
		},
	}
	
	vt.Screenshot("form_theme_light", form, vfyne.WithSize(400, 300))
	vt.SetTheme(theme.DarkTheme())
	vt.Screenshot("form_theme_dark", form, vfyne.WithSize(400, 300))
	vt.SetTheme(theme.LightTheme())
}

// TestResponsiveDesign tests UI at different sizes
func TestResponsiveDesign(t *testing.T) {
	vt := vfyne.New(t)
	
	// Create a responsive layout
	responsive := container.NewBorder(
		widget.NewToolbar(
			widget.NewToolbarAction(theme.MenuIcon(), func() {}),
			widget.NewToolbarSeparator(),
			widget.NewToolbarAction(theme.SettingsIcon(), func() {}),
		),
		widget.NewLabel("Footer"),
		nil,
		nil,
		container.NewGridWithColumns(2,
			widget.NewCard("Card 1", "Content", widget.NewButton("Action", func() {})),
			widget.NewCard("Card 2", "Content", widget.NewButton("Action", func() {})),
		),
	)
	
	// Test at different viewport sizes
	sizes := []struct {
		name   string
		width  float32
		height float32
	}{
		{"mobile", 375, 667},
		{"tablet", 768, 1024},
		{"desktop", 1920, 1080},
	}
	
	for _, size := range sizes {
		vt.Screenshot("responsive_"+size.name, responsive, vfyne.WithSize(size.width, size.height))
	}
}