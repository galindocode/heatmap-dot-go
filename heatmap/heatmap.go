package heatmap

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
)

// Heatmap represents a heatmap visualization with configuration and data points.
// It can generate standalone heatmap images or overlay heatmaps on existing images.
type Heatmap struct {
	// Deprecated: Use builder methods instead. Will be removed in v1.0
	// MaxValue is the maximum value for color gradient normalization
	MaxValue float64

	// Deprecated: Use builder methods instead. Will be removed in v1.0
	// Width is the width of the heatmap in pixels
	Width int

	// Deprecated: Use builder methods instead. Will be removed in v1.0
	// Height is the height of the heatmap in pixels
	Height int

	// Deprecated: Use AddPoint() method instead. Will be removed in v1.0
	// Data contains the heatmap data points
	Data []HeatmapData

	// Deprecated: Use builder methods instead. Will be removed in v1.0
	// ColorScheme defines the gradient colors
	ColorScheme []string

	// Deprecated: Use builder methods instead. Will be removed in v1.0
	// PointsSize is the radius of each point
	PointsSize int

	// Internal configuration (immutable after Build)
	config *Config

	// Internal data storage
	points []Point
}

// HeatmapData represents a single data point (legacy structure)
// Deprecated: Use Point struct instead. Will be removed in v1.0
type HeatmapData struct {
	X     int
	Y     int
	Value float64
}

// New creates a new Heatmap with the specified dimensions.
// This is the simple constructor for basic use cases.
//
// Example:
//
//	hm := heatmap.New(800, 600)
//	hm.AddPoint(400, 300, 1.0)
//	hm.SavePNG("output.png")
func New(width, height int) *Heatmap {
	cfg := newDefaultConfig()
	cfg.Width = width
	cfg.Height = height

	return &Heatmap{
		// Set legacy fields for backward compatibility
		Width:       width,
		Height:      height,
		ColorScheme: cfg.ColorScheme,
		PointsSize:  cfg.PointSize,
		MaxValue:    0, // Will be auto-calculated

		// Internal fields
		config: cfg,
		points: make([]Point, 0),
	}
}

// AddPoint adds a single data point to the heatmap.
// Coordinates can be fractional for sub-pixel precision.
// Points outside the heatmap bounds will be silently discarded during rendering.
func (h *Heatmap) AddPoint(x, y, value float64) {
	h.points = append(h.points, Point{
		X:     x,
		Y:     y,
		Value: value,
	})
}

// AddData adds a data point using the legacy HeatmapData structure.
// Deprecated: Use AddPoint() instead. Will be removed in v1.0
func (h *Heatmap) AddData(x, y int, value float64) {
	h.Data = append(h.Data, HeatmapData{X: x, Y: y, Value: value})
	h.AddPoint(float64(x), float64(y), value)
}

// SetMaxValue sets the maximum value for gradient normalization.
// If not set, the maximum will be auto-calculated from the data.
func (h *Heatmap) SetMaxValue(maxValue float64) {
	h.MaxValue = maxValue
	h.config.MaxValue = &maxValue
}

// SetSize sets the dimensions of the heatmap.
func (h *Heatmap) SetSize(width, height int) {
	h.Width = width
	h.Height = height
	h.config.Width = width
	h.config.Height = height
}

// SetAlpha sets the transparency level for the heatmap (0-255).
// 0 is fully transparent, 255 is fully opaque.
func (h *Heatmap) SetAlpha(alpha uint8) {
	h.config.Alpha = alpha
}

// SetPointSize sets the radius of each point in pixels.
func (h *Heatmap) SetPointSize(size int) {
	h.PointsSize = size
	h.config.PointSize = size
}

// SetColorScheme sets the color gradient using hex color strings.
// Accepts formats: #RGB, #RRGGBB, or #RRGGBBAA
func (h *Heatmap) SetColorScheme(colors []string) {
	h.ColorScheme = colors
	h.config.ColorScheme = colors
}

// SetBackground sets the background color.
// Pass nil for transparent background (default).
func (h *Heatmap) SetBackground(bg color.Color) {
	h.config.Background = &bg
}

