// Standalone heatmap — no background image, dark canvas.
// Run from the project root: go run ./examples/simple/
package main

import (
	"image/color"
	"log"
	"math"
	"math/rand"

	"github.com/galindocode/heatmap-dot-go/heatmap"
)

func main() {
	hm, err := heatmap.NewBuilder().
		Size(800, 500).
		Colors("#3b82f6", "#22c55e", "#eab308", "#ef4444").
		PointSize(120).
		Alpha(230).
		MaxValue(1.0).
		Gaussian(true).
		Background(color.RGBA{12, 12, 22, 255}).
		Build()
	if err != nil {
		log.Fatalf("build: %v", err)
	}

	addCluster(hm, 200, 150, 70, 120, 1.0)
	addCluster(hm, 600, 340, 55, 90, 0.85)
	addCluster(hm, 400, 260, 65, 100, 0.7)

	// Sparse background noise.
	r := rand.New(rand.NewSource(42))
	for range 80 {
		hm.AddPoint(r.Float64()*800, r.Float64()*500, r.Float64()*0.25)
	}

	if err := hm.SavePNG("img/demo_simple.png"); err != nil {
		log.Fatalf("save: %v", err)
	}
	log.Println("saved → img/demo_simple.png")
}

func addCluster(hm *heatmap.Heatmap, cx, cy, spread float64, n int, peak float64) {
	r := rand.New(rand.NewSource(int64(cx + cy*1000)))
	for range n {
		angle := r.Float64() * 2 * math.Pi
		d := r.Float64() * spread
		value := peak * (1 - d/spread*0.6)
		hm.AddPoint(cx+d*math.Cos(angle), cy+d*math.Sin(angle), value)
	}
}
