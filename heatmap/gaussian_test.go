package heatmap

import (
	"image"
	"image/color"
	"math"
	"testing"
)

// makeBuffer allocates a zeroed 2-D float64 buffer (height rows × width cols).
func makeBuffer(width, height int) [][]float64 {
	buf := make([][]float64, height)
	for y := range buf {
		buf[y] = make([]float64, width)
	}
	return buf
}

// --- splatGaussian unit tests ---

func TestSplatGaussian_CenterGetsMaxWeight(t *testing.T) {
	buf := makeBuffer(100, 100)
	splatGaussian(buf, 100, 100, 50, 50, 1.0, 15)

	center := buf[50][50]
	if center < 0.99 {
		t.Errorf("center pixel should have weight ~1.0, got %f", center)
	}
}

func TestSplatGaussian_Falloff(t *testing.T) {
	buf := makeBuffer(100, 100)
	splatGaussian(buf, 100, 100, 50, 50, 1.0, 15)

	center := buf[50][50]
	edge := buf[50][65] // distance == radius

	if edge >= center {
		t.Errorf("edge pixel (%f) should be less than center (%f)", edge, center)
	}

	// sigma = 15/3 = 5, twoSigma2 = 50; at d=15: exp(-225/50) ≈ 0.011
	if edge > 0.05 {
		t.Errorf("edge weight should be very low (~0.011), got %f", edge)
	}
}

func TestSplatGaussian_ValueScales(t *testing.T) {
	buf1 := makeBuffer(100, 100)
	splatGaussian(buf1, 100, 100, 50, 50, 1.0, 15)

	buf2 := makeBuffer(100, 100)
	splatGaussian(buf2, 100, 100, 50, 50, 3.0, 15)

	ratio := buf2[50][50] / buf1[50][50]
	if math.Abs(ratio-3.0) > 0.001 {
		t.Errorf("value should scale contribution linearly: expected ratio 3.0, got %f", ratio)
	}
}

func TestSplatGaussian_Accumulation(t *testing.T) {
	buf := makeBuffer(100, 100)
	splatGaussian(buf, 100, 100, 50, 50, 1.0, 15)
	single := buf[50][50]

	bufDouble := makeBuffer(100, 100)
	splatGaussian(bufDouble, 100, 100, 50, 50, 1.0, 15)
	splatGaussian(bufDouble, 100, 100, 50, 50, 1.0, 15)
	double := bufDouble[50][50]

	if math.Abs(double-2*single) > 0.001 {
		t.Errorf("two overlapping splats should double center value: got %f, expected %f", double, 2*single)
	}
}

func TestSplatGaussian_SymmetryAroundCenter(t *testing.T) {
	buf := makeBuffer(200, 200)
	splatGaussian(buf, 200, 200, 100, 100, 1.0, 20)

	tests := [][2]int{
		{90, 100}, {110, 100},
		{100, 90}, {100, 110},
		{92, 92}, {108, 92}, {92, 108}, {108, 108},
	}
	for i := 0; i < len(tests)-1; i += 2 {
		a := buf[tests[i][1]][tests[i][0]]
		b := buf[tests[i+1][1]][tests[i+1][0]]
		if math.Abs(a-b) > 1e-9 {
			t.Errorf("symmetric pixels differ: %v=%f vs %v=%f", tests[i], a, tests[i+1], b)
		}
	}
}

func TestSplatGaussian_EdgePoint(t *testing.T) {
	buf := makeBuffer(100, 100)
	// Should not panic; kernel is partially inside the image.
	splatGaussian(buf, 100, 100, 0, 0, 1.0, 15)
	splatGaussian(buf, 100, 100, 99, 99, 1.0, 15)

	// Center pixels should have received weight.
	if buf[0][0] < 0.99 {
		t.Errorf("top-left corner should still get max weight when point is there")
	}
}

