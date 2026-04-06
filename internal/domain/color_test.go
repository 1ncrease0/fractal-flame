package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//test mostly check valid ranges

func TestColor_Vibrant(t *testing.T) {
	t.Parallel()
	t.Run("basic vibrant", func(t *testing.T) {
		t.Parallel()
		color := Color{R: 0.5, G: 0.5, B: 0.5}

		color.Vibrant(2.0, 1.0)

		assert.GreaterOrEqual(t, color.R, 0.0)
		assert.LessOrEqual(t, color.R, 1.0)
		assert.GreaterOrEqual(t, color.G, 0.0)
		assert.LessOrEqual(t, color.G, 1.0)
		assert.GreaterOrEqual(t, color.B, 0.0)
		assert.LessOrEqual(t, color.B, 1.0)
	})

	t.Run("red color", func(t *testing.T) {
		t.Parallel()
		color := Color{R: 1.0, G: 0.0, B: 0.0}
		color.Vibrant(2.0, 1.0)

		assert.GreaterOrEqual(t, color.R, 0.0)
		assert.LessOrEqual(t, color.R, 1.0)
	})

	t.Run("zero color", func(t *testing.T) {
		t.Parallel()
		color := Color{R: 0.0, G: 0.0, B: 0.0}
		color.Vibrant(2.0, 1.0)

		assert.GreaterOrEqual(t, color.R, 0.0)
		assert.LessOrEqual(t, color.R, 1.0)
	})
}

func TestRGBToHSL(t *testing.T) {
	t.Parallel()
	t.Run("red", func(t *testing.T) {
		t.Parallel()
		color := Color{R: 1.0, G: 0.0, B: 0.0}
		hsl := rgbToHSL(&color)

		assert.InDelta(t, 0.0, hsl.H, 1.0)
		assert.Greater(t, hsl.S, 0.0)
		assert.Greater(t, hsl.L, 0.0)
	})

	t.Run("green", func(t *testing.T) {
		t.Parallel()
		color := Color{R: 0.0, G: 1.0, B: 0.0}
		hsl := rgbToHSL(&color)

		assert.Greater(t, hsl.H, 0.0)
		assert.Greater(t, hsl.S, 0.0)
		assert.Greater(t, hsl.L, 0.0)
	})

	t.Run("blue", func(t *testing.T) {
		t.Parallel()
		color := Color{R: 0.0, G: 0.0, B: 1.0}
		hsl := rgbToHSL(&color)

		assert.Greater(t, hsl.H, 0.0)
		assert.Greater(t, hsl.S, 0.0)
		assert.Greater(t, hsl.L, 0.0)
	})

	t.Run("white", func(t *testing.T) {
		t.Parallel()
		color := Color{R: 1.0, G: 1.0, B: 1.0}
		hsl := rgbToHSL(&color)

		assert.Greater(t, hsl.L, 50.0)
	})

	t.Run("black", func(t *testing.T) {
		t.Parallel()
		color := Color{R: 0.0, G: 0.0, B: 0.0}
		hsl := rgbToHSL(&color)

		assert.Less(t, hsl.L, 50.0)
	})
}

func TestHSLToRGB(t *testing.T) {
	t.Parallel()
	t.Run("red hue", func(t *testing.T) {
		t.Parallel()
		hsl := HSL{H: 0.0, S: 100.0, L: 50.0}
		rgb := hslToRGB(hsl)

		assert.Greater(t, rgb.R, 0.0)
		assert.Less(t, rgb.G, rgb.R)
		assert.Less(t, rgb.B, rgb.R)
	})

	t.Run("green hue", func(t *testing.T) {
		t.Parallel()
		hsl := HSL{H: 120.0, S: 100.0, L: 50.0}
		rgb := hslToRGB(hsl)

		assert.Greater(t, rgb.G, 0.0)
		assert.Less(t, rgb.R, rgb.G)
		assert.Less(t, rgb.B, rgb.G)
	})

	t.Run("blue hue", func(t *testing.T) {
		t.Parallel()
		hsl := HSL{H: 240.0, S: 100.0, L: 50.0}
		rgb := hslToRGB(hsl)

		assert.Greater(t, rgb.B, 0.0)
		assert.Less(t, rgb.R, rgb.B)
		assert.Less(t, rgb.G, rgb.B)
	})

	t.Run("gray (zero saturation)", func(t *testing.T) {
		t.Parallel()
		hsl := HSL{H: 0.0, S: 0.0, L: 50.0}
		rgb := hslToRGB(hsl)

		assert.InDelta(t, rgb.R, rgb.G, 1.0)
		assert.InDelta(t, rgb.G, rgb.B, 1.0)
	})
}

func TestHueToRGB(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		p    float64
		q    float64
		t    float64
	}{
		{"zero", 0.0, 1.0, 0.0},
		{"one sixth", 0.0, 1.0, 1.0 / 6.0},
		{"one half", 0.0, 1.0, 0.5},
		{"two thirds", 0.0, 1.0, 2.0 / 3.0},
		{"one", 0.0, 1.0, 1.0},
		{"negative", 0.0, 1.0, -0.1},
		{"greater than one", 0.0, 1.0, 1.1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := hueToRGB(tt.p, tt.q, tt.t)
			assert.GreaterOrEqual(t, result, 0.0)
			assert.LessOrEqual(t, result, 1.0)
		})
	}
}

func TestColorConversionRoundTrip(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		color Color
	}{
		{"red", Color{R: 1.0, G: 0.0, B: 0.0}},
		{"green", Color{R: 0.0, G: 1.0, B: 0.0}},
		{"blue", Color{R: 0.0, G: 0.0, B: 1.0}},
		{"white", Color{R: 1.0, G: 1.0, B: 1.0}},
		{"gray", Color{R: 0.5, G: 0.5, B: 0.5}},
		{"mixed", Color{R: 0.7, G: 0.3, B: 0.9}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			original := tt.color
			hsl := rgbToHSL(&original)
			rgb := hslToRGB(hsl)

			rgbNormalized := Color{
				R: rgb.R / 255.0,
				G: rgb.G / 255.0,
				B: rgb.B / 255.0,
			}

			assert.InDelta(t, original.R, rgbNormalized.R, 0.1)
			assert.InDelta(t, original.G, rgbNormalized.G, 0.1)
			assert.InDelta(t, original.B, rgbNormalized.B, 0.1)
		})
	}
}
