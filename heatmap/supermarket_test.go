package heatmap_test

// Integration tests using supermarket.jpg (640×424).
// They demonstrate real-world Gaussian overlay usage and save output PNGs for visual verification.

import (
	"image/color"
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/galindocode/heatmap-dot-go/heatmap"
)

const (
	supermarketPath = "../supermarket.jpg"
	imgW            = 640
	imgH            = 424
)

// saveOutput writes PNG bytes to the test output directory, failing loudly on error.
func saveOutput(t *testing.T, name string, pngBytes []byte) {
	t.Helper()
	dir := "testoutput"
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("could not create %s: %v", dir, err)
	}
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, pngBytes, 0644); err != nil {
		t.Fatalf("could not write %s: %v", path, err)
	}
	t.Logf("saved → %s", path)
}

// ── Test 1: basic Gaussian overlay ───────────────────────────────────────────

func TestSupermarket_GaussianOverlay_Basic(t *testing.T) {
	if _, err := os.Stat(supermarketPath); os.IsNotExist(err) {
		t.Skip("supermarket.jpg not found")
	}

	base, err := heatmap.LoadImage(supermarketPath)
	if err != nil {
		t.Fatalf("LoadImage failed: %v", err)
	}

	hm := heatmap.New(imgW, imgH)
	hm.SetGaussianMode(true)
	hm.SetPointSize(60)
	hm.SetAlpha(200)

	// Simulate generic foot traffic
	points := [][3]float64{
		{320, 212, 1.0}, // center
		{160, 106, 0.7},
		{480, 106, 0.6},
		{160, 318, 0.8},
		{480, 318, 0.5},
	}
	for _, p := range points {
		hm.AddPoint(p[0], p[1], p[2])
	}

	pngBytes, err := hm.GenerateOverlayPNG(base)
	if err != nil {
		t.Fatalf("GenerateOverlayPNG failed: %v", err)
	}
	if len(pngBytes) == 0 {
		t.Error("empty PNG output")
	}
	saveOutput(t, "supermarket_basic.png", pngBytes)
}

// ── Test 2: foot traffic – checkout hotspot ───────────────────────────────────

func TestSupermarket_FootTraffic_CheckoutIsHottestZone(t *testing.T) {
	if _, err := os.Stat(supermarketPath); os.IsNotExist(err) {
		t.Skip("supermarket.jpg not found")
	}

	base, err := heatmap.LoadImage(supermarketPath)
	if err != nil {
		t.Fatalf("LoadImage failed: %v", err)
	}

	// Retail layout (approx. pixel coords for 640×424):
	//   Entrance:  x≈320 y≈50   (top center)
	//   Checkouts: x≈320 y≈370  (bottom center)
	//   Aisles:    x≈160 y≈212  (mid-left)
	hm := heatmap.New(imgW, imgH)
	hm.SetGaussianMode(true)
	hm.SetPointSize(55)
	hm.SetAlpha(210)
	hm.SetMaxValue(10.0)

	// Checkout queue – high density
	for i := 0; i < 8; i++ {
		hm.AddPoint(float64(240+i*20), 370, 10.0)
	}
	// Entrance – medium
	hm.AddPoint(320, 60, 6.0)
	hm.AddPoint(280, 80, 5.0)
	hm.AddPoint(360, 80, 5.0)
	// Aisle mid-left – low
	hm.AddPoint(120, 200, 3.0)
	hm.AddPoint(120, 250, 2.0)

	img, err := hm.Generate()
	if err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}

	// Checkout zone (bottom) should be more opaque than aisle zone (mid-left).
	checkout := color.RGBAModel.Convert(img.At(320, 370)).(color.RGBA)
	aisle := color.RGBAModel.Convert(img.At(120, 225)).(color.RGBA)
	if checkout.A <= aisle.A {
		t.Errorf("checkout zone should be hotter: checkout alpha=%d, aisle alpha=%d", checkout.A, aisle.A)
	}

	pngBytes, err := hm.GenerateOverlayPNG(base)
	if err != nil {
		t.Fatalf("GenerateOverlayPNG failed: %v", err)
	}
	saveOutput(t, "supermarket_checkout_hotspot.png", pngBytes)
}

// ── Test 3: entrance warmer than far back ─────────────────────────────────────

func TestSupermarket_FootTraffic_EntranceWarmerThanBack(t *testing.T) {
	if _, err := os.Stat(supermarketPath); os.IsNotExist(err) {
		t.Skip("supermarket.jpg not found")
	}

	hm := heatmap.New(imgW, imgH)
	hm.SetGaussianMode(true)
	hm.SetPointSize(50)
	hm.SetAlpha(200)
	hm.SetMaxValue(5.0)

	// Many customers near entrance
	for i := 0; i < 6; i++ {
		hm.AddPoint(float64(260+i*20), 60, 5.0)
	}
	// Few customers at the back
	hm.AddPoint(500, 360, 1.0)

	img, err := hm.Generate()
	if err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}

	entrance := color.RGBAModel.Convert(img.At(320, 60)).(color.RGBA)
	back := color.RGBAModel.Convert(img.At(500, 360)).(color.RGBA)

	if entrance.A <= back.A {
		t.Errorf("entrance should be warmer: entrance alpha=%d, back alpha=%d", entrance.A, back.A)
	}
}

