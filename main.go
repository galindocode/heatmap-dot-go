package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"

	"github.com/galindocode/heatmap-dot-go/heatmap"
)

func main() {
	fmt.Println("Heatmap Dot Go - Example Usage")
	fmt.Println("================================")

	// Example 1: Simple API
	simpleExample()

	// Example 2: Builder API
	builderExample()

	// Example 3: Custom configuration
	customExample()

	// Example 4: Overlay on image
	// Commented out - requires base image file
	// overlayExample()

	fmt.Println("\nAll examples completed successfully!")
	fmt.Println("Check the output directory for generated images.")
}

// simpleExample demonstrates the basic usage with the simple API
func simpleExample() {
	fmt.Println("\n1. Simple API Example")
	fmt.Println("   Creating basic heatmap...")

	// Create heatmap with simple constructor
	hm := heatmap.New(800, 600)

	// Add some random data points
	for i := 0; i < 100; i++ {
		x := float64(rand.Intn(800))
		y := float64(rand.Intn(600))
		value := rand.Float64()
		hm.AddPoint(x, y, value)
	}

	// Generate and save
	err := hm.SavePNG("output_simple.png")
	if err != nil {
		log.Printf("   Error: %v\n", err)
		return
	}

	fmt.Println("   ✓ Saved to: output_simple.png")
}

// builderExample demonstrates the builder pattern for more control
func builderExample() {
	fmt.Println("\n2. Builder Pattern Example")
	fmt.Println("   Creating heatmap with builder...")

	// Use builder for fluent configuration
	hm, err := heatmap.NewBuilder().
		Size(1920, 1080).
		MaxValue(100).
		Colors("#3b82f6", "#22c55e", "#eab308", "#ef4444").
		PointSize(20).
		Alpha(200).
		Build()

	if err != nil {
		log.Printf("   Error building heatmap: %v\n", err)
		return
	}

	// Add clustered points to create hotspots
	addClusteredPoints(hm, 960, 540, 50, 200) // Center
	addClusteredPoints(hm, 400, 300, 30, 100) // Top-left
	addClusteredPoints(hm, 1520, 780, 40, 150) // Bottom-right

	// Save the result
	err = hm.SavePNG("output_builder.png")
	if err != nil {
		log.Printf("   Error: %v\n", err)
		return
	}

	fmt.Println("   ✓ Saved to: output_builder.png")
}

// customExample shows advanced customization options
func customExample() {
	fmt.Println("\n3. Custom Configuration Example")
	fmt.Println("   Creating heatmap with custom settings...")

	// Create with custom gradient and white background
	hm, err := heatmap.NewBuilder().
		Size(800, 600).
		MaxValue(1.0).
		Colors("#000080", "#0000ff", "#00ffff", "#00ff00", "#ffff00", "#ff0000").
		PointSize(30).
		Alpha(180).
		Background(color.White).
		Build()

	if err != nil {
		log.Printf("   Error: %v\n", err)
		return
	}

	// Add points in a pattern
	for x := 100; x < 700; x += 100 {
		for y := 100; y < 500; y += 100 {
			value := float64(x+y) / 1200.0
			hm.AddPoint(float64(x), float64(y), value)
		}
	}

	err = hm.SavePNG("output_custom.png")
	if err != nil {
		log.Printf("   Error: %v\n", err)
		return
	}

	fmt.Println("   ✓ Saved to: output_custom.png")
	fmt.Printf("   ✓ Total points: %d\n", hm.PointCount())
}

// overlayExample demonstrates overlaying a heatmap on an existing image
// Commented out by default - requires a base image file
/*
func overlayExample() {
	fmt.Println("\n4. Overlay Example")
	fmt.Println("   Creating heatmap overlay...")

	// Load base image
	baseImg, err := heatmap.LoadImage("base_image.png")
	if err != nil {
		fmt.Println("   ⚠ Skipping overlay example (base_image.png not found)")
		return
	}

	bounds := baseImg.Bounds()
	hm := heatmap.New(bounds.Dx(), bounds.Dy())
	hm.SetAlpha(150) // More transparent for overlay

	// Add points
	for i := 0; i < 200; i++ {
		x := float64(rand.Intn(bounds.Dx()))
		y := float64(rand.Intn(bounds.Dy()))
		value := rand.Float64()
		hm.AddPoint(x, y, value)
	}

	// Generate overlay
	pngBytes, err := hm.GenerateOverlayPNG(baseImg)
	if err != nil {
		log.Printf("   Error: %v\n", err)
		return
	}

	// Save to file
	err = os.WriteFile("output_overlay.png", pngBytes, 0644)
	if err != nil {
		log.Printf("   Error: %v\n", err)
		return
	}

	fmt.Println("   ✓ Saved to: output_overlay.png")
}
*/

// Helper function to add clustered points around a center
func addClusteredPoints(hm *heatmap.Heatmap, centerX, centerY float64, radius int, count int) {
	for i := 0; i < count; i++ {
		// Random offset within radius
		angle := rand.Float64() * 2 * 3.14159
		distance := rand.Float64() * float64(radius)

		x := centerX + distance*cos(angle)
		y := centerY + distance*sin(angle)

		// Value decreases with distance from center
		value := 1.0 - (distance / float64(radius))

		hm.AddPoint(x, y, value)
	}
}

func cos(angle float64) float64 {
	return math.Cos(angle)
}

func sin(angle float64) float64 {
	return math.Sin(angle)
}
