package domain

import (
	"math"
)

const (
	freqThreshold = 10.0
	scale         = 3
)

type Histogram struct {
	width        int
	height       int
	scaledWidth  int
	scaledHeight int
	gamma        float64
	pixels       [][]Pixel
}

func NewHistogram(width, height int, gamma float64) *Histogram {
	if width < 0 {
		width = 0
	}
	if height < 0 {
		height = 0
	}

	scaledWidth := width * scale
	scaledHeight := height * scale

	pixels := make([][]Pixel, scaledHeight)
	for i := range pixels {
		pixels[i] = make([]Pixel, scaledWidth)
	}

	return &Histogram{
		width:        width,
		height:       height,
		scaledWidth:  scaledWidth,
		scaledHeight: scaledHeight,
		pixels:       pixels,
		gamma:        gamma,
	}
}

func (h *Histogram) InBounds(x, y int) bool {
	if x < 0 || x >= h.scaledWidth || y < 0 || y >= h.scaledHeight {
		return false
	}
	return true
}

func (h *Histogram) Ratio() float64 {
	if h.height <= 0 {
		return 0
	}
	return float64(h.width) / float64(h.height)
}

func (h *Histogram) Width() int {
	return h.width
}

func (h *Histogram) Height() int {
	return h.height
}

func (h *Histogram) ScaledHeight() int {
	return h.scaledHeight
}

func (h *Histogram) ScaledWidth() int {
	return h.scaledWidth
}

func (h *Histogram) UpdatePixel(x, y int, color Color) {
	if h.InBounds(x, y) {
		h.pixels[y][x].UpdatePixel(color)
	}

}

func (h *Histogram) avgCellFreq(x, y int) float64 {
	var sum int64
	sx := x * scale
	sy := y * scale

	for hy := sy; hy < sy+scale && hy < h.scaledHeight; hy++ {
		for hx := sx; hx < sx+scale && hx < h.scaledWidth; hx++ {
			sum += h.pixels[hy][hx].Alpha
		}
	}

	return float64(sum) / float64(scale*scale)
}

func (h *Histogram) avgCellColor(x, y int) Color {
	var sumR, sumG, sumB float64
	sx := x * scale
	sy := y * scale

	for hy := sy; hy < sy+scale && hy < h.scaledHeight; hy++ {
		for hx := sx; hx < sx+scale && hx < h.scaledWidth; hx++ {
			sumR += h.pixels[hy][hx].Color.R
			sumG += h.pixels[hy][hx].Color.G
			sumB += h.pixels[hy][hx].Color.B
		}
	}

	return Color{
		R: sumR / float64(scale*scale),
		G: sumG / float64(scale*scale),
		B: sumB / float64(scale*scale),
	}
}

func (h *Histogram) Correction() [][]Pixel {
	if h.width <= 0 || h.height <= 0 {
		return make([][]Pixel, 0)
	}
	if h.gamma <= 0 {
		return nil
	}

	freqs := make([]float64, h.width*h.height)
	maxAvgFreq := 0.0
	for y := 0; y < h.height; y++ {
		rowOff := y * h.width
		for x := 0; x < h.width; x++ {
			f := h.avgCellFreq(x, y)
			freqs[rowOff+x] = f
			if f > maxAvgFreq {
				maxAvgFreq = f
			}
		}
	}

	maxFreq := maxAvgFreq / freqThreshold
	if maxFreq <= 1.0 {
		finalPixels := make([][]Pixel, h.height)
		for y := 0; y < h.height; y++ {
			finalPixels[y] = make([]Pixel, h.width)
		}
		return finalPixels
	}

	maxLogFreq := math.Log(maxFreq)
	invGamma := 1.0 / h.gamma

	finalPixels := make([][]Pixel, h.height)
	for y := 0; y < h.height; y++ {
		finalPixels[y] = make([]Pixel, h.width)
		rowOff := y * h.width
		for x := 0; x < h.width; x++ {
			avgFreq := freqs[rowOff+x]
			if avgFreq <= 0 {
				finalPixels[y][x] = Pixel{Color: Color{R: 0, G: 0, B: 0}, Alpha: 0}
				continue
			}

			avgColor := h.avgCellColor(x, y)
			avgColor.Vibrant(satShift, lightShift)

			alpha := math.Log(avgFreq) / maxLogFreq
			if alpha < 0 {
				alpha = 0
			} else if alpha > 1 {
				alpha = 1
			}

			alphaGamma := math.Pow(alpha, invGamma)
			c := Color{
				R: avgColor.R * alphaGamma,
				G: avgColor.G * alphaGamma,
				B: avgColor.B * alphaGamma,
			}
			c.R = math.Min(1.0, math.Max(0.0, c.R))
			c.G = math.Min(1.0, math.Max(0.0, c.G))
			c.B = math.Min(1.0, math.Max(0.0, c.B))

			finalPixels[y][x] = Pixel{Color: c, Alpha: int64(avgFreq)}
		}
	}

	return finalPixels
}
