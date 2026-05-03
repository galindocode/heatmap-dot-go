package heatmap

import (
	"image"
	"image/color"
	"os"
	"testing"
)

// Test backward compatibility with legacy API
func TestWithoutImage(t *testing.T) {
	hm := New(1000, 1000)
	hm.SetPointSize(50)

	hm.AddData(100, 150, 10)
	hm.AddData(100, 200, 9)
	hm.AddData(200, 100, 8)
	hm.AddData(100, 300, 7)
	hm.AddData(300, 100, 6)
	hm.AddData(200, 200, 5)
	hm.AddData(300, 100, 3)
	hm.AddData(200, 400, 2)
	hm.AddData(500, 500, 3)
	hm.AddData(900, 900, 1)
	hm.SetMaxValue(10)

	img, err := hm.PNGImage()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if img == nil {
		t.Errorf("Expected image to be generated, got nil")
	}

	// Save image to file
	f, err := os.Create("test_without_image.png")
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	defer f.Close()
	f.Write(img)
}

func TestNew(t *testing.T) {
	hm := New(800, 600)

	if hm == nil {
		t.Fatal("New() returned nil")
	}

	if hm.config.Width != 800 {
		t.Errorf("Expected width 800, got %d", hm.config.Width)
	}

	if hm.config.Height != 600 {
		t.Errorf("Expected height 600, got %d", hm.config.Height)
	}

	if hm.config.PointSize != 10 {
		t.Errorf("Expected default point size 10, got %d", hm.config.PointSize)
	}

	if hm.config.Alpha != 180 {
		t.Errorf("Expected default alpha 180, got %d", hm.config.Alpha)
	}
}

func TestAddPoint(t *testing.T) {
	hm := New(800, 600)

	hm.AddPoint(100, 100, 0.5)
	hm.AddPoint(200, 150, 1.0)

	if hm.PointCount() != 2 {
		t.Errorf("Expected 2 points, got %d", hm.PointCount())
	}

	if hm.points[0].X != 100 || hm.points[0].Y != 100 {
		t.Errorf("First point coordinates incorrect")
	}

	if hm.points[1].Value != 1.0 {
		t.Errorf("Second point value incorrect")
	}
}

func TestSetters(t *testing.T) {
	hm := New(800, 600)

	hm.SetMaxValue(100)
	if hm.config.MaxValue == nil || *hm.config.MaxValue != 100 {
		t.Errorf("SetMaxValue failed")
	}

	hm.SetAlpha(200)
	if hm.config.Alpha != 200 {
		t.Errorf("SetAlpha failed")
	}

	hm.SetPointSize(20)
	if hm.config.PointSize != 20 {
		t.Errorf("SetPointSize failed")
	}

	colors := []string{"#ff0000", "#00ff00", "#0000ff"}
	hm.SetColorScheme(colors)
	if len(hm.config.ColorScheme) != 3 {
		t.Errorf("SetColorScheme failed")
	}
}

func TestClear(t *testing.T) {
	hm := New(800, 600)
	hm.AddPoint(100, 100, 0.5)
	hm.AddPoint(200, 200, 1.0)

	if hm.PointCount() != 2 {
		t.Errorf("Expected 2 points before clear")
	}

	hm.Clear()

	if hm.PointCount() != 0 {
		t.Errorf("Expected 0 points after clear, got %d", hm.PointCount())
	}
}

