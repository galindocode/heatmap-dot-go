// Retail foot traffic — simulates person tracking over a store camera feed.
// Shows hotspots at checkout, entrance and aisles using the default gradient.
// Run from the project root: go run ./examples/foot_traffic/
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
	hm.SetPointSize(90)
	hm.SetAlpha(220)
	hm.SetMaxValue(10.0)

	// Each entry is [x, y, dwell_time_score].
	traffic := [][3]float64{
		// Checkout — longest dwell time
		{220, 375, 10.0}, {255, 370, 9.7}, {290, 373, 9.4},
		{325, 370, 9.7}, {360, 374, 9.1}, {395, 371, 8.7},
		// Customer service desk
		{575, 375, 7.5}, {590, 360, 6.8},
		// Entrance / exit
		{305, 52, 6.5}, {330, 60, 6.0}, {275, 58, 5.8},
		// Centre aisle
		{320, 200, 4.0}, {320, 280, 3.5},
		// Left aisle (produce)
		{100, 160, 3.0}, {115, 230, 2.8}, {120, 300, 2.5},
		// Right aisle (frozen)
		{530, 170, 2.8}, {545, 250, 2.4},
		// Display end-caps
		{200, 320, 4.5}, {440, 320, 4.0},
	}
	for _, p := range traffic {
		hm.AddPoint(p[0], p[1], p[2])
	}

	pngBytes, err := hm.GenerateOverlayPNG(base)
	if err != nil {
		log.Fatalf("overlay: %v", err)
	}
	if err := os.WriteFile("img/demo_foot_traffic.png", pngBytes, 0644); err != nil {
		log.Fatalf("save: %v", err)
	}
	log.Println("saved → img/demo_foot_traffic.png")
}
