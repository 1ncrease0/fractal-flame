package domain

import "sync"

type Pixel struct {
	Color Color
	Alpha int64
	mu    sync.Mutex
}

func (p *Pixel) UpdatePixel(color Color) {
	p.mu.Lock()
	p.Color.R = (p.Color.R + color.R) / 2
	p.Color.G = (p.Color.G + color.G) / 2
	p.Color.B = (p.Color.B + color.B) / 2
	p.Alpha += 10
	p.mu.Unlock()
}
