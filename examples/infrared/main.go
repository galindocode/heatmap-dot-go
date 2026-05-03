// Infrared / thermal-camera color scheme overlaid on a real image.
// Run from the project root: go run ./examples/infrared/
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

	hm, err := heatmap.NewBuilder().
		Size(640, 424).
		Colors("#000033", "#0000ff", "#00ffff", "#00ff00", "#ffff00", "#ff4500", "#ff0000").
		PointSize(110).
		Alpha(235).
		MaxValue(10.0).
		Gaussian(true).
		Build()
	if err != nil {
		log.Fatalf("build: %v", err)
	}

	for _, p := range [][3]float64{
		{225, 370, 10.0}, {260, 365, 9.8}, {300, 368, 9.5},
		{340, 370, 9.8}, {380, 365, 9.2}, {415, 372, 8.8},
		{315, 58, 8.0}, {280, 68, 7.0}, {350, 62, 7.5},
		{140, 190, 5.5}, {155, 270, 5.0},
		{490, 175, 4.5}, {505, 255, 4.0},
		{520, 345, 3.5},
	} {
		hm.AddPoint(p[0], p[1], p[2])
	}

	pngBytes, err := hm.GenerateOverlayPNG(base)
	if err != nil {
		log.Fatalf("overlay: %v", err)
	}
	if err := os.WriteFile("img/demo_infrared.png", pngBytes, 0644); err != nil {
		log.Fatalf("save: %v", err)
	}
	log.Println("saved → img/demo_infrared.png")
}