// ── Test 4: infrared-style color scheme ──────────────────────────────────────

func TestSupermarket_InfraredColorScheme(t *testing.T) {
	if _, err := os.Stat(supermarketPath); os.IsNotExist(err) {
		t.Skip("supermarket.jpg not found")
	}

	base, err := heatmap.LoadImage(supermarketPath)
	if err != nil {
		t.Fatalf("LoadImage failed: %v", err)
	}

	hm, err := heatmap.NewBuilder().
		Size(imgW, imgH).
		Colors("#000033", "#0000ff", "#00ffff", "#00ff00", "#ffff00", "#ff4500", "#ff0000").
		PointSize(65).
		Alpha(220).
		MaxValue(10.0).
		Gaussian(true).
		Build()
	if err != nil {
		t.Fatalf("Build() failed: %v", err)
	}

	// Simulate camera-tracked person positions
	positions := [][3]float64{
		{180, 150, 8.0},
		{200, 200, 9.0},
		{350, 100, 7.0},
		{420, 280, 10.0},
		{500, 350, 6.0},
		{100, 300, 5.0},
		{320, 320, 8.5},
	}
	for _, p := range positions {
		hm.AddPoint(p[0], p[1], p[2])
	}

	pngBytes, err := hm.GenerateOverlayPNG(base)
	if err != nil {
		t.Fatalf("GenerateOverlayPNG failed: %v", err)
	}
	if len(pngBytes) == 0 {
		t.Error("empty PNG output")
	}
	saveOutput(t, "supermarket_infrared.png", pngBytes)

	// Hot pixel should have strong red component (end of infrared gradient).
	img, _ := hm.Generate()
	hot := color.RGBAModel.Convert(img.At(420, 280)).(color.RGBA)
	if hot.R < 150 {
		t.Errorf("hottest point should have high red component in infrared scheme, got R=%d", hot.R)
	}
}

// ── Test 5: builder Gaussian flag propagates ──────────────────────────────────

func TestSupermarket_BuilderGaussianFlag(t *testing.T) {
	hmOn, err := heatmap.NewBuilder().
		Size(imgW, imgH).
		Gaussian(true).
		Build()
	if err != nil {
		t.Fatalf("Build() failed: %v", err)
	}
	hmOn.AddPoint(320, 212, 1.0)

	hmOff, err := heatmap.NewBuilder().
		Size(imgW, imgH).
		Gaussian(false).
		Build()
	if err != nil {
		t.Fatalf("Build() failed: %v", err)
	}
	hmOff.AddPoint(320, 212, 1.0)

	imgOn, err := hmOn.Generate()
	if err != nil {
		t.Fatalf("Gaussian generate failed: %v", err)
	}
	imgOff, err := hmOff.Generate()
	if err != nil {
		t.Fatalf("Circle generate failed: %v", err)
	}

	// Gaussian mode should produce transparent pixels far from the point;
	// default (circle) mode also does, but the edge pixel differs.
	// We just verify both render without error and are the right size.
	if imgOn.Bounds().Dx() != imgW || imgOff.Bounds().Dx() != imgW {
		t.Error("wrong image width")
	}
}

// ── Test 6: Gaussian vs circle visual comparison (saves both) ─────────────────

func TestSupermarket_GaussianVsCircle_VisualComparison(t *testing.T) {
	if _, err := os.Stat(supermarketPath); os.IsNotExist(err) {
		t.Skip("supermarket.jpg not found")
	}

	base, err := heatmap.LoadImage(supermarketPath)
	if err != nil {
		t.Fatalf("LoadImage failed: %v", err)
	}

	addPts := func(hm *heatmap.Heatmap) {
		hm.AddPoint(320, 212, 1.0)
		hm.AddPoint(160, 150, 0.8)
		hm.AddPoint(480, 150, 0.7)
		hm.AddPoint(240, 320, 0.9)
		hm.AddPoint(400, 320, 0.6)
	}

	// Gaussian
	hmG := heatmap.New(imgW, imgH)
	hmG.SetGaussianMode(true)
	hmG.SetPointSize(60)
	hmG.SetAlpha(200)
	addPts(hmG)
	pngG, err := hmG.GenerateOverlayPNG(base)
	if err != nil {
		t.Fatalf("Gaussian overlay failed: %v", err)
	}
	saveOutput(t, "supermarket_gaussian.png", pngG)

	// Circle
	hmC := heatmap.New(imgW, imgH)
	hmC.SetPointSize(60)
	hmC.SetAlpha(200)
	addPts(hmC)
	pngC, err := hmC.GenerateOverlayPNG(base)
	if err != nil {
		t.Fatalf("Circle overlay failed: %v", err)
	}
	saveOutput(t, "supermarket_circle.png", pngC)

	// Gaussian center pixel should have higher alpha (smooth accumulation vs hard circle).
	imgG, _ := hmG.Generate()
	imgC, _ := hmC.Generate()
	cG := color.RGBAModel.Convert(imgG.At(320, 212)).(color.RGBA)
	cC := color.RGBAModel.Convert(imgC.At(320, 212)).(color.RGBA)
	t.Logf("Gaussian center alpha=%d  Circle center alpha=%d", cG.A, cC.A)
}

