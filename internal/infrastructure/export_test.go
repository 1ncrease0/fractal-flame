package infrastructure

import (
	"os"
	"testing"

	"fractal-flame/internal/application"
	"fractal-flame/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExportParamsToConfig(t *testing.T) {
	t.Parallel()
	t.Run("full params", func(t *testing.T) {
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
			Color:  domain.Color{R: 1.0, G: 0.5, B: 0.25},
		}
		_ = affineList.Add(affine)

		variationList := domain.NewVariationList()
		v, _ := domain.GetVariation("linear")
		_ = variationList.Add(v, 1.0)

		params, err := application.NewParams(1000000, 42, 2, 4, 1920, 1080, 2.2, *affineList, *variationList)
		require.NoError(t, err)

		cfg, err := ExportParamsToConfig(params)
		require.NoError(t, err)
		require.NotNil(t, cfg)

		assert.Equal(t, 1920, cfg.Size.Width)
		assert.Equal(t, 1080, cfg.Size.Height)
		assert.NotNil(t, cfg.Seed)
		assert.Equal(t, int64(42), *cfg.Seed)
		assert.NotNil(t, cfg.IterationCount)
		assert.Equal(t, int64(1000000), *cfg.IterationCount)
		assert.NotNil(t, cfg.Threads)
		assert.Equal(t, 4, *cfg.Threads)
		assert.NotNil(t, cfg.SymmetryLevel)
		assert.Equal(t, 2, *cfg.SymmetryLevel)
		assert.NotNil(t, cfg.Gamma)
		assert.InDelta(t, 2.2, *cfg.Gamma, 1e-10)

		require.Len(t, cfg.AffineParams, 1)
		assert.Equal(t, 1.0, cfg.AffineParams[0].A)
		assert.NotNil(t, cfg.AffineParams[0].Color)
		assert.Equal(t, 1.0, cfg.AffineParams[0].Color.R)
		assert.Equal(t, 0.5, cfg.AffineParams[0].Color.G)
		assert.Equal(t, 0.25, cfg.AffineParams[0].Color.B)

		require.Len(t, cfg.Functions, 1)
		assert.Equal(t, "linear", cfg.Functions[0].Name)
		assert.Equal(t, 1.0, cfg.Functions[0].Weight)
	})

	t.Run("gamma 1.0", func(t *testing.T) {
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
			Color:  domain.Color{R: 1.0, G: 0.5, B: 0.25},
		}
		_ = affineList.Add(affine)

		variationList := domain.NewVariationList()
		v, _ := domain.GetVariation("linear")
		_ = variationList.Add(v, 1.0)

		params, err := application.NewParams(1000000, 42, 2, 4, 1920, 1080, 1.0, *affineList, *variationList)
		require.NoError(t, err)

		cfg, err := ExportParamsToConfig(params)
		require.NoError(t, err)
		require.NotNil(t, cfg)

		assert.NotNil(t, cfg.GammaCorrection)
		assert.False(t, *cfg.GammaCorrection)
	})

	t.Run("multiple affines and variations", func(t *testing.T) {
		t.Parallel()
		affineList := domain.NewAffineList()
		affine1 := domain.Affine{
			A:      1.0,
			B:      0.0,
			C:      0.0,
			D:      0.0,
			E:      1.0,
			F:      0.0,
			Weight: 1.0,
			Color:  domain.Color{R: 1.0, G: 0.0, B: 0.0},
		}
		affine2 := domain.Affine{
			A:      2.0,
			B:      0.0,
			C:      0.0,
			D:      0.0,
			E:      2.0,
			F:      0.0,
			Weight: 2.0,
			Color:  domain.Color{R: 0.0, G: 1.0, B: 0.0},
		}
		_ = affineList.Add(affine1)
		_ = affineList.Add(affine2)

		variationList := domain.NewVariationList()
		v1, _ := domain.GetVariation("linear")
		v2, _ := domain.GetVariation("sinusoidal")
		_ = variationList.Add(v1, 1.0)
		_ = variationList.Add(v2, 2.0)

		params, err := application.NewParams(1000000, 42, 2, 4, 1920, 1080, 2.2, *affineList, *variationList)
		require.NoError(t, err)

		cfg, err := ExportParamsToConfig(params)
		require.NoError(t, err)
		require.NotNil(t, cfg)

		assert.Len(t, cfg.AffineParams, 2)
		assert.Len(t, cfg.Functions, 2)
	})
}

