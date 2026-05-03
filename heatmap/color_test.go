package heatmap

import (
	"image/color"
	"testing"
)

func TestColorGradient_GetColor(t *testing.T) {
	cg := ColorGradient{
		Colors:   []string{"#0000ff", "#ff0000"},
		MaxValue: 100,
	}

	// Test at minimum
	c := cg.GetColor(0, 255)
	rgba := color.RGBAModel.Convert(c).(color.RGBA)
	if rgba.R != 0 || rgba.G != 0 || rgba.B != 255 {
		t.Errorf("Expected blue at min value, got %v", rgba)
	}

	// Test at maximum
	c = cg.GetColor(100, 255)
	rgba = color.RGBAModel.Convert(c).(color.RGBA)
	if rgba.R != 255 || rgba.G != 0 || rgba.B != 0 {
		t.Errorf("Expected red at max value, got %v", rgba)
	}

	// Test alpha
	c = cg.GetColor(50, 128)
	rgba = color.RGBAModel.Convert(c).(color.RGBA)
	if rgba.A != 128 {
		t.Errorf("Expected alpha 128, got %d", rgba.A)
	}
}

func TestColorGradient_HexParsing(t *testing.T) {
	tests := []struct {
		name     string
		hex      string
		expected color.RGBA
	}{
		{"RGB format", "#f00", color.RGBA{255, 0, 0, 255}},
		{"RRGGBB format", "#ff0000", color.RGBA{255, 0, 0, 255}},
		{"RRGGBBAA format", "#ff000080", color.RGBA{255, 0, 0, 128}},
		{"Blue", "#0000ff", color.RGBA{0, 0, 255, 255}},
		{"Green", "#00ff00", color.RGBA{0, 255, 0, 255}},
	}

	cg := ColorGradient{MaxValue: 100}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cg.hexToRGBA(tt.hex)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestColorGradient_Interpolation(t *testing.T) {
	cg := ColorGradient{
		Colors:   []string{"#000000", "#ffffff"},
		MaxValue: 100,
	}

	// Middle value should be gray
	c := cg.GetColor(50, 255)
	rgba := color.RGBAModel.Convert(c).(color.RGBA)

	// Should be approximately gray (127, 127, 127)
	if rgba.R < 120 || rgba.R > 135 {
		t.Errorf("Expected gray around 127, got R=%d", rgba.R)
	}
}
