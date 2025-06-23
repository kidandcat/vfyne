package visual_test

import (
	"testing"
	
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	vfyne "github.com/jairo/vfyne/testing"
)

func TestBasicWidgets(t *testing.T) {
	vt := vfyne.New(t)
	
	t.Run("Label", func(t *testing.T) {
		label := widget.NewLabel("Hello, VFyne!")
		vt.Screenshot("basic_label", label)
	})
	
	t.Run("Button", func(t *testing.T) {
		button := widget.NewButton("Click Me", func() {})
		vt.Screenshot("basic_button", button)
	})
	
	t.Run("Entry", func(t *testing.T) {
		entry := widget.NewEntry()
		entry.SetPlaceHolder("Enter text here...")
		vt.Screenshot("basic_entry", entry)
	})
	
	t.Run("Form", func(t *testing.T) {
		form := widget.NewForm(
			widget.NewFormItem("Name", widget.NewEntry()),
			widget.NewFormItem("Email", widget.NewEntry()),
			widget.NewFormItem("Message", widget.NewMultiLineEntry()),
		)
		vt.Screenshot("basic_form", form, vfyne.WithSize(400, 300))
	})
}

func TestSnapshotComparison(t *testing.T) {
	vt := vfyne.New(t)
	
	t.Run("ButtonSnapshot", func(t *testing.T) {
		button := widget.NewButton("Snapshot Test", func() {})
		vt.Snapshot("button_snapshot", button)
	})
	
	t.Run("CardSnapshot", func(t *testing.T) {
		card := widget.NewCard(
			"Card Title",
			"This is a card subtitle",
			widget.NewLabel("Card content goes here"),
		)
		vt.Snapshot("card_snapshot", card)
	})
}

func TestResponsiveLayouts(t *testing.T) {
	vt := vfyne.New(t)
	
	content := container.NewVBox(
		widget.NewLabel("Responsive Layout Test"),
		widget.NewButton("Action", func() {}),
		widget.NewEntry(),
	)
	
	t.Run("Desktop", func(t *testing.T) {
		vt.Screenshot("responsive_desktop", content, vfyne.WithSize(800, 600))
	})
	
	t.Run("Tablet", func(t *testing.T) {
		vt.Screenshot("responsive_tablet", content, vfyne.WithTabletSize())
	})
	
	t.Run("Mobile", func(t *testing.T) {
		vt.Screenshot("responsive_mobile", content, vfyne.WithMobileSize())
	})
}

func TestQuickAssertions(t *testing.T) {
	vfyne.AssertScreenshot(t, "quick_label", widget.NewLabel("Quick test"))
	
	vfyne.AssertSnapshot(t, "quick_button", widget.NewButton("Quick snapshot", func() {}))
}