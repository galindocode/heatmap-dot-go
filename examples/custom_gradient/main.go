// Custom 6-stop gradient on a dark background — standalone heatmap.
// Run from the project root: go run ./examples/custom_gradient/
package main

import (
	"image/color"
	"log"

	"github.com/galindocode/heatmap-dot-go/heatmap"
)

func main() {
	hm, err := heatmap.NewBuilder().
		Size(800, 500).
		Colors("#0d0887", "#6a00a8", "#b12a90", "#e16462", "#fca636", "#f0f921").
		PointSize(70).
		Alpha(215).
		MaxValue(1.0).
		Gaussian(true).
		Background(color.RGBA{8, 8, 20, 255}).
		Build()
	if err != nil {
		log.Fatalf("build: %v", err)
	}

	// Grid of points with value increasing diagonally.
	for col := 0; col < 8; col++ {
		for row := 0; row < 5; row++ {
			x := float64(80 + col*90)
			y := float64(70 + row*90)
			value := float64(col+row) / 11.0
			hm.AddPoint(x, y, value)
		}
	}

	if err := hm.SavePNG("img/demo_custom_gradient.png"); err != nil {
		log.Fatalf("save: %v", err)
	}
	log.Println("saved → img/demo_custom_gradient.png")
}
