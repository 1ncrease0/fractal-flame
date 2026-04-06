package application

import (
	"io"
	"log/slog"
	"testing"
	"time"

	"fractal-flame/internal/domain"
)

// go test ./internal/application -bench=BenchmarkRenderer_Render  -v -run='^$'

func createBenchmarkParams(threads int) *Params {
	affineList := domain.NewAffineList()

	aff1 := domain.Affine{
		A:      0.5,
		B:      0.0,
		C:      0.0,
		D:      0.0,
		E:      0.5,
		F:      0.1,
		Weight: 1.0,
		Color:  domain.Color{R: 1.0, G: 0.5, B: 0.25},
	}
	_ = affineList.Add(aff1)

	aff2 := domain.Affine{
		A:      -0.3,
		B:      0.0,
		C:      0.0,
		D:      0.0,
		E:      -0.3,
		F:      0.0,
		Weight: 1.0,
		Color:  domain.Color{R: 0.2, G: 0.6, B: 0.8},
	}
	_ = affineList.Add(aff2)

	varList := domain.NewVariationList()

	if v, err := domain.GetVariation("spherical"); err == nil {
		_ = varList.Add(v, 1.0)
	}

	params, _ := NewParams(100000000, 5, 1, threads, 1920, 1080, 2.2, *affineList, *varList)
	return params
}
func BenchmarkRenderer_Render(b *testing.B) {
	threadCounts := []int{1, 2, 4, 8}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	results := make(map[int]time.Duration)

	b.StopTimer()

	for _, tc := range threadCounts {
		params := createBenchmarkParams(tc)

		start := time.Now()
		NewRenderer(*params, logger).Render()
		elapsed := time.Since(start)

		results[tc] = elapsed
	}

	b.StartTimer()

	b.Log("\n========== Renderer Benchmark Results ==========")
	b.Logf("%-8s %18s %12s", "Threads", "Time (ms)", "Speedup")
	b.Log("========================================================")

	base := results[1]
	for _, tc := range threadCounts {
		elapsed := results[tc]
		ms := float64(elapsed) / float64(time.Millisecond)
		speedup := float64(base) / float64(elapsed)
		b.Logf("%-8d %18.2f %12.2fx", tc, ms, speedup)
	}
	b.Log("========================================================")
}
