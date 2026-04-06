package infrastructure

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	t.Parallel()
	t.Run("valid config file", func(t *testing.T) {
		t.Parallel()
		configData := `{
			"size": {
				"width": 1920,
				"height": 1080
			},
			"seed": 42,
			"iteration_count": 1000000,
			"threads": 4,
			"symmetry_level": 2,
			"gamma": 2.2,
			"functions": [
				{"name": "linear", "weight": 1.0}
			],
			"affine_params": [
				{"a": 1.0, "b": 0.0, "c": 0.0, "d": 0.0, "e": 1.0, "f": 0.0}
			]
		}`

		tmpFile, err := os.CreateTemp("", "test_config_*.json")
		require.NoError(t, err)
		defer func(name string) {
			_ = os.Remove(name)
		}(tmpFile.Name())

		_, err = tmpFile.WriteString(configData)
		require.NoError(t, err)
		require.NoError(t, tmpFile.Close())

		cfg, err := LoadConfig(tmpFile.Name())
		require.NoError(t, err)
		require.NotNil(t, cfg)
		assert.Equal(t, 1920, cfg.Size.Width)
		assert.Equal(t, 1080, cfg.Size.Height)
		assert.NotNil(t, cfg.Seed)
		assert.Equal(t, int64(42), *cfg.Seed)
	})

	t.Run("non-existent file", func(t *testing.T) {
		t.Parallel()
		cfg, err := LoadConfig("non_existent_file.json")
		assert.Error(t, err)
		assert.Nil(t, cfg)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		t.Parallel()
		tmpFile, err := os.CreateTemp("", "test_config_*.json")
		require.NoError(t, err)
		defer func(name string) {
			_ = os.Remove(name)
		}(tmpFile.Name())

		_, err = tmpFile.WriteString("invalid json")
		require.NoError(t, err)
		require.NoError(t, tmpFile.Close())

		cfg, err := LoadConfig(tmpFile.Name())
		assert.Error(t, err)
		assert.Nil(t, cfg)
	})
}

func TestConfig_BuildParams(t *testing.T) {
	t.Parallel()
	t.Run("full config", func(t *testing.T) {
		t.Parallel()
		cfg := &Config{}
		cfg.Size.Width = 1920
		cfg.Size.Height = 1080
		seed := int64(42)
		cfg.Seed = &seed
		iters := int64(1000000)
		cfg.IterationCount = &iters
		threads := 4
		cfg.Threads = &threads
		symmetry := 2
		cfg.SymmetryLevel = &symmetry
		gamma := 2.2
		cfg.Gamma = &gamma
		outputPath := "test.png"
		cfg.OutputPath = &outputPath

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

		params, path, err := cfg.BuildParams()
		require.NoError(t, err)
		require.NotNil(t, params)
		assert.Equal(t, "test.png", path)
		assert.Equal(t, 1920, params.Width())
		assert.Equal(t, 1080, params.Height())
		assert.Equal(t, int64(42), params.Seed())
		assert.Equal(t, int64(1000000), params.Iterations())
		assert.Equal(t, 4, params.Threads())
		assert.Equal(t, 2, params.Symmetry())
		assert.InDelta(t, 2.2, params.Gamma(), 1e-10)
	})

	t.Run("default values", func(t *testing.T) {
		t.Parallel()
		cfg := &Config{}
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

		params, path, err := cfg.BuildParams()
		require.NoError(t, err)
		require.NotNil(t, params)
		assert.Equal(t, DefaultOutputPath, path)
		assert.Equal(t, DefaultWidth, params.Width())
		assert.Equal(t, DefaultHeight, params.Height())
		assert.Equal(t, DefaultSeed, params.Seed())
		assert.Equal(t, DefaultIterationCount, params.Iterations())
		assert.Equal(t, DefaultThreads, params.Threads())
		assert.Equal(t, DefaultSymmetryLevel, params.Symmetry())
		assert.InDelta(t, DefaultGamma, params.Gamma(), 1e-10)
	})

	t.Run("empty affine params", func(t *testing.T) {
		t.Parallel()
		cfg := &Config{}
		cfg.Functions = []struct {
			Name   string  `json:"name"`
			Weight float64 `json:"weight"`
		}{
			{Name: "linear", Weight: 1.0},
		}

		params, _, err := cfg.BuildParams()
		assert.Error(t, err)
		assert.Equal(t, ErrEmptyAffineParams, err)
		assert.Nil(t, params)
	})

	t.Run("empty variation params", func(t *testing.T) {
		t.Parallel()
		cfg := &Config{}
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

		params, _, err := cfg.BuildParams()
		assert.Error(t, err)
		assert.Equal(t, ErrEmptyVariationParams, err)
		assert.Nil(t, params)
	})
}