// SetGaussianMode enables smooth Gaussian falloff rendering.
// When true, each point is rendered as a soft blob using a Gaussian kernel
// accumulated into a float64 buffer, producing smooth blending like a thermal camera.
// When false (default), points are rendered as hard filled circles.
func (h *Heatmap) SetGaussianMode(enabled bool) {
	h.config.GaussianMode = enabled
}

// Clear removes all data points from the heatmap.
func (h *Heatmap) Clear() {
	h.points = make([]Point, 0)
	h.Data = make([]HeatmapData, 0)
}

// PointCount returns the number of data points in the heatmap.
func (h *Heatmap) PointCount() int {
	return len(h.points)
}

// Generate creates and returns the heatmap as an image.Image.
// This gives you full control over the image for further processing.
func (h *Heatmap) Generate() (image.Image, error) {
	if err := h.validate(); err != nil {
		return nil, err
	}

	return h.generateImage()
}

// GeneratePNG generates the heatmap and returns it as PNG bytes.
// This is convenient for direct file writing or HTTP responses.
func (h *Heatmap) GeneratePNG() ([]byte, error) {
	img, err := h.Generate()
	if err != nil {
		return nil, err
	}

	return imageToPNGBytes(img)
}

// PNGImage generates the heatmap and returns it as PNG bytes.
// Deprecated: Use GeneratePNG() instead. Will be removed in v1.0
func (h *Heatmap) PNGImage() ([]byte, error) {
	return h.GeneratePNG()
}

// SavePNG generates the heatmap and saves it as a PNG file.
func (h *Heatmap) SavePNG(filepath string) error {
	pngBytes, err := h.GeneratePNG()
	if err != nil {
		return err
	}

	return saveBytesToFile(filepath, pngBytes)
}

// GenerateOverlay creates a heatmap overlay on top of the provided base image.
// The heatmap dimensions must match the base image dimensions.
// Returns the composited image with the heatmap overlaid on top.
func (h *Heatmap) GenerateOverlay(baseImage image.Image) (image.Image, error) {
	if baseImage == nil {
		return nil, ErrInvalidImage
	}

	// Validate that sizes match
	bounds := baseImage.Bounds()
	if bounds.Dx() != h.config.Width || bounds.Dy() != h.config.Height {
		return nil, ErrSizeMismatch
	}

	// Generate the heatmap on a transparent canvas
	heatmapImg, err := h.Generate()
	if err != nil {
		return nil, err
	}

	// Composite the heatmap over the base image
	return compositeNormal(baseImage, heatmapImg), nil
}

// GenerateOverlayPNG creates a heatmap overlay and returns it as PNG bytes.
func (h *Heatmap) GenerateOverlayPNG(baseImage image.Image) ([]byte, error) {
	img, err := h.GenerateOverlay(baseImage)
	if err != nil {
		return nil, err
	}

	return imageToPNGBytes(img)
}

// DrawOverlay creates a heatmap overlay on the provided image.
// Deprecated: Use GenerateOverlay() instead for better error handling.
func (h *Heatmap) DrawOverlay(img image.Image) (image.Image, error) {
	return h.GenerateOverlay(img)
}

// validate checks that the heatmap configuration and data are valid
func (h *Heatmap) validate() error {
	// Validate dimensions
	if h.config.Width <= 0 || h.config.Height <= 0 {
		return ErrInvalidSize
	}

	// Validate point size
	if h.config.PointSize <= 0 {
		return ErrInvalidPointSize
	}

	// Validate MaxValue if set
	if h.config.MaxValue != nil && *h.config.MaxValue <= 0 {
		return ErrInvalidMaxValue
	}

	// Validate that we have data
	if len(h.points) == 0 {
		return ErrNoData
	}

	// Validate color scheme
	if len(h.config.ColorScheme) == 0 {
		return ErrInvalidColor
	}

	return nil
}

// generateImage creates the actual heatmap image
func (h *Heatmap) generateImage() (*image.RGBA, error) {
	if h.config.GaussianMode {
		return h.generateGaussianImage()
	}

	// Create empty canvas
	canvas := image.NewRGBA(image.Rect(0, 0, h.config.Width, h.config.Height))

	// Fill background if specified
	if h.config.Background != nil {
		h.fillBackground(canvas, *h.config.Background)
	}

	// Calculate max value if not set
	maxValue := h.getMaxValue()

	// Create color gradient
	gradient := ColorGradient{
		Colors:   h.config.ColorScheme,
		MaxValue: maxValue,
	}

	// Render each point
	for _, point := range h.points {
		// Skip points outside bounds
		if !h.isInBounds(int(point.X), int(point.Y)) {
			continue
		}

		// Get color for this point's value
		pointColor := gradient.GetColor(point.Value, h.config.Alpha)

		// Draw the point
		h.drawCircle(canvas, int(point.X), int(point.Y), h.config.PointSize, pointColor)
	}

	return canvas, nil
}

