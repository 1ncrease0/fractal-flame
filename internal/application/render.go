package application

import (
	"fmt"
	"fractal-flame/internal/domain"
	"log/slog"
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

const (
	warmUp           = 20
	updateProgressMs = 500
	batch            = int64(10000)
)

type Renderer struct {
	affineList domain.AffineList
	variations domain.VariationList
	histogram  *domain.Histogram
	iters      int64
	threads    int
	symmetry   int
	rands      []*rand.Rand
	done       int64
	log        *slog.Logger
}

func NewRenderer(params Params, logger *slog.Logger) *Renderer {
	rands := make([]*rand.Rand, params.Threads())
	for i := 0; i < params.Threads(); i++ {
		rands[i] = rand.New(rand.NewSource(params.Seed() + int64(i)))
	}

	return &Renderer{
		affineList: params.Affine(),
		variations: params.Variations(),
		histogram:  domain.NewHistogram(params.Width(), params.Height(), params.Gamma()),
		iters:      params.Iterations(),
		symmetry:   params.Symmetry(),
		threads:    params.Threads(),
		rands:      rands,
		log:        logger,
	}
}

func (r *Renderer) Render() *domain.Histogram {

	minX, maxX := -1.5*r.histogram.Ratio(), 1.5*r.histogram.Ratio()
	minY, maxY := -1.5, 1.5

	perThread := r.iters / int64(r.threads)
	var wg sync.WaitGroup

	ticker := time.NewTicker(updateProgressMs * time.Millisecond)
	doneCh := make(chan struct{})
	var wgProgress sync.WaitGroup
	wgProgress.Add(1)
	go func() {
		defer wgProgress.Done()
		lastPercent := -1
		for {
			select {
			case <-ticker.C:
				done := atomic.LoadInt64(&r.done)
				percent := float64(done) / float64(r.iters) * 100.0
				currentPercent := int(percent)
				if currentPercent != lastPercent {
					r.log.Info("Progress", "percent", fmt.Sprintf("%.2f%%", percent), "done", done, "total", r.iters)
					lastPercent = currentPercent
				}
			case <-doneCh:
				done := atomic.LoadInt64(&r.done)
				r.log.Info("Progress", "percent", "100.00%", "done", done, "total", r.iters)
				return
			}
		}
	}()

	for i := 0; i < r.threads; i++ {
		iters := perThread
		if i == r.threads-1 {
			iters = r.iters - perThread*int64(r.threads-1)
		}
		wg.Add(1)
		idx := i
		go func() {
			defer wg.Done()
			r.run(idx, iters, minX, maxX, minY, maxY)
		}()
	}

	wg.Wait()
	ticker.Stop()
	close(doneCh)
	wgProgress.Wait()
	return r.histogram
}

func (r *Renderer) run(idx int, iters int64, minX, maxX, minY, maxY float64) {
	var x, y float64
	x = minX + r.rands[idx].Float64()*(maxX-minX)
	y = minY + r.rands[idx].Float64()*(maxY-minY)

	color := domain.Color{
		R: r.rands[idx].Float64(),
		G: r.rands[idx].Float64(),
		B: r.rands[idx].Float64(),
	}

	angle := 2 * math.Pi / float64(r.symmetry)

	var local int64

	for i := int64(0); i < iters; i++ {
		affine := r.affineList.Random(r.rands[idx])

		color.R = (color.R + affine.Color.R) / 2
		color.G = (color.G + affine.Color.G) / 2
		color.B = (color.B + affine.Color.B) / 2

		x, y = affine.Apply(x, y)
		x, y = r.variations.Apply(x, y)

		if i > warmUp && x >= minX && y >= minY && x <= maxX && y <= maxY {
			for k := 0; k < r.symmetry; k++ {
				rotAngle := float64(k) * angle
				rotX, rotY := rotate(x, y, rotAngle)

				if rotX >= minX && rotY >= minY && rotX <= maxX && rotY <= maxY {
					hx := int((rotX - minX) / (maxX - minX) * float64(r.histogram.ScaledWidth()))
					hy := int((rotY - minY) / (maxY - minY) * float64(r.histogram.ScaledHeight()))

					if hx < 0 {
						hx = 0
					}
					if hx >= r.histogram.ScaledWidth() {
						hx = r.histogram.ScaledWidth() - 1
					}
					if hy < 0 {
						hy = 0
					}
					if hy >= r.histogram.ScaledHeight() {
						hy = r.histogram.ScaledHeight() - 1
					}

					r.histogram.UpdatePixel(hx, hy, color)
				}
			}
		}

		local++
		if local >= batch {
			atomic.AddInt64(&r.done, local)
			local = 0
		}
	}

	if local > 0 {
		atomic.AddInt64(&r.done, local)
	}

}

func rotate(x, y, angle float64) (float64, float64) {
	cosA := math.Cos(angle)
	sinA := math.Sin(angle)
	rotX := x*cosA - y*sinA
	rotY := x*sinA + y*cosA
	return rotX, rotY
}
