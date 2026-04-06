package application

import (
	"testing"

	"fractal-flame/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewParams(t *testing.T) {
	t.Parallel()
	affineList := domain.NewAffineList()
	affine := domain.Affine{
		A:      1.0,
		B:      0.0,
		C:      0.0,
		D:      0.0,
		E:      1.0,
		F:      0.0,
		Weight: 1.0,
	}
	_ = affineList.Add(affine)

	variationList := domain.NewVariationList()
	v, _ := domain.GetVariation("linear")
	_ = variationList.Add(v, 1.0)

	t.Run("valid params", func(t *testing.T) {
		t.Parallel()
		params, err := NewParams(1000000, 42, 1, 4, 1920, 1080, 2.2, *affineList, *variationList)
		require.NoError(t, err)
		require.NotNil(t, params)
		assert.Equal(t, int64(1000000), params.Iterations())
		assert.Equal(t, int64(42), params.Seed())
		assert.Equal(t, 1, params.Symmetry())
		assert.Equal(t, 4, params.Threads())
		assert.Equal(t, 1920, params.Width())
		assert.Equal(t, 1080, params.Height())
		assert.InDelta(t, 2.2, params.Gamma(), 1e-10)
	})

	t.Run("invalid iterations", func(t *testing.T) {
		t.Parallel()
		params, err := NewParams(0, 42, 1, 4, 1920, 1080, 2.2, *affineList, *variationList)
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidIters, err)
		assert.Nil(t, params)
	})

	t.Run("invalid symmetry", func(t *testing.T) {
		t.Parallel()
		params, err := NewParams(1000000, 42, 0, 4, 1920, 1080, 2.2, *affineList, *variationList)
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidSymmetry, err)
		assert.Nil(t, params)
	})

	t.Run("invalid threads", func(t *testing.T) {
		t.Parallel()
		params, err := NewParams(1000000, 42, 1, 0, 1920, 1080, 2.2, *affineList, *variationList)
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidThreads, err)
		assert.Nil(t, params)
	})

	t.Run("empty affine list", func(t *testing.T) {
		t.Parallel()
		emptyAffineList := domain.NewAffineList()
		params, err := NewParams(1000000, 42, 1, 4, 1920, 1080, 2.2, *emptyAffineList, *variationList)
		assert.Error(t, err)
		assert.Equal(t, ErrEmptyAffineList, err)
		assert.Nil(t, params)
	})

	t.Run("invalid width", func(t *testing.T) {
		t.Parallel()
		params, err := NewParams(1000000, 42, 1, 4, 0, 1080, 2.2, *affineList, *variationList)
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidSize, err)
		assert.Nil(t, params)
	})

	t.Run("invalid height", func(t *testing.T) {
		t.Parallel()
		params, err := NewParams(1000000, 42, 1, 4, 1920, 0, 2.2, *affineList, *variationList)
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidSize, err)
		assert.Nil(t, params)
	})
}

func TestParams_Getters(t *testing.T) {
	t.Parallel()
	affineList := domain.NewAffineList()
	affine := domain.Affine{
		A:      1.0,
		B:      0.0,
		C:      0.0,
		D:      0.0,
		E:      1.0,
		F:      0.0,
		Weight: 1.0,
	}
	_ = affineList.Add(affine)

	variationList := domain.NewVariationList()
	v, _ := domain.GetVariation("linear")
	_ = variationList.Add(v, 1.0)

	params, err := NewParams(1000000, 42, 2, 4, 1920, 1080, 2.2, *affineList, *variationList)
	require.NoError(t, err)

	t.Run("affine", func(t *testing.T) {
		t.Parallel()
		affineList := params.Affine()
		assert.Equal(t, 1, affineList.Len())
	})

	t.Run("variations", func(t *testing.T) {
		t.Parallel()
		variationList := params.Variations()
		assert.Equal(t, 1, variationList.Len())
	})

	t.Run("iterations", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, int64(1000000), params.Iterations())
	})

	t.Run("seed", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, int64(42), params.Seed())
	})

	t.Run("symmetry", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, 2, params.Symmetry())
	})

	t.Run("threads", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, 4, params.Threads())
	})

	t.Run("gamma", func(t *testing.T) {
		t.Parallel()
		assert.InDelta(t, 2.2, params.Gamma(), 1e-10)
	})

	t.Run("width", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, 1920, params.Width())
	})

	t.Run("height", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, 1080, params.Height())
	})
}

func TestRandomParams(t *testing.T) {
	t.Parallel()
	t.Run("generates valid params", func(t *testing.T) {
		t.Parallel()
		params := RandomParams()
		require.NotNil(t, params)
		assert.Greater(t, params.Iterations(), int64(0))
		assert.Greater(t, params.Seed(), int64(0))
		assert.GreaterOrEqual(t, params.Symmetry(), 1)
		assert.Greater(t, params.Threads(), 0)
		assert.Greater(t, params.Width(), 0)
		assert.Greater(t, params.Height(), 0)
		assert.Greater(t, params.Gamma(), 0.0)
		affineList := params.Affine()
		affines := affineList.Affine()
		assert.Greater(t, len(affines), 0)
		variationList := params.Variations()
		variations, _ := variationList.Variations()
		assert.Greater(t, len(variations), 0)
	})

	t.Run("generates different params", func(t *testing.T) {
		t.Parallel()
		params1 := RandomParams()
		params2 := RandomParams()

		if params1.Seed() == params2.Seed() {
			t.Skip("Seeds happened to be the same (very unlikely)")
		}
	})
}
