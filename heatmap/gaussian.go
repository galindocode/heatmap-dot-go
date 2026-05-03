package heatmap

import (
	"image"
	"image/color"
	"math"
)

// splatGaussian adds a point's contribution to the accumulation buffer using Gaussian falloff.
// Sigma is set to radius/3 so ~99.7% of the weight falls within the radius.
func splatGaussian(buf [][]float64, width, height int, cx, cy, value float64, radius int) {
	sigma := float64(radius) / 3.0
	twoSigma2 := 2.0 * sigma * sigma

	x0 := int(cx) - radius
	x1 := int(cx) + radius
	y0 := int(cy) - radius
	y1 := int(cy) + radius

	// Early exit when the kernel is entirely outside the image.
	if x1 < 0 || y1 < 0 || x0 >= width || y0 >= height {
		return
	}

	if x0 < 0 {
		x0 = 0
	}
	if y0 < 0 {
		y0 = 0
	}
	if x1 >= width {
		x1 = width - 1
	}
	if y1 >= height {
		y1 = height - 1
	}

	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
			dx := float64(x) - cx
			dy := float64(y) - cy
			weight := math.Exp(-(dx*dx+dy*dy) / twoSigma2)
			buf[y][x] += value * weight
		}
	}
}

// bufferMax returns the maximum value in a 2-D float64 buffer.
func bufferMax(buf [][]float64, width, height int) float64 {
	max := 0.0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if buf[y][x] > max {
				max = buf[y][x]
			}
		}
	}
	return max
}

// generateGaussianImage renders the heatmap using an accumulation buffer and Gaussian falloff.
// Each point's intensity is splatted into the buffer; the buffer is then normalized and mapped
// to the color gradient. Alpha scales proportionally with intensity, giving smooth transparent edges.
func (h *Heatmap) generateGaussianImage() (*image.RGBA, error) {
	width, height := h.config.Width, h.config.Height

	buf := make([][]float64, height)
	for y := range buf {
		buf[y] = make([]float64, width)
	}

	for _, p := range h.points {
		splatGaussian(buf, width, height, p.X, p.Y, p.Value, h.config.PointSize)
	}

	maxBuf := bufferMax(buf, width, height)

	canvas := image.NewRGBA(image.Rect(0, 0, width, height))
	if h.config.Background != nil {
		h.fillBackground(canvas, *h.config.Background)
	}

	if maxBuf == 0 {
		return canvas, nil
	}

	maxVal := h.getMaxValue()
	gradient := &ColorGradient{
		Colors:   h.config.ColorScheme,
		MaxValue: maxVal,
	}

	hasBg := h.config.Background != nil

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			v := buf[y][x]
			if v <= 0 {
				continue
			}
			ratio := v / maxBuf
			colorValue := ratio * maxVal
			pixelAlpha := uint8(float64(h.config.Alpha) * ratio)
			if pixelAlpha == 0 {
				continue // leave background intact rather than overwriting with invisible pixel
			}
			src := color.RGBAModel.Convert(gradient.GetColor(colorValue, pixelAlpha)).(color.RGBA)
			if hasBg {
				// Blend heatmap over the already-filled background so the output
				// pixel is fully opaque. Transparent edges in a dark-background
				// image otherwise composite to white in PNG viewers.
				bg := canvas.RGBAAt(x, y)
				a := float64(src.A) / 255.0
				canvas.SetRGBA(x, y, color.RGBA{
					R: uint8(float64(src.R)*a + float64(bg.R)*(1-a)),
					G: uint8(float64(src.G)*a + float64(bg.G)*(1-a)),
					B: uint8(float64(src.B)*a + float64(bg.B)*(1-a)),
					A: 255,
				})
			} else {
				canvas.Set(x, y, src)
			}
		}
	}

	return canvas, nil
}