// getMaxValue returns the maximum value for normalization
func (h *Heatmap) getMaxValue() float64 {
	// Use configured MaxValue if set
	if h.config.MaxValue != nil {
		return *h.config.MaxValue
	}

	// Auto-calculate from data
	maxVal := 0.0
	for _, p := range h.points {
		if p.Value > maxVal {
			maxVal = p.Value
		}
	}

	// Avoid division by zero
	if maxVal == 0 {
		maxVal = 1.0
	}

	return maxVal
}

// isInBounds checks if coordinates are within the heatmap bounds
func (h *Heatmap) isInBounds(x, y int) bool {
	return x >= 0 && x < h.config.Width && y >= 0 && y < h.config.Height
}

// drawCircle renders a filled circle at the specified position
func (h *Heatmap) drawCircle(img *image.RGBA, x, y, radius int, c color.Color) {
	rect := img.Bounds()
	r2 := radius * radius

	for dy := -radius; dy <= radius; dy++ {
		for dx := -radius; dx <= radius; dx++ {
			// Check if point is inside the circle
			if dx*dx+dy*dy <= r2 {
				nx := x + dx
				ny := y + dy

				// Check bounds
				if nx >= rect.Min.X && nx < rect.Max.X && ny >= rect.Min.Y && ny < rect.Max.Y {
					// Blend colors additively
					existing := img.RGBAAt(nx, ny)
					blended := blendAdditive(existing, c)
					img.Set(nx, ny, blended)
				}
			}
		}
	}
}

// fillBackground fills the entire canvas with the specified color
func (h *Heatmap) fillBackground(img *image.RGBA, bgColor color.Color) {
	rect := img.Bounds()
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			img.Set(x, y, bgColor)
		}
	}
}

// blendAdditive performs additive color blending
func blendAdditive(dst color.Color, src color.Color) color.RGBA {
	dstRGBA := color.RGBAModel.Convert(dst).(color.RGBA)
	srcRGBA := color.RGBAModel.Convert(src).(color.RGBA)

	return color.RGBA{
		R: clampUint8(int(dstRGBA.R) + int(srcRGBA.R)),
		G: clampUint8(int(dstRGBA.G) + int(srcRGBA.G)),
		B: clampUint8(int(dstRGBA.B) + int(srcRGBA.B)),
		A: maxUint8(dstRGBA.A, srcRGBA.A),
	}
}

// compositeNormal performs normal alpha blending to composite overlay on base
func compositeNormal(base, overlay image.Image) *image.RGBA {
	bounds := base.Bounds()
	result := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			baseColor := color.RGBAModel.Convert(base.At(x, y)).(color.RGBA)
			overlayColor := color.RGBAModel.Convert(overlay.At(x, y)).(color.RGBA)

			// Alpha blending formula: result = overlay * alpha + base * (1 - alpha)
			alpha := float64(overlayColor.A) / 255.0

			result.SetRGBA(x, y, color.RGBA{
				R: uint8(float64(overlayColor.R)*alpha + float64(baseColor.R)*(1-alpha)),
				G: uint8(float64(overlayColor.G)*alpha + float64(baseColor.G)*(1-alpha)),
				B: uint8(float64(overlayColor.B)*alpha + float64(baseColor.B)*(1-alpha)),
				A: maxUint8(baseColor.A, overlayColor.A),
			})
		}
	}

	return result
}

// imageToPNGBytes converts an image to PNG bytes
func imageToPNGBytes(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Helper functions

func clampUint8(val int) uint8 {
	if val < 0 {
		return 0
	}
	if val > 255 {
		return 255
	}
	return uint8(val)
}

func maxUint8(a, b uint8) uint8 {
	if a > b {
		return a
	}
	return b
}
