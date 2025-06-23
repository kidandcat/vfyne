package visual_test

import (
	"fmt"
	"strings"
	"testing"
	"time"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	vfyne "github.com/jairo/vfyne/testing"
)

func TestComplexDashboard(t *testing.T) {
	vt := vfyne.New(t)
	vt.SetRenderWait(200 * time.Millisecond)
	
	dashboard := createDashboard()
	vt.Snapshot("complex_dashboard", dashboard, vfyne.WithSize(1200, 800))
}

func TestTabContainer(t *testing.T) {
	vt := vfyne.New(t)
	
	tabs := container.NewAppTabs(
		container.NewTabItem("Overview", createOverviewTab()),
		container.NewTabItem("Settings", createSettingsTab()),
		container.NewTabItem("Logs", createLogsTab()),
	)
	
	vt.Snapshot("tab_container", tabs, vfyne.WithSize(800, 600))
}

func TestFormValidation(t *testing.T) {
	vt := vfyne.New(t)
	
	t.Run("EmptyForm", func(t *testing.T) {
		form := createValidationForm()
		vt.Snapshot("form_empty", form)
	})
	
	t.Run("InvalidForm", func(t *testing.T) {
		form := createValidationForm()
		form.(*widget.Form).Items[0].Widget.(*widget.Entry).SetText("a")
		form.(*widget.Form).Items[1].Widget.(*widget.Entry).SetText("invalid-email")
		vt.Snapshot("form_invalid", form)
	})
	
	t.Run("ValidForm", func(t *testing.T) {
		form := createValidationForm()
		form.(*widget.Form).Items[0].Widget.(*widget.Entry).SetText("John Doe")
		form.(*widget.Form).Items[1].Widget.(*widget.Entry).SetText("john@example.com")
		vt.Snapshot("form_valid", form)
	})
}

func createDashboard() fyne.CanvasObject {
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.HomeIcon(), func() {}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {}),
		widget.NewToolbarAction(theme.LogoutIcon(), func() {}),
	)
	
	stats := container.NewGridWithColumns(4,
		createStatCard("Users", "1,234", theme.AccountIcon()),
		createStatCard("Revenue", "$12,345", theme.StorageIcon()),
		createStatCard("Orders", "567", theme.DocumentIcon()),
		createStatCard("Growth", "+12%", theme.MoveUpIcon()),
	)
	
	chart := widget.NewCard("Performance", "Last 7 days",
		widget.NewProgressBar(),
	)
	
	list := widget.NewList(
		func() int { return 5 },
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil,
				widget.NewIcon(theme.DocumentIcon()),
				widget.NewLabel("Status"),
				widget.NewLabel("Item"),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			border := o.(*fyne.Container)
			border.Objects[3].(*widget.Label).SetText(fmt.Sprintf("Order #%d", 1000+i))
			border.Objects[4].(*widget.Label).SetText("Completed")
		},
	)
	
	return container.NewBorder(
		toolbar, nil, nil, nil,
		container.NewVSplit(
			container.NewVBox(stats, chart),
			widget.NewCard("Recent Orders", "", list),
		),
	)
}

func createStatCard(title, value string, icon fyne.Resource) fyne.CanvasObject {
	return widget.NewCard("", "",
		container.NewBorder(nil, nil,
			widget.NewIcon(icon),
			nil,
			container.NewVBox(
				widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{}),
				widget.NewLabelWithStyle(value, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			),
		),
	)
}

func createOverviewTab() fyne.CanvasObject {
	return container.NewVBox(
		widget.NewLabel("Overview Content"),
		widget.NewSeparator(),
		widget.NewLabel("Dashboard statistics and charts would go here"),
	)
}

func createSettingsTab() fyne.CanvasObject {
	return widget.NewForm(
		widget.NewFormItem("Theme", widget.NewSelect([]string{"Light", "Dark", "Auto"}, func(string) {})),
		widget.NewFormItem("Language", widget.NewSelect([]string{"English", "Spanish", "French"}, func(string) {})),
		widget.NewFormItem("Notifications", widget.NewCheck("Enable notifications", func(bool) {})),
	)
}

func createLogsTab() fyne.CanvasObject {
	entry := widget.NewMultiLineEntry()
	entry.SetText("2024-01-20 10:30:45 - Application started\n" +
		"2024-01-20 10:30:46 - Connected to database\n" +
		"2024-01-20 10:30:47 - Server listening on port 8080\n")
	entry.Disable()
	return entry
}

func createValidationForm() fyne.CanvasObject {
	nameEntry := widget.NewEntry()
	nameEntry.Validator = func(s string) error {
		if len(s) < 3 {
			return fmt.Errorf("Name must be at least 3 characters")
		}
		return nil
	}
	
	emailEntry := widget.NewEntry()
	emailEntry.Validator = func(s string) error {
		if !strings.Contains(s, "@") {
			return fmt.Errorf("Invalid email address")
		}
		return nil
	}
	
	return widget.NewForm(
		widget.NewFormItem("Name", nameEntry),
		widget.NewFormItem("Email", emailEntry),
	)
}