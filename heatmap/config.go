package heatmap

import "image/color"

// Config holds the immutable configuration for a heatmap.
// Once a heatmap is built, its configuration cannot be changed.
type Config struct {
	// Width is the width of the heatmap in pixels
	Width int

	// Height is the height of the heatmap in pixels
	Height int

	// MaxValue is the maximum value for normalization.
	// If nil, it will be auto-calculated from the data.
	MaxValue *float64

	// ColorScheme defines the colors used for the gradient.
	// Uses hex color format: #RGB, #RRGGBB, or #RRGGBBAA
	ColorScheme []string

	// PointSize is the radius of each point in pixels
	PointSize int

	// Alpha is the transparency level (0-255)
	// 0 = fully transparent, 255 = fully opaque
	Alpha uint8

	// Background is the background color.
	// If nil, the background will be transparent.
	Background *color.Color

	// GaussianMode enables smooth Gaussian falloff rendering instead of hard circles.
	GaussianMode bool
}

// Point represents a single data point in the heatmap
type Point struct {
	// X coordinate (can be fractional for sub-pixel precision)
	X float64

	// Y coordinate (can be fractional for sub-pixel precision)
	Y float64

	// Value is the intensity/weight of this point
	Value float64
}

// newDefaultConfig creates a config with sensible defaults
func newDefaultConfig() *Config {
	return &Config{
		Width:       800,
		Height:      600,
		MaxValue:    nil, // Auto-calculate
		ColorScheme: []string{"#3b82f6", "#22c55e", "#eab308", "#ef4444"},
		PointSize:   10,
		Alpha:        180, // ~70% opacity
		Background:   nil, // Transparent
		GaussianMode: false,
	}
}

// copy creates a deep copy of the config
func (c *Config) copy() *Config {
	cfg := &Config{
		Width:       c.Width,
		Height:      c.Height,
		PointSize:   c.PointSize,
		Alpha:       c.Alpha,
		ColorScheme: make([]string, len(c.ColorScheme)),
	}

	copy(cfg.ColorScheme, c.ColorScheme)

	if c.MaxValue != nil {
		val := *c.MaxValue
		cfg.MaxValue = &val
	}

	if c.Background != nil {
		bg := *c.Background
		cfg.Background = &bg
	}

	cfg.GaussianMode = c.GaussianMode

	return cfg
}
