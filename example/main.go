package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	fynetest "github.com/jairo/vfyne"
)

func main() {
	// Create a test suite with custom configuration
	suite := fynetest.NewSuite().
		WithConfig(func(config *fynetest.SuiteConfig) {
			config.Name = "Fyne UI Component Gallery"
			config.Verbose = true
			config.ReportTitle = "Fyne Component Visual Tests"
		})

	// Add tests using the builder pattern
	suite.
		// Simple tests
		AddBuilder(
			fynetest.NewTest("label_simple").
				WithDescription("A simple text label").
				WithSetup(func() fyne.CanvasObject {
					return widget.NewLabel("Hello, FyneTest!")
				}).
				WithTags("basic", "text"),
		).
		AddBuilder(
			fynetest.NewTest("label_multiline").
				WithDescription("Multi-line label with wrapping").
				WithSetup(func() fyne.CanvasObject {
					label := widget.NewLabel("This is a longer text that should wrap to multiple lines when the window is narrow enough. It demonstrates text wrapping capabilities.")
					label.Wrapping = fyne.TextWrapWord
					return label
				}).
				WithSize(300, 200).
				WithTags("basic", "text"),
		).
		
		// Button variations
		AddBuilder(
			fynetest.NewTest("buttons_showcase").
				WithDescription("Various button styles and states").
				WithSetup(func() fyne.CanvasObject {
					return container.NewVBox(
						widget.NewLabel("Button Showcase:"),
						widget.NewButton("Default Button", func() {}),
						widget.NewButtonWithIcon("Icon Button", theme.DocumentCreateIcon(), func() {}),
						&widget.Button{
							Text:       "High Importance",
							Importance: widget.HighImportance,
							OnTapped:   func() {},
						},
						&widget.Button{
							Text:       "Medium Importance",
							Importance: widget.MediumImportance,
							OnTapped:   func() {},
						},
						&widget.Button{
							Text:       "Low Importance",
							Importance: widget.LowImportance,
							OnTapped:   func() {},
						},
						&widget.Button{
							Text:       "Warning Button",
							Importance: widget.WarningImportance,
							OnTapped:   func() {},
						},
						&widget.Button{
							Text:       "Danger Button",
							Importance: widget.DangerImportance,
							OnTapped:   func() {},
						},
						func() fyne.CanvasObject {
							btn := widget.NewButton("Disabled Button", func() {})
							btn.Disable()
							return btn
						}(),
					)
				}).
				WithTags("buttons", "interactive"),
		).
		
		// Form examples
		AddBuilder(
			fynetest.NewTest("form_basic").
				WithDescription("Basic form with various input types").
				WithSetup(func() fyne.CanvasObject {
					nameEntry := widget.NewEntry()
					nameEntry.SetPlaceHolder("Enter your name")
					
					emailEntry := widget.NewEntry()
					emailEntry.SetPlaceHolder("email@example.com")
					
					passwordEntry := widget.NewPasswordEntry()
					passwordEntry.SetPlaceHolder("Password")
					
					ageSelect := widget.NewSelect([]string{"18-25", "26-35", "36-45", "46-55", "56+"}, func(string) {})
					
					agreeCheck := widget.NewCheck("I agree to the terms and conditions", func(bool) {})
					newsletterCheck := widget.NewCheck("Subscribe to newsletter", func(bool) {})
					
					form := widget.NewForm(
						widget.NewFormItem("Name", nameEntry),
						widget.NewFormItem("Email", emailEntry),
						widget.NewFormItem("Password", passwordEntry),
						widget.NewFormItem("Age Group", ageSelect),
						widget.NewFormItem("", container.NewVBox(agreeCheck, newsletterCheck)),
					)
					
					form.OnSubmit = func() {}
					form.OnCancel = func() {}
					
					return form
				}).
				WithTags("forms", "input"),
		).
		AddBuilder(
			fynetest.NewTest("form_validation").
				WithDescription("Form with validation errors displayed").
				WithSetup(func() fyne.CanvasObject {
					emailEntry := widget.NewEntry()
					emailEntry.SetText("invalid-email")
					
					passwordEntry := widget.NewPasswordEntry()
					passwordEntry.SetText("123") // Too short
					
					form := widget.NewForm(
						widget.NewFormItem("Email", emailEntry),
						widget.NewFormItem("Password", passwordEntry),
					)
					
					
					return form
				}).
				WithTags("forms", "validation", "error-state"),
		).
		
		// Card layouts
		AddBuilder(
			fynetest.NewTest("card_simple").
				WithDescription("Simple card with content").
				WithSetup(func() fyne.CanvasObject {
					content := container.NewVBox(
						widget.NewLabel("This is the card content"),
						widget.NewLabel("It can contain multiple elements"),
						widget.NewButton("Action", func() {}),
					)
					
					card := widget.NewCard(
						"Card Title",
						"This is a subtitle that provides more context",
						content,
					)
					
					return container.NewPadded(card)
				}).
				WithTags("layout", "card"),
		).
		
		// Progress indicators
		AddBuilder(
			fynetest.NewTest("progress_indicators").
				WithDescription("Various progress indicators").
				WithSetup(func() fyne.CanvasObject {
					progress := widget.NewProgressBar()
					progress.SetValue(0.65)
					
					infinite := widget.NewProgressBarInfinite()
					
					// Activity indicator placeholder
					activity := widget.NewProgressBarInfinite()
					
					return container.NewVBox(
						widget.NewLabel("Download Progress (65%):"),
						progress,
						widget.NewSeparator(),
						widget.NewLabel("Processing:"),
						infinite,
						widget.NewSeparator(),
						widget.NewLabel("Activity Indicator:"),
						container.NewCenter(activity),
					)
				}).
				WithTags("progress", "feedback"),
		).
		
		// Toolbar example
		AddBuilder(
			fynetest.NewTest("toolbar_comprehensive").
				WithDescription("Toolbar with various action types").
				WithSetup(func() fyne.CanvasObject {
					toolbar := widget.NewToolbar(
						widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {}),
						widget.NewToolbarAction(theme.FolderOpenIcon(), func() {}),
						widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {}),
						widget.NewToolbarSeparator(),
						widget.NewToolbarAction(theme.ContentCutIcon(), func() {}),
						widget.NewToolbarAction(theme.ContentCopyIcon(), func() {}),
						widget.NewToolbarAction(theme.ContentPasteIcon(), func() {}),
						widget.NewToolbarSeparator(),
						widget.NewToolbarAction(theme.ContentUndoIcon(), func() {}),
						widget.NewToolbarAction(theme.ContentRedoIcon(), func() {}),
						widget.NewToolbarSpacer(),
						widget.NewToolbarAction(theme.SettingsIcon(), func() {}),
					)
					
					content := widget.NewLabel("Document content area")
					
					return container.NewBorder(toolbar, nil, nil, nil, 
						container.NewPadded(content))
				}).
				WithTags("toolbar", "navigation"),
		).
		
		// Tab container
		AddBuilder(
			fynetest.NewTest("tabs_with_icons").
				WithDescription("Tab container with icons and content").
				WithSetup(func() fyne.CanvasObject {
					tabs := container.NewAppTabs(
						container.NewTabItemWithIcon("Home", theme.HomeIcon(), 
							widget.NewLabel("Welcome to the home tab")),
						container.NewTabItemWithIcon("Profile", theme.AccountIcon(),
							container.NewVBox(
								widget.NewLabel("User Profile"),
								widget.NewEntry(),
								widget.NewButton("Save", func() {}),
							)),
						container.NewTabItemWithIcon("Settings", theme.SettingsIcon(),
							widget.NewLabel("Application settings")),
						container.NewTabItemWithIcon("Help", theme.HelpIcon(),
							widget.NewLabel("Help and documentation")),
					)
					
					return tabs
				}).
				WithSize(600, 400).
				WithTags("tabs", "navigation"),
		).
		
		// Lists
		AddBuilder(
			fynetest.NewTest("list_simple").
				WithDescription("Simple list with items").
				WithSetup(func() fyne.CanvasObject {
					data := []string{
						"First Item",
						"Second Item",
						"Third Item",
						"Fourth Item",
						"Fifth Item",
					}
					
					list := widget.NewList(
						func() int { return len(data) },
						func() fyne.CanvasObject {
							return container.NewHBox(
								widget.NewIcon(theme.DocumentIcon()),
								widget.NewLabel("Template"),
							)
						},
						func(i widget.ListItemID, o fyne.CanvasObject) {
							o.(*fyne.Container).Objects[1].(*widget.Label).SetText(data[i])
						},
					)
					
					return list
				}).
				WithSize(400, 300).
				WithTags("list", "data"),
		).
		
		// Dark theme tests
		AddBuilder(
			fynetest.NewTest("dark_theme_showcase").
				WithDescription("UI components in dark theme").
				WithTheme(theme.DarkTheme()).
				WithSetup(func() fyne.CanvasObject {
					return container.NewVBox(
						widget.NewLabel("Dark Theme Showcase"),
						widget.NewButton("Dark Button", func() {}),
						widget.NewEntry(),
						widget.NewCheck("Dark checkbox", func(bool) {}),
						widget.NewRadioGroup([]string{"Option 1", "Option 2"}, func(string) {}),
						widget.NewSlider(0, 100),
					)
				}).
				WithTags("theme", "dark"),
		).
		
		// Mobile size test
		AddBuilder(
			fynetest.NewTest("mobile_layout").
				WithDescription("UI optimized for mobile screen size").
				WithSize(375, 667). // iPhone SE size
				WithSetup(func() fyne.CanvasObject {
					header := widget.NewToolbar(
						widget.NewToolbarAction(theme.MenuIcon(), func() {}),
						widget.NewToolbarSpacer(),
						widget.NewToolbarAction(theme.MoreVerticalIcon(), func() {}),
					)
					
					searchEntry := widget.NewEntry()
					searchEntry.SetPlaceHolder("Search...")
					
					list := widget.NewList(
						func() int { return 10 },
						func() fyne.CanvasObject {
							return container.NewHBox(
								widget.NewIcon(theme.AccountIcon()),
								widget.NewLabel("Contact Name"),
								layout.NewSpacer(),
								widget.NewLabel("Time"),
							)
						},
						func(i widget.ListItemID, o fyne.CanvasObject) {
							container := o.(*fyne.Container)
							container.Objects[1].(*widget.Label).SetText("Contact " + string(rune('A'+i)))
							container.Objects[3].(*widget.Label).SetText("2:30 PM")
						},
					)
					
					return container.NewBorder(
						container.NewVBox(header, searchEntry),
						nil, nil, nil,
						list,
					)
				}).
				WithTags("mobile", "responsive"),
		).
		
		// Complex layout
		AddBuilder(
			fynetest.NewTest("dashboard_layout").
				WithDescription("Complex dashboard layout").
				WithSize(1024, 768).
				WithSetup(func() fyne.CanvasObject {
					// Sidebar
					sidebar := container.NewVBox(
						widget.NewLabel("Dashboard"),
						widget.NewButton("Overview", func() {}),
						widget.NewButton("Analytics", func() {}),
						widget.NewButton("Reports", func() {}),
						widget.NewButton("Settings", func() {}),
					)
					
					// Stats cards
					statCard := func(title, value string) fyne.CanvasObject {
						return widget.NewCard(title, "", 
							container.NewCenter(widget.NewLabel(value)))
					}
					
					stats := container.NewGridWithColumns(3,
						statCard("Total Users", "1,234"),
						statCard("Revenue", "$45,678"),
						statCard("Growth", "+12.5%"),
					)
					
					// Chart placeholder
					chartPlaceholder := widget.NewCard("Sales Chart", "Last 30 days",
						container.NewCenter(
							widget.NewLabel("📊 Chart visualization would go here")))
					
					// Main content
					mainContent := container.NewVBox(
						stats,
						chartPlaceholder,
					)
					
					// Combine everything
					content := container.NewHSplit(
						container.NewPadded(sidebar),
						container.NewPadded(mainContent),
					)
					content.SetOffset(0.2) // 20% for sidebar
					
					return content
				}).
				WithTags("layout", "dashboard", "complex"),
		)

	// Using quick helper functions for simple tests
	suite.Add(fynetest.QuickTestWithDescription(
		"hyperlink_example",
		"Clickable hyperlink",
		func() fyne.CanvasObject {
			return widget.NewHyperlink("Visit Fyne.io", nil)
		},
	))

	// Run the test suite with CLI support
	suite.RunCLI()
}