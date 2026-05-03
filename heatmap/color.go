package heatmap

import (
	"image/color"
	"strconv"
	"strings"
)

// ColorGradient represents a color gradient with customizable maximum value.
// It maps values to colors using linear interpolation between gradient stops.
type ColorGradient struct {
	// Colors are the gradient colors in hexadecimal format (#RGB, #RRGGBB, or #RRGGBBAA)
	Colors []string

	// MaxValue is the maximum value that corresponds to 100% of the gradient
	MaxValue float64
}

// GetColor returns the color corresponding to a specific value with the given alpha.
// The value is normalized against MaxValue to determine position in the gradient.
// Alpha should be in range 0-255 (0 = fully transparent, 255 = fully opaque).
func (cg *ColorGradient) GetColor(value float64, alpha uint8) color.Color {
	if len(cg.Colors) == 0 {
		return color.RGBA{0, 0, 0, alpha}
	}

	// Handle boundary cases
	if value <= 0 {
		return cg.hexToColorWithAlpha(cg.Colors[0], alpha)
	}
	if value >= cg.MaxValue {
		return cg.hexToColorWithAlpha(cg.Colors[len(cg.Colors)-1], alpha)
	}

	// Calculate relative percentage (0-100)
	percent := (value / cg.MaxValue) * 100.0

	// Determine which gradient segment we're in
	segmentSize := 100.0 / float64(len(cg.Colors)-1)
	segment := int(percent / segmentSize)

	// Prevent index out of bounds
	if segment >= len(cg.Colors)-1 {
		segment = len(cg.Colors) - 2
	}

	// Calculate position within the segment (0.0 to 1.0)
	segmentPercent := (percent - float64(segment)*segmentSize) / segmentSize

	// Get colors for current segment
	c1 := cg.hexToRGBA(cg.Colors[segment])
	c2 := cg.hexToRGBA(cg.Colors[segment+1])

	// Interpolate between colors
	r := cg.interpolate(c1.R, c2.R, segmentPercent)
	g := cg.interpolate(c1.G, c2.G, segmentPercent)
	b := cg.interpolate(c1.B, c2.B, segmentPercent)

	return color.RGBA{uint8(r), uint8(g), uint8(b), alpha}
}

// hexToColorWithAlpha converts hexadecimal color to color.Color with specified alpha
func (cg *ColorGradient) hexToColorWithAlpha(hex string, alpha uint8) color.Color {
	rgba := cg.hexToRGBA(hex)
	return color.RGBA{rgba.R, rgba.G, rgba.B, alpha}
}

// hexToRGBA converts hexadecimal color string to RGBA components.
// Supports formats: #RGB, #RRGGBB, and #RRGGBBAA
func (cg *ColorGradient) hexToRGBA(hex string) color.RGBA {
	hex = strings.TrimPrefix(hex, "#")

	// Default alpha to 255 (fully opaque)
	a := uint8(255)

	// Handle different hex formats
	switch len(hex) {
	case 3: // #RGB - short form
		r, _ := strconv.ParseInt(hex[0:1]+hex[0:1], 16, 64)
		g, _ := strconv.ParseInt(hex[1:2]+hex[1:2], 16, 64)
		b, _ := strconv.ParseInt(hex[2:3]+hex[2:3], 16, 64)
		return color.RGBA{uint8(r), uint8(g), uint8(b), a}

	case 6: // #RRGGBB - standard form
		r, _ := strconv.ParseInt(hex[0:2], 16, 64)
		g, _ := strconv.ParseInt(hex[2:4], 16, 64)
		b, _ := strconv.ParseInt(hex[4:6], 16, 64)
		return color.RGBA{uint8(r), uint8(g), uint8(b), a}

	case 8: // #RRGGBBAA - with alpha
		r, _ := strconv.ParseInt(hex[0:2], 16, 64)
		g, _ := strconv.ParseInt(hex[2:4], 16, 64)
		b, _ := strconv.ParseInt(hex[4:6], 16, 64)
		a, _ := strconv.ParseInt(hex[6:8], 16, 64)
		return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}

	default:
		// Return black on error
		return color.RGBA{0, 0, 0, 255}
	}
}

// interpolate performs linear interpolation between two uint8 values
// percent should be in range 0.0 to 1.0
func (cg *ColorGradient) interpolate(a, b uint8, percent float64) float64 {
	return float64(a) + (float64(b)-float64(a))*percent
}
