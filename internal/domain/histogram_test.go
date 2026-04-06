package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHistogram(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		width  int
		height int
		gamma  float64
	}{
		{"small", 10, 10, 2.2},
		{"medium", 100, 100, 2.2},
		{"large", 1920, 1080, 2.2},
		{"gamma 1", 100, 100, 1.0},
		{"1x1", 1, 1, 2.2},
		{"0x0", 0, 0, 2.2},
		{"negative sizes", -1, -2, 2.2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			hist := NewHistogram(tt.width, tt.height, tt.gamma)
			require.NotNil(t, hist)
			expW := tt.width
			expH := tt.height
			if expW < 0 {
				expW = 0
			}
			if expH < 0 {
				expH = 0
			}
			assert.Equal(t, expW, hist.Width())
			assert.Equal(t, expH, hist.Height())
			assert.Equal(t, expW*scale, hist.ScaledWidth())
			assert.Equal(t, expH*scale, hist.ScaledHeight())
		})
	}
}

func TestHistogram_Ratio(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		width  int
		height int
		ratio  float64
	}{
		{"square", 100, 100, 1.0},
		{"wide", 1920, 1080, 1920.0 / 1080.0},
		{"tall", 100, 200, 0.5},
		{"zero height", 10, 0, 0.0},
		{"0/0", 0, 0, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			hist := NewHistogram(tt.width, tt.height, 2.2)
			assert.InDelta(t, tt.ratio, hist.Ratio(), 1e-10)
		})
	}
}

func TestHistogram_InBounds(t *testing.T) {
	t.Parallel()
	hist := NewHistogram(100, 100, 2.2)

	tests := []struct {
		name string
		x    int
		y    int
		want bool
	}{
		{"inside", 50, 50, true},
		{"top left", 0, 0, true},
		{"bottom right", hist.ScaledWidth() - 1, hist.ScaledHeight() - 1, true},
		{"negative x", -1, 50, false},
		{"negative y", 50, -1, false},
		{"too large x", hist.ScaledWidth(), 50, false},
		{"too large y", 50, hist.ScaledHeight(), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, hist.InBounds(tt.x, tt.y))
		})
	}
}

func TestHistogram_UpdatePixel(t *testing.T) {
	t.Parallel()
	t.Run("valid pixel", func(t *testing.T) {
		t.Parallel()
		hist := NewHistogram(100, 100, 2.2)
		color := Color{R: 1.0, G: 0.5, B: 0.25}
		for i := 0; i < 100; i++ {
			hist.UpdatePixel(50, 50, color)
		}

		pixels := hist.Correction()
		assert.NotNil(t, pixels)
	})

	t.Run("out of bounds", func(t *testing.T) {
		t.Parallel()
		hist := NewHistogram(100, 100, 2.2)
		color := Color{R: 1.0, G: 0.5, B: 0.25}
		hist.UpdatePixel(-1, -1, color)
		hist.UpdatePixel(1000, 1000, color)
	})

	t.Run("multiple updates", func(t *testing.T) {
		t.Parallel()
		hist := NewHistogram(100, 100, 2.2)
		color1 := Color{R: 1.0, G: 0.0, B: 0.0}
		color2 := Color{R: 0.0, G: 1.0, B: 0.0}
		for i := 0; i < 50; i++ {
			hist.UpdatePixel(50, 50, color1)
			hist.UpdatePixel(50, 50, color2)
		}

		pixels := hist.Correction()
		assert.NotNil(t, pixels)
	})
}

func TestHistogram_Correction_Output(t *testing.T) {
	t.Parallel()
	t.Run("empty histogram", func(t *testing.T) {
		t.Parallel()
		hist := NewHistogram(10, 10, 2.2)
		color := Color{R: 1.0, G: 0.5, B: 0.25}
		for i := 0; i < 100; i++ {
			hist.UpdatePixel(5*scale, 5*scale, color)
		}
		pixels := hist.Correction()

		require.NotNil(t, pixels)
		assert.Equal(t, 10, len(pixels))
		if len(pixels) > 0 {
			assert.Equal(t, 10, len(pixels[0]))
		}
	})

	t.Run("with updates", func(t *testing.T) {
		t.Parallel()
		hist := NewHistogram(10, 10, 2.2)
		color := Color{R: 1.0, G: 0.5, B: 0.25}

		for i := 0; i < 10; i++ {
			for y := 0; y < 5; y++ {
				for x := 0; x < 5; x++ {
					hist.UpdatePixel(x*scale, y*scale, color)
				}
			}
		}

		pixels := hist.Correction()
		require.NotNil(t, pixels)
		assert.Equal(t, 10, len(pixels))
		if len(pixels) > 0 {
			assert.Equal(t, 10, len(pixels[0]))
		}
	})
}

