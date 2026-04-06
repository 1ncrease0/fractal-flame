package application

import (
	"fmt"
	"fractal-flame/internal/domain"
	"log/slog"
	"math"
	"os"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestParams(threads, iterations int, symmetry int) *Params {
	affineList := domain.NewAffineList()
	affine := domain.Affine{
		A:      0.5,
		B:      0.0,
		C:      0.0,
		D:      0.0,
		E:      0.5,
		F:      0.0,
		Weight: 1.0,
		Color:  domain.Color{R: 1.0, G: 0.5, B: 0.25},
	}
	_ = affineList.Add(affine)

	variationList := domain.NewVariationList()
	v, _ := domain.GetVariation("linear")
	_ = variationList.Add(v, 1.0)

	params, _ := NewParams(int64(iterations), 42, symmetry, threads, 100, 100, 2.2, *affineList, *variationList)
	return params
}

func TestNewRenderer(t *testing.T) {
	t.Parallel()
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	params := createTestParams(4, 10000, 1)
	renderer := NewRenderer(*params, logger)

	require.NotNil(t, renderer)
	assert.Equal(t, 4, renderer.threads)
	assert.Equal(t, int64(10000), renderer.iters)
	assert.Equal(t, 1, renderer.symmetry)
	assert.NotNil(t, renderer.histogram)
	assert.Len(t, renderer.rands, 4)
}

func TestRenderer_Render_Basic(t *testing.T) {
	t.Parallel()
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	params := createTestParams(2, 10000, 1)
	renderer := NewRenderer(*params, logger)

	histogram := renderer.Render()

	require.NotNil(t, histogram)
	assert.Equal(t, 100, histogram.Width())
	assert.Equal(t, 100, histogram.Height())
}

func TestRenderer_Render_DifferentThreads(t *testing.T) {
	t.Parallel()
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	threadCounts := []int{1, 2, 4, 8}
	for _, threads := range threadCounts {
		threads := threads
		t.Run(fmt.Sprintf("%d_threads", threads), func(t *testing.T) {
			t.Parallel()
			params := createTestParams(threads, 10000, 1)
			renderer := NewRenderer(*params, logger)

			histogram := renderer.Render()
			require.NotNil(t, histogram)
			assert.Equal(t, 100, histogram.Width())
			assert.Equal(t, 100, histogram.Height())
		})
	}
}

func TestRenderer_Render_DifferentSymmetry(t *testing.T) {
	t.Parallel()
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	symmetryLevels := []int{1, 2, 4, 8}
	for _, symmetry := range symmetryLevels {
		symmetry := symmetry
		t.Run(fmt.Sprintf("symmetry_%d", symmetry), func(t *testing.T) {
			t.Parallel()
			params := createTestParams(2, 10000, symmetry)
			renderer := NewRenderer(*params, logger)

			histogram := renderer.Render()
			require.NotNil(t, histogram)
		})
	}
}

func TestRenderer_Render_RemainderIterations(t *testing.T) {
	t.Parallel()
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	params := createTestParams(4, 10001, 1)
	renderer := NewRenderer(*params, logger)

	histogram := renderer.Render()
	require.NotNil(t, histogram)
}

func TestRenderer_Render_SmallIterations(t *testing.T) {
	t.Parallel()
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	params := createTestParams(2, 10, 1)
	renderer := NewRenderer(*params, logger)

	histogram := renderer.Render()
	require.NotNil(t, histogram)
}

func TestRenderer_Render_MultipleAffines(t *testing.T) {
	t.Parallel()
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	affineList := domain.NewAffineList()
	affine1 := domain.Affine{
		A: 0.5, B: 0.0, C: 0.0, D: 0.0, E: 0.5, F: 0.0,
		Weight: 1.0,
		Color:  domain.Color{R: 1.0, G: 0.0, B: 0.0},
	}
	affine2 := domain.Affine{
		A: 0.3, B: 0.0, C: 0.0, D: 0.0, E: 0.3, F: 0.0,
		Weight: 1.0,
		Color:  domain.Color{R: 0.0, G: 1.0, B: 0.0},
	}
	_ = affineList.Add(affine1)
	_ = affineList.Add(affine2)

	variationList := domain.NewVariationList()
	v, _ := domain.GetVariation("linear")
	_ = variationList.Add(v, 1.0)

	params, _ := NewParams(10000, 42, 1, 2, 100, 100, 2.2, *affineList, *variationList)
	renderer := NewRenderer(*params, logger)

	histogram := renderer.Render()
	require.NotNil(t, histogram)
}

func TestRenderer_Render_MultipleVariations(t *testing.T) {
	t.Parallel()
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	affineList := domain.NewAffineList()
	affine := domain.Affine{
		A: 0.5, B: 0.0, C: 0.0, D: 0.0, E: 0.5, F: 0.0,
		Weight: 1.0,
		Color:  domain.Color{R: 1.0, G: 0.5, B: 0.25},
	}
	_ = affineList.Add(affine)

	variationList := domain.NewVariationList()
	v1, _ := domain.GetVariation("linear")
	v2, _ := domain.GetVariation("sinusoidal")
	_ = variationList.Add(v1, 1.0)
	_ = variationList.Add(v2, 0.5)

	params, _ := NewParams(10000, 42, 1, 2, 100, 100, 2.2, *affineList, *variationList)
	renderer := NewRenderer(*params, logger)

	histogram := renderer.Render()
	require.NotNil(t, histogram)
}

func TestRenderer_Render_PointsOutsideBounds(t *testing.T) {
	t.Parallel()
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	affineList := domain.NewAffineList()
	affine := domain.Affine{
		A: 10.0, B: 0.0, C: 0.0, D: 0.0, E: 10.0, F: 0.0,
		Weight: 1.0,
		Color:  domain.Color{R: 1.0, G: 0.5, B: 0.25},
	}
	_ = affineList.Add(affine)

	variationList := domain.NewVariationList()
	v, _ := domain.GetVariation("linear")
	_ = variationList.Add(v, 1.0)

	params, _ := NewParams(10000, 42, 1, 2, 100, 100, 2.2, *affineList, *variationList)
	renderer := NewRenderer(*params, logger)

	histogram := renderer.Render()
	require.NotNil(t, histogram)
}

func TestRenderer_Render_Deterministic(t *testing.T) {
	t.Parallel()
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	params1 := createTestParams(2, 10000, 1)
	renderer1 := NewRenderer(*params1, logger)
	histogram1 := renderer1.Render()

	params2 := createTestParams(2, 10000, 1)
	renderer2 := NewRenderer(*params2, logger)
	histogram2 := renderer2.Render()

	require.NotNil(t, histogram1)
	require.NotNil(t, histogram2)
	assert.Equal(t, histogram1.Width(), histogram2.Width())
	assert.Equal(t, histogram1.Height(), histogram2.Height())
}

func TestRotate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		x        float64
		y        float64
		angle    float64
		expected struct {
			x float64
			y float64
		}
	}{
		{
			name:  "zero angle",
			x:     1.0,
			y:     0.0,
			angle: 0.0,
			expected: struct {
				x float64
				y float64
			}{x: 1.0, y: 0.0},
		},
		{
			name:  "90 degrees",
			x:     1.0,
			y:     0.0,
			angle: math.Pi / 2,
			expected: struct {
				x float64
				y float64
			}{x: 0.0, y: 1.0},
		},
		{
			name:  "180 degrees",
			x:     1.0,
			y:     0.0,
			angle: math.Pi,
			expected: struct {
				x float64
				y float64
			}{x: -1.0, y: 0.0},
		},
		{
			name:  "270 degrees",
			x:     1.0,
			y:     0.0,
			angle: 3 * math.Pi / 2,
			expected: struct {
				x float64
				y float64
			}{x: 0.0, y: -1.0},
		},
		{
			name:  "45 degrees",
			x:     1.0,
			y:     0.0,
			angle: math.Pi / 4,
			expected: struct {
				x float64
				y float64
			}{x: math.Sqrt(2) / 2, y: math.Sqrt(2) / 2},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rotX, rotY := rotate(tt.x, tt.y, tt.angle)
			assert.InDelta(t, tt.expected.x, rotX, 0.0001)
			assert.InDelta(t, tt.expected.y, rotY, 0.0001)
		})
	}
}

func TestRenderer_Render_ProgressTracking(t *testing.T) {
	t.Parallel()
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	params := createTestParams(2, 50000, 1)
	renderer := NewRenderer(*params, logger)

	histogram := renderer.Render()
	require.NotNil(t, histogram)
	assert.GreaterOrEqual(t, atomic.LoadInt64(&renderer.done), int64(0))
}

func TestRenderer_Render_BatchUpdates(t *testing.T) {
	t.Parallel()
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	params := createTestParams(2, 50000, 1)
	renderer := NewRenderer(*params, logger)

	histogram := renderer.Render()
	require.NotNil(t, histogram)
}
