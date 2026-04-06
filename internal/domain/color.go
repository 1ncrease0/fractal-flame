package domain

import "math"

const (
	satShift   = 2.0
	lightShift = 1.0
)

type Color struct {
	R float64
	G float64
	B float64
}

func (c *Color) Vibrant(satShift, lightShift float64) {
	hsl := rgbToHSL(c)
	hsl.S = (hsl.S + (math.Max(satShift, 0) * 100)) / (math.Abs(satShift) + 1)
	hsl.L = (hsl.L + (math.Max(lightShift, 0) * 100)) / (math.Abs(lightShift) + 1)
	rgb := hslToRGB(hsl)
	c.R = rgb.R / 255.0
	c.G = rgb.G / 255.0
	c.B = rgb.B / 255.0
}

type HSL struct {
	H float64
	S float64
	L float64
}

func rgbToHSL(c *Color) HSL {
	r := c.R
	g := c.G
	b := c.B
	maxColor := math.Max(math.Max(r, g), b)
	minColor := math.Min(math.Min(r, g), b)
	h := 0.0
	s := 0.0
	l := (maxColor + minColor) / 2
	if maxColor != minColor {
		if l < 0.5 {
			s = (maxColor - minColor) / (maxColor + minColor)
		} else {
			s = (maxColor - minColor) / (2.0 - maxColor - minColor)
		}

		if r == maxColor {
			h = (g - b) / (maxColor - minColor)
		} else if g == maxColor {
			h = 2.0 + (b-r)/(maxColor-minColor)
		} else {
			h = 4.0 + (r-g)/(maxColor-minColor)
		}
	}

	h *= 60
	s *= 100
	l *= 100

	if h < 0 {
		h += 360
	}
	return HSL{
		H: h,
		S: s,
		L: l,
	}
}

func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}
	if t < 1.0/6.0 {
		return p + (q-p)*6*t
	}
	if t < 1.0/2.0 {
		return q
	}
	if t < 2.0/3.0 {
		return p + (q-p)*(2.0/3.0-t)*6
	}
	return p
}

func hslToRGB(hsl HSL) *Color {
	h := hsl.H / 360.0
	s := hsl.S / 100.0
	l := hsl.L / 100.0

	var r, g, b float64
	if s == 0 {
		r, g, b = l, l, l
	} else {
		var p, q float64
		if l < 0.5 {
			q = l * (1 + s)
		} else {
			q = l + s - l*s
		}
		p = 2*l - q
		r = hueToRGB(p, q, h+1.0/3.0)
		g = hueToRGB(p, q, h)
		b = hueToRGB(p, q, h-1.0/3.0)
	}
	return &Color{
		R: math.Round(r * 255),
		G: math.Round(g * 255),
		B: math.Round(b * 255),
	}
}
