# Heatmap Dot Go

![](./img/banner.jpg)

A powerful and flexible Go library for generating heatmap visualizations. Create beautiful heatmaps as standalone images or overlay them on existing images with ease.

[![Go Reference](https://pkg.go.dev/badge/github.com/galindocode/heatmap-dot-go.svg)](https://pkg.go.dev/github.com/galindocode/heatmap-dot-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/galindocode/heatmap-dot-go)](https://goreportcard.com/report/github.com/galindocode/heatmap-dot-go)

## Features

- 🎨 **Flexible Color Gradients**: Use predefined color schemes or create custom gradients with hex colors
- 🏗️ **Dual API**: Simple constructor for quick use, Builder pattern for advanced configuration
- 🖼️ **Standalone or Overlay**: Generate heatmaps as independent images or overlay them on existing images
- ⚡ **High Performance**: Efficient rendering with optimized algorithms
- 🎯 **Sub-pixel Precision**: Support for fractional coordinates for smooth visualizations
- 🔧 **Highly Configurable**: Control point size, alpha transparency, max values, and more
- ✅ **Well Tested**: Comprehensive test suite with >85% coverage
- 📚 **Full Documentation**: Complete GoDoc documentation with examples

## Installation

```bash
go get github.com/galindocode/heatmap-dot-go
```

## Quick Start

### Simple API

```go
package main

import "github.com/galindocode/heatmap-dot-go/heatmap"

func main() {
    // Create a new heatmap
    hm := heatmap.New(800, 600)

    // Add data points (x, y, value)
    hm.AddPoint(400, 300, 1.0)
    hm.AddPoint(200, 150, 0.5)
    hm.AddPoint(600, 450, 0.8)

    // Save as PNG
    hm.SavePNG("heatmap.png")
}
```

### Builder API

For more control, use the fluent builder pattern:

```go
hm, err := heatmap.NewBuilder().
    Size(1920, 1080).
    MaxValue(100).
    Colors("#3b82f6", "#22c55e", "#eab308", "#ef4444").
    PointSize(20).
    Alpha(200).
    AddPoint(960, 540, 75).
    AddPoint(400, 300, 50).
    Build()

if err != nil {
    log.Fatal(err)
}

err = hm.SavePNG("heatmap.png")
```

## Usage Examples

### Custom Color Gradient

```go
hm, _ := heatmap.NewBuilder().
    Size(800, 600).
    Colors("#000080", "#0000ff", "#00ffff", "#00ff00", "#ffff00", "#ff0000").
    PointSize(30).
    Build()

// Add your data points
for _, dataPoint := range myData {
    hm.AddPoint(dataPoint.X, dataPoint.Y, dataPoint.Value)
}

hm.SavePNG("output.png")
```

### Overlay on Existing Image

```go
// Load base image
baseImg, err := heatmap.LoadImage("map.png")
if err != nil {
    log.Fatal(err)
}

// Create heatmap matching image size
bounds := baseImg.Bounds()
hm := heatmap.New(bounds.Dx(), bounds.Dy())
hm.SetAlpha(150) // More transparent for better visibility

// Add data points
hm.AddPoint(100, 100, 1.0)
hm.AddPoint(200, 200, 0.8)

// Generate overlay
overlayImg, err := hm.GenerateOverlay(baseImg)
if err != nil {
    log.Fatal(err)
}

// Save result
pngBytes, _ := heatmap.GeneratePNG(overlayImg)
os.WriteFile("overlay.png", pngBytes, 0644)
```

### Configuring Alpha Transparency

```go
hm := heatmap.New(800, 600)
hm.SetAlpha(180) // 0 = fully transparent, 255 = fully opaque

// Or use builder
hm, _ := heatmap.NewBuilder().
    Size(800, 600).
    Alpha(200).
    Build()
```

### Setting Max Value for Normalization

```go
hm := heatmap.New(800, 600)
hm.SetMaxValue(100) // All values will be normalized against this

// Or let it auto-calculate
hm := heatmap.New(800, 600)
// MaxValue will be automatically determined from your data
```

### Background Color

```go
import "image/color"

hm, _ := heatmap.NewBuilder().
    Size(800, 600).
    Background(color.White). // Set white background
    Build()

// Or transparent (default)
hm, _ := heatmap.NewBuilder().
    Size(800, 600).
    TransparentBackground().
    Build()
```

## API Reference

### Core Types

#### Heatmap

The main structure for creating heatmaps.

```go
type Heatmap struct {
    // Configuration and data
}
```

**Methods:**
- `New(width, height int) *Heatmap` - Create new heatmap with simple API
- `AddPoint(x, y, value float64)` - Add a data point
- `SetMaxValue(maxValue float64)` - Set max value for normalization
- `SetAlpha(alpha uint8)` - Set transparency (0-255)
- `SetPointSize(size int)` - Set point radius in pixels
- `SetColorScheme(colors []string)` - Set color gradient
- `SetBackground(bg color.Color)` - Set background color
- `Clear()` - Remove all data points
- `PointCount() int` - Get number of points
- `Generate() (image.Image, error)` - Generate heatmap image
- `GeneratePNG() ([]byte, error)` - Generate PNG bytes
- `SavePNG(filepath string) error` - Save as PNG file
- `GenerateOverlay(baseImage image.Image) (image.Image, error)` - Create overlay

#### Builder

Fluent interface for complex configurations.

```go
type Builder struct {
    // Builder state
}
```

**Methods:**
- `NewBuilder() *Builder` - Create new builder
- `Size(width, height int) *Builder` - Set dimensions
- `MaxValue(maxValue float64) *Builder` - Set max value
- `ColorScheme(colors []string) *Builder` - Set colors
- `Colors(colors ...string) *Builder` - Set colors (variadic)
- `PointSize(size int) *Builder` - Set point size
- `Alpha(alpha uint8) *Builder` - Set transparency
- `Background(bg color.Color) *Builder` - Set background
- `TransparentBackground() *Builder` - Set transparent background
- `AddPoint(x, y, value float64) *Builder` - Add point
- `AddPoints(points []Point) *Builder` - Add multiple points
- `Build() (*Heatmap, error)` - Build heatmap
- `MustBuild() *Heatmap` - Build or panic

#### Point

Represents a data point.

```go
type Point struct {
    X     float64 // X coordinate
    Y     float64 // Y coordinate
    Value float64 // Intensity/weight
}
```

### Error Types

```go
var (
    ErrInvalidSize      = errors.New("heatmap: invalid size")
    ErrNoData           = errors.New("heatmap: no data points")
    ErrInvalidMaxValue  = errors.New("heatmap: MaxValue must be > 0")
    ErrInvalidPointSize = errors.New("heatmap: PointSize must be > 0")
    ErrInvalidColor     = errors.New("heatmap: invalid color format")
    ErrInvalidImage     = errors.New("heatmap: invalid base image")
    ErrSizeMismatch     = errors.New("heatmap: size mismatch")
)
```

## Configuration Options

### Color Schemes

The library supports hex color formats:
- `#RGB` - Short form (e.g., `#f00` for red)
- `#RRGGBB` - Standard form (e.g., `#ff0000`)
- `#RRGGBBAA` - With alpha channel (e.g., `#ff0000ff`)

**Default gradient:** `["#3b82f6", "#22c55e", "#eab308", "#ef4444"]` (Blue → Green → Yellow → Red)

### Point Size

Default: `10` pixels

Controls the radius of each data point circle.

### Alpha Transparency

Default: `180` (~70% opaque)

Range: `0` (fully transparent) to `255` (fully opaque)

### Max Value

Default: `auto-calculated from data`

Used for normalizing values to the color gradient.

## Performance

The library is optimized for performance:

- Efficient circle rendering with bounding box optimization
- Additive color blending for overlapping points
- Minimal memory allocations
- Suitable for real-time applications

**Benchmarks** (on typical hardware):
- 1,000 points @ 800x600: ~50ms
- 10,000 points @ 1920x1080: ~200ms

## Examples

Check out the `main.go` file for complete working examples:

```bash
go run main.go
```

This will generate three example heatmaps:
1. `output_simple.png` - Basic usage with simple API
2. `output_builder.png` - Advanced usage with builder pattern
3. `output_custom.png` - Custom gradient and configuration

## Use Cases

- **Geographic Data**: Visualize density, temperature, or other spatial data
- **User Analytics**: Click heatmaps, gaze tracking, user interaction patterns
- **Scientific Visualization**: Temperature distributions, particle density
- **Business Intelligence**: Sales density, foot traffic analysis
- **Performance Monitoring**: Server latency, resource usage hotspots

## Development

### Running Tests

```bash
cd heatmap
go test -v
```

### Running Benchmarks

```bash
cd heatmap
go test -bench=. -benchmem
```

### Test Coverage

```bash
cd heatmap
go test -cover
```

## Design Documentation

For detailed architecture and design decisions, see the [design documentation](./design/README.md).

## Roadmap

See [ROADMAP.md](./design/06-migration-roadmap.md) for planned features and improvements.

### Upcoming Features

- 🔄 **Gaussian Renderer**: Smooth gradient rendering with blur
- 📊 **Density Renderer**: KDE-based visualization
- 🗺️ **Contour Lines**: Iso-line generation
- ⚡ **Parallel Rendering**: Multi-core optimization
- 🎨 **More Color Schemes**: Predefined palettes (Viridis, Plasma, etc.)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

[MIT License](LICENSE)

## Acknowledgments

Inspired by heatmap libraries in other languages and designed following Go best practices.

---

**Made with ❤️ by [@galindocode](https://github.com/galindocode)**