func TestValidation(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *Heatmap
		wantErr error
	}{
		{
			name: "valid heatmap",
			setup: func() *Heatmap {
				hm := New(800, 600)
				hm.AddPoint(100, 100, 1.0)
				return hm
			},
			wantErr: nil,
		},
		{
			name: "no data points",
			setup: func() *Heatmap {
				return New(800, 600)
			},
			wantErr: ErrNoData,
		},
		{
			name: "invalid size",
			setup: func() *Heatmap {
				hm := New(0, 0)
				hm.AddPoint(100, 100, 1.0)
				return hm
			},
			wantErr: ErrInvalidSize,
		},
		{
			name: "invalid point size",
			setup: func() *Heatmap {
				hm := New(800, 600)
				hm.SetPointSize(0)
				hm.AddPoint(100, 100, 1.0)
				return hm
			},
			wantErr: ErrInvalidPointSize,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hm := tt.setup()
			err := hm.validate()

			if err != tt.wantErr {
				t.Errorf("Expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestGenerate(t *testing.T) {
	hm := New(100, 100)
	hm.AddPoint(50, 50, 1.0)
	hm.AddPoint(25, 75, 0.5)

	img, err := hm.Generate()
	if err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}

	if img == nil {
		t.Fatal("Generate() returned nil image")
	}

	bounds := img.Bounds()
	if bounds.Dx() != 100 || bounds.Dy() != 100 {
		t.Errorf("Image dimensions incorrect: got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestGeneratePNG(t *testing.T) {
	hm := New(100, 100)
	hm.AddPoint(50, 50, 1.0)

	pngBytes, err := hm.GeneratePNG()
	if err != nil {
		t.Fatalf("GeneratePNG() failed: %v", err)
	}

	if len(pngBytes) == 0 {
		t.Error("GeneratePNG() returned empty bytes")
	}

	// Check PNG header
	if len(pngBytes) < 8 || pngBytes[0] != 0x89 || pngBytes[1] != 'P' {
		t.Error("GeneratePNG() did not return valid PNG data")
	}
}

func TestGenerateOverlay(t *testing.T) {
	// Create a base image
	baseImg := image.NewRGBA(image.Rect(0, 0, 100, 100))
	// Fill with white
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			baseImg.Set(x, y, color.White)
		}
	}

	hm := New(100, 100)
	hm.AddPoint(50, 50, 1.0)

	overlayImg, err := hm.GenerateOverlay(baseImg)
	if err != nil {
		t.Fatalf("GenerateOverlay() failed: %v", err)
	}

	if overlayImg == nil {
		t.Fatal("GenerateOverlay() returned nil")
	}

	bounds := overlayImg.Bounds()
	if bounds.Dx() != 100 || bounds.Dy() != 100 {
		t.Errorf("Overlay image dimensions incorrect")
	}
}

func TestGenerateOverlay_SizeMismatch(t *testing.T) {
	// Create a base image with different size
	baseImg := image.NewRGBA(image.Rect(0, 0, 200, 200))

	hm := New(100, 100)
	hm.AddPoint(50, 50, 1.0)

	_, err := hm.GenerateOverlay(baseImg)
	if err != ErrSizeMismatch {
		t.Errorf("Expected ErrSizeMismatch, got %v", err)
	}
}

func TestGenerateOverlay_NilImage(t *testing.T) {
	hm := New(100, 100)
	hm.AddPoint(50, 50, 1.0)

	_, err := hm.GenerateOverlay(nil)
	if err != ErrInvalidImage {
		t.Errorf("Expected ErrInvalidImage, got %v", err)
	}
}

func TestGetMaxValue(t *testing.T) {
	hm := New(100, 100)
	hm.AddPoint(10, 10, 0.5)
	hm.AddPoint(20, 20, 1.5)
	hm.AddPoint(30, 30, 1.0)

	maxVal := hm.getMaxValue()
	if maxVal != 1.5 {
		t.Errorf("Expected max value 1.5, got %f", maxVal)
	}
}

func TestGetMaxValue_WithConfigured(t *testing.T) {
	hm := New(100, 100)
	hm.SetMaxValue(10.0)
	hm.AddPoint(10, 10, 5.0)

	maxVal := hm.getMaxValue()
	if maxVal != 10.0 {
		t.Errorf("Expected configured max value 10.0, got %f", maxVal)
	}
}

func TestIsInBounds(t *testing.T) {
	hm := New(100, 100)

	tests := []struct {
		x, y     int
		expected bool
	}{
		{50, 50, true},
		{0, 0, true},
		{99, 99, true},
		{100, 100, false},
		{-1, 50, false},
		{50, -1, false},
		{150, 50, false},
	}

	for _, tt := range tests {
		result := hm.isInBounds(tt.x, tt.y)
		if result != tt.expected {
			t.Errorf("isInBounds(%d, %d) = %v, want %v", tt.x, tt.y, result, tt.expected)
		}
	}
}

func TestBlendAdditive(t *testing.T) {
	dst := color.RGBA{100, 100, 100, 200}
	src := color.RGBA{50, 50, 50, 100}

	result := blendAdditive(dst, src)

	if result.R != 150 || result.G != 150 || result.B != 150 {
		t.Errorf("Additive blend incorrect: got %v", result)
	}

	if result.A != 200 {
		t.Errorf("Additive blend alpha incorrect: got %d", result.A)
	}
}

func TestBlendAdditive_Clamping(t *testing.T) {
	dst := color.RGBA{200, 200, 200, 255}
	src := color.RGBA{100, 100, 100, 255}

	result := blendAdditive(dst, src)

	// Should clamp to 255
	if result.R != 255 || result.G != 255 || result.B != 255 {
		t.Errorf("Additive blend should clamp to 255: got %v", result)
	}
}

func TestCompositeNormal(t *testing.T) {
	base := image.NewRGBA(image.Rect(0, 0, 10, 10))
	overlay := image.NewRGBA(image.Rect(0, 0, 10, 10))

	// Fill base with white
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			base.Set(x, y, color.RGBA{255, 255, 255, 255})
			overlay.Set(x, y, color.RGBA{255, 0, 0, 128}) // 50% transparent red
		}
	}

	result := compositeNormal(base, overlay)

	if result == nil {
		t.Fatal("compositeNormal returned nil")
	}

	// Check that blending occurred
	c := result.RGBAAt(5, 5)
	if c.R == 255 && c.G == 255 && c.B == 255 {
		t.Error("Expected blending to occur, but got pure white")
	}
}

func BenchmarkGenerate(b *testing.B) {
	hm := New(800, 600)

	// Add 1000 points
	for i := 0; i < 1000; i++ {
		hm.AddPoint(float64(i%800), float64(i%600), float64(i%100)/100.0)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := hm.Generate()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAddPoint(b *testing.B) {
	hm := New(800, 600)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hm.AddPoint(100, 100, 0.5)
	}
}