func TestSplatGaussian_OutOfBounds_NoContribution(t *testing.T) {
	buf := makeBuffer(100, 100)
	splatGaussian(buf, 100, 100, -50, -50, 1.0, 15) // entirely outside
	splatGaussian(buf, 100, 100, 200, 200, 1.0, 15) // entirely outside

	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			if buf[y][x] != 0 {
				t.Errorf("out-of-bounds point should not affect buffer at (%d,%d)", x, y)
				return
			}
		}
	}
}

func TestSplatGaussian_PartiallyOutOfBounds(t *testing.T) {
	buf := makeBuffer(100, 100)
	// Point at (-5, 50) with radius 15 — left side clipped.
	splatGaussian(buf, 100, 100, -5, 50, 1.0, 15)

	// Should have nonzero contribution near x=0..10
	if buf[50][0] == 0 {
		t.Error("partially out-of-bounds point should still contribute within image bounds")
	}
	// Should not contribute far right
	if buf[50][20] != 0 {
		t.Errorf("contribution should not reach x=20 when point is at x=-5 with radius 15")
	}
}

// --- bufferMax tests ---

func TestBufferMax_Basic(t *testing.T) {
	buf := makeBuffer(10, 10)
	buf[5][5] = 3.14
	buf[2][7] = 1.0

	max := bufferMax(buf, 10, 10)
	if math.Abs(max-3.14) > 1e-9 {
		t.Errorf("expected max 3.14, got %f", max)
	}
}

func TestBufferMax_Empty(t *testing.T) {
	buf := makeBuffer(10, 10)
	max := bufferMax(buf, 10, 10)
	if max != 0 {
		t.Errorf("empty buffer max should be 0, got %f", max)
	}
}

func TestBufferMax_SingleValue(t *testing.T) {
	buf := makeBuffer(5, 5)
	buf[0][0] = 7.0
	if bufferMax(buf, 5, 5) != 7.0 {
		t.Error("single nonzero value should be the max")
	}
}

// --- Gaussian mode integration tests ---

func TestGaussianMode_DefaultIsFalse(t *testing.T) {
	hm := New(100, 100)
	if hm.config.GaussianMode {
		t.Error("GaussianMode should be false by default")
	}
}

func TestSetGaussianMode(t *testing.T) {
	hm := New(100, 100)
	hm.SetGaussianMode(true)
	if !hm.config.GaussianMode {
		t.Error("SetGaussianMode(true) did not enable Gaussian mode")
	}
	hm.SetGaussianMode(false)
	if hm.config.GaussianMode {
		t.Error("SetGaussianMode(false) did not disable Gaussian mode")
	}
}

func TestGaussianMode_Generate(t *testing.T) {
	hm := New(100, 100)
	hm.SetGaussianMode(true)
	hm.AddPoint(50, 50, 1.0)

	img, err := hm.Generate()
	if err != nil {
		t.Fatalf("Generate() with GaussianMode failed: %v", err)
	}
	if img == nil {
		t.Fatal("Generate() returned nil")
	}
	bounds := img.Bounds()
	if bounds.Dx() != 100 || bounds.Dy() != 100 {
		t.Errorf("wrong dimensions: %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestGaussianMode_GeneratePNG(t *testing.T) {
	hm := New(100, 100)
	hm.SetGaussianMode(true)
	hm.AddPoint(50, 50, 1.0)

	pngBytes, err := hm.GeneratePNG()
	if err != nil {
		t.Fatalf("GeneratePNG() failed: %v", err)
	}
	if len(pngBytes) < 8 || pngBytes[0] != 0x89 || pngBytes[1] != 'P' {
		t.Error("GeneratePNG() did not return valid PNG data")
	}
}

func TestGaussianMode_CenterIsVisible(t *testing.T) {
	hm := New(200, 200)
	hm.SetGaussianMode(true)
	hm.SetPointSize(30)
	hm.AddPoint(100, 100, 1.0)

	img, err := hm.Generate()
	if err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}

	center := color.RGBAModel.Convert(img.At(100, 100)).(color.RGBA)
	if center.A == 0 {
		t.Error("center pixel should be visible (alpha > 0)")
	}
}

func TestGaussianMode_CornerIsTransparent(t *testing.T) {
	hm := New(200, 200)
	hm.SetGaussianMode(true)
	hm.SetPointSize(20)
	hm.AddPoint(100, 100, 1.0) // far from corners

	img, err := hm.Generate()
	if err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}

	corner := color.RGBAModel.Convert(img.At(0, 0)).(color.RGBA)
	if corner.A != 0 {
		t.Errorf("corner pixel should be transparent, got alpha=%d", corner.A)
	}
}

