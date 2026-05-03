// Gaussian overlay on a real camera image.
// Run from the project root: go run ./examples/gaussian_overlay/
package main

import (
	"log"
	"os"

	"github.com/galindocode/heatmap-dot-go/heatmap"
)

func main() {
	base, err := heatmap.LoadImage("supermarket.jpg")
	if err != nil {
		log.Fatalf("load image: %v", err)
	}

	hm := heatmap.New(640, 424)
	hm.SetGaussianMode(true)
	hm.SetPointSize(100)
	hm.SetAlpha(230)
	hm.SetMaxValue(10.0)

	// Tracked person positions across the store.
	for _, p := range [][3]float64{
		// Checkout area — highest traffic
		{230, 375, 10.0}, {270, 370, 9.5}, {310, 372, 9.0},
		{350, 368, 9.5}, {390, 373, 8.5},
		// Entrance
		{310, 55, 7.0}, {340, 65, 6.5}, {280, 60, 6.0},
		// Aisles
		{160, 200, 5.0}, {160, 280, 4.5},
		{480, 180, 4.0}, {500, 260, 3.5},
		// Produce section
		{90, 150, 3.0}, {110, 220, 2.5},
	} {
		hm.AddPoint(p[0], p[1], p[2])
	}

	pngBytes, err := hm.GenerateOverlayPNG(base)
	if err != nil {
		log.Fatalf("overlay: %v", err)
	}
	if err := os.WriteFile("img/demo_gaussian.png", pngBytes, 0644); err != nil {
		log.Fatalf("save: %v", err)
	}
	log.Println("saved → img/demo_gaussian.png")
}