func TestSaveConfig(t *testing.T) {
	t.Parallel()
	t.Run("save and load", func(t *testing.T) {
		t.Parallel()
		cfg := &Config{}
		cfg.Size.Width = 1920
		cfg.Size.Height = 1080
		seed := int64(42)
		cfg.Seed = &seed
		cfg.Functions = []struct {
			Name   string  `json:"name"`
			Weight float64 `json:"weight"`
		}{
			{Name: "linear", Weight: 1.0},
		}
		cfg.AffineParams = []struct {
			A      float64  `json:"a"`
			B      float64  `json:"b"`
			C      float64  `json:"c"`
			D      float64  `json:"d"`
			E      float64  `json:"e"`
			F      float64  `json:"f"`
			Weight *float64 `json:"weight,omitempty"`
			Color  *struct {
				R float64 `json:"r"`
				G float64 `json:"g"`
				B float64 `json:"b"`
			} `json:"color,omitempty"`
		}{
			{A: 1.0, B: 0.0, C: 0.0, D: 0.0, E: 1.0, F: 0.0},
		}

		tmpFile, err := os.CreateTemp("", "test_export_*.json")
		require.NoError(t, err)
		defer func(name string) {
			_ = os.Remove(name)
		}(tmpFile.Name())
		require.NoError(t, tmpFile.Close())

		err = SaveConfig(cfg, tmpFile.Name())
		require.NoError(t, err)

		loadedCfg, err := LoadConfig(tmpFile.Name())
		require.NoError(t, err)
		require.NotNil(t, loadedCfg)

		assert.Equal(t, cfg.Size.Width, loadedCfg.Size.Width)
		assert.Equal(t, cfg.Size.Height, loadedCfg.Size.Height)
		assert.Equal(t, *cfg.Seed, *loadedCfg.Seed)
		assert.Len(t, loadedCfg.Functions, 1)
		assert.Len(t, loadedCfg.AffineParams, 1)
	})

	t.Run("invalid path", func(t *testing.T) {
		t.Parallel()
		cfg := &Config{}
		err := SaveConfig(cfg, "/invalid/path/config.json")
		_ = err
	})
}

func TestExportRoundTrip(t *testing.T) {
	t.Parallel()
	t.Run("export and import", func(t *testing.T) {
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
			Color:  domain.Color{R: 1.0, G: 0.5, B: 0.25},
		}
		_ = affineList.Add(affine)

		variationList := domain.NewVariationList()
		v, _ := domain.GetVariation("linear")
		_ = variationList.Add(v, 1.0)

		originalParams, err := application.NewParams(1000000, 42, 2, 4, 1920, 1080, 2.2, *affineList, *variationList)
		require.NoError(t, err)

		cfg, err := ExportParamsToConfig(originalParams)
		require.NoError(t, err)

		tmpFile, err := os.CreateTemp("", "test_roundtrip_*.json")
		require.NoError(t, err)
		defer func(name string) {
			_ = os.Remove(name)
		}(tmpFile.Name())
		require.NoError(t, tmpFile.Close())

		err = SaveConfig(cfg, tmpFile.Name())
		require.NoError(t, err)

		loadedCfg, err := LoadConfig(tmpFile.Name())
		require.NoError(t, err)

		loadedParams, _, err := loadedCfg.BuildParams()
		require.NoError(t, err)

		assert.Equal(t, originalParams.Width(), loadedParams.Width())
		assert.Equal(t, originalParams.Height(), loadedParams.Height())
		assert.Equal(t, originalParams.Seed(), loadedParams.Seed())
		assert.Equal(t, originalParams.Iterations(), loadedParams.Iterations())
		assert.Equal(t, originalParams.Threads(), loadedParams.Threads())
		assert.Equal(t, originalParams.Symmetry(), loadedParams.Symmetry())
		assert.InDelta(t, originalParams.Gamma(), loadedParams.Gamma(), 1e-10)
	})
}
