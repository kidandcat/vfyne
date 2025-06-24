package testing

import (
	"flag"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
)

var updateSnapshots = flag.Bool("update-snapshots", false, "Update snapshot images")

type VFyneTest struct {
	t              *testing.T
	app            fyne.App
	window         fyne.Window
	snapshotDir    string
	screenshotDir  string
	renderWait     time.Duration
}

func New(t *testing.T) *VFyneTest {
	t.Helper()
	
	testDir := filepath.Dir(t.Name())
	if testDir == "." {
		testDir = "testdata"
	}
	
	return &VFyneTest{
		t:             t,
		app:           test.NewApp(),
		snapshotDir:   filepath.Join(testDir, "snapshots"),
		screenshotDir: filepath.Join(testDir, "screenshots"),
		renderWait:    100 * time.Millisecond,
	}
}

func (v *VFyneTest) SetTheme(theme fyne.Theme) {
	v.app.Settings().SetTheme(theme)
}

func (v *VFyneTest) SetRenderWait(duration time.Duration) {
	v.renderWait = duration
}

func (v *VFyneTest) Screenshot(name string, content fyne.CanvasObject, opts ...ScreenshotOption) {
	v.t.Helper()
	
	options := &screenshotOptions{
		size: fyne.NewSize(800, 600),
	}
	
	for _, opt := range opts {
		opt(options)
	}
	
	v.window = test.NewWindow(content)
	v.window.Resize(options.size)
	
	// Wait for rendering
	time.Sleep(v.renderWait)
	
	// Capture the canvas
	canvas := v.window.Canvas()
	img := canvas.Capture()
	
	filename := sanitizeFilename(name) + ".png"
	path := filepath.Join(v.screenshotDir, filename)
	
	if err := os.MkdirAll(v.screenshotDir, 0755); err != nil {
		v.t.Fatalf("Failed to create screenshot directory: %v", err)
	}
	
	if err := saveImage(path, img); err != nil {
		v.t.Fatalf("Failed to save screenshot: %v", err)
	}
	
	v.t.Logf("Screenshot saved: %s", path)
	
	v.window.Close()
}

func (v *VFyneTest) Snapshot(name string, content fyne.CanvasObject, opts ...ScreenshotOption) {
	v.t.Helper()
	
	options := &screenshotOptions{
		size: fyne.NewSize(800, 600),
	}
	
	for _, opt := range opts {
		opt(options)
	}
	
	v.window = test.NewWindow(content)
	v.window.Resize(options.size)
	
	// Wait for rendering
	time.Sleep(v.renderWait)
	
	// Capture the canvas
	canvas := v.window.Canvas()
	img := canvas.Capture()
	
	filename := sanitizeFilename(name) + ".png"
	snapshotPath := filepath.Join(v.snapshotDir, filename)
	
	if *updateSnapshots {
		if err := os.MkdirAll(v.snapshotDir, 0755); err != nil {
			v.t.Fatalf("Failed to create snapshot directory: %v", err)
		}
		
		if err := saveImage(snapshotPath, img); err != nil {
			v.t.Fatalf("Failed to save snapshot: %v", err)
		}
		
		v.t.Logf("Snapshot updated: %s", snapshotPath)
	} else {
		if _, err := os.Stat(snapshotPath); os.IsNotExist(err) {
			v.t.Errorf("Snapshot does not exist: %s (run with -update-snapshots to create)", snapshotPath)
			
			tempPath := filepath.Join(v.screenshotDir, "failed_"+filename)
			if err := os.MkdirAll(v.screenshotDir, 0755); err == nil {
				saveImage(tempPath, img)
				v.t.Logf("Actual output saved to: %s", tempPath)
			}
		} else {
			expected, err := loadImage(snapshotPath)
			if err != nil {
				v.t.Fatalf("Failed to load snapshot: %v", err)
			}
			
			if !imagesEqual(expected, img) {
				v.t.Errorf("Snapshot mismatch for %s", name)
				
				diffPath := filepath.Join(v.screenshotDir, "diff_"+filename)
				actualPath := filepath.Join(v.screenshotDir, "actual_"+filename)
				
				if err := os.MkdirAll(v.screenshotDir, 0755); err == nil {
					saveImage(actualPath, img)
					if diff := createDiffImage(expected, img); diff != nil {
						saveImage(diffPath, diff)
						v.t.Logf("Diff saved to: %s", diffPath)
					}
					v.t.Logf("Actual output saved to: %s", actualPath)
				}
			} else {
				v.t.Logf("Snapshot matched: %s", name)
			}
		}
	}
	
	v.window.Close()
}

type screenshotOptions struct {
	size fyne.Size
}

type ScreenshotOption func(*screenshotOptions)

func WithSize(width, height float32) ScreenshotOption {
	return func(o *screenshotOptions) {
		o.size = fyne.NewSize(width, height)
	}
}

func WithMobileSize() ScreenshotOption {
	return func(o *screenshotOptions) {
		o.size = fyne.NewSize(375, 667)
	}
}

func WithTabletSize() ScreenshotOption {
	return func(o *screenshotOptions) {
		o.size = fyne.NewSize(768, 1024)
	}
}

func sanitizeFilename(name string) string {
	reg := regexp.MustCompile(`[^a-zA-Z0-9_-]+`)
	sanitized := reg.ReplaceAllString(name, "_")
	sanitized = strings.Trim(sanitized, "_")
	return strings.ToLower(sanitized)
}

func saveImage(path string, img image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	
	return png.Encode(file, img)
}

func loadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	return png.Decode(file)
}

func imagesEqual(a, b image.Image) bool {
	if a.Bounds() != b.Bounds() {
		return false
	}
	
	bounds := a.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if a.At(x, y) != b.At(x, y) {
				return false
			}
		}
	}
	
	return true
}

func createDiffImage(expected, actual image.Image) image.Image {
	bounds := expected.Bounds()
	if bounds != actual.Bounds() {
		return nil
	}
	
	diff := image.NewRGBA(bounds)
	
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			e := expected.At(x, y)
			a := actual.At(x, y)
			
			if e != a {
				diff.Set(x, y, theme.ErrorColor())
			} else {
				diff.Set(x, y, e)
			}
		}
	}
	
	return diff
}

func AssertScreenshot(t *testing.T, name string, content fyne.CanvasObject, opts ...ScreenshotOption) {
	t.Helper()
	vt := New(t)
	vt.Screenshot(name, content, opts...)
}

func AssertSnapshot(t *testing.T, name string, content fyne.CanvasObject, opts ...ScreenshotOption) {
	t.Helper()
	vt := New(t)
	vt.Snapshot(name, content, opts...)
}