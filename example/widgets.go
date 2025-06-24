package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// CreateLoginForm creates a sample login form
func CreateLoginForm() fyne.CanvasObject {
	username := widget.NewEntry()
	username.SetPlaceHolder("Username")
	
	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Password")
	
	remember := widget.NewCheck("Remember me", nil)
	
	loginBtn := widget.NewButton("Login", func() {})
	loginBtn.Importance = widget.HighImportance
	
	form := container.NewVBox(
		widget.NewLabel("Login"),
		username,
		password,
		remember,
		loginBtn,
	)
	
	return container.NewPadded(form)
}

// CreateDashboard creates a sample dashboard
func CreateDashboard() fyne.CanvasObject {
	// Progress card
	progress := widget.NewProgressBar()
	progress.SetValue(0.7)
	progressCard := widget.NewCard("Progress", "70% Complete", progress)
	
	// List widget
	data := []string{"Item 1", "Item 2", "Item 3", "Item 4", "Item 5"}
	list := widget.NewList(
		func() int { return len(data) },
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(data[i])
		},
	)
	listCard := widget.NewCard("Recent Items", "", list)
	
	// Action buttons
	actions := container.NewHBox(
		widget.NewButton("Add", func() {}),
		widget.NewButton("Edit", func() {}),
		widget.NewButton("Delete", func() {}),
	)
	
	// Stats
	stat1 := widget.NewCard("Users", "1,234", widget.NewIcon(theme.AccountIcon()))
	stat2 := widget.NewCard("Revenue", "$5,678", widget.NewIcon(theme.ComputerIcon()))
	stat3 := widget.NewCard("Orders", "456", widget.NewIcon(theme.DocumentIcon()))
	
	stats := container.NewGridWithColumns(3, stat1, stat2, stat3)
	
	return container.NewVBox(
		stats,
		container.NewGridWithColumns(2,
			progressCard,
			listCard,
		),
		actions,
	)
}

// CreateSettingsForm creates a settings form with various widgets
func CreateSettingsForm() fyne.CanvasObject {
	// Theme selection
	themeSelect := widget.NewSelect([]string{"Light", "Dark", "Auto"}, nil)
	themeSelect.SetSelected("Auto")
	
	// Notification settings
	notifyEmail := widget.NewCheck("Email notifications", nil)
	notifyEmail.SetChecked(true)
	notifyPush := widget.NewCheck("Push notifications", nil)
	
	// Slider
	volumeSlider := binding.NewFloat()
	volumeSlider.Set(0.5)
	slider := widget.NewSliderWithData(0, 1, volumeSlider)
	
	// Radio group
	language := widget.NewRadioGroup([]string{"English", "Spanish", "French"}, nil)
	language.SetSelected("English")
	
	// Form
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Theme", Widget: themeSelect},
			{Text: "Notifications", Widget: container.NewVBox(notifyEmail, notifyPush)},
			{Text: "Volume", Widget: slider},
			{Text: "Language", Widget: language},
		},
		OnCancel: func() {},
		OnSubmit: func() {},
	}
	
	return container.NewPadded(form)
}

// CreateComplexLayout creates a complex layout with tabs and split containers
func CreateComplexLayout() fyne.CanvasObject {
	// Tabs
	tabs := container.NewAppTabs(
		container.NewTabItem("Login", CreateLoginForm()),
		container.NewTabItem("Dashboard", CreateDashboard()),
		container.NewTabItem("Settings", CreateSettingsForm()),
	)
	
	// Sidebar
	sidebar := container.NewVBox(
		widget.NewCard("Navigation", "", 
			widget.NewList(
				func() int { return 5 },
				func() fyne.CanvasObject {
					return widget.NewButtonWithIcon("", theme.DocumentIcon(), nil)
				},
				func(i widget.ListItemID, o fyne.CanvasObject) {
					labels := []string{"Home", "Profile", "Messages", "Settings", "Logout"}
					o.(*widget.Button).SetText(labels[i])
				},
			),
		),
	)
	
	// Create split container
	split := container.NewHSplit(sidebar, tabs)
	split.SetOffset(0.2)
	
	return split
}

// CreateDataTable creates a sample data table
func CreateDataTable() fyne.CanvasObject {
	table := widget.NewTable(
		func() (int, int) { return 5, 3 },
		func() fyne.CanvasObject {
			return widget.NewLabel("Cell")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			label := o.(*widget.Label)
			if i.Row == 0 {
				headers := []string{"Name", "Age", "City"}
				label.SetText(headers[i.Col])
				label.TextStyle = fyne.TextStyle{Bold: true}
			} else {
				data := [][]string{
					{"John Doe", "30", "New York"},
					{"Jane Smith", "25", "London"},
					{"Bob Johnson", "35", "Paris"},
					{"Alice Brown", "28", "Tokyo"},
				}
				label.SetText(data[i.Row-1][i.Col])
				label.TextStyle = fyne.TextStyle{}
			}
		},
	)
	
	table.SetColumnWidth(0, 150)
	table.SetColumnWidth(1, 80)
	table.SetColumnWidth(2, 120)
	
	return widget.NewCard("User Data", "", table)
}

// CreateErrorDialog creates an error dialog example
func CreateErrorDialog() fyne.CanvasObject {
	// Simulate an error message
	errorContent := container.NewVBox(
		widget.NewLabelWithStyle("An error occurred!", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabel("Unable to connect to the server. Please check your internet connection and try again."),
		container.NewHBox(
			widget.NewButton("Retry", func() {}),
			widget.NewButton("Cancel", func() {}),
		),
	)
	
	return widget.NewCard("Error", "", errorContent)
}