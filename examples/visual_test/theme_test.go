package visual_test

import (
	"testing"
	
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	vfyne "github.com/jairo/vfyne/testing"
)

func TestThemes(t *testing.T) {
	createContent := func() fyne.CanvasObject {
		return container.NewVBox(
			widget.NewLabel("Theme Test"),
			widget.NewButton("Primary", func() {}),
			widget.NewButtonWithIcon("Icon Button", theme.ConfirmIcon(), func() {}),
			widget.NewEntry(),
			widget.NewCheck("Check me", func(bool) {}),
			widget.NewRadioGroup([]string{"Option 1", "Option 2", "Option 3"}, func(string) {}),
			widget.NewProgressBar(),
			widget.NewSlider(0, 100),
		)
	}
	
	t.Run("LightTheme", func(t *testing.T) {
		vt := vfyne.New(t)
		vt.SetTheme(theme.LightTheme())
		vt.Screenshot("theme_light", createContent())
	})
	
	t.Run("DarkTheme", func(t *testing.T) {
		vt := vfyne.New(t)
		vt.SetTheme(theme.DarkTheme())
		vt.Screenshot("theme_dark", createContent())
	})
}

func TestButtonImportance(t *testing.T) {
	vt := vfyne.New(t)
	
	buttons := container.NewVBox(
		widget.NewButton("Default Button", func() {}),
		widget.NewButtonWithIcon("High Importance", theme.WarningIcon(), func() {}),
		widget.NewButtonWithIcon("Success", theme.ConfirmIcon(), func() {}),
		widget.NewButtonWithIcon("Danger", theme.DeleteIcon(), func() {}),
	)
	
	buttons.Objects[1].(*widget.Button).Importance = widget.HighImportance
	buttons.Objects[2].(*widget.Button).Importance = widget.SuccessImportance
	buttons.Objects[3].(*widget.Button).Importance = widget.DangerImportance
	
	vt.Screenshot("button_importance", buttons)
}