func TestGaussianMode_CenterHasMaxAlpha(t *testing.T) {
	hm := New(200, 200)
	hm.SetGaussianMode(true)
	hm.SetPointSize(30)
	hm.SetAlpha(200)
	hm.AddPoint(100, 100, 1.0)

	img, err := hm.Generate()
	if err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}

	// The hotspot pixel is normalized to ratio=1.0, so alpha = uint8(200 * 1.0) = 200.
	center := color.RGBAModel.Convert(img.At(100, 100)).(color.RGBA)
	if center.A < 195 {
		t.Errorf("center should have alpha near max (200), got %d", center.A)
	}
}

func TestGaussianMode_AlphaDecreasesWithDistance(t *testing.T) {
	hm := New(200, 200)
	hm.SetGaussianMode(true)
	hm.SetPointSize(30)
	hm.SetAlpha(200)
	hm.AddPoint(100, 100, 1.0)

	img, err := hm.Generate()
	if err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}

	center := color.RGBAModel.Convert(img.At(100, 100)).(color.RGBA)
	near := color.RGBAModel.Convert(img.At(115, 100)).(color.RGBA) // 15 px from center

	if near.A >= center.A {
		t.Errorf("alpha should decrease with distance: center=%d, near=%d", center.A, near.A)
	}
}

func TestGaussianMode_MidpointBetweenTwoPoints(t *testing.T) {
	// A pixel exactly between two flanking points accumulates contributions from both,
	// so it should be brighter than a single point gives it from the same distance.
	const radius = 40

	// Single point at x=80; check pixel at x=100 (20px away).
	hm1 := New(200, 200)
	hm1.SetGaussianMode(true)
	hm1.SetPointSize(radius)
	hm1.AddPoint(80, 100, 1.0)
	img1, _ := hm1.Generate()

	// Two flanking points at x=80 and x=120; midpoint is x=100.
	hm2 := New(200, 200)
	hm2.SetGaussianMode(true)
	hm2.SetPointSize(radius)
	hm2.AddPoint(80, 100, 1.0)
	hm2.AddPoint(120, 100, 1.0)
	img2, _ := hm2.Generate()

	mid1 := color.RGBAModel.Convert(img1.At(100, 100)).(color.RGBA)
	mid2 := color.RGBAModel.Convert(img2.At(100, 100)).(color.RGBA)

	if mid2.A <= mid1.A {
		t.Errorf("midpoint between two flanking points should have higher alpha than a single-point: single=%d, flanked=%d", mid1.A, mid2.A)
	}
}

func TestGaussianMode_AllPointsOutOfBounds_TransparentImage(t *testing.T) {
	hm := New(100, 100)
	hm.SetGaussianMode(true)
	hm.AddPoint(500, 500, 1.0)
	hm.AddPoint(-200, -200, 1.0)

	img, err := hm.Generate()
	if err != nil {
		t.Fatalf("should not error for out-of-bounds points: %v", err)
	}

	c := color.RGBAModel.Convert(img.At(50, 50)).(color.RGBA)
	if c.A != 0 {
		t.Errorf("image should be fully transparent when all points are out of bounds, got alpha=%d", c.A)
	}
}

