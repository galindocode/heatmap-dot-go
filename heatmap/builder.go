package heatmap

import "image/color"

// Builder provides a fluent interface for constructing Heatmap instances
// with complex configurations.
//
// Example:
//
//	hm, err := heatmap.NewBuilder().
//	    Size(800, 600).
//	    MaxValue(100).
//	    ColorScheme([]string{"#0000ff", "#00ff00", "#ffff00", "#ff0000"}).
//	    PointSize(15).
//	    Alpha(200).
//	    AddPoint(100, 100, 50).
//	    AddPoint(200, 150, 75).
//	    Build()
type Builder struct {
	config *Config
	points []Point
	errors []error
}

// NewBuilder creates a new Builder with default configuration.
func NewBuilder() *Builder {
	return &Builder{
		config: newDefaultConfig(),
		points: make([]Point, 0),
		errors: make([]error, 0),
	}
}

// Size sets the dimensions of the heatmap in pixels.
// Both width and height must be greater than 0.
func (b *Builder) Size(width, height int) *Builder {
	b.config.Width = width
	b.config.Height = height
	return b
}

// Width sets the width of the heatmap in pixels.
func (b *Builder) Width(width int) *Builder {
	b.config.Width = width
	return b
}

// Height sets the height of the heatmap in pixels.
func (b *Builder) Height(height int) *Builder {
	b.config.Height = height
	return b
}

// MaxValue sets the maximum value for gradient normalization.
// If not set, the maximum will be auto-calculated from the data.
func (b *Builder) MaxValue(maxValue float64) *Builder {
	b.config.MaxValue = &maxValue
	return b
}

// ColorScheme sets the color gradient using hex color strings.
// Accepts formats: #RGB, #RRGGBB, or #RRGGBBAA
//
// Example:
//
//	builder.ColorScheme([]string{"#3b82f6", "#22c55e", "#eab308", "#ef4444"})
func (b *Builder) ColorScheme(colors []string) *Builder {
	b.config.ColorScheme = colors
	return b
}

// Colors is an alias for ColorScheme that accepts variadic arguments.
//
// Example:
//
//	builder.Colors("#0000ff", "#00ff00", "#ff0000")
func (b *Builder) Colors(colors ...string) *Builder {
	b.config.ColorScheme = colors
	return b
}

// PointSize sets the radius of each point in pixels.
func (b *Builder) PointSize(size int) *Builder {
	b.config.PointSize = size
	return b
}

// Alpha sets the transparency level for the heatmap (0-255).
// 0 is fully transparent, 255 is fully opaque.
// Default is 180 (~70% opacity).
func (b *Builder) Alpha(alpha uint8) *Builder {
	b.config.Alpha = alpha
	return b
}

// Background sets the background color.
// If not called, the background will be transparent (default).
//
// Example:
//
//	builder.Background(color.White)
func (b *Builder) Background(bg color.Color) *Builder {
	b.config.Background = &bg
	return b
}

// TransparentBackground explicitly sets the background to transparent.
// This is the default behavior, so calling this is optional.
func (b *Builder) TransparentBackground() *Builder {
	b.config.Background = nil
	return b
}

// Gaussian enables or disables Gaussian falloff rendering.
// When true, each point is rendered as a smooth blob instead of a hard circle,
// producing thermal-camera-style heatmaps suitable for overlay on real images.
func (b *Builder) Gaussian(enabled bool) *Builder {
	b.config.GaussianMode = enabled
	return b
}

// AddPoint adds a single data point to the heatmap.
// Coordinates can be fractional for sub-pixel precision.
func (b *Builder) AddPoint(x, y, value float64) *Builder {
	b.points = append(b.points, Point{
		X:     x,
		Y:     y,
		Value: value,
	})
	return b
}

// AddPoints adds multiple data points to the heatmap.
func (b *Builder) AddPoints(points []Point) *Builder {
	b.points = append(b.points, points...)
	return b
}

// validate checks that the configuration is valid
func (b *Builder) validate() error {
	// Validate dimensions
	if b.config.Width <= 0 || b.config.Height <= 0 {
		return ErrInvalidSize
	}

	// Validate point size
	if b.config.PointSize <= 0 {
		return ErrInvalidPointSize
	}

	// Validate MaxValue if set
	if b.config.MaxValue != nil && *b.config.MaxValue <= 0 {
		return ErrInvalidMaxValue
	}

	// Validate color scheme
	if len(b.config.ColorScheme) == 0 {
		return ErrInvalidColor
	}

	// Validate that we have at least one data point
	// Note: We only validate this in Build(), not during construction
	// This allows building the heatmap first and adding points later

	return nil
}

// Build constructs the final Heatmap instance.
// Returns an error if the configuration is invalid.
func (b *Builder) Build() (*Heatmap, error) {
	if err := b.validate(); err != nil {
		return nil, err
	}

	// Create heatmap with copied config
	hm := &Heatmap{
		// Set legacy fields for backward compatibility
		Width:       b.config.Width,
		Height:      b.config.Height,
		ColorScheme: b.config.ColorScheme,
		PointsSize:  b.config.PointSize,
		MaxValue:    0, // Will be calculated or set from config

		// Set internal fields
		config: b.config.copy(),
		points: make([]Point, len(b.points)),
	}

	// Copy MaxValue if set
	if b.config.MaxValue != nil {
		hm.MaxValue = *b.config.MaxValue
	}

	// Copy points
	copy(hm.points, b.points)

	return hm, nil
}

// MustBuild is like Build but panics if there's an error.
// This is useful for static configurations that are known to be valid.
//
// Example:
//
//	hm := heatmap.NewBuilder().
//	    Size(800, 600).
//	    MustBuild()
func (b *Builder) MustBuild() *Heatmap {
	hm, err := b.Build()
	if err != nil {
		panic(err)
	}
	return hm
}