func TestHistogram_Correction(t *testing.T) {
	t.Parallel()
	t.Run("zero gamma", func(t *testing.T) {
		t.Parallel()
		hist := NewHistogram(10, 10, 0.0)
		assert.Nil(t, hist.Correction())
	})

	t.Run("negative gamma", func(t *testing.T) {
		t.Parallel()
		hist := NewHistogram(10, 10, -1.0)
		assert.Nil(t, hist.Correction())
	})

	t.Run("normal gamma", func(t *testing.T) {
		t.Parallel()
		hist := NewHistogram(10, 10, 2.2)
		color := Color{R: 1.0, G: 0.5, B: 0.25}
		for i := 0; i < 10; i++ {
			for y := 0; y < 10; y++ {
				for x := 0; x < 10; x++ {
					hist.UpdatePixel(x*scale, y*scale, color)
				}
			}
		}

		pixels := hist.Correction()
		require.NotNil(t, pixels)
		for y := range pixels {
			for x := range pixels[y] {
				assert.GreaterOrEqual(t, pixels[y][x].Color.R, 0.0)
				assert.LessOrEqual(t, pixels[y][x].Color.R, 1.0)
				assert.GreaterOrEqual(t, pixels[y][x].Color.G, 0.0)
				assert.LessOrEqual(t, pixels[y][x].Color.G, 1.0)
				assert.GreaterOrEqual(t, pixels[y][x].Color.B, 0.0)
				assert.LessOrEqual(t, pixels[y][x].Color.B, 1.0)
			}
		}
	})

	t.Run("gamma 1.0", func(t *testing.T) {
		t.Parallel()
		hist := NewHistogram(10, 10, 1.0)
		color := Color{R: 1.0, G: 0.5, B: 0.25}
		for i := 0; i < 10; i++ {
			for y := 0; y < 10; y++ {
				for x := 0; x < 10; x++ {
					hist.UpdatePixel(x*scale, y*scale, color)
				}
			}
		}

		pixels := hist.Correction()
		assert.NotNil(t, pixels)
	})
}

func TestHistogram_AvgCellFreq(t *testing.T) {
	t.Parallel()
	hist := NewHistogram(10, 10, 2.2)

	color := Color{R: 1.0, G: 0.5, B: 0.25}
	for y := 0; y < scale; y++ {
		for x := 0; x < scale; x++ {
			hist.UpdatePixel(x, y, color)
		}
	}

	freq := hist.avgCellFreq(0, 0)
	assert.GreaterOrEqual(t, freq, 0.0)
}

func TestHistogram_AvgCellColor(t *testing.T) {
	t.Parallel()
	hist := NewHistogram(10, 10, 2.2)

	color := Color{R: 1.0, G: 0.5, B: 0.25}
	for y := 0; y < scale; y++ {
		for x := 0; x < scale; x++ {
			hist.UpdatePixel(x, y, color)
		}
	}

	avgColor := hist.avgCellColor(0, 0)
	assert.GreaterOrEqual(t, avgColor.R, 0.0)
	assert.LessOrEqual(t, avgColor.R, 1.0)
	assert.GreaterOrEqual(t, avgColor.G, 0.0)
	assert.LessOrEqual(t, avgColor.G, 1.0)
	assert.GreaterOrEqual(t, avgColor.B, 0.0)
	assert.LessOrEqual(t, avgColor.B, 1.0)
}

func TestHistogram_EdgeCases(t *testing.T) {
	t.Parallel()
	t.Run("single pixel", func(t *testing.T) {
		t.Parallel()
		hist := NewHistogram(1, 1, 2.2)
		color := Color{R: 1.0, G: 0.5, B: 0.25}
		for i := 0; i < 100; i++ {
			hist.UpdatePixel(0, 0, color)
		}
		pixels := hist.Correction()
		require.NotNil(t, pixels)
		assert.Equal(t, 1, len(pixels))
		if len(pixels) > 0 {
			assert.Equal(t, 1, len(pixels[0]))
		}
	})

	t.Run("very small histogram", func(t *testing.T) {
		t.Parallel()
		hist := NewHistogram(2, 2, 2.2)
		color := Color{R: 1.0, G: 0.5, B: 0.25}
		for i := 0; i < 100; i++ {
			hist.UpdatePixel(0, 0, color)
		}
		pixels := hist.Correction()
		assert.NotNil(t, pixels)
	})

	t.Run("corner pixels", func(t *testing.T) {
		t.Parallel()
		hist := NewHistogram(10, 10, 2.2)
		color := Color{R: 1.0, G: 0.5, B: 0.25}

		for i := 0; i < 30; i++ {
			hist.UpdatePixel(0, 0, color)
			hist.UpdatePixel(hist.ScaledWidth()-1, 0, color)
			hist.UpdatePixel(0, hist.ScaledHeight()-1, color)
			hist.UpdatePixel(hist.ScaledWidth()-1, hist.ScaledHeight()-1, color)
		}

		pixels := hist.Correction()
		assert.NotNil(t, pixels)
	})

	t.Run("0x0 histogram", func(t *testing.T) {
		t.Parallel()
		hist := NewHistogram(0, 0, 2.2)
		assert.Equal(t, 0.0, hist.Ratio())
		pixels := hist.Correction()
		require.NotNil(t, pixels)
		assert.Len(t, pixels, 0)
	})

	t.Run("negative sizes do not panic", func(t *testing.T) {
		t.Parallel()
		hist := NewHistogram(-10, -20, 2.2)
		assert.Equal(t, 0, hist.Width())
		assert.Equal(t, 0, hist.Height())
		assert.Equal(t, 0, hist.ScaledWidth())
		assert.Equal(t, 0, hist.ScaledHeight())
		pixels := hist.Correction()
		require.NotNil(t, pixels)
		assert.Len(t, pixels, 0)
	})
}
