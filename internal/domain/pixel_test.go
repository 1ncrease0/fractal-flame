package domain

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPixel_UpdatePixel(t *testing.T) {
	t.Parallel()
	t.Run("initial update", func(t *testing.T) {
		t.Parallel()
		pixel := Pixel{
			Color: Color{R: 0, G: 0, B: 0},
			Alpha: 0,
		}
		color := Color{R: 1.0, G: 0.5, B: 0.25}
		pixel.UpdatePixel(color)

		assert.InDelta(t, 0.5, pixel.Color.R, 1e-10)
		assert.InDelta(t, 0.25, pixel.Color.G, 1e-10)
		assert.InDelta(t, 0.125, pixel.Color.B, 1e-10)
		assert.Equal(t, int64(10), pixel.Alpha)
	})

	t.Run("multiple updates", func(t *testing.T) {
		t.Parallel()
		pixel := Pixel{
			Color: Color{R: 0, G: 0, B: 0},
			Alpha: 0,
		}
		color1 := Color{R: 1.0, G: 0.0, B: 0.0}
		color2 := Color{R: 0.0, G: 1.0, B: 0.0}

		pixel.UpdatePixel(color1)
		pixel.UpdatePixel(color2)

		assert.InDelta(t, 0.25, pixel.Color.R, 1e-10)
		assert.InDelta(t, 0.5, pixel.Color.G, 1e-10)
		assert.Equal(t, int64(20), pixel.Alpha)
	})

	t.Run("concurrent updates", func(t *testing.T) {
		t.Parallel()
		pixel := Pixel{
			Color: Color{R: 0, G: 0, B: 0},
			Alpha: 0,
		}

		var wg sync.WaitGroup
		numGoroutines := 10
		updatesPerGoroutine := 100

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				color := Color{R: 1.0, G: 1.0, B: 1.0}
				for j := 0; j < updatesPerGoroutine; j++ {
					pixel.UpdatePixel(color)
				}
			}()
		}

		wg.Wait()

		expectedAlpha := int64(numGoroutines * updatesPerGoroutine * 10)
		assert.Equal(t, expectedAlpha, pixel.Alpha)
		assert.GreaterOrEqual(t, pixel.Color.R, 0.0)
		assert.LessOrEqual(t, pixel.Color.R, 1.0)
	})
}