func TestGaussianMode_Overlay(t *testing.T) {
	base := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			base.Set(x, y, color.White)
		}
	}

	hm := New(100, 100)
	hm.SetGaussianMode(true)
	hm.SetPointSize(15)
	hm.AddPoint(50, 50, 1.0)

	result, err := hm.GenerateOverlay(base)
	if err != nil {
		t.Fatalf("GenerateOverlay() failed: %v", err)
	}
	if result == nil {
		t.Fatal("GenerateOverlay() returned nil")
	}

	// Corners are far from the point so base (white) should show through.
	corner := color.RGBAModel.Convert(result.At(0, 0)).(color.RGBA)
	if corner.R < 250 || corner.G < 250 || corner.B < 250 {
		t.Errorf("corner should remain white (base showing through), got %v", corner)
	}

	// Center should be blended (not pure white).
	center := color.RGBAModel.Convert(result.At(50, 50)).(color.RGBA)
	if center.R == 255 && center.G == 255 && center.B == 255 {
		t.Error("center should be blended with heatmap color, not pure white")
	}
}

func TestGaussianMode_OverlayPNG(t *testing.T) {
	base := image.NewRGBA(image.Rect(0, 0, 50, 50))
	hm := New(50, 50)
	hm.SetGaussianMode(true)
	hm.AddPoint(25, 25, 1.0)

	pngBytes, err := hm.GenerateOverlayPNG(base)
	if err != nil {
		t.Fatalf("GenerateOverlayPNG() failed: %v", err)
	}
	if len(pngBytes) == 0 {
		t.Error("GenerateOverlayPNG() returned empty bytes")
	}
}

func TestGaussianMode_MultipleColorsInGradient(t *testing.T) {
	// Low-value point should be blueish, high-value point should be reddish.
	hm := New(300, 100)
	hm.SetGaussianMode(true)
	hm.SetPointSize(20)
	hm.SetMaxValue(10)
	hm.AddPoint(75, 50, 1.0)  // low intensity → blue end
	hm.AddPoint(225, 50, 9.0) // high intensity → red end

	img, err := hm.Generate()
	if err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}

	low := color.RGBAModel.Convert(img.At(75, 50)).(color.RGBA)
	high := color.RGBAModel.Convert(img.At(225, 50)).(color.RGBA)

	// High-value point should have more red component.
	if high.R <= low.R {
		t.Errorf("high-value point should be redder: low.R=%d, high.R=%d", low.R, high.R)
	}
	// Low-value point should have more blue component.
	if low.B <= high.B {
		t.Errorf("low-value point should be bluer: low.B=%d, high.B=%d", low.B, high.B)
	}
}

// --- Benchmarks ---

func BenchmarkSplatGaussian(b *testing.B) {
	buf := makeBuffer(800, 600)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		splatGaussian(buf, 800, 600, 400, 300, 1.0, 50)
	}
}

func BenchmarkGenerateGaussian_1000Points(b *testing.B) {
	hm := New(800, 600)
	hm.SetGaussianMode(true)
	hm.SetPointSize(50)
	for i := 0; i < 1000; i++ {
		hm.AddPoint(float64(i%800), float64(i%600), float64(i%100)/100.0)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := hm.Generate()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGenerateGaussian_vs_Circle(b *testing.B) {
	addPoints := func(hm *Heatmap) {
		for i := 0; i < 500; i++ {
			hm.AddPoint(float64(i%400), float64(i%300), float64(i%10)/10.0)
		}
	}

	b.Run("Circle", func(b *testing.B) {
		hm := New(400, 300)
		hm.SetPointSize(30)
		addPoints(hm)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = hm.Generate()
		}
	})

	b.Run("Gaussian", func(b *testing.B) {
		hm := New(400, 300)
		hm.SetGaussianMode(true)
		hm.SetPointSize(30)
		addPoints(hm)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = hm.Generate()
		}
	})
}