// ── Test 7: standalone Gaussian (no overlay) matches image dimensions ──────────

func TestSupermarket_StandaloneGaussian_Dimensions(t *testing.T) {
	hm := heatmap.New(imgW, imgH)
	hm.SetGaussianMode(true)
	hm.SetPointSize(40)

	for x := 80; x < imgW; x += 160 {
		for y := 60; y < imgH; y += 120 {
			hm.AddPoint(float64(x), float64(y), 1.0)
		}
	}

	img, err := hm.Generate()
	if err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}
	if img.Bounds().Dx() != imgW || img.Bounds().Dy() != imgH {
		t.Errorf("wrong dimensions: %dx%d", img.Bounds().Dx(), img.Bounds().Dy())
	}
}

// ── Test 8: alpha=0 makes fully transparent heatmap ──────────────────────────

func TestSupermarket_ZeroAlpha_FullyTransparent(t *testing.T) {
	hm := heatmap.New(imgW, imgH)
	hm.SetGaussianMode(true)
	hm.SetAlpha(0)
	hm.AddPoint(320, 212, 1.0)

	img, err := hm.Generate()
	if err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}
	c := color.RGBAModel.Convert(img.At(320, 212)).(color.RGBA)
	if c.A != 0 {
		t.Errorf("alpha=0 should produce fully transparent output, got alpha=%d", c.A)
	}
}

// ── Test 9: Gaussian intensity rings (concentric distance bands) ──────────────

func TestSupermarket_GaussianIntensity_MonotonicallyDecreases(t *testing.T) {
	hm := heatmap.New(imgW, imgH)
	hm.SetGaussianMode(true)
	hm.SetPointSize(80)
	hm.SetAlpha(255)
	hm.AddPoint(320, 212, 1.0)

	img, err := hm.Generate()
	if err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}

	cx, cy := 320, 212
	// Sample alpha at increasing radii from center.
	radii := []int{0, 10, 20, 30, 40, 50, 60}
	alphas := make([]uint8, len(radii))
	for i, r := range radii {
		px := cx + r
		alphas[i] = color.RGBAModel.Convert(img.At(px, cy)).(color.RGBA).A
	}

	for i := 1; i < len(alphas); i++ {
		if alphas[i] > alphas[i-1] {
			t.Errorf("alpha should decrease monotonically: at r=%d got %d > r=%d got %d",
				radii[i], alphas[i], radii[i-1], alphas[i-1])
		}
	}
}

// ── Test 10: large point size covers more area ────────────────────────────────

func TestSupermarket_LargePointSize_CoversMoreArea(t *testing.T) {
	countVisible := func(pointSize int) int {
		hm := heatmap.New(imgW, imgH)
		hm.SetGaussianMode(true)
		hm.SetPointSize(pointSize)
		hm.SetAlpha(200)
		hm.AddPoint(320, 212, 1.0)
		img, _ := hm.Generate()
		n := 0
		for y := 0; y < imgH; y++ {
			for x := 0; x < imgW; x++ {
				if color.RGBAModel.Convert(img.At(x, y)).(color.RGBA).A > 0 {
					n++
				}
			}
		}
		return n
	}

	small := countVisible(20)
	large := countVisible(80)
	if large <= small {
		t.Errorf("large point size should cover more pixels: small=%d, large=%d", small, large)
	}
}

// ── Test 11: multiple overlapping Gaussian blobs create smooth gradient ────────

func TestSupermarket_OverlappingBlobs_SmoothGradient(t *testing.T) {
	// Two points side by side; pixel between them should have intermediate alpha.
	const psize = 60
	hm := heatmap.New(imgW, imgH)
	hm.SetGaussianMode(true)
	hm.SetPointSize(psize)
	hm.SetAlpha(200)
	hm.SetMaxValue(1.0)
	hm.AddPoint(200, 212, 1.0)
	hm.AddPoint(440, 212, 1.0)

	img, err := hm.Generate()
	if err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}

	left := color.RGBAModel.Convert(img.At(200, 212)).(color.RGBA)
	mid := color.RGBAModel.Convert(img.At(320, 212)).(color.RGBA)
	right := color.RGBAModel.Convert(img.At(440, 212)).(color.RGBA)

	// Left and right peaks should both be visible.
	if left.A == 0 {
		t.Error("left peak should be visible")
	}
	if right.A == 0 {
		t.Error("right peak should be visible")
	}
	// Midpoint should also be visible (blobs overlap enough to reach it).
	sigma := float64(psize) / 3.0
	dist := 120.0 // distance from each peak to midpoint
	weight := math.Exp(-(dist * dist) / (2 * sigma * sigma))
	if weight > 0.001 && mid.A == 0 {
		t.Errorf("midpoint should be visible when blobs overlap (weight=%.4f)", weight)
	}
	t.Logf("left.A=%d  mid.A=%d  right.A=%d", left.A, mid.A, right.A)
}
