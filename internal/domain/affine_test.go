package domain

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAffine_Apply(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		affine Affine
		x      float64
		y      float64
		expX   float64
		expY   float64
	}{
		{
			name:   "identity",
			affine: Affine{A: 1, B: 0, C: 0, D: 0, E: 1, F: 0},
			x:      1.0,
			y:      2.0,
			expX:   1.0,
			expY:   2.0,
		},
		{
			name:   "translation",
			affine: Affine{A: 1, B: 0, C: 5, D: 0, E: 1, F: 10},
			x:      1.0,
			y:      2.0,
			expX:   6.0,
			expY:   12.0,
		},
		{
			name:   "rotation and scale",
			affine: Affine{A: 2, B: -1, C: 0, D: 1, E: 2, F: 0},
			x:      1.0,
			y:      1.0,
			expX:   1.0,
			expY:   3.0,
		},
		{
			name:   "zero input",
			affine: Affine{A: 1, B: 2, C: 3, D: 4, E: 5, F: 6},
			x:      0.0,
			y:      0.0,
			expX:   3.0,
			expY:   6.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			x, y := tt.affine.Apply(tt.x, tt.y)
			assert.InDelta(t, tt.expX, x, 1e-10, "x coordinate")
			assert.InDelta(t, tt.expY, y, 1e-10, "y coordinate")
		})
	}
}

func TestAffineList_Add(t *testing.T) {
	t.Parallel()
	t.Run("add valid affine", func(t *testing.T) {
		t.Parallel()
		list := NewAffineList()
		affine := Affine{
			A:      1.0,
			B:      0.0,
			C:      0.0,
			D:      0.0,
			E:      1.0,
			F:      0.0,
			Weight: 1.0,
		}
		err := list.Add(affine)
		assert.NoError(t, err)
		assert.Equal(t, 1, list.Len())
	})

	t.Run("add multiple affines", func(t *testing.T) {
		t.Parallel()
		list := NewAffineList()
		affine1 := Affine{A: 1, B: 0, C: 0, D: 0, E: 1, F: 0, Weight: 1.0}
		affine2 := Affine{A: 2, B: 0, C: 0, D: 0, E: 2, F: 0, Weight: 2.0}
		err1 := list.Add(affine1)
		err2 := list.Add(affine2)
		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.Equal(t, 2, list.Len())
	})

	t.Run("add negative weight", func(t *testing.T) {
		t.Parallel()
		list := NewAffineList()
		affine := Affine{
			A:      1.0,
			B:      0.0,
			C:      0.0,
			D:      0.0,
			E:      1.0,
			F:      0.0,
			Weight: -1.0,
		}
		err := list.Add(affine)
		assert.Error(t, err)
		assert.Equal(t, ErrNegativeAffineWeight, err)
		assert.Equal(t, 0, list.Len())
	})
}

func TestAffineList_Random(t *testing.T) {
	t.Parallel()
	t.Run("empty list", func(t *testing.T) {
		t.Parallel()
		list := NewAffineList()
		rnd := rand.New(rand.NewSource(42))
		affine := list.Random(rnd)
		assert.Nil(t, affine)
	})

	t.Run("single affine", func(t *testing.T) {
		t.Parallel()
		list := NewAffineList()
		affine1 := Affine{A: 1, B: 0, C: 0, D: 0, E: 1, F: 0, Weight: 1.0}
		_ = list.Add(affine1)
		rnd := rand.New(rand.NewSource(42))
		affine := list.Random(rnd)
		require.NotNil(t, affine)
		assert.Equal(t, affine1.A, affine.A)
		assert.Equal(t, affine1.B, affine.B)
		assert.Equal(t, affine1.C, affine.C)
		assert.Equal(t, affine1.D, affine.D)
		assert.Equal(t, affine1.E, affine.E)
		assert.Equal(t, affine1.F, affine.F)
	})

	t.Run("multiple affines with weights", func(t *testing.T) {
		t.Parallel()
		list := NewAffineList()
		affine1 := Affine{A: 1, B: 0, C: 0, D: 0, E: 1, F: 0, Weight: 1.0}
		affine2 := Affine{A: 2, B: 0, C: 0, D: 0, E: 2, F: 0, Weight: 9.0}
		_ = list.Add(affine1)
		_ = list.Add(affine2)
		rnd := rand.New(rand.NewSource(42))

		results := make(map[*Affine]int)
		for i := 0; i < 100; i++ {
			affine := list.Random(rnd)
			require.NotNil(t, affine)
			results[affine]++
		}

		assert.Greater(t, len(results), 0)
	})

	t.Run("weighted selection", func(t *testing.T) {
		t.Parallel()
		list := NewAffineList()
		affine1 := Affine{A: 1, B: 0, C: 0, D: 0, E: 1, F: 0, Weight: 0.1}
		affine2 := Affine{A: 2, B: 0, C: 0, D: 0, E: 2, F: 0, Weight: 0.9}
		_ = list.Add(affine1)
		_ = list.Add(affine2)
		rnd := rand.New(rand.NewSource(42))

		affine2Count := 0
		for i := 0; i < 100; i++ {
			affine := list.Random(rnd)
			if affine.A == 2 {
				affine2Count++
			}
		}
		assert.Greater(t, affine2Count, 50)
	})
}

func TestAffineList_Affine(t *testing.T) {
	t.Parallel()
	list := NewAffineList()
	affine1 := Affine{A: 1, B: 0, C: 0, D: 0, E: 1, F: 0, Weight: 1.0}
	affine2 := Affine{A: 2, B: 0, C: 0, D: 0, E: 2, F: 0, Weight: 2.0}
	_ = list.Add(affine1)
	_ = list.Add(affine2)

	affines := list.Affine()
	require.Equal(t, 2, len(affines))
	assert.Equal(t, affine1.A, affines[0].A)
	assert.Equal(t, affine2.A, affines[1].A)
}

func TestAffineList_Len(t *testing.T) {
	t.Parallel()
	list := NewAffineList()
	assert.Equal(t, 0, list.Len())

	affine1 := Affine{A: 1, B: 0, C: 0, D: 0, E: 1, F: 0, Weight: 1.0}
	_ = list.Add(affine1)
	assert.Equal(t, 1, list.Len())

	affine2 := Affine{A: 2, B: 0, C: 0, D: 0, E: 2, F: 0, Weight: 2.0}
	_ = list.Add(affine2)
	assert.Equal(t, 2, list.Len())
}
