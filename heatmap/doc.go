/*
Package heatmap provides a powerful and flexible library for generating heatmap visualizations in Go.

# Overview

Heatmap Dot Go allows you to create beautiful heatmap visualizations either as standalone images
or as overlays on existing images. It provides a dual API: a simple constructor for quick usage
and a fluent builder pattern for advanced configurations.

# Quick Start

The simplest way to create a heatmap:

	hm := heatmap.New(800, 600)
	hm.AddPoint(400, 300, 1.0)
	hm.AddPoint(200, 150, 0.5)
	hm.SavePNG("heatmap.png")

# Builder Pattern

For more control over the heatmap configuration:

	hm, err := heatmap.NewBuilder().
		Size(1920, 1080).
		MaxValue(100).
		Colors("#3b82f6", "#22c55e", "#eab308", "#ef4444").
		PointSize(20).
		Alpha(200).
		AddPoint(960, 540, 75).
		Build()

	if err != nil {
		log.Fatal(err)
	}

	err = hm.SavePNG("heatmap.png")

# Configuration Options

The library provides extensive configuration options:

Color Gradients: Support for hex colors in #RGB, #RRGGBB, or #RRGGBBAA format.
Default gradient is Blue → Green → Yellow → Red.

Point Size: Control the radius of each data point (default: 10 pixels)

Alpha Transparency: Set transparency from 0 (fully transparent) to 255 (fully opaque).
Default is 180 (~70% opaque)

Max Value: Control value normalization. Auto-calculated if not specified.

Background: Set a solid background color or use transparent (default)

# Standalone Heatmaps

Generate a heatmap as an independent image:

	hm := heatmap.New(800, 600)

	// Add data points
	for i := 0; i < 100; i++ {
		x := float64(rand.Intn(800))
		y := float64(rand.Intn(600))
		value := rand.Float64()
		hm.AddPoint(x, y, value)
	}

	// Generate and save
	err := hm.SavePNG("output.png")

# Overlay Heatmaps

Overlay a heatmap on an existing image:

	// Load base image
	baseImg, err := heatmap.LoadImage("map.png")
	if err != nil {
		log.Fatal(err)
	}

	// Create heatmap matching image size
	bounds := baseImg.Bounds()
	hm := heatmap.New(bounds.Dx(), bounds.Dy())
	hm.SetAlpha(150) // More transparent for overlay

	// Add data points
	hm.AddPoint(100, 100, 1.0)
	hm.AddPoint(200, 200, 0.8)

	// Generate overlay
	overlayImg, err := hm.GenerateOverlay(baseImg)
	if err != nil {
		log.Fatal(err)
	}

# Custom Color Gradients

Create custom color gradients with multiple colors:

	hm, _ := heatmap.NewBuilder().
		Size(800, 600).
		Colors("#000080", "#0000ff", "#00ffff", "#00ff00", "#ffff00", "#ff0000").
		PointSize(30).
		Build()

Colors are interpolated linearly between the specified stops.

# Data Points

Data points support sub-pixel precision using float64 coordinates:

	hm.AddPoint(100.5, 200.7, 0.8)

Points outside the heatmap bounds are silently discarded during rendering.

# Error Handling

The library provides specific error types for different failure scenarios:

	hm := heatmap.New(0, 0) // Invalid size
	_, err := hm.Generate()
	if errors.Is(err, heatmap.ErrInvalidSize) {
		// Handle invalid size error
	}

Common errors include:
  - ErrInvalidSize: Width or height is <= 0
  - ErrNoData: No data points have been added
  - ErrInvalidPointSize: Point size is <= 0
  - ErrInvalidMaxValue: Max value is <= 0
  - ErrSizeMismatch: Overlay image size doesn't match heatmap
  - ErrInvalidImage: Base image is nil or invalid

# Performance

The library is optimized for performance:
  - Efficient circle rendering with bounding box optimization
  - Additive color blending for overlapping points
  - Minimal memory allocations
  - Suitable for real-time applications

Typical performance:
  - 1,000 points @ 800x600: ~50ms
  - 10,000 points @ 1920x1080: ~200ms

# Thread Safety

Heatmap instances are not thread-safe. If you need to access a heatmap
from multiple goroutines, you must synchronize access yourself.

# Examples

For more examples, see the repository's main.go file or the README.md.
*/
package heatmap