func TestConfig_BuildAffineList(t *testing.T) {
	t.Parallel()
	t.Run("with colors", func(t *testing.T) {
		t.Parallel()
		cfg := &Config{}
		seed := int64(42)
		cfg.Seed = &seed

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
			{
				A: 1.0, B: 0.0, C: 0.0, D: 0.0, E: 1.0, F: 0.0,
				Color: &struct {
					R float64 `json:"r"`
					G float64 `json:"g"`
					B float64 `json:"b"`
				}{R: 1.0, G: 0.5, B: 0.25},
			},
		}

		affineList, err := cfg.BuildAffineList()
		require.NoError(t, err)
		require.NotNil(t, affineList)
		assert.Equal(t, 1, affineList.Len())

		affines := affineList.Affine()
		assert.Equal(t, 1.0, affines[0].Color.R)
		assert.Equal(t, 0.5, affines[0].Color.G)
		assert.Equal(t, 0.25, affines[0].Color.B)
	})

	t.Run("without colors", func(t *testing.T) {
		t.Parallel()
		cfg := &Config{}
		seed := int64(42)
		cfg.Seed = &seed

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

		affineList, err := cfg.BuildAffineList()
		require.NoError(t, err)
		require.NotNil(t, affineList)
		assert.Equal(t, 1, affineList.Len())

		affines := affineList.Affine()
		assert.GreaterOrEqual(t, affines[0].Color.R, 0.0)
		assert.LessOrEqual(t, affines[0].Color.R, 1.0)
	})

	t.Run("with weights", func(t *testing.T) {
		t.Parallel()
		cfg := &Config{}
		seed := int64(42)
		cfg.Seed = &seed
		weight := 2.0

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
			{
				A: 1.0, B: 0.0, C: 0.0, D: 0.0, E: 1.0, F: 0.0,
				Weight: &weight,
			},
		}

		affineList, err := cfg.BuildAffineList()
		require.NoError(t, err)
		affines := affineList.Affine()
		assert.Equal(t, 2.0, affines[0].Weight)
	})
}

func TestConfig_BuildVariationList(t *testing.T) {
	t.Parallel()
	t.Run("valid variations", func(t *testing.T) {
		t.Parallel()
		cfg := &Config{}
		cfg.Functions = []struct {
			Name   string  `json:"name"`
			Weight float64 `json:"weight"`
		}{
			{Name: "linear", Weight: 1.0},
			{Name: "sinusoidal", Weight: 2.0},
		}

		variationList, err := cfg.BuildVariationList()
		require.NoError(t, err)
		require.NotNil(t, variationList)
		assert.Equal(t, 2, variationList.Len())
		assert.InDelta(t, 3.0, variationList.TotalWeight(), 1e-10)
	})

	t.Run("invalid variation name", func(t *testing.T) {
		t.Parallel()
		cfg := &Config{}
		cfg.Functions = []struct {
			Name   string  `json:"name"`
			Weight float64 `json:"weight"`
		}{
			{Name: "invalid_variation", Weight: 1.0},
		}

		variationList, err := cfg.BuildVariationList()
		assert.Error(t, err)
		assert.Nil(t, variationList)
	})

	t.Run("negative weight", func(t *testing.T) {
		t.Parallel()
		cfg := &Config{}
		cfg.Functions = []struct {
			Name   string  `json:"name"`
			Weight float64 `json:"weight"`
		}{
			{Name: "linear", Weight: -1.0},
		}

		variationList, err := cfg.BuildVariationList()
		assert.Error(t, err)
		assert.Nil(t, variationList)
	})
}

func TestConfig_Seed(t *testing.T) {
	t.Parallel()
	t.Run("with seed", func(t *testing.T) {
		t.Parallel()
		cfg := &Config{}
		seed := int64(42)
		cfg.Seed = &seed
		assert.Equal(t, int64(42), cfg.seed())
	})

	t.Run("without seed", func(t *testing.T) {
		t.Parallel()
		cfg := &Config{}
		assert.Equal(t, DefaultSeed, cfg.seed())
	})
}

func TestConfig_JSON(t *testing.T) {
	t.Parallel()
	t.Run("marshal and unmarshal", func(t *testing.T) {
		t.Parallel()
		cfg := &Config{}
		cfg.Size.Width = 1920
		cfg.Size.Height = 1080
		seed := int64(42)
		cfg.Seed = &seed

		data, err := json.Marshal(cfg)
		require.NoError(t, err)

		var cfg2 Config
		err = json.Unmarshal(data, &cfg2)
		require.NoError(t, err)

		assert.Equal(t, cfg.Size.Width, cfg2.Size.Width)
		assert.Equal(t, cfg.Size.Height, cfg2.Size.Height)
		assert.Equal(t, *cfg.Seed, *cfg2.Seed)
	})
}